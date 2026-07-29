// Package main is the entry point for the bragdoc MCP server binary.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vagnerclementino/bragdoc/config"
	"github.com/vagnerclementino/bragdoc/internal/database"
	mcpserver "github.com/vagnerclementino/bragdoc/internal/mcp"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	dbPath := getDatabasePath(cfg)

	db, err := database.New(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close database: %v\n", err)
		}
	}()

	ctx := context.Background()
	if err := db.Migrate(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	sqliteDB := database.NewSQLiteDB(db.Conn())

	// Initialize repositories
	userRepo := database.NewUserRepository(sqliteDB)
	categoryRepo := database.NewCategoryRepository(sqliteDB)
	jobTitleRepo := database.NewJobTitleRepository(sqliteDB, userRepo)
	bragRepo := database.NewBragRepository(sqliteDB, userRepo, categoryRepo, jobTitleRepo)
	tagRepo := database.NewTagRepository(sqliteDB)

	// Initialize services
	bragService := service.NewBragService(bragRepo)
	userService := service.NewUserService(userRepo)
	tagService := service.NewTagService(tagRepo)
	jobTitleService := service.NewJobTitleService(jobTitleRepo)
	docService := service.NewDocumentService(userService)

	// Create and run MCP server
	srv := mcpserver.NewServer(bragService, tagService, userService, docService, jobTitleService)
	if err := srv.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
		os.Exit(1)
	}
}

func loadConfig() (*config.Config, error) {
	mgr := config.NewManager()
	cfg, err := mgr.Load(context.Background())
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func getDatabasePath(cfg *config.Config) string {
	if cfg.Database.Path != "" {
		return expandPath(cfg.Database.Path)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".bragdoc", "bragdoc.db")
}

func expandPath(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if path == "~" {
		return homeDir
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}

	return path
}
