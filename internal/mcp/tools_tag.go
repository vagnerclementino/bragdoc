package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/vagnerclementino/bragdoc/internal/domain"
)

// handleTagCreate creates a new tag.
func (s *Server) handleTagCreate(ctx context.Context, _ *mcp.CallToolRequest, params TagCreateParams) (*mcp.CallToolResult, any, error) {
	tag := &domain.Tag{
		Name:    params.Name,
		OwnerID: params.OwnerID,
	}

	created, err := s.tagService.Create(ctx, tag)
	if err != nil {
		return toolError(err), nil, nil
	}

	return marshalResult(toTagResponse(created))
}

// handleTagList lists all tags for a user.
func (s *Server) handleTagList(ctx context.Context, _ *mcp.CallToolRequest, params TagListParams) (*mcp.CallToolResult, any, error) {
	tags, err := s.tagService.ListByUser(ctx, params.UserID)
	if err != nil {
		return toolError(err), nil, nil
	}

	responses := make([]TagResponse, 0, len(tags))
	for _, t := range tags {
		responses = append(responses, toTagResponse(t))
	}

	return marshalResult(responses)
}

// handleTagAttach attaches tags to a brag entry.
func (s *Server) handleTagAttach(ctx context.Context, _ *mcp.CallToolRequest, params TagAttachParams) (*mcp.CallToolResult, any, error) {
	if err := s.tagService.AttachToBrag(ctx, params.BragID, params.TagIDs); err != nil {
		return toolError(err), nil, nil
	}

	return marshalResult(SuccessResponse{
		Success: true,
		Message: fmt.Sprintf("attached %d tag(s) to brag %d", len(params.TagIDs), params.BragID),
	})
}

// handleTagDetach detaches tags from a brag entry.
func (s *Server) handleTagDetach(ctx context.Context, _ *mcp.CallToolRequest, params TagDetachParams) (*mcp.CallToolResult, any, error) {
	if err := s.tagService.DetachFromBrag(ctx, params.BragID, params.TagIDs); err != nil {
		return toolError(err), nil, nil
	}

	return marshalResult(SuccessResponse{
		Success: true,
		Message: fmt.Sprintf("detached %d tag(s) from brag %d", len(params.TagIDs), params.BragID),
	})
}

// handleTagDelete deletes a tag by ID.
func (s *Server) handleTagDelete(ctx context.Context, _ *mcp.CallToolRequest, params TagDeleteParams) (*mcp.CallToolResult, any, error) {
	if err := s.tagService.Delete(ctx, params.ID); err != nil {
		return toolError(err), nil, nil
	}

	return marshalResult(SuccessResponse{
		Success: true,
		Message: fmt.Sprintf("tag %d deleted", params.ID),
	})
}

// handleTagGetOrCreate gets an existing tag by name or creates it.
func (s *Server) handleTagGetOrCreate(ctx context.Context, _ *mcp.CallToolRequest, params TagGetOrCreateParams) (*mcp.CallToolResult, any, error) {
	tag, err := s.tagService.GetOrCreate(ctx, params.OwnerID, params.Name)
	if err != nil {
		return toolError(err), nil, nil
	}

	return marshalResult(toTagResponse(tag))
}

// toTagResponse converts a domain.Tag to a TagResponse.
func toTagResponse(t *domain.Tag) TagResponse {
	return TagResponse{
		ID:        t.ID,
		Name:      t.Name,
		OwnerID:   t.OwnerID,
		CreatedAt: t.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
