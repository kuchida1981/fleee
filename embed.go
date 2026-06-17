package fleee

import "embed"

//go:embed migrations/*.sql
var MigrationFS embed.FS

//go:embed all:web/dist
var WebDistFS embed.FS
