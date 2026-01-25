package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vagnerclementino/bragdoc/internal/database/queries"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// DB wraps the database connection and queries
type DB struct {
	conn    *sql.DB
	queries *queries.Queries
}

// New creates a new database connection
func New(dbPath string) (*DB, error) {
	// Validate database path
	if err := validateDatabasePath(dbPath); err != nil {
		return nil, err
	}

	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Verify directory is writable
	if err := checkWritePermission(dir); err != nil {
		return nil, fmt.Errorf("database directory is not writable: %w", err)
	}

	// Open SQLite connection
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(1) // SQLite works best with single connection
	conn.SetMaxIdleConns(1)

	db := &DB{
		conn:    conn,
		queries: queries.New(conn),
	}

	return db, nil
}

// validateDatabasePath validates the database path
func validateDatabasePath(dbPath string) error {
	if dbPath == "" {
		return fmt.Errorf("database path cannot be empty")
	}

	// Check if path is a directory (not allowed)
	if info, err := os.Stat(dbPath); err == nil && info.IsDir() {
		return fmt.Errorf("database path cannot be a directory: %s", dbPath)
	}

	return nil
}

// checkWritePermission checks if the directory has write permissions
func checkWritePermission(dir string) error {
	// Try to create a temporary file to test write permissions
	testFile := filepath.Join(dir, ".write_test")
	// #nosec G304 - Test file path is constructed from validated directory
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("no write permission in directory %s: %w", dir, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close test file: %w", err)
	}
	if err := os.Remove(testFile); err != nil {
		return fmt.Errorf("failed to remove test file: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Queries returns the SQLC generated queries
func (db *DB) Queries() *queries.Queries {
	return db.queries
}

// Conn returns the underlying database connection
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// Migrate runs all pending migrations
func (db *DB) Migrate(ctx context.Context) error {
	// Create migrations table if it doesn't exist
	_, err := db.conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := db.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Read migration files
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort migration files
	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Apply pending migrations
	for _, filename := range migrationFiles {
		version := strings.TrimSuffix(filename, ".sql")

		// Skip if already applied
		if appliedMigrations[version] {
			continue
		}

		// Read migration file
		content, err := migrationsFS.ReadFile(filepath.Join("migrations", filename))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", filename, err)
		}

		// Execute migration in a transaction
		tx, err := db.conn.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %s: %w", filename, err)
		}

		// Execute migration SQL
		if _, err := tx.ExecContext(ctx, string(content)); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("failed to execute migration %s: %w (rollback error: %v)", filename, err, rbErr)
			}
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		// Record migration
		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("failed to record migration %s: %w (rollback error: %v)", filename, err, rbErr)
			}
			return fmt.Errorf("failed to record migration %s: %w", filename, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", filename, err)
		}
	}

	return nil
}

// getAppliedMigrations returns a map of applied migration versions
func (db *DB) getAppliedMigrations(ctx context.Context) (map[string]bool, error) {
	rows, err := db.conn.QueryContext(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close rows: %v\n", err)
		}
	}()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// Transaction executes a function within a database transaction
func (db *DB) Transaction(ctx context.Context, fn func(*queries.Queries) error) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	qtx := db.queries.WithTx(tx)

	if err := fn(qtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// SetupDatabase creates and initializes a new database with migrations
func SetupDatabase(dbPath string) (*DB, error) {
	db, err := New(dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(context.Background()); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, fmt.Errorf("failed to run migrations: %w (close error: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// NewNullString creates a sql.NullString from a string
// Returns a valid NullString if the string is not empty, otherwise returns an invalid NullString
func NewNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
