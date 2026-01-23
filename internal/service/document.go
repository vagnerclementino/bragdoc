package service

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/vagnerclementino/bragdoc/internal/domain"
)

// DocumentService provides business logic for document generation
type DocumentService struct {
	userService *UserService
}

// NewDocumentService creates a new document service
func NewDocumentService(userService *UserService) *DocumentService {
	return &DocumentService{
		userService: userService,
	}
}

// GenerateOptions contains options for document generation
type GenerateOptions struct {
	Format        domain.DocumentFormat
	Template      string
	EnhanceWithAI bool
	Language      string
}

// Generate generates a document from brags
func (s *DocumentService) Generate(ctx context.Context, brags []*domain.Brag, userID int64, opts GenerateOptions) (*domain.Document, error) {
	if len(brags) == 0 {
		return nil, fmt.Errorf("no brags provided for document generation")
	}

	// Get user information
	user, err := s.userService.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user information: %w", err)
	}

	// For MVP, only Markdown format is supported
	if opts.Format != domain.FormatMarkdown {
		return nil, fmt.Errorf("format %s not yet supported (MVP supports only Markdown)", opts.Format.String())
	}

	// Generate document content based on format
	var content []byte
	switch opts.Format {
	case domain.FormatMarkdown:
		content, err = s.generateMarkdown(brags, user, opts.Template)
		if err != nil {
			return nil, fmt.Errorf("failed to generate markdown: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", opts.Format.String())
	}

	// Build metadata
	metadata := s.buildMetadata(brags, user)

	doc := &domain.Document{
		Format:   opts.Format,
		Template: opts.Template,
		Content:  content,
		Metadata: metadata,
	}

	return doc, nil
}

// generateMarkdown generates a markdown document from brags
func (s *DocumentService) generateMarkdown(brags []*domain.Brag, user *domain.User, templateName string) ([]byte, error) {
	// Get template
	tmpl, err := s.getTemplate(templateName)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// Group brags by category
	bragsByCategory := s.groupBragsByCategory(brags)

	// Prepare template data
	data := map[string]interface{}{
		"User":            user,
		"Brags":           brags,
		"BragsByCategory": bragsByCategory,
		"TotalBrags":      len(brags),
		"GeneratedAt":     time.Now().Format("January 2, 2006"),
		"Categories":      s.getCategoryList(bragsByCategory),
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

// getTemplate returns the template for the given name
func (s *DocumentService) getTemplate(name string) (*template.Template, error) {
	var tmplContent string

	switch name {
	case "default", "":
		tmplContent = defaultMarkdownTemplate
	case "executive":
		return nil, fmt.Errorf("executive template not yet implemented")
	case "technical":
		return nil, fmt.Errorf("technical template not yet implemented")
	default:
		return nil, fmt.Errorf("unknown template: %s", name)
	}

	// Create template with custom functions
	tmpl := template.New("document").Funcs(template.FuncMap{
		"title": strings.Title,
	})

	tmpl, err := tmpl.Parse(tmplContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl, nil
}

// groupBragsByCategory groups brags by their category
func (s *DocumentService) groupBragsByCategory(brags []*domain.Brag) map[string][]*domain.Brag {
	grouped := make(map[string][]*domain.Brag)

	for _, brag := range brags {
		category := brag.Category.String()
		grouped[category] = append(grouped[category], brag)
	}

	return grouped
}

// getCategoryList returns a sorted list of categories
func (s *DocumentService) getCategoryList(bragsByCategory map[string][]*domain.Brag) []string {
	categories := make([]string, 0, len(bragsByCategory))
	for category := range bragsByCategory {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}

// buildMetadata builds document metadata from brags and user
func (s *DocumentService) buildMetadata(brags []*domain.Brag, user *domain.User) domain.DocumentMetadata {
	// Extract unique categories
	categoryMap := make(map[string]bool)
	for _, brag := range brags {
		categoryMap[brag.Category.String()] = true
	}
	categories := make([]string, 0, len(categoryMap))
	for cat := range categoryMap {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	// Extract unique tags
	tagMap := make(map[string]bool)
	for _, brag := range brags {
		for _, tag := range brag.Tags {
			tagMap[tag.Name] = true
		}
	}
	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	return domain.DocumentMetadata{
		Title:       fmt.Sprintf("Professional Achievements - %s", user.Name),
		Author:      user.Name,
		BragCount:   len(brags),
		Categories:  categories,
		Tags:        tags,
		GeneratedAt: time.Now().Format("January 2, 2006"),
	}
}

// defaultMarkdownTemplate is the default template for markdown documents
const defaultMarkdownTemplate = `# Professional Achievements

**{{.User.Name}}**
{{if .User.JobTitle}}*{{.User.JobTitle}}*{{end}}
{{if .User.Company}}*{{.User.Company}}*{{end}}

Generated on {{.GeneratedAt}}

---

## Summary

This document contains {{.TotalBrags}} professional achievements across {{len .Categories}} categories.

{{range .Categories}}
### {{. | title}}

{{$brags := index $.BragsByCategory .}}
{{range $brags}}
#### {{.Title}}

{{.Description}}

{{if .Tags}}**Tags:** {{range $i, $tag := .Tags}}{{if $i}}, {{end}}{{$tag.Name}}{{end}}{{end}}

*Added on {{.CreatedAt.Format "January 2, 2006"}}*

---

{{end}}
{{end}}

## About This Document

This brag document was generated using Bragdoc CLI, a tool for tracking and documenting professional achievements.
`
