/*
 * main.go
 * This file is part of the gekka-dashboard project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	gekka "github.com/sopranoworks/gekka"
	gcluster "github.com/sopranoworks/gekka/cluster"
	config "github.com/sopranoworks/gekka-config"
	gekkaotel "github.com/sopranoworks/gekka-extensions-telemetry-otel"
	"github.com/sopranoworks/gekka-dashboard/notify"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	otelmetric "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func main() {
	flagConfig := flag.String("config", "", "Path to HOCON application.conf (required)")
	flagOtlp := flag.String("otlp", "", "OTLP/HTTP collector endpoint (overrides config)")
	flagListen := flag.String("listen", "", "Dashboard HTTP listen address (overrides config, default :9000)")
	flagHeadless := flag.Bool("headless", false, "Disable UI, run as metrics-only exporter with notifications")
	flagDisableManagement := flag.Bool("disable-management", false,
		"Do not auto-enable the HTTP management API")
	flag.Parse()

	if *flagConfig == "" {
		fmt.Fprintln(os.Stderr, "--config is required")
		os.Exit(1)
	}

	cfg, err := gekka.LoadConfig(*flagConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	cfg.Roles = appendIfMissing(cfg.Roles, "dashboard")

	applyManagementDefaults(&cfg, *flagDisableManagement)

	otlpEndpoint := cfg.Telemetry.OtlpEndpoint
	if *flagOtlp != "" {
		otlpEndpoint = *flagOtlp
	}

	// Set log level
	var level slog.Level
	switch strings.ToUpper(cfg.LogLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))

	// OTel SDK
	ctx := context.Background()
	mp, shutdown, err := initMeterProvider(ctx, otlpEndpoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init meter provider: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdown(shutCtx)
	}()
	otel.SetMeterProvider(mp)

	// Join cluster
	node, err := gekka.NewCluster(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create cluster node: %v\n", err)
		os.Exit(1)
	}
	defer node.Shutdown()

	if err := node.JoinSeeds(); err != nil {
		fmt.Fprintf(os.Stderr, "join cluster: %v\n", err)
		os.Exit(1)
	}

	// Register OTel gauge
	meter := mp.Meter("github.com/sopranoworks/gekka-dashboard")
	cm := node.ClusterManager()

	_, _ = meter.Int64ObservableGauge(
		"gekka.cluster.members",
		otelmetric.WithDescription("Number of cluster members in each status/dc combination"),
		otelmetric.WithUnit("{members}"),
		otelmetric.WithInt64Callback(func(_ context.Context, obs otelmetric.Int64Observer) error {
			cm.Mu.RLock()
			gossip := cm.State
			cm.Mu.RUnlock()

			if gossip == nil {
				return nil
			}

			type groupKey struct{ status, dc string }
			counts := make(map[groupKey]int64)
			for _, m := range gossip.GetMembers() {
				status := statusString(m.GetStatus().String())
				dc := gcluster.DataCenterForMember(gossip, m)
				counts[groupKey{status, dc}]++
			}
			for k, n := range counts {
				obs.Observe(n,
					otelmetric.WithAttributes(
						attribute.String("status", k.status),
						attribute.String("dc", k.dc),
					),
				)
			}
			return nil
		}),
	)

	// Start notification engine
	notifyCfg, notifyErr := notify.ParseNotifyConfigFromFile(*flagConfig)
	if notifyErr != nil {
		slog.Warn("notify: config parse error, notifications disabled", "err", notifyErr)
	}

	if notifyCfg != nil && len(notifyCfg.Rules) > 0 {
		channels := notify.BuildChannels(notifyCfg)
		eng := notify.NewEngine(notifyCfg.Rules, channels)

		sub := cm.SubscribeChannel()
		go eng.Run(ctx)
		go notify.BridgeClusterEvents(ctx, sub, cm, eng)

		slog.Info("notify: engine started", "rules", len(notifyCfg.Rules), "channels", len(channels))
	}

	// Resolve dashboard listen address: flag > config > default
	listenAddr := resolveDashboardListen(*flagConfig, *flagListen)

	// Start server or block
	if *flagHeadless {
		slog.Info("dashboard: running in headless mode (metrics + notifications only)")
		select {}
	}

	hub := NewHub()
	slog.Info("dashboard: starting", "listen", listenAddr)
	if err := startServer(listenAddr, hub); err != nil {
		slog.Error("dashboard: server failed", "err", err)
		os.Exit(1)
	}
}

func statusString(s string) string {
	return strings.TrimPrefix(s, "MemberStatus_")
}

func appendIfMissing(roles []string, role string) []string {
	for _, r := range roles {
		if r == role {
			return roles
		}
	}
	return append(roles, role)
}

// applyManagementDefaults flips the HTTP management API on by default unless
// disable is true. Hostname and Port are set to gekka-dashboard defaults
// (127.0.0.1:8559) whenever they hold their zero value OR the upstream
// gekka.LoadConfig default (8558).
// resolveDashboardListen determines the HTTP listen address using precedence:
//  1. --listen flag (non-empty means explicitly set)
//  2. gekka.dashboard.listen from HOCON config
//  3. Default ":9000"
func resolveDashboardListen(configPath, flagListen string) string {
	const defaultListen = ":9000"

	// 1. CLI flag wins
	if flagListen != "" {
		return flagListen
	}

	// 2. Try HOCON config
	if configPath != "" {
		if addr := readDashboardListenFromConfig(configPath); addr != "" {
			return addr
		}
	}

	// 3. Default
	return defaultListen
}

func readDashboardListenFromConfig(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	c, err := config.ParseString(string(data))
	if err != nil {
		return ""
	}
	dashCfg, err := c.GetConfig("gekka.dashboard")
	if err != nil {
		return ""
	}
	listen, err := dashCfg.GetString("listen")
	if err != nil || listen == "" {
		return ""
	}
	return listen
}

func applyManagementDefaults(cfg *gekka.ClusterConfig, disable bool) {
	if disable {
		return
	}
	cfg.Management.Enabled = true
	if cfg.Management.Hostname == "" || cfg.Management.Hostname == "127.0.0.1" {
		cfg.Management.Hostname = "127.0.0.1"
	}
	if cfg.Management.Port == 0 || cfg.Management.Port == 8558 {
		cfg.Management.Port = 8559
	}
}

func initMeterProvider(ctx context.Context, otlpEndpoint string) (*sdkmetric.MeterProvider, func(context.Context) error, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName("gekka-dashboard")),
		resource.WithFromEnv(),
	)
	if err != nil {
		return nil, nil, err
	}

	var readerOpt sdkmetric.Option
	if otlpEndpoint != "" {
		host, secure := gekkaotel.ParseEndpoint(otlpEndpoint)

		opts := []otlpmetrichttp.Option{
			otlpmetrichttp.WithEndpoint(host),
		}
		if !secure {
			opts = append(opts, otlpmetrichttp.WithInsecure())
		}

		exp, err := otlpmetrichttp.New(ctx, opts...)
		if err != nil {
			return nil, nil, err
		}
		readerOpt = sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp))
	} else {
		readerOpt = sdkmetric.WithReader(sdkmetric.NewManualReader())
	}

	mp := sdkmetric.NewMeterProvider(sdkmetric.WithResource(res), readerOpt)
	return mp, mp.Shutdown, nil
}
