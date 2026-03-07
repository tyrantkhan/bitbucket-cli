package repo

import (
	"context"
	"fmt"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/git"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/urfave/cli/v3"
)

func newCmdClone() *cli.Command {
	return &cli.Command{
		Name:      "clone",
		Usage:     "Clone a repository",
		ArgsUsage: "<slug> [directory]",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			&cli.StringFlag{
				Name:  "protocol",
				Usage: "Clone protocol: https, ssh",
				Value: "https",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			f := cmdutil.GetFactory(ctx)
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			workspace, err := cmdutil.ResolveWorkspace(ctx, cmd)
			if err != nil {
				return err
			}

			repoSlug := cmd.Args().First()
			if repoSlug == "" {
				return fmt.Errorf("repository slug is required. Usage: bb repo clone <slug> [directory]")
			}

			if err := api.ValidateSlug("repo", repoSlug); err != nil {
				return err
			}

			// Optional directory argument.
			dir := ""
			if cmd.Args().Len() > 1 {
				dir = cmd.Args().Get(1)
			}

			protocol := cmd.String("protocol")
			if protocol != "https" && protocol != "ssh" {
				return fmt.Errorf("invalid protocol %q: must be \"https\" or \"ssh\"", protocol)
			}

			// Fetch repo details to get the clone URL.
			path := fmt.Sprintf("/2.0/repositories/%s/%s", workspace, repoSlug)
			resp, err := client.Get(path)
			if err != nil {
				return err
			}

			var repo models.Repository
			if err := api.DecodeJSON(resp, &repo); err != nil {
				return fmt.Errorf("failed to decode repository: %w", err)
			}

			cloneURL := repo.CloneURL(protocol)
			if cloneURL == "" {
				return fmt.Errorf("no %s clone URL found for %s", protocol, repo.FullName)
			}

			return git.Clone(cloneURL, dir)
		},
	}
}
