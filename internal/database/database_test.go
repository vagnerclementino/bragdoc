package database

import (
	"context"
	"os"
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

// TestNew_ParentDirectoryCreation tests the parent directory creation property
// Feature: cli-architecture-refactor, Property 3: Parent directories are created
// For any valid database path, if the parent directory doesn't exist, the system should create it
func TestNew_ParentDirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()

	testCases := []struct {
		name   string
		dbPath string
	}{
		{
			name:   "single_level_directory",
			dbPath: filepath.Join(tmpDir, "data", "test.db"),
		},
		{
			name:   "nested_directories",
			dbPath: filepath.Join(tmpDir, "data", "nested", "deep", "test.db"),
		},
		{
			name:   "very_deep_nesting",
			dbPath: filepath.Join(tmpDir, "a", "b", "c", "d", "e", "f", "test.db"),
		},
		{
			name:   "with_dots_in_path",
			dbPath: filepath.Join(tmpDir, ".config", "bragdoc", "test.db"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Verify parent directory doesn't exist initially
			parentDir := filepath.Dir(tc.dbPath)
			_, err := os.Stat(parentDir)
			assert.True(t, os.IsNotExist(err), "parent directory should not exist initially")

			// Create database - should create parent directories
			db, err := New(tc.dbPath)
			require.NoError(t, err, "New should create parent directories")
			require.NotNil(t, db)
			defer db.Close()

			// Property 1: Parent directory should now exist
			stat, err := os.Stat(parentDir)
			require.NoError(t, err, "parent directory should exist after New()")
			assert.True(t, stat.IsDir(), "parent path should be a directory")

			// Property 2: Database should be functional (can execute queries)
			err = db.conn.Ping()
			assert.NoError(t, err, "database connection should work")

			// Property 3: After running migrations, database file should exist
			err = db.Migrate(context.Background())
			require.NoError(t, err, "migrations should run successfully")

			_, err = os.Stat(tc.dbPath)
			require.NoError(t, err, "database file should exist after migrations")
		})
	}
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

// TestNew_PathValidation tests the path validation property
// Feature: cli-architecture-refactor, Property 4: Inaccessible paths are rejected
// For any database path that is not accessible, the system should return a descriptive error
func TestNew_PathValidation(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("empty_path_rejected", func(t *testing.T) {
		// Property: Empty paths should be rejected
		db, err := New("")
		assert.Error(t, err, "empty path should be rejected")
		assert.Nil(t, db)
		assert.Contains(t, err.Error(), "empty", "error should mention empty path")
	})

	t.Run("directory_as_database_rejected", func(t *testing.T) {
		// Property: Directories cannot be used as database files
		dirPath := filepath.Join(tmpDir, "testdir")
		err := os.MkdirAll(dirPath, 0755)
		require.NoError(t, err)

		db, err := New(dirPath)
		assert.Error(t, err, "directory path should be rejected")
		assert.Nil(t, db)
		assert.Contains(t, err.Error(), "directory", "error should mention directory")
	})

	t.Run("read_only_directory_rejected", func(t *testing.T) {
		// Property: Paths in read-only directories should be rejected
		readOnlyDir := filepath.Join(tmpDir, "readonly")
		err := os.MkdirAll(readOnlyDir, 0755)
		require.NoError(t, err)

		// Make directory read-only
		err = os.Chmod(readOnlyDir, 0555)
		require.NoError(t, err)
		defer os.Chmod(readOnlyDir, 0755) // Restore permissions for cleanup

		dbPath := filepath.Join(readOnlyDir, "test.db")
		db, err := New(dbPath)
		assert.Error(t, err, "read-only directory should be rejected")
		assert.Nil(t, db)
		assert.Contains(t, err.Error(), "not writable", "error should mention write permission")
	})

	t.Run("valid_paths_accepted", func(t *testing.T) {
		// Property: Valid writable paths should be accepted
		validPaths := []string{
			filepath.Join(tmpDir, "valid1.db"),
			filepath.Join(tmpDir, "subdir", "valid2.db"),
			filepath.Join(tmpDir, ".hidden", "valid3.db"),
		}

		for _, dbPath := range validPaths {
			t.Run(dbPath, func(t *testing.T) {
				db, err := New(dbPath)
				assert.NoError(t, err, "valid path should be accepted")
				assert.NotNil(t, db)
				if db != nil {
					db.Close()
				}
			})
		}
	})
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
				Name:   "Test User",
				Email:  "test@example.com",
				Locale: "en-US",
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
				Name:   "Another User",
				Email:  "another@example.com",
				Locale: "en-US",
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
