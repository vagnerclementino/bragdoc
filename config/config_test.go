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
				Database: DatabaseConfig{
					Path: "/path/to/db",
				},
			},
			wantErr: false,
		},
		{
			name: "missing database path",
			config: &Config{
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
	config := GetDefaultConfig("/test/dir")

	assert.Equal(t, "/test/dir/bragdoc.db", config.Database.Path)
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
