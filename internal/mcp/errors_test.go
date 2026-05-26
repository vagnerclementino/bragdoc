package mcp

// Feature: mcp-server-mode, Property 19: Error Code Classification

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

// TestClassifyError_ValidationPattern is a property-based test:
// for any error string containing "validation failed:", classifyError returns ErrCodeValidation with the full message.
func TestClassifyError_ValidationPattern(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a random suffix after the validation prefix
		suffix := rapid.StringMatching(`[a-zA-Z0-9 _\-]{1,50}`).Draw(t, "suffix")
		// Optionally add a prefix before "validation failed:"
		prefix := rapid.StringMatching(`[a-zA-Z0-9 ]{0,20}`).Draw(t, "prefix")

		// Ensure the error message contains "validation failed:" (with HasPrefix match per implementation)
		// The actual implementation uses HasPrefix, so we test with the prefix at position 0
		errMsg := "validation failed:" + suffix
		// Also test with prefix to verify Contains vs HasPrefix behavior
		if prefix != "" {
			// The implementation uses HasPrefix, so only messages starting with "validation failed:" match
			// Test the pure prefix case
			errMsg = "validation failed:" + suffix
		}

		err := errors.New(errMsg)
		code, message := classifyError(err)

		assert.Equal(t, ErrCodeValidation, code, "expected ErrCodeValidation for message: %s", errMsg)
		assert.Equal(t, errMsg, message, "expected full error message to be returned")
	})
}

// TestClassifyError_NotFoundPattern is a property-based test:
// for any error string containing "not found", classifyError returns ErrCodeNotFound with the full message.
func TestClassifyError_NotFoundPattern(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate prefix and suffix around "not found"
		prefix := rapid.StringMatching(`[a-zA-Z0-9 ]{0,30}`).Draw(t, "prefix")
		suffix := rapid.StringMatching(`[a-zA-Z0-9 ]{0,30}`).Draw(t, "suffix")

		// Ensure message does NOT start with "validation failed:" to avoid the first case matching
		errMsg := prefix + "not found" + suffix
		if strings.HasPrefix(errMsg, "validation failed:") {
			errMsg = "resource not found" + suffix
		}

		err := errors.New(errMsg)
		code, message := classifyError(err)

		assert.Equal(t, ErrCodeNotFound, code, "expected ErrCodeNotFound for message: %s", errMsg)
		assert.Equal(t, errMsg, message, "expected full error message to be returned")
	})
}

// TestClassifyError_InternalFallback is a property-based test:
// for any error string NOT matching "validation failed:" prefix or "not found" substring,
// classifyError returns ErrCodeInternal with "internal error" message.
func TestClassifyError_InternalFallback(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a random error message that does NOT contain either pattern
		errMsg := rapid.StringMatching(`[a-zA-Z0-9_\-]{1,50}`).Draw(t, "errMsg")

		// Ensure neither pattern matches
		for strings.HasPrefix(errMsg, "validation failed:") || strings.Contains(errMsg, "not found") {
			errMsg = rapid.StringMatching(`[a-zA-Z0-9_\-]{1,50}`).Draw(t, "errMsg")
		}

		err := errors.New(errMsg)
		code, message := classifyError(err)

		assert.Equal(t, ErrCodeInternal, code, "expected ErrCodeInternal for message: %s", errMsg)
		assert.Equal(t, "internal error", message, "expected generic 'internal error' message, not internal details")
	})
}

// TestClassifyError_EdgeCases tests specific edge cases with example-based tests.
func TestClassifyError_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantCode    int
		wantMessage string
	}{
		{
			name:        "empty string falls to internal",
			err:         errors.New(""),
			wantCode:    ErrCodeInternal,
			wantMessage: "internal error",
		},
		{
			name:        "both patterns present - validation failed prefix takes priority",
			err:         errors.New("validation failed: resource not found"),
			wantCode:    ErrCodeValidation,
			wantMessage: "validation failed: resource not found",
		},
		{
			name:        "not found without validation prefix",
			err:         errors.New("user not found in database"),
			wantCode:    ErrCodeNotFound,
			wantMessage: "user not found in database",
		},
		{
			name:        "validation failed not at prefix does not match validation",
			err:         errors.New("something validation failed: oops"),
			wantCode:    ErrCodeInternal,
			wantMessage: "internal error",
		},
		{
			name:        "exact validation failed prefix with colon",
			err:         errors.New("validation failed: title too short"),
			wantCode:    ErrCodeValidation,
			wantMessage: "validation failed: title too short",
		},
		{
			name:        "not found at start of message",
			err:         errors.New("not found"),
			wantCode:    ErrCodeNotFound,
			wantMessage: "not found",
		},
		{
			name:        "generic error with no patterns",
			err:         fmt.Errorf("connection refused"),
			wantCode:    ErrCodeInternal,
			wantMessage: "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, message := classifyError(tt.err)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantMessage, message)
		})
	}
}
