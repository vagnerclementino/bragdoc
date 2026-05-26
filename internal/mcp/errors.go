package mcp

import (
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	ErrCodeValidation = -32602 // Invalid params (JSON-RPC standard)
	ErrCodeNotFound   = -32001 // Resource not found (application-defined)
	ErrCodeInternal   = -32603 // Internal error (JSON-RPC standard)
)

// classifyError maps a service error to an appropriate JSON-RPC error code and safe message.
func classifyError(err error) (code int, message string) {
	msg := err.Error()
	switch {
	case strings.HasPrefix(msg, "validation failed:"):
		return ErrCodeValidation, msg
	case strings.Contains(msg, "not found"):
		return ErrCodeNotFound, msg
	default:
		return ErrCodeInternal, "internal error"
	}
}

// toolError builds an error CallToolResult from a service error using classifyError.
func toolError(err error) *mcp.CallToolResult {
	_, message := classifyError(err)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: message},
		},
		IsError: true,
	}
}
