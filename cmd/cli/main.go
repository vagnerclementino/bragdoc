package main

import (
	"log"
	"os"

	"github.com/vagnerclementino/bragdoc/cmd/cli/commands"
	"github.com/vagnerclementino/bragdoc/internal/database"
	"github.com/vagnerclementino/bragdoc/internal/repository"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

func main() {
	// Setup database
	db, err := database.SetupDatabase()
	if err != nil {
		log.Fatalf("failed to setup database: %v", err)
	}

	// Initialize repositories
	bragRepo := repository.NewBragRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	bragService := service.NewBragService(bragRepo)
	userService := service.NewUserService(userRepo)

	// Create root command with dependencies
	rootCmd := commands.NewRootCmd(bragService, userService)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
