// Package domain provides
package domain

import (
	"errors"
	"time"
)

type Document struct {
	ID string
	User
	Brags     []Brag
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func (d *Document) Validate() error {

	if len(d.Brags) == 0 {
		return errors.New("Document.Brags: the document's brag list cannot be empty. Please provide at least one brag")

	}
	return nil
}
