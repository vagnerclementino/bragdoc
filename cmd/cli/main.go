package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/vagnerclementino/bragdoc/cmd/cli/commands"
	"github.com/vagnerclementino/bragdoc/internal/database"
	"github.com/vagnerclementino/bragdoc/internal/repository"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

func main() {
	// Setup database path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get home directory: %v", err)
	}
	dbPath := filepath.Join(homeDir, ".bragdoc", "bragdoc.db")

	// Open database connection (without running migrations)
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	bragRepo := repository.NewBragRepository(db.Conn())
	userRepo := repository.NewUserRepository(db.Conn())
	tagRepo := repository.NewTagRepository(db.Conn())

	// Initialize services
	bragService := service.NewBragService(bragRepo)
	userService := service.NewUserService(userRepo)
	tagService := service.NewTagService(tagRepo)

	// Create root command with dependencies
	rootCmd := commands.NewRootCmd(bragService, userService, tagService)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
