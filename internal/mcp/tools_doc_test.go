package mcp

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"pgregory.net/rapid"
)

// Feature: mcp-server-mode, Property 14: Unsupported Format/Template Error
//
// **Validates: Requirements 4.2, 4.4**

func TestHandleDocGenerate_Property_UnsupportedFormat(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		fix := newTestServer()
		ctx := context.Background()

		// Generate a non-empty string that is NOT "markdown"
		format := rapid.StringMatching(`[a-zA-Z]{1,20}`).
			Filter(func(s string) bool { return s != "markdown" }).
			Draw(t, "format")

		result, _, err := fix.server.handleDocGenerate(ctx, nil, DocGenerateParams{
			UserID: 1,
			Format: format,
		})
		if err != nil {
			t.Fatalf("handleDocGenerate returned unexpected error: %v", err)
		}

		// Property: unsupported format always returns a validation error
		assert.True(t, result.IsError)
		assert.Contains(t, extractText(result), "validation failed")
		assert.Contains(t, extractText(result), "markdown")
	})
}

func TestHandleDocGenerate_Property_UnsupportedTemplate(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		fix := newTestServer()
		ctx := context.Background()

		// Generate a non-empty string that is NOT "default"
		tmpl := rapid.StringMatching(`[a-zA-Z]{1,20}`).
			Filter(func(s string) bool { return s != "default" }).
			Draw(t, "template")

		result, _, err := fix.server.handleDocGenerate(ctx, nil, DocGenerateParams{
			UserID:   1,
			Format:   "markdown", // valid format so we reach template check
			Template: tmpl,
		})
		if err != nil {
			t.Fatalf("handleDocGenerate returned unexpected error: %v", err)
		}

		// Property: unsupported template always returns a validation error
		assert.True(t, result.IsError)
		assert.Contains(t, extractText(result), "validation failed")
		assert.Contains(t, extractText(result), "template")
	})
}

func TestHandleDocGenerate_NoBragsForUser(t *testing.T) {
	fix := newTestServer()
	ctx := context.Background()

	fix.bragRepo.On("SelectAll", mock.Anything, int64(42)).Return([]*domain.Brag{}, nil)

	result, _, err := fix.server.handleDocGenerate(ctx, nil, DocGenerateParams{
		UserID: 42,
		Format: "markdown",
	})
	assert.NoError(t, err)
	assert.True(t, result.IsError)
	assert.Contains(t, extractText(result), "no brags found")
}

func TestHandleDocGenerate_DefaultFormat(t *testing.T) {
	fix := newTestServer()
	ctx := context.Background()

	userID := int64(1)
	user := &domain.User{
		ID:        userID,
		Name:      "Alice",
		Email:     "alice@example.com",
		JobTitle:  "Engineer",
		Company:   "Acme",
		Locale:    domain.Locale("en-US"),
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}

	brags := []*domain.Brag{
		{
			ID:          1,
			Owner:       *user,
			Title:       "Shipped Feature X",
			Description: "Delivered a complex feature ahead of schedule with full test coverage",
			Category:    domain.Category{Name: domain.CategoryNameProject},
			Tags:        []*domain.Tag{},
			CreatedAt:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	fix.bragRepo.On("SelectAll", mock.Anything, userID).Return(brags, nil)
	fix.userRepo.On("Select", mock.Anything, userID).Return(user, nil)

	// Format omitted — should default to markdown and succeed
	result, _, err := fix.server.handleDocGenerate(ctx, nil, DocGenerateParams{
		UserID: userID,
	})
	assert.NoError(t, err)
	assert.False(t, result.IsError, "expected success but got error: %s", extractText(result))

	var resp DocResponse
	assert.NoError(t, unmarshalResult(result, &resp))
	assert.Equal(t, "Alice", resp.Author)
	assert.Equal(t, 1, resp.BragCount)
	assert.Contains(t, resp.Content, "Shipped Feature X")
}

func TestHandleDocGenerate_SelectAllError(t *testing.T) {
	fix := newTestServer()
	ctx := context.Background()

	fix.bragRepo.On("SelectAll", mock.Anything, int64(5)).Return(nil, fmt.Errorf("database error"))

	result, _, err := fix.server.handleDocGenerate(ctx, nil, DocGenerateParams{
		UserID: 5,
		Format: "markdown",
	})
	assert.NoError(t, err)
	assert.True(t, result.IsError)
}
