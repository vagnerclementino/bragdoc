package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vagnerclementino/bragdoc/internal/domain"
)

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	users map[int64]*domain.User
	nextID int64
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[int64]*domain.User),
		nextID: 1,
	}
}

func (m *MockUserRepository) Select(ctx context.Context, id int64) (*domain.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (m *MockUserRepository) SelectByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) SelectAll(ctx context.Context) ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) Insert(ctx context.Context, user *domain.User) (*domain.User, error) {
	user.ID = m.nextID
	user.CreatedAt = time.Now()
	m.users[user.ID] = user
	m.nextID++
	return user, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return user, nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	delete(m.users, id)
	return nil
}

func TestUserService_Create(t *testing.T) {
	tests := []struct {
		name    string
		user    *domain.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user with all fields",
			user: &domain.User{
				Name:     "John Doe",
				Email:    "john@example.com",
				JobTitle: "Developer",
				Company:  "Tech Corp",
				Locale:   domain.LocaleEnglishUS,
			},
			wantErr: false,
		},
		{
			name: "valid user with required fields only",
			user: &domain.User{
				Name:   "Jane Doe",
				Email:  "jane@example.com",
				Locale: domain.LocaleEnglishUS,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			user: &domain.User{
				Email:  "test@example.com",
				Locale: domain.LocaleEnglishUS,
			},
			wantErr: true,
			errMsg:  "validation failed",
		},
		{
			name: "missing email",
			user: &domain.User{
				Name:   "Test User",
				Locale: domain.LocaleEnglishUS,
			},
			wantErr: true,
			errMsg:  "validation failed",
		},
		{
			name: "missing locale defaults to en-US",
			user: &domain.User{
				Name:  "Test User",
				Email: "test@example.com",
			},
			wantErr: false, // Locale defaults to "en-US"
		},
		{
			name: "invalid email format",
			user: &domain.User{
				Name:   "Test User",
				Email:  "invalid-email",
				Locale: domain.LocaleEnglishUS,
			},
			wantErr: true,
			errMsg:  "validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUserRepository()
			service := NewUserService(repo)
			ctx := context.Background()

			created, err := service.Create(ctx, tt.user)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, created)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, created)
				assert.NotZero(t, created.ID)
				assert.Equal(t, tt.user.Name, created.Name)
				assert.Equal(t, tt.user.Email, created.Email)
				assert.Equal(t, tt.user.JobTitle, created.JobTitle)
				assert.Equal(t, tt.user.Company, created.Company)
				assert.Equal(t, tt.user.Locale, created.Locale)
				assert.False(t, created.CreatedAt.IsZero())
			}
		})
	}
}

func TestUserService_Create_DuplicateEmail(t *testing.T) {
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	ctx := context.Background()

	// Create first user
	user1 := &domain.User{
		Name:   "John Doe",
		Email:  "john@example.com",
		Locale: "en-US",
	}

	created1, err := service.Create(ctx, user1)
	require.NoError(t, err)
	assert.NotNil(t, created1)

	// Try to create second user with same email
	user2 := &domain.User{
		Name:   "Jane Doe",
		Email:  "john@example.com", // Same email
		Locale: "en-US",
	}

	created2, err := service.Create(ctx, user2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	assert.Nil(t, created2)

	// Verify only one user exists
	users, err := service.List(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 1)
}

func TestUserService_GetByID(t *testing.T) {
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	ctx := context.Background()

	// Create a user
	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Locale: "en-US",
	}

	created, err := service.Create(ctx, user)
	require.NoError(t, err)

	// Get user by ID
	retrieved, err := service.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Name, retrieved.Name)
	assert.Equal(t, created.Email, retrieved.Email)

	// Get non-existent user
	notFound, err := service.GetByID(ctx, 999)
	assert.NoError(t, err) // Mock returns nil, nil for not found
	assert.Nil(t, notFound)
}

func TestUserService_GetByEmail(t *testing.T) {
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	ctx := context.Background()

	// Create a user
	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Locale: "en-US",
	}

	created, err := service.Create(ctx, user)
	require.NoError(t, err)

	// Get user by email
	retrieved, err := service.GetByEmail(ctx, "john@example.com")
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Email, retrieved.Email)

	// Get non-existent user
	notFound, err := service.GetByEmail(ctx, "notfound@example.com")
	assert.NoError(t, err) // Mock returns nil, nil for not found
	assert.Nil(t, notFound)
}

func TestUserService_Update(t *testing.T) {
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	ctx := context.Background()

	// Create a user
	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Locale: "en-US",
	}

	created, err := service.Create(ctx, user)
	require.NoError(t, err)

	// Update user
	created.Name = "John Updated"
	created.JobTitle = "Senior Developer"
	created.Company = "New Corp"

	updated, err := service.Update(ctx, created)
	require.NoError(t, err)
	assert.Equal(t, "John Updated", updated.Name)
	assert.Equal(t, "Senior Developer", updated.JobTitle)
	assert.Equal(t, "New Corp", updated.Company)
	assert.False(t, updated.UpdatedAt.IsZero())
}

func TestUserService_Delete(t *testing.T) {
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	ctx := context.Background()

	// Create a user
	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Locale: "en-US",
	}

	created, err := service.Create(ctx, user)
	require.NoError(t, err)

	// Delete user
	err = service.Delete(ctx, created.ID)
	require.NoError(t, err)

	// Verify user is deleted
	retrieved, err := service.GetByID(ctx, created.ID)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

func TestUserService_List(t *testing.T) {
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	ctx := context.Background()

	// Create multiple users
	users := []*domain.User{
		{Name: "User 1", Email: "user1@example.com", Locale: "en-US"},
		{Name: "User 2", Email: "user2@example.com", Locale: "pt-BR"},
		{Name: "User 3", Email: "user3@example.com", Locale: "en-US"},
	}

	for _, user := range users {
		_, err := service.Create(ctx, user)
		require.NoError(t, err)
	}

	// List all users
	allUsers, err := service.List(ctx)
	require.NoError(t, err)
	assert.Len(t, allUsers, 3)
}
