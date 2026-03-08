package pr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/urfave/cli/v3"
)

func newCmdDraft() *cli.Command {
	return &cli.Command{
		Name:      "draft",
		Usage:     "Convert a pull request to draft",
		ArgsUsage: "<id>",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			f := cmdutil.GetFactory(ctx)
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			workspace, repo, err := cmdutil.ResolveWorkspaceAndRepo(ctx, cmd)
			if err != nil {
				return err
			}

			idStr := cmd.Args().First()
			if idStr == "" {
				return fmt.Errorf("pull request ID is required")
			}
			prID, err := strconv.Atoi(idStr)
			if err != nil {
				return fmt.Errorf("invalid pull request ID: %s", idStr)
			}

			path := fmt.Sprintf("/2.0/repositories/%s/%s/pullrequests/%d", workspace, repo, prID)
			return setDraftStatus(client, f.IOOut, path, prID, true)
		},
	}
}
