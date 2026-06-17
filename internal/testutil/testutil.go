package testutil

import (
	"testing"

	"github.com/kosuke/fleee"
	"github.com/kosuke/fleee/internal/store"
)

// NewTestDB creates an in-memory SQLite database with migrations applied.
func NewTestDB(t *testing.T) *store.DB {
	t.Helper()
	db, err := store.NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}
	if err := db.Migrate(fleee.MigrationFS); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}
