package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	ErrCodeValidation = -32602 // Invalid params
	ErrCodeNotFound   = -32001 // Resource not found
	ErrCodeInternal   = -32603 // Internal error
)

// ErrorResponse is the structured error payload returned in CallToolResult content.
// Clients can parse the JSON text content to branch on the "code" field.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// classifyError maps a service error to an appropriate error code and safe message.
func classifyError(err error) (code int, message string) {
	msg := err.Error()
	switch {
	case strings.HasPrefix(msg, "validation failed:"):
		return ErrCodeValidation, msg
	case strings.Contains(msg, "not found"):
		return ErrCodeNotFound, msg
	default:
		// Log internal errors to stderr for debugging
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		return ErrCodeInternal, "internal error"
	}
}

// toolError builds an error CallToolResult from a service error.
// The content is a JSON object with "code" and "message" fields so clients
// can programmatically handle different error types.
func toolError(err error) *mcp.CallToolResult {
	code, message := classifyError(err)
	resp := ErrorResponse{Code: code, Message: message}
	data, _ := json.Marshal(resp)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
		IsError: true,
	}
}
