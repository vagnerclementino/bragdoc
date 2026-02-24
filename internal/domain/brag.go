package domain

import "time"

// Brag represents a professional achievement or accomplishment.
type Brag struct {
    ID          int64
    Owner       User
    JobTitle    *JobTitle
    Category    Category
    Title       string
    Description string
    Tags        []*Tag
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
