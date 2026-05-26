package mcp

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

// Server wraps the MCP SDK server and registers bragdoc tools.
type Server struct {
	mcpServer   *mcp.Server
	bragService *service.BragService
	tagService  *service.TagService
	userService *service.UserService
	docService  *service.DocumentService
	jobService  *service.JobTitleService
}

// NewServer creates an MCP server with all tools registered.
func NewServer(
	bragService *service.BragService,
	tagService *service.TagService,
	userService *service.UserService,
	docService *service.DocumentService,
	jobService *service.JobTitleService,
) *Server {
	version := "dev"
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		version = info.Main.Version
	}

	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "bragdoc",
		Version: version,
	}, nil)

	s := &Server{
		mcpServer:   mcpServer,
		bragService: bragService,
		tagService:  tagService,
		userService: userService,
		docService:  docService,
		jobService:  jobService,
	}

	s.registerTools()
	return s
}

// Run starts the MCP server on stdio transport, blocking until EOF or error.
func (s *Server) Run(ctx context.Context) error {
	fmt.Fprintln(os.Stderr, "bragdoc MCP server ready")
	return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}

// registerTools registers all 17 bragdoc tools with the MCP server.
func (s *Server) registerTools() {
	// Brag tools
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "brag_create",
		Description: "Create a new brag entry (achievement record)",
	}, s.handleBragCreate)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "brag_get",
		Description: "Retrieve a brag entry by ID",
	}, s.handleBragGet)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "brag_list",
		Description: "List all brag entries for a user",
	}, s.handleBragList)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "brag_search_by_tags",
		Description: "Search brag entries by tag names",
	}, s.handleBragSearchByTags)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "brag_search_by_category",
		Description: "Search brag entries by category",
	}, s.handleBragSearchByCategory)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "brag_update",
		Description: "Update an existing brag entry",
	}, s.handleBragUpdate)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "brag_delete",
		Description: "Delete a brag entry by ID",
	}, s.handleBragDelete)

	// Tag tools
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "tag_create",
		Description: "Create a new tag",
	}, s.handleTagCreate)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "tag_list",
		Description: "List all tags for a user",
	}, s.handleTagList)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "tag_attach",
		Description: "Attach tags to a brag entry",
	}, s.handleTagAttach)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "tag_detach",
		Description: "Detach tags from a brag entry",
	}, s.handleTagDetach)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "tag_delete",
		Description: "Delete a tag by ID",
	}, s.handleTagDelete)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "tag_get_or_create",
		Description: "Get an existing tag by name or create it if it does not exist",
	}, s.handleTagGetOrCreate)

	// Doc tools
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "doc_generate",
		Description: "Generate a brag document for a user",
	}, s.handleDocGenerate)

	// User tools
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "user_get",
		Description: "Get a user profile by ID",
	}, s.handleUserGet)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "user_list",
		Description: "List all users",
	}, s.handleUserList)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "user_get_by_email",
		Description: "Get a user profile by email address",
	}, s.handleUserGetByEmail)
}






