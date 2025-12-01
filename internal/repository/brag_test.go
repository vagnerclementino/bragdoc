package repository

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vagnerclementino/bragdoc/internal/database"
	"github.com/vagnerclementino/bragdoc/internal/domain"
)

func setupTestDB(t *testing.T) *database.DB {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := database.SetupDatabase(dbPath)
	require.NoError(t, err)
	return db
}

func createTestUser(t *testing.T, db *database.DB) int64 {
	ctx := context.Background()
	userRepo := NewUserRepository(db.Conn())
	user, err := userRepo.Insert(ctx, &domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Language: "en",
	})
	require.NoError(t, err)
	return user.ID
}

func TestBragRepository_Insert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewBragRepository(db.Conn())
	ctx := context.Background()

	brag := &domain.Brag{
		OwnerID:     userID,
		Title:       "Test Achievement",
		Description: "This is a test achievement with sufficient description length",
		Category:    domain.CategoryProject,
	}

	created, err := repo.Insert(ctx, brag)
	require.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, brag.Title, created.Title)
	assert.Equal(t, brag.Description, created.Description)
	assert.Equal(t, brag.Category, created.Category)
}

func TestBragRepository_Select(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewBragRepository(db.Conn())
	ctx := context.Background()

	// Create a brag
	brag := &domain.Brag{
		OwnerID:     userID,
		Title:       "Test Achievement",
		Description: "This is a test achievement with sufficient description length",
		Category:    domain.CategoryProject,
	}
	created, err := repo.Insert(ctx, brag)
	require.NoError(t, err)

	// Retrieve the brag
	retrieved, err := repo.Select(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Title, retrieved.Title)
	assert.Equal(t, created.Description, retrieved.Description)
}

func TestBragRepository_SelectAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewBragRepository(db.Conn())
	ctx := context.Background()

	// Create multiple brags
	for i := 0; i < 3; i++ {
		brag := &domain.Brag{
			OwnerID:     userID,
			Title:       "Test Achievement",
			Description: "This is a test achievement with sufficient description length",
			Category:    domain.CategoryProject,
		}
		_, err := repo.Insert(ctx, brag)
		require.NoError(t, err)
	}

	// Retrieve all brags
	brags, err := repo.SelectAll(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, brags, 3)
}

func TestBragRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewBragRepository(db.Conn())
	ctx := context.Background()

	// Create a brag
	brag := &domain.Brag{
		OwnerID:     userID,
		Title:       "Original Title",
		Description: "This is the original description with sufficient length",
		Category:    domain.CategoryProject,
	}
	created, err := repo.Insert(ctx, brag)
	require.NoError(t, err)

	// Update the brag
	created.Title = "Updated Title"
	created.Description = "This is the updated description with sufficient length"
	created.Category = domain.CategoryAchievement

	updated, err := repo.Update(ctx, created)
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.Title)
	assert.Equal(t, "This is the updated description with sufficient length", updated.Description)
	assert.Equal(t, domain.CategoryAchievement, updated.Category)
}

func TestBragRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewBragRepository(db.Conn())
	ctx := context.Background()

	// Create a brag
	brag := &domain.Brag{
		OwnerID:     userID,
		Title:       "Test Achievement",
		Description: "This is a test achievement with sufficient description length",
		Category:    domain.CategoryProject,
	}
	created, err := repo.Insert(ctx, brag)
	require.NoError(t, err)

	// Delete the brag
	err = repo.Delete(ctx, created.ID)
	require.NoError(t, err)

	// Verify it's deleted
	_, err = repo.Select(ctx, created.ID)
	assert.Error(t, err)
}

func TestBragRepository_SelectByCategory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewBragRepository(db.Conn())
	ctx := context.Background()

	// Create brags with different categories
	categories := []domain.Category{
		domain.CategoryProject,
		domain.CategoryAchievement,
		domain.CategoryProject,
	}

	for _, cat := range categories {
		brag := &domain.Brag{
			OwnerID:     userID,
			Title:       "Test Achievement",
			Description: "This is a test achievement with sufficient description length",
			Category:    cat,
		}
		_, err := repo.Insert(ctx, brag)
		require.NoError(t, err)
	}

	// Retrieve brags by category
	projectBrags, err := repo.SelectByCategory(ctx, userID, domain.CategoryProject)
	require.NoError(t, err)
	assert.Len(t, projectBrags, 2)

	achievementBrags, err := repo.SelectByCategory(ctx, userID, domain.CategoryAchievement)
	require.NoError(t, err)
	assert.Len(t, achievementBrags, 1)
}
