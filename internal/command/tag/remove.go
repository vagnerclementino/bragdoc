package tag

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

// NewRemoveCmd creates a new command for removing tags.
func NewRemoveCmd(tagService *service.TagService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <id>",
		Short: "Remove a tag",
		Long:  `Remove a tag by ID. This will also remove all associations with brags.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemove(cmd.Context(), tagService, cmd, args)
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	return cmd
}

func runRemove(ctx context.Context, tagService *service.TagService, cmd *cobra.Command, args []string) error {
	idStr := args[0]
	confirm, _ := cmd.Flags().GetBool("force")

	// Parse tag ID
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid tag ID: %s", idStr)
	}

	// Get tag to show what will be deleted
	tag, err := tagService.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get tag: %w", err)
	}

	// Confirmation prompt
	if !confirm {
		fmt.Printf("⚠️  About to remove tag: %s (ID: %d)\n", tag.Name, tag.ID)
		fmt.Println("This will also remove all associations with brags.")
		fmt.Print("Are you sure? (y/N): ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read confirmation: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("❌ Operation cancelled.")
			return nil
		}
	}

	// Delete tag
	if err := tagService.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	fmt.Printf("✅ Tag '%s' (ID: %d) removed successfully!\n", tag.Name, tag.ID)
	return nil
}
