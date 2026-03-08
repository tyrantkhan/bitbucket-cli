package pr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdReady() *cli.Command {
	return &cli.Command{
		Name:      "ready",
		Usage:     "Mark a draft pull request as ready for review",
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

			resp, err := client.Get(path)
			if err != nil {
				return err
			}

			var pr models.PullRequest
			if err := api.DecodeJSON(resp, &pr); err != nil {
				return fmt.Errorf("failed to decode pull request: %w", err)
			}

			if !pr.Draft {
				fmt.Fprintln(f.IOOut, output.Muted.Render(
					fmt.Sprintf("Pull request #%d is already ready for review.", prID),
				))
				return nil
			}

			body := buildPRUpdateBody(pr)
			body["draft"] = false

			resp, err = client.Put(path, body)
			if err != nil {
				return err
			}
			_ = resp.Body.Close()

			fmt.Fprintln(f.IOOut, output.Success.Render(
				fmt.Sprintf("Pull request #%d is now ready for review.", prID),
			))

			return nil
		},
	}
}
