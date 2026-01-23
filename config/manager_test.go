package config

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()

	assert.NotNil(t, manager)
	assert.NotEmpty(t, manager.configDir)
	assert.NotNil(t, manager.viper)
}

func TestManager_IsInitialized(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	manager := &Manager{
		configDir: tempDir,
	}

	// Should not be initialized initially
	assert.False(t, manager.IsInitialized())

	// Create a config file
	configFile := filepath.Join(tempDir, "config.yaml")
	err := os.WriteFile(configFile, []byte("user:\n  name: Test\n"), 0644)
	require.NoError(t, err)

	// Should be initialized now
	assert.True(t, manager.IsInitialized())
}

func TestManager_Initialize(t *testing.T) {
	tempDir := t.TempDir()

	manager := &Manager{
		configDir: tempDir,
	}

	config := GetDefaultConfig(tempDir)
	ctx := context.Background()

	// Test YAML initialization (v1 only supported format)
	err := manager.Initialize(ctx, config, FormatYAML)
	require.NoError(t, err)

	// Check that config file was created
	assert.FileExists(t, filepath.Join(tempDir, "config.yaml"))

	// Check that logs directory was created
	assert.DirExists(t, filepath.Join(tempDir, "logs"))
}

func TestManager_Initialize_UnsupportedFormat(t *testing.T) {
	tempDir := t.TempDir()

	manager := &Manager{
		configDir: tempDir,
	}

	config := GetDefaultConfig(tempDir)
	ctx := context.Background()

	// Test JSON initialization (not supported in v1)
	err := manager.Initialize(ctx, config, FormatJSON)
	require.Error(t, err)
	assert.IsType(t, ErrUnsupportedFormat{}, err)

	// Test TOML initialization (not supported in v1)
	err = manager.Initialize(ctx, config, FormatTOML)
	require.Error(t, err)
	assert.IsType(t, ErrUnsupportedFormat{}, err)
}

func TestManager_Initialize_InvalidConfig(t *testing.T) {
	tempDir := t.TempDir()

	manager := &Manager{
		configDir: tempDir,
	}

	// Config with missing required fields
	config := &Config{
		Database: DatabaseConfig{
			// Missing path
		},
	}

	ctx := context.Background()

	err := manager.Initialize(ctx, config, FormatYAML)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid configuration")
}

func TestManager_SaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()

	manager := &Manager{
		configDir:  tempDir,
		configFile: filepath.Join(tempDir, "config.yaml"),
		format:     FormatYAML,
	}

	originalConfig := GetDefaultConfig(tempDir)
	ctx := context.Background()

	// Save configuration
	err := manager.Save(ctx, originalConfig)
	require.NoError(t, err)

	// Load configuration
	loadedConfig, err := manager.Load(ctx)
	require.NoError(t, err)

	// Verify loaded configuration matches original
	assert.Equal(t, originalConfig.Database.Path, loadedConfig.Database.Path)
}

func TestManager_Load_NotInitialized(t *testing.T) {
	tempDir := t.TempDir()

	manager := &Manager{
		configDir: tempDir,
	}

	ctx := context.Background()

	// Load configuration when not initialized
	loadedConfig, err := manager.Load(ctx)
	require.NoError(t, err)
	assert.NotNil(t, loadedConfig)

	// Should return default config
	defaultConfig := manager.GetDefaultConfig()
	assert.Equal(t, defaultConfig.Database.Path, loadedConfig.Database.Path)
}

func TestManager_GetConfigPath(t *testing.T) {
	tempDir := t.TempDir()

	manager := &Manager{
		configDir:  tempDir,
		configFile: filepath.Join(tempDir, "config.yaml"),
	}

	path := manager.GetConfigPath()
	assert.Equal(t, filepath.Join(tempDir, "config.yaml"), path)
}

func TestManager_GetDatabasePath(t *testing.T) {
	tempDir := t.TempDir()

	manager := &Manager{
		configDir: tempDir,
	}

	path := manager.GetDatabasePath()
	assert.Equal(t, filepath.Join(tempDir, "bragdoc.db"), path)
}

func TestManager_GetDefaultConfig(t *testing.T) {
	tempDir := t.TempDir()

	manager := &Manager{
		configDir: tempDir,
	}

	config := manager.GetDefaultConfig()

	assert.NotNil(t, config)
	assert.Contains(t, config.Database.Path, tempDir)
}

func TestManager_detectConfigFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name           string
		createFile     string
		expectedFormat ConfigFormat
	}{
		{
			name:           "detect yaml",
			createFile:     "config.yaml",
			expectedFormat: FormatYAML,
		},
		{
			name:           "detect yml",
			createFile:     "config.yml",
			expectedFormat: FormatYAML,
		},
		{
			name:           "detect json",
			createFile:     "config.json",
			expectedFormat: FormatJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh temp directory for each test
			testDir := filepath.Join(tempDir, tt.name)
			err := os.MkdirAll(testDir, 0755)
			require.NoError(t, err)

			manager := &Manager{
				configDir: testDir,
			}

			// Create the config file
			configPath := filepath.Join(testDir, tt.createFile)
			err = os.WriteFile(configPath, []byte("test"), 0644)
			require.NoError(t, err)

			// Detect config file
			err = manager.detectConfigFile()
			require.NoError(t, err)

			assert.Equal(t, configPath, manager.configFile)
			assert.Equal(t, tt.expectedFormat, manager.format)
		})
	}
}

func TestManager_detectConfigFile_NotFound(t *testing.T) {
	tempDir := t.TempDir()

	manager := &Manager{
		configDir: tempDir,
	}

	err := manager.detectConfigFile()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no configuration file found")
}

func TestExpandHomeDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "expand tilde",
			input: "~/.bragdoc/config.yaml",
			want:  filepath.Join(homeDir, ".bragdoc/config.yaml"),
		},
		{
			name:  "tilde only",
			input: "~",
			want:  homeDir,
		},
		{
			name:  "no tilde",
			input: "/absolute/path",
			want:  "/absolute/path",
		},
		{
			name:  "relative path",
			input: "relative/path",
			want:  "relative/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpandHomeDir(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestExpandHomeDir_Property tests the tilde expansion property
// Feature: cli-architecture-refactor, Property 2: Tilde expansion works correctly
// For any path that begins with ~/, the system should expand the til to the home directory
func TestExpandHomeDir_Property(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	// Property 1: Paths starting with ~/ should always be expanded to homeDir + rest of path
	t.Run("tilde_slash_expansion", func(t *testing.T) {
		testCases := []string{
			"~/.bragdoc",
			"~/config.yaml",
			"~/.config/bragdoc/db.sqlite",
			"~/Documents/bragdoc",
			"~/a",
			"~/very/long/path/to/some/file.txt",
		}

		for _, path := range testCases {
			t.Run(path, func(t *testing.T) {
				result := ExpandHomeDir(path)

				// Property: Result should start with home directory
				assert.True(t, filepath.IsAbs(result), "expanded path should be absolute")
				assert.True(t, strings.HasPrefix(result, homeDir),
					"expanded path should start with home directory")

				// Property: The suffix after ~/ should be preserved
				expectedSuffix := path[2:] // Remove "~/"
				assert.True(t, strings.HasSuffix(result, expectedSuffix),
					"expanded path should preserve the suffix after ~/")

				// Property: Result should equal homeDir + suffix
				expected := filepath.Join(homeDir, expectedSuffix)
				assert.Equal(t, expected, result)
			})
		}
	})

	// Property 2: Paths not starting with ~ should remain unchanged
	t.Run("non_tilde_paths_unchanged", func(t *testing.T) {
		testCases := []string{
			"/absolute/path",
			"relative/path",
			"./current/dir",
			"../parent/dir",
			"",
			"path/with/~tilde/in/middle",
		}

		for _, path := range testCases {
			t.Run(path, func(t *testing.T) {
				result := ExpandHomeDir(path)
				assert.Equal(t, path, result, "non-tilde paths should remain unchanged")
			})
		}
	})

	// Property 3: Tilde alone should expand to home directory
	t.Run("tilde_alone", func(t *testing.T) {
		result := ExpandHomeDir("~")
		assert.Equal(t, homeDir, result, "tilde alone should expand to home directory")
	})

	// Property 4: Idempotence - expanding twice should give same result
	t.Run("idempotence", func(t *testing.T) {
		testCases := []string{
			"~/.bragdoc",
			"~/config.yaml",
			"/absolute/path",
			"relative/path",
		}

		for _, path := range testCases {
			t.Run(path, func(t *testing.T) {
				firstExpansion := ExpandHomeDir(path)
				secondExpansion := ExpandHomeDir(firstExpansion)
				assert.Equal(t, firstExpansion, secondExpansion,
					"expanding an already expanded path should not change it")
			})
		}
	})
}
