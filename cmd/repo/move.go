package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdMove() *cli.Command {
	return &cli.Command{
		Name:      "move",
		Usage:     "Move a repository to a different project",
		ArgsUsage: "<slug>",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.FormatFlag,
			&cli.StringFlag{
				Name:  "project",
				Usage: "Destination project key",
			},
			&cli.StringFlag{
				Name:  "prefix",
				Usage: "Prefix to add to the repo name (e.g. \"Archived-\")",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			f := cmdutil.GetFactory(ctx)
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			workspace, err := cmdutil.ResolveWorkspace(ctx, cmd)
			if err != nil {
				return err
			}

			slug := cmd.Args().First()
			if slug == "" {
				return fmt.Errorf("repository slug is required")
			}

			if err := api.ValidateSlug("repo", slug); err != nil {
				return err
			}

			projectKey := cmd.String("project")
			if projectKey == "" {
				return fmt.Errorf("--project flag is required")
			}

			// Fetch current repo to get existing name.
			path := fmt.Sprintf("/2.0/repositories/%s/%s", workspace, slug)
			resp, err := client.Get(path)
			if err != nil {
				return err
			}

			var repo models.Repository
			if err := api.DecodeJSON(resp, &repo); err != nil {
				return fmt.Errorf("failed to decode repository: %w", err)
			}

			body := map[string]interface{}{
				"project": map[string]string{
					"key": projectKey,
				},
			}

			// Add prefix to name if specified and not already present.
			prefix := cmd.String("prefix")
			newName := repo.Name
			if prefix != "" && !strings.HasPrefix(repo.Name, prefix) {
				newName = prefix + repo.Name
				body["name"] = newName
			}

			resp, err = client.Put(path, body)
			if err != nil {
				return err
			}

			var updated models.Repository
			if err := api.DecodeJSON(resp, &updated); err != nil {
				return fmt.Errorf("failed to decode updated repository: %w", err)
			}

			format := cmdutil.GetFormat(ctx, cmd)
			if format == "json" {
				return output.RenderJSON(updated)
			}

			fmt.Fprintln(f.IOOut, output.Success.Render(
				fmt.Sprintf("Repository moved: %s → project %s", updated.FullName, projectKey),
			))
			if newName != repo.Name {
				fmt.Fprintf(f.IOOut, "%s  %s → %s\n", output.Bold.Render("Renamed:"), repo.Name, updated.Name)
			}
			fmt.Fprintf(f.IOOut, "%s  %s\n", output.Bold.Render("Slug:"), updated.Slug)
			if updated.Project != nil {
				fmt.Fprintf(f.IOOut, "%s  %s (%s)\n", output.Bold.Render("Project:"), updated.Project.Name, updated.Project.Key)
			}

			return nil
		},
	}
}
