//go:build !cgo
// +build !cgo

package db

import (
	"fmt"
	"os"
	"path/filepath"
)

// This file provides stubs for when CGO is disabled.
// The real implementation is in db.go, but it requires CGO.

// NewNoCGO creates a new database connection
// When CGO is disabled, this will provide a meaningful error message
func NewNoCGO(dbPath string) (*DB, error) {
	// Ensure directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	return nil, fmt.Errorf("SQLite support requires CGO. Please install a C compiler (MinGW-w64 on Windows or GCC on Linux/macOS) or use a pre-built binary")
}
