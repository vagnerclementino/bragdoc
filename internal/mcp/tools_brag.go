package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/vagnerclementino/bragdoc/internal/domain"
)

// handleBragCreate creates a new brag entry with optional tag attachment.
func (s *Server) handleBragCreate(ctx context.Context, _ *mcp.CallToolRequest, params BragCreateParams) (*mcp.CallToolResult, any, error) {
	cat, err := domain.ParseCategory(params.Category)
	if err != nil {
		return toolError(fmt.Errorf("validation failed: %v", err)), nil, nil
	}

	brag := &domain.Brag{
		Owner:       domain.User{ID: params.UserID},
		Title:       params.Title,
		Description: params.Description,
		Category:    cat,
	}

	created, err := s.bragService.Create(ctx, brag)
	if err != nil {
		return toolError(err), nil, nil
	}

	// Attach optional tags
	if len(params.Tags) > 0 {
		for _, tagName := range params.Tags {
			tag, err := s.tagService.GetOrCreate(ctx, params.UserID, tagName)
			if err != nil {
				return toolError(err), nil, nil
			}
			if err := s.tagService.AttachToBrag(ctx, created.ID, []int64{tag.ID}); err != nil {
				return toolError(err), nil, nil
			}
		}
	}

	// Re-fetch to get tags populated
	result, err := s.bragService.GetByID(ctx, created.ID)
	if err != nil {
		return toolError(err), nil, nil
	}

	return marshalResult(toBragResponse(result))
}

// handleBragGet retrieves a brag entry by ID.
func (s *Server) handleBragGet(ctx context.Context, _ *mcp.CallToolRequest, params BragGetParams) (*mcp.CallToolResult, any, error) {
	brag, err := s.bragService.GetByID(ctx, params.ID)
	if err != nil {
		return toolError(err), nil, nil
	}

	return marshalResult(toBragResponse(brag))
}

// handleBragList lists all brag entries for a user.
func (s *Server) handleBragList(ctx context.Context, _ *mcp.CallToolRequest, params BragListParams) (*mcp.CallToolResult, any, error) {
	brags, err := s.bragService.List(ctx, params.UserID)
	if err != nil {
		return toolError(err), nil, nil
	}

	responses := make([]BragResponse, 0, len(brags))
	for _, b := range brags {
		responses = append(responses, toBragResponse(b))
	}

	return marshalResult(responses)
}

// handleBragSearchByTags searches brag entries by tag names.
func (s *Server) handleBragSearchByTags(ctx context.Context, _ *mcp.CallToolRequest, params BragSearchByTagsParams) (*mcp.CallToolResult, any, error) {
	brags, err := s.bragService.SearchByTags(ctx, params.UserID, params.TagNames)
	if err != nil {
		return toolError(err), nil, nil
	}

	responses := make([]BragResponse, 0, len(brags))
	for _, b := range brags {
		responses = append(responses, toBragResponse(b))
	}

	return marshalResult(responses)
}

// handleBragSearchByCategory searches brag entries by category.
func (s *Server) handleBragSearchByCategory(ctx context.Context, _ *mcp.CallToolRequest, params BragSearchByCategoryParams) (*mcp.CallToolResult, any, error) {
	cat, err := domain.ParseCategory(params.Category)
	if err != nil {
		return toolError(fmt.Errorf("validation failed: %v", err)), nil, nil
	}

	brags, err := s.bragService.SearchByCategory(ctx, params.UserID, cat)
	if err != nil {
		return toolError(err), nil, nil
	}

	responses := make([]BragResponse, 0, len(brags))
	for _, b := range brags {
		responses = append(responses, toBragResponse(b))
	}

	return marshalResult(responses)
}

// handleBragUpdate updates an existing brag entry.
func (s *Server) handleBragUpdate(ctx context.Context, _ *mcp.CallToolRequest, params BragUpdateParams) (*mcp.CallToolResult, any, error) {
	existing, err := s.bragService.GetByID(ctx, params.ID)
	if err != nil {
		return toolError(err), nil, nil
	}

	if params.Title != "" {
		existing.Title = params.Title
	}
	if params.Description != "" {
		existing.Description = params.Description
	}
	if params.Category != "" {
		cat, err := domain.ParseCategory(params.Category)
		if err != nil {
			return toolError(fmt.Errorf("validation failed: %v", err)), nil, nil
		}
		existing.Category = cat
	}

	updated, err := s.bragService.Update(ctx, existing)
	if err != nil {
		return toolError(err), nil, nil
	}

	// Handle tag replacement if provided
	if params.Tags != nil {
		// Detach all existing tags
		if len(existing.Tags) > 0 {
			existingTagIDs := make([]int64, 0, len(existing.Tags))
			for _, t := range existing.Tags {
				existingTagIDs = append(existingTagIDs, t.ID)
			}
			if err := s.tagService.DetachFromBrag(ctx, updated.ID, existingTagIDs); err != nil {
				return toolError(err), nil, nil
			}
		}

		// Attach new tags
		for _, tagName := range params.Tags {
			tag, err := s.tagService.GetOrCreate(ctx, existing.Owner.ID, tagName)
			if err != nil {
				return toolError(err), nil, nil
			}
			if err := s.tagService.AttachToBrag(ctx, updated.ID, []int64{tag.ID}); err != nil {
				return toolError(err), nil, nil
			}
		}

		// Re-fetch to get updated tags
		updated, err = s.bragService.GetByID(ctx, updated.ID)
		if err != nil {
			return toolError(err), nil, nil
		}
	}

	return marshalResult(toBragResponse(updated))
}

// handleBragDelete deletes a brag entry by ID.
func (s *Server) handleBragDelete(ctx context.Context, _ *mcp.CallToolRequest, params BragDeleteParams) (*mcp.CallToolResult, any, error) {
	if err := s.bragService.Delete(ctx, params.ID); err != nil {
		return toolError(err), nil, nil
	}

	return marshalResult(SuccessResponse{
		Success: true,
		Message: fmt.Sprintf("brag %d deleted", params.ID),
	})
}

// toBragResponse converts a domain Brag to a BragResponse.
func toBragResponse(b *domain.Brag) BragResponse {
	tags := make([]string, 0, len(b.Tags))
	for _, t := range b.Tags {
		tags = append(tags, t.Name)
	}

	return BragResponse{
		ID:          b.ID,
		Title:       b.Title,
		Description: b.Description,
		Category:    string(b.Category.Name),
		Tags:        tags,
		CreatedAt:   b.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   b.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// marshalResult serializes a value to JSON and wraps it in a CallToolResult.
func marshalResult(v any) (*mcp.CallToolResult, any, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return toolError(fmt.Errorf("internal error: failed to marshal response")), nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}
