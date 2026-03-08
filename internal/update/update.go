package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/tyrantkhan/bb/internal/config"
)

const checkInterval = 24 * time.Hour

type state struct {
	CheckedAt     int64  `json:"checked_at"`
	LatestVersion string `json:"latest_version"`
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

// CheckForUpdate checks whether a newer version of bb is available.
// Returns the latest version string if newer, or empty string otherwise.
func CheckForUpdate(currentVersion string) (string, error) {
	stateFile := config.StateFilePath()

	// Try reading cached state first.
	if data, err := os.ReadFile(stateFile); err == nil {
		var s state
		if json.Unmarshal(data, &s) == nil {
			if time.Since(time.Unix(s.CheckedAt, 0)) < checkInterval {
				if isNewer(s.LatestVersion, currentVersion) {
					return s.LatestVersion, nil
				}
				return "", nil
			}
		}
	}

	// Fetch latest release from GitHub.
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/tyrantkhan/bitbucket-cli/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github API returned %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	latest := strings.TrimPrefix(release.TagName, "v")

	// Persist state.
	s := state{
		CheckedAt:     time.Now().Unix(),
		LatestVersion: latest,
	}
	if data, err := json.Marshal(s); err == nil {
		_ = config.EnsureConfigDir()
		_ = os.WriteFile(stateFile, data, 0600)
	}

	if isNewer(latest, currentVersion) {
		return latest, nil
	}
	return "", nil
}

// UpgradeCommand returns a user-friendly upgrade command based on
// how bb was installed.
func UpgradeCommand() string {
	exe, err := os.Executable()
	if err != nil {
		return "https://github.com/tyrantkhan/bitbucket-cli/releases"
	}

	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		resolved = exe
	}

	switch {
	case strings.Contains(resolved, "/opt/homebrew/") || strings.Contains(resolved, "/usr/local/Cellar/"):
		return "brew upgrade bb"
	case strings.Contains(resolved, "/.local/share/mise/"):
		return "mise upgrade ubi:tyrantkhan/bitbucket-cli"
	case strings.Contains(resolved, "/go/bin/") || isInGoPath(resolved):
		return "go install github.com/tyrantkhan/bb@latest"
	case strings.Contains(resolved, "/.local/bin/"):
		return "curl -sSL https://tyrantkhan.github.io/bitbucket-cli/install.sh | sh"
	default:
		return "https://github.com/tyrantkhan/bitbucket-cli/releases"
	}
}

func isInGoPath(path string) bool {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return false
	}
	return strings.HasPrefix(path, filepath.Join(gopath, "bin"))
}

// isNewer returns true if latest is a higher semver than current.
func isNewer(latest, current string) bool {
	l := parseSemver(latest)
	c := parseSemver(current)
	if l == nil || c == nil {
		return false
	}
	for i := 0; i < 3; i++ {
		if l[i] > c[i] {
			return true
		}
		if l[i] < c[i] {
			return false
		}
	}
	return false
}

func parseSemver(v string) []int {
	v = strings.TrimPrefix(v, "v")
	// Strip any pre-release suffix (e.g. "0.0.5-rc1").
	if idx := strings.IndexByte(v, '-'); idx != -1 {
		v = v[:idx]
	}
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 3 {
		return nil
	}
	nums := make([]int, 3)
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil
		}
		nums[i] = n
	}
	return nums
}
