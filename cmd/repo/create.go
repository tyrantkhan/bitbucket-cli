package repo

import (
	"context"
	"fmt"

	"charm.land/huh/v2"
	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdCreate() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new repository",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.FormatFlag,
			&cli.StringFlag{
				Name:  "name",
				Usage: "Repository name",
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Repository description",
			},
			&cli.BoolFlag{
				Name:  "private",
				Usage: "Make the repository private",
				Value: true,
			},
			&cli.StringFlag{
				Name:  "project",
				Usage: "Project key to assign the repository to",
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

			name := cmd.String("name")
			description := cmd.String("description")
			isPrivate := cmd.Bool("private")
			projectKey := cmd.String("project")

			// Interactive mode: prompt for values if name is not provided.
			if name == "" {
				err := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Repository Name").
							Description("The name for the new repository").
							Value(&name),
						huh.NewText().
							Title("Description").
							Description("A short description of the repository (optional)").
							Value(&description),
						huh.NewConfirm().
							Title("Private").
							Description("Should this repository be private?").
							Value(&isPrivate),
						huh.NewInput().
							Title("Project Key").
							Description("Bitbucket project key (optional)").
							Value(&projectKey),
					),
				).Run()
				if err != nil {
					return err
				}
			}

			if name == "" {
				return fmt.Errorf("repository name is required")
			}

			if err := api.ValidateSlug("repo", name); err != nil {
				return err
			}

			// Build the request body.
			body := map[string]interface{}{
				"scm":         "git",
				"is_private":  isPrivate,
				"description": description,
			}

			if projectKey != "" {
				body["project"] = map[string]string{
					"key": projectKey,
				}
			}

			path := fmt.Sprintf("/2.0/repositories/%s/%s", workspace, name)

			resp, err := client.Post(path, body)
			if err != nil {
				return err
			}

			var repo models.Repository
			if err := api.DecodeJSON(resp, &repo); err != nil {
				return fmt.Errorf("failed to decode created repository: %w", err)
			}

			format := cmdutil.GetFormat(ctx, cmd)
			if format == "json" {
				return output.RenderJSON(repo)
			}

			// Show success message and repository details.
			fmt.Fprintln(f.IOOut, output.Success.Render(
				fmt.Sprintf("Repository %s created successfully!", repo.FullName),
			))
			fmt.Fprintln(f.IOOut)
			fmt.Fprintf(f.IOOut, "%s  %s\n", output.Bold.Render("Name:"), repo.Name)
			fmt.Fprintf(f.IOOut, "%s  %s\n", output.Bold.Render("Slug:"), repo.Slug)
			fmt.Fprintf(f.IOOut, "%s  %s\n", output.Bold.Render("Visibility:"), repo.Visibility())

			if repo.Description != "" {
				fmt.Fprintf(f.IOOut, "%s  %s\n", output.Bold.Render("Description:"), repo.Description)
			}

			if repo.Project != nil {
				fmt.Fprintf(f.IOOut, "%s  %s (%s)\n", output.Bold.Render("Project:"), repo.Project.Name, repo.Project.Key)
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

			return nil
		}),
	}
}
