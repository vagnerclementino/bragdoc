package tag

import (
	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

// NewTagCmd creates the root command for tag management.
func NewTagCmd(tagService *service.TagService) *cobra.Command {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags for organizing brags",
		Long:  `Create, list and manage tags to categorize your professional achievements`,
		// Check initialization before running any tag subcommand
		PersistentPreRunE: requiresInitialization(),
	}

	tagCmd.AddCommand(
		NewListCmd(tagService),
		NewAddCmd(tagService),
		NewRemoveCmd(tagService),
	)

	return tagCmd
}
