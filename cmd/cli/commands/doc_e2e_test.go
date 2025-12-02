package commands

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vagnerclementino/bragdoc/config"
	"github.com/vagnerclementino/bragdoc/internal/database"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/repository"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

// TestDocCommands_E2E tests the complete workflow of doc commands
func TestDocCommands_E2E(t *testing.T) {
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
	user, err := userService.Create(ctx, &domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Locale:   "en-US",
		JobTitle: "Software Engineer",
		Company:  "Test Corp",
	})
	require.NoError(t, err)

	// Setup services
	bragRepo := repository.NewBragRepository(db.Conn())
	tagRepo := repository.NewTagRepository(db.Conn())
	bragService := service.NewBragService(bragRepo)
	tagService := service.NewTagService(tagRepo)
	docService := service.NewDocumentService(userService)

	// Create test brags
	brag1, err := bragService.Create(ctx, &domain.Brag{
		OwnerID:     user.ID,
		Title:       "Test Achievement",
		Description: "This is a test achievement for document generation",
		Category:    domain.CategoryAchievement,
	})
	require.NoError(t, err)

	brag2, err := bragService.Create(ctx, &domain.Brag{
		OwnerID:     user.ID,
		Title:       "Leadership Project",
		Description: "Led a team to deliver a critical project on time",
		Category:    domain.CategoryLeadership,
	})
	require.NoError(t, err)

	// Create and attach tags
	tag1, err := tagService.Create(ctx, &domain.Tag{
		Name:    "test",
		OwnerID: user.ID,
	})
	require.NoError(t, err)

	tag2, err := tagService.Create(ctx, &domain.Tag{
		Name:    "leadership",
		OwnerID: user.ID,
	})
	require.NoError(t, err)

	err = tagService.AttachToBrag(ctx, brag1.ID, []int64{tag1.ID})
	require.NoError(t, err)

	err = tagService.AttachToBrag(ctx, brag2.ID, []int64{tag2.ID})
	require.NoError(t, err)

	t.Run("Document generation workflow", func(t *testing.T) {
		// Test 1: Generate document with all brags
		t.Run("generate document with all brags", func(t *testing.T) {
			outputFile := filepath.Join(tempDir, "test-doc.md")
			
			cmd := NewDocGenerateCmd(docService, bragService, tagService)
			cmd.SetArgs([]string{
				"--output", outputFile,
			})

			err := cmd.Execute()
			assert.NoError(t, err)

			// Verify file was created
			assert.FileExists(t, outputFile)

			// Read and verify content
			content, err := os.ReadFile(outputFile)
			require.NoError(t, err)

			contentStr := string(content)
			
			// Verify header
			assert.Contains(t, contentStr, "Professional Achievements")
			assert.Contains(t, contentStr, "Test User")
			assert.Contains(t, contentStr, "Software Engineer")
			assert.Contains(t, contentStr, "Test Corp")
			
			// Verify summary
			assert.Contains(t, contentStr, "2 professional achievements")
			
			// Verify brags
			assert.Contains(t, contentStr, "Test Achievement")
			assert.Contains(t, contentStr, "This is a test achievement for document generation")
			assert.Contains(t, contentStr, "Leadership Project")
			assert.Contains(t, contentStr, "Led a team to deliver a critical project on time")
			
			// Verify tags
			assert.Contains(t, contentStr, "test")
			assert.Contains(t, contentStr, "leadership")
			
			// Verify categories
			assert.Contains(t, contentStr, "Achievement")
			assert.Contains(t, contentStr, "Leadership")
			
			// Verify footer
			assert.Contains(t, contentStr, "Bragdoc CLI")
		})

		// Test 2: Generate document with specific brags
		t.Run("generate document with specific brags", func(t *testing.T) {
			outputFile := filepath.Join(tempDir, "test-doc-filtered.md")
			
			cmd := NewDocGenerateCmd(docService, bragService, tagService)
			cmd.SetArgs([]string{
				"--output", outputFile,
				"--brags", "1",
			})

			err := cmd.Execute()
			assert.NoError(t, err)

			// Verify file was created
			assert.FileExists(t, outputFile)

			// Read and verify content
			content, err := os.ReadFile(outputFile)
			require.NoError(t, err)

			contentStr := string(content)
			assert.Contains(t, contentStr, "Test Achievement")
			assert.NotContains(t, contentStr, "Leadership Project")
		})

		// Test 3: Generate document filtered by category
		t.Run("generate document filtered by category", func(t *testing.T) {
			outputFile := filepath.Join(tempDir, "test-doc-category.md")
			
			cmd := NewDocGenerateCmd(docService, bragService, tagService)
			cmd.SetArgs([]string{
				"--output", outputFile,
				"--category", "leadership",
			})

			err := cmd.Execute()
			assert.NoError(t, err)

			// Verify file was created
			assert.FileExists(t, outputFile)

			// Read and verify content
			content, err := os.ReadFile(outputFile)
			require.NoError(t, err)

			contentStr := string(content)
			assert.Contains(t, contentStr, "Leadership Project")
			assert.NotContains(t, contentStr, "Test Achievement")
		})

		// Test 4: Generate document filtered by tags
		t.Run("generate document filtered by tags", func(t *testing.T) {
			outputFile := filepath.Join(tempDir, "test-doc-tags.md")
			
			cmd := NewDocGenerateCmd(docService, bragService, tagService)
			cmd.SetArgs([]string{
				"--output", outputFile,
				"--tags", "test",
			})

			err := cmd.Execute()
			assert.NoError(t, err)

			// Verify file was created
			assert.FileExists(t, outputFile)

			// Read and verify content
			content, err := os.ReadFile(outputFile)
			require.NoError(t, err)

			contentStr := string(content)
			assert.Contains(t, contentStr, "Test Achievement")
		})

		// Test 5: Reject unsupported format
		t.Run("reject unsupported format", func(t *testing.T) {
			cmd := NewDocGenerateCmd(docService, bragService, tagService)
			cmd.SetArgs([]string{
				"--format", "pdf",
			})

			err := cmd.Execute()
			assert.Error(t, err)
			assert.Contains(t, strings.ToLower(err.Error()), "not yet supported")
		})

		// Test 6: Reject AI enhancement (not implemented)
		t.Run("reject AI enhancement", func(t *testing.T) {
			cmd := NewDocGenerateCmd(docService, bragService, tagService)
			cmd.SetArgs([]string{
				"--enhance-with-ai",
			})

			err := cmd.Execute()
			assert.Error(t, err)
			assert.Contains(t, strings.ToLower(err.Error()), "not yet implemented")
		})

		// Test 7: Error when no brags match criteria
		t.Run("error when no brags match criteria", func(t *testing.T) {
			cmd := NewDocGenerateCmd(docService, bragService, tagService)
			cmd.SetArgs([]string{
				"--category", "innovation",
			})

			err := cmd.Execute()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "no brags found")
		})

		// Test 8: Generate with invalid brag ID
		t.Run("error with invalid brag ID", func(t *testing.T) {
			cmd := NewDocGenerateCmd(docService, bragService, tagService)
			cmd.SetArgs([]string{
				"--brags", "999",
			})

			err := cmd.Execute()
			assert.Error(t, err)
		})

		// Test 9: Generate with invalid category
		t.Run("error with invalid category", func(t *testing.T) {
			cmd := NewDocGenerateCmd(docService, bragService, tagService)
			cmd.SetArgs([]string{
				"--category", "invalid-category",
			})

			err := cmd.Execute()
			assert.Error(t, err)
			assert.Contains(t, strings.ToLower(err.Error()), "invalid category")
		})

		// Test 10: Generate with multiple filters (category + tags)
		t.Run("generate with combined filters", func(t *testing.T) {
			outputFile := filepath.Join(tempDir, "test-doc-combined.md")
			
			cmd := NewDocGenerateCmd(docService, bragService, tagService)
			cmd.SetArgs([]string{
				"--output", outputFile,
				"--category", "achievement",
				"--tags", "test",
			})

			err := cmd.Execute()
			assert.NoError(t, err)

			// Verify file was created
			assert.FileExists(t, outputFile)

			// Read and verify content
			content, err := os.ReadFile(outputFile)
			require.NoError(t, err)

			contentStr := string(content)
			// Should only have the achievement with test tag
			assert.Contains(t, contentStr, "Test Achievement")
			assert.NotContains(t, contentStr, "Leadership Project")
		})

		// Test 11: Verify document structure and formatting
		t.Run("verify document structure", func(t *testing.T) {
			outputFile := filepath.Join(tempDir, "test-doc-structure.md")
			
			cmd := NewDocGenerateCmd(docService, bragService, tagService)
			cmd.SetArgs([]string{
				"--output", outputFile,
			})

			err := cmd.Execute()
			require.NoError(t, err)

			content, err := os.ReadFile(outputFile)
			require.NoError(t, err)

			contentStr := string(content)
			
			// Verify markdown structure
			assert.Contains(t, contentStr, "# Professional Achievements")
			assert.Contains(t, contentStr, "## Summary")
			assert.Contains(t, contentStr, "### Achievement")
			assert.Contains(t, contentStr, "### Leadership")
			assert.Contains(t, contentStr, "#### Test Achievement")
			assert.Contains(t, contentStr, "#### Leadership Project")
			assert.Contains(t, contentStr, "---")
			
			// Verify metadata format
			assert.Contains(t, contentStr, "**Test User**")
			assert.Contains(t, contentStr, "*Software Engineer*")
			assert.Contains(t, contentStr, "*Test Corp*")
			assert.Contains(t, contentStr, "Generated on")
			
			// Verify tags format
			assert.Contains(t, contentStr, "**Tags:**")
		})

		// Test 12: Generate to stdout (no output file)
		t.Run("generate to stdout", func(t *testing.T) {
			// This test verifies the command doesn't error when no output file is specified
			// In real execution, it would print to stdout
			cmd := NewDocGenerateCmd(docService, bragService, tagService)
			cmd.SetArgs([]string{
				"--brags", "1",
			})

			// Redirect stdout to capture output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := cmd.Execute()
			
			w.Close()
			os.Stdout = oldStdout

			assert.NoError(t, err)

			// Read captured output
			output, err := io.ReadAll(r)
			require.NoError(t, err)
			outputStr := string(output)

			// Verify output contains document content
			assert.Contains(t, outputStr, "Professional Achievements")
			assert.Contains(t, outputStr, "Test Achievement")
		})
	})
}
