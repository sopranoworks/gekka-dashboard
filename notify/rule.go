/*
 * rule.go
 * This file is part of the gekka-metrics project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package notify

import "time"

// Rule defines a notification rule: which events to match, optional role
// filter, which channels to dispatch to, and a throttle interval.
type Rule struct {
	Name     string
	Events   []EventKind
	Roles    []string      // empty = match all roles
	Channels []string      // channel names to dispatch to
	Throttle time.Duration // minimum interval between repeated alerts
}

// Matches returns true if the event satisfies this rule's event type and role filters.
func (r *Rule) Matches(evt NotifyEvent) bool {
	if !r.matchesEvent(evt.Kind) {
		return false
	}
	if !r.matchesRoles(evt.Roles) {
		return false
	}
	return true
}

func (r *Rule) matchesEvent(kind EventKind) bool {
	for _, e := range r.Events {
		if e == kind {
			return true
		}
	}
	return false
}

func (r *Rule) matchesRoles(memberRoles []string) bool {
	if len(r.Roles) == 0 {
		return true // no filter = match all
	}
	for _, required := range r.Roles {
		for _, has := range memberRoles {
			if required == has {
				return true
			}
		}
	}
	return false
}
