package pr

import (
	"context"
	"fmt"

	"charm.land/huh/v2"
	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/git"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdCreate() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a pull request",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
			&cli.StringFlag{
				Name:  "title",
				Usage: "Pull request title",
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Pull request description",
			},
			&cli.StringFlag{
				Name:  "source",
				Usage: "Source branch name (default: current git branch)",
			},
			&cli.StringFlag{
				Name:  "destination",
				Usage: "Destination branch name",
				Value: "main",
			},
			&cli.BoolFlag{
				Name:  "close-source-branch",
				Usage: "Close the source branch after merge",
			},
			&cli.StringSliceFlag{
				Name:  "reviewer",
				Usage: "Reviewer UUID (repeatable)",
			},
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

			title := cmd.String("title")
			description := cmd.String("description")
			source := cmd.String("source")
			destination := cmd.String("destination")
			closeSource := cmd.Bool("close-source-branch")
			reviewerUUIDs := cmd.StringSlice("reviewer")

			// Default source to current git branch.
			if source == "" {
				branch, err := git.CurrentBranch()
				if err != nil {
					return fmt.Errorf("could not detect current branch: %w. Use --source to specify", err)
				}
				source = branch
			}

			// Interactive mode: prompt for values if title is not provided.
			if title == "" {
				err := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Title").
							Description("Pull request title").
							Value(&title),
						huh.NewText().
							Title("Description").
							Description("Pull request description (optional)").
							Value(&description),
						huh.NewInput().
							Title("Source Branch").
							Description("The branch to merge from").
							Value(&source),
						huh.NewInput().
							Title("Destination Branch").
							Description("The branch to merge into").
							Value(&destination),
						huh.NewConfirm().
							Title("Close source branch after merge?").
							Value(&closeSource),
					),
				).Run()
				if err != nil {
					return err
				}
			}

			if title == "" {
				return fmt.Errorf("pull request title is required")
			}

			// Build the request body.
			body := map[string]interface{}{
				"title":       title,
				"description": description,
				"source": map[string]interface{}{
					"branch": map[string]string{
						"name": source,
					},
				},
				"destination": map[string]interface{}{
					"branch": map[string]string{
						"name": destination,
					},
				},
				"close_source_branch": closeSource,
			}

			if len(reviewerUUIDs) > 0 {
				reviewers := make([]map[string]string, len(reviewerUUIDs))
				for i, uuid := range reviewerUUIDs {
					reviewers[i] = map[string]string{"uuid": uuid}
				}
				body["reviewers"] = reviewers
			}

			path := fmt.Sprintf("/2.0/repositories/%s/%s/pullrequests", workspace, repo)

			resp, err := client.Post(path, body)
			if err != nil {
				return err
			}

			var pr models.PullRequest
			if err := api.DecodeJSON(resp, &pr); err != nil {
				return fmt.Errorf("failed to decode pull request: %w", err)
			}

			fmt.Fprintln(f.IOOut, output.Success.Render(
				fmt.Sprintf("Pull request #%d created successfully!", pr.ID),
			))
			fmt.Fprintln(f.IOOut)
			fmt.Fprintf(f.IOOut, "%s  %s\n", output.Bold.Render("Title:"), pr.Title)
			fmt.Fprintf(f.IOOut, "%s %s → %s\n",
				output.Bold.Render("Branch:"),
				output.Cyan.Render(pr.Source.Branch.Name),
				output.Cyan.Render(pr.Destination.Branch.Name),
			)

			if pr.Links.HTML != nil {
				fmt.Fprintf(f.IOOut, "%s    %s\n", output.Bold.Render("URL:"), pr.Links.HTML.Href)
			}

			return nil
		}),
	}
}
