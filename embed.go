package fleee

import "embed"

// MigrationFS embeds the database migrations directory
//
//go:embed migrations/*.sql
var MigrationFS embed.FS
