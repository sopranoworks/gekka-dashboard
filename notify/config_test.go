/*
 * config_test.go
 * This file is part of the gekka-metrics project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package notify

import (
	"testing"
	"time"
)

func TestParseConfig_BasicRule(t *testing.T) {
	hocon := `
gekka.notifications {
  rules {
    test-rule {
      events = ["node.unreachable", "node.downed"]
      roles = ["cart"]
      channels = ["slack"]
      throttle = "5m"
    }
  }
  channels {
    slack {
      webhook-url = "https://hooks.slack.com/test"
    }
  }
}
`
	cfg, err := ParseNotifyConfig([]byte(hocon))
	if err != nil {
		t.Fatalf("ParseNotifyConfig: %v", err)
	}
	if len(cfg.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(cfg.Rules))
	}
	r := cfg.Rules[0]
	if r.Name != "test-rule" {
		t.Errorf("name = %q", r.Name)
	}
	if len(r.Events) != 2 {
		t.Errorf("events count = %d", len(r.Events))
	}
	if r.Events[0] != EventNodeUnreachable {
		t.Errorf("events[0] = %q", r.Events[0])
	}
	if len(r.Roles) != 1 || r.Roles[0] != "cart" {
		t.Errorf("roles = %v", r.Roles)
	}
	if r.Throttle != 5*time.Minute {
		t.Errorf("throttle = %v", r.Throttle)
	}
	if cfg.Slack == nil {
		t.Fatal("slack config is nil")
	}
	if cfg.Slack.WebhookURL != "https://hooks.slack.com/test" {
		t.Errorf("slack webhook = %q", cfg.Slack.WebhookURL)
	}
}

func TestParseConfig_EmailChannel(t *testing.T) {
	hocon := `
gekka.notifications {
  rules {}
  channels {
    email {
      smtp-host = "smtp.example.com"
      smtp-port = 587
      from = "alerts@example.com"
      to = ["ops@example.com", "team@example.com"]
    }
  }
}
`
	cfg, err := ParseNotifyConfig([]byte(hocon))
	if err != nil {
		t.Fatalf("ParseNotifyConfig: %v", err)
	}
	if cfg.Email == nil {
		t.Fatal("email config is nil")
	}
	if cfg.Email.Host != "smtp.example.com" {
		t.Errorf("host = %q", cfg.Email.Host)
	}
	if cfg.Email.Port != 587 {
		t.Errorf("port = %d", cfg.Email.Port)
	}
	if len(cfg.Email.To) != 2 {
		t.Errorf("to count = %d", len(cfg.Email.To))
	}
}

func TestParseConfig_NoNotificationsSection(t *testing.T) {
	hocon := `pekko.cluster.seed-nodes = []`
	cfg, err := ParseNotifyConfig([]byte(hocon))
	if err != nil {
		t.Fatalf("ParseNotifyConfig: %v", err)
	}
	if len(cfg.Rules) != 0 {
		t.Errorf("expected 0 rules, got %d", len(cfg.Rules))
	}
}
