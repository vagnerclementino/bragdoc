package config

import (
	"context"
	"os"
	"path/filepath"
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

	user := UserConfig{
		Name:   "Test User",
		Email:  "test@example.com",
		Locale: "en-US",
	}

	config := GetDefaultConfig(user, tempDir)
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

	user := UserConfig{
		Name:   "Test User",
		Email:  "test@example.com",
		Locale: "en-US",
	}

	config := GetDefaultConfig(user, tempDir)
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
		User: UserConfig{
			Name: "Test User",
			// Missing email
		},
		Database: DatabaseConfig{
			Path: tempDir + "/bragdoc.db",
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

	user := UserConfig{
		Name:     "Test User",
		Email:    "test@example.com",
		JobTitle: "Developer",
		Company:  "Test Corp",
		Locale:   "en-US",
	}

	originalConfig := GetDefaultConfig(user, tempDir)
	ctx := context.Background()

	// Save configuration
	err := manager.Save(ctx, originalConfig)
	require.NoError(t, err)

	// Load configuration
	loadedConfig, err := manager.Load(ctx)
	require.NoError(t, err)

	// Verify loaded configuration matches original (core fields)
	assert.Equal(t, originalConfig.User.Name, loadedConfig.User.Name)
	assert.Equal(t, originalConfig.User.Email, loadedConfig.User.Email)
	assert.Equal(t, originalConfig.User.Locale, loadedConfig.User.Locale)
	assert.Equal(t, originalConfig.AI.Provider, loadedConfig.AI.Provider)
	assert.Equal(t, originalConfig.AI.Model, loadedConfig.AI.Model)
	assert.Equal(t, originalConfig.Server.Port, loadedConfig.Server.Port)
	
	// Optional fields should be preserved if set
	assert.Equal(t, originalConfig.User.JobTitle, loadedConfig.User.JobTitle)
	assert.Equal(t, originalConfig.User.Company, loadedConfig.User.Company)
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

	user := UserConfig{
		Name:   "Test User",
		Email:  "test@example.com",
		Locale: "en-US",
	}

	config := manager.GetDefaultConfig(user)

	assert.NotNil(t, config)
	assert.Equal(t, user.Name, config.User.Name)
	assert.Equal(t, user.Email, config.User.Email)
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
			got := expandHomeDir(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
