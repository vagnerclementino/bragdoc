package config

import (
	"time"
)

// Config represents the complete application configuration
// User data (name, email, job_title, company, locale) is stored in the database
// and should not be duplicated here
type Config struct {
	Database DatabaseConfig `yaml:"database" json:"database" toml:"database"`
	// Future configurations can be added here as needed:
	// AI       AIConfig       `yaml:"ai,omitempty" json:"ai,omitempty" toml:"ai,omitempty"`
	// Server   ServerConfig   `yaml:"server,omitempty" json:"server,omitempty" toml:"server,omitempty"`
	// Logging  LoggingConfig  `yaml:"logging,omitempty" json:"logging,omitempty" toml:"logging,omitempty"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Path string `yaml:"path" json:"path" toml:"path"` // Path to SQLite database file
}

// Commented out configurations for future use:
// These are not currently used by the application but are documented here
// for future implementation

// AIConfig represents AI provider configuration (not currently used)
// type AIConfig struct {
// 	Provider  string `yaml:"provider" json:"provider" toml:"provider"`       // AI provider (e.g., "openai", "anthropic")
// 	APIKey    string `yaml:"api_key" json:"api_key" toml:"api_key"`         // API key for the provider
// 	Model     string `yaml:"model" json:"model" toml:"model"`               // Model to use (e.g., "gpt-4")
// 	MaxTokens int    `yaml:"max_tokens" json:"max_tokens" toml:"max_tokens"` // Maximum tokens for generation
// }

// ServerConfig represents web server configuration (not currently used)
// type ServerConfig struct {
// 	Port        int    `yaml:"port" json:"port" toml:"port"`                         // Server port
// 	StaticDir   string `yaml:"static_dir" json:"static_dir" toml:"static_dir"`       // Static files directory
// 	CORSEnabled bool   `yaml:"cors_enabled" json:"cors_enabled" toml:"cors_enabled"` // Enable CORS
// }

// LoggingConfig represents logging configuration (not currently used)
// type LoggingConfig struct {
// 	Level    string `yaml:"level" json:"level" toml:"level"`             // Log level: debug, info, warn, error
// 	FilePath string `yaml:"file_path" json:"file_path" toml:"file_path"` // Log file path
// 	MaxSize  int    `yaml:"max_size" json:"max_size" toml:"max_size"`   // Max log file size in MB
// 	MaxAge   int    `yaml:"max_age" json:"max_age" toml:"max_age"`     // Max age of log files in days
// 	Console  bool   `yaml:"console" json:"console" toml:"console"`     // Also log to console
// }

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Database.Path == "" {
		return ErrInvalidConfig{Field: "database.path", Reason: "database path is required"}
	}
	return nil
}

// ErrInvalidConfig represents a configuration validation error
type ErrInvalidConfig struct {
	Field  string
	Reason string
}

func (e ErrInvalidConfig) Error() string {
	return "invalid config field " + e.Field + ": " + e.Reason
}

// GetDefaultConfig returns a default configuration
// User data is stored in the database, not in the config file
func GetDefaultConfig(configDir string) *Config {
	return &Config{
		Database: DatabaseConfig{
			Path: configDir + "/bragdoc.db",
		},
	}
}

// ConfigFormat represents supported configuration file formats
type ConfigFormat string

const (
	FormatYAML ConfigFormat = "yaml"
	FormatJSON ConfigFormat = "json"
	FormatTOML ConfigFormat = "toml"
)

// String returns the string representation of the format
func (f ConfigFormat) String() string {
	return string(f)
}

// Extension returns the file extension for the format
func (f ConfigFormat) Extension() string {
	switch f {
	case FormatYAML:
		return ".yaml"
	case FormatJSON:
		return ".json"
	case FormatTOML:
		return ".toml"
	default:
		return ".yaml"
	}
}

// ParseFormat parses a string into a ConfigFormat
func ParseFormat(s string) (ConfigFormat, error) {
	switch s {
	case "yaml", "yml", "":
		return FormatYAML, nil
	case "json":
		return FormatJSON, nil
	case "toml":
		return FormatTOML, nil
	default:
		return "", ErrUnsupportedFormat{Format: s}
	}
}

// ErrUnsupportedFormat represents an unsupported configuration format error
type ErrUnsupportedFormat struct {
	Format string
}

func (e ErrUnsupportedFormat) Error() string {
	return "unsupported configuration format: " + e.Format + " (v1 supports only YAML)"
}

// ConfigMetadata stores metadata about the configuration
type ConfigMetadata struct {
	Version   string    `yaml:"version" json:"version" toml:"version"`
	CreatedAt time.Time `yaml:"created_at" json:"created_at" toml:"created_at"`
	UpdatedAt time.Time `yaml:"updated_at" json:"updated_at" toml:"updated_at"`
}
