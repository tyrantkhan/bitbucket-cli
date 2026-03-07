package git

import (
	"os/exec"
	"regexp"
	"strings"
)

var (
	// Match SSH remotes: git@bitbucket.org:workspace/repo.git
	sshRemoteRe = regexp.MustCompile(`git@bitbucket\.org:([^/]+)/([^/.]+?)(?:\.git)?$`)
	// Match HTTPS remotes: https://bitbucket.org/workspace/repo.git
	httpsRemoteRe = regexp.MustCompile(`https?://(?:[^@]+@)?bitbucket\.org/([^/]+)/([^/.]+?)(?:\.git)?$`)
)

// RepoContext holds a workspace and repo slug parsed from git remotes.
type RepoContext struct {
	Workspace string
	RepoSlug  string
}

// DetectRepo attempts to detect the Bitbucket workspace and repo from git remotes.
func DetectRepo() (*RepoContext, error) {
	out, err := exec.Command("git", "remote", "-v").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		url := fields[1]

		if ctx := parseRemoteURL(url); ctx != nil {
			return ctx, nil
		}
	}

	return nil, nil
}

func parseRemoteURL(url string) *RepoContext {
	if m := sshRemoteRe.FindStringSubmatch(url); len(m) == 3 {
		return &RepoContext{Workspace: m[1], RepoSlug: m[2]}
	}
	if m := httpsRemoteRe.FindStringSubmatch(url); len(m) == 3 {
		return &RepoContext{Workspace: m[1], RepoSlug: m[2]}
	}
	return nil
}

// CurrentBranch returns the current git branch name.
func CurrentBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
