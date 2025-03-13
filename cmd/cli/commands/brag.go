package commands

import (
	"github.com/spf13/cobra"
)

func NewBragCmd() *cobra.Command {
	bragCmd := &cobra.Command{
		Use:   "brag",
		Short: "Manage your brag entries",
		Long:  `Create, list and manage your professional achievements`,
	}

	bragCmd.AddCommand(
		NewCreateCmd(),
		NewListCmd(),
	)

	return bragCmd
}
