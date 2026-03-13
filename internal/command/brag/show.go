package brag

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

// NewShowCmd creates a new command for showing detailed brag information.
func NewShowCmd(bragService *service.BragService, tagService *service.TagService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <ids>",
		Short: "Show detailed information about brag entries",
		Long: `Show detailed information about one or more brag entries by ID.
Supports multiple IDs and ranges:
  - Single ID: bragdoc brag show 1
  - Multiple IDs: bragdoc brag show 1,2,3
  - Range: bragdoc brag show 1-5`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShow(cmd.Context(), bragService, tagService, args)
		},
	}

	return cmd
}

func runShow(ctx context.Context, bragService *service.BragService, tagService *service.TagService, args []string) error {
	// Parse IDs
	ids, err := parseIDs(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse IDs: %w", err)
	}

	if len(ids) == 0 {
		return fmt.Errorf("no valid IDs provided")
	}

	// Fetch and display each brag
	for i, id := range ids {
		if i > 0 {
			fmt.Println("\n" + strings.Repeat("-", 80))
		}

		brag, err := bragService.GetByID(ctx, id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to get brag %d: %v\n", id, err)
			continue
		}

		// Load tags
		tags, err := tagService.ListByBrag(ctx, brag.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Warning: failed to load tags for brag %d: %v\n", brag.ID, err)
		} else {
			brag.Tags = tags
		}

		displayBrag(brag)
	}

	return nil
}

func displayBrag(brag *domain.Brag) {
	fmt.Printf("\n📝 Brag #%d\n", brag.ID)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Title:       %s\n", brag.Title)
	fmt.Printf("Category:    %s\n", brag.Category.String())

	if brag.JobTitle != nil {
		fmt.Printf("Job Title:   %s", brag.JobTitle.Title)
		if brag.JobTitle.Company != "" {
			fmt.Printf(" at %s", brag.JobTitle.Company)
		}
		if brag.JobTitle.StartDate != nil {
			fmt.Printf(" (since %s)", brag.JobTitle.StartDate.Format("Jan 2006"))
		}
		fmt.Println()
	}

	fmt.Printf("Created:     %s\n", brag.CreatedAt.Format("2006-01-02 15:04:05"))
	if !brag.UpdatedAt.IsZero() {
		fmt.Printf("Updated:     %s\n", brag.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	if len(brag.Tags) > 0 {
		tagNames := make([]string, len(brag.Tags))
		for i, tag := range brag.Tags {
			tagNames[i] = tag.Name
		}
		fmt.Printf("Tags:        %s\n", strings.Join(tagNames, ", "))
	}

	fmt.Println("\nDescription:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println(brag.Description)
	fmt.Println(strings.Repeat("-", 80))
}
