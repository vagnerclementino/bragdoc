package command

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/config"
	mcpserver "github.com/vagnerclementino/bragdoc/internal/mcp"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

// NewMCPCmd serves MCP requests over stdio using the CLI's application services.
func NewMCPCmd(brags *service.BragService, users *service.UserService, tags *service.TagService, jobs *service.JobTitleService, docs *service.DocumentService) *cobra.Command {
	return &cobra.Command{
		Use:          "mcp",
		Short:        "Serve MCP tools over stdin/stdout",
		Long:         "Run the MCP server over stdio until the client disconnects. Initialize Bragdoc with 'bragdoc init' before connecting your IDE.",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if !config.NewManager().IsInitialized() {
				return fmt.Errorf("bragdoc is not initialized. Please run 'bragdoc init' first")
			}
			if brags == nil || users == nil || tags == nil || jobs == nil || docs == nil {
				return fmt.Errorf("failed to initialize MCP services; check the Bragdoc configuration")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return mcpserver.NewServer(brags, tags, users, docs, jobs).Run(cmd.Context())
		},
		// Override the root hook: MCP must not run CLI update checks on disconnect.
		PersistentPostRun: func(_ *cobra.Command, _ []string) {},
	}
}
