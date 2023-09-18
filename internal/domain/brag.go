// Package domain provides ...
package domain

import (
	"time"
)

type Brag struct {
	ID         string
	descrption string
	details    *string
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

func (b *Brag) Validate() error {

	return nil
}
