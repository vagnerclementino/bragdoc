package commands

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

func NewTagListCmd(tagService *service.TagService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tags",
		Long:  `List all tags created by the user`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagList(cmd.Context(), tagService, cmd)
		},
	}

	cmd.Flags().StringP("format", "f", "table", "Output format (table|json|yaml)")

	return cmd
}

func runTagList(ctx context.Context, tagService *service.TagService, cmd *cobra.Command) error {
	format, _ := cmd.Flags().GetString("format")

	// TODO: Get actual user ID from config/session
	userID := int64(1)

	tags, err := tagService.ListByUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}

	if len(tags) == 0 {
		fmt.Println("No tags found.")
		return nil
	}

	// Format output
	switch strings.ToLower(format) {
	case "json":
		return outputTagsJSON(tags)
	case "yaml":
		return outputTagsYAML(tags)
	case "table":
		return outputTagsTable(tags)
	default:
		return fmt.Errorf("invalid format: %s. Valid options: table, json, yaml", format)
	}
}

func outputTagsTable(tags []*domain.Tag) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tCREATED")
	fmt.Fprintln(w, "--\t----\t-------")

	for _, tag := range tags {
		fmt.Fprintf(w, "%d\t%s\t%s\n",
			tag.ID,
			tag.Name,
			tag.CreatedAt.Format("2006-01-02"),
		)
	}

	return w.Flush()
}

func outputTagsJSON(tags []*domain.Tag) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tags)
}

func outputTagsYAML(tags []*domain.Tag) error {
	encoder := yaml.NewEncoder(os.Stdout)
	defer encoder.Close()
	return encoder.Encode(tags)
}
