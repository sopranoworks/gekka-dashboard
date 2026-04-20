/*
 * channel.go
 * This file is part of the gekka-metrics project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package notify

import "context"

// Channel is the interface for notification dispatch targets.
type Channel interface {
	Name() string
	Send(ctx context.Context, evt NotifyEvent, ruleName string) error
}
