package cmdutil

import (
	"context"
	"fmt"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/git"
	"github.com/urfave/cli/v3"
)

// ResolveWorkspaceAndRepo resolves workspace and repo from flags, git remote, or config.
func ResolveWorkspaceAndRepo(ctx context.Context, cmd *cli.Command) (string, string, error) {
	f := GetFactory(ctx)

	workspace := cmd.String("workspace")
	repo := cmd.String("repo")

	// Try git remote detection for missing values
	if workspace == "" || repo == "" {
		if detected, _ := git.DetectRepo(); detected != nil {
			if workspace == "" {
				workspace = detected.Workspace
			}
			if repo == "" {
				repo = detected.RepoSlug
			}
		}
	}

	// Fall back to config default workspace
	if workspace == "" && f != nil && f.Config.DefaultWorkspace != "" {
		workspace = f.Config.DefaultWorkspace
	}

	if workspace == "" {
		return "", "", fmt.Errorf("workspace is required. Use --workspace flag, or run from a Bitbucket repo")
	}

	if err := api.ValidateSlug("workspace", workspace); err != nil {
		return "", "", err
	}
	if repo == "" {
		return "", "", fmt.Errorf("repository is required. Use --repo flag, or run from a Bitbucket repo")
	}

	if err := api.ValidateSlug("repo", repo); err != nil {
		return "", "", err
	}

	return workspace, repo, nil
}

// ResolveWorkspace resolves just the workspace (for commands that don't need a repo).
func ResolveWorkspace(ctx context.Context, cmd *cli.Command) (string, error) {
	f := GetFactory(ctx)

	workspace := cmd.String("workspace")

	if workspace == "" {
		if detected, _ := git.DetectRepo(); detected != nil {
			workspace = detected.Workspace
		}
	}

	if workspace == "" && f != nil && f.Config.DefaultWorkspace != "" {
		workspace = f.Config.DefaultWorkspace
	}

	if workspace == "" {
		return "", fmt.Errorf("workspace is required. Use --workspace flag, or run from a Bitbucket repo")
	}

	if err := api.ValidateSlug("workspace", workspace); err != nil {
		return "", err
	}

	return workspace, nil
}

// GetFormat resolves the output format from flags or config.
func GetFormat(ctx context.Context, cmd *cli.Command) string {
	format := cmd.String("format")
	if format != "" {
		return format
	}
	f := GetFactory(ctx)
	if f != nil && f.Config.DefaultFormat != "" {
		return f.Config.DefaultFormat
	}
	return "table"
}
