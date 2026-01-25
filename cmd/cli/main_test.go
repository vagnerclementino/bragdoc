package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vagnerclementino/bragdoc/config"
)

var (
	update     = flag.Bool("update", false, "update golden files")
	binaryName = "bragdoc-test"
	binaryPath = ""
)

func TestMain(m *testing.M) {
	// Change to project root
	if err := os.Chdir("../.."); err != nil {
		fmt.Printf("could not change dir: %v\n", err)
		os.Exit(1)
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("could not get current dir: %v\n", err)
		os.Exit(1)
	}

	binaryPath = filepath.Join(dir, binaryName)

	// Build test binary with coverage
	version := "0.1.0"
	build := "test"
	ldflags := fmt.Sprintf("-X 'github.com/vagnerclementino/bragdoc/internal/command.Version=%s' -X 'github.com/vagnerclementino/bragdoc/internal/command.Build=%s'", version, build)

	// #nosec G204 - Command arguments are controlled by test code
	buildCmd := exec.Command("go", "build", "-cover", "-o", binaryName, "-ldflags", ldflags, "./cmd/cli")
	buildCmd.Env = append(os.Environ(), "CGO_ENABLED=1")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		fmt.Printf("Failed to build test binary: %v\n%s\n", err, output)
		os.Exit(1)
	}

	// Create coverage directory
	if err := os.MkdirAll(".coverdata", 0750); err != nil {
		fmt.Printf("could not create coverage dir: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	if err := os.Remove(binaryPath); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Failed to remove binary: %v\n", err)
	}
	if err := os.RemoveAll(".coverdata"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to remove coverdata: %v\n", err)
	}

	os.Exit(code)
}

// Helper functions for golden files
func runBinary(args []string, env map[string]string) ([]byte, error) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, binaryPath, args...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR=.coverdata")

	// Add custom environment variables
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	return cmd.CombinedOutput()
}

func loadFixture(t *testing.T, filename string) string {
	t.Helper()

	// #nosec G304 - Test fixture path is controlled by test code
	content, err := os.ReadFile(filepath.Join("testdata", "golden", filename))
	if err != nil {
		t.Fatalf("could not read fixture %s: %v", filename, err)
	}

	return string(content)
}

func writeFixture(t *testing.T, filename string, content []byte) {
	t.Helper()

	fixtureDir := filepath.Join("testdata", "golden")
	if err := os.MkdirAll(fixtureDir, 0750); err != nil {
		t.Fatalf("could not create fixture dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(fixtureDir, filename), content, 0600); err != nil {
		t.Fatalf("could not write fixture %s: %v", filename, err)
	}
}

// Unit tests
func TestExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "expand tilde only",
			input:    "~",
			expected: homeDir,
		},
		{
			name:     "expand tilde with path",
			input:    "~/.bragdoc/bragdoc.db",
			expected: filepath.Join(homeDir, ".bragdoc", "bragdoc.db"),
		},
		{
			name:     "no tilde - absolute path",
			input:    "/var/data/bragdoc.db",
			expected: "/var/data/bragdoc.db",
		},
		{
			name:     "no tilde - relative path",
			input:    "./data/bragdoc.db",
			expected: "./data/bragdoc.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandPath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDatabasePath_UsesConfigPath(t *testing.T) {
	tests := []struct {
		name           string
		configPath     string
		expectedResult string
	}{
		{
			name:           "uses config path with tilde",
			configPath:     "~/custom/path/db.sqlite",
			expectedResult: expandPath("~/custom/path/db.sqlite"),
		},
		{
			name:           "uses config path absolute",
			configPath:     "/var/data/bragdoc.db",
			expectedResult: "/var/data/bragdoc.db",
		},
		{
			name:           "uses config path relative",
			configPath:     "./data/bragdoc.db",
			expectedResult: "./data/bragdoc.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Database: config.DatabaseConfig{
					Path: tt.configPath,
				},
			}

			result := getDatabasePath(cfg)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGetDatabasePath_UsesDefaultWhenEmpty(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Path: "",
		},
	}

	result := getDatabasePath(cfg)

	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)
	expected := filepath.Join(homeDir, ".bragdoc", "bragdoc.db")

	assert.Equal(t, expected, result)
}

// Integration tests
func TestCLIVersion(t *testing.T) {
	output, err := runBinary([]string{"version"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if *update {
		writeFixture(t, "version.golden", output)
	}

	actual := string(output)
	expected := loadFixture(t, "version.golden")

	if actual != expected {
		t.Fatalf("actual = %s, expected = %s", actual, expected)
	}
}

func TestCLIHelp(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		fixture string
	}{
		{"root help", []string{"--help"}, "help.golden"},
		{"brag help", []string{"brag", "--help"}, "brag-help.golden"},
		{"tag help", []string{"tag", "--help"}, "tag-help.golden"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runBinary(tt.args, nil)
			if err != nil {
				t.Fatal(err)
			}

			if *update {
				writeFixture(t, tt.fixture, output)
			}

			actual := string(output)
			expected := loadFixture(t, tt.fixture)

			if actual != expected {
				t.Fatalf("actual = %s, expected = %s", actual, expected)
			}
		})
	}
}

func TestCLIRequiresInit(t *testing.T) {
	// Create temporary HOME directory
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		args    []string
		fixture string
	}{
		{"brag list without init", []string{"brag", "list"}, "brag-list-no-init.golden"},
		{"tag list without init", []string{"tag", "list"}, "tag-list-no-init.golden"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, _ := runBinary(tt.args, map[string]string{"HOME": tmpDir})

			if *update {
				writeFixture(t, tt.fixture, output)
			}

			actual := string(output)
			expected := loadFixture(t, tt.fixture)

			if actual != expected {
				t.Fatalf("actual = %s, expected = %s", actual, expected)
			}
		})
	}
}

func TestCLIFullWorkflow(t *testing.T) {
	tmpDir := t.TempDir()
	env := map[string]string{"HOME": tmpDir}

	// Test version
	t.Run("version", func(t *testing.T) {
		output, err := runBinary([]string{"version"}, env)
		assert.NoError(t, err)
		assert.Contains(t, string(output), "Bragdoc 0.1.0")
	})

	// Test help
	t.Run("help", func(t *testing.T) {
		output, err := runBinary([]string{"--help"}, env)
		assert.NoError(t, err)
		assert.Contains(t, string(output), "Bragdoc is a powerful command-line interface")
	})

	// Initialize bragdoc
	t.Run("init", func(t *testing.T) {
		output, err := runBinary([]string{"init", "--name", "TestUser", "--email", "test@example.com"}, env)
		assert.NoError(t, err)
		assert.Contains(t, string(output), "Bragdoc initialized successfully")
	})

	// CRUD: Create brags
	t.Run("create brags", func(t *testing.T) {
		tests := []struct {
			title       string
			description string
			category    string
			tags        string
		}{
			{"Achievement One", "This is my first achievement with enough characters", "achievement", "test,smoke"},
			{"Leadership Project", "Led a team to deliver important project successfully", "leadership", "leadership,team"},
			{"Innovation Work", "Created innovative solution for complex problem here", "innovation", ""},
		}

		for _, tt := range tests {
			args := []string{"brag", "add", "--title", tt.title, "--description", tt.description, "--category", tt.category}
			if tt.tags != "" {
				args = append(args, "--tags", tt.tags)
			}
			output, err := runBinary(args, env)
			assert.NoError(t, err)
			assert.Contains(t, string(output), "Brag created successfully")
		}
	})

	// CRUD: List brags
	t.Run("list brags", func(t *testing.T) {
		output, err := runBinary([]string{"brag", "list"}, env)
		assert.NoError(t, err)
		assert.Contains(t, string(output), "Achievement One")
		assert.Contains(t, string(output), "Leadership Project")
		assert.Contains(t, string(output), "Innovation Work")
	})

	// CRUD: Show brag
	t.Run("show brag", func(t *testing.T) {
		output, err := runBinary([]string{"brag", "show", "1"}, env)
		assert.NoError(t, err)
		assert.Contains(t, string(output), "Achievement One")
		assert.Contains(t, string(output), "This is my first achievement")
	})

	// CRUD: Edit brag
	t.Run("edit brag", func(t *testing.T) {
		output, err := runBinary([]string{"brag", "edit", "1", "--title", "Updated Achievement"}, env)
		assert.NoError(t, err)
		assert.Contains(t, string(output), "Brag updated successfully")
	})

	// CRUD: Remove brag
	t.Run("remove brag", func(t *testing.T) {
		output, err := runBinary([]string{"brag", "remove", "3", "--force"}, env)
		assert.NoError(t, err)
		assert.Contains(t, string(output), "Successfully removed")
	})

	// Validation: Invalid brags
	t.Run("invalid brag - short title", func(t *testing.T) {
		output, _ := runBinary([]string{"brag", "add", "--title", "Hi", "--description", "Valid description with enough characters here", "--category", "achievement"}, env)
		assert.Contains(t, string(output), "title must be at least 5 characters")
	})

	t.Run("invalid brag - short description", func(t *testing.T) {
		output, _ := runBinary([]string{"brag", "add", "--title", "Valid Title", "--description", "Short", "--category", "achievement"}, env)
		assert.Contains(t, string(output), "description must be at least 20 characters")
	})

	t.Run("invalid brag - invalid category", func(t *testing.T) {
		output, _ := runBinary([]string{"brag", "add", "--title", "Valid Title", "--description", "Valid description with enough characters", "--category", "invalid"}, env)
		assert.Contains(t, string(output), "invalid category")
	})

	// Tag CRUD: Create tags
	t.Run("create tags", func(t *testing.T) {
		tags := []string{"golang", "python", "aws"}
		for _, tag := range tags {
			output, err := runBinary([]string{"tag", "add", "--name", tag}, env)
			assert.NoError(t, err)
			assert.Contains(t, string(output), "Tag created successfully")
		}
	})

	// Tag CRUD: List tags
	t.Run("list tags", func(t *testing.T) {
		output, err := runBinary([]string{"tag", "list"}, env)
		assert.NoError(t, err)
		assert.Contains(t, string(output), "golang")
		assert.Contains(t, string(output), "python")
		assert.Contains(t, string(output), "aws")
	})

	// Tag CRUD: Remove tag
	t.Run("remove tag", func(t *testing.T) {
		// Remove tag by ID (golang should be ID 5 based on creation order)
		output, err := runBinary([]string{"tag", "remove", "5", "--force"}, env)
		assert.NoError(t, err)
		assert.Contains(t, string(output), "removed successfully")
	})

	// Tag Validation: Invalid tags
	t.Run("invalid tag - too short", func(t *testing.T) {
		output, _ := runBinary([]string{"tag", "add", "--name", "a"}, env)
		assert.Contains(t, string(output), "must be at least 2 characters")
	})

	t.Run("invalid tag - too long", func(t *testing.T) {
		output, _ := runBinary([]string{"tag", "add", "--name", "this-is-way-too-long-x"}, env)
		assert.Contains(t, string(output), "cannot exceed 20 characters")
	})

	t.Run("invalid tag - duplicate", func(t *testing.T) {
		output, _ := runBinary([]string{"tag", "add", "--name", "python"}, env)
		assert.Contains(t, string(output), "already exists")
	})

	// Generate document
	t.Run("generate document", func(t *testing.T) {
		output, err := runBinary([]string{"doc", "generate", "--output", filepath.Join(tmpDir, "bragdoc.md")}, env)
		assert.NoError(t, err)
		assert.Contains(t, string(output), "Document generated successfully")

		// Verify file was created
		// #nosec G304 - Test file path is controlled by test code
		content, err := os.ReadFile(filepath.Join(tmpDir, "bragdoc.md"))
		assert.NoError(t, err)
		assert.Contains(t, string(content), "Updated Achievement")
		assert.Contains(t, string(content), "Leadership Project")
	})
}
