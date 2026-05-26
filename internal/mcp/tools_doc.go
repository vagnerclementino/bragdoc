package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

// handleDocGenerate generates a brag document for a user.
func (s *Server) handleDocGenerate(ctx context.Context, _ *mcp.CallToolRequest, params DocGenerateParams) (*mcp.CallToolResult, any, error) {
	// Validate format
	if params.Format != "" && params.Format != "markdown" {
		return toolError(fmt.Errorf("validation failed: only markdown format is supported")), nil, nil
	}

	// Validate template
	if params.Template != "" && params.Template != "default" {
		return toolError(fmt.Errorf("validation failed: only default template is supported")), nil, nil
	}

	// Fetch all brags for user
	brags, err := s.bragService.List(ctx, params.UserID)
	if err != nil {
		return toolError(err), nil, nil
	}

	if len(brags) == 0 {
		return toolError(fmt.Errorf("not found: no brags found for user")), nil, nil
	}

	// Generate document
	opts := service.GenerateOptions{
		Format:   domain.FormatMarkdown,
		Template: params.Template,
	}

	doc, err := s.docService.Generate(ctx, brags, params.UserID, opts)
	if err != nil {
		return toolError(err), nil, nil
	}

	resp := DocResponse{
		Content:     string(doc.Content),
		Title:       doc.Metadata.Title,
		Author:      doc.Metadata.Author,
		BragCount:   doc.Metadata.BragCount,
		Categories:  doc.Metadata.Categories,
		Tags:        doc.Metadata.Tags,
		GeneratedAt: doc.Metadata.GeneratedAt,
	}

	return marshalResult(resp)
}
