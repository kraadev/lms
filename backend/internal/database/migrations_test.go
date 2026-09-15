package database

import (
	"testing"

	"lms/internal/config"
)

func TestRunMigrations(t *testing.T) {
	cfg := &config.Config{
		DBDriver:     "sqlite",
		DBSQLitePath: ":memory:",
	}
	db, err := Connect(cfg)
	if err != nil {
		t.Fatalf("failed to connect in-memory sqlite: %v", err)
	}
	defer db.Close()

	if err := RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Verify organizations table exists
	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='organizations'").Scan(&name)
	if err != nil || name != "organizations" {
		t.Errorf("expected organizations table to exist")
	}

	// Verify system_settings table exists
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='system_settings'").Scan(&name)
	if err != nil || name != "system_settings" {
		t.Errorf("expected system_settings table to exist")
	}
}
