package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
