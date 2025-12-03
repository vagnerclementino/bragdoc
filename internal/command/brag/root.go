package brag

import (
	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

func NewBragCmd(bragService *service.BragService, userService *service.UserService, tagService *service.TagService) *cobra.Command {
	bragCmd := &cobra.Command{
		Use:   "brag",
		Short: "Manage your brag entries",
		Long:  `Create, list and manage your professional achievements`,
		// Check initialization before running any brag subcommand
		PersistentPreRunE: requiresInitialization(),
	}

	bragCmd.AddCommand(
		NewAddCmd(bragService, userService, tagService),
		NewListCmd(bragService, tagService),
		NewEditCmd(bragService, tagService),
		NewRemoveCmd(bragService),
		NewShowCmd(bragService, tagService),
	)

	return bragCmd
}
