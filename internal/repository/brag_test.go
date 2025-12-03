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
		Name:   "Test User",
		Email:  "test@example.com",
		Locale: "en-US",
	})
	require.NoError(t, err)
	return user.ID
}

func createTestUserObject(t *testing.T, db *database.DB) *domain.User {
	ctx := context.Background()
	userRepo := NewUserRepository(db.Conn())
	user, err := userRepo.Insert(ctx, &domain.User{
		Name:   "Test User",
		Email:  "test@example.com",
		Locale: "en-US",
	})
	require.NoError(t, err)
	return user
}

func TestBragRepository_Insert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepository(db.Conn())
	user := createTestUserObject(t, db)
	repo := NewBragRepository(db.Conn(), userRepo)
	ctx := context.Background()

	brag := &domain.Brag{
		Owner:       *user,
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

	userRepo := NewUserRepository(db.Conn())
	user := createTestUserObject(t, db)
	repo := NewBragRepository(db.Conn(), userRepo)
	ctx := context.Background()

	// Create a brag
	brag := &domain.Brag{
		Owner:       *user,
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

func TestBragRepository_Select_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepository(db.Conn())
	repo := NewBragRepository(db.Conn(), userRepo)
	ctx := context.Background()

	// Try to retrieve a non-existent brag
	_, err := repo.Select(ctx, 99999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "brag not found")
}

func TestBragRepository_SelectAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepository(db.Conn())
	user := createTestUserObject(t, db)
	repo := NewBragRepository(db.Conn(), userRepo)
	ctx := context.Background()

	// Create multiple brags
	for i := 0; i < 3; i++ {
		brag := &domain.Brag{
			Owner:       *user,
			Title:       "Test Achievement",
			Description: "This is a test achievement with sufficient description length",
			Category:    domain.CategoryProject,
		}
		_, err := repo.Insert(ctx, brag)
		require.NoError(t, err)
	}

	// Retrieve all brags
	brags, err := repo.SelectAll(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, brags, 3)
}

func TestBragRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepository(db.Conn())
	user := createTestUserObject(t, db)
	repo := NewBragRepository(db.Conn(), userRepo)
	ctx := context.Background()

	// Create a brag
	brag := &domain.Brag{
		Owner:       *user,
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

	userRepo := NewUserRepository(db.Conn())
	user := createTestUserObject(t, db)
	repo := NewBragRepository(db.Conn(), userRepo)
	ctx := context.Background()

	// Create a brag
	brag := &domain.Brag{
		Owner:       *user,
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

	userRepo := NewUserRepository(db.Conn())
	user := createTestUserObject(t, db)
	repo := NewBragRepository(db.Conn(), userRepo)
	ctx := context.Background()

	// Create brags with different categories
	categories := []domain.Category{
		domain.CategoryProject,
		domain.CategoryAchievement,
		domain.CategoryProject,
	}

	for _, cat := range categories {
		brag := &domain.Brag{
			Owner:       *user,
			Title:       "Test Achievement",
			Description: "This is a test achievement with sufficient description length",
			Category:    cat,
		}
		_, err := repo.Insert(ctx, brag)
		require.NoError(t, err)
	}

	// Retrieve brags by category
	projectBrags, err := repo.SelectByCategory(ctx, user.ID, domain.CategoryProject)
	require.NoError(t, err)
	assert.Len(t, projectBrags, 2)

	achievementBrags, err := repo.SelectByCategory(ctx, user.ID, domain.CategoryAchievement)
	require.NoError(t, err)
	assert.Len(t, achievementBrags, 1)
}

func TestBragRepository_SelectByTags(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepository(db.Conn())
	user := createTestUserObject(t, db)
	bragRepo := NewBragRepository(db.Conn(), userRepo)
	tagRepo := NewTagRepository(db.Conn())
	ctx := context.Background()

	// Create tags
	tag1, err := tagRepo.Insert(ctx, &domain.Tag{
		Name:    "golang",
		OwnerID: user.ID,
	})
	require.NoError(t, err)

	tag2, err := tagRepo.Insert(ctx, &domain.Tag{
		Name:    "testing",
		OwnerID: user.ID,
	})
	require.NoError(t, err)

	// Create brags
	brag1 := &domain.Brag{
		Owner:       *user,
		Title:       "Brag with golang tag",
		Description: "This is a test achievement with sufficient description length",
		Category:    domain.CategoryProject,
	}
	created1, err := bragRepo.Insert(ctx, brag1)
	require.NoError(t, err)

	brag2 := &domain.Brag{
		Owner:       *user,
		Title:       "Brag with both tags",
		Description: "This is a test achievement with sufficient description length",
		Category:    domain.CategoryProject,
	}
	created2, err := bragRepo.Insert(ctx, brag2)
	require.NoError(t, err)

	brag3 := &domain.Brag{
		Owner:       *user,
		Title:       "Brag with no tags",
		Description: "This is a test achievement with sufficient description length",
		Category:    domain.CategoryProject,
	}
	_, err = bragRepo.Insert(ctx, brag3)
	require.NoError(t, err)

	// Attach tags to brags
	err = tagRepo.AttachToBrag(ctx, created1.ID, []int64{tag1.ID})
	require.NoError(t, err)

	err = tagRepo.AttachToBrag(ctx, created2.ID, []int64{tag1.ID, tag2.ID})
	require.NoError(t, err)

	// Search by single tag - should return brags with that tag
	brags, err := bragRepo.SelectByTags(ctx, user.ID, []string{"golang"})
	require.NoError(t, err)
	assert.Len(t, brags, 2) // Both brag1 and brag2 have golang tag

	// Search by multiple tags - should return brags with ANY of the tags (OR logic)
	brags, err = bragRepo.SelectByTags(ctx, user.ID, []string{"golang", "testing"})
	require.NoError(t, err)
	assert.Len(t, brags, 2) // Both brag1 and brag2 have at least one of these tags

	// Search by tag that only one brag has
	brags, err = bragRepo.SelectByTags(ctx, user.ID, []string{"testing"})
	require.NoError(t, err)
	assert.Len(t, brags, 1)
	assert.Equal(t, "Brag with both tags", brags[0].Title)
}
