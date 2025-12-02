package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Tag represents a label for categorizing brags
type Tag struct {
	ID        int64
	Name      string
	OwnerID   int64
	CreatedAt time.Time
}

// IsValid performs basic structural validation
func (t *Tag) IsValid() bool {
	return t != nil && strings.TrimSpace(t.Name) != ""
}

// Validate performs comprehensive validation with detailed error messages
func (t *Tag) Validate() error {
	if t == nil {
		return errors.New("tag cannot be nil")
	}

	name := strings.TrimSpace(t.Name)
	if name == "" {
		return errors.New("tag name cannot be empty")
	}

	if len(name) < 2 {
		return fmt.Errorf("tag name must be at least 2 characters, got %d", len(name))
	}

	if len(name) > 20 {
		return fmt.Errorf("tag name cannot exceed 20 characters, got %d", len(name))
	}

	return nil
}
