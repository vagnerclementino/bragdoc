package brag

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

func NewEditCmd(bragService *service.BragService, tagService *service.TagService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit an existing brag entry",
		Long:  `Edit an existing brag entry by ID`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEdit(cmd.Context(), bragService, tagService, cmd, args)
		},
	}

	// Optional flags - at least one must be provided
	cmd.Flags().StringP("title", "t", "", "New brag title")
	cmd.Flags().StringP("description", "d", "", "New brag description")
	cmd.Flags().StringP("category", "c", "", "New brag category (project|achievement|skill|leadership|innovation)")
	cmd.Flags().StringSliceP("tags", "", []string{}, "New comma-separated list of tags (replaces existing tags)")

	return cmd
}

func runEdit(ctx context.Context, bragService *service.BragService, tagService *service.TagService, cmd *cobra.Command, args []string) error {
	// Parse brag ID
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid brag ID: %w", err)
	}

	// Get existing brag
	brag, err := bragService.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get brag: %w", err)
	}

	// Check if any flags were provided
	title, _ := cmd.Flags().GetString("title")
	description, _ := cmd.Flags().GetString("description")
	categoryStr, _ := cmd.Flags().GetString("category")
	tagNames, _ := cmd.Flags().GetStringSlice("tags")
	tagsChanged := cmd.Flags().Changed("tags")

	if title == "" && description == "" && categoryStr == "" && !tagsChanged {
		return fmt.Errorf("at least one field must be provided to update")
	}

	// Update fields if provided
	if title != "" {
		brag.Title = title
	}

	if description != "" {
		brag.Description = description
	}

	if categoryStr != "" {
		category, err := domain.ParseCategory(categoryStr)
		if err != nil {
			return fmt.Errorf("invalid category: %w. Valid options: project, achievement, skill, leadership, innovation", err)
		}
		brag.Category = category
	}

	brag.UpdatedAt = time.Now()

	// Update brag
	updated, err := bragService.Update(ctx, brag)
	if err != nil {
		return fmt.Errorf("failed to update brag: %w", err)
	}

	// Handle tags if provided
	if tagsChanged {
		// First, get current tags and detach them
		currentTags, err := tagService.ListByBrag(ctx, brag.ID)
		if err == nil && len(currentTags) > 0 {
			currentTagIDs := make([]int64, len(currentTags))
			for i, tag := range currentTags {
				currentTagIDs[i] = tag.ID
			}
			// Ignore error if detaching fails
			_ = tagService.DetachFromBrag(ctx, brag.ID, currentTagIDs)
		}

		// Attach new tags
		if len(tagNames) > 0 {
			// TODO: Get actual user ID from config/session
			userID := int64(1)
			if err := attachTags(ctx, tagService, brag.ID, userID, tagNames); err != nil {
				fmt.Printf("⚠️  Warning: failed to attach tags: %v\n", err)
			}
		}
	}

	fmt.Printf("✅ Brag updated successfully! ID: %d\n", updated.ID)
	fmt.Printf("   Title: %s\n", updated.Title)
	fmt.Printf("   Category: %s\n", updated.Category.String())

	return nil
}
