package commands

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vagnerclementino/bragdoc/config"
	"github.com/vagnerclementino/bragdoc/internal/database"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/repository"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

// TestBragCommands_E2E tests the complete workflow of brag commands
func TestBragCommands_E2E(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	ctx := context.Background()

	// Initialize bragdoc
	configDir := filepath.Join(tempDir, ".bragdoc")
	dbPath := filepath.Join(configDir, "bragdoc.db")

	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	// Create and setup database
	db, err := database.New(dbPath)
	require.NoError(t, err)
	defer db.Close()

	err = db.Migrate(ctx)
	require.NoError(t, err)

	// Create config
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Path: dbPath,
		},
	}

	manager := config.NewManager()
	err = manager.Initialize(ctx, cfg, config.FormatYAML)
	require.NoError(t, err)

	// Create user in database
	userRepo := repository.NewUserRepository(db.Conn())
	userService := service.NewUserService(userRepo)
	_, err = userService.Create(ctx, &domain.User{
		Name:   "Test User",
		Email:  "test@example.com",
		Locale: "en-US",
	})
	require.NoError(t, err)

	// Setup services
	bragRepo := repository.NewBragRepository(db.Conn())
	tagRepo := repository.NewTagRepository(db.Conn())
	bragService := service.NewBragService(bragRepo)
	tagService := service.NewTagService(tagRepo)

	t.Run("Complete brag workflow", func(t *testing.T) {
		// Test 1: Add a brag
		t.Run("add brag with tags", func(t *testing.T) {
			cmd := NewBragAddCmd(bragService, tagService)
			cmd.SetArgs([]string{
				"--title", "Test Achievement",
				"--description", "This is a test achievement for e2e testing",
				"--category", "achievement",
				"--tags", "test,e2e,automation",
			})

			err := cmd.Execute()
			assert.NoError(t, err)
			
			// Verify brag was created by checking database
			brags, err := bragService.List(ctx, 1)
			assert.NoError(t, err)
			assert.Len(t, brags, 1)
			assert.Equal(t, "Test Achievement", brags[0].Title)
		})

		// Test 2: List brags
		t.Run("list brags", func(t *testing.T) {
			cmd := NewBragListCmd(bragService, tagService)
			cmd.SetArgs([]string{"--format", "table"})

			err := cmd.Execute()
			assert.NoError(t, err)
		})

		// Test 3: Show brag
		t.Run("show brag", func(t *testing.T) {
			cmd := NewBragShowCmd(bragService, tagService)
			cmd.SetArgs([]string{"1"})

			err := cmd.Execute()
			assert.NoError(t, err)
		})

		// Test 4: Edit brag
		t.Run("edit brag", func(t *testing.T) {
			cmd := NewBragEditCmd(bragService, tagService)
			cmd.SetArgs([]string{
				"1",
				"--title", "Updated Test Achievement",
				"--category", "skill",
			})

			err := cmd.Execute()
			assert.NoError(t, err)
			
			// Verify edit
			brag, err := bragService.GetByID(ctx, 1)
			assert.NoError(t, err)
			assert.Equal(t, "Updated Test Achievement", brag.Title)
			assert.Equal(t, domain.CategorySkill, brag.Category)
		})

		// Test 5: Verify edit persisted
		t.Run("verify edit persisted", func(t *testing.T) {
			brag, err := bragService.GetByID(ctx, 1)
			assert.NoError(t, err)
			assert.Equal(t, "Updated Test Achievement", brag.Title)
			assert.Equal(t, domain.CategorySkill, brag.Category)
		})

		// Test 6: Add more brags for filtering
		t.Run("add multiple brags", func(t *testing.T) {
			brags := []struct {
				title       string
				description string
				category    string
				tags        string
			}{
				{
					title:       "Leadership Project",
					description: "Led a team of developers to success",
					category:    "leadership",
					tags:        "leadership,team,project",
				},
				{
					title:       "Innovation Initiative",
					description: "Implemented a new innovative solution",
					category:    "innovation",
					tags:        "innovation,creative",
				},
			}

			for _, brag := range brags {
				cmd := NewBragAddCmd(bragService, tagService)
				cmd.SetArgs([]string{
					"--title", brag.title,
					"--description", brag.description,
					"--category", brag.category,
					"--tags", brag.tags,
				})

				var buf bytes.Buffer
				cmd.SetOut(&buf)
				cmd.SetErr(&buf)

				err := cmd.Execute()
				assert.NoError(t, err)
			}
		})

		// Test 7: Filter by category
		t.Run("filter by category", func(t *testing.T) {
			cmd := NewBragListCmd(bragService, tagService)
			cmd.SetArgs([]string{"--category", "leadership"})

			err := cmd.Execute()
			assert.NoError(t, err)
			
			// Verify filter works
			brags, err := bragService.SearchByCategory(ctx, 1, domain.CategoryLeadership)
			assert.NoError(t, err)
			assert.Len(t, brags, 1)
			assert.Equal(t, "Leadership Project", brags[0].Title)
		})

		// Test 8: Filter by tags
		t.Run("filter by tags", func(t *testing.T) {
			cmd := NewBragListCmd(bragService, tagService)
			cmd.SetArgs([]string{"--tags", "innovation"})

			err := cmd.Execute()
			assert.NoError(t, err)
			
			// Verify filter works
			brags, err := bragService.SearchByTags(ctx, 1, []string{"innovation"})
			assert.NoError(t, err)
			assert.Len(t, brags, 1)
			assert.Equal(t, "Innovation Initiative", brags[0].Title)
		})

		// Test 9: Show multiple brags
		t.Run("show multiple brags", func(t *testing.T) {
			cmd := NewBragShowCmd(bragService, tagService)
			cmd.SetArgs([]string{"1,2"})

			err := cmd.Execute()
			assert.NoError(t, err)
		})

		// Test 10: Show range of brags
		t.Run("show range of brags", func(t *testing.T) {
			cmd := NewBragShowCmd(bragService, tagService)
			cmd.SetArgs([]string{"1-3"})

			err := cmd.Execute()
			assert.NoError(t, err)
		})

		// Test 11: Remove brag with force flag
		t.Run("remove brag with force", func(t *testing.T) {
			cmd := NewBragRemoveCmd(bragService)
			cmd.SetArgs([]string{"3", "--force"})

			err := cmd.Execute()
			assert.NoError(t, err)
			
			// Verify removal
			_, err = bragService.GetByID(ctx, 3)
			assert.Error(t, err) // Should not find deleted brag
		})

		// Test 12: Verify removal persisted
		t.Run("verify removal persisted", func(t *testing.T) {
			brags, err := bragService.List(ctx, 1)
			assert.NoError(t, err)
			assert.Len(t, brags, 2) // Should have 2 brags left
		})

		// Test 13: JSON output format
		t.Run("json output format", func(t *testing.T) {
			cmd := NewBragListCmd(bragService, tagService)
			cmd.SetArgs([]string{"--format", "json"})

			err := cmd.Execute()
			assert.NoError(t, err)
		})

		// Test 14: YAML output format
		t.Run("yaml output format", func(t *testing.T) {
			cmd := NewBragListCmd(bragService, tagService)
			cmd.SetArgs([]string{"--format", "yaml"})

			err := cmd.Execute()
			assert.NoError(t, err)
		})
	})
}

// TestBragCommands_E2E_InitializationCheck tests that commands fail without initialization
// This test is covered by TestRequiresInitialization in initialization_check_test.go
// and by manual testing, so we don't need to duplicate it here
