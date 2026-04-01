package brag

import (
	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

// NewBragCmd creates the root command for brag management.
func NewBragCmd(bragService *service.BragService, userService *service.UserService, tagService *service.TagService, jobTitleService *service.JobTitleService) *cobra.Command {
	bragCmd := &cobra.Command{
		Use:   "brag",
		Short: "Manage your brag entries",
		Long:  `Create, list and manage your professional achievements`,
		// Check initialization before running any brag subcommand
		PersistentPreRunE: requiresInitialization(),
	}

	bragCmd.AddCommand(
		NewAddCmd(bragService, userService, tagService, jobTitleService),
		NewListCmd(bragService, tagService),
		NewEditCmd(bragService, userService, tagService, jobTitleService),
		NewRemoveCmd(bragService),
		NewShowCmd(bragService, tagService),
	)

	return bragCmd
}
