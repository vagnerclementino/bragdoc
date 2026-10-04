package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vagnerclementino/bragdoc/config"
	mcpserver "github.com/vagnerclementino/bragdoc/internal/mcp"
)

func TestMCPStdioLifecycle(t *testing.T) {
	dataHome := t.TempDir()
	t.Setenv(config.BragdocHomeEnv, dataHome)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, stderr, err := runBinary([]string{"init", "--name", "MCP Test", "--email", "mcp@example.com"}, map[string]string{config.BragdocHomeEnv: dataHome})
	require.NoError(t, err, "%s", stderr)
	// Make an update check due, so an inherited post-run hook would change the
	// configuration and fail the snapshot assertion on disconnect.
	manager := config.NewManager()
	cfg, err := manager.Load(ctx)
	require.NoError(t, err)
	cfg.UpdateChecker.Enabled = true
	cfg.UpdateChecker.LastCheckedAt = time.Time{}
	require.NoError(t, manager.Save(ctx, cfg))
	configPath := filepath.Join(dataHome, "config.yaml")
	// #nosec G304 -- Test fixture path is inside t.TempDir().
	before, err := os.ReadFile(configPath)
	require.NoError(t, err)

	cmd := exec.CommandContext(ctx, binaryPath, "mcp")
	cmd.Env = append(os.Environ(), config.BragdocHomeEnv+"="+dataHome, "GOCOVERDIR=.coverdata")
	var diagnostics bytes.Buffer
	cmd.Stderr = &diagnostics
	client := mcp.NewClient(&mcp.Implementation{Name: "stdio-test", Version: "1"}, nil)
	session, err := client.Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })
	assert.Equal(t, "bragdoc", session.InitializeResult().ServerInfo.Name)
	tools, err := session.ListTools(ctx, nil)
	require.NoError(t, err)
	assert.Len(t, tools.Tools, 17)

	users, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "user_list", Arguments: map[string]any{}})
	require.NoError(t, err)
	require.False(t, users.IsError)
	var profiles []mcpserver.UserResponse
	require.NoError(t, json.Unmarshal([]byte(users.Content[0].(*mcp.TextContent).Text), &profiles))
	require.Len(t, profiles, 1)

	created, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "brag_create", Arguments: mcpserver.BragCreateParams{
		UserID: profiles[0].ID, Title: "Shipped MCP command", Description: "Integrated the MCP runtime into the existing Bragdoc executable", Category: "PROJECT",
	}})
	require.NoError(t, err)
	require.False(t, created.IsError)
	var brag mcpserver.BragResponse
	require.NoError(t, json.Unmarshal([]byte(created.Content[0].(*mcp.TextContent).Text), &brag))
	require.NotZero(t, brag.ID)
	got, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "brag_get", Arguments: mcpserver.BragGetParams{ID: brag.ID}})
	require.NoError(t, err)
	assert.Equal(t, created.Content[0].(*mcp.TextContent).Text, got.Content[0].(*mcp.TextContent).Text)

	// Closing the session closes stdin; the command must exit without signals.
	require.NoError(t, session.Close())
	require.NotNil(t, cmd.ProcessState)
	assert.True(t, cmd.ProcessState.Success(), "%s", diagnostics.String())
	assert.Contains(t, diagnostics.String(), "bragdoc MCP server ready")
	// #nosec G304 -- Test fixture path is inside t.TempDir().
	after, err := os.ReadFile(configPath)
	require.NoError(t, err)
	assert.Equal(t, before, after, "disconnect must not trigger the CLI update checker")
}
