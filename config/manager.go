package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// UseCase defines the configuration management interface
type UseCase interface {
	IsInitialized() bool
	Initialize(ctx context.Context, config *Config, format ConfigFormat) error
	Load(ctx context.Context) (*Config, error)
	Save(ctx context.Context, config *Config) error
	GetConfigPath() string
	GetDatabasePath() string
	GetDefaultConfig() *Config
}

// Manager implements the configuration management use case
type Manager struct {
	configDir  string
	configFile string
	format     ConfigFormat
	viper      *viper.Viper
}

// BragdocHomeEnv is the environment variable to override the default data directory.
// When set, bragdoc uses its value instead of ~/.bragdoc.
// This follows Viper's convention with the BRAGDOC prefix.
const BragdocHomeEnv = "BRAGDOC_HOME"

// ResolveBragdocHome returns the bragdoc data directory following Viper's
// precedence: env var (BRAGDOC_HOME) > default (~/.bragdoc).
// The config file itself lives inside this directory, so this resolution
// must happen before Viper loads the config file.
func ResolveBragdocHome() string {
	if envHome := os.Getenv(BragdocHomeEnv); envHome != "" {
		return expandHomeDir(envHome)
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".bragdoc")
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	return &Manager{
		configDir: ResolveBragdocHome(),
		viper:     viper.New(),
	}
}

// IsInitialized checks if bragdoc is already initialized
func (m *Manager) IsInitialized() bool {
	// Check for any configuration file format
	configFiles := []string{
		filepath.Join(m.configDir, "config.yaml"),
		filepath.Join(m.configDir, "config.yml"),
		filepath.Join(m.configDir, "config.json"),
		filepath.Join(m.configDir, "config.toml"),
	}

	for _, configFile := range configFiles {
		if _, err := os.Stat(configFile); err == nil {
			return true
		}
	}

	return false
}

// Initialize creates the configuration directory and file
func (m *Manager) Initialize(ctx context.Context, config *Config, format ConfigFormat) error {
	// V1: Only YAML format supported
	if format != FormatYAML && format != "" {
		return ErrUnsupportedFormat{Format: string(format)}
	}

	// Default to YAML if not specified
	if format == "" {
		format = FormatYAML
	}

	// Create .bragdoc directory
	if err := os.MkdirAll(m.configDir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create logs directory
	logsDir := filepath.Join(m.configDir, "logs")
	if err := os.MkdirAll(logsDir, 0750); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Set configuration file path
	m.format = format
	m.configFile = filepath.Join(m.configDir, "config"+format.Extension())

	// Validate configuration
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Save configuration file
	if err := m.Save(ctx, config); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	return nil
}

// Load reads the configuration from file
// Returns an empty config with defaults if not initialized
func (m *Manager) Load(_ context.Context) (*Config, error) {
	// If not initialized, return default config
	if !m.IsInitialized() {
		return m.GetDefaultConfig(), nil
	}

	// Detect configuration file format
	if err := m.detectConfigFile(); err != nil {
		return nil, err
	}

	// Configure viper
	m.viper.SetConfigFile(m.configFile)
	m.viper.SetConfigType(string(m.format))

	// Set defaults for new config fields
	m.viper.SetDefault("update_checker.enabled", true)

	// Enable environment variable substitution
	m.viper.AutomaticEnv()
	m.viper.SetEnvPrefix("BRAGDOC")
	m.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read configuration
	if err := m.viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal into Config struct
	var config Config
	if err := m.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Apply defaults for fields not present in config file
	if !m.viper.InConfig("update_checker") {
		config.UpdateChecker.Enabled = true
	}

	// Expand environment variables in paths
	config.Database.Path = os.ExpandEnv(config.Database.Path)

	// Expand ~ to home directory
	config.Database.Path = expandHomeDir(config.Database.Path)

	return &config, nil
}

// Save writes the configuration to file
func (m *Manager) Save(_ context.Context, config *Config) error {
	if m.configFile == "" {
		return fmt.Errorf("configuration file path not set")
	}

	// Validate configuration before saving
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Create file
	file, err := os.Create(m.configFile)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close config file: %w", closeErr)
		}
	}()

	// Encode based on format
	switch m.format {
	case FormatYAML:
		// Use viper to write YAML
		if m.viper == nil {
			m.viper = viper.New()
		}
		m.viper.SetConfigType("yaml")
		m.setViperConfig(config)
		if err := m.viper.WriteConfigAs(m.configFile); err != nil {
			return fmt.Errorf("failed to write YAML config: %w", err)
		}

	case FormatJSON:
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(config); err != nil {
			return fmt.Errorf("failed to encode JSON config: %w", err)
		}

	case FormatTOML:
		// TOML encoding would go here (future version)
		return ErrUnsupportedFormat{Format: string(m.format)}

	default:
		return ErrUnsupportedFormat{Format: string(m.format)}
	}

	return nil
}

// GetConfigPath returns the path to the configuration file
func (m *Manager) GetConfigPath() string {
	if m.configFile != "" {
		return m.configFile
	}

	// Try to detect existing config file
	if err := m.detectConfigFile(); err == nil {
		return m.configFile
	}

	// Return default YAML path
	return filepath.Join(m.configDir, "config.yaml")
}

// GetDatabasePath returns the path to the database file
func (m *Manager) GetDatabasePath() string {
	return filepath.Join(m.configDir, "bragdoc.db")
}

// GetDefaultConfig returns a default configuration
func (m *Manager) GetDefaultConfig() *Config {
	return GetDefaultConfig(m.configDir)
}

// detectConfigFile finds and sets the configuration file path and format
func (m *Manager) detectConfigFile() error {
	// Check for configuration files in order of preference
	configFiles := []struct {
		path   string
		format ConfigFormat
	}{
		{filepath.Join(m.configDir, "config.yaml"), FormatYAML},
		{filepath.Join(m.configDir, "config.yml"), FormatYAML},
		{filepath.Join(m.configDir, "config.json"), FormatJSON},
		{filepath.Join(m.configDir, "config.toml"), FormatTOML},
	}

	for _, cf := range configFiles {
		if _, err := os.Stat(cf.path); err == nil {
			m.configFile = cf.path
			m.format = cf.format
			return nil
		}
	}

	return fmt.Errorf("no configuration file found in %s", m.configDir)
}

// setViperConfig sets all configuration values in viper
func (m *Manager) setViperConfig(config *Config) {
	// Database configuration
	m.viper.Set("database.path", config.Database.Path)

	// Update checker configuration
	m.viper.Set("update_checker.enabled", config.UpdateChecker.Enabled)
	if !config.UpdateChecker.LastCheckedAt.IsZero() {
		m.viper.Set("update_checker.last_checked_at", config.UpdateChecker.LastCheckedAt)
	}
}

// ExpandHomeDir expands ~ to the user's home directory
// This function is exported for testing purposes
func ExpandHomeDir(path string) string {
	if !strings.HasPrefix(path, "~/") && path != "~" {
		return path
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if path == "~" {
		return homeDir
	}

	return filepath.Join(homeDir, path[2:])
}

// expandHomeDir is a private wrapper for backward compatibility
func expandHomeDir(path string) string {
	return ExpandHomeDir(path)
}
