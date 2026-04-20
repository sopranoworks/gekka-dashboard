/*
 * server.go
 * This file is part of the gekka-dashboard project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"io/fs"
	"log/slog"
	"net/http"
)

func startServer(addr string, wsHub *Hub) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", wsHub.HandleWebSocket)

	staticSub, err := fs.Sub(staticFS, "static")
	if err != nil {
		return err
	}
	fileServer := http.FileServer(http.FS(staticSub))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		f, err := staticSub.Open(path[1:])
		if err == nil {
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})

	slog.Info("dashboard: HTTP server starting", "addr", addr)
	return http.ListenAndServe(addr, mux)
}
