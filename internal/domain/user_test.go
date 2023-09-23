package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_validate(t *testing.T) {
	tests := []struct {
		name     string
		scenario func(t *testing.T)
	}{
		{
			name: "should returns a error when the use's name was not provided",
			scenario: func(t *testing.T) {

				u := &User{
					ID:    "952050e6-7f8c-45f9-a7d3-1eca6bcd9fe6",
					Email: "user@test.com",
				}

				err := u.Validate()

				assert.EqualError(t, err, "User.Name: the user's name cannot be empty")
			},
		},
		{
			name: "should returns a error when the use's name is empty",
			scenario: func(t *testing.T) {

				u := &User{
					ID:    "30f71dfe-b569-4e8f-879d-53b5df73929a",
					Name:  "",
					Email: "user@test.com",
				}

				err := u.Validate()

				assert.EqualError(t, err, "User.Name: the user's name cannot be empty")
			},
		},
		{
			name: "should returns a error when the use's name is less or equal 3 characters",
			scenario: func(t *testing.T) {

				u := &User{
					ID:    "30f71dfe-b569-4e8f-879d-53b5df73929a",
					Name:  "joe",
					Email: "user@test.com",
				}

				err := u.Validate()

				assert.EqualError(t, err, "User.Name: the user's name has an unexpected size")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.scenario(t)
		})
	}
}
