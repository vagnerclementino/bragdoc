package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the application version, set at build time
	Version = "unknown"
	// Build is the git commit hash, set at build time
	Build = "unknown"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Bragdoc",
		Long:  `All software has versions. This is Bragdoc's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Bragdoc %s (build: %s)\n", Version, Build)
		},
	}
}
