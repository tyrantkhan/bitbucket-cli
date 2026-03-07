package pr

import (
	"context"
	"fmt"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdList() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List pull requests",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
			cmdutil.FormatFlag,
			cmdutil.LimitFlag,
			&cli.StringFlag{
				Name:  "state",
				Usage: "PR state: OPEN, MERGED, DECLINED",
				Value: "OPEN",
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

			state := cmd.String("state")
			limit := int(cmd.Int("limit"))

			path := fmt.Sprintf("/2.0/repositories/%s/%s/pullrequests?state=%s", workspace, repo, state)

			prs, err := api.Paginate[models.PullRequest](client, path, limit)
			if err != nil {
				return err
			}

			format := cmdutil.GetFormat(ctx, cmd)

			headers := []string{"ID", "Title", "Author", "Branch", "State", "Updated"}
			rows := make([][]string, len(prs))
			for i, pr := range prs {
				branch := fmt.Sprintf("%s → %s", pr.Source.Branch.Name, pr.Destination.Branch.Name)
				rows[i] = []string{
					fmt.Sprintf("%d", pr.ID),
					pr.Title,
					pr.Author.DisplayName,
					branch,
					output.StatusColor(pr.State).Render(pr.State),
					models.FormatTime(pr.UpdatedOn),
				}
			}

			return output.Format(format, prs, headers, rows)
		},
	}
}
