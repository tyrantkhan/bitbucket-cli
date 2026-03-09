package pr

import (
	"context"
	"fmt"
	"strconv"

	"charm.land/huh/v2"
	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdMerge() *cli.Command {
	return &cli.Command{
		Name:      "merge",
		Usage:     "Merge a pull request",
		ArgsUsage: "<id>",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
			cmdutil.FormatFlag,
			&cli.StringFlag{
				Name:  "strategy",
				Usage: "Merge strategy: merge_commit, squash, fast_forward",
				Value: "merge_commit",
			},
			&cli.StringFlag{
				Name:  "message",
				Usage: "Merge commit message",
			},
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

			strategy := cmd.String("strategy")
			message := cmd.String("message")

			// Confirmation prompt.
			var confirmed bool
			err = huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title(fmt.Sprintf("Merge pull request #%d?", prID)).
						Description(fmt.Sprintf("Strategy: %s", strategy)).
						Value(&confirmed),
				),
			).Run()
			if err != nil {
				return err
			}
			if !confirmed {
				fmt.Fprintln(f.IOOut, output.Warning.Render("Merge cancelled."))
				return nil
			}

			// Build the merge request body.
			body := map[string]interface{}{
				"merge_strategy":      strategy,
				"close_source_branch": true,
			}
			if message != "" {
				body["message"] = message
			}

			path := fmt.Sprintf("/2.0/repositories/%s/%s/pullrequests/%d/merge", workspace, repo, prID)

			resp, err := client.Post(path, body)
			if err != nil {
				return err
			}

			var pr models.PullRequest
			if err := api.DecodeJSON(resp, &pr); err != nil {
				return fmt.Errorf("failed to decode merge response: %w", err)
			}

			format := cmdutil.GetFormat(ctx, cmd)
			if format == "json" {
				return output.RenderJSON(pr)
			}

			fmt.Fprintln(f.IOOut, output.Success.Render(
				fmt.Sprintf("Pull request #%d merged successfully!", pr.ID),
			))

			return nil
		},
	}
}
