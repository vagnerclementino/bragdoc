package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"pgregory.net/rapid"
)

// --- Helpers ---

// extractText gets the text content from a CallToolResult.
func extractText(result *mcp.CallToolResult) string {
	if len(result.Content) == 0 {
		return ""
	}
	if tc, ok := result.Content[0].(*mcp.TextContent); ok {
		return tc.Text
	}
	return ""
}

// unmarshalResult parses the JSON text content from a CallToolResult into the target.
// Returns an error instead of calling t.Fatal so it works with both testing.T and rapid.T.
func unmarshalResult(result *mcp.CallToolResult, target any) error {
	text := extractText(result)
	return json.Unmarshal([]byte(text), target)
}

// --- Property Tests ---

// Feature: mcp-server-mode, Property 3: Brag Create-Get Round-Trip
//
// **Validates: Requirements 2.1, 2.2**
func TestProperty_BragCreateGetRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate valid inputs (no spaces to avoid trim issues)
		title := rapid.StringMatching(`[a-zA-Z]{5,30}`).Draw(t, "title")
		description := rapid.StringMatching(`[a-zA-Z]{20,60}`).Draw(t, "description")
		categoryIdx := rapid.IntRange(0, len(validCategories)-1).Draw(t, "categoryIdx")
		category := validCategories[categoryIdx]

		ctx := context.Background()
		fx := newTestServer()

		var bragID int64 = 42
		parsedCat, _ := domain.ParseCategory(category)

		returnedBrag := &domain.Brag{
			ID:          bragID,
			Owner:       domain.User{ID: 1},
			Title:       title,
			Description: description,
			Category:    parsedCat,
			Tags:        []*domain.Tag{},
			CreatedAt:   fixedTime,
			UpdatedAt:   fixedTime,
		}

		// Mock Insert
		fx.bragRepo.On("Insert", mock.Anything, mock.MatchedBy(func(b *domain.Brag) bool {
			return b.Title == title && b.Description == description
		})).Return(returnedBrag, nil).Once()

		// Mock Select (re-fetch after create, and for get)
		fx.bragRepo.On("Select", mock.Anything, bragID).Return(returnedBrag, nil)

		// Create
		createResult, _, err := fx.server.handleBragCreate(ctx, nil, BragCreateParams{
			UserID:      1,
			Title:       title,
			Description: description,
			Category:    category,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if createResult.IsError {
			t.Fatalf("create returned error: %s", extractText(createResult))
		}

		var createResp BragResponse
		if err := unmarshalResult(createResult, &createResp); err != nil { t.Fatalf("unmarshal: %v", err) }

		// Get
		getResult, _, err := fx.server.handleBragGet(ctx, nil, BragGetParams{ID: bragID})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if getResult.IsError {
			t.Fatalf("get returned error: %s", extractText(getResult))
		}

		var getResp BragResponse
		if err := unmarshalResult(getResult, &getResp); err != nil { t.Fatalf("unmarshal: %v", err) }

		// Verify round-trip
		if createResp.Title != getResp.Title {
			t.Fatalf("title mismatch: %q vs %q", createResp.Title, getResp.Title)
		}
		if createResp.Description != getResp.Description {
			t.Fatalf("description mismatch: %q vs %q", createResp.Description, getResp.Description)
		}
		if createResp.Category != getResp.Category {
			t.Fatalf("category mismatch: %q vs %q", createResp.Category, getResp.Category)
		}
		if createResp.ID != getResp.ID {
			t.Fatalf("ID mismatch: %d vs %d", createResp.ID, getResp.ID)
		}
	})
}

// Feature: mcp-server-mode, Property 4: Brag List Invariant
//
// **Validates: Requirements 2.3**
func TestProperty_BragListInvariant(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		n := rapid.IntRange(1, 10).Draw(t, "n")

		ctx := context.Background()
		fx := newTestServer()

		var userID int64 = 1
		createdBrags := make([]*domain.Brag, 0, n)

		for i := 0; i < n; i++ {
			bragID := int64(i + 1)
			title := fmt.Sprintf("Title %d valid", i)
			desc := fmt.Sprintf("Description number %d that is long enough to pass validation", i)
			cat, _ := domain.ParseCategory("PROJECT")

			brag := &domain.Brag{
				ID:          bragID,
				Owner:       domain.User{ID: userID},
				Title:       title,
				Description: desc,
				Category:    cat,
				Tags:        []*domain.Tag{},
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			}
			createdBrags = append(createdBrags, brag)

			fx.bragRepo.On("Insert", mock.Anything, mock.Anything).Return(brag, nil).Once()
			fx.bragRepo.On("Select", mock.Anything, bragID).Return(brag, nil)
		}

		fx.bragRepo.On("SelectAll", mock.Anything, userID).Return(createdBrags, nil)

		// Create N brags
		for i := 0; i < n; i++ {
			title := fmt.Sprintf("Title %d valid", i)
			desc := fmt.Sprintf("Description number %d that is long enough to pass validation", i)
			result, _, err := fx.server.handleBragCreate(ctx, nil, BragCreateParams{
				UserID:      userID,
				Title:       title,
				Description: desc,
				Category:    "PROJECT",
			})
			if err != nil {
				t.Fatalf("unexpected error on create %d: %v", i, err)
			}
			if result.IsError {
				t.Fatalf("create %d returned error: %s", i, extractText(result))
			}
		}

		// List
		listResult, _, err := fx.server.handleBragList(ctx, nil, BragListParams{UserID: userID})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if listResult.IsError {
			t.Fatalf("list returned error: %s", extractText(listResult))
		}

		var listResp []BragResponse
		if err := unmarshalResult(listResult, &listResp); err != nil { t.Fatalf("unmarshal: %v", err) }

		// Verify count
		if len(listResp) != n {
			t.Fatalf("expected %d brags, got %d", n, len(listResp))
		}

		// Verify unique IDs
		ids := make(map[int64]bool)
		for _, b := range listResp {
			ids[b.ID] = true
		}
		if len(ids) != n {
			t.Fatalf("expected %d unique IDs, got %d", n, len(ids))
		}
	})
}

// Feature: mcp-server-mode, Property 9: Invalid Category Error
//
// **Validates: Requirements 2.8**
func TestProperty_InvalidCategoryError(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a lowercase string that cannot be a valid category
		invalidCat := rapid.StringMatching(`[a-z]{3,15}`).Draw(t, "invalidCategory")

		ctx := context.Background()
		fx := newTestServer()

		result, _, err := fx.server.handleBragCreate(ctx, nil, BragCreateParams{
			UserID:      1,
			Title:       "Valid Title Here",
			Description: "This is a valid description with enough characters",
			Category:    invalidCat,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.IsError {
			t.Fatal("expected error result for invalid category")
		}

		text := extractText(result)
		if text == "" {
			t.Fatal("expected non-empty error text")
		}
	})
}

// --- Example-based Tests ---

func TestBragCreate_ValidationFailure_ShortTitle(t *testing.T) {
	ctx := context.Background()
	fx := newTestServer()

	result, _, err := fx.server.handleBragCreate(ctx, nil, BragCreateParams{
		UserID:      1,
		Title:       "Bad",
		Description: "This is a valid description that is long enough",
		Category:    "PROJECT",
	})
	assert.NoError(t, err)
	assert.True(t, result.IsError)

	text := extractText(result)
	assert.Contains(t, text, "validation failed")
	assert.Contains(t, text, "title")
}

func TestBragCreate_ValidationFailure_ShortDescription(t *testing.T) {
	ctx := context.Background()
	fx := newTestServer()

	result, _, err := fx.server.handleBragCreate(ctx, nil, BragCreateParams{
		UserID:      1,
		Title:       "Valid Title Here",
		Description: "Too short",
		Category:    "PROJECT",
	})
	assert.NoError(t, err)
	assert.True(t, result.IsError)

	text := extractText(result)
	assert.Contains(t, text, "validation failed")
	assert.Contains(t, text, "description")
}

func TestBragGet_NotFound(t *testing.T) {
	ctx := context.Background()
	fx := newTestServer()

	fx.bragRepo.On("Select", mock.Anything, int64(999)).Return(nil, fmt.Errorf("brag not found"))

	result, _, err := fx.server.handleBragGet(ctx, nil, BragGetParams{ID: 999})
	assert.NoError(t, err)
	assert.True(t, result.IsError)

	text := extractText(result)
	assert.Contains(t, text, "not found")
}

func TestBragUpdate_NotFound(t *testing.T) {
	ctx := context.Background()
	fx := newTestServer()

	fx.bragRepo.On("Select", mock.Anything, int64(999)).Return(nil, fmt.Errorf("brag not found"))

	result, _, err := fx.server.handleBragUpdate(ctx, nil, BragUpdateParams{
		ID:    999,
		Title: "New Title Here",
	})
	assert.NoError(t, err)
	assert.True(t, result.IsError)

	text := extractText(result)
	assert.Contains(t, text, "not found")
}

func TestBragDelete_NotFound(t *testing.T) {
	ctx := context.Background()
	fx := newTestServer()

	fx.bragRepo.On("Delete", mock.Anything, int64(999)).Return(fmt.Errorf("brag not found"))

	result, _, err := fx.server.handleBragDelete(ctx, nil, BragDeleteParams{ID: 999})
	assert.NoError(t, err)
	assert.True(t, result.IsError)

	text := extractText(result)
	assert.Contains(t, text, "not found")
}
