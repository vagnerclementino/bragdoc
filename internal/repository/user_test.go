package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vagnerclementino/bragdoc/internal/domain"
)

func TestUserRepository_Insert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db.Conn())
	ctx := context.Background()

	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		JobTitle: "Software Engineer",
		Company:  "Tech Corp",
		Locale:   domain.LocaleEnglishUS,
	}

	created, err := repo.Insert(ctx, user)
	require.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, user.Name, created.Name)
	assert.Equal(t, user.Email, created.Email)
	assert.Equal(t, user.JobTitle, created.JobTitle)
	assert.Equal(t, user.Company, created.Company)
	assert.Equal(t, user.Locale, created.Locale)
	assert.NotZero(t, created.CreatedAt)
}

func TestUserRepository_Select(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db.Conn())
	ctx := context.Background()

	// Create a user
	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		JobTitle: "Software Engineer",
		Company:  "Tech Corp",
		Locale:   domain.LocaleEnglishUS,
	}
	created, err := repo.Insert(ctx, user)
	require.NoError(t, err)

	// Retrieve the user
	retrieved, err := repo.Select(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Name, retrieved.Name)
	assert.Equal(t, created.Email, retrieved.Email)
	assert.Equal(t, created.JobTitle, retrieved.JobTitle)
	assert.Equal(t, created.Company, retrieved.Company)
}

func TestUserRepository_Select_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db.Conn())
	ctx := context.Background()

	// Try to retrieve a non-existent user
	_, err := repo.Select(ctx, 99999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_SelectByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db.Conn())
	ctx := context.Background()

	// Create a user
	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		JobTitle: "Software Engineer",
		Company:  "Tech Corp",
		Locale:   domain.LocaleEnglishUS,
	}
	created, err := repo.Insert(ctx, user)
	require.NoError(t, err)

	// Retrieve by email
	retrieved, err := repo.SelectByEmail(ctx, "john@example.com")
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Email, retrieved.Email)
}

func TestUserRepository_SelectByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db.Conn())
	ctx := context.Background()

	// Try to retrieve a non-existent user
	_, err := repo.SelectByEmail(ctx, "nonexistent@example.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_SelectAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db.Conn())
	ctx := context.Background()

	// Create multiple users
	users := []*domain.User{
		{
			Name:   "John Doe",
			Email:  "john@example.com",
			Locale: domain.LocaleEnglishUS,
		},
		{
			Name:   "Jane Smith",
			Email:  "jane@example.com",
			Locale: domain.LocalePortugueseBR,
		},
		{
			Name:   "Bob Johnson",
			Email:  "bob@example.com",
			Locale: domain.LocaleEnglishUS,
		},
	}

	for _, user := range users {
		_, err := repo.Insert(ctx, user)
		require.NoError(t, err)
	}

	// Retrieve all users
	allUsers, err := repo.SelectAll(ctx)
	require.NoError(t, err)
	assert.Len(t, allUsers, 3)
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db.Conn())
	ctx := context.Background()

	// Create a user
	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		JobTitle: "Software Engineer",
		Company:  "Tech Corp",
		Locale:   domain.LocaleEnglishUS,
	}
	created, err := repo.Insert(ctx, user)
	require.NoError(t, err)

	// Update the user
	created.Name = "John Updated"
	created.JobTitle = "Senior Software Engineer"
	created.Company = "New Tech Corp"
	created.Locale = domain.LocalePortugueseBR

	updated, err := repo.Update(ctx, created)
	require.NoError(t, err)
	assert.Equal(t, "John Updated", updated.Name)
	assert.Equal(t, "Senior Software Engineer", updated.JobTitle)
	assert.Equal(t, "New Tech Corp", updated.Company)
	assert.Equal(t, domain.LocalePortugueseBR, updated.Locale)
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db.Conn())
	ctx := context.Background()

	// Create a user
	user := &domain.User{
		Name:   "John Doe",
		Email:  "john@example.com",
		Locale: domain.LocaleEnglishUS,
	}
	created, err := repo.Insert(ctx, user)
	require.NoError(t, err)

	// Delete the user
	err = repo.Delete(ctx, created.ID)
	require.NoError(t, err)

	// Verify it's deleted
	_, err = repo.Select(ctx, created.ID)
	assert.Error(t, err)
}

func TestUserRepository_UniqueEmailConstraint(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db.Conn())
	ctx := context.Background()

	// Create a user
	user := &domain.User{
		Name:   "John Doe",
		Email:  "john@example.com",
		Locale: domain.LocaleEnglishUS,
	}
	_, err := repo.Insert(ctx, user)
	require.NoError(t, err)

	// Try to create another user with the same email
	duplicateUser := &domain.User{
		Name:   "Jane Doe",
		Email:  "john@example.com", // Same email
		Locale: domain.LocaleEnglishUS,
	}
	_, err = repo.Insert(ctx, duplicateUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE constraint failed")
}

func TestUserRepository_InsertWithOptionalFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db.Conn())
	ctx := context.Background()

	// Create a user without optional fields
	user := &domain.User{
		Name:   "John Doe",
		Email:  "john@example.com",
		Locale: domain.LocaleEnglishUS,
		// JobTitle and Company are empty
	}

	created, err := repo.Insert(ctx, user)
	require.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, user.Name, created.Name)
	assert.Equal(t, user.Email, created.Email)
	assert.Empty(t, created.JobTitle)
	assert.Empty(t, created.Company)
}
