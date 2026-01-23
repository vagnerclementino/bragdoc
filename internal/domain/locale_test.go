package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocale_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		locale Locale
		valid  bool
	}{
		{"en-US is valid", LocaleEnglishUS, true},
		{"pt-BR is valid", LocalePortugueseBR, true},
		{"en-GB is invalid", Locale("en-GB"), false},
		{"pt-PT is invalid", Locale("pt-PT"), false},
		{"es-ES is invalid", Locale("es-ES"), false},
		{"de-DE is invalid", Locale("de-DE"), false},
		{"empty is invalid", Locale(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.locale.IsValid())
		})
	}
}

func TestLocale_GetLanguageCode(t *testing.T) {
	tests := []struct {
		locale Locale
		want   string
	}{
		{LocaleEnglishUS, "en"},
		{LocalePortugueseBR, "pt"},
	}

	for _, tt := range tests {
		t.Run(string(tt.locale), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.locale.GetLanguageCode())
		})
	}
}

func TestLocale_GetCountryCode(t *testing.T) {
	tests := []struct {
		locale Locale
		want   string
	}{
		{LocaleEnglishUS, "US"},
		{LocalePortugueseBR, "BR"},
	}

	for _, tt := range tests {
		t.Run(string(tt.locale), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.locale.GetCountryCode())
		})
	}
}

func TestLocale_GetLocaleName(t *testing.T) {
	tests := []struct {
		locale Locale
		name   string
	}{
		{LocaleEnglishUS, "English (United States)"},
		{LocalePortugueseBR, "Portuguese (Brazil)"},
		{Locale("de-DE"), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.locale), func(t *testing.T) {
			assert.Equal(t, tt.name, tt.locale.GetLocaleName())
		})
	}
}

func TestParseLocale(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Locale
		wantErr bool
	}{
		{
			name:    "valid en-US",
			input:   "en-US",
			want:    LocaleEnglishUS,
			wantErr: false,
		},
		{
			name:    "valid lowercase",
			input:   "pt-br",
			want:    LocalePortugueseBR,
			wantErr: false,
		},
		{
			name:    "valid with spaces",
			input:   " en-US ",
			want:    LocaleEnglishUS,
			wantErr: false,
		},
		{
			name:    "invalid locale pt-PT",
			input:   "pt-PT",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid locale de-DE",
			input:   "de-DE",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLocale(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "unsupported locale")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSupportedLocales(t *testing.T) {
	locales := SupportedLocales()

	assert.Len(t, locales, 2)
	assert.Contains(t, locales, LocaleEnglishUS)
	assert.Contains(t, locales, LocalePortugueseBR)

	// All locales should be valid
	for _, locale := range locales {
		assert.True(t, locale.IsValid(), "locale %s should be valid", locale)
	}
}

func TestSupportedLocalesString(t *testing.T) {
	str := SupportedLocalesString()

	assert.Equal(t, "en-US, pt-BR", str)
	assert.Contains(t, str, "en-US")
	assert.Contains(t, str, "pt-BR")
}

func TestLocale_String(t *testing.T) {
	assert.Equal(t, "en-US", LocaleEnglishUS.String())
	assert.Equal(t, "pt-BR", LocalePortugueseBR.String())
}
