package migrations

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrationRunnerSQLiteLifecycle(t *testing.T) {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	defer db.Close()

	migrationsDir := t.TempDir()

	upSQL := `CREATE TABLE users (id TEXT PRIMARY KEY, email TEXT NOT NULL UNIQUE);`
	downSQL := `DROP TABLE users;`

	if err := os.WriteFile(filepath.Join(migrationsDir, "000001_create_users.up.sql"), []byte(upSQL), 0o644); err != nil {
		t.Fatalf("write up migration: %v", err)
	}

	if err := os.WriteFile(filepath.Join(migrationsDir, "000001_create_users.down.sql"), []byte(downSQL), 0o644); err != nil {
		t.Fatalf("write down migration: %v", err)
	}

	runner, err := NewSQLiteMigrationRunner(db, migrationsDir)
	if err != nil {
		t.Fatalf("new sqlite migration runner: %v", err)
	}
	defer runner.Close()

	var migrationsTable string
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='migrations'`).Scan(&migrationsTable)
	if err != nil {
		t.Fatalf("verify migrations table exists: %v", err)
	}

	statusBefore, err := runner.Status()
	if err != nil {
		t.Fatalf("status before up: %v", err)
	}

	if statusBefore.TotalUpFiles != 1 {
		t.Fatalf("expected 1 up migration file, got %d", statusBefore.TotalUpFiles)
	}

	if err := runner.Up(); err != nil {
		t.Fatalf("run up: %v", err)
	}

	var tableName string
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users'`).Scan(&tableName)
	if err != nil {
		t.Fatalf("users table not created: %v", err)
	}

	statusAfterUp, err := runner.Status()
	if err != nil {
		t.Fatalf("status after up: %v", err)
	}

	if statusAfterUp.CurrentVersion == nil || *statusAfterUp.CurrentVersion != 1 {
		t.Fatalf("expected current version 1 after up, got %+v", statusAfterUp.CurrentVersion)
	}

	if err := runner.Down(); err != nil {
		t.Fatalf("run down: %v", err)
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users'`).Scan(&tableName)
	if err == nil {
		t.Fatalf("users table should be dropped after down")
	}

	statusAfterDown, err := runner.Status()
	if err != nil {
		t.Fatalf("status after down: %v", err)
	}

	if statusAfterDown.CurrentVersion != nil {
		t.Fatalf("expected nil version after down, got %d", *statusAfterDown.CurrentVersion)
	}
}
