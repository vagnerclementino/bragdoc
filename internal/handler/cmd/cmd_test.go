package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdHandler_register(t *testing.T) {
	tests := []struct {
		name     string
		scenario func(t *testing.T)
	}{
		{
			name: "should returns an error to a unknown command",
			scenario: func(t *testing.T) {
				_, err := NewCli([]string{"xyz"})
				assert.EqualError(t, err, "the command 'xyz' cannot be registered")

			},
		},
		{
			name: "should returns no  error with a valid command",
			scenario: func(t *testing.T) {
				_, err := NewCli([]string{"version"})
				assert.NoError(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.scenario(t)
		})
	}
}
