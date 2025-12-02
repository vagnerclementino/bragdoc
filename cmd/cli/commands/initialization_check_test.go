package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInitialized(t *testing.T) {
	// Save original home directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary directory for testing
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	t.Run("should return false when .bragdoc directory does not exist", func(t *testing.T) {
		assert.False(t, isInitialized())
	})

	t.Run("should return false when .bragdoc directory exists but database does not", func(t *testing.T) {
		configDir := filepath.Join(tempDir, ".bragdoc")
		err := os.MkdirAll(configDir, 0755)
		assert.NoError(t, err)

		assert.False(t, isInitialized())
	})

	t.Run("should return true when both .bragdoc directory and database exist", func(t *testing.T) {
		configDir := filepath.Join(tempDir, ".bragdoc")
		dbPath := filepath.Join(configDir, "bragdoc.db")

		err := os.MkdirAll(configDir, 0755)
		assert.NoError(t, err)

		// Create empty database file
		file, err := os.Create(dbPath)
		assert.NoError(t, err)
		file.Close()

		assert.True(t, isInitialized())
	})
}

func TestRequiresInitialization(t *testing.T) {
	// Save original home directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary directory for testing
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	t.Run("should return error when not initialized", func(t *testing.T) {
		checkFunc := requiresInitialization()
		err := checkFunc(nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "bragdoc is not initialized")
		assert.Contains(t, err.Error(), "bragdoc init")
	})

	t.Run("should return nil when initialized", func(t *testing.T) {
		// Setup initialized state
		configDir := filepath.Join(tempDir, ".bragdoc")
		dbPath := filepath.Join(configDir, "bragdoc.db")

		err := os.MkdirAll(configDir, 0755)
		assert.NoError(t, err)

		file, err := os.Create(dbPath)
		assert.NoError(t, err)
		file.Close()

		checkFunc := requiresInitialization()
		err = checkFunc(nil, nil)
		assert.NoError(t, err)
	})
}
