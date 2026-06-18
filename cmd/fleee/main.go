package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/kosuke/fleee"
	"github.com/kosuke/fleee/internal/handler"
	"github.com/kosuke/fleee/internal/importer"
	"github.com/kosuke/fleee/internal/server"
	"github.com/kosuke/fleee/internal/store"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]
	switch subcommand {
	case "serve":
		serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
		port := serveCmd.String("port", "8080", "Port to bind HTTP server to")
		dbPath := serveCmd.String("db", "fleee.db", "SQLite database file path")

		if err := serveCmd.Parse(os.Args[2:]); err != nil {
			log.Fatalf("Failed to parse flags: %v", err)
		}

		runServer(*port, *dbPath)

	default:
		fmt.Printf("Unknown subcommand: %s\n", subcommand)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: fleee <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  serve    Start the API server")
	fmt.Println("    -port  Port number (default: 8080)")
	fmt.Println("    -db    Database file path (default: fleee.db)")
}

func runServer(port, dbPath string) {
	// 1. Initialize SQLite Database
	db, err := store.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer func() { _ = db.Close() }()

	// 2. Execute pending database migrations
	if err := db.Migrate(fleee.MigrationFS); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 3. Set up repository, importer and handler layers
	accountStore := store.NewAccountStore(db)
	journalEntryStore := store.NewJournalEntryStore(db)
	journalEntryHandler := handler.NewJournalEntryHandler(journalEntryStore)
	accountImporter := importer.NewAccountImporter(accountStore)
	accountHandler := handler.NewAccountHandler(accountStore, accountImporter)

	// 4. Setup and start Server
	webFS, err := fs.Sub(fleee.WebDistFS, "web/dist")
	if err != nil {
		log.Fatalf("Failed to get web dist FS: %v", err)
	}
	srv := server.NewServer(port, accountHandler, journalEntryHandler, webFS)
	if err := srv.Start(); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
