package command

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/config"
	"github.com/vagnerclementino/bragdoc/internal/database"
	"github.com/vagnerclementino/bragdoc/internal/database/queries"
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize Bragdoc configuration and database",
		Long:  `Initialize Bragdoc by creating the configuration directory and setting up the database with migrations`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(cmd.Context(), cmd)
		},
	}

	// Required flags
	cmd.Flags().StringP("name", "n", "", "Your full name (required)")
	cmd.Flags().StringP("email", "e", "", "Your email (required)")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("email")

	// Optional flags
	cmd.Flags().StringP("job-title", "j", "", "Your job title (optional)")
	cmd.Flags().StringP("company", "c", "", "Your company (optional)")
	cmd.Flags().StringP("locale", "l", "en-US", "Locale (language-COUNTRY): en-US or pt-BR")

	return cmd
}

func runInit(ctx context.Context, cmd *cobra.Command) error {
	fmt.Println("🚀 Initializing Bragdoc...")

	// Create config manager
	configManager := config.NewManager()

	// Check if already initialized
	if configManager.IsInitialized() {
		fmt.Println("⚠️  Bragdoc is already initialized!")
		fmt.Printf("📁 Configuration: %s\n", configManager.GetConfigPath())
		fmt.Printf("🗄️  Database: %s\n", configManager.GetDatabasePath())
		return nil
	}

	// Get user information from flags
	name, _ := cmd.Flags().GetString("name")
	email, _ := cmd.Flags().GetString("email")
	jobTitle, _ := cmd.Flags().GetString("job-title")
	company, _ := cmd.Flags().GetString("company")
	locale, _ := cmd.Flags().GetString("locale")

	// Validate locale before proceeding
	if locale != "en-US" && locale != "pt-BR" {
		return fmt.Errorf("invalid locale: %s (supported: en-US, pt-BR)", locale)
	}

	// Generate default configuration (user data will be stored in database)
	defaultConfig := configManager.GetDefaultConfig()

	// Create configuration file (YAML only in v1)
	if err := configManager.Initialize(ctx, defaultConfig, config.FormatYAML); err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	// Setup database
	dbPath := configManager.GetDatabasePath()
	db, err := database.New(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(ctx); err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Create user in database
	q := queries.New(db.Conn())
	createdUser, err := q.CreateUser(ctx, queries.CreateUserParams{
		Name:     name,
		Email:    email,
		JobTitle: newNullString(jobTitle),
		Company:  newNullString(company),
		Locale:   locale,
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Println("✅ Bragdoc initialized successfully!")
	fmt.Printf("📁 Configuration: %s\n", configManager.GetConfigPath())
	fmt.Printf("🗄️  Database: %s\n", dbPath)
	fmt.Printf("👤 User created: %s (ID: %d)\n", createdUser.Name, createdUser.ID)
	fmt.Println("\n💡 Next steps:")
	fmt.Println("   - Use 'bragdoc brag add' to create your first brag")
	fmt.Println("   - Use 'bragdoc brag list' to view your brags")

	return nil
}

// newNullString creates a sql.NullString from a string
// Returns a valid NullString if the string is not empty, otherwise returns an invalid NullString
func newNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
