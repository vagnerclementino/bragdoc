package domain

import (
	"time"
)

// User represents a user of the system
// This is a pure data structure with no validation logic
type User struct {
	ID        int64
	Name      string
	Email     string
	JobTitle  string
	Company   string
	Locale    Locale // Locale in format language-COUNTRY (e.g., en-US, pt-BR, pt-PT)
	CreatedAt time.Time
	UpdatedAt time.Time
}
