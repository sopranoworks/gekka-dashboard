/*
 * throttle_test.go
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

func TestThrottle_AllowsFirst(t *testing.T) {
	tr := NewThrottleTracker()
	if !tr.Allow("rule-a", 5*time.Minute) {
		t.Error("first call should always be allowed")
	}
}

func TestThrottle_BlocksWithinInterval(t *testing.T) {
	tr := NewThrottleTracker()
	tr.Allow("rule-a", 5*time.Minute)
	if tr.Allow("rule-a", 5*time.Minute) {
		t.Error("second call within interval should be blocked")
	}
}

func TestThrottle_AllowsAfterInterval(t *testing.T) {
	tr := NewThrottleTracker()
	tr.last["rule-a"] = time.Now().Add(-6 * time.Minute)
	if !tr.Allow("rule-a", 5*time.Minute) {
		t.Error("call after interval elapsed should be allowed")
	}
}

func TestThrottle_ZeroDurationAlwaysAllows(t *testing.T) {
	tr := NewThrottleTracker()
	tr.Allow("rule-b", 0)
	if !tr.Allow("rule-b", 0) {
		t.Error("zero throttle should always allow")
	}
}

func TestThrottle_IndependentRules(t *testing.T) {
	tr := NewThrottleTracker()
	tr.Allow("rule-a", 5*time.Minute)
	if !tr.Allow("rule-b", 5*time.Minute) {
		t.Error("different rules should throttle independently")
	}
}
