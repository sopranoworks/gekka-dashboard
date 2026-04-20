/*
 * slack.go
 * This file is part of the gekka-metrics project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package notify

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
)

// SlackChannel dispatches notifications via Slack incoming webhook.
type SlackChannel struct {
	WebhookURL string
}

func (c *SlackChannel) Name() string { return "slack" }

func (c *SlackChannel) Send(ctx context.Context, evt NotifyEvent, ruleName string) error {
	payload := c.formatPayload(evt, ruleName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.WebhookURL,
		bytes.NewBufferString(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("slack webhook returned %d", resp.StatusCode)
	}
	return nil
}

func (c *SlackChannel) formatPayload(evt NotifyEvent, ruleName string) string {
	roles := strings.Join(evt.Roles, ", ")
	if roles == "" {
		roles = "-"
	}
	text := fmt.Sprintf("*[%s]* %s\n>Address: `%s`\n>Roles: %s | DC: %s\n>Rule: %s | %s",
		evt.Kind, evt.Address, evt.Address, roles, evt.DC, ruleName,
		evt.Timestamp.Format("15:04:05 UTC"))
	return fmt.Sprintf(`{"text":%q}`, text)
}
