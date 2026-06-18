package store_test

import (
	"testing"

	"github.com/kosuke/fleee"
	"github.com/kosuke/fleee/internal/store"
)

func TestDB_Migrate(t *testing.T) {
	// Setup db without migrating automatically
	db, err := store.NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer func() { _ = db.Close() }()

	// Execute migrations
	if err := db.Migrate(fleee.MigrationFS); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Verify accounts table exists
	var count int
	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='accounts'").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query sqlite_master: %v", err)
	}
	if count != 1 {
		t.Errorf("expected accounts table to exist (count 1), got %d", count)
	}
}

func TestDB_Migrate_Idempotent(t *testing.T) {
	db, err := store.NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer func() { _ = db.Close() }()

	// 1st migration
	if err := db.Migrate(fleee.MigrationFS); err != nil {
		t.Fatalf("failed to run first migration: %v", err)
	}

	// 2nd migration should succeed without error
	if err := db.Migrate(fleee.MigrationFS); err != nil {
		t.Fatalf("failed to run second migration: %v", err)
	}

	// Check table again
	var count int
	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='accounts'").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query sqlite_master: %v", err)
	}
	if count != 1 {
		t.Errorf("expected accounts table to exist (count 1) after idempotent run, got %d", count)
	}
}
