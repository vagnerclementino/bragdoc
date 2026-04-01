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

// NewEditCmd creates a new command for editing brag entries.
func NewEditCmd(bragService *service.BragService, userService *service.UserService, tagService *service.TagService, jobTitleService *service.JobTitleService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit an existing brag entry",
		Long:  `Edit an existing brag entry by ID`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEdit(cmd.Context(), bragService, userService, tagService, jobTitleService, cmd, args)
		},
	}

	// Optional flags - at least one must be provided
	cmd.Flags().StringP("title", "t", "", "New brag title")
	cmd.Flags().StringP("description", "d", "", "New brag description")
	cmd.Flags().StringP("category", "c", "", "New brag category (project|achievement|skill|leadership|innovation)")
	cmd.Flags().StringP("job", "j", "", "New job title name")
	cmd.Flags().StringSliceP("tags", "", []string{}, "New comma-separated list of tags (replaces existing tags)")
	cmd.Flags().StringP("date", "D", "", "New date of the event (format based on locale: DD/MM/YYYY for pt-BR, MM/DD/YYYY for en-US)")

	return cmd
}

func runEdit(ctx context.Context, bragService *service.BragService, userService *service.UserService, tagService *service.TagService, jobTitleService *service.JobTitleService, cmd *cobra.Command, args []string) error {
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
	jobName, _ := cmd.Flags().GetString("job")
	tagNames, _ := cmd.Flags().GetStringSlice("tags")
	dateStr, _ := cmd.Flags().GetString("date")
	tagsChanged := cmd.Flags().Changed("tags")
	jobChanged := cmd.Flags().Changed("job")
	dateChanged := cmd.Flags().Changed("date")

	if title == "" && description == "" && categoryStr == "" && !tagsChanged && !jobChanged && !dateChanged {
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
			return fmt.Errorf("invalid category: %w. Valid options: PROJECT, ACHIEVEMENT, SKILL, LEADERSHIP, INNOVATION", err)
		}
		brag.Category = category
	}

	// Update job title if provided
	if jobChanged {
		if jobName != "" {
			// TODO: Get actual user ID from config/session
			userID := int64(1)
			jobTitle, err := jobTitleService.GetOrCreate(ctx, userID, jobName, brag.Owner.Company)
			if err != nil {
				return fmt.Errorf("failed to get or create job title: %w", err)
			}
			brag.JobTitle = jobTitle
		} else {
			// Empty string means remove job title
			brag.JobTitle = nil
		}
	}

	// Update date if provided
	if dateChanged && dateStr != "" {
		userID := int64(1)
		user, err := userService.GetByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		parsed, err := time.Parse(user.Locale.DateFormat(), dateStr)
		if err != nil {
			return fmt.Errorf("invalid date '%s': expected format %s", dateStr, user.Locale.DateFormatHint())
		}
		brag.CreatedAt = parsed
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
	if updated.JobTitle != nil {
		fmt.Printf("   Job Title: %s\n", updated.JobTitle.Title)
	}

	return nil
}
