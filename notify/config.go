/*
 * config.go
 * This file is part of the gekka-metrics project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package notify

import (
	"os"
	"time"

	config "github.com/sopranoworks/gekka-config"
)

// NotifyConfig holds the parsed notification configuration.
type NotifyConfig struct {
	Rules []*Rule
	Email *EmailConfig // nil if not configured
	Slack *SlackConfig // nil if not configured
}

// SlackConfig holds Slack webhook configuration.
type SlackConfig struct {
	WebhookURL string
}

// ParseNotifyConfigFromFile reads a HOCON file and parses notification config.
func ParseNotifyConfigFromFile(path string) (*NotifyConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseNotifyConfig(data)
}

// ParseNotifyConfig parses HOCON bytes and extracts the gekka.notifications section.
// Returns an empty config (no error) if the section is absent.
//
// Expected HOCON structure:
//
//	gekka.notifications {
//	  rules {
//	    rule-name {
//	      events = ["node.unreachable", "node.downed"]
//	      roles = ["cart"]
//	      channels = ["slack"]
//	      throttle = "5m"
//	    }
//	  }
//	  channels {
//	    email { smtp-host = "...", smtp-port = 587, from = "...", to = ["..."] }
//	    slack { webhook-url = "https://..." }
//	  }
//	}
func ParseNotifyConfig(data []byte) (*NotifyConfig, error) {
	c, err := config.ParseString(string(data))
	if err != nil {
		return nil, err
	}

	result := &NotifyConfig{}

	notif, err := c.GetConfig("gekka.notifications")
	if err != nil {
		// Section not found — return empty config
		return result, nil
	}

	// Parse rules (each rule is a named object under "rules")
	rulesConfig, err := notif.GetConfig("rules")
	if err == nil {
		for _, name := range rulesConfig.Keys() {
			rc, err := rulesConfig.GetConfig(name)
			if err != nil {
				continue
			}
			rule := parseRule(name, rc)
			result.Rules = append(result.Rules, rule)
		}
	}

	// Parse channels
	channels, err := notif.GetConfig("channels")
	if err == nil {
		if emailCfg, err := channels.GetConfig("email"); err == nil {
			result.Email = parseEmailConfig(emailCfg)
		}
		if slackCfg, err := channels.GetConfig("slack"); err == nil {
			webhookURL, _ := slackCfg.GetString("webhook-url")
			result.Slack = &SlackConfig{
				WebhookURL: webhookURL,
			}
		}
	}

	return result, nil
}

func parseRule(name string, rc config.Config) *Rule {
	rule := &Rule{
		Name: name,
	}

	// Parse events list
	var eventsHolder struct {
		Events []string `hocon:"events"`
	}
	if err := rc.Unmarshal(&eventsHolder); err == nil {
		for _, e := range eventsHolder.Events {
			rule.Events = append(rule.Events, EventKind(e))
		}
	}

	// Parse roles list
	var rolesHolder struct {
		Roles []string `hocon:"roles"`
	}
	if err := rc.Unmarshal(&rolesHolder); err == nil {
		rule.Roles = rolesHolder.Roles
	}

	// Parse channels list
	var channelsHolder struct {
		Channels []string `hocon:"channels"`
	}
	if err := rc.Unmarshal(&channelsHolder); err == nil {
		rule.Channels = channelsHolder.Channels
	}

	// Parse throttle
	if d, err := rc.GetDuration("throttle"); err == nil {
		rule.Throttle = d
	} else if s, err := rc.GetString("throttle"); err == nil && s != "" {
		if d, err := time.ParseDuration(s); err == nil {
			rule.Throttle = d
		}
	}

	return rule
}

func parseEmailConfig(emailCfg config.Config) *EmailConfig {
	host, _ := emailCfg.GetString("smtp-host")
	port, _ := emailCfg.GetInt("smtp-port")
	from, _ := emailCfg.GetString("from")

	ec := &EmailConfig{
		Host: host,
		Port: port,
		From: from,
	}

	var toHolder struct {
		To []string `hocon:"to"`
	}
	if err := emailCfg.Unmarshal(&toHolder); err == nil {
		ec.To = toHolder.To
	}

	if user, err := emailCfg.GetString("username"); err == nil {
		ec.Username = user
	}
	if pass, err := emailCfg.GetString("password"); err == nil {
		ec.Password = pass
	}

	return ec
}

// BuildChannels creates Channel instances from a NotifyConfig.
func BuildChannels(cfg *NotifyConfig) map[string]Channel {
	channels := make(map[string]Channel)
	if cfg.Email != nil {
		channels["email"] = &EmailChannel{
			Host:     cfg.Email.Host,
			Port:     cfg.Email.Port,
			From:     cfg.Email.From,
			To:       cfg.Email.To,
			Username: cfg.Email.Username,
			Password: cfg.Email.Password,
		}
	}
	if cfg.Slack != nil {
		channels["slack"] = &SlackChannel{
			WebhookURL: cfg.Slack.WebhookURL,
		}
	}
	return channels
}
