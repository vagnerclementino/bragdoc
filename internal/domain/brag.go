package domain

import (
	"fmt"
	"strings"
	"time"
)

// Category represents the category of a brag
type Category int

const (
	CategoryUnknown Category = iota
	CategoryProject
	CategoryAchievement
	CategorySkill
	CategoryLeadership
	CategoryInnovation
)

var categoryStrings = map[Category]string{
	CategoryUnknown:     "unknown",
	CategoryProject:     "project",
	CategoryAchievement: "achievement",
	CategorySkill:       "skill",
	CategoryLeadership:  "leadership",
	CategoryInnovation:  "innovation",
}

var stringToCategory = map[string]Category{
	"unknown":     CategoryUnknown,
	"project":     CategoryProject,
	"achievement": CategoryAchievement,
	"skill":       CategorySkill,
	"leadership":  CategoryLeadership,
	"innovation":  CategoryInnovation,
}

// String returns the string representation of the category
func (c Category) String() string {
	if str, ok := categoryStrings[c]; ok {
		return str
	}
	return "unknown"
}

// ParseCategory parses a string into a Category
func ParseCategory(s string) (Category, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if cat, ok := stringToCategory[s]; ok {
		return cat, nil
	}
	return CategoryUnknown, fmt.Errorf("invalid category: %s", s)
}

// Brag represents a professional achievement or accomplishment
// This is a pure data structure with no validation logic
type Brag struct {
	ID          int64
	OwnerID     int64
	Title       string
	Description string
	Category    Category
	Tags        []*Tag
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
