/*
 * throttle.go
 * This file is part of the gekka-metrics project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package notify

import (
	"sync"
	"time"
)

// ThrottleTracker tracks the last dispatch time per rule name and prevents
// repeated notifications within the configured interval.
type ThrottleTracker struct {
	mu   sync.Mutex
	last map[string]time.Time
}

// NewThrottleTracker creates a new throttle tracker.
func NewThrottleTracker() *ThrottleTracker {
	return &ThrottleTracker{last: make(map[string]time.Time)}
}

// Allow returns true if the rule is allowed to fire (not within throttle interval).
// If allowed, records the current time as the last fire time.
func (t *ThrottleTracker) Allow(ruleName string, interval time.Duration) bool {
	if interval <= 0 {
		return true
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	last, ok := t.last[ruleName]
	if ok && time.Since(last) < interval {
		return false
	}
	t.last[ruleName] = time.Now()
	return true
}
