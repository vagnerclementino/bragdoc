package brag

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/service"
	"gopkg.in/yaml.v3"
)

func NewListCmd(bragService *service.BragService, tagService *service.TagService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List brag entries",
		Long:  `List all your documented professional achievements with optional filters`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(cmd.Context(), bragService, tagService, cmd)
		},
	}

	cmd.Flags().StringP("category", "c", "", "Filter by category (project|achievement|skill|leadership|innovation)")
	cmd.Flags().StringSliceP("tags", "t", []string{}, "Filter by tags (comma-separated)")
	cmd.Flags().StringP("format", "f", "table", "Output format (table|json|yaml)")
	cmd.Flags().IntP("limit", "l", 50, "Maximum number of results")

	return cmd
}

func runList(ctx context.Context, bragService *service.BragService, tagService *service.TagService, cmd *cobra.Command) error {
	categoryStr, _ := cmd.Flags().GetString("category")
	tagNames, _ := cmd.Flags().GetStringSlice("tags")
	format, _ := cmd.Flags().GetString("format")
	limit, _ := cmd.Flags().GetInt("limit")

	// TODO: Get actual user ID from config/session
	userID := int64(1)

	var brags []*domain.Brag
	var err error

	// Apply filters
	if categoryStr != "" {
		category, err := domain.ParseCategory(categoryStr)
		if err != nil {
			return fmt.Errorf("invalid category: %w. Valid options: project, achievement, skill, leadership, innovation", err)
		}
		brags, err = bragService.SearchByCategory(ctx, userID, category)
		if err != nil {
			return fmt.Errorf("failed to search brags by category: %w", err)
		}
	} else if len(tagNames) > 0 {
		brags, err = bragService.SearchByTags(ctx, userID, tagNames)
		if err != nil {
			return fmt.Errorf("failed to search brags by tags: %w", err)
		}
	} else {
		brags, err = bragService.List(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to list brags: %w", err)
		}
	}

	// Apply limit
	if limit > 0 && len(brags) > limit {
		brags = brags[:limit]
	}

	for _, brag := range brags {
		tags, err := tagService.ListByBrag(ctx, brag.ID)
		if err != nil {
			// Don't fail, just log warning
			fmt.Fprintf(os.Stderr, "⚠️  Warning: failed to load tags for brag %d: %v\n", brag.ID, err)
		} else {
			brag.Tags = tags
		}
	}

	if len(brags) == 0 {
		fmt.Println("No brags found.")
		return nil
	}

	switch strings.ToLower(format) {
	case "json":
		return outputJSON(brags)
	case "yaml":
		return outputYAML(brags)
	case "table":
		return outputTable(brags)
	default:
		return fmt.Errorf("invalid format: %s. Valid options: table, json, yaml", format)
	}
}

func outputTable(brags []*domain.Brag) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() {
		if err := w.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to flush output: %v\n", err)
		}
	}()

	if _, err := fmt.Fprintln(w, "ID\tTITLE\tCATEGORY\tTAGS\tCREATED"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "--\t-----\t--------\t----\t-------"); err != nil {
		return err
	}

	for _, brag := range brags {
		tagNames := make([]string, len(brag.Tags))
		for i, tag := range brag.Tags {
			tagNames[i] = tag.Name
		}
		tagsStr := strings.Join(tagNames, ", ")
		if tagsStr == "" {
			tagsStr = "-"
		}

		// Truncate title if too long
		title := brag.Title
		if len(title) > 50 {
			title = title[:47] + "..."
		}

		if _, err := fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			brag.ID,
			title,
			brag.Category.String(),
			tagsStr,
			brag.CreatedAt.Format("2006-01-02"),
		); err != nil {
			return err
		}
	}

	return nil
}

func outputJSON(brags []*domain.Brag) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(brags)
}

func outputYAML(brags []*domain.Brag) error {
	encoder := yaml.NewEncoder(os.Stdout)
	defer func() {
		if err := encoder.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close YAML encoder: %v\n", err)
		}
	}()
	return encoder.Encode(brags)
}
