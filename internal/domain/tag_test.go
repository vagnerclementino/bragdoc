package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTag_Validate(t *testing.T) {
	tests := []struct {
		name        string
		tag         *Tag
		expectError bool
		errorMsg    string
	}{
		{
			name:        "should return error when tag is nil",
			tag:         nil,
			expectError: true,
			errorMsg:    "tag cannot be nil",
		},
		{
			name: "should return error when name is empty",
			tag: &Tag{
				ID:      1,
				Name:    "",
				OwnerID: 1,
			},
			expectError: true,
			errorMsg:    "tag name cannot be empty",
		},
		{
			name: "should return error when name is too short",
			tag: &Tag{
				ID:      1,
				Name:    "a",
				OwnerID: 1,
			},
			expectError: true,
			errorMsg:    "tag name must be at least 2 characters",
		},
		{
			name: "should return error when name is too long",
			tag: &Tag{
				ID:      1,
				Name:    "this-is-a-very-long-tag-name",
				OwnerID: 1,
			},
			expectError: true,
			errorMsg:    "tag name cannot exceed 20 characters",
		},
		{
			name: "should pass validation with valid data",
			tag: &Tag{
				ID:      1,
				Name:    "golang",
				OwnerID: 1,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tag.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTag_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		tag      *Tag
		expected bool
	}{
		{
			name:     "should return false when tag is nil",
			tag:      nil,
			expected: false,
		},
		{
			name: "should return false when name is empty",
			tag: &Tag{
				ID:      1,
				Name:    "",
				OwnerID: 1,
			},
			expected: false,
		},
		{
			name: "should return true when name is present",
			tag: &Tag{
				ID:      1,
				Name:    "golang",
				OwnerID: 1,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tag.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}
