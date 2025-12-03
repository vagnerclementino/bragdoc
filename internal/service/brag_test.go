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

// MockBragRepository is a mock implementation of BragRepository
type MockBragRepository struct {
	mock.Mock
}

func (m *MockBragRepository) Insert(ctx context.Context, brag *domain.Brag) (*domain.Brag, error) {
	args := m.Called(ctx, brag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Brag), args.Error(1)
}

func (m *MockBragRepository) Select(ctx context.Context, id int64) (*domain.Brag, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Brag), args.Error(1)
}

func (m *MockBragRepository) SelectAll(ctx context.Context, userID int64) ([]*domain.Brag, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Brag), args.Error(1)
}

func (m *MockBragRepository) SelectByTags(ctx context.Context, userID int64, tagNames []string) ([]*domain.Brag, error) {
	args := m.Called(ctx, userID, tagNames)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Brag), args.Error(1)
}

func (m *MockBragRepository) SelectByCategory(ctx context.Context, userID int64, category domain.Category) ([]*domain.Brag, error) {
	args := m.Called(ctx, userID, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Brag), args.Error(1)
}

func (m *MockBragRepository) Update(ctx context.Context, brag *domain.Brag) (*domain.Brag, error) {
	args := m.Called(ctx, brag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Brag), args.Error(1)
}

func (m *MockBragRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// TestBragService_Create_Success tests successful brag creation
// **Feature: cli-architecture-refactor, Property 5: Services validate before persisting**
func TestBragService_Create_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	brag := &domain.Brag{
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Valid Title",
		Description: "This is a valid description with more than 20 characters",
		Category:    domain.CategoryAchievement,
		CreatedAt:   time.Now(),
	}

	expectedBrag := &domain.Brag{
		ID: 1,
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Valid Title",
		Description: "This is a valid description with more than 20 characters",
		Category:    domain.CategoryAchievement,
		CreatedAt:   time.Now(),
	}

	mockRepo.On("Insert", mock.Anything, brag).Return(expectedBrag, nil)

	// Act
	created, err := service.Create(context.Background(), brag)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, int64(1), created.ID)
	assert.Equal(t, "Valid Title", created.Title)
	mockRepo.AssertExpectations(t)
}

// TestBragService_Create_ValidationError_TitleTooShort tests validation for short titles
// **Feature: cli-architecture-refactor, Property 6: Validation errors are descriptive**
func TestBragService_Create_ValidationError_TitleTooShort(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	brag := &domain.Brag{
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Bad", // Too short (< 5 characters)
		Description: "This is a valid description with more than 20 characters",
		Category:    domain.CategoryAchievement,
	}

	// Act
	created, err := service.Create(context.Background(), brag)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "title must be at least 5 characters")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestBragService_Create_ValidationError_DescriptionTooShort tests validation for short descriptions
// **Feature: cli-architecture-refactor, Property 6: Validation errors are descriptive**
func TestBragService_Create_ValidationError_DescriptionTooShort(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	brag := &domain.Brag{
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Valid Title",
		Description: "Too short", // Too short (< 20 characters)
		Category:    domain.CategoryAchievement,
	}

	// Act
	created, err := service.Create(context.Background(), brag)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "description must be at least 20 characters")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestBragService_Create_ValidationError_InvalidCategory tests validation for invalid categories
// **Feature: cli-architecture-refactor, Property 7: Business rules are enforced**
func TestBragService_Create_ValidationError_InvalidCategory(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	brag := &domain.Brag{
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Valid Title",
		Description: "This is a valid description with more than 20 characters",
		Category:    domain.Category(999), // Invalid category
	}

	// Act
	created, err := service.Create(context.Background(), brag)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "invalid brag category")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestBragService_Create_ValidationError_NilBrag tests validation for nil brag
func TestBragService_Create_ValidationError_NilBrag(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	// Act
	created, err := service.Create(context.Background(), nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "brag cannot be nil")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestBragService_Create_ValidationError_InvalidOwnerID tests validation for invalid Owner.ID
// **Feature: brag-owner-refactor, Property 4: Validation rejects invalid Owner**
func TestBragService_Create_ValidationError_InvalidOwnerID(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	brag := &domain.Brag{
		Owner: domain.User{
			ID: 0, // Invalid Owner.ID
		},
		Title:       "Valid Title",
		Description: "This is a valid description with more than 20 characters",
		Category:    domain.CategoryAchievement,
	}

	// Act
	created, err := service.Create(context.Background(), brag)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "owner ID cannot be empty")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestBragService_Create_ValidationError_EmptyTitle tests validation for empty title
func TestBragService_Create_ValidationError_EmptyTitle(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	brag := &domain.Brag{
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "   ", // Empty after trim
		Description: "This is a valid description with more than 20 characters",
		Category:    domain.CategoryAchievement,
	}

	// Act
	created, err := service.Create(context.Background(), brag)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "title cannot be empty")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestBragService_Create_ValidationError_EmptyDescription tests validation for empty description
func TestBragService_Create_ValidationError_EmptyDescription(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	brag := &domain.Brag{
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Valid Title",
		Description: "   ", // Empty after trim
		Category:    domain.CategoryAchievement,
	}

	// Act
	created, err := service.Create(context.Background(), brag)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "description cannot be empty")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestBragService_Update_Success tests successful brag update
func TestBragService_Update_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	brag := &domain.Brag{
		ID: 1,
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Updated Title",
		Description: "This is an updated description with more than 20 characters",
		Category:    domain.CategoryProject,
	}

	expectedBrag := &domain.Brag{
		ID: 1,
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Updated Title",
		Description: "This is an updated description with more than 20 characters",
		Category:    domain.CategoryProject,
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("Update", mock.Anything, brag).Return(expectedBrag, nil)

	// Act
	updated, err := service.Update(context.Background(), brag)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, int64(1), updated.ID)
	assert.Equal(t, "Updated Title", updated.Title)
	mockRepo.AssertExpectations(t)
}

// TestBragService_Update_ValidationError tests validation during update
func TestBragService_Update_ValidationError(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	brag := &domain.Brag{
		ID: 1,
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Bad", // Too short
		Description: "This is a valid description with more than 20 characters",
		Category:    domain.CategoryProject,
	}

	// Act
	updated, err := service.Update(context.Background(), brag)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Contains(t, err.Error(), "title must be at least 5 characters")
	mockRepo.AssertNotCalled(t, "Update")
}

// TestBragService_GetByID_Success tests successful retrieval by ID
func TestBragService_GetByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	expectedBrag := &domain.Brag{
		ID: 1,
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Test Brag",
		Description: "This is a test brag description",
		Category:    domain.CategoryAchievement,
	}

	mockRepo.On("Select", mock.Anything, int64(1)).Return(expectedBrag, nil)

	// Act
	brag, err := service.GetByID(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, brag)
	assert.Equal(t, int64(1), brag.ID)
	assert.Equal(t, "Test Brag", brag.Title)
	mockRepo.AssertExpectations(t)
}

// TestBragService_GetByID_NotFound tests retrieval of non-existent brag
func TestBragService_GetByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	mockRepo.On("Select", mock.Anything, int64(999)).Return(nil, errors.New("brag not found"))

	// Act
	brag, err := service.GetByID(context.Background(), 999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, brag)
	mockRepo.AssertExpectations(t)
}

// TestBragService_List_Success tests successful listing of brags
func TestBragService_List_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	expectedBrags := []*domain.Brag{
		{
			ID: 1,
			Owner: domain.User{
				ID:    1,
				Name:  "Test User",
				Email: "test@example.com",
			},
			Title:       "First Brag",
			Description: "This is the first brag description",
			Category:    domain.CategoryAchievement,
		},
		{
			ID: 2,
			Owner: domain.User{
				ID:    1,
				Name:  "Test User",
				Email: "test@example.com",
			},
			Title:       "Second Brag",
			Description: "This is the second brag description",
			Category:    domain.CategoryProject,
		},
	}

	mockRepo.On("SelectAll", mock.Anything, int64(1)).Return(expectedBrags, nil)

	// Act
	brags, err := service.List(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, brags)
	assert.Len(t, brags, 2)
	assert.Equal(t, "First Brag", brags[0].Title)
	assert.Equal(t, "Second Brag", brags[1].Title)
	mockRepo.AssertExpectations(t)
}

// TestBragService_SearchByTags_Success tests successful search by tags
func TestBragService_SearchByTags_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	tagNames := []string{"golang", "backend"}
	expectedBrags := []*domain.Brag{
		{
			ID: 1,
			Owner: domain.User{
				ID:    1,
				Name:  "Test User",
				Email: "test@example.com",
			},
			Title:       "Backend Project",
			Description: "This is a backend project description",
			Category:    domain.CategoryProject,
		},
	}

	mockRepo.On("SelectByTags", mock.Anything, int64(1), tagNames).Return(expectedBrags, nil)

	// Act
	brags, err := service.SearchByTags(context.Background(), 1, tagNames)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, brags)
	assert.Len(t, brags, 1)
	assert.Equal(t, "Backend Project", brags[0].Title)
	mockRepo.AssertExpectations(t)
}

// TestBragService_SearchByTags_EmptyTags tests search with empty tag list
func TestBragService_SearchByTags_EmptyTags(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	// Act
	brags, err := service.SearchByTags(context.Background(), 1, []string{})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, brags)
	assert.Contains(t, err.Error(), "at least one tag name is required")
	mockRepo.AssertNotCalled(t, "SelectByTags")
}

// TestBragService_SearchByCategory_Success tests successful search by category
func TestBragService_SearchByCategory_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	expectedBrags := []*domain.Brag{
		{
			ID: 1,
			Owner: domain.User{
				ID:    1,
				Name:  "Test User",
				Email: "test@example.com",
			},
			Title:       "Achievement Brag",
			Description: "This is an achievement description",
			Category:    domain.CategoryAchievement,
		},
	}

	mockRepo.On("SelectByCategory", mock.Anything, int64(1), domain.CategoryAchievement).Return(expectedBrags, nil)

	// Act
	brags, err := service.SearchByCategory(context.Background(), 1, domain.CategoryAchievement)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, brags)
	assert.Len(t, brags, 1)
	assert.Equal(t, "Achievement Brag", brags[0].Title)
	assert.Equal(t, domain.CategoryAchievement, brags[0].Category)
	mockRepo.AssertExpectations(t)
}

// TestBragService_Delete_Success tests successful deletion
func TestBragService_Delete_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	// Act
	err := service.Delete(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBragService_Delete_Error tests deletion error
func TestBragService_Delete_Error(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(999)).Return(errors.New("brag not found"))

	// Act
	err := service.Delete(context.Background(), 999)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
