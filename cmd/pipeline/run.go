package pipeline

import (
	"context"
	"fmt"
	"strings"

	"charm.land/huh/v2"
	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/git"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdRun() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Run a pipeline",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
			cmdutil.FormatFlag,
			&cli.StringFlag{
				Name:  "branch",
				Usage: "Branch to run pipeline on (default: current git branch)",
			},
			&cli.StringFlag{
				Name:  "custom",
				Usage: "Custom pipeline name to run",
			},
			&cli.StringSliceFlag{
				Name:  "variable",
				Usage: "Pipeline variable as KEY=VALUE (repeatable)",
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

			branch := cmd.String("branch")
			custom := cmd.String("custom")
			variables := cmd.StringSlice("variable")

			// Default branch to current git branch.
			if branch == "" {
				b, err := git.CurrentBranch()
				if err != nil {
					return fmt.Errorf("could not detect current branch: %w. Use --branch to specify", err)
				}
				branch = b
			}

			// Build the request body.
			target := map[string]interface{}{
				"ref_type": "branch",
				"type":     "pipeline_ref_target",
				"ref_name": branch,
			}

			if custom != "" {
				target["selector"] = map[string]interface{}{
					"type":    "custom",
					"pattern": custom,
				}
			}

			body := map[string]interface{}{
				"target": target,
			}

			// Parse and add variables.
			if len(variables) > 0 {
				vars := make([]map[string]string, 0, len(variables))
				for _, v := range variables {
					parts := strings.SplitN(v, "=", 2)
					if len(parts) != 2 {
						return fmt.Errorf("invalid variable format %q, expected KEY=VALUE", v)
					}
					vars = append(vars, map[string]string{
						"key":   parts[0],
						"value": parts[1],
					})
				}
				body["variables"] = vars
			}

			// Confirmation prompt.
			description := fmt.Sprintf("Branch: %s", branch)
			if custom != "" {
				description = fmt.Sprintf("Branch: %s, Custom: %s", branch, custom)
			}

			var confirmed bool
			err = huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title("Run pipeline?").
						Description(description).
						Value(&confirmed),
				),
			).Run()
			if err != nil {
				return err
			}
			if !confirmed {
				fmt.Fprintln(f.IOOut, output.Warning.Render("Pipeline run cancelled."))
				return nil
			}

			path := fmt.Sprintf("/2.0/repositories/%s/%s/pipelines/", workspace, repo)

			resp, err := client.Post(path, body)
			if err != nil {
				return err
			}

			var p models.Pipeline
			if err := api.DecodeJSON(resp, &p); err != nil {
				return fmt.Errorf("failed to decode pipeline: %w", err)
			}

			format := cmdutil.GetFormat(ctx, cmd)
			if format == "json" {
				return output.RenderJSON(p)
			}

			status := p.StatusText()
			fmt.Fprintln(f.IOOut, output.Success.Render(
				fmt.Sprintf("Pipeline #%d triggered successfully!", p.BuildNumber),
			))
			fmt.Fprintln(f.IOOut)
			fmt.Fprintf(f.IOOut, "%s %s\n", output.Bold.Render("Branch:"), output.Cyan.Render(p.Target.RefName))
			fmt.Fprintf(f.IOOut, "%s %s\n", output.Bold.Render("Status:"), output.StatusColor(status).Render(status))

			return nil
		}),
	}
}
