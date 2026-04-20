/*
 * server_test.go
 * This file is part of the gekka-dashboard project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"context"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"nhooyr.io/websocket"
)

func TestHub_WebSocketConnect(t *testing.T) {
	hub := NewHub()
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.HandleWebSocket)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wsURL := "ws" + srv.URL[4:] + "/ws"
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var msg map[string]any
	if err := json.Unmarshal(data, &msg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if msg["type"] != "connected" {
		t.Errorf("expected type=connected, got %v", msg["type"])
	}

	if hub.ClientCount() != 1 {
		t.Errorf("expected 1 client, got %d", hub.ClientCount())
	}
}

func TestHub_Broadcast(t *testing.T) {
	hub := NewHub()
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.HandleWebSocket)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wsURL := "ws" + srv.URL[4:] + "/ws"
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	_, _, _ = conn.Read(ctx) // discard welcome

	hub.Broadcast(map[string]string{"type": "test", "data": "hello"})

	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read broadcast: %v", err)
	}
	var msg map[string]string
	if err := json.Unmarshal(data, &msg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if msg["type"] != "test" || msg["data"] != "hello" {
		t.Errorf("unexpected broadcast: %v", msg)
	}
}

func TestStaticFS_HasIndexHTML(t *testing.T) {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		t.Fatalf("fs.Sub: %v", err)
	}
	f, err := sub.Open("index.html")
	if err != nil {
		t.Fatalf("static/index.html not found: %v", err)
	}
	f.Close()
}
