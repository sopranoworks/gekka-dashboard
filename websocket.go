/*
 * websocket.go
 * This file is part of the gekka-dashboard project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*websocket.Conn]struct{})}
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		slog.Error("websocket: accept failed", "err", err)
		return
	}

	h.mu.Lock()
	h.clients[conn] = struct{}{}
	h.mu.Unlock()

	slog.Info("websocket: client connected", "remote", r.RemoteAddr)

	welcome := map[string]any{"type": "connected", "version": "0.1.0"}
	data, _ := json.Marshal(welcome)
	_ = conn.Write(r.Context(), websocket.MessageText, data)

	for {
		_, _, err := conn.Read(r.Context())
		if err != nil {
			break
		}
	}

	h.mu.Lock()
	delete(h.clients, conn)
	h.mu.Unlock()

	slog.Info("websocket: client disconnected", "remote", r.RemoteAddr)
}

func (h *Hub) Broadcast(msg any) {
	data, err := json.Marshal(msg)
	if err != nil {
		slog.Error("websocket: marshal failed", "err", err)
		return
	}

	h.mu.RLock()
	clients := make([]*websocket.Conn, 0, len(h.clients))
	for c := range h.clients {
		clients = append(clients, c)
	}
	h.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, c := range clients {
		if err := c.Write(ctx, websocket.MessageText, data); err != nil {
			slog.Debug("websocket: write failed, removing client", "err", err)
			h.mu.Lock()
			delete(h.clients, c)
			h.mu.Unlock()
			c.Close(websocket.StatusGoingAway, "write failed")
		}
	}
}

func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
