package brag

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

// NewRemoveCmd creates a new command for removing brag entries.
func NewRemoveCmd(bragService *service.BragService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <ids>",
		Short: "Remove brag entries",
		Long: `Remove one or more brag entries by ID.
Supports multiple IDs and ranges:
  - Single ID: bragdoc brag remove 1
  - Multiple IDs: bragdoc brag remove 1,2,3
  - Range: bragdoc brag remove 1-5
  - Combined: bragdoc brag remove 1,3,5-8`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemove(cmd.Context(), bragService, cmd, args)
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	return cmd
}

func runRemove(ctx context.Context, bragService *service.BragService, cmd *cobra.Command, args []string) error {
	// Parse IDs
	ids, err := parseIDs(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse IDs: %w", err)
	}

	if len(ids) == 0 {
		return fmt.Errorf("no valid IDs provided")
	}

	// Get confirmation unless --force flag is set
	force, _ := cmd.Flags().GetBool("force")
	if !force {
		fmt.Printf("About to remove %d brag(s) with IDs: %v\n", len(ids), ids)
		fmt.Print("Are you sure? (y/N): ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read confirmation: %w", err)
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			fmt.Println("Operation cancelled.")
			return nil
		}
	}

	// Remove brags
	var failed []int64
	var succeeded []int64

	for _, id := range ids {
		if err := bragService.Delete(ctx, id); err != nil {
			failed = append(failed, id)
			fmt.Fprintf(os.Stderr, "❌ Failed to remove brag %d: %v\n", id, err)
		} else {
			succeeded = append(succeeded, id)
		}
	}

	// Report results
	if len(succeeded) > 0 {
		fmt.Printf("✅ Successfully removed %d brag(s): %v\n", len(succeeded), succeeded)
	}

	if len(failed) > 0 {
		return fmt.Errorf("failed to remove %d brag(s): %v", len(failed), failed)
	}

	return nil
}

// parseIDs parses a string containing IDs and ranges into a slice of int64
// Supports formats like: "1", "1,2,3", "1-5", "1,3,5-8"
func parseIDs(idsStr string) ([]int64, error) {
	if idsStr == "" {
		return nil, fmt.Errorf("no IDs provided")
	}

	var result []int64
	parts := strings.Split(idsStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if strings.Contains(part, "-") {
			// Handle range (e.g., "1-5")
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}

			start, err := strconv.ParseInt(strings.TrimSpace(rangeParts[0]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid start ID in range %s: %w", part, err)
			}

			end, err := strconv.ParseInt(strings.TrimSpace(rangeParts[1]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid end ID in range %s: %w", part, err)
			}

			if start > end {
				return nil, fmt.Errorf("invalid range %s: start must be <= end", part)
			}

			if end-start > 1000 {
				return nil, fmt.Errorf("invalid range %s: range too large (max 1000)", part)
			}

			for i := start; i <= end; i++ {
				result = append(result, i)
			}
		} else {
			// Handle single ID
			id, err := strconv.ParseInt(part, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid ID: %s", part)
			}
			result = append(result, id)
		}
	}

	// Remove duplicates and sort
	seen := make(map[int64]bool)
	var unique []int64
	for _, id := range result {
		if !seen[id] {
			seen[id] = true
			unique = append(unique, id)
		}
	}

	sort.Slice(unique, func(i, j int) bool {
		return unique[i] < unique[j]
	})

	return unique, nil
}
