package mcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var errToolFailure = errors.New("tool execution failed")

// transactionalTool preserves tool failures while asking storage to roll back.
// Successful results reach the client only after the transaction commits.
func transactionalTool[In any](s *Server, handler mcp.ToolHandlerFor[In, any]) mcp.ToolHandlerFor[In, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, params In) (*mcp.CallToolResult, any, error) {
		if s.transactions == nil {
			return toolError(fmt.Errorf("write transaction runner is unavailable")), nil, nil
		}
		var result *mcp.CallToolResult
		var output any
		err := s.transactions.WithinTransaction(ctx, func(txCtx context.Context) error {
			var err error
			result, output, err = handler(txCtx, req, params)
			if err != nil {
				return err
			}
			if result != nil && result.IsError {
				return errToolFailure
			}
			return nil
		})
		if errors.Is(err, errToolFailure) {
			return result, nil, nil
		}
		if err != nil {
			return toolError(err), nil, nil
		}
		return result, output, nil
	}
}
