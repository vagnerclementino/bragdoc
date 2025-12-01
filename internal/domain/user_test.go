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
			name: "should return error when language code is invalid",
			user: &User{
				ID:       1,
				Name:     "John Doe",
				Email:    "john@test.com",
				Language: "english",
			},
			expectError: true,
			errorMsg:    "language must be a 2-letter ISO 639-1 code",
		},
		{
			name: "should pass validation with valid data",
			user: &User{
				ID:       1,
				Name:     "John Doe",
				Email:    "john@test.com",
				Language: "en",
			},
			expectError: false,
		},
		{
			name: "should default to 'en' when language is empty",
			user: &User{
				ID:       1,
				Name:     "John Doe",
				Email:    "john@test.com",
				Language: "",
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
				if tt.user != nil && tt.user.Language == "" {
					// After validation, empty language should be set to "en"
					assert.Equal(t, "en", tt.user.Language)
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
