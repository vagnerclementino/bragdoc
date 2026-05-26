package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/config"
)

// NewResetCmd creates a new command for resetting bragdoc to a clean state.
func NewResetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset bragdoc to a clean state",
		Long: `Reset bragdoc by removing the configuration directory and database.
This deletes all data including users, brags, tags, and configuration.
After reset, you can run 'bragdoc init' to start fresh.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runReset(cmd)
		},
	}

	cmd.Flags().Bool("force", false, "Skip confirmation prompt")

	return cmd
}

func runReset(cmd *cobra.Command) error {
	configManager := config.NewManager()

	if !configManager.IsInitialized() {
		fmt.Println("ℹ️  Bragdoc is not initialized. Nothing to reset.")
		return nil
	}

	force, _ := cmd.Flags().GetBool("force")

	configDir := configManager.GetConfigDir()
	fmt.Printf("⚠️  This will delete all bragdoc data:\n")
	fmt.Printf("   📁 %s\n", configDir)

	if !force {
		fmt.Print("\nAre you sure? (y/N): ")
		var answer string
		if _, err := fmt.Scanln(&answer); err != nil || (answer != "y" && answer != "Y") {
			fmt.Println("❌ Reset cancelled.")
			return nil
		}
	}

	if err := os.RemoveAll(configDir); err != nil {
		return fmt.Errorf("failed to remove bragdoc data: %w", err)
	}

	fmt.Println("✅ Bragdoc has been reset. Run 'bragdoc init' to start fresh.")
	return nil
}
