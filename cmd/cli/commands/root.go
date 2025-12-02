package commands

import (
	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

func NewRootCmd(bragService *service.BragService, userService *service.UserService, tagService *service.TagService) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "bragdoc",
		Short: "Bragdoc - Document your professional achievements",
		Long: `Bragdoc is a powerful command-line interface (CLI) tool designed to help individuals
build their own "Brag Documents" to track and showcase their professional achievements.`,
	}

	rootCmd.AddCommand(
		NewBragCmd(bragService, tagService),
		NewTagCmd(tagService),
		NewInitCmd(),
		NewVersionCmd(),
	)

	return rootCmd
}
