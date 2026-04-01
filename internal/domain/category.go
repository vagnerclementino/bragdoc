package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type CategoryName string

const (
	CategoryNameUnknown     CategoryName = "UNKNOWN"
	CategoryNameProject     CategoryName = "PROJECT"
	CategoryNameAchievement CategoryName = "ACHIEVEMENT"
	CategoryNameSkill       CategoryName = "SKILL"
	CategoryNameLeadership  CategoryName = "LEADERSHIP"
	CategoryNameInnovation  CategoryName = "INNOVATION"
	CategoryNameDelivery    CategoryName = "DELIVERY"
)

var defaultCategoryDescriptions = map[CategoryName]string{
	CategoryNameUnknown:     "GENERAL CATEGORY",
	CategoryNameProject:     "PROJECT DELIVERABLES",
	CategoryNameAchievement: "MEASURABLE ACHIEVEMENTS",
	CategoryNameSkill:       "SKILLS AND LEARNING",
	CategoryNameLeadership:  "TEAM OR LEADERSHIP ACTS",
	CategoryNameInnovation:  "INNOVATIONS AND IMPROVEMENTS",
	CategoryNameDelivery:    "VALUE DELIVERIES",
}

type Category struct {
	ID          int64
	Name        CategoryName
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (c Category) String() string {
	if c.Name == "" {
		return string(CategoryNameUnknown)
	}
	return string(c.Name)
}

func (c Category) Validate() error {
	if c.Name == "" {
		return errors.New("category name is required")
	}
	if len(c.Name) > 20 {
		return fmt.Errorf("category name must be at most 20 characters, got %d", len(c.Name))
	}
	if string(c.Name) != strings.ToUpper(string(c.Name)) {
		return errors.New("category name must be uppercase")
	}
	return nil
}

func NormalizeCategoryName(value string) CategoryName {
	return CategoryName(strings.ToUpper(strings.TrimSpace(value)))
}

func ParseCategory(input string) (Category, error) {
	name := NormalizeCategoryName(input)
	desc, ok := defaultCategoryDescriptions[name]
	if !ok {
		return Category{}, fmt.Errorf("invalid category: %s", input)
	}
	return Category{Name: name, Description: desc}, nil
}

func DefaultCategories() []Category {
	categories := make([]Category, 0, len(defaultCategoryDescriptions))
	for name, desc := range defaultCategoryDescriptions {
		categories = append(categories, Category{Name: name, Description: desc})
	}
	return categories
}
