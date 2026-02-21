# Migrations

## Overview
Database migrations for the Kanban application using [golang-migrate](https://github.com/golang-migrate/migrate).

## Database Schema
The application uses the following 12 tables:
- users
- boards
- columns
- tasks
- comments
- labels
- attachments
- notifications
- task_labels (junction table)
- members (junction table)
- refresh_tokens
- audit_logs

## Migration Pattern
golang-migrate uses versioned SQL files with up/down migrations.

### File Naming Convention
- Up migration: `00001_create_table.up.sql`
- Down migration: `00001_create_table.down.sql`

The version number (e.g., `00001`) must be unique and increment with each new migration.

### Up/Down Pattern
- **Up file**: Contains SQL to apply the migration (CREATE TABLE, ALTER TABLE, etc.)
- **Down file**: Contains SQL to rollback the migration (DROP TABLE, etc.)

## Running Migrations

### CLI Tool (to be implemented)
```bash
# Apply all pending migrations
go run migrations/main.go up

# Rollback the last migration
go run migrations/main.go down

# Show migration status
go run migrations/main.go status
```

### Using migrate CLI
```bash
# Create a new migration
migrate create -ext sql -dir migrations/sql -seq create_users_table

# Apply migrations
migrate -path migrations/sql -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up

# Rollback
migrate -path migrations/sql -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" down 1
```

## Database Connection
Migrations use the same PostgreSQL connection configured in `config/database.go`.

## Notes
- Migrations are versioned and tracked in a `schema_migrations` table
- Always test down migrations before committing
- Do not modify existing migration files once deployed
