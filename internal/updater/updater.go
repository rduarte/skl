package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	RepoOwner = "rduarte"
	RepoName  = "skl"
)

// Release represents a GitHub release.
type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// FetchLatestRelease fetches the latest release from GitHub API.
func FetchLatestRelease(timeout time.Duration) (*Release, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", RepoOwner, RepoName)

	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API status: %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// CheckLatestVersion is a light version of FetchLatestRelease with 1.5s timeout.
func CheckLatestVersion() (string, error) {
	rel, err := FetchLatestRelease(1500 * time.Millisecond)
	if err != nil {
		return "", err
	}
	return rel.TagName, nil
}
