// Package domain provides ...
package domain

import (
	"errors"
	"fmt"
	"strings"
)

type Brag struct {
	ID          string
	Description string
	Details     *string
	CreatedAt   int64  // Unix timestamp
	UpdatedAt   *int64 // Unix timestamp pointer
}

func (b *Brag) Validate() error {

	if strings.TrimSpace(b.Description) == "" {
		return errors.New("Brag.Description: the brag's description cannot be empty")
	}

	if len(strings.TrimSpace(b.Description)) < 10 {
		return fmt.Errorf("Brag.Description: the brag's description is very short. Please provide a text with a minimum size of %d.", 10)

	}

	if b.Details != nil && len(strings.TrimSpace(*b.Details)) < 20 {
		return fmt.Errorf("Brag.Details: the brag's details is very short. Please provide a text with a minimum size of %d.", 20)
	}

	return nil
}
