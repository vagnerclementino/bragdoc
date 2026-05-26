package mcp

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"pgregory.net/rapid"
)

// Feature: mcp-server-mode, Property 10: Tag Create-List Round-Trip
//
// **Validates: Requirements 3.1, 3.4**
func TestProperty_TagCreateListRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		n := rapid.IntRange(1, 8).Draw(t, "n")

		ctx := context.Background()
		fx := newTestServer()

		var ownerID int64 = 1
		createdTags := make([]*domain.Tag, 0, n)

		// Generate N unique names
		nameSet := make(map[string]bool)
		names := make([]string, 0, n)
		for len(names) < n {
			name := rapid.StringMatching(`[a-z]{2,20}`).Draw(t, fmt.Sprintf("name_%d", len(names)))
			if !nameSet[name] {
				nameSet[name] = true
				names = append(names, name)
			}
		}

		for i, name := range names {
			tagID := int64(i + 1)
			tag := &domain.Tag{
				ID:        tagID,
				Name:      name,
				OwnerID:   ownerID,
				CreatedAt: fixedTime,
			}
			createdTags = append(createdTags, tag)

			// Mock SelectByName to return nil (tag doesn't exist yet)
			fx.tagRepo.On("SelectByName", mock.Anything, ownerID, name).Return(nil, fmt.Errorf("not found")).Once()
			// Mock Insert
			fx.tagRepo.On("Insert", mock.Anything, mock.MatchedBy(func(t *domain.Tag) bool {
				return t.Name == name && t.OwnerID == ownerID
			})).Return(tag, nil).Once()
		}

		// Mock SelectAll for list
		fx.tagRepo.On("SelectAll", mock.Anything, ownerID).Return(createdTags, nil)

		// Create N tags
		for _, name := range names {
			result, _, err := fx.server.handleTagCreate(ctx, nil, TagCreateParams{
				OwnerID: ownerID,
				Name:    name,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.IsError {
				t.Fatalf("create returned error: %s", extractText(result))
			}
		}

		// List
		listResult, _, err := fx.server.handleTagList(ctx, nil, TagListParams{UserID: ownerID})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if listResult.IsError {
			t.Fatalf("list returned error: %s", extractText(listResult))
		}

		var listResp []TagResponse
		if err := unmarshalResult(listResult, &listResp); err != nil { t.Fatalf("unmarshal: %v", err) }

		if len(listResp) != n {
			t.Fatalf("expected %d tags, got %d", n, len(listResp))
		}
	})
}

// Feature: mcp-server-mode, Property 11: Duplicate Tag Name Error
//
// **Validates: Requirements 3.2**
func TestProperty_DuplicateTagNameError(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		name := rapid.StringMatching(`[a-z]{2,20}`).Draw(t, "name")

		ctx := context.Background()
		fx := newTestServer()

		var ownerID int64 = 1
		existingTag := &domain.Tag{
			ID:        1,
			Name:      name,
			OwnerID:   ownerID,
			CreatedAt: fixedTime,
		}

		// First create: tag doesn't exist, then insert succeeds
		fx.tagRepo.On("SelectByName", mock.Anything, ownerID, name).Return(nil, fmt.Errorf("not found")).Once()
		fx.tagRepo.On("Insert", mock.Anything, mock.MatchedBy(func(t *domain.Tag) bool {
			return t.Name == name
		})).Return(existingTag, nil).Once()

		// First create succeeds
		result1, _, err := fx.server.handleTagCreate(ctx, nil, TagCreateParams{
			OwnerID: ownerID,
			Name:    name,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result1.IsError {
			t.Fatalf("first create should succeed: %s", extractText(result1))
		}

		// Second create: tag now exists
		fx.tagRepo.On("SelectByName", mock.Anything, ownerID, name).Return(existingTag, nil)

		result2, _, err := fx.server.handleTagCreate(ctx, nil, TagCreateParams{
			OwnerID: ownerID,
			Name:    name,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result2.IsError {
			t.Fatal("second create should return error for duplicate name")
		}

		text := extractText(result2)
		if text == "" {
			t.Fatal("expected non-empty error message")
		}
	})
}

// Feature: mcp-server-mode, Property 13: Tag Get-or-Create Idempotence
//
// **Validates: Requirements 3.9**
func TestProperty_TagGetOrCreateIdempotence(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		name := rapid.StringMatching(`[a-z]{2,20}`).Draw(t, "name")

		ctx := context.Background()
		fx := newTestServer()

		var ownerID int64 = 1
		var tagID int64 = 7

		existingTag := &domain.Tag{
			ID:        tagID,
			Name:      name,
			OwnerID:   ownerID,
			CreatedAt: fixedTime,
		}

		// First call: GetOrCreate calls SelectByName (not found), then Create also calls SelectByName (not found), then Insert
		fx.tagRepo.On("SelectByName", mock.Anything, ownerID, name).Return(nil, fmt.Errorf("not found")).Times(2)
		fx.tagRepo.On("Insert", mock.Anything, mock.MatchedBy(func(t *domain.Tag) bool {
			return t.Name == name
		})).Return(existingTag, nil).Once()

		result1, _, err := fx.server.handleTagGetOrCreate(ctx, nil, TagGetOrCreateParams{
			OwnerID: ownerID,
			Name:    name,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result1.IsError {
			t.Fatalf("first call should succeed: %s", extractText(result1))
		}

		var resp1 TagResponse
		if err := unmarshalResult(result1, &resp1); err != nil { t.Fatalf("unmarshal: %v", err) }

		// Second call: tag now exists → SelectByName returns the tag
		fx.tagRepo.ExpectedCalls = nil // Reset mock expectations for clean second call
		fx.tagRepo.On("SelectByName", mock.Anything, ownerID, name).Return(existingTag, nil)

		result2, _, err := fx.server.handleTagGetOrCreate(ctx, nil, TagGetOrCreateParams{
			OwnerID: ownerID,
			Name:    name,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result2.IsError {
			t.Fatalf("second call should succeed: %s", extractText(result2))
		}

		var resp2 TagResponse
		if err := unmarshalResult(result2, &resp2); err != nil { t.Fatalf("unmarshal: %v", err) }

		// Same ID both times
		if resp1.ID != resp2.ID {
			t.Fatalf("expected same ID, got %d and %d", resp1.ID, resp2.ID)
		}
	})
}

// --- Example-based Tests ---

func TestTagCreate_NameTooShort(t *testing.T) {
	ctx := context.Background()
	fx := newTestServer()

	result, _, err := fx.server.handleTagCreate(ctx, nil, TagCreateParams{
		OwnerID: 1,
		Name:    "a",
	})
	assert.NoError(t, err)
	assert.True(t, result.IsError)

	text := extractText(result)
	assert.Contains(t, text, "validation failed")
	assert.Contains(t, text, "2 characters")
}

func TestTagCreate_NameTooLong(t *testing.T) {
	ctx := context.Background()
	fx := newTestServer()

	result, _, err := fx.server.handleTagCreate(ctx, nil, TagCreateParams{
		OwnerID: 1,
		Name:    "this-name-is-way-too-long-for-a-tag",
	})
	assert.NoError(t, err)
	assert.True(t, result.IsError)

	text := extractText(result)
	assert.Contains(t, text, "validation failed")
	assert.Contains(t, text, "20 characters")
}

func TestTagAttach_EmptyTagIDs(t *testing.T) {
	ctx := context.Background()
	fx := newTestServer()

	result, _, err := fx.server.handleTagAttach(ctx, nil, TagAttachParams{
		BragID: 1,
		TagIDs: []int64{},
	})
	assert.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestTagDetach_EmptyTagIDs(t *testing.T) {
	ctx := context.Background()
	fx := newTestServer()

	result, _, err := fx.server.handleTagDetach(ctx, nil, TagDetachParams{
		BragID: 1,
		TagIDs: []int64{},
	})
	assert.NoError(t, err)
	assert.True(t, result.IsError)
}
