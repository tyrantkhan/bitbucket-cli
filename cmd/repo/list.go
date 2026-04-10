package repo

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

type repoDetails struct {
	LastCommit string
	OpenPRs    int
}

func fetchRepoDetails(client *api.Client, workspace, slug string) repoDetails {
	var d repoDetails

	// Last commit on default branch.
	commits, err := api.Paginate[models.Commit](client, fmt.Sprintf("/2.0/repositories/%s/%s/commits", workspace, slug), 1)
	if err == nil && len(commits) > 0 {
		d.LastCommit = commits[0].Date
	}

	// Open PR count.
	count, err := api.Count(client, fmt.Sprintf("/2.0/repositories/%s/%s/pullrequests?state=OPEN", workspace, slug))
	if err == nil {
		d.OpenPRs = count
	}

	return d
}

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
			&cli.BoolFlag{
				Name:  "details",
				Usage: "Include last commit date and open PR count (slower)",
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

			var q string
			if project != "" {
				q = fmt.Sprintf(`project.key="%s"`, project)
			} else if excludeProject != "" {
				q = fmt.Sprintf(`project.key!="%s"`, excludeProject)
			}
			if q != "" {
				path += "?q=" + url.QueryEscape(q)
			}

			repos, err := api.Paginate[models.Repository](client, path, limit)
			if err != nil {
				return err
			}

			format := cmdutil.GetFormat(ctx, cmd)
			details := cmd.Bool("details")

			if details {
				// Fetch details concurrently with bounded parallelism.
				detailsMap := make([]repoDetails, len(repos))
				var wg sync.WaitGroup
				sem := make(chan struct{}, 10)

				for i, r := range repos {
					wg.Add(1)
					go func(idx int, slug string) {
						defer wg.Done()
						sem <- struct{}{}
						detailsMap[idx] = fetchRepoDetails(client, workspace, slug)
						<-sem
					}(i, r.Slug)
				}
				wg.Wait()

				headers := []string{"Name", "Project", "Slug", "Language", "Last Commit", "Open PRs", "Updated"}
				rows := make([][]string, len(repos))
				for i, r := range repos {
					projectKey := ""
					if r.Project != nil {
						projectKey = r.Project.Key
					}
					lastCommit := "-"
					if detailsMap[i].LastCommit != "" {
						lastCommit = models.FormatTime(detailsMap[i].LastCommit)
					}
					rows[i] = []string{
						r.Name,
						projectKey,
						r.Slug,
						r.Language,
						lastCommit,
						fmt.Sprintf("%d", detailsMap[i].OpenPRs),
						models.FormatTime(r.UpdatedOn),
					}
				}
				return output.Format(format, repos, headers, rows)
			}

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
