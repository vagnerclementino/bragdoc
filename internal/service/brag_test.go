package service

import (
	"context"
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

func (m *MockBragRepository) Insert(ctx context.Context, brag *domain.Brag) (*domain.Brag, error) {
	args := m.Called(ctx, brag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Brag), args.Error(1)
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

// TestBragService_Create tests brag creation
func TestBragService_Create(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	achievementCategory, _ := domain.ParseCategory("ACHIEVEMENT")

	brag := &domain.Brag{
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Test Brag",
		Description: "This is a test brag with more than 20 characters",
		Category:    achievementCategory,
	}

	expectedBrag := &domain.Brag{
		ID: 1,
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Test Brag",
		Description: "This is a test brag with more than 20 characters",
		Category:    achievementCategory,
		CreatedAt:   time.Now(),
	}

	mockRepo.On("Insert", mock.Anything, brag).Return(expectedBrag, nil)

	// Act
	created, err := service.Create(context.Background(), brag)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, int64(1), created.ID)
	assert.Equal(t, "Test Brag", created.Title)
	mockRepo.AssertExpectations(t)
}

// TestBragService_Create_ValidationError tests validation during creation
func TestBragService_Create_ValidationError(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	achievementCategory, _ := domain.ParseCategory("ACHIEVEMENT")

	brag := &domain.Brag{
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Bad", // Too short
		Description: "This is a valid description with more than 20 characters",
		Category:    achievementCategory,
	}

	// Act
	created, err := service.Create(context.Background(), brag)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "brag title must be at least 5 characters")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestBragService_GetByID tests brag retrieval by ID
func TestBragService_GetByID(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	achievementCategory, _ := domain.ParseCategory("ACHIEVEMENT")

	expectedBrag := &domain.Brag{
		ID: 1,
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Test Brag",
		Description: "This is a test brag with more than 20 characters",
		Category:    achievementCategory,
		CreatedAt:   time.Now(),
	}

	mockRepo.On("Select", mock.Anything, int64(1)).Return(expectedBrag, nil)

	// Act
	brag, err := service.GetByID(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, brag)
	assert.Equal(t, int64(1), brag.ID)
	mockRepo.AssertExpectations(t)
}

// TestBragService_GetByID_NotFound tests brag retrieval when not found
func TestBragService_GetByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	mockRepo.On("Select", mock.Anything, int64(999)).Return(nil, assert.AnError)

	// Act
	brag, err := service.GetByID(context.Background(), 999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, brag)
	mockRepo.AssertExpectations(t)
}

// TestBragService_List tests listing brags for a user
func TestBragService_List(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	achievementCategory, _ := domain.ParseCategory("ACHIEVEMENT")

	expectedBrags := []*domain.Brag{
		{
			ID: 1,
			Owner: domain.User{
				ID:    1,
				Name:  "Test User",
				Email: "test@example.com",
			},
			Title:       "Test Brag 1",
			Description: "This is test brag 1 with more than 20 characters",
			Category:    achievementCategory,
			CreatedAt:   time.Now(),
		},
		{
			ID: 2,
			Owner: domain.User{
				ID:    1,
				Name:  "Test User",
				Email: "test@example.com",
			},
			Title:       "Test Brag 2",
			Description: "This is test brag 2 with more than 20 characters",
			Category:    achievementCategory,
			CreatedAt:   time.Now(),
		},
	}

	mockRepo.On("SelectAll", mock.Anything, int64(1)).Return(expectedBrags, nil)

	// Act
	brags, err := service.List(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, brags)
	assert.Len(t, brags, 2)
	mockRepo.AssertExpectations(t)
}

// TestBragService_SearchByTags tests searching brags by tags
func TestBragService_SearchByTags(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	achievementCategory, _ := domain.ParseCategory("ACHIEVEMENT")

	expectedBrags := []*domain.Brag{
		{
			ID: 1,
			Owner: domain.User{
				ID:    1,
				Name:  "Test User",
				Email: "test@example.com",
			},
			Title:       "Test Brag with Tag",
			Description: "This is a test brag with more than 20 characters",
			Category:    achievementCategory,
			CreatedAt:   time.Now(),
		},
	}

	mockRepo.On("SelectByTags", mock.Anything, int64(1), []string{"go", "testing"}).Return(expectedBrags, nil)

	// Act
	brags, err := service.SearchByTags(context.Background(), 1, []string{"go", "testing"})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, brags)
	assert.Len(t, brags, 1)
	mockRepo.AssertExpectations(t)
}

// TestBragService_SearchByTags_EmptyTags tests searching with empty tags
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

// TestBragService_SearchByCategory tests searching brags by category
func TestBragService_SearchByCategory(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	achievementCategory, _ := domain.ParseCategory("ACHIEVEMENT")

	expectedBrags := []*domain.Brag{
		{
			ID: 1,
			Owner: domain.User{
				ID:    1,
				Name:  "Test User",
				Email: "test@example.com",
			},
			Title:       "Test Achievement",
			Description: "This is a test achievement with more than 20 characters",
			Category:    achievementCategory,
			CreatedAt:   time.Now(),
		},
	}

	mockRepo.On("SelectByCategory", mock.Anything, int64(1), achievementCategory).Return(expectedBrags, nil)

	// Act
	brags, err := service.SearchByCategory(context.Background(), 1, achievementCategory)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, brags)
	assert.Len(t, brags, 1)
	assert.Equal(t, achievementCategory.Name, brags[0].Category.Name)
	mockRepo.AssertExpectations(t)
}

// TestBragService_Update tests brag update
func TestBragService_Update(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	projectCategory, _ := domain.ParseCategory("PROJECT")

	brag := &domain.Brag{
		ID: 1,
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Updated Title",
		Description: "This is an updated description with more than 20 characters",
		Category:    projectCategory,
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
		Category:    projectCategory,
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

	projectCategory, _ := domain.ParseCategory("PROJECT")

	brag := &domain.Brag{
		ID: 1,
		Owner: domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		Title:       "Bad", // Too short
		Description: "This is a valid description with more than 20 characters",
		Category:    projectCategory,
	}

	// Act
	updated, err := service.Update(context.Background(), brag)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Contains(t, err.Error(), "brag title must be at least 5 characters")
	mockRepo.AssertNotCalled(t, "Update")
}

// TestBragService_Delete tests brag deletion
func TestBragService_Delete(t *testing.T) {
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

// TestBragService_Delete_Error tests brag deletion with error
func TestBragService_Delete_Error(t *testing.T) {
	// Arrange
	mockRepo := new(MockBragRepository)
	service := NewBragService(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(1)).Return(assert.AnError)

	// Act
	err := service.Delete(context.Background(), 1)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
