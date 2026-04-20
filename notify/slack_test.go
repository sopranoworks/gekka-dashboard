/*
 * slack_test.go
 * This file is part of the gekka-metrics project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package notify

import (
	"strings"
	"testing"
	"time"
)

func TestSlackChannel_FormatPayload(t *testing.T) {
	ch := &SlackChannel{WebhookURL: "https://hooks.slack.com/test"}
	evt := NotifyEvent{
		Kind:      EventNodeDowned,
		Address:   "pekko://System@10.0.0.2:2552",
		Roles:     []string{"payment"},
		DC:        "eu-west",
		Timestamp: time.Date(2026, 4, 20, 14, 30, 0, 0, time.UTC),
	}
	payload := ch.formatPayload(evt, "payment-down")
	if !strings.Contains(payload, "node.downed") {
		t.Errorf("payload should contain event kind, got: %s", payload)
	}
	if !strings.Contains(payload, "10.0.0.2:2552") {
		t.Errorf("payload should contain address, got: %s", payload)
	}
	if !strings.Contains(payload, "payment") {
		t.Errorf("payload should contain role, got: %s", payload)
	}
}
