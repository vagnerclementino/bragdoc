package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vagnerclementino/bragdoc/config"
)

// TestExpandPath tests the expandPath function
func TestExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "expand tilde only",
			input:    "~",
			expected: homeDir,
		},
		{
			name:     "expand tilde with path",
			input:    "~/.bragdoc/bragdoc.db",
			expected: filepath.Join(homeDir, ".bragdoc", "bragdoc.db"),
		},
		{
			name:     "no tilde - absolute path",
			input:    "/var/data/bragdoc.db",
			expected: "/var/data/bragdoc.db",
		},
		{
			name:     "no tilde - relative path",
			input:    "./data/bragdoc.db",
			expected: "./data/bragdoc.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandPath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetDatabasePath_UsesConfigPath tests that getDatabasePath uses the config path when provided
// **Validates: Requirements 3.1**
// Property 1: Config database path is used when provided
func TestGetDatabasePath_UsesConfigPath(t *testing.T) {
	tests := []struct {
		name           string
		configPath     string
		expectedResult string
	}{
		{
			name:           "uses config path with tilde",
			configPath:     "~/custom/path/db.sqlite",
			expectedResult: expandPath("~/custom/path/db.sqlite"),
		},
		{
			name:           "uses config path absolute",
			configPath:     "/var/data/bragdoc.db",
			expectedResult: "/var/data/bragdoc.db",
		},
		{
			name:           "uses config path relative",
			configPath:     "./data/bragdoc.db",
			expectedResult: "./data/bragdoc.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Database: config.DatabaseConfig{
					Path: tt.configPath,
				},
			}

			result := getDatabasePath(cfg)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

// TestGetDatabasePath_UsesDefaultWhenEmpty tests that getDatabasePath uses default when config path is empty
func TestGetDatabasePath_UsesDefaultWhenEmpty(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Path: "",
		},
	}

	result := getDatabasePath(cfg)

	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)
	expected := filepath.Join(homeDir, ".bragdoc", "bragdoc.db")

	assert.Equal(t, expected, result)
}
