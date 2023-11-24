// Package  provides ...
package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocument_Validate(t *testing.T) {
	tests := []struct {
		name     string
		scenario func(t *testing.T)
	}{
		{
			name: "should returns a error if a document is empty",
			scenario: func(t *testing.T) {

				d := Document{}

				err := d.Validate()

				assert.EqualError(t, err, "Document.Brags: the document's brag list cannot be empty. Please provide at least one brag")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.scenario(t)

		})
	}
}
