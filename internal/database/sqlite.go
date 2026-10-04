package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

// Queries uses the transaction bound to ctx, or the normal connection otherwise.
// All repositories sharing this SQLiteDB must propagate the request context.
func (s *SQLiteDB) Queries(ctx context.Context) *queries.Queries {
	if tx, ok := ctx.Value(transactionKey{db: s}).(*sql.Tx); ok {
		return s.queries.WithTx(tx)
	}
	return s.queries
}

type transactionKey struct{ db *SQLiteDB }

// WithinTransaction commits a complete operation or rolls back on error/panic.
// The transaction belongs to this request; other requests keep their own contexts.
func (s *SQLiteDB) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	if ctx.Value(transactionKey{db: s}) != nil {
		return fmt.Errorf("nested transactions are not supported")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Also releases the transaction on panic.
	if err := fn(context.WithValue(ctx, transactionKey{db: s}, tx)); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			return fmt.Errorf("rollback transaction after %v: %w", err, rollbackErr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

// BeginTx starts a new transaction
func (s *SQLiteDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, opts)
}

// Close closes the database connection
func (s *SQLiteDB) Close() error {
	return s.db.Close()
}
