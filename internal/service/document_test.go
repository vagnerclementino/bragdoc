package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vagnerclementino/bragdoc/internal/domain"
)

func TestDocumentService_Generate_EmptyBrags(t *testing.T) {
	docService := &DocumentService{}

	opts := GenerateOptions{
		Format:   domain.FormatMarkdown,
		Template: "default",
	}

	_, err := docService.Generate(context.Background(), []*domain.Brag{}, 1, opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no brags provided")
}

// TestDocumentService_Generate_UnsupportedFormat is tested in E2E tests
// because it requires a full setup with userService

func TestDocumentService_GroupBragsByCategory(t *testing.T) {
	docService := &DocumentService{}

	tests := []struct {
		name     string
		brags    []*domain.Brag
		expected map[string]int // category -> count
	}{
		{
			name: "Multiple categories",
			brags: []*domain.Brag{
				{ID: 1, Category: domain.CategoryAchievement, Title: "Achievement 1"},
				{ID: 2, Category: domain.CategoryLeadership, Title: "Leadership 1"},
				{ID: 3, Category: domain.CategoryAchievement, Title: "Achievement 2"},
				{ID: 4, Category: domain.CategoryProject, Title: "Project 1"},
			},
			expected: map[string]int{
				"achievement": 2,
				"leadership":  1,
				"project":     1,
			},
		},
		{
			name: "Single category",
			brags: []*domain.Brag{
				{ID: 1, Category: domain.CategoryAchievement, Title: "Achievement 1"},
				{ID: 2, Category: domain.CategoryAchievement, Title: "Achievement 2"},
			},
			expected: map[string]int{
				"achievement": 2,
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
			for category, expectedCount := range tt.expected {
				assert.Len(t, grouped[category], expectedCount, "Category %s should have %d brags", category, expectedCount)
			}
		})
	}
}

func TestDocumentService_GetCategoryList(t *testing.T) {
	docService := &DocumentService{}

	tests := []struct {
		name            string
		bragsByCategory map[string][]*domain.Brag
		expected        []string
	}{
		{
			name: "Multiple categories sorted",
			bragsByCategory: map[string][]*domain.Brag{
				"leadership":  {{ID: 2}},
				"achievement": {{ID: 1}},
				"project":     {{ID: 3}},
			},
			expected: []string{"achievement", "leadership", "project"},
		},
		{
			name: "Single category",
			bragsByCategory: map[string][]*domain.Brag{
				"achievement": {{ID: 1}},
			},
			expected: []string{"achievement"},
		},
		{
			name:            "Empty map",
			bragsByCategory: map[string][]*domain.Brag{},
			expected:        []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categories := docService.getCategoryList(tt.bragsByCategory)

			assert.Equal(t, tt.expected, categories, "Categories should be sorted alphabetically")
		})
	}
}

func TestDocumentService_BuildMetadata(t *testing.T) {
	docService := &DocumentService{}

	user := &domain.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	}

	tests := []struct {
		name               string
		brags              []*domain.Brag
		expectedBragCount  int
		expectedCategories []string
		expectedTags       []string
	}{
		{
			name: "Multiple brags with tags",
			brags: []*domain.Brag{
				{
					ID:       1,
					Category: domain.CategoryAchievement,
					Tags: []*domain.Tag{
						{Name: "test"},
						{Name: "automation"},
					},
				},
				{
					ID:       2,
					Category: domain.CategoryLeadership,
					Tags: []*domain.Tag{
						{Name: "leadership"},
						{Name: "test"}, // Duplicate should be deduplicated
					},
				},
			},
			expectedBragCount:  2,
			expectedCategories: []string{"achievement", "leadership"},
			expectedTags:       []string{"automation", "leadership", "test"}, // Sorted
		},
		{
			name: "Brags without tags",
			brags: []*domain.Brag{
				{ID: 1, Category: domain.CategoryProject},
			},
			expectedBragCount:  1,
			expectedCategories: []string{"project"},
			expectedTags:       []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := docService.buildMetadata(tt.brags, user)

			assert.Equal(t, tt.expectedBragCount, metadata.BragCount)
			assert.Equal(t, tt.expectedCategories, metadata.Categories)
			assert.Equal(t, tt.expectedTags, metadata.Tags)
			assert.Equal(t, "Test User", metadata.Author)
			assert.Contains(t, metadata.Title, "Test User")
			assert.NotEmpty(t, metadata.GeneratedAt)
		})
	}
}

func TestDocumentService_GetTemplate(t *testing.T) {
	docService := &DocumentService{}

	tests := []struct {
		name         string
		templateName string
		shouldError  bool
		errorMsg     string
	}{
		{
			name:         "Default template",
			templateName: "default",
			shouldError:  false,
		},
		{
			name:         "Empty template name (defaults to default)",
			templateName: "",
			shouldError:  false,
		},
		{
			name:         "Executive template (not implemented)",
			templateName: "executive",
			shouldError:  true,
			errorMsg:     "not yet implemented",
		},
		{
			name:         "Technical template (not implemented)",
			templateName: "technical",
			shouldError:  true,
			errorMsg:     "not yet implemented",
		},
		{
			name:         "Unknown template",
			templateName: "unknown",
			shouldError:  true,
			errorMsg:     "unknown template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := docService.getTemplate(tt.templateName)

			if tt.shouldError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, tmpl)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, tmpl)
			}
		})
	}
}

func TestDocumentService_GenerateMarkdown_TemplateContent(t *testing.T) {
	docService := &DocumentService{}

	user := &domain.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		JobTitle: "Senior Engineer",
		Company:  "Tech Corp",
	}

	brags := []*domain.Brag{
		{
			ID:          1,
			Title:       "Implemented Feature X",
			Description: "Successfully implemented a critical feature that improved performance by 50%",
			Category:    domain.CategoryAchievement,
			Tags: []*domain.Tag{
				{Name: "performance"},
				{Name: "backend"},
			},
			CreatedAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          2,
			Title:       "Led Team Migration",
			Description: "Led a team of 5 engineers to migrate legacy system to microservices",
			Category:    domain.CategoryLeadership,
			Tags: []*domain.Tag{
				{Name: "leadership"},
				{Name: "migration"},
			},
			CreatedAt: time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	content, err := docService.generateMarkdown(brags, user, "default")
	require.NoError(t, err)
	require.NotEmpty(t, content)

	contentStr := string(content)

	// Verify header information
	assert.Contains(t, contentStr, "Professional Achievements")
	assert.Contains(t, contentStr, "John Doe")
	assert.Contains(t, contentStr, "Senior Engineer")
	assert.Contains(t, contentStr, "Tech Corp")

	// Verify summary
	assert.Contains(t, contentStr, "2 professional achievements")

	// Verify brags are included
	assert.Contains(t, contentStr, "Implemented Feature X")
	assert.Contains(t, contentStr, "Successfully implemented a critical feature")
	assert.Contains(t, contentStr, "Led Team Migration")
	assert.Contains(t, contentStr, "Led a team of 5 engineers")

	// Verify tags are included
	assert.Contains(t, contentStr, "performance")
	assert.Contains(t, contentStr, "backend")
	assert.Contains(t, contentStr, "leadership")
	assert.Contains(t, contentStr, "migration")

	// Verify categories
	assert.Contains(t, contentStr, "Achievement")
	assert.Contains(t, contentStr, "Leadership")

	// Verify footer
	assert.Contains(t, contentStr, "Bragdoc CLI")
}

func TestDocumentService_GenerateMarkdown_WithoutOptionalFields(t *testing.T) {
	docService := &DocumentService{}

	user := &domain.User{
		ID:    1,
		Name:  "Jane Smith",
		Email: "jane@example.com",
		// No JobTitle or Company
	}

	brags := []*domain.Brag{
		{
			ID:          1,
			Title:       "Simple Achievement",
			Description: "A simple achievement without tags",
			Category:    domain.CategoryAchievement,
			Tags:        []*domain.Tag{}, // No tags
			CreatedAt:   time.Now(),
		},
	}

	content, err := docService.generateMarkdown(brags, user, "default")
	require.NoError(t, err)
	require.NotEmpty(t, content)

	contentStr := string(content)

	// Should still have user name
	assert.Contains(t, contentStr, "Jane Smith")

	// Should have the brag
	assert.Contains(t, contentStr, "Simple Achievement")

	// Should not have empty job title or company lines
	lines := strings.Split(contentStr, "\n")
	for _, line := range lines {
		// Make sure we don't have lines with just asterisks (empty optional fields)
		assert.NotEqual(t, "**", strings.TrimSpace(line))
	}
}
