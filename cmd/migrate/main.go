package main

import (
	"fmt"
	"os"

	"github.com/touros-platform/api/internal/config"
	"github.com/touros-platform/api/internal/database"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/migrate/main.go [up|down]")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	db, err := database.NewConnection(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "up":
		if err := database.AutoMigrate(db); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to run migrations: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations completed successfully")
	case "down":
		fmt.Println("Down migrations not implemented - use manual rollback if needed")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

