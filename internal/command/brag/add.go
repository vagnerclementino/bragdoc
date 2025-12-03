package brag

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

func NewAddCmd(bragService *service.BragService, tagService *service.TagService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new brag entry",
		Long:  `Add a new brag entry to document your professional achievements`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(cmd.Context(), bragService, tagService, cmd)
		},
	}

	// Required flags
	cmd.Flags().StringP("title", "t", "", "Brag title (required)")
	cmd.Flags().StringP("description", "d", "", "Brag description (required)")
	cmd.MarkFlagRequired("title")
	cmd.MarkFlagRequired("description")

	// Optional flags
	cmd.Flags().StringP("category", "c", "achievement", "Brag category (project|achievement|skill|leadership|innovation)")
	cmd.Flags().StringSliceP("tags", "", []string{}, "Comma-separated list of tags")

	return cmd
}

func runAdd(ctx context.Context, bragService *service.BragService, tagService *service.TagService, cmd *cobra.Command) error {
	title, _ := cmd.Flags().GetString("title")
	description, _ := cmd.Flags().GetString("description")
	categoryStr, _ := cmd.Flags().GetString("category")
	tagNames, _ := cmd.Flags().GetStringSlice("tags")

	// Parse category
	category, err := domain.ParseCategory(categoryStr)
	if err != nil {
		return fmt.Errorf("invalid category: %w. Valid options: project, achievement, skill, leadership, innovation", err)
	}

	// TODO: Get actual user ID from config/session
	// For now, using a default user ID of 1
	userID := int64(1)

	// Create brag
	newBrag := &domain.Brag{
		OwnerID:     userID,
		Title:       title,
		Description: description,
		Category:    category,
		CreatedAt:   time.Now(),
	}

	created, err := bragService.Create(ctx, newBrag)
	if err != nil {
		return fmt.Errorf("failed to create brag: %w", err)
	}

	// Handle tags if provided
	if len(tagNames) > 0 {
		if err := attachTags(ctx, tagService, created.ID, userID, tagNames); err != nil {
			fmt.Printf("⚠️  Warning: failed to attach some tags: %v\n", err)
		}
	}

	fmt.Printf("✅ Brag created successfully! ID: %d\n", created.ID)
	fmt.Printf("   Title: %s\n", created.Title)
	fmt.Printf("   Category: %s\n", created.Category.String())
	if len(tagNames) > 0 {
		fmt.Printf("   Tags: %s\n", strings.Join(tagNames, ", "))
	}

	return nil
}

// attachTags handles tag creation and attachment to a brag
func attachTags(ctx context.Context, tagService *service.TagService, bragID int64, userID int64, tagNames []string) error {
	var tagIDs []int64

	for _, tagName := range tagNames {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}

		// Get or create tag
		tag, err := tagService.GetOrCreate(ctx, userID, tagName)
		if err != nil {
			return fmt.Errorf("failed to get or create tag '%s': %w", tagName, err)
		}

		tagIDs = append(tagIDs, tag.ID)
	}

	if len(tagIDs) > 0 {
		if err := tagService.AttachToBrag(ctx, bragID, tagIDs); err != nil {
			return fmt.Errorf("failed to attach tags to brag: %w", err)
		}
	}

	return nil
}
