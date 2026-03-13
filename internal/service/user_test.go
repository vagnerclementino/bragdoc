package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vagnerclementino/bragdoc/internal/domain"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Insert(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Select(ctx context.Context, id int64) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) SelectByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) SelectAll(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// TestUserService_Create_Success tests successful user creation
// **Feature: cli-architecture-refactor, Property 5: Services validate before persisting**
func TestUserService_Create_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		JobTitle:  "Software Engineer",
		Company:   "Tech Corp",
		Locale:    domain.LocaleEnglishUS,
		CreatedAt: time.Now(),
	}

	expectedUser := &domain.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		JobTitle:  "Software Engineer",
		Company:   "Tech Corp",
		Locale:    domain.LocaleEnglishUS,
		CreatedAt: time.Now(),
	}

	mockRepo.On("SelectByEmail", mock.Anything, "john.doe@example.com").Return(nil, errors.New("not found"))
	mockRepo.On("Insert", mock.Anything, user).Return(expectedUser, nil)

	// Act
	created, err := service.Create(context.Background(), user)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, int64(1), created.ID)
	assert.Equal(t, "John Doe", created.Name)
	assert.Equal(t, "john.doe@example.com", created.Email)
	mockRepo.AssertExpectations(t)
}

// TestUserService_Create_ValidationError_NameTooShort tests validation for short names
// **Feature: cli-architecture-refactor, Property 6: Validation errors are descriptive**
func TestUserService_Create_ValidationError_NameTooShort(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{
		Name:   "J", // Too short (< 2 characters)
		Email:  "john@example.com",
		Locale: domain.LocaleEnglishUS,
	}

	// Act
	created, err := service.Create(context.Background(), user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "user name must be at least 2 characters")
	mockRepo.AssertNotCalled(t, "SelectByEmail")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestUserService_Create_ValidationError_InvalidEmail tests validation for invalid email
// **Feature: cli-architecture-refactor, Property 6: Validation errors are descriptive**
func TestUserService_Create_ValidationError_InvalidEmail(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{
		Name:   "John Doe",
		Email:  "invalid-email", // Invalid email format
		Locale: domain.LocaleEnglishUS,
	}

	// Act
	created, err := service.Create(context.Background(), user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "invalid email format")
	mockRepo.AssertNotCalled(t, "SelectByEmail")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestUserService_Create_ValidationError_InvalidLocale tests validation for invalid locale
// **Feature: cli-architecture-refactor, Property 7: Business rules are enforced**
func TestUserService_Create_ValidationError_InvalidLocale(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{
		Name:   "John Doe",
		Email:  "john@example.com",
		Locale: domain.Locale("fr-FR"), // Invalid locale (not supported)
	}

	// Act
	created, err := service.Create(context.Background(), user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "invalid locale")
	mockRepo.AssertNotCalled(t, "SelectByEmail")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestUserService_Create_ValidationError_NilUser tests validation for nil user
func TestUserService_Create_ValidationError_NilUser(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Act
	created, err := service.Create(context.Background(), nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "user cannot be nil")
	mockRepo.AssertNotCalled(t, "SelectByEmail")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestUserService_Create_ValidationError_EmptyName tests validation for empty name
func TestUserService_Create_ValidationError_EmptyName(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{
		Name:   "   ", // Empty after trim
		Email:  "john@example.com",
		Locale: domain.LocaleEnglishUS,
	}

	// Act
	created, err := service.Create(context.Background(), user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "user name cannot be empty")
	mockRepo.AssertNotCalled(t, "SelectByEmail")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestUserService_Create_ValidationError_EmptyEmail tests validation for empty email
func TestUserService_Create_ValidationError_EmptyEmail(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{
		Name:   "John Doe",
		Email:  "   ", // Empty after trim
		Locale: domain.LocaleEnglishUS,
	}

	// Act
	created, err := service.Create(context.Background(), user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "user email cannot be empty")
	mockRepo.AssertNotCalled(t, "SelectByEmail")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestUserService_Create_ValidationError_DuplicateEmail tests validation for duplicate email
// **Feature: cli-architecture-refactor, Property 7: Business rules are enforced**
func TestUserService_Create_ValidationError_DuplicateEmail(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{
		Name:   "John Doe",
		Email:  "john@example.com",
		Locale: domain.LocaleEnglishUS,
	}

	existingUser := &domain.User{
		ID:     1,
		Name:   "Jane Doe",
		Email:  "john@example.com",
		Locale: domain.LocaleEnglishUS,
	}

	mockRepo.On("SelectByEmail", mock.Anything, "john@example.com").Return(existingUser, nil)

	// Act
	created, err := service.Create(context.Background(), user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "user with email john@example.com already exists")
	mockRepo.AssertNotCalled(t, "Insert")
	mockRepo.AssertExpectations(t)
}

// TestUserService_Create_DefaultLocale tests that empty locale defaults to en-US
func TestUserService_Create_DefaultLocale(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{
		Name:   "John Doe",
		Email:  "john@example.com",
		Locale: "", // Empty locale should default to en-US
	}

	expectedUser := &domain.User{
		ID:     1,
		Name:   "John Doe",
		Email:  "john@example.com",
		Locale: domain.LocaleEnglishUS,
	}

	mockRepo.On("SelectByEmail", mock.Anything, "john@example.com").Return(nil, errors.New("not found"))
	mockRepo.On("Insert", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Locale == domain.LocaleEnglishUS
	})).Return(expectedUser, nil)

	// Act
	created, err := service.Create(context.Background(), user)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, domain.LocaleEnglishUS, created.Locale)
	mockRepo.AssertExpectations(t)
}

// TestUserService_Update_Success tests successful user update
// **Feature: cli-architecture-refactor, Property 5: Services validate before persisting**
func TestUserService_Update_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{
		ID:        1,
		Name:      "John Updated",
		Email:     "john.updated@example.com",
		JobTitle:  "Senior Engineer",
		Company:   "New Corp",
		Locale:    domain.LocalePortugueseBR,
		UpdatedAt: time.Now(),
	}

	expectedUser := &domain.User{
		ID:        1,
		Name:      "John Updated",
		Email:     "john.updated@example.com",
		JobTitle:  "Senior Engineer",
		Company:   "New Corp",
		Locale:    domain.LocalePortugueseBR,
		UpdatedAt: time.Now(),
	}

	mockRepo.On("Update", mock.Anything, user).Return(expectedUser, nil)

	// Act
	updated, err := service.Update(context.Background(), user)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, int64(1), updated.ID)
	assert.Equal(t, "John Updated", updated.Name)
	assert.Equal(t, domain.LocalePortugueseBR, updated.Locale)
	mockRepo.AssertExpectations(t)
}

// TestUserService_Update_ValidationError tests validation during update
// **Feature: cli-architecture-refactor, Property 6: Validation errors are descriptive**
func TestUserService_Update_ValidationError(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{
		ID:     1,
		Name:   "J", // Too short
		Email:  "john@example.com",
		Locale: domain.LocaleEnglishUS,
	}

	// Act
	updated, err := service.Update(context.Background(), user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Contains(t, err.Error(), "user name must be at least 2 characters")
	mockRepo.AssertNotCalled(t, "Update")
}

// TestUserService_Update_ValidationError_InvalidEmail tests email validation during update
func TestUserService_Update_ValidationError_InvalidEmail(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{
		ID:     1,
		Name:   "John Doe",
		Email:  "not-an-email", // Invalid email
		Locale: domain.LocaleEnglishUS,
	}

	// Act
	updated, err := service.Update(context.Background(), user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Contains(t, err.Error(), "invalid email format")
	mockRepo.AssertNotCalled(t, "Update")
}

// TestUserService_Update_ValidationError_InvalidLocale tests locale validation during update
func TestUserService_Update_ValidationError_InvalidLocale(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &domain.User{
		ID:     1,
		Name:   "John Doe",
		Email:  "john@example.com",
		Locale: domain.Locale("es-ES"), // Invalid locale
	}

	// Act
	updated, err := service.Update(context.Background(), user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Contains(t, err.Error(), "invalid locale")
	mockRepo.AssertNotCalled(t, "Update")
}

// TestUserService_GetByID_Success tests successful retrieval by ID
func TestUserService_GetByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	expectedUser := &domain.User{
		ID:     1,
		Name:   "John Doe",
		Email:  "john@example.com",
		Locale: domain.LocaleEnglishUS,
	}

	mockRepo.On("Select", mock.Anything, int64(1)).Return(expectedUser, nil)

	// Act
	user, err := service.GetByID(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "John Doe", user.Name)
	mockRepo.AssertExpectations(t)
}

// TestUserService_GetByID_NotFound tests retrieval of non-existent user
func TestUserService_GetByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("Select", mock.Anything, int64(999)).Return(nil, errors.New("user not found"))

	// Act
	user, err := service.GetByID(context.Background(), 999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// TestUserService_GetByEmail_Success tests successful retrieval by email
func TestUserService_GetByEmail_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	expectedUser := &domain.User{
		ID:     1,
		Name:   "John Doe",
		Email:  "john@example.com",
		Locale: domain.LocaleEnglishUS,
	}

	mockRepo.On("SelectByEmail", mock.Anything, "john@example.com").Return(expectedUser, nil)

	// Act
	user, err := service.GetByEmail(context.Background(), "john@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "john@example.com", user.Email)
	mockRepo.AssertExpectations(t)
}

// TestUserService_List_Success tests successful listing of users
func TestUserService_List_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	expectedUsers := []*domain.User{
		{
			ID:     1,
			Name:   "John Doe",
			Email:  "john@example.com",
			Locale: domain.LocaleEnglishUS,
		},
		{
			ID:     2,
			Name:   "Jane Smith",
			Email:  "jane@example.com",
			Locale: domain.LocalePortugueseBR,
		},
	}

	mockRepo.On("SelectAll", mock.Anything).Return(expectedUsers, nil)

	// Act
	users, err := service.List(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	assert.Equal(t, "John Doe", users[0].Name)
	assert.Equal(t, "Jane Smith", users[1].Name)
	mockRepo.AssertExpectations(t)
}

// TestUserService_Delete_Success tests successful deletion
func TestUserService_Delete_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	// Act
	err := service.Delete(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserService_Delete_Error tests deletion error
func TestUserService_Delete_Error(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(999)).Return(errors.New("user not found"))

	// Act
	err := service.Delete(context.Background(), 999)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
