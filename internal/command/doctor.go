// Package command provides diagnostic and debugging commands
package command

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vagnerclementino/bragdoc/config"
	"github.com/vagnerclementino/bragdoc/internal/database"
)

// NewDoctorCmd creates a hidden command for debugging and fixing issues
func NewDoctorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "doctor",
		Short:  "Diagnose and fix bragdoc issues",
		Long:   `Advanced commands for diagnosing and fixing bragdoc database and configuration issues.`,
		Hidden: true, // Hidden from normal help
	}

	cmd.AddCommand(newDoctorMigrateCmd())
	cmd.AddCommand(newDoctorCheckCmd())

	return cmd
}

// newDoctorMigrateCmd creates the migrate subcommand
func newDoctorMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Manage database migrations",
		Long:  `Advanced migration management. Migrations run automatically, use this only for debugging.`,
	}

	cmd.AddCommand(newMigrateStatusCmd())
	cmd.AddCommand(newMigrateUpCmd())
	cmd.AddCommand(newMigrateDownCmd())
	cmd.AddCommand(newMigrateForceCmd())

	return cmd
}

// newDoctorCheckCmd creates the check subcommand
func newDoctorCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Check bragdoc health",
		Long:  `Verify database integrity, configuration, and overall health of bragdoc.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runDoctorCheck(cmd.Context())
		},
	}
}

func newMigrateStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current migration status",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runMigrateStatus(cmd.Context())
		},
	}
}

func newMigrateUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Apply all pending migrations",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runMigrateUp(cmd.Context())
		},
	}
}

func newMigrateDownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Rollback the last migration",
		Long:  `Rollback the most recently applied migration. Use with caution!`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runMigrateDown(cmd.Context())
		},
	}
}

func newMigrateForceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "force [VERSION]",
		Short: "Force migration version (dangerous!)",
		Long: `Force the migration version without running migrations.
This is useful for recovering from dirty state, but should be used with extreme caution.
Only use this if you know what you're doing!`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var version int
			if _, err := fmt.Sscanf(args[0], "%d", &version); err != nil {
				return fmt.Errorf("invalid version: %s", args[0])
			}
			return runMigrateForce(cmd.Context(), version)
		},
	}
	return cmd
}

func runDoctorCheck(ctx context.Context) error {
	fmt.Println("🔍 Running bragdoc health check...")
	fmt.Println()

	configManager := config.NewManager()

	// Check 1: Initialization
	fmt.Print("✓ Checking initialization... ")
	if !configManager.IsInitialized() {
		fmt.Println("❌ FAILED")
		fmt.Println("  → Run 'bragdoc init' to initialize")
		return fmt.Errorf("bragdoc not initialized")
	}
	fmt.Println("✅ OK")

	dbPath := configManager.GetDatabasePath()

	// Check 2: Database connection + integrity
	fmt.Print("✓ Checking database connection... ")
	db, err := database.New(dbPath)
	if err != nil {
		fmt.Println("❌ FAILED")
		fmt.Printf("  → Error: %v\n", err)
		return err
	}
	fmt.Println("✅ OK")

	fmt.Print("✓ Checking database integrity... ")
	if err := checkDatabaseIntegrity(ctx, db); err != nil {
		fmt.Println("❌ FAILED")
		fmt.Printf("  → Error: %v\n", err)
	} else {
		fmt.Println("✅ OK")
	}
	if err := db.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to close database: %v\n", err)
	}

	// Check 3: Migration status (opens its own connection via newMigrate)
	fmt.Print("✓ Checking migration status... ")
	db2, err := database.New(dbPath)
	if err != nil {
		fmt.Println("❌ FAILED")
		fmt.Printf("  → Error: %v\n", err)
		return err
	}
	version, dirty, err := db2.MigrateVersion()
	if closeErr := db2.Close(); closeErr != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to close database: %v\n", closeErr)
	}
	if err != nil {
		fmt.Println("⚠️  WARNING")
		fmt.Printf("  → Error: %v\n", err)
	} else if dirty {
		fmt.Println("❌ DIRTY")
		fmt.Printf("  → Database is in dirty state (version %d)\n", version)
		fmt.Println("  → Run 'bragdoc doctor migrate force <version>' to fix")
	} else {
		fmt.Printf("✅ OK (version %d)\n", version)
	}

	fmt.Println("\n✅ Health check complete!")
	return nil
}

func runMigrateStatus(_ context.Context) error {
	db, err := openDatabase()
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close database: %v\n", closeErr)
		}
	}()

	version, dirty, err := db.MigrateVersion()
	if err != nil {
		return err
	}

	fmt.Println("📊 Migration Status:")
	fmt.Printf("   Version: %d\n", version)
	fmt.Printf("   Dirty: %v\n", dirty)

	if dirty {
		fmt.Println("\n⚠️  WARNING: Database is in dirty state!")
		fmt.Println("   A migration failed partway through.")
		fmt.Println("   Options:")
		fmt.Println("   1. Fix the issue manually in the database")
		fmt.Println("   2. Use 'bragdoc doctor migrate force <version>' to mark as clean")
		fmt.Println("   3. Restore from backup if available")
	} else {
		fmt.Println("\n✅ Database is clean and up to date")
	}

	return nil
}

func runMigrateUp(ctx context.Context) error {
	db, err := openDatabase()
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close database: %v\n", closeErr)
		}
	}()

	fmt.Println("🔄 Applying pending migrations...")

	if err := db.Migrate(ctx); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	version, _, err := db.MigrateVersion()
	if err != nil {
		return err
	}

	fmt.Printf("✅ Migrations applied successfully!\n")
	fmt.Printf("📊 Current version: %d\n", version)
	return nil
}

func runMigrateDown(ctx context.Context) error {
	db, err := openDatabase()
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close database: %v\n", closeErr)
		}
	}()

	fmt.Println("⚠️  WARNING: This will rollback the last migration!")
	fmt.Println("⏪ Rolling back...")

	if err := db.MigrateDown(ctx); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	version, _, err := db.MigrateVersion()
	if err != nil {
		return err
	}

	fmt.Printf("✅ Rollback complete!\n")
	fmt.Printf("📊 Current version: %d\n", version)
	return nil
}

func runMigrateForce(_ context.Context, version int) error {
	db, err := openDatabase()
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close database: %v\n", closeErr)
		}
	}()

	fmt.Printf("⚠️  WARNING: Forcing migration version to %d\n", version)
	fmt.Println("   This does NOT run migrations, only updates the version marker!")
	fmt.Println("   Use this only if you know what you're doing!")

	if err := db.MigrateForce(version); err != nil {
		return fmt.Errorf("force failed: %w", err)
	}

	fmt.Printf("\n✅ Version forced to %d\n", version)
	fmt.Println("💡 Run 'bragdoc doctor migrate status' to verify")
	return nil
}

func openDatabase() (*database.DB, error) {
	configManager := config.NewManager()

	if !configManager.IsInitialized() {
		return nil, fmt.Errorf("bragdoc not initialized. Run 'bragdoc init' first")
	}

	dbPath := configManager.GetDatabasePath()
	return database.New(dbPath)
}

func checkDatabaseIntegrity(ctx context.Context, db *database.DB) error {
	// Check if required tables exist
	requiredTables := []string{
		"users",
		"brags",
		"tags",
		"brag_tags",
		"categories",
		"job_titles",
	}

	for _, table := range requiredTables {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='%s'", table)
		if err := db.Conn().QueryRowContext(ctx, query).Scan(&count); err != nil {
			return fmt.Errorf("failed to check table %s: %w", table, err)
		}
		if count == 0 {
			return fmt.Errorf("required table '%s' not found", table)
		}
	}

	return nil
}
