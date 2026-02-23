package domain

import "time"

// Position represents a user's role or job at a specific time.
type Position struct {
    ID        int64
    User      User
    Title     string
    Company   string
    StartDate *time.Time
    EndDate   *time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
}
