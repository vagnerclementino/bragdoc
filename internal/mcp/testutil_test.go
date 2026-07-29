package mcp

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"github.com/vagnerclementino/bragdoc/internal/service"
)

// --- Mock Repositories ---

// MockBragRepository is a testify mock for repository.BragRepository.
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

// MockTagRepository is a testify mock for repository.TagRepository.
type MockTagRepository struct {
	mock.Mock
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

func (m *MockTagRepository) Insert(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	args := m.Called(ctx, tag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tag), args.Error(1)
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

// MockUserRepository is a testify mock for repository.UserRepository.
type MockUserRepository struct {
	mock.Mock
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

func (m *MockUserRepository) Insert(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
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

// MockJobTitleRepository is a testify mock for repository.JobTitleRepository.
type MockJobTitleRepository struct {
	mock.Mock
}

func (m *MockJobTitleRepository) Get(ctx context.Context, id int64) (*domain.JobTitle, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.JobTitle), args.Error(1)
}

func (m *MockJobTitleRepository) GetActive(ctx context.Context, userID int64) (*domain.JobTitle, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.JobTitle), args.Error(1)
}

func (m *MockJobTitleRepository) GetByName(ctx context.Context, userID int64, title string) (*domain.JobTitle, error) {
	args := m.Called(ctx, userID, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.JobTitle), args.Error(1)
}

func (m *MockJobTitleRepository) ListByUser(ctx context.Context, userID int64) ([]*domain.JobTitle, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.JobTitle), args.Error(1)
}

func (m *MockJobTitleRepository) Create(ctx context.Context, jobTitle *domain.JobTitle) (*domain.JobTitle, error) {
	args := m.Called(ctx, jobTitle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.JobTitle), args.Error(1)
}

func (m *MockJobTitleRepository) Update(ctx context.Context, jobTitle *domain.JobTitle) (*domain.JobTitle, error) {
	args := m.Called(ctx, jobTitle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.JobTitle), args.Error(1)
}

func (m *MockJobTitleRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- Test Fixtures ---

// testFixtures holds all mock repos and services for a test server.
type testFixtures struct {
	bragRepo     *MockBragRepository
	tagRepo      *MockTagRepository
	userRepo     *MockUserRepository
	jobTitleRepo *MockJobTitleRepository
	server       *Server
}

// newTestServer creates a test Server backed by real services with mock repositories.
// It creates the Server struct directly without calling NewServer (which would trigger
// MCP SDK tool registration and potentially panic on schema validation).
func newTestServer() *testFixtures {
	bragRepo := new(MockBragRepository)
	tagRepo := new(MockTagRepository)
	userRepo := new(MockUserRepository)
	jobTitleRepo := new(MockJobTitleRepository)

	bragService := service.NewBragService(bragRepo)
	tagService := service.NewTagService(tagRepo)
	userService := service.NewUserService(userRepo)
	docService := service.NewDocumentService(userService)
	jobService := service.NewJobTitleService(jobTitleRepo)

	srv := &Server{
		bragService: bragService,
		tagService:  tagService,
		userService: userService,
		docService:  docService,
		jobService:  jobService,
	}

	return &testFixtures{
		bragRepo:     bragRepo,
		tagRepo:      tagRepo,
		userRepo:     userRepo,
		jobTitleRepo: jobTitleRepo,
		server:       srv,
	}
}

// --- Valid categories ---

var validCategories = []string{"PROJECT", "ACHIEVEMENT", "SKILL", "LEADERSHIP", "INNOVATION"}

// fixedTime is a stable time for mock returns.
var fixedTime = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
