/*
 * event.go
 * This file is part of the gekka-metrics project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package notify

import "time"

// EventKind identifies the type of cluster event for rule matching.
type EventKind string

const (
	EventNodeJoined          EventKind = "node.joined"
	EventNodeLeft            EventKind = "node.left"
	EventNodeUnreachable     EventKind = "node.unreachable"
	EventNodeDowned          EventKind = "node.downed"
	EventNodeRemoved         EventKind = "node.removed"
	EventNodeReachable       EventKind = "node.reachable"
	EventClusterSplit        EventKind = "cluster.split"
	EventClusterConverged    EventKind = "cluster.converged"
	EventShardRebalanceStart EventKind = "shard.rebalance.started"
	EventShardRebalanceDone  EventKind = "shard.rebalance.completed"
	EventShardAllocFailed    EventKind = "shard.allocation.failed"
	EventHealthHBTimeout     EventKind = "health.heartbeat.timeout"
	EventHealthRTTDegraded   EventKind = "health.rtt.degraded"
)

// NotifyEvent is a cluster event enriched with metadata for notification dispatch.
type NotifyEvent struct {
	Kind      EventKind
	Address   string   // member address (e.g. "pekko://System@10.0.0.1:2552")
	Roles     []string // roles of the affected member
	DC        string   // data-center
	Timestamp time.Time
}
