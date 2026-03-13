package domain

import (
	"time"
)

// Tag represents a label for categorizing brags
// This is a pure data structure with no validation logic
type Tag struct {
	ID        int64
	Name      string
	OwnerID   int64
	CreatedAt time.Time
}
