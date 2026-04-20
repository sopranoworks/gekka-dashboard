/*
 * rule_test.go
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

func TestRule_Matches_EventType(t *testing.T) {
	r := Rule{
		Name:     "test",
		Events:   []EventKind{EventNodeUnreachable, EventNodeDowned},
		Channels: []string{"slack"},
	}
	evt := NotifyEvent{Kind: EventNodeUnreachable, Timestamp: time.Now()}
	if !r.Matches(evt) {
		t.Error("expected rule to match EventNodeUnreachable")
	}
	evt2 := NotifyEvent{Kind: EventNodeJoined, Timestamp: time.Now()}
	if r.Matches(evt2) {
		t.Error("expected rule to NOT match EventNodeJoined")
	}
}

func TestRule_Matches_RoleFilter(t *testing.T) {
	r := Rule{
		Name:     "cart-only",
		Events:   []EventKind{EventNodeUnreachable},
		Roles:    []string{"cart", "payment"},
		Channels: []string{"email"},
	}
	evt := NotifyEvent{Kind: EventNodeUnreachable, Roles: []string{"cart", "worker"}, Timestamp: time.Now()}
	if !r.Matches(evt) {
		t.Error("expected rule to match member with 'cart' role")
	}
	evt2 := NotifyEvent{Kind: EventNodeUnreachable, Roles: []string{"api"}, Timestamp: time.Now()}
	if r.Matches(evt2) {
		t.Error("expected rule to NOT match member without cart/payment role")
	}
}

func TestRule_Matches_EmptyRolesMatchesAll(t *testing.T) {
	r := Rule{
		Name:     "all-nodes",
		Events:   []EventKind{EventNodeDowned},
		Channels: []string{"slack"},
	}
	evt := NotifyEvent{Kind: EventNodeDowned, Roles: []string{"anything"}, Timestamp: time.Now()}
	if !r.Matches(evt) {
		t.Error("rule with no role filter should match any member")
	}
}
