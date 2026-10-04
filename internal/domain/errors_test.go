package domain

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClassifiedErrorsPreserveCauseAndMessage(t *testing.T) {
	cause := fmt.Errorf("user not found: %w", sql.ErrNoRows)
	err := fmt.Errorf("lookup failed: %w", NotFoundError(cause))
	assert.ErrorIs(t, err, ErrNotFound)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Equal(t, "lookup failed: "+cause.Error(), err.Error())
	assert.False(t, errors.Is(err, ErrValidation))
}
