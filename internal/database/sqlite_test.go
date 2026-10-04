package database

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vagnerclementino/bragdoc/internal/database/queries"
)

func TestWithinTransaction(t *testing.T) {
	for _, outcome := range []string{"commit", "error", "cancel", "panic"} {
		t.Run(outcome, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			db, err := New(filepath.Join(t.TempDir(), "transaction.db"))
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, db.Close()) })
			require.NoError(t, db.Migrate(ctx))
			sqlite := NewSQLiteDB(db.Conn())
			txCtx, cancelTx := context.WithCancel(ctx)
			defer cancelTx()
			failure := errors.New("injected failure")
			run := func() error {
				return sqlite.WithinTransaction(txCtx, func(requestCtx context.Context) error {
					_, err := sqlite.Queries(requestCtx).CreateUser(requestCtx, queries.CreateUserParams{
						Name: "Transaction User", Email: "transaction@example.com", Locale: "en-US",
					})
					if err != nil {
						return err
					}
					switch outcome {
					case "error":
						return failure
					case "cancel":
						cancelTx() // A successful callback must not hide a failed commit.
					case "panic":
						panic(failure)
					}
					return nil
				})
			}
			if outcome == "panic" {
				assert.Panics(t, func() { _ = run() })
			} else {
				err := run()
				if outcome == "commit" {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
					if outcome == "error" {
						assert.ErrorIs(t, err, failure)
					}
				}
			}
			// A different context uses the pool, not a completed transaction.
			users, err := sqlite.Queries(ctx).ListUsers(ctx)
			require.NoError(t, err)
			if outcome == "commit" {
				assert.Len(t, users, 1)
			} else {
				assert.Empty(t, users)
			}
		})
	}
}
