package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	// SQLite driver
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

// newMigrate creates a new migrate instance
func (db *DB) newMigrate() (*migrate.Migrate, error) {
	// Create source from embedded filesystem
	sourceFS, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create sub filesystem: %w", err)
	}

	sourceDriver, err := iofs.New(sourceFS, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to create source driver: %w", err)
	}

	// Create database driver with NoLock to prevent closing the connection
	dbDriver, err := sqlite3.WithInstance(db.conn, &sqlite3.Config{
		NoTxWrap: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create database driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite3", dbDriver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return m, nil
}

// Migrate runs all pending migrations (up)
func (db *DB) Migrate(_ context.Context) error {
	m, err := db.newMigrate()
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Don't close - let the DB instance manage the connection lifecycle
	return nil
}

// MigrateDown rolls back the last migration
func (db *DB) MigrateDown(ctx context.Context) error {
	m, err := db.newMigrate()
	if err != nil {
		return err
	}

	if err := m.Steps(-1); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}

// MigrateVersion returns the current migration version
func (db *DB) MigrateVersion() (uint, bool, error) {
	m, err := db.newMigrate()
	if err != nil {
		return 0, false, err
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

// MigrateTo migrates to a specific version
func (db *DB) MigrateTo(ctx context.Context, version uint) error {
	m, err := db.newMigrate()
	if err != nil {
		return err
	}

	if err := m.Migrate(version); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}

	return nil
}

// MigrateForce forces the migration version (use with caution)
func (db *DB) MigrateForce(version int) error {
	m, err := db.newMigrate()
	if err != nil {
		return err
	}

	if err := m.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	return nil
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
