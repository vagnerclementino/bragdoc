package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrag_validate(t *testing.T) {
	tests := []struct {
		name     string
		scenario func(t *testing.T)
	}{
		{
			name: "should returns an error when brag's description is empty",
			scenario: func(t *testing.T) {

				b := Brag{}
				err := b.Validate()
				assert.EqualError(t, err, "Brag.Description: the brag's description cannot be empty")
			},
		},
		{
			name: "should returns an error when brag's description is short",
			scenario: func(t *testing.T) {

				b := Brag{
					Description: "hi",
				}
				err := b.Validate()
				assert.EqualError(t, err, "Brag.Description: the brag's description is very short. Please provide a text with a minimum size of 10.")
			},
		},
		{
			name: "should returns an error when brag's details is short",
			scenario: func(t *testing.T) {

				description := "only a presentation"
				b := Brag{
					Description: "technical presentation for the team",
					Details:     &description,
				}
				err := b.Validate()
				assert.EqualError(t, err, "Brag.Details: the brag's details is very short. Please provide a text with a minimum size of 20.")

			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.scenario(t)

		})
	}
}
