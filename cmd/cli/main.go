package main

import (
    "context"
    "log"
    "os"
    "path/filepath"
    "strings"
    "github.com/vagnerclementino/bragdoc/config"
    "github.com/vagnerclementino/bragdoc/internal/command"
    "github.com/vagnerclementino/bragdoc/internal/database"
    "github.com/vagnerclementino/bragdoc/internal/service"
)

func main() {
    cfg, err := loadConfig()
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    dbPath := getDatabasePath(cfg)

    db, err := database.New(dbPath)
    if err != nil {
        log.Fatalf("failed to open database: %v", err)
    }
    defer func(db *database.DB) {
        err := db.Close()
        if err != nil {
            log.Fatalf("failed to close database: %v", err)
        }
    }(db)

    sqliteDB := database.NewSQLiteDB(db.Conn())

    // Initialize repositories
    userRepo := database.NewUserRepository(sqliteDB)
    categoryRepo := database.NewCategoryRepository(sqliteDB)
    positionRepo := database.NewPositionRepository(sqliteDB, userRepo)
    bragRepo := database.NewBragRepository(sqliteDB, userRepo, categoryRepo, positionRepo)
    tagRepo := database.NewTagRepository(sqliteDB)

    // Initialize services
    bragService := service.NewBragService(bragRepo)
    userService := service.NewUserService(userRepo)
    tagService := service.NewTagService(tagRepo)
    docService := service.NewDocumentService(userService)

    // Create root command with dependencies
    rootCmd := command.NewRootCmd(bragService, userService, tagService, docService)
    if err := rootCmd.Execute(); err != nil {
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
