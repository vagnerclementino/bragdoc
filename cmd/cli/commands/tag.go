package commands

import (
	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

func NewTagCmd(tagService *service.TagService) *cobra.Command {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags for organizing brags",
		Long:  `Create, list and manage tags to categorize your professional achievements`,
		// Check initialization before running any tag subcommand
		PersistentPreRunE: requiresInitialization(),
	}

	tagCmd.AddCommand(
		NewTagListCmd(tagService),
		NewTagAddCmd(tagService),
		NewTagRemoveCmd(tagService),
	)

	return tagCmd
}
