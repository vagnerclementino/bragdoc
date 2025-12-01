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
	SetupDatabase(ctx context.Context) error
	GetDefaultConfig(user UserConfig) *Config
	GetDefaultPrompts() PromptsConfig
}

// Manager implements the configuration management use case
type Manager struct {
	configDir  string
	configFile string
	format     ConfigFormat
	viper      *viper.Viper
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home is not available
		homeDir = "."
	}

	return &Manager{
		configDir: filepath.Join(homeDir, ".bragdoc"),
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
	if err := os.MkdirAll(m.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create logs directory
	logsDir := filepath.Join(m.configDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
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
func (m *Manager) Load(ctx context.Context) (*Config, error) {
	// Detect configuration file format
	if err := m.detectConfigFile(); err != nil {
		return nil, err
	}

	// Configure viper
	m.viper.SetConfigFile(m.configFile)
	m.viper.SetConfigType(string(m.format))

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

	// Expand environment variables in paths
	config.Database.Path = os.ExpandEnv(config.Database.Path)
	config.Logging.FilePath = os.ExpandEnv(config.Logging.FilePath)

	// Expand ~ to home directory
	config.Database.Path = expandHomeDir(config.Database.Path)
	config.Logging.FilePath = expandHomeDir(config.Logging.FilePath)

	return &config, nil
}

// Save writes the configuration to file
func (m *Manager) Save(ctx context.Context, config *Config) error {
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
	defer file.Close()

	// Encode based on format
	switch m.format {
	case FormatYAML:
		// Use viper to write YAML
		if m.viper == nil {
			m.viper = viper.New()
		}
		m.viper.SetConfigType("yaml")
		if err := m.setViperConfig(config); err != nil {
			return err
		}
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

// SetupDatabase is a placeholder for database setup
// The actual database setup is handled by the database package
func (m *Manager) SetupDatabase(ctx context.Context) error {
	// This method is kept for interface compatibility
	// The actual database setup happens in the database package
	return nil
}

// GetDefaultConfig returns a default configuration with user data
func (m *Manager) GetDefaultConfig(user UserConfig) *Config {
	return GetDefaultConfig(user, m.configDir)
}

// GetDefaultPrompts returns default AI prompts
func (m *Manager) GetDefaultPrompts() PromptsConfig {
	return GetDefaultPrompts()
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
func (m *Manager) setViperConfig(config *Config) error {
	// User configuration
	m.viper.Set("user.name", config.User.Name)
	m.viper.Set("user.email", config.User.Email)
	if config.User.JobTitle != "" {
		m.viper.Set("user.job_title", config.User.JobTitle)
	}
	if config.User.Company != "" {
		m.viper.Set("user.company", config.User.Company)
	}
	m.viper.Set("user.locale", config.User.Locale)

	// Database configuration
	m.viper.Set("database.path", config.Database.Path)

	// AI configuration
	m.viper.Set("ai.provider", config.AI.Provider)
	m.viper.Set("ai.api_key", config.AI.APIKey)
	m.viper.Set("ai.model", config.AI.Model)
	m.viper.Set("ai.max_tokens", config.AI.MaxTokens)

	// Server configuration
	m.viper.Set("server.port", config.Server.Port)
	m.viper.Set("server.static_dir", config.Server.StaticDir)
	m.viper.Set("server.cors_enabled", config.Server.CORSEnabled)

	// Prompts configuration
	m.viper.Set("prompts.enhance_description", config.Prompts.EnhanceDescription)
	m.viper.Set("prompts.generate_document", config.Prompts.GenerateDocument)
	m.viper.Set("prompts.suggest_tags", config.Prompts.SuggestTags)
	m.viper.Set("prompts.translate_brag", config.Prompts.TranslateBrag)

	// Logging configuration
	m.viper.Set("logging.level", config.Logging.Level)
	m.viper.Set("logging.file_path", config.Logging.FilePath)
	m.viper.Set("logging.max_size", config.Logging.MaxSize)
	m.viper.Set("logging.max_age", config.Logging.MaxAge)
	m.viper.Set("logging.console", config.Logging.Console)

	// I18n configuration
	m.viper.Set("i18n.language", config.I18n.Language)
	m.viper.Set("i18n.locale", config.I18n.Locale)

	return nil
}

// expandHomeDir expands ~ to the user's home directory
func expandHomeDir(path string) string {
	if !strings.HasPrefix(path, "~") {
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
