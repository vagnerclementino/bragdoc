package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/vagnerclementino/bragdoc/internal/domain"
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
	case errors.Is(err, domain.ErrValidation):
		return ErrCodeValidation, msg
	case errors.Is(err, domain.ErrNotFound):
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
	return errorResult(code, message)
}

func errorResult(code int, message string) *mcp.CallToolResult {
	resp := ErrorResponse{Code: code, Message: message}
	data, _ := json.Marshal(resp)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
		IsError:           true,
		StructuredContent: resp,
	}
}

// normalizeInputErrors gives SDK schema/decoding errors the same envelope as
// application errors. Our handlers always return structured errors; the SDK's
// input-validation errors occur before the handler and contain only text.
func normalizeInputErrors(next mcp.MethodHandler) mcp.MethodHandler {
	return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
		result, err := next(ctx, method, req)
		if err != nil || method != "tools/call" {
			return result, err
		}
		toolResult, ok := result.(*mcp.CallToolResult)
		if !ok || !toolResult.IsError || toolResult.StructuredContent != nil {
			return result, nil
		}
		message := "invalid tool arguments"
		if len(toolResult.Content) > 0 {
			if content, ok := toolResult.Content[0].(*mcp.TextContent); ok {
				message = content.Text
			}
		}
		return errorResult(ErrCodeValidation, message), nil
	}
}
