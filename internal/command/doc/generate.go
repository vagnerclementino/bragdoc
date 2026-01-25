package doc

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

func NewGenerateCmd(docService *service.DocumentService, bragService *service.BragService, tagService *service.TagService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a brag document",
		Long: `Generate a professional achievement document from your brags.

You can filter which brags to include using IDs, categories, or tags.
If no filters are specified, all brags will be included.

Examples:
  # Generate markdown document with all brags
  bragdoc doc generate

  # Generate and save to file
  bragdoc doc generate --output achievements.md

  # Generate with specific brags
  bragdoc doc generate --brags 1,2,3

  # Generate with brags in specific categories
  bragdoc doc generate --category project,leadership

  # Generate with brags having specific tags
  bragdoc doc generate --tags promotion,review`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(cmd.Context(), docService, bragService, tagService, cmd)
		},
	}

	cmd.Flags().StringP("format", "f", "markdown", "Document format (markdown|pdf|docx) - MVP supports only markdown")
	cmd.Flags().StringP("output", "o", "", "Output file path (if not specified, prints to stdout)")
	cmd.Flags().StringSliceP("brags", "b", []string{}, "Specific brag IDs to include (comma-separated)")
	cmd.Flags().StringSliceP("category", "c", []string{}, "Include only these categories (comma-separated)")
	cmd.Flags().StringSliceP("tags", "t", []string{}, "Include only brags with these tags (comma-separated)")
	cmd.Flags().String("template", "default", "Document template (default|executive|technical) - MVP supports only default")
	cmd.Flags().Bool("enhance-with-ai", false, "Enhance descriptions using AI (not yet implemented)")

	return cmd
}

func runGenerate(ctx context.Context, docService *service.DocumentService, bragService *service.BragService, tagService *service.TagService, cmd *cobra.Command) error {
	formatStr, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")
	bragIDs, _ := cmd.Flags().GetStringSlice("brags")
	categories, _ := cmd.Flags().GetStringSlice("category")
	tagNames, _ := cmd.Flags().GetStringSlice("tags")
	template, _ := cmd.Flags().GetString("template")
	enhanceWithAI, _ := cmd.Flags().GetBool("enhance-with-ai")

	// Parse format
	format, err := domain.ParseDocumentFormat(formatStr)
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
	}

	// MVP: Only markdown is supported
	if format != domain.FormatMarkdown {
		return fmt.Errorf("format %s not yet supported in MVP (only markdown is currently available)", format.String())
	}

	// MVP: AI enhancement not yet implemented
	if enhanceWithAI {
		return fmt.Errorf("AI enhancement not yet implemented (coming in next phase)")
	}

	// TODO: Get actual user ID from config/session
	userID := int64(1)

	// Get brags based on filters
	var brags []*domain.Brag

	if len(bragIDs) > 0 {
		// Get specific brags by IDs
		brags, err = getBragsByIDs(ctx, bragService, bragIDs)
		if err != nil {
			return fmt.Errorf("failed to get brags by IDs: %w", err)
		}
	} else if len(categories) > 0 || len(tagNames) > 0 {
		// Get brags by filters
		brags, err = getBragsByFilters(ctx, bragService, userID, categories, tagNames)
		if err != nil {
			return fmt.Errorf("failed to get brags by filters: %w", err)
		}
	} else {
		// Get all brags for user
		brags, err = bragService.List(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get brags: %w", err)
		}
	}

	if len(brags) == 0 {
		return fmt.Errorf("no brags found matching the criteria")
	}

	// Load tags for each brag
	for _, brag := range brags {
		tags, err := tagService.ListByBrag(ctx, brag.ID)
		if err != nil {
			// Don't fail, just log warning
			fmt.Fprintf(os.Stderr, "⚠️  Warning: failed to load tags for brag %d: %v\n", brag.ID, err)
		} else {
			brag.Tags = tags
		}
	}

	// Generate document
	opts := service.GenerateOptions{
		Format:        format,
		Template:      template,
		EnhanceWithAI: enhanceWithAI,
	}

	doc, err := docService.Generate(ctx, brags, userID, opts)
	if err != nil {
		return fmt.Errorf("failed to generate document: %w", err)
	}

	// Save or output document
	if output != "" {
		if err := os.WriteFile(output, doc.Content, 0600); err != nil {
			return fmt.Errorf("failed to write document: %w", err)
		}
		fmt.Printf("✅ Document generated successfully!\n")
		fmt.Printf("📄 File: %s\n", output)
		fmt.Printf("📊 Brags included: %d\n", len(brags))
		fmt.Printf("📁 Categories: %s\n", strings.Join(doc.Metadata.Categories, ", "))
	} else {
		// Output to stdout
		fmt.Print(string(doc.Content))
	}

	return nil
}

// getBragsByIDs retrieves brags by their IDs
func getBragsByIDs(ctx context.Context, bragService *service.BragService, bragIDStrs []string) ([]*domain.Brag, error) {
	brags := make([]*domain.Brag, 0, len(bragIDStrs))

	for _, idStr := range bragIDStrs {
		id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid brag ID '%s': %w", idStr, err)
		}

		brag, err := bragService.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get brag %d: %w", id, err)
		}

		brags = append(brags, brag)
	}

	return brags, nil
}

// getBragsByFilters retrieves brags filtered by categories and/or tags
func getBragsByFilters(ctx context.Context, bragService *service.BragService, userID int64, categories []string, tagNames []string) ([]*domain.Brag, error) {
	allBrags := make([]*domain.Brag, 0)
	bragMap := make(map[int64]*domain.Brag)

	// If categories specified, get brags by category
	if len(categories) > 0 {
		for _, catStr := range categories {
			category, err := domain.ParseCategory(catStr)
			if err != nil {
				return nil, fmt.Errorf("invalid category '%s': %w", catStr, err)
			}

			brags, err := bragService.SearchByCategory(ctx, userID, category)
			if err != nil {
				return nil, fmt.Errorf("failed to search brags by category %s: %w", catStr, err)
			}

			// Add to map to avoid duplicates
			for _, brag := range brags {
				bragMap[brag.ID] = brag
			}
		}
	}

	// If tags specified, get brags by tags
	if len(tagNames) > 0 {
		brags, err := bragService.SearchByTags(ctx, userID, tagNames)
		if err != nil {
			return nil, fmt.Errorf("failed to search brags by tags: %w", err)
		}

		if len(categories) > 0 {
			// If we also filtered by category, only include brags that match both
			for _, brag := range brags {
				if _, exists := bragMap[brag.ID]; exists {
					// Keep this brag as it matches both filters
				} else {
					// Remove from map as it doesn't match category filter
					delete(bragMap, brag.ID)
				}
			}
		} else {
			// No category filter, just add all tag-matched brags
			for _, brag := range brags {
				bragMap[brag.ID] = brag
			}
		}
	}

	// Convert map to slice
	for _, brag := range bragMap {
		allBrags = append(allBrags, brag)
	}

	return allBrags, nil
}
