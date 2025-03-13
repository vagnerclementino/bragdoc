package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [description]",
		Short: "Create a new brag entry",
		Long:  `Create a new brag entry to document your professional achievements`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Creating new brag entry...")
			// TODO: Implement brag creation
		},
	}
}
