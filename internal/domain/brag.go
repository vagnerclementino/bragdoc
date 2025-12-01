package domain

import (
	"errors"
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

// IsValid performs basic structural validation
func (b *Brag) IsValid() bool {
	return b != nil &&
		strings.TrimSpace(b.Title) != "" &&
		strings.TrimSpace(b.Description) != ""
}

// Validate performs comprehensive validation with detailed error messages
func (b *Brag) Validate() error {
	if b == nil {
		return errors.New("brag cannot be nil")
	}

	if strings.TrimSpace(b.Title) == "" {
		return errors.New("brag title cannot be empty")
	}

	if len(strings.TrimSpace(b.Title)) < 5 {
		return fmt.Errorf("brag title must be at least 5 characters, got %d", len(strings.TrimSpace(b.Title)))
	}

	if strings.TrimSpace(b.Description) == "" {
		return errors.New("brag description cannot be empty")
	}

	if len(strings.TrimSpace(b.Description)) < 20 {
		return fmt.Errorf("brag description must be at least 20 characters, got %d", len(strings.TrimSpace(b.Description)))
	}

	if b.Category < CategoryProject || b.Category > CategoryInnovation {
		return fmt.Errorf("invalid brag category: %d", b.Category)
	}

	return nil
}
