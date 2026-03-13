// Package command provides version upgrade functionality
//
//nolint:errcheck,gosec // Upgrade code with acceptable error handling
package command

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	githubAPIURL  = "https://api.github.com/repos/vagnerclementino/bragdoc/releases/latest"
	githubTimeout = 30 * time.Second
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// NewVersionUpgradeCmd creates a command to upgrade bragdoc
func NewVersionUpgradeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade bragdoc to the latest version",
		Long: `Check for the latest version on GitHub and upgrade bragdoc.
This will download the latest release, replace the current binary, and run any pending migrations.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runVersionUpgrade(cmd.Context())
		},
	}
}

func runVersionUpgrade(ctx context.Context) error {
	fmt.Println("🔍 Checking for updates...")

	// Get current version
	currentVersion := Version
	if currentVersion == "" {
		currentVersion = "dev"
	}

	// Fetch latest release from GitHub
	release, err := fetchLatestRelease(ctx)
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersionClean := strings.TrimPrefix(currentVersion, "v")

	fmt.Printf("📦 Current version: %s\n", currentVersion)
	fmt.Printf("📦 Latest version: %s\n", release.TagName)

	if currentVersionClean == latestVersion {
		fmt.Println("✅ You are already running the latest version!")
		return nil
	}

	// Find appropriate asset for current platform
	assetURL, assetName, err := findAssetForPlatform(release)
	if err != nil {
		return err
	}

	fmt.Printf("\n⬇️  Downloading %s...\n", assetName)

	// Download the asset
	tmpFile, err := downloadAsset(ctx, assetURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer func() {
		if err := os.Remove(tmpFile); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to remove temp file: %v\n", err)
		}
	}()

	fmt.Println("📦 Extracting binary...")

	// Extract binary
	binaryPath, err := extractBinary(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to extract binary: %w", err)
	}
	defer func() {
		if err := os.Remove(binaryPath); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to remove temp binary: %v\n", err)
		}
	}()

	// Get current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	fmt.Println("🔄 Replacing binary...")

	// Replace current binary
	if err := replaceBinary(currentExe, binaryPath); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	fmt.Printf("\n✅ Successfully upgraded to version %s!\n", release.TagName)
	fmt.Println("🔄 Migrations will run automatically on next command execution.")
	fmt.Println("\n💡 Restart your terminal or run any bragdoc command to complete the upgrade.")

	return nil
}

func fetchLatestRelease(ctx context.Context) (*GitHubRelease, error) {
	client := &http.Client{Timeout: githubTimeout}

	req, err := http.NewRequestWithContext(ctx, "GET", githubAPIURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func findAssetForPlatform(release *GitHubRelease) (string, string, error) {
	platform := runtime.GOOS
	arch := runtime.GOARCH

	// Map Go arch to release naming
	archMap := map[string]string{
		"amd64": "x86_64",
		"arm64": "arm64",
	}

	releaseArch, ok := archMap[arch]
	if !ok {
		return "", "", fmt.Errorf("unsupported architecture: %s", arch)
	}

	// Expected format: bragdoc_v1.0.0_darwin_x86_64.tar.gz
	expectedPattern := fmt.Sprintf("bragdoc_%s_%s_%s.tar.gz", release.TagName, platform, releaseArch)

	for _, asset := range release.Assets {
		if asset.Name == expectedPattern {
			return asset.BrowserDownloadURL, asset.Name, nil
		}
	}

	return "", "", fmt.Errorf("no release found for %s/%s", platform, arch)
}

func downloadAsset(ctx context.Context, url string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Minute}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "bragdoc-*.tar.gz")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

func extractBinary(tarGzPath string) (string, error) {
	file, err := os.Open(tarGzPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// Look for the bragdoc binary
		if header.Name == "bragdoc" || filepath.Base(header.Name) == "bragdoc" {
			tmpFile, err := os.CreateTemp("", "bragdoc-new-*")
			if err != nil {
				return "", err
			}
			defer tmpFile.Close()

			if _, err := io.Copy(tmpFile, tr); err != nil {
				os.Remove(tmpFile.Name())
				return "", err
			}

			// Make executable
			if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
				os.Remove(tmpFile.Name())
				return "", err
			}

			return tmpFile.Name(), nil
		}
	}

	return "", fmt.Errorf("binary not found in archive")
}

func replaceBinary(currentPath, newPath string) error {
	// Backup current binary
	backupPath := currentPath + ".backup"
	if err := copyFile(currentPath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Replace with new binary
	if err := copyFile(newPath, currentPath); err != nil {
		// Restore backup on failure
		copyFile(backupPath, currentPath)
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	// Remove backup on success
	os.Remove(backupPath)

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Copy permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}
