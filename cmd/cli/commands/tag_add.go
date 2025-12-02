package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

func NewTagAddCmd(tagService *service.TagService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new tag",
		Long:  `Create a new tag for organizing your brags`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagAdd(cmd.Context(), tagService, cmd)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Tag name (required)")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runTagAdd(ctx context.Context, tagService *service.TagService, cmd *cobra.Command) error {
	name, _ := cmd.Flags().GetString("name")

	// Trim and validate name
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("tag name cannot be empty")
	}

	// TODO: Get actual user ID from config/session
	userID := int64(1)

	// Create tag
	tag := &domain.Tag{
		Name:      name,
		OwnerID:   userID,
		CreatedAt: time.Now(),
	}

	created, err := tagService.Create(ctx, tag)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	fmt.Printf("✅ Tag created successfully! ID: %d, Name: %s\n", created.ID, created.Name)
	return nil
}
