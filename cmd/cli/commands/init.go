package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize user configuration",
		Long:  `Initialize user configuration by setting up name and email`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Initializing user configuration...")
			// TODO: Implement user configuration setup
		},
	}
}
