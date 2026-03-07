package pr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdApprove() *cli.Command {
	return &cli.Command{
		Name:      "approve",
		Usage:     "Approve a pull request",
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

			path := fmt.Sprintf("/2.0/repositories/%s/%s/pullrequests/%d/approve", workspace, repo, prID)

			resp, err := client.Post(path, nil)
			if err != nil {
				return err
			}
			_ = resp.Body.Close()

			fmt.Fprintln(f.IOOut, output.Success.Render(
				fmt.Sprintf("Pull request #%d approved.", prID),
			))

			return nil
		},
	}
}
