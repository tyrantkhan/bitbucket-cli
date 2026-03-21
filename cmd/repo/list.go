package repo

import (
	"context"
	"fmt"
	"net/url"

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
			&cli.StringFlag{
				Name:    "project",
				Aliases: []string{"p"},
				Usage:   "Filter by project key",
			},
			&cli.StringFlag{
				Name:  "exclude-project",
				Usage: "Exclude repos in this project",
			},
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

			project := cmd.String("project")
			excludeProject := cmd.String("exclude-project")
			if project != "" && excludeProject != "" {
				return fmt.Errorf("cannot use --project and --exclude-project together")
			}

			limit := int(cmd.Int("limit"))
			path := fmt.Sprintf("/2.0/repositories/%s", workspace)

			if project != "" {
				path += fmt.Sprintf("?q=%s", url.QueryEscape(fmt.Sprintf(`project.key="%s"`, project)))
			} else if excludeProject != "" {
				path += fmt.Sprintf("?q=%s", url.QueryEscape(fmt.Sprintf(`project.key!="%s"`, excludeProject)))
			}

			repos, err := api.Paginate[models.Repository](client, path, limit)
			if err != nil {
				return err
			}

			format := cmdutil.GetFormat(ctx, cmd)

			headers := []string{"Name", "Project", "Slug", "Visibility", "Language", "Updated"}
			rows := make([][]string, len(repos))
			for i, r := range repos {
				projectKey := ""
				if r.Project != nil {
					projectKey = r.Project.Key
				}
				rows[i] = []string{
					r.Name,
					projectKey,
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
