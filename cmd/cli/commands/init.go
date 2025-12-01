package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/internal/database"
)

func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize Bragdoc configuration and database",
		Long:  `Initialize Bragdoc by creating the configuration directory and setting up the database with migrations`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runInit(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func runInit() error {
	fmt.Println("🚀 Initializing Bragdoc...")

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create .bragdoc directory
	bragdocDir := filepath.Join(homeDir, ".bragdoc")
	if err := os.MkdirAll(bragdocDir, 0755); err != nil {
		return fmt.Errorf("failed to create .bragdoc directory: %w", err)
	}

	// Database path
	dbPath := filepath.Join(bragdocDir, "bragdoc.db")

	// Check if database already exists
	if _, err := os.Stat(dbPath); err == nil {
		fmt.Println("⚠️  Bragdoc is already initialized!")
		fmt.Printf("📁 Configuration directory: %s\n", bragdocDir)
		fmt.Printf("🗄️  Database: %s\n", dbPath)
		return nil
	}

	// Create database and run migrations
	db, err := database.New(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	defer db.Close()

	// Run migrations (silently)
	if err := db.Migrate(context.Background()); err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	fmt.Println("✅ Bragdoc initialized successfully!")
	fmt.Printf("📁 Configuration directory: %s\n", bragdocDir)
	fmt.Printf("🗄️  Database: %s\n", dbPath)
	fmt.Println("\n💡 Next steps:")
	fmt.Println("   - Use 'bragdoc brag add' to create your first brag")
	fmt.Println("   - Use 'bragdoc brag list' to view your brags")

	return nil
}
