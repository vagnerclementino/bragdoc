package command

import (
	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/command/brag"
	"github.com/vagnerclementino/bragdoc/internal/command/doc"
	"github.com/vagnerclementino/bragdoc/internal/command/tag"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

// NewRootCmd creates the root command for the bragdoc CLI.
func NewRootCmd(bragService *service.BragService, userService *service.UserService, tagService *service.TagService, docService *service.DocumentService) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "bragdoc",
		Short: "Bragdoc - Document your professional achievements",
		Long: `Bragdoc is a powerful command-line interface (CLI) tool designed to help individuals
build their own "Brag Documents" to track and showcase their professional achievements.`,
	}

	rootCmd.AddCommand(
		brag.NewBragCmd(bragService, userService, tagService),
		tag.NewTagCmd(tagService),
		doc.NewDocCmd(docService, bragService, tagService),
		NewInitCmd(),
		NewVersionCmd(),
	)

	return rootCmd
}
