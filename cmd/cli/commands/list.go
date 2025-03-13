package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list [flags]",
		Short: "List all brag entries",
		Long:  `List all your documented professional achievements`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listing brag entries...")
			// TODO: Implement brag listing
		},
	}
}
