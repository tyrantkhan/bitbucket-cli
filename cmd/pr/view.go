package pr

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdView() *cli.Command {
	return &cli.Command{
		Name:      "view",
		Usage:     "View pull request details",
		ArgsUsage: "<id>",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
			cmdutil.FormatFlag,
			cmdutil.WebFlag,
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

			// Open in browser if --web flag is set.
			if cmd.Bool("web") {
				url := ""
				if pr.Links.HTML != nil {
					url = pr.Links.HTML.Href
				}
				if url == "" {
					return fmt.Errorf("no web URL available for this pull request")
				}
				fmt.Fprintln(f.IOOut, output.Muted.Render("Opening in browser: "+url))
				return exec.Command("open", url).Run()
			}

			format := cmdutil.GetFormat(ctx, cmd)

			if format == "json" {
				return output.RenderJSON(pr)
			}

			// Rich detail view.
			fmt.Fprintf(f.IOOut, "%s #%d\n", output.Header.Render(pr.Title), pr.ID)
			fmt.Fprintf(f.IOOut, "%s  %s\n", output.Bold.Render("State:"), output.StatusColor(pr.State).Render(pr.State))
			fmt.Fprintf(f.IOOut, "%s %s\n", output.Bold.Render("Author:"), pr.Author.DisplayName)
			fmt.Fprintf(f.IOOut, "%s %s → %s\n",
				output.Bold.Render("Branch:"),
				output.Cyan.Render(pr.Source.Branch.Name),
				output.Cyan.Render(pr.Destination.Branch.Name),
			)

			// Description
			if pr.Description != "" {
				fmt.Fprintln(f.IOOut)
				output.RenderMarkdown(pr.Description)
			}

			// Reviewers
			if len(pr.Participants) > 0 {
				fmt.Fprintln(f.IOOut)
				fmt.Fprintln(f.IOOut, output.Bold.Render("Reviewers:"))
				for _, p := range pr.Participants {
					if p.Role != "REVIEWER" {
						continue
					}
					status := "pending"
					statusStyle := output.Yellow
					if p.Approved {
						status = "approved"
						statusStyle = output.Green
					} else if p.State == "changes_requested" {
						status = "changes requested"
						statusStyle = output.Red
					}
					fmt.Fprintf(f.IOOut, "  %s  %s\n",
						p.User.DisplayName,
						statusStyle.Render(status),
					)
				}
			}

			// Counts
			fmt.Fprintln(f.IOOut)
			fmt.Fprintf(f.IOOut, "%s   %d\n", output.Bold.Render("Tasks:"), pr.TaskCount)
			fmt.Fprintf(f.IOOut, "%s %d\n", output.Bold.Render("Comments:"), pr.CommentCount)

			// Timestamps
			fmt.Fprintln(f.IOOut)
			fmt.Fprintf(f.IOOut, "%s  %s\n", output.Muted.Render("Created:"), models.FormatTime(pr.CreatedOn))
			fmt.Fprintf(f.IOOut, "%s  %s\n", output.Muted.Render("Updated:"), models.FormatTime(pr.UpdatedOn))

			return nil
		},
	}
}
