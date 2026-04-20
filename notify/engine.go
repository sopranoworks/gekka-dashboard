/*
 * engine.go
 * This file is part of the gekka-metrics project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package notify

import (
	"context"
	"log/slog"
)

// Engine is the notification dispatch engine. It receives cluster events,
// evaluates them against rules, and dispatches to matching channels.
type Engine struct {
	rules    []*Rule
	channels map[string]Channel
	throttle *ThrottleTracker
	events   chan NotifyEvent
}

// NewEngine creates a notification engine with the given rules and channels.
func NewEngine(rules []*Rule, channels map[string]Channel) *Engine {
	return &Engine{
		rules:    rules,
		channels: channels,
		throttle: NewThrottleTracker(),
		events:   make(chan NotifyEvent, 256),
	}
}

// HandleEvent enqueues a cluster event for evaluation. Non-blocking; drops
// events if the internal buffer is full (logs a warning).
func (e *Engine) HandleEvent(evt NotifyEvent) {
	select {
	case e.events <- evt:
	default:
		slog.Warn("notify: event buffer full, dropping event", "kind", evt.Kind, "address", evt.Address)
	}
}

// Run starts the dispatch loop. Blocks until ctx is cancelled.
func (e *Engine) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case evt := <-e.events:
			e.evaluate(ctx, evt)
		}
	}
}

func (e *Engine) evaluate(ctx context.Context, evt NotifyEvent) {
	for _, rule := range e.rules {
		if !rule.Matches(evt) {
			continue
		}
		if !e.throttle.Allow(rule.Name, rule.Throttle) {
			slog.Debug("notify: throttled", "rule", rule.Name, "kind", evt.Kind)
			continue
		}
		for _, chName := range rule.Channels {
			ch, ok := e.channels[chName]
			if !ok {
				slog.Warn("notify: unknown channel in rule", "channel", chName, "rule", rule.Name)
				continue
			}
			if err := ch.Send(ctx, evt, rule.Name); err != nil {
				slog.Error("notify: dispatch failed", "channel", chName, "rule", rule.Name, "err", err)
			}
		}
	}
}
