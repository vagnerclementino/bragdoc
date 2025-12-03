package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vagnerclementino/bragdoc/internal/domain"
)

func TestTagRepository_Insert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewTagRepository(db.Conn())
	ctx := context.Background()

	tag := &domain.Tag{
		Name:    "golang",
		OwnerID: userID,
	}

	created, err := repo.Insert(ctx, tag)
	require.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, tag.Name, created.Name)
	assert.Equal(t, tag.OwnerID, created.OwnerID)
	assert.NotZero(t, created.CreatedAt)
}

func TestTagRepository_Select(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewTagRepository(db.Conn())
	ctx := context.Background()

	// Create a tag
	tag := &domain.Tag{
		Name:    "golang",
		OwnerID: userID,
	}
	created, err := repo.Insert(ctx, tag)
	require.NoError(t, err)

	// Retrieve the tag
	retrieved, err := repo.Select(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Name, retrieved.Name)
	assert.Equal(t, created.OwnerID, retrieved.OwnerID)
}

func TestTagRepository_Select_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewTagRepository(db.Conn())
	ctx := context.Background()

	// Try to retrieve a non-existent tag
	_, err := repo.Select(ctx, 99999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tag not found")
}

func TestTagRepository_SelectByName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewTagRepository(db.Conn())
	ctx := context.Background()

	// Create a tag
	tag := &domain.Tag{
		Name:    "golang",
		OwnerID: userID,
	}
	created, err := repo.Insert(ctx, tag)
	require.NoError(t, err)

	// Retrieve by name
	retrieved, err := repo.SelectByName(ctx, userID, "golang")
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Name, retrieved.Name)
}

func TestTagRepository_SelectByName_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewTagRepository(db.Conn())
	ctx := context.Background()

	// Try to retrieve a non-existent tag
	_, err := repo.SelectByName(ctx, userID, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tag not found")
}

func TestTagRepository_SelectAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewTagRepository(db.Conn())
	ctx := context.Background()

	// Create multiple tags
	tagNames := []string{"golang", "testing", "cli"}
	for _, name := range tagNames {
		tag := &domain.Tag{
			Name:    name,
			OwnerID: userID,
		}
		_, err := repo.Insert(ctx, tag)
		require.NoError(t, err)
	}

	// Retrieve all tags
	tags, err := repo.SelectAll(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, tags, 3)
}

func TestTagRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewTagRepository(db.Conn())
	ctx := context.Background()

	// Create a tag
	tag := &domain.Tag{
		Name:    "golang",
		OwnerID: userID,
	}
	created, err := repo.Insert(ctx, tag)
	require.NoError(t, err)

	// Delete the tag
	err = repo.Delete(ctx, created.ID)
	require.NoError(t, err)

	// Verify it's deleted
	_, err = repo.Select(ctx, created.ID)
	assert.Error(t, err)
}

func TestTagRepository_AttachToBrag(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := createTestUserObject(t, db)
	userRepo := NewUserRepository(db.Conn())
	tagRepo := NewTagRepository(db.Conn())
	bragRepo := NewBragRepository(db.Conn(), userRepo)
	ctx := context.Background()

	// Create a brag
	brag := &domain.Brag{
		Owner:       *user,
		Title:       "Test Achievement",
		Description: "This is a test achievement with sufficient description length",
		Category:    domain.CategoryProject,
	}
	createdBrag, err := bragRepo.Insert(ctx, brag)
	require.NoError(t, err)

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

	// Attach tags to brag
	err = tagRepo.AttachToBrag(ctx, createdBrag.ID, []int64{tag1.ID, tag2.ID})
	require.NoError(t, err)

	// Verify tags are attached
	tags, err := tagRepo.SelectByBrag(ctx, createdBrag.ID)
	require.NoError(t, err)
	assert.Len(t, tags, 2)
}

func TestTagRepository_DetachFromBrag(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := createTestUserObject(t, db)
	userRepo := NewUserRepository(db.Conn())
	tagRepo := NewTagRepository(db.Conn())
	bragRepo := NewBragRepository(db.Conn(), userRepo)
	ctx := context.Background()

	// Create a brag
	brag := &domain.Brag{
		Owner:       *user,
		Title:       "Test Achievement",
		Description: "This is a test achievement with sufficient description length",
		Category:    domain.CategoryProject,
	}
	createdBrag, err := bragRepo.Insert(ctx, brag)
	require.NoError(t, err)

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

	// Attach tags to brag
	err = tagRepo.AttachToBrag(ctx, createdBrag.ID, []int64{tag1.ID, tag2.ID})
	require.NoError(t, err)

	// Detach one tag
	err = tagRepo.DetachFromBrag(ctx, createdBrag.ID, []int64{tag1.ID})
	require.NoError(t, err)

	// Verify only one tag remains
	tags, err := tagRepo.SelectByBrag(ctx, createdBrag.ID)
	require.NoError(t, err)
	assert.Len(t, tags, 1)
	assert.Equal(t, "testing", tags[0].Name)
}

func TestTagRepository_UniqueConstraint(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewTagRepository(db.Conn())
	ctx := context.Background()

	// Create a tag
	tag := &domain.Tag{
		Name:    "golang",
		OwnerID: userID,
	}
	_, err := repo.Insert(ctx, tag)
	require.NoError(t, err)

	// Try to create a duplicate tag (same name + owner_id)
	duplicateTag := &domain.Tag{
		Name:    "golang",
		OwnerID: userID,
	}
	_, err = repo.Insert(ctx, duplicateTag)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE constraint failed")
}

func TestTagRepository_SelectByBrag(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := createTestUserObject(t, db)
	userRepo := NewUserRepository(db.Conn())
	tagRepo := NewTagRepository(db.Conn())
	bragRepo := NewBragRepository(db.Conn(), userRepo)
	ctx := context.Background()

	// Create a brag
	brag := &domain.Brag{
		Owner:       *user,
		Title:       "Test Achievement",
		Description: "This is a test achievement with sufficient description length",
		Category:    domain.CategoryProject,
	}
	createdBrag, err := bragRepo.Insert(ctx, brag)
	require.NoError(t, err)

	// Create tags
	tagNames := []string{"golang", "testing", "cli"}
	var tagIDs []int64
	for _, name := range tagNames {
		tag, err := tagRepo.Insert(ctx, &domain.Tag{
			Name:    name,
			OwnerID: user.ID,
		})
		require.NoError(t, err)
		tagIDs = append(tagIDs, tag.ID)
	}

	// Attach tags to brag
	err = tagRepo.AttachToBrag(ctx, createdBrag.ID, tagIDs)
	require.NoError(t, err)

	// Retrieve tags by brag
	tags, err := tagRepo.SelectByBrag(ctx, createdBrag.ID)
	require.NoError(t, err)
	assert.Len(t, tags, 3)
}
