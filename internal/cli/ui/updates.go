package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func CheckForUpdates(currentVersion string) {
	isDebug := os.Getenv("PLUTUS_DEBUG_UPDATE") != ""

	if currentVersion == "dev" && !isDebug {
		return
	}

	home, _ := os.UserHomeDir()
	cachePath := filepath.Join(home, ".plutus_update_check")

	if info, err := os.Stat(cachePath); err == nil {
		if time.Since(info.ModTime()) < 24*time.Hour {
			return
		}
	}

	client := &http.Client{
		Timeout: 1200 * time.Millisecond,
	}

	req, _ := http.NewRequest("GET", "https://api.github.com/repos/j-lewandowski/plutus-cli/releases/latest", nil)
	req.Header.Set("User-Agent", "plutus-cli-updater")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return
	}

	if release.TagName != currentVersion {
		fmt.Printf("\n\033[33m🚀 A new version of Plutus is available: %s (current: %s)\033[0m\n", release.TagName, currentVersion)
		fmt.Println("\033[33m👉 To update, run:\033[0m")

		installCmd := "curl -fsSL https://raw.githubusercontent.com/j-lewandowski/plutus-cli/main/install.sh | bash"
		fmt.Printf("\033[36m%s\033[0m\n\n", installCmd)
	}

	os.WriteFile(cachePath, []byte(release.TagName), 0644)
}
