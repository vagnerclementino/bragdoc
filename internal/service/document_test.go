package service

import (
    "context"
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/vagnerclementino/bragdoc/internal/domain"
)

func TestDocumentService_GroupBragsByCategory(t *testing.T) {
    docService := &DocumentService{}

    achievementCategory, _ := domain.ParseCategory("ACHIEVEMENT")
    leadershipCategory, _ := domain.ParseCategory("LEADERSHIP")
    projectCategory, _ := domain.ParseCategory("PROJECT")

    tests := []struct {
        name     string
        brags    []*domain.Brag
        expected map[string]int // category name -> count
    }{
        {
            name: "Multiple categories",
            brags: []*domain.Brag{
                {ID: 1, Category: achievementCategory, Title: "Achievement 1"},
                {ID: 2, Category: leadershipCategory, Title: "Leadership 1"},
                {ID: 3, Category: achievementCategory, Title: "Achievement 2"},
                {ID: 4, Category: projectCategory, Title: "Project 1"},
            },
            expected: map[string]int{
                "ACHIEVEMENT": 2,
                "LEADERSHIP":  1,
                "PROJECT":     1,
            },
        },
        {
            name: "Single category",
            brags: []*domain.Brag{
                {ID: 1, Category: achievementCategory, Title: "Achievement 1"},
                {ID: 2, Category: achievementCategory, Title: "Achievement 2"},
            },
            expected: map[string]int{
                "ACHIEVEMENT": 2,
            },
        },
        {
            name:     "Empty brags",
            brags:    []*domain.Brag{},
            expected: map[string]int{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            grouped := docService.groupBragsByCategory(tt.brags)

            assert.Len(t, grouped, len(tt.expected))
            for categoryName, expectedCount := range tt.expected {
                assert.Len(t, grouped[categoryName], expectedCount, "Category %s should have %d brags", categoryName, expectedCount)
            }
        })
    }
}

func TestDocumentService_GetCategoryList(t *testing.T) {
    docService := &DocumentService{}

    achievementCategory, _ := domain.ParseCategory("ACHIEVEMENT")
    leadershipCategory, _ := domain.ParseCategory("LEADERSHIP")
    projectCategory, _ := domain.ParseCategory("PROJECT")

    tests := []struct {
        name     string
        brags    []*domain.Brag
        expected []string
    }{
        {
            name: "Multiple categories",
            brags: []*domain.Brag{
                {ID: 1, Category: achievementCategory, Title: "Achievement 1"},
                {ID: 2, Category: leadershipCategory, Title: "Leadership 1"},
                {ID: 3, Category: achievementCategory, Title: "Achievement 2"},
                {ID: 4, Category: projectCategory, Title: "Project 1"},
            },
            expected: []string{"ACHIEVEMENT", "LEADERSHIP", "PROJECT"},
        },
        {
            name: "Single category",
            brags: []*domain.Brag{
                {ID: 1, Category: achievementCategory, Title: "Achievement 1"},
                {ID: 2, Category: achievementCategory, Title: "Achievement 2"},
            },
            expected: []string{"ACHIEVEMENT"},
        },
        {
            name:     "Empty brags",
            brags:    []*domain.Brag{},
            expected: []string{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            grouped := docService.groupBragsByCategory(tt.brags)
            categories := docService.getCategoryList(grouped)
            assert.ElementsMatch(t, tt.expected, categories)
        })
    }
}

func TestDocumentService_Generate(t *testing.T) {
    mockUserService := &MockUserService{}
    docService := NewDocumentService(mockUserService)

    achievementCategory, _ := domain.ParseCategory("ACHIEVEMENT")
    leadershipCategory, _ := domain.ParseCategory("LEADERSHIP")

    user := &domain.User{
        ID:        1,
        Name:      "John Doe",
        Email:     "john@example.com",
        JobTitle:  "Senior Developer",
        Company:   "Tech Corp",
        Locale:    "en-US",
        CreatedAt: time.Now(),
    }

    brags := []*domain.Brag{
        {
            ID:          1,
            Owner:       *user,
            Title:       "Major Project Delivery",
            Description: "Successfully delivered a major project ahead of schedule with 99.9% uptime",
            Category:    achievementCategory,
            CreatedAt:   time.Now().Add(-24 * time.Hour),
        },
        {
            ID:          2,
            Owner:       *user,
            Title:       "Team Leadership",
            Description: "Led a team of 5 developers through a complex migration project",
            Category:    leadershipCategory,
            CreatedAt:   time.Now().Add(-48 * time.Hour),
        },
    }

    mockUserService.On("GetByID", mock.Anything, user.ID).Return(user, nil)

    doc, err := docService.Generate(context.Background(), brags, user.ID, GenerateOptions{
        Format: domain.FormatMarkdown,
    })

    assert.NoError(t, err)
    assert.NotNil(t, doc)
    assert.NotEmpty(t, doc.Content)
    assert.Contains(t, string(doc.Content), "John Doe")
    assert.Contains(t, string(doc.Content), "Tech Corp")
    assert.Contains(t, string(doc.Content), "Senior Developer")
    assert.Contains(t, string(doc.Content), "Major Project Delivery")
    assert.Contains(t, string(doc.Content), "Team Leadership")

    mockUserService.AssertExpectations(t)
}

func TestDocumentService_Generate_EmptyBrags(t *testing.T) {
    mockUserService := &MockUserService{}
    docService := NewDocumentService(mockUserService)

    doc, err := docService.Generate(context.Background(), []*domain.Brag{}, 1, GenerateOptions{
        Format: domain.FormatMarkdown,
    })

    assert.Error(t, err)
    assert.Nil(t, doc)
    assert.Contains(t, err.Error(), "no brags provided")
}

func TestDocumentService_Generate_UserNotFound(t *testing.T) {
    mockUserService := &MockUserService{}
    docService := NewDocumentService(mockUserService)

    achievementCategory, _ := domain.ParseCategory("ACHIEVEMENT")
    user := &domain.User{ID: 1, Name: "Test", Email: "test@example.com", Locale: "en-US"}
    brags := []*domain.Brag{
        {ID: 1, Owner: *user, Title: "Test", Description: "Test", Category: achievementCategory},
    }

    mockUserService.On("GetByID", mock.Anything, int64(999)).Return(nil, assert.AnError)

    doc, err := docService.Generate(context.Background(), brags, 999, GenerateOptions{
        Format: domain.FormatMarkdown,
    })

    assert.Error(t, err)
    assert.Nil(t, doc)

    mockUserService.AssertExpectations(t)
}

// MockUserService is a mock implementation of UserService for testing
type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetOrCreate(ctx context.Context, name, email string) (*domain.User, error) {
    args := m.Called(ctx, name, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
    args := m.Called(ctx, user)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}
