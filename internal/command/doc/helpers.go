package doc

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// requiresInitialization returns a PreRunE function that checks if bragdoc is initialized
func requiresInitialization() func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, _ []string) error {
		if !isInitialized() {
			return fmt.Errorf(`bragdoc is not initialized. Please run 'bragdoc init' first.

Example:
  bragdoc init --name "Your Name" --email "your.email@example.com"`)
		}
		return nil
	}
}

// isInitialized checks if bragdoc has been initialized by looking for the config directory
func isInitialized() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	configDir := filepath.Join(homeDir, ".bragdoc")
	dbPath := filepath.Join(configDir, "bragdoc.db")

	// Check if both config directory and database exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return false
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return false
	}

	return true
}
