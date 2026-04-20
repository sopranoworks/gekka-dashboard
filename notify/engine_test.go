/*
 * engine_test.go
 * This file is part of the gekka-metrics project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package notify

import (
	"context"
	"sync"
	"testing"
	"time"
)

type mockChannel struct {
	mu    sync.Mutex
	name  string
	sent  []NotifyEvent
	rules []string
}

func (m *mockChannel) Name() string { return m.name }
func (m *mockChannel) Send(_ context.Context, evt NotifyEvent, ruleName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sent = append(m.sent, evt)
	m.rules = append(m.rules, ruleName)
	return nil
}
func (m *mockChannel) sentCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.sent)
}

func TestEngine_DispatchesMatchingEvent(t *testing.T) {
	mock := &mockChannel{name: "test-ch"}
	rules := []*Rule{{
		Name:     "r1",
		Events:   []EventKind{EventNodeUnreachable},
		Channels: []string{"test-ch"},
	}}

	eng := NewEngine(rules, map[string]Channel{"test-ch": mock})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go eng.Run(ctx)

	eng.HandleEvent(NotifyEvent{
		Kind:      EventNodeUnreachable,
		Address:   "pekko://System@10.0.0.1:2552",
		Timestamp: time.Now(),
	})

	time.Sleep(50 * time.Millisecond)

	if mock.sentCount() != 1 {
		t.Errorf("expected 1 dispatch, got %d", mock.sentCount())
	}
}

func TestEngine_DoesNotDispatchNonMatching(t *testing.T) {
	mock := &mockChannel{name: "test-ch"}
	rules := []*Rule{{
		Name:     "r1",
		Events:   []EventKind{EventNodeDowned},
		Channels: []string{"test-ch"},
	}}

	eng := NewEngine(rules, map[string]Channel{"test-ch": mock})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go eng.Run(ctx)

	eng.HandleEvent(NotifyEvent{
		Kind:      EventNodeJoined,
		Timestamp: time.Now(),
	})

	time.Sleep(50 * time.Millisecond)

	if mock.sentCount() != 0 {
		t.Errorf("expected 0 dispatches, got %d", mock.sentCount())
	}
}

func TestEngine_ThrottlesRepeatedEvents(t *testing.T) {
	mock := &mockChannel{name: "test-ch"}
	rules := []*Rule{{
		Name:     "r1",
		Events:   []EventKind{EventNodeUnreachable},
		Channels: []string{"test-ch"},
		Throttle: 1 * time.Hour,
	}}

	eng := NewEngine(rules, map[string]Channel{"test-ch": mock})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go eng.Run(ctx)

	evt := NotifyEvent{Kind: EventNodeUnreachable, Timestamp: time.Now()}
	eng.HandleEvent(evt)
	eng.HandleEvent(evt)

	time.Sleep(50 * time.Millisecond)

	if mock.sentCount() != 1 {
		t.Errorf("expected 1 dispatch (throttled), got %d", mock.sentCount())
	}
}
