package pipeline

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
		Usage: "List pipelines",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
			cmdutil.FormatFlag,
			cmdutil.LimitFlag,
		},
		Action: cmdutil.NoArgs(func(ctx context.Context, cmd *cli.Command) error {
			f := cmdutil.GetFactory(ctx)
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			workspace, repo, err := cmdutil.ResolveWorkspaceAndRepo(ctx, cmd)
			if err != nil {
				return err
			}

			limit := int(cmd.Int("limit"))

			path := fmt.Sprintf("/2.0/repositories/%s/%s/pipelines/?sort=-created_on", workspace, repo)

			pipelines, err := api.Paginate[models.Pipeline](client, path, limit)
			if err != nil {
				return err
			}

			format := cmdutil.GetFormat(ctx, cmd)

			headers := []string{"#", "Branch", "Status", "Duration", "Creator", "Created"}
			rows := make([][]string, len(pipelines))
			for i, p := range pipelines {
				status := p.StatusText()
				rows[i] = []string{
					fmt.Sprintf("%d", p.BuildNumber),
					p.Target.RefName,
					output.StatusColor(status).Render(status),
					p.Duration(),
					p.Creator.DisplayName,
					models.FormatTime(p.CreatedOn),
				}
			}

			return output.Format(format, pipelines, headers, rows)
		}),
	}
}
