// Package brag provides commands for managing brag entries
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

func NewAddCmd(bragService *service.BragService, userService *service.UserService, tagService *service.TagService, jobTitleService *service.JobTitleService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new brag entry",
		Long:  `Add a new brag entry to document your professional achievements`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runAdd(cmd.Context(), bragService, userService, tagService, jobTitleService, cmd)
		},
	}

	cmd.Flags().StringP("title", "t", "", "Brag title (required)")
	cmd.Flags().StringP("description", "d", "", "Brag description (required)")
	if err := cmd.MarkFlagRequired("title"); err != nil {
		panic(fmt.Sprintf("failed to mark title flag as required: %v", err))
	}
	if err := cmd.MarkFlagRequired("description"); err != nil {
		panic(fmt.Sprintf("failed to mark description flag as required: %v", err))
	}

	cmd.Flags().StringP("category", "c", "ACHIEVEMENT", "Brag category (UPPERCASE: PROJECT|ACHIEVEMENT|SKILL|LEADERSHIP|INNOVATION)")
	cmd.Flags().StringSliceP("tags", "", []string{}, "Comma-separated list of tags")
	cmd.Flags().StringP("job", "j", "", "Job Title name (optional, uses active job title if not specified)")

	return cmd
}

func runAdd(ctx context.Context, bragService *service.BragService, userService *service.UserService, tagService *service.TagService, jobTitleService *service.JobTitleService, cmd *cobra.Command) error {
	title, _ := cmd.Flags().GetString("title")
	description, _ := cmd.Flags().GetString("description")
	categoryStr, _ := cmd.Flags().GetString("category")
	tagNames, _ := cmd.Flags().GetStringSlice("tags")
	jobTitleName, _ := cmd.Flags().GetString("job")

	// Parse category
	category, err := domain.ParseCategory(categoryStr)
	if err != nil {
		return fmt.Errorf("invalid category: %w. Valid options: PROJECT, ACHIEVEMENT, SKILL, LEADERSHIP, INNOVATION", err)
	}

	// TODO: Get actual user ID from config/session
	userID := int64(1)

	// Fetch the user to populate Owner
	user, err := userService.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Handle job title
	var jobTitle *domain.JobTitle
	if jobTitleName != "" {
		// Get or create job title by name
		jobTitle, err = jobTitleService.GetOrCreate(ctx, userID, jobTitleName, user.Company)
		if err != nil {
			return fmt.Errorf("failed to get or create job title: %w", err)
		}
	} else {
		// Try to get active job title
		jobTitle, err = jobTitleService.GetActive(ctx, userID)
		if err != nil {
			// It's ok if there's no active job title
			jobTitle = nil
		}
	}

	// Create brag
	newBrag := &domain.Brag{
		Owner:       *user,
		Title:       title,
		Description: description,
		Category:    category,
		JobTitle:    jobTitle,
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
	if created.JobTitle != nil {
		fmt.Printf("   Job Title: %s\n", created.JobTitle.Title)
	}
	if len(tagNames) > 0 {
		fmt.Printf("   Tags: %s\n", strings.Join(tagNames, ", "))
	}

	return nil
}

func attachTags(ctx context.Context, tagService *service.TagService, bragID int64, userID int64, tagNames []string) error {
	tagIDs := make([]int64, 0, len(tagNames))

	for _, tagName := range tagNames {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}

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
