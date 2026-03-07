package repo

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
		Usage: "List repositories in a workspace",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.FormatFlag,
			cmdutil.LimitFlag,
		},
		Action: cmdutil.NoArgs(func(ctx context.Context, cmd *cli.Command) error {
			f := cmdutil.GetFactory(ctx)
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			workspace, err := cmdutil.ResolveWorkspace(ctx, cmd)
			if err != nil {
				return err
			}

			limit := int(cmd.Int("limit"))
			path := fmt.Sprintf("/2.0/repositories/%s", workspace)

			repos, err := api.Paginate[models.Repository](client, path, limit)
			if err != nil {
				return err
			}

			format := cmdutil.GetFormat(ctx, cmd)

			headers := []string{"Name", "Slug", "Visibility", "Language", "Updated"}
			rows := make([][]string, len(repos))
			for i, r := range repos {
				rows[i] = []string{
					r.Name,
					r.Slug,
					r.Visibility(),
					r.Language,
					models.FormatTime(r.UpdatedOn),
				}
			}

			return output.Format(format, repos, headers, rows)
		}),
	}
}
