package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name        string
		user        *User
		expectError bool
		errorMsg    string
	}{
		{
			name:        "should return error when user is nil",
			user:        nil,
			expectError: true,
			errorMsg:    "user cannot be nil",
		},
		{
			name: "should return error when name is empty",
			user: &User{
				ID:    1,
				Name:  "",
				Email: "user@test.com",
			},
			expectError: true,
			errorMsg:    "user name cannot be empty",
		},
		{
			name: "should return error when name is too short",
			user: &User{
				ID:    1,
				Name:  "Jo",
				Email: "user@test.com",
			},
			expectError: true,
			errorMsg:    "user name must be at least 3 characters",
		},
		{
			name: "should return error when email is empty",
			user: &User{
				ID:    1,
				Name:  "John Doe",
				Email: "",
			},
			expectError: true,
			errorMsg:    "user email cannot be empty",
		},
		{
			name: "should return error when email is invalid",
			user: &User{
				ID:    1,
				Name:  "John Doe",
				Email: "@test.com",
			},
			expectError: true,
			errorMsg:    "invalid email address",
		},
		{
			name: "should return error when locale is invalid",
			user: &User{
				ID:     1,
				Name:   "John Doe",
				Email:  "john@test.com",
				Locale: Locale("de-DE"),
			},
			expectError: true,
			errorMsg:    "unsupported locale",
		},
		{
			name: "should pass validation with valid data",
			user: &User{
				ID:     1,
				Name:   "John Doe",
				Email:  "john@test.com",
				Locale: LocaleEnglishUS,
			},
			expectError: false,
		},
		{
			name: "should default to 'en-US' when locale is empty",
			user: &User{
				ID:     1,
				Name:   "John Doe",
				Email:  "john@test.com",
				Locale: "",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				if tt.user != nil && tt.user.Locale == "" {
					// After validation, empty locale should be set to "en-US"
					assert.Equal(t, LocaleEnglishUS, tt.user.Locale)
				}
			}
		})
	}
}

func TestUser_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		expected bool
	}{
		{
			name:     "should return false when user is nil",
			user:     nil,
			expected: false,
		},
		{
			name: "should return false when name is empty",
			user: &User{
				ID:    1,
				Name:  "",
				Email: "user@test.com",
			},
			expected: false,
		},
		{
			name: "should return false when email is empty",
			user: &User{
				ID:    1,
				Name:  "John Doe",
				Email: "",
			},
			expected: false,
		},
		{
			name: "should return true when both name and email are present",
			user: &User{
				ID:    1,
				Name:  "John Doe",
				Email: "john@test.com",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}
