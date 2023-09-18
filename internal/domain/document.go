// Package domain provides
package domain

import "time"

type Document struct {
	ID string
	User
	Brags     []Brag
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func (d *Document) Validate() error {

	return nil
}
