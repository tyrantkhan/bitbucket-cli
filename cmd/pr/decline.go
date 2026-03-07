package pr

import (
	"context"
	"fmt"
	"strconv"

	"charm.land/huh/v2"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdDecline() *cli.Command {
	return &cli.Command{
		Name:      "decline",
		Usage:     "Decline a pull request",
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

			// Confirmation prompt.
			var confirmed bool
			err = huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title(fmt.Sprintf("Decline pull request #%d?", prID)).
						Description("This action cannot be easily undone.").
						Value(&confirmed),
				),
			).Run()
			if err != nil {
				return err
			}
			if !confirmed {
				fmt.Fprintln(f.IOOut, output.Warning.Render("Decline cancelled."))
				return nil
			}

			path := fmt.Sprintf("/2.0/repositories/%s/%s/pullrequests/%d/decline", workspace, repo, prID)

			resp, err := client.Post(path, nil)
			if err != nil {
				return err
			}
			resp.Body.Close()

			fmt.Fprintln(f.IOOut, output.Success.Render(
				fmt.Sprintf("Pull request #%d declined.", prID),
			))

			return nil
		},
	}
}
