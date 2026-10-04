package mcp

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"pgregory.net/rapid"
)

func TestClassifyError_WrappedCategories(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		message := rapid.String().Draw(t, "message")
		for _, test := range []struct {
			err  error
			code int
		}{
			{domain.ValidationError(errors.New(message)), ErrCodeValidation},
			{domain.NotFoundError(errors.New(message)), ErrCodeNotFound},
		} {
			wrapped := fmt.Errorf("operation failed: %w", test.err)
			code, safeMessage := classifyError(wrapped)
			assert.Equal(t, test.code, code)
			assert.Equal(t, wrapped.Error(), safeMessage)
		}
	})
}

func TestClassifyError_InternalMessagesAreHidden(t *testing.T) {
	for _, message := range []string{"", "database password secret", "validation failed: database corruption", "database file not found"} {
		code, safeMessage := classifyError(errors.New(message))
		assert.Equal(t, ErrCodeInternal, code)
		assert.Equal(t, "internal error", safeMessage)
	}
}

func TestToolError_StructuredAndTextContentAgree(t *testing.T) {
	result := toolError(domain.ValidationError(errors.New("at least one tag name is required")))
	assert.True(t, result.IsError)
	assert.Equal(t, ErrorResponse{Code: ErrCodeValidation, Message: "at least one tag name is required"}, result.StructuredContent)
	var response ErrorResponse
	assert.NoError(t, unmarshalResult(result, &response))
	assert.Equal(t, result.StructuredContent, response)
}
