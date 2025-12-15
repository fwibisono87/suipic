package main

import (
	"flag"
	"log"

	"github.com/suipic/backend/config"
	"github.com/suipic/backend/db"
)

func main() {
	var action string
	flag.StringVar(&action, "action", "up", "Migration action: up or down")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := db.Connect(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Printf("Running migrations (%s)...", action)

	switch action {
	case "up":
		if err := db.RunMigrations(db.DB, "./db/migrations"); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations completed successfully")
	case "down":
		if err := db.RollbackMigration(db.DB, "./db/migrations"); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		log.Println("Rollback completed successfully")
	default:
		log.Fatalf("Unknown action: %s (use 'up' or 'down')", action)
	}
}
