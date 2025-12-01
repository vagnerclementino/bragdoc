package config

import (
	"time"
)

// Config represents the complete application configuration
type Config struct {
	User     UserConfig     `yaml:"user" json:"user" toml:"user"`
	Database DatabaseConfig `yaml:"database" json:"database" toml:"database"`
	AI       AIConfig       `yaml:"ai" json:"ai" toml:"ai"`
	Server   ServerConfig   `yaml:"server" json:"server" toml:"server"`
	Prompts  PromptsConfig  `yaml:"prompts" json:"prompts" toml:"prompts"`
	Logging  LoggingConfig  `yaml:"logging" json:"logging" toml:"logging"`
	I18n     I18nConfig     `yaml:"i18n" json:"i18n" toml:"i18n"`
}

// UserConfig represents user profile information
type UserConfig struct {
	Name     string `yaml:"name" json:"name" toml:"name" mapstructure:"name"`
	Email    string `yaml:"email" json:"email" toml:"email" mapstructure:"email"`
	JobTitle string `yaml:"job_title,omitempty" json:"job_title,omitempty" toml:"job_title,omitempty" mapstructure:"job_title"`
	Company  string `yaml:"company,omitempty" json:"company,omitempty" toml:"company,omitempty" mapstructure:"company"`
	Locale   string `yaml:"locale" json:"locale" toml:"locale" mapstructure:"locale"` // Locale format: language-COUNTRY (e.g., en-US, pt-BR, pt-PT)
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Path string `yaml:"path" json:"path" toml:"path"`
}

// AIConfig represents AI provider configuration
type AIConfig struct {
	Provider  string `yaml:"provider" json:"provider" toml:"provider"`
	APIKey    string `yaml:"api_key" json:"api_key" toml:"api_key"`
	Model     string `yaml:"model" json:"model" toml:"model"`
	MaxTokens int    `yaml:"max_tokens" json:"max_tokens" toml:"max_tokens"`
}

// ServerConfig represents web server configuration
type ServerConfig struct {
	Port        int    `yaml:"port" json:"port" toml:"port"`
	StaticDir   string `yaml:"static_dir" json:"static_dir" toml:"static_dir"`
	CORSEnabled bool   `yaml:"cors_enabled" json:"cors_enabled" toml:"cors_enabled"`
}

// PromptsConfig represents AI prompts configuration
type PromptsConfig struct {
	EnhanceDescription string `yaml:"enhance_description" json:"enhance_description" toml:"enhance_description"`
	GenerateDocument   string `yaml:"generate_document" json:"generate_document" toml:"generate_document"`
	SuggestTags        string `yaml:"suggest_tags" json:"suggest_tags" toml:"suggest_tags"`
	TranslateBrag      string `yaml:"translate_brag" json:"translate_brag" toml:"translate_brag"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level    string `yaml:"level" json:"level" toml:"level"`         // debug, info, warn, error
	FilePath string `yaml:"file_path" json:"file_path" toml:"file_path"` // ~/.bragdoc/logs/bragdoc.log
	MaxSize  int    `yaml:"max_size" json:"max_size" toml:"max_size"`   // MB
	MaxAge   int    `yaml:"max_age" json:"max_age" toml:"max_age"`     // days
	Console  bool   `yaml:"console" json:"console" toml:"console"`     // also log to console
}

// I18nConfig represents internationalization configuration
type I18nConfig struct {
	Language string `yaml:"language" json:"language" toml:"language"` // en, pt, es, fr, etc.
	Locale   string `yaml:"locale" json:"locale" toml:"locale"`     // en_US, pt_BR, es_ES, etc.
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.User.Name == "" {
		return ErrInvalidConfig{Field: "user.name", Reason: "name is required"}
	}
	if c.User.Email == "" {
		return ErrInvalidConfig{Field: "user.email", Reason: "email is required"}
	}
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

// GetDefaultConfig returns a default configuration with user data
func GetDefaultConfig(user UserConfig, configDir string) *Config {
	return &Config{
		User: user,
		Database: DatabaseConfig{
			Path: configDir + "/bragdoc.db",
		},
		AI: AIConfig{
			Provider:  "openai",
			APIKey:    "${OPENAI_API_KEY}",
			Model:     "gpt-4",
			MaxTokens: 2000,
		},
		Server: ServerConfig{
			Port:        8080,
			StaticDir:   "./web/static",
			CORSEnabled: true,
		},
		Prompts:  GetDefaultPrompts(),
		Logging:  GetDefaultLoggingConfig(configDir),
		I18n:     GetDefaultI18nConfig(),
	}
}

// GetDefaultPrompts returns default AI prompts
func GetDefaultPrompts() PromptsConfig {
	return PromptsConfig{
		EnhanceDescription: `You are an assistant specialized in improving professional achievement descriptions.

Task: Improve the following achievement description, making it more impactful and specific for use in professional promotion documents.

Original description: {{.Description}}
Target language: {{.Language}}

Guidelines:
- Use action verbs in past tense
- Include metrics when possible
- Highlight impact and results
- Maintain professional tone
- Maximum 200 words
- Write the improved description in {{.Language}}

Improved description:`,

		GenerateDocument: `You are an expert in creating professional achievement documents (brag documents).

Task: Create a professional document based on the following achievements in {{.Language}}.

{{range .Brags}}
- **{{.Title}}** ({{.Category}}): {{.Description}}
{{if .Details}}Details: {{.Details}}{{end}}
{{end}}

User profile:
- Name: {{.User.Name}}
- Job Title: {{.User.JobTitle}}
- Company: {{.User.Company}}

Guidelines:
- Write the entire document in {{.Language}}
- Organize by categories
- Use professional language
- Highlight impact and results
- Include introduction and conclusion
- Format suitable for performance reviews

Document:`,

		SuggestTags: `Analyze the following achievement and suggest 3-5 relevant tags:

Title: {{.Title}}
Description: {{.Description}}

Suggest tags that are:
- Specific and relevant
- Useful for categorization
- In English
- One word or short term

Suggested tags (comma-separated):`,

		TranslateBrag: `You are a professional translator specialized in career achievements and professional documents.

Task: Translate and adapt the following professional achievement to {{.TargetLanguage}}.

Original Title: {{.Title}}
Original Description: {{.Description}}

Guidelines:
- Translate to {{.TargetLanguage}} maintaining professional tone
- Adapt cultural references if needed
- Keep technical terms accurate
- Maintain the impact and metrics
- Use appropriate action verbs for {{.TargetLanguage}}

Provide the translation in this exact format:
Title: [translated title]
Description: [translated description]`,
	}
}

// GetDefaultLoggingConfig returns default logging configuration
func GetDefaultLoggingConfig(configDir string) LoggingConfig {
	return LoggingConfig{
		Level:    "info",
		FilePath: configDir + "/logs/bragdoc.log",
		MaxSize:  10,
		MaxAge:   30,
		Console:  false,
	}
}

// GetDefaultI18nConfig returns default internationalization configuration
func GetDefaultI18nConfig() I18nConfig {
	return I18nConfig{
		Language: "en",
		Locale:   "en_US",
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
