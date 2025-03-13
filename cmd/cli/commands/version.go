package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/handler"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Bragdoc",
		Long:  `All software has versions. This is Bragdoc's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Bragdoc", handler.Version)
		},
	}
}
