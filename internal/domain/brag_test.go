package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrag_Validate(t *testing.T) {
	tests := []struct {
		name        string
		brag        *Brag
		expectError bool
		errorMsg    string
	}{
		{
			name:        "should return error when brag is nil",
			brag:        nil,
			expectError: true,
			errorMsg:    "brag cannot be nil",
		},
		{
			name: "should return error when title is empty",
			brag: &Brag{
				Title:       "",
				Description: "This is a valid description with sufficient length",
				Category:    CategoryProject,
			},
			expectError: true,
			errorMsg:    "brag title cannot be empty",
		},
		{
			name: "should return error when title is too short",
			brag: &Brag{
				Title:       "Hi",
				Description: "This is a valid description with sufficient length",
				Category:    CategoryProject,
			},
			expectError: true,
			errorMsg:    "brag title must be at least 5 characters",
		},
		{
			name: "should return error when description is empty",
			brag: &Brag{
				Title:       "Valid Title",
				Description: "",
				Category:    CategoryProject,
			},
			expectError: true,
			errorMsg:    "brag description cannot be empty",
		},
		{
			name: "should return error when description is too short",
			brag: &Brag{
				Title:       "Valid Title",
				Description: "Short",
				Category:    CategoryProject,
			},
			expectError: true,
			errorMsg:    "brag description must be at least 20 characters",
		},
		{
			name: "should return error when category is invalid",
			brag: &Brag{
				Title:       "Valid Title",
				Description: "This is a valid description with sufficient length",
				Category:    Category(999),
			},
			expectError: true,
			errorMsg:    "invalid brag category",
		},
		{
			name: "should pass validation with valid data",
			brag: &Brag{
				Title:       "Valid Title",
				Description: "This is a valid description with sufficient length",
				Category:    CategoryProject,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.brag.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBrag_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		brag     *Brag
		expected bool
	}{
		{
			name:     "should return false when brag is nil",
			brag:     nil,
			expected: false,
		},
		{
			name: "should return false when title is empty",
			brag: &Brag{
				Title:       "",
				Description: "Valid description",
			},
			expected: false,
		},
		{
			name: "should return false when description is empty",
			brag: &Brag{
				Title:       "Valid Title",
				Description: "",
			},
			expected: false,
		},
		{
			name: "should return true when both title and description are present",
			brag: &Brag{
				Title:       "Valid Title",
				Description: "Valid description",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.brag.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseCategory(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    Category
		expectError bool
	}{
		{"project", "project", CategoryProject, false},
		{"achievement", "achievement", CategoryAchievement, false},
		{"skill", "skill", CategorySkill, false},
		{"leadership", "leadership", CategoryLeadership, false},
		{"innovation", "innovation", CategoryInnovation, false},
		{"uppercase", "PROJECT", CategoryProject, false},
		{"mixed case", "AcHiEvEmEnT", CategoryAchievement, false},
		{"with spaces", "  project  ", CategoryProject, false},
		{"invalid", "invalid", CategoryUnknown, true},
		{"empty", "", CategoryUnknown, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseCategory(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestCategory_String(t *testing.T) {
	tests := []struct {
		category Category
		expected string
	}{
		{CategoryProject, "project"},
		{CategoryAchievement, "achievement"},
		{CategorySkill, "skill"},
		{CategoryLeadership, "leadership"},
		{CategoryInnovation, "innovation"},
		{CategoryUnknown, "unknown"},
		{Category(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.category.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}
