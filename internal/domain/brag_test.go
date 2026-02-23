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
        {"project", "project", Category{Name: CategoryNameProject, Description: "PROJECT DELIVERABLES"}, false},
        {"achievement", "achievement", Category{Name: CategoryNameAchievement, Description: "MEASURABLE ACHIEVEMENTS"}, false},
        {"skill", "skill", Category{Name: CategoryNameSkill, Description: "SKILLS AND LEARNING"}, false},
        {"leadership", "leadership", Category{Name: CategoryNameLeadership, Description: "TEAM OR LEADERSHIP ACTS"}, false},
        {"innovation", "innovation", Category{Name: CategoryNameInnovation, Description: "INNOVATIONS AND IMPROVEMENTS"}, false},
        {"uppercase", "PROJECT", Category{Name: CategoryNameProject, Description: "PROJECT DELIVERABLES"}, false},
        {"mixed case", "AcHiEvEmEnT", Category{Name: CategoryNameAchievement, Description: "MEASURABLE ACHIEVEMENTS"}, false},
        {"with spaces", "  project  ", Category{Name: CategoryNameProject, Description: "PROJECT DELIVERABLES"}, false},
        {"invalid", "invalid", Category{}, true},
        {"empty", "", Category{}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := ParseCategory(tt.input)
            if tt.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected.Name, result.Name)
                assert.Equal(t, tt.expected.Description, result.Description)
            }
        })
    }
}

func TestCategory_String(t *testing.T) {
    tests := []struct {
        category Category
        expected string
    }{
        {Category{Name: CategoryNameProject}, "PROJECT"},
        {Category{Name: CategoryNameAchievement}, "ACHIEVEMENT"},
        {Category{Name: CategoryNameSkill}, "SKILL"},
        {Category{Name: CategoryNameLeadership}, "LEADERSHIP"},
        {Category{Name: CategoryNameInnovation}, "INNOVATION"},
        {Category{Name: CategoryNameUnknown}, "UNKNOWN"},
        {Category{}, "UNKNOWN"},
    }

    for _, tt := range tests {
        t.Run(tt.expected, func(t *testing.T) {
            result := tt.category.String()
            assert.Equal(t, tt.expected, result)
        })
    }
}
