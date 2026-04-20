/*
 * email_test.go
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

func TestEmailChannel_FormatMessage(t *testing.T) {
	ch := &EmailChannel{
		From: "alerts@example.com",
		To:   []string{"ops@example.com"},
	}
	evt := NotifyEvent{
		Kind:      EventNodeUnreachable,
		Address:   "pekko://System@10.0.0.1:2552",
		Roles:     []string{"cart", "worker"},
		DC:        "us-east",
		Timestamp: time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC),
	}
	subject, body := ch.formatMessage(evt, "critical-alert")
	if !strings.Contains(subject, "node.unreachable") {
		t.Errorf("subject should contain event kind, got: %s", subject)
	}
	if !strings.Contains(body, "10.0.0.1:2552") {
		t.Errorf("body should contain address, got: %s", body)
	}
	if !strings.Contains(body, "cart") {
		t.Errorf("body should contain roles, got: %s", body)
	}
}
