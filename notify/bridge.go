/*
 * bridge.go
 * This file is part of the gekka-metrics project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package notify

import (
	"context"
	"time"

	"github.com/sopranoworks/gekka/cluster"
)

// BridgeClusterEvents reads from a ClusterManager channel subscription and
// feeds enriched NotifyEvents into the Engine. Blocks until ctx is cancelled.
func BridgeClusterEvents(ctx context.Context, sub *cluster.ChanSubscription, cm *cluster.ClusterManager, eng *Engine) {
	defer sub.Cancel()
	for {
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-sub.C:
			if !ok {
				return
			}
			ne := translateEvent(evt, cm)
			if ne != nil {
				eng.HandleEvent(*ne)
			}
		}
	}
}

func translateEvent(evt cluster.ClusterDomainEvent, cm *cluster.ClusterManager) *NotifyEvent {
	var kind EventKind
	var memberAddr cluster.MemberAddress

	switch e := evt.(type) {
	case cluster.MemberUp:
		kind = EventNodeJoined
		memberAddr = e.Member
	case cluster.MemberLeft:
		kind = EventNodeLeft
		memberAddr = e.Member
	case cluster.MemberExited:
		kind = EventNodeLeft
		memberAddr = e.Member
	case cluster.MemberDowned:
		kind = EventNodeDowned
		memberAddr = e.Member
	case cluster.MemberRemoved:
		kind = EventNodeRemoved
		memberAddr = e.Member
	case cluster.UnreachableMember:
		kind = EventNodeUnreachable
		memberAddr = e.Member
	case cluster.ReachableMember:
		kind = EventNodeReachable
		memberAddr = e.Member
	default:
		return nil
	}

	roles := lookupRoles(cm, memberAddr)

	return &NotifyEvent{
		Kind:      kind,
		Address:   memberAddr.String(),
		Roles:     roles,
		DC:        memberAddr.DataCenter,
		Timestamp: time.Now(),
	}
}

func lookupRoles(cm *cluster.ClusterManager, addr cluster.MemberAddress) []string {
	cm.Mu.RLock()
	gossip := cm.State
	cm.Mu.RUnlock()

	if gossip == nil {
		return nil
	}

	for _, m := range gossip.GetMembers() {
		a := gossip.GetAllAddresses()[m.GetAddressIndex()].GetAddress()
		if a.GetHostname() == addr.Host && a.GetPort() == addr.Port {
			return cluster.GetRolesForMember(gossip, m)
		}
	}
	return nil
}
