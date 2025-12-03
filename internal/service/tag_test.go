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

// MockTagRepository is a mock implementation of TagRepository
type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) Insert(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	args := m.Called(ctx, tag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tag), args.Error(1)
}

func (m *MockTagRepository) Select(ctx context.Context, id int64) (*domain.Tag, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tag), args.Error(1)
}

func (m *MockTagRepository) SelectByName(ctx context.Context, ownerID int64, name string) (*domain.Tag, error) {
	args := m.Called(ctx, ownerID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tag), args.Error(1)
}

func (m *MockTagRepository) SelectAll(ctx context.Context, ownerID int64) ([]*domain.Tag, error) {
	args := m.Called(ctx, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Tag), args.Error(1)
}

func (m *MockTagRepository) SelectByBrag(ctx context.Context, bragID int64) ([]*domain.Tag, error) {
	args := m.Called(ctx, bragID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Tag), args.Error(1)
}

func (m *MockTagRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTagRepository) AttachToBrag(ctx context.Context, bragID int64, tagIDs []int64) error {
	args := m.Called(ctx, bragID, tagIDs)
	return args.Error(0)
}

func (m *MockTagRepository) DetachFromBrag(ctx context.Context, bragID int64, tagIDs []int64) error {
	args := m.Called(ctx, bragID, tagIDs)
	return args.Error(0)
}

// TestTagService_Create_Success tests successful tag creation
// **Feature: cli-architecture-refactor, Property 5: Services validate before persisting**
func TestTagService_Create_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	tag := &domain.Tag{
		Name:      "golang",
		OwnerID:   1,
		CreatedAt: time.Now(),
	}

	expectedTag := &domain.Tag{
		ID:        1,
		Name:      "golang",
		OwnerID:   1,
		CreatedAt: time.Now(),
	}

	mockRepo.On("SelectByName", mock.Anything, int64(1), "golang").Return(nil, errors.New("not found"))
	mockRepo.On("Insert", mock.Anything, tag).Return(expectedTag, nil)

	// Act
	created, err := service.Create(context.Background(), tag)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, int64(1), created.ID)
	assert.Equal(t, "golang", created.Name)
	mockRepo.AssertExpectations(t)
}

// TestTagService_Create_ValidationError_NameTooShort tests validation for short names
// **Feature: cli-architecture-refactor, Property 6: Validation errors are descriptive**
func TestTagService_Create_ValidationError_NameTooShort(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	tag := &domain.Tag{
		Name:    "a", // Too short (< 2 characters)
		OwnerID: 1,
	}

	// Act
	created, err := service.Create(context.Background(), tag)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "tag name must be at least 2 characters")
	mockRepo.AssertNotCalled(t, "SelectByName")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestTagService_Create_ValidationError_NameTooLong tests validation for long names
// **Feature: cli-architecture-refactor, Property 6: Validation errors are descriptive**
func TestTagService_Create_ValidationError_NameTooLong(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	tag := &domain.Tag{
		Name:    "this-is-a-very-long-tag-name-that-exceeds-limit", // Too long (> 20 characters)
		OwnerID: 1,
	}

	// Act
	created, err := service.Create(context.Background(), tag)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "tag name cannot exceed 20 characters")
	mockRepo.AssertNotCalled(t, "SelectByName")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestTagService_Create_ValidationError_DuplicateTag tests validation for duplicate tags
// **Feature: cli-architecture-refactor, Property 7: Business rules are enforced**
func TestTagService_Create_ValidationError_DuplicateTag(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	tag := &domain.Tag{
		Name:    "golang",
		OwnerID: 1,
	}

	existingTag := &domain.Tag{
		ID:      1,
		Name:    "golang",
		OwnerID: 1,
	}

	mockRepo.On("SelectByName", mock.Anything, int64(1), "golang").Return(existingTag, nil)

	// Act
	created, err := service.Create(context.Background(), tag)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "tag 'golang' already exists for this user")
	mockRepo.AssertNotCalled(t, "Insert")
	mockRepo.AssertExpectations(t)
}

// TestTagService_Create_ValidationError_NilTag tests validation for nil tag
func TestTagService_Create_ValidationError_NilTag(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	// Act
	created, err := service.Create(context.Background(), nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "tag cannot be nil")
	mockRepo.AssertNotCalled(t, "SelectByName")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestTagService_Create_ValidationError_EmptyName tests validation for empty name
func TestTagService_Create_ValidationError_EmptyName(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	tag := &domain.Tag{
		Name:    "   ", // Empty after trim
		OwnerID: 1,
	}

	// Act
	created, err := service.Create(context.Background(), tag)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Contains(t, err.Error(), "tag name cannot be empty")
	mockRepo.AssertNotCalled(t, "SelectByName")
	mockRepo.AssertNotCalled(t, "Insert")
}

// TestTagService_GetOrCreate_ExistingTag tests getting an existing tag
func TestTagService_GetOrCreate_ExistingTag(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	existingTag := &domain.Tag{
		ID:      1,
		Name:    "golang",
		OwnerID: 1,
	}

	mockRepo.On("SelectByName", mock.Anything, int64(1), "golang").Return(existingTag, nil)

	// Act
	tag, err := service.GetOrCreate(context.Background(), 1, "golang")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, int64(1), tag.ID)
	assert.Equal(t, "golang", tag.Name)
	mockRepo.AssertNotCalled(t, "Insert")
	mockRepo.AssertExpectations(t)
}

// TestTagService_GetOrCreate_NewTag tests creating a new tag
func TestTagService_GetOrCreate_NewTag(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	newTag := &domain.Tag{
		ID:      1,
		Name:    "python",
		OwnerID: 1,
	}

	mockRepo.On("SelectByName", mock.Anything, int64(1), "python").Return(nil, errors.New("not found"))
	mockRepo.On("Insert", mock.Anything, mock.MatchedBy(func(tag *domain.Tag) bool {
		return tag.Name == "python" && tag.OwnerID == 1
	})).Return(newTag, nil)

	// Act
	tag, err := service.GetOrCreate(context.Background(), 1, "python")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, int64(1), tag.ID)
	assert.Equal(t, "python", tag.Name)
	mockRepo.AssertExpectations(t)
}

// TestTagService_GetByID_Success tests successful retrieval by ID
func TestTagService_GetByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	expectedTag := &domain.Tag{
		ID:      1,
		Name:    "golang",
		OwnerID: 1,
	}

	mockRepo.On("Select", mock.Anything, int64(1)).Return(expectedTag, nil)

	// Act
	tag, err := service.GetByID(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, int64(1), tag.ID)
	assert.Equal(t, "golang", tag.Name)
	mockRepo.AssertExpectations(t)
}

// TestTagService_GetByName_Success tests successful retrieval by name
func TestTagService_GetByName_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	expectedTag := &domain.Tag{
		ID:      1,
		Name:    "golang",
		OwnerID: 1,
	}

	mockRepo.On("SelectByName", mock.Anything, int64(1), "golang").Return(expectedTag, nil)

	// Act
	tag, err := service.GetByName(context.Background(), 1, "golang")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, int64(1), tag.ID)
	assert.Equal(t, "golang", tag.Name)
	mockRepo.AssertExpectations(t)
}

// TestTagService_ListByUser_Success tests successful listing of user tags
func TestTagService_ListByUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	expectedTags := []*domain.Tag{
		{ID: 1, Name: "golang", OwnerID: 1},
		{ID: 2, Name: "python", OwnerID: 1},
	}

	mockRepo.On("SelectAll", mock.Anything, int64(1)).Return(expectedTags, nil)

	// Act
	tags, err := service.ListByUser(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tags)
	assert.Len(t, tags, 2)
	assert.Equal(t, "golang", tags[0].Name)
	assert.Equal(t, "python", tags[1].Name)
	mockRepo.AssertExpectations(t)
}

// TestTagService_ListByBrag_Success tests successful listing of brag tags
func TestTagService_ListByBrag_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	expectedTags := []*domain.Tag{
		{ID: 1, Name: "golang", OwnerID: 1},
		{ID: 2, Name: "backend", OwnerID: 1},
	}

	mockRepo.On("SelectByBrag", mock.Anything, int64(1)).Return(expectedTags, nil)

	// Act
	tags, err := service.ListByBrag(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tags)
	assert.Len(t, tags, 2)
	mockRepo.AssertExpectations(t)
}

// TestTagService_AttachToBrag_Success tests successful tag attachment
func TestTagService_AttachToBrag_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	tagIDs := []int64{1, 2}
	mockRepo.On("AttachToBrag", mock.Anything, int64(1), tagIDs).Return(nil)

	// Act
	err := service.AttachToBrag(context.Background(), 1, tagIDs)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestTagService_AttachToBrag_EmptyTags tests attachment with empty tag list
func TestTagService_AttachToBrag_EmptyTags(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	// Act
	err := service.AttachToBrag(context.Background(), 1, []int64{})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one tag ID is required")
	mockRepo.AssertNotCalled(t, "AttachToBrag")
}

// TestTagService_DetachFromBrag_Success tests successful tag detachment
func TestTagService_DetachFromBrag_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	tagIDs := []int64{1, 2}
	mockRepo.On("DetachFromBrag", mock.Anything, int64(1), tagIDs).Return(nil)

	// Act
	err := service.DetachFromBrag(context.Background(), 1, tagIDs)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestTagService_DetachFromBrag_EmptyTags tests detachment with empty tag list
func TestTagService_DetachFromBrag_EmptyTags(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	// Act
	err := service.DetachFromBrag(context.Background(), 1, []int64{})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one tag ID is required")
	mockRepo.AssertNotCalled(t, "DetachFromBrag")
}

// TestTagService_Delete_Success tests successful deletion
func TestTagService_Delete_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	// Act
	err := service.Delete(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestTagService_Delete_Error tests deletion error
func TestTagService_Delete_Error(t *testing.T) {
	// Arrange
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(999)).Return(errors.New("tag not found"))

	// Act
	err := service.Delete(context.Background(), 999)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
