package commands

import (
	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

func NewBragCmd(bragService *service.BragService, tagService *service.TagService) *cobra.Command {
	bragCmd := &cobra.Command{
		Use:   "brag",
		Short: "Manage your brag entries",
		Long:  `Create, list and manage your professional achievements`,
		// Check initialization before running any brag subcommand
		PersistentPreRunE: requiresInitialization(),
	}

	bragCmd.AddCommand(
		NewBragAddCmd(bragService, tagService),
		NewBragListCmd(bragService, tagService),
		NewBragEditCmd(bragService, tagService),
		NewBragRemoveCmd(bragService),
		NewBragShowCmd(bragService, tagService),
	)

	return bragCmd
}
