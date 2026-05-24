package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var databaseDSN, migrationsPath, command string
	var migrationsTable string

	flag.StringVar(&databaseDSN, "database-dsn", "", "PostgreSQL DSN (e.g., 'postgres://user:password@host:port/dbname?sslmode=disable')")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations files (e.g., './migrations')")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table (optional, usually set in DSN or default)")
	flag.StringVar(&command, "command", "up", "migration command (up, down, force)")
	flag.Parse()

	if databaseDSN == "" {
		log.Fatal("database-dsn (PostgreSQL DSN) is required")
	}
	if migrationsPath == "" {
		log.Fatal("migrations-path is required")
	}

	finalDSN := fmt.Sprintf("%s&x-migrations-table=%s", databaseDSN, migrationsTable)

	m, err := migrate.New(
		"file://"+migrationsPath,
		finalDSN,
	)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	switch command {
	case "up":
		log.Println("Applying migrations UP...")
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("No migrations to apply")
				return
			}
			log.Fatalf("Failed to apply migrations UP: %v", err)
		}
		fmt.Println("Migrations applied successfully (UP)")
	case "down":
		log.Println("Applying migrations DOWN...")
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("No migrations to rollback")
				return
			}
			log.Fatalf("Failed to apply migrations DOWN: %v", err)
		}
		fmt.Println("Migrations applied successfully (DOWN)")
	case "force":
		var version int
		fmt.Print("Enter migration version to force: ")
		_, err := fmt.Scanf("%d", &version)
		if err != nil {
			log.Fatalf("Invalid version: %v", err)
		}
		log.Printf("Forcing migration version %d...", version)
		if err := m.Force(version); err != nil {
			log.Fatalf("Failed to force migration version %d: %v", version, err)
		}
		fmt.Printf("Migration version %d forced successfully.\n", version)
	default:
		log.Fatalf("Unknown command: %s. Use 'up' or 'down' or 'force'.", command)
	}
}
