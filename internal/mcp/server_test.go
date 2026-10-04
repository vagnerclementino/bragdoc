package mcp

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vagnerclementino/bragdoc/internal/database"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

// --- Real SQLite setup helper ---

// integrationServer holds a Server backed by a real temporary SQLite database.
type integrationServer struct {
	server      *Server
	userService *service.UserService
	bragService *service.BragService
	tagService  *service.TagService
	db          *sql.DB
}

// newIntegrationServer creates a Server backed by a real temporary SQLite DB with
// migrations applied. This allows full end-to-end testing without mocks.
func newIntegrationServer(t *testing.T) *integrationServer {
	t.Helper()

	// Use a temp file for the database that gets cleaned up
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"

	db, err := database.New(dbPath)
	require.NoError(t, err)

	ctx := context.Background()
	err = db.Migrate(ctx)
	require.NoError(t, err)

	sqliteDB := database.NewSQLiteDB(db.Conn())

	// Initialize repositories
	userRepo := database.NewUserRepository(sqliteDB)
	categoryRepo := database.NewCategoryRepository(sqliteDB)
	jobTitleRepo := database.NewJobTitleRepository(sqliteDB, userRepo)
	bragRepo := database.NewBragRepository(sqliteDB, userRepo, categoryRepo, jobTitleRepo)
	tagRepo := database.NewTagRepository(sqliteDB)

	// Initialize services
	bragService := service.NewBragService(bragRepo)
	userService := service.NewUserService(userRepo)
	tagService := service.NewTagService(tagRepo)
	jobTitleService := service.NewJobTitleService(jobTitleRepo)
	docService := service.NewDocumentService(userService)

	srv := NewServer(bragService, tagService, userService, docService, jobTitleService)

	t.Cleanup(func() {
		_ = db.Close()
	})

	return &integrationServer{
		server:      srv,
		userService: userService,
		bragService: bragService,
		tagService:  tagService,
		db:          db.Conn(),
	}
}

// --- Integration Tests ---

// TestIntegration_FullRoundTrip tests creating a user via service, then using
// handleBragCreate → handleBragGet to verify the full chain works end-to-end
// with a real SQLite database.
func TestIntegration_FullRoundTrip(t *testing.T) {
	is := newIntegrationServer(t)
	ctx, client := newIntegrationClient(t, is.server)

	// Create a user directly via service
	user, err := is.userService.Create(ctx, &domain.User{
		Name:  "Alice Integration",
		Email: "alice@integration.test",
	})
	require.NoError(t, err)
	require.NotZero(t, user.ID)

	// Create a brag via handler
	createResult, err := client.CallTool(ctx, &mcp.CallToolParams{Name: "brag_create", Arguments: BragCreateParams{
		UserID:      user.ID,
		Title:       "Shipped MCP Server",
		Description: "Implemented the full MCP server mode with all 17 tools registered",
		Category:    "PROJECT",
	}})
	require.NoError(t, err)
	require.False(t, createResult.IsError, "expected success, got: %s", extractText(createResult))

	var createResp BragResponse
	err = json.Unmarshal([]byte(extractText(createResult)), &createResp)
	require.NoError(t, err)
	assert.Equal(t, "Shipped MCP Server", createResp.Title)
	assert.Equal(t, "PROJECT", createResp.Category)
	assert.NotZero(t, createResp.ID)

	// Get the brag via handler
	getResult, err := client.CallTool(ctx, &mcp.CallToolParams{Name: "brag_get", Arguments: BragGetParams{ID: createResp.ID}})
	require.NoError(t, err)
	require.False(t, getResult.IsError, "expected success, got: %s", extractText(getResult))

	var getResp BragResponse
	err = json.Unmarshal([]byte(extractText(getResult)), &getResp)
	require.NoError(t, err)

	// Verify round-trip consistency
	assert.Equal(t, createResp.ID, getResp.ID)
	assert.Equal(t, createResp.Title, getResp.Title)
	assert.Equal(t, createResp.Description, getResp.Description)
	assert.Equal(t, createResp.Category, getResp.Category)
}

// newIntegrationClient exercises SDK initialization, routing, and serialization.
func newIntegrationClient(t *testing.T, srv *Server) (context.Context, *mcp.ClientSession) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := srv.mcpServer.Connect(ctx, serverTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, serverSession.Close()) })
	client := mcp.NewClient(&mcp.Implementation{Name: "bragdoc-test", Version: "1"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, session.Close()) })
	return ctx, session
}

func TestIntegration_InitializeAndDiscoverTools(t *testing.T) {
	is := newIntegrationServer(t)
	ctx, client := newIntegrationClient(t, is.server)
	initialized := client.InitializeResult()
	require.NotNil(t, initialized)
	assert.Equal(t, "bragdoc", initialized.ServerInfo.Name)
	assert.NotEmpty(t, initialized.ProtocolVersion)
	require.NotNil(t, initialized.Capabilities.Tools)
	result, err := client.ListTools(ctx, nil)
	require.NoError(t, err)
	names := make([]string, 0, len(result.Tools))
	for _, tool := range result.Tools {
		names = append(names, tool.Name)
		assert.NotNil(t, tool.InputSchema)
	}
	assert.ElementsMatch(t, []string{
		"brag_create", "brag_get", "brag_list", "brag_search_by_tags", "brag_search_by_category", "brag_update", "brag_delete",
		"tag_create", "tag_list", "tag_attach", "tag_detach", "tag_delete", "tag_get_or_create",
		"doc_generate", "user_get", "user_list", "user_get_by_email",
	}, names)
}

func TestIntegration_UnknownToolName(t *testing.T) {
	is := newIntegrationServer(t)
	ctx, client := newIntegrationClient(t, is.server)
	_, err := client.CallTool(ctx, &mcp.CallToolParams{Name: "unknown_tool", Arguments: map[string]any{}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown_tool")
	// A protocol error must not end the session.
	result, err := client.CallTool(ctx, &mcp.CallToolParams{Name: "user_list", Arguments: map[string]any{}})
	require.NoError(t, err)
	assert.False(t, result.IsError)
}

func TestIntegration_InvalidArguments(t *testing.T) {
	is := newIntegrationServer(t)
	ctx, client := newIntegrationClient(t, is.server)
	for _, arguments := range []map[string]any{{}, {"id": "invalid"}, {"id": 1, "unexpected": true}} {
		result, err := client.CallTool(ctx, &mcp.CallToolParams{Name: "brag_get", Arguments: arguments})
		require.NoError(t, err)
		assert.True(t, result.IsError)
	}
}

// TestIntegration_HandlersStateless verifies that the Server struct doesn't hold
// mutable state between handler calls — each call is independent and works
// correctly with the same server instance.
func TestIntegration_HandlersStateless(t *testing.T) {
	is := newIntegrationServer(t)
	ctx := context.Background()

	// Create a user
	user, err := is.userService.Create(ctx, &domain.User{
		Name:  "Stateless Test User",
		Email: "stateless@test.com",
	})
	require.NoError(t, err)

	// Call handleBragCreate multiple times — each should work independently
	for i := 0; i < 3; i++ {
		result, _, err := is.server.handleBragCreate(ctx, nil, BragCreateParams{
			UserID:      user.ID,
			Title:       "Stateless Brag Entry",
			Description: "This verifies that handlers do not hold state between calls",
			Category:    "ACHIEVEMENT",
		})
		require.NoError(t, err)
		require.False(t, result.IsError, "call %d failed: %s", i, extractText(result))

		var resp BragResponse
		err = json.Unmarshal([]byte(extractText(result)), &resp)
		require.NoError(t, err)
		assert.Equal(t, "Stateless Brag Entry", resp.Title)
	}

	// Verify all 3 brags were created (list should return 3)
	listResult, _, err := is.server.handleBragList(ctx, nil, BragListParams{UserID: user.ID})
	require.NoError(t, err)
	require.False(t, listResult.IsError)

	var listResp []BragResponse
	err = json.Unmarshal([]byte(extractText(listResult)), &listResp)
	require.NoError(t, err)
	assert.Len(t, listResp, 3, "expected 3 brags from stateless calls")

	// Each brag should have a unique ID
	ids := make(map[int64]bool)
	for _, b := range listResp {
		ids[b.ID] = true
	}
	assert.Len(t, ids, 3, "expected 3 unique brag IDs")
}
