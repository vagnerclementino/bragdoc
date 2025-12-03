package tag

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

func NewListCmd(tagService *service.TagService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tags",
		Long:  `List all tags created by the user`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(cmd.Context(), tagService, cmd)
		},
	}

	cmd.Flags().StringP("format", "f", "table", "Output format (table|json|yaml)")

	return cmd
}

func runList(ctx context.Context, tagService *service.TagService, cmd *cobra.Command) error {
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
		return outputJSON(tags)
	case "yaml":
		return outputYAML(tags)
	case "table":
		return outputTable(tags)
	default:
		return fmt.Errorf("invalid format: %s. Valid options: table, json, yaml", format)
	}
}

func outputTable(tags []*domain.Tag) error {
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

func outputJSON(tags []*domain.Tag) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tags)
}

func outputYAML(tags []*domain.Tag) error {
	encoder := yaml.NewEncoder(os.Stdout)
	defer encoder.Close()
	return encoder.Encode(tags)
}
