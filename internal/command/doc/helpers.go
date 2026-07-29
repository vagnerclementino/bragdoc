package doc

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/config"
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
	m := config.NewManager()
	return m.IsInitialized()
}
