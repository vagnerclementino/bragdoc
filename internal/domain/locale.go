package domain

import (
	"fmt"
	"strings"
)

// Locale represents a supported locale (language-COUNTRY format)
// Format: ISO 639-1 (language) + ISO 3166-1 alpha-2 (country)
type Locale string

// Supported locales (v1 supports only en-US and pt-BR)
const (
	LocaleEnglishUS    Locale = "en-US" // English (United States)
	LocalePortugueseBR Locale = "pt-BR" // Portuguese (Brazil)
)

// String returns the string representation of the locale
func (l Locale) String() string {
	return string(l)
}

// IsValid checks if the locale is supported
func (l Locale) IsValid() bool {
	switch l {
	case LocaleEnglishUS, LocalePortugueseBR:
		return true
	default:
		return false
	}
}

// GetLanguageCode returns the language part of the locale (e.g., "en" from "en-US")
func (l Locale) GetLanguageCode() string {
	parts := strings.Split(string(l), "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return "en"
}

// GetCountryCode returns the country part of the locale (e.g., "US" from "en-US")
func (l Locale) GetCountryCode() string {
	parts := strings.Split(string(l), "-")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

// GetLocaleName returns the full locale name
func (l Locale) GetLocaleName() string {
	switch l {
	case LocaleEnglishUS:
		return "English (United States)"
	case LocalePortugueseBR:
		return "Portuguese (Brazil)"
	default:
		return "Unknown"
	}
}

// ParseLocale parses a string into a Locale
func ParseLocale(s string) (Locale, error) {
	// Normalize: convert to lowercase, trim spaces
	normalized := strings.ToLower(strings.TrimSpace(s))
	
	// Try to match with supported locales (case-insensitive)
	for _, locale := range SupportedLocales() {
		if strings.ToLower(string(locale)) == normalized {
			return locale, nil
		}
	}
	
	return "", fmt.Errorf("unsupported locale: %s (supported: %s)", s, SupportedLocalesString())
}

// SupportedLocales returns a list of all supported locales
func SupportedLocales() []Locale {
	return []Locale{
		LocaleEnglishUS,
		LocalePortugueseBR,
	}
}

// SupportedLocalesString returns a comma-separated string of supported locales
func SupportedLocalesString() string {
	locales := SupportedLocales()
	strs := make([]string, len(locales))
	for i, locale := range locales {
		strs[i] = string(locale)
	}
	return strings.Join(strs, ", ")
}
