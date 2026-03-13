package doc

import (
	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

// NewDocCmd creates the root command for document management.
func NewDocCmd(docService *service.DocumentService, bragService *service.BragService, tagService *service.TagService) *cobra.Command {
	docCmd := &cobra.Command{
		Use:   "doc",
		Short: "Generate brag documents",
		Long:  `Generate professional achievement documents in various formats`,
		// Check initialization before running any doc subcommand
		PersistentPreRunE: requiresInitialization(),
	}

	docCmd.AddCommand(
		NewGenerateCmd(docService, bragService, tagService),
	)

	return docCmd
}
