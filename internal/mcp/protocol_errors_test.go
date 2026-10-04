package mcp

import (
	"encoding/json"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vagnerclementino/bragdoc/internal/domain"
)

func requireErrorEnvelope(t *testing.T, result *mcp.CallToolResult, code int) ErrorResponse {
	t.Helper()
	require.True(t, result.IsError)
	var response ErrorResponse
	require.NoError(t, unmarshalResult(result, &response))
	assert.Equal(t, code, response.Code)
	structured, err := json.Marshal(result.StructuredContent)
	require.NoError(t, err)
	assert.JSONEq(t, extractText(result), string(structured))
	return response
}

func TestIntegration_ErrorEnvelopes(t *testing.T) {
	is := newIntegrationServer(t)
	ctx, client := newIntegrationClient(t, is.server)
	user, err := is.userService.Create(ctx, &domain.User{Name: "Error Test", Email: "errors@example.com"})
	require.NoError(t, err)
	_, err = is.tagService.Create(ctx, &domain.Tag{OwnerID: user.ID, Name: "existing"})
	require.NoError(t, err)

	for _, test := range []struct {
		name      string
		tool      string
		arguments map[string]any
		code      int
	}{
		{"missing required argument", "brag_get", map[string]any{}, ErrCodeValidation},
		{"wrong argument type", "brag_get", map[string]any{"id": "one"}, ErrCodeValidation},
		{"unknown field", "brag_get", map[string]any{"id": 1, "extra": true}, ErrCodeValidation},
		{"missing resource", "brag_get", map[string]any{"id": 999}, ErrCodeNotFound},
		{"empty tag search", "brag_search_by_tags", map[string]any{"user_id": user.ID, "tag_names": []string{}}, ErrCodeValidation},
		{"empty attachment", "tag_attach", map[string]any{"brag_id": 1, "tag_ids": []int64{}}, ErrCodeValidation},
		{"empty detachment", "tag_detach", map[string]any{"brag_id": 1, "tag_ids": []int64{}}, ErrCodeValidation},
		{"duplicate tag", "tag_create", map[string]any{"owner_id": user.ID, "name": "existing"}, ErrCodeValidation},
		{"short title", "brag_create", map[string]any{"user_id": user.ID, "title": "Hi", "description": "A sufficiently detailed achievement description", "category": "PROJECT"}, ErrCodeValidation},
		{"invalid category", "brag_search_by_category", map[string]any{"user_id": user.ID, "category": "invalid"}, ErrCodeValidation},
		{"invalid document format", "doc_generate", map[string]any{"user_id": user.ID, "format": "pdf"}, ErrCodeValidation},
		{"empty document", "doc_generate", map[string]any{"user_id": user.ID}, ErrCodeNotFound},
	} {
		t.Run(test.name, func(t *testing.T) {
			result, err := client.CallTool(ctx, &mcp.CallToolParams{Name: test.tool, Arguments: test.arguments})
			require.NoError(t, err)
			requireErrorEnvelope(t, result, test.code)
		})
	}
	// Infrastructure errors must not expose connection details to the client.
	require.NoError(t, is.db.Close())
	result, err := client.CallTool(ctx, &mcp.CallToolParams{Name: "user_list", Arguments: map[string]any{}})
	require.NoError(t, err)
	response := requireErrorEnvelope(t, result, ErrCodeInternal)
	assert.Equal(t, "internal error", response.Message)
}
