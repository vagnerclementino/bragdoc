package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid configuration",
			config: &Config{
				User: UserConfig{
					Name:   "John Doe",
					Email:  "john@example.com",
					Locale: "en-US",
				},
				Database: DatabaseConfig{
					Path: "/path/to/db",
				},
			},
			wantErr: false,
		},
		{
			name: "missing user name",
			config: &Config{
				User: UserConfig{
					Email: "john@example.com",
				},
				Database: DatabaseConfig{
					Path: "/path/to/db",
				},
			},
			wantErr: true,
			errMsg:  "user.name",
		},
		{
			name: "missing user email",
			config: &Config{
				User: UserConfig{
					Name: "John Doe",
				},
				Database: DatabaseConfig{
					Path: "/path/to/db",
				},
			},
			wantErr: true,
			errMsg:  "user.email",
		},
		{
			name: "missing database path",
			config: &Config{
				User: UserConfig{
					Name:  "John Doe",
					Email: "john@example.com",
				},
				Database: DatabaseConfig{},
			},
			wantErr: true,
			errMsg:  "database.path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetDefaultConfig(t *testing.T) {
	user := UserConfig{
		Name:     "Test User",
		Email:    "test@example.com",
		JobTitle: "Developer",
		Company:  "Test Corp",
		Locale:   "en-US",
	}

	config := GetDefaultConfig(user, "/test/dir")

	assert.Equal(t, user.Name, config.User.Name)
	assert.Equal(t, user.Email, config.User.Email)
	assert.Equal(t, user.JobTitle, config.User.JobTitle)
	assert.Equal(t, user.Company, config.User.Company)
	assert.Equal(t, user.Locale, config.User.Locale)

	assert.Equal(t, "/test/dir/bragdoc.db", config.Database.Path)
	assert.Equal(t, "openai", config.AI.Provider)
	assert.Equal(t, "gpt-4", config.AI.Model)
	assert.Equal(t, 2000, config.AI.MaxTokens)
	assert.Equal(t, 8080, config.Server.Port)
	assert.Equal(t, "info", config.Logging.Level)
	assert.Equal(t, "en", config.I18n.Language)
}

func TestGetDefaultPrompts(t *testing.T) {
	prompts := GetDefaultPrompts()

	assert.NotEmpty(t, prompts.EnhanceDescription)
	assert.NotEmpty(t, prompts.GenerateDocument)
	assert.NotEmpty(t, prompts.SuggestTags)
	assert.NotEmpty(t, prompts.TranslateBrag)

	// Check that prompts contain expected placeholders
	assert.Contains(t, prompts.EnhanceDescription, "{{.Description}}")
	assert.Contains(t, prompts.GenerateDocument, "{{range .Brags}}")
	assert.Contains(t, prompts.SuggestTags, "{{.Title}}")
	assert.Contains(t, prompts.TranslateBrag, "{{.TargetLanguage}}")
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ConfigFormat
		wantErr bool
	}{
		{
			name:    "yaml format",
			input:   "yaml",
			want:    FormatYAML,
			wantErr: false,
		},
		{
			name:    "yml format",
			input:   "yml",
			want:    FormatYAML,
			wantErr: false,
		},
		{
			name:    "empty defaults to yaml",
			input:   "",
			want:    FormatYAML,
			wantErr: false,
		},
		{
			name:    "json format",
			input:   "json",
			want:    FormatJSON,
			wantErr: false,
		},
		{
			name:    "toml format",
			input:   "toml",
			want:    FormatTOML,
			wantErr: false,
		},
		{
			name:    "unsupported format",
			input:   "xml",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				assert.IsType(t, ErrUnsupportedFormat{}, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestConfigFormat_Extension(t *testing.T) {
	tests := []struct {
		format ConfigFormat
		want   string
	}{
		{FormatYAML, ".yaml"},
		{FormatJSON, ".json"},
		{FormatTOML, ".toml"},
		{ConfigFormat("unknown"), ".yaml"},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			got := tt.format.Extension()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfigFormat_String(t *testing.T) {
	tests := []struct {
		format ConfigFormat
		want   string
	}{
		{FormatYAML, "yaml"},
		{FormatJSON, "json"},
		{FormatTOML, "toml"},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			got := tt.format.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
