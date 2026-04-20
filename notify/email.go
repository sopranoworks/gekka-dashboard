/*
 * email.go
 * This file is part of the gekka-metrics project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package notify

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

// EmailConfig holds SMTP configuration for the email channel.
type EmailConfig struct {
	Host     string
	Port     int
	From     string
	To       []string
	Username string
	Password string
}

// EmailChannel dispatches notifications via SMTP email.
type EmailChannel struct {
	Host     string
	Port     int
	From     string
	To       []string
	Username string
	Password string
}

func (c *EmailChannel) Name() string { return "email" }

func (c *EmailChannel) Send(ctx context.Context, evt NotifyEvent, ruleName string) error {
	subject, body := c.formatMessage(evt, ruleName)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		c.From, strings.Join(c.To, ", "), subject, body)

	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	var auth smtp.Auth
	if c.Username != "" {
		auth = smtp.PlainAuth("", c.Username, c.Password, c.Host)
	}
	return smtp.SendMail(addr, auth, c.From, c.To, []byte(msg))
}

func (c *EmailChannel) formatMessage(evt NotifyEvent, ruleName string) (subject, body string) {
	subject = fmt.Sprintf("[gekka] %s — %s", evt.Kind, evt.Address)
	body = fmt.Sprintf("Rule: %s\nEvent: %s\nAddress: %s\nRoles: %s\nDC: %s\nTime: %s\n",
		ruleName, evt.Kind, evt.Address, strings.Join(evt.Roles, ", "), evt.DC,
		evt.Timestamp.Format("2006-01-02 15:04:05 UTC"))
	return
}
