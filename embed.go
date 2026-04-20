/*
 * embed.go
 * This file is part of the gekka-dashboard project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */

package main

import "embed"

//go:embed static/*
var staticFS embed.FS
