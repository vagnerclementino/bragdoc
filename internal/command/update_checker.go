package command

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vagnerclementino/bragdoc/config"
)

// CheckForUpdates verifies if a new version is available
func CheckForUpdates(cfg *config.Config, configMgr *config.Manager) {
	// Skip if disabled
	if !cfg.UpdateChecker.Enabled {
		return
	}

	// Check only once per day
	if time.Since(cfg.UpdateChecker.LastCheckedAt) < 24*time.Hour {
		return
	}

	// Update last checked timestamp
	cfg.UpdateChecker.LastCheckedAt = time.Now()
	_ = configMgr.Save(context.Background(), cfg) // Ignore error

	// Fetch latest version
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", githubAPIURL, nil)
	if err != nil {
		return
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close() // Ignore close error for update checker
	}()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return
	}

	// Compare versions
	currentVersion := strings.TrimPrefix(Version, "v")
	latestVersion := strings.TrimPrefix(release.TagName, "v")

	if currentVersion != "unknown" && currentVersion != latestVersion {
		fmt.Printf("\n💡 New version available: %s (current: %s)\n", release.TagName, Version)
		fmt.Println("   Run 'bragdoc version upgrade' to update")
	}
}
