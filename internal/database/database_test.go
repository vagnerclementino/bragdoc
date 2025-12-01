package database

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vagnerclementino/bragdoc/internal/database/queries"
)

func TestNew(t *testing.T) {
	// Create temporary directory for test database (real SQLite file)
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create database
	db, err := New(dbPath)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	// Verify connection is working
	err = db.conn.Ping()
	assert.NoError(t, err)
}

func TestMigrate(t *testing.T) {
	// Create temporary directory for test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create database
	db, err := New(dbPath)
	require.NoError(t, err)
	defer db.Close()

	// Run migrations
	err = db.Migrate(context.Background())
	require.NoError(t, err)

	// Verify migrations table exists
	var count int
	err = db.conn.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	require.NoError(t, err)
	assert.Greater(t, count, 0, "At least one migration should be applied")

	// Verify main tables exist
	tables := []string{"users", "brags", "tags", "brag_tags"}
	for _, table := range tables {
		var tableExists int
		query := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?"
		err = db.conn.QueryRow(query, table).Scan(&tableExists)
		require.NoError(t, err)
		assert.Equal(t, 1, tableExists, "Table %s should exist", table)
	}
}

func TestSetupDatabase(t *testing.T) {
	// Create temporary directory for test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Setup database (creates and migrates)
	db, err := SetupDatabase(dbPath)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	// Verify database is ready to use
	assert.NotNil(t, db.Queries())
	assert.NotNil(t, db.Conn())
}

func TestTransaction(t *testing.T) {
	// Create temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := SetupDatabase(dbPath)
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()

	// Test successful transaction
	t.Run("Successful Transaction", func(t *testing.T) {
		err := db.Transaction(ctx, func(q *queries.Queries) error {
			// Create a user within transaction
			_, err := q.CreateUser(ctx, queries.CreateUserParams{
				Name:     "Test User",
				Email:    "test@example.com",
				Language: "en",
			})
			return err
		})
		require.NoError(t, err)

		// Verify user was created
		users, err := db.Queries().ListUsers(ctx)
		require.NoError(t, err)
		assert.Len(t, users, 1)
	})

	// Test failed transaction (should rollback)
	t.Run("Failed Transaction Rollback", func(t *testing.T) {
		initialCount := 0
		users, err := db.Queries().ListUsers(ctx)
		require.NoError(t, err)
		initialCount = len(users)

		// Attempt transaction that will fail
		err = db.Transaction(ctx, func(q *queries.Queries) error {
			_, err := q.CreateUser(ctx, queries.CreateUserParams{
				Name:     "Another User",
				Email:    "another@example.com",
				Language: "en",
			})
			if err != nil {
				return err
			}
			// Force an error to trigger rollback
			return assert.AnError
		})
		require.Error(t, err)

		// Verify user count hasn't changed (rollback worked)
		users, err = db.Queries().ListUsers(ctx)
		require.NoError(t, err)
		assert.Equal(t, initialCount, len(users))
	})
}
