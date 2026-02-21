package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"kanban-backend/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

type MigrationStatus struct {
	CurrentVersion *uint
	Dirty          bool
	TotalUpFiles   int
}

type MigrationRunner struct {
	migrator       *migrate.Migrate
	db             *sql.DB
	dialect        string
	migrationsPath string
}

func NewMigrationRunnerFromConfig(migrationsPath string) (*MigrationRunner, error) {
	if hasPostgresConfig() {
		if config.DB == nil {
			config.ConnectDB()
		}

		return NewPostgresMigrationRunner(config.DB, migrationsPath)
	}

	sqliteDB, err := sql.Open("sqlite3", "file:migrations_dev.db?cache=shared")
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	return NewSQLiteMigrationRunner(sqliteDB, migrationsPath)
}

func hasPostgresConfig() bool {
	return os.Getenv("DB_HOST") != "" &&
		os.Getenv("DB_PORT") != "" &&
		os.Getenv("DB_USER") != "" &&
		os.Getenv("DB_NAME") != ""
}

func NewPostgresMigrationRunner(gormDB *gorm.DB, migrationsPath string) (*MigrationRunner, error) {
	if gormDB == nil {
		return nil, errors.New("gorm database is nil")
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db from gorm: %w", err)
	}

	if err := ensureMigrationsTable(sqlDB); err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{MigrationsTable: "migrations"})
	if err != nil {
		return nil, fmt.Errorf("create postgres migration driver: %w", err)
	}

	return newRunnerWithDriver(driver, sqlDB, "postgres", migrationsPath)
}

func NewSQLiteMigrationRunner(sqlDB *sql.DB, migrationsPath string) (*MigrationRunner, error) {
	if sqlDB == nil {
		return nil, errors.New("sql db is nil")
	}

	if err := ensureMigrationsTable(sqlDB); err != nil {
		return nil, err
	}

	driver, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{MigrationsTable: "migrations"})
	if err != nil {
		return nil, fmt.Errorf("create sqlite migration driver: %w", err)
	}

	return newRunnerWithDriver(driver, sqlDB, "sqlite3", migrationsPath)
}

func newRunnerWithDriver(driver database.Driver, sqlDB *sql.DB, dialect string, migrationsPath string) (*MigrationRunner, error) {
	absolutePath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("resolve migrations path: %w", err)
	}

	if _, err := loadMigrationFiles(absolutePath); err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+absolutePath, dialect, driver)
	if err != nil {
		return nil, fmt.Errorf("initialize migrator: %w", err)
	}

	return &MigrationRunner{
		migrator:       m,
		db:             sqlDB,
		dialect:        dialect,
		migrationsPath: absolutePath,
	}, nil
}

func loadMigrationFiles(migrationsPath string) ([]string, error) {
	entries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("read migrations directory: %w", err)
	}

	files := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasSuffix(entry.Name(), ".up.sql") || strings.HasSuffix(entry.Name(), ".down.sql") {
			files = append(files, entry.Name())
		}
	}

	sort.Strings(files)
	return files, nil
}

func ensureMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			version BIGINT NOT NULL PRIMARY KEY,
			dirty BOOLEAN NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	return nil
}

func (r *MigrationRunner) Up() error {
	err := r.migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations up: %w", err)
	}

	return nil
}

func (r *MigrationRunner) Down() error {
	err := r.migrator.Steps(-1)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("rollback migration: %w", err)
	}

	return nil
}

func (r *MigrationRunner) Status() (MigrationStatus, error) {
	upFiles, err := r.UpFiles()
	if err != nil {
		return MigrationStatus{}, err
	}

	version, dirty, err := r.migrator.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return MigrationStatus{
				CurrentVersion: nil,
				Dirty:          false,
				TotalUpFiles:   len(upFiles),
			}, nil
		}

		return MigrationStatus{}, fmt.Errorf("get migration status: %w", err)
	}

	v := version
	return MigrationStatus{
		CurrentVersion: &v,
		Dirty:          dirty,
		TotalUpFiles:   len(upFiles),
	}, nil
}

func (r *MigrationRunner) UpFiles() ([]string, error) {
	allFiles, err := loadMigrationFiles(r.migrationsPath)
	if err != nil {
		return nil, err
	}

	upFiles := make([]string, 0)
	for _, file := range allFiles {
		if strings.HasSuffix(file, ".up.sql") {
			upFiles = append(upFiles, file)
		}
	}

	return upFiles, nil
}

func (r *MigrationRunner) Close() error {
	sourceErr, dbErr := r.migrator.Close()
	if sourceErr != nil {
		return sourceErr
	}

	if dbErr != nil {
		return dbErr
	}

	return nil
}
