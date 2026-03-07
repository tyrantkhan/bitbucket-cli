package repo

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdView() *cli.Command {
	return &cli.Command{
		Name:      "view",
		Usage:     "View repository details",
		ArgsUsage: "[slug]",
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

			workspace, repoSlug, err := cmdutil.ResolveWorkspaceAndRepo(ctx, cmd)
			if err != nil {
				return err
			}

			// Allow the slug as a positional argument.
			if cmd.Args().Len() > 0 {
				repoSlug = cmd.Args().First()
			}

			if repoSlug == "" {
				return fmt.Errorf("repository slug is required. Use --repo flag or pass it as an argument")
			}

			if err := api.ValidateSlug("repo", repoSlug); err != nil {
				return err
			}

			path := fmt.Sprintf("/2.0/repositories/%s/%s", workspace, repoSlug)

			resp, err := client.Get(path)
			if err != nil {
				return err
			}

			var repo models.Repository
			if err := api.DecodeJSON(resp, &repo); err != nil {
				return fmt.Errorf("failed to decode repository: %w", err)
			}

			// Open in browser if --web flag is set.
			if cmd.Bool("web") {
				url := ""
				if repo.Links.HTML != nil {
					url = repo.Links.HTML.Href
				}
				if url == "" {
					return fmt.Errorf("no web URL available for this repository")
				}
				fmt.Fprintln(f.IOOut, output.Muted.Render("Opening in browser: "+url))
				return exec.Command("open", url).Run()
			}

			format := cmdutil.GetFormat(ctx, cmd)

			if format == "json" {
				return output.RenderJSON(repo)
			}

			// Rich detail view.
			fmt.Fprintln(f.IOOut, output.Header.Render(repo.FullName))

			if repo.Description != "" {
				fmt.Fprintln(f.IOOut)
				output.RenderMarkdown(repo.Description)
			}

			fmt.Fprintln(f.IOOut)
			fmt.Fprintf(f.IOOut, "%s  %s\n", output.Bold.Render("Visibility:"), repo.Visibility())

			if repo.Language != "" {
				fmt.Fprintf(f.IOOut, "%s   %s\n", output.Bold.Render("Language:"), repo.Language)
			}

			if repo.MainBranch != nil {
				fmt.Fprintf(f.IOOut, "%s %s\n", output.Bold.Render("Main Branch:"), repo.MainBranch.Name)
			}

			if repo.Project != nil {
				fmt.Fprintf(f.IOOut, "%s    %s (%s)\n", output.Bold.Render("Project:"), repo.Project.Name, repo.Project.Key)
			}

			httpsURL := repo.CloneURL("https")
			sshURL := repo.CloneURL("ssh")

			if httpsURL != "" || sshURL != "" {
				fmt.Fprintln(f.IOOut)
				fmt.Fprintln(f.IOOut, output.Bold.Render("Clone URLs:"))
				if httpsURL != "" {
					fmt.Fprintf(f.IOOut, "  %s  %s\n", output.Muted.Render("HTTPS:"), httpsURL)
				}
				if sshURL != "" {
					fmt.Fprintf(f.IOOut, "  %s    %s\n", output.Muted.Render("SSH:"), sshURL)
				}
			}

			fmt.Fprintln(f.IOOut)
			fmt.Fprintf(f.IOOut, "%s    %s\n", output.Muted.Render("Created:"), models.FormatTime(repo.CreatedOn))
			fmt.Fprintf(f.IOOut, "%s    %s\n", output.Muted.Render("Updated:"), models.FormatTime(repo.UpdatedOn))

			return nil
		},
	}
}
