//go:build ignore

package main

import (
	"fmt"
	"log"
	"os"

	"kanban-backend/migrations"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	runner, err := migrations.NewMigrationRunnerFromConfig(defaultMigrationsPath())
	if err != nil {
		log.Fatalf("failed to initialize migration runner: %v", err)
	}
	defer runner.Close()

	command := os.Args[1]
	switch command {
	case "up":
		if err := runner.Up(); err != nil {
			log.Fatalf("migration up failed: %v", err)
		}
		fmt.Println("migrations applied")
	case "down":
		if err := runner.Down(); err != nil {
			log.Fatalf("migration down failed: %v", err)
		}
		fmt.Println("last migration rolled back")
	case "status":
		status, err := runner.Status()
		if err != nil {
			log.Fatalf("status check failed: %v", err)
		}

		if status.CurrentVersion == nil {
			fmt.Printf("current_version: none, dirty: %t, total_migrations: %d\n", status.Dirty, status.TotalUpFiles)
		} else {
			fmt.Printf("current_version: %d, dirty: %t, total_migrations: %d\n", *status.CurrentVersion, status.Dirty, status.TotalUpFiles)
		}
	default:
		printUsage()
		os.Exit(1)
	}
}

func defaultMigrationsPath() string {
	return "migrations"
}

func printUsage() {
	fmt.Println("usage: go run migrations/main.go [up|down|status]")
}
