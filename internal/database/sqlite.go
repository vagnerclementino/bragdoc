package database

import (
	"context"
	"database/sql"

	"github.com/vagnerclementino/bragdoc/internal/database/queries"
)

// SQLiteDB wraps sql.DB and provides access to generated queries
type SQLiteDB struct {
	db      *sql.DB
	queries *queries.Queries
}

// NewSQLiteDB creates a new SQLite database wrapper
func NewSQLiteDB(db *sql.DB) *SQLiteDB {
	return &SQLiteDB{
		db:      db,
		queries: queries.New(db),
	}
}

// DB returns the underlying sql.DB
func (s *SQLiteDB) DB() *sql.DB {
	return s.db
}

// Queries returns the generated queries
func (s *SQLiteDB) Queries() *queries.Queries {
	return s.queries
}

// BeginTx starts a new transaction
func (s *SQLiteDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, opts)
}

// Close closes the database connection
func (s *SQLiteDB) Close() error {
	return s.db.Close()
}
