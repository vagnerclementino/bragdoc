package domain

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

// User represents a user of the system
type User struct {
	ID        int64
	Name      string
	Email     string
	JobTitle  string
	Company   string
	Locale    Locale // Locale in format language-COUNTRY (e.g., en-US, pt-BR, pt-PT)
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsValid performs basic structural validation
func (u *User) IsValid() bool {
	return u != nil &&
		strings.TrimSpace(u.Name) != "" &&
		strings.TrimSpace(u.Email) != ""
}

// Validate performs comprehensive validation with detailed error messages
func (u *User) Validate() error {
	if u == nil {
		return errors.New("user cannot be nil")
	}

	if strings.TrimSpace(u.Name) == "" {
		return errors.New("user name cannot be empty")
	}

	if len(strings.TrimSpace(u.Name)) < 3 {
		return fmt.Errorf("user name must be at least 3 characters, got %d", len(strings.TrimSpace(u.Name)))
	}

	if strings.TrimSpace(u.Email) == "" {
		return errors.New("user email cannot be empty")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return fmt.Errorf("invalid email address: %s", u.Email)
	}

	// Validate locale if provided
	if u.Locale != "" {
		if !u.Locale.IsValid() {
			return fmt.Errorf("unsupported locale: %s (supported: %s)", 
				u.Locale, SupportedLocalesString())
		}
	} else {
		u.Locale = LocaleEnglishUS // Default to en-US
	}

	return nil
}
