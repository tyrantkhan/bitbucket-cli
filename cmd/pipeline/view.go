package pipeline

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdView() *cli.Command {
	return &cli.Command{
		Name:      "view",
		Usage:     "View pipeline details",
		ArgsUsage: "<uuid>",
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

			pipelineUUID := cmd.Args().First()
			if pipelineUUID == "" {
				return fmt.Errorf("pipeline UUID is required")
			}

			path := fmt.Sprintf("/2.0/repositories/%s/%s/pipelines/%s", workspace, repo, pipelineUUID)

			resp, err := client.Get(path)
			if err != nil {
				return err
			}

			var p models.Pipeline
			if err := api.DecodeJSON(resp, &p); err != nil {
				return fmt.Errorf("failed to decode pipeline: %w", err)
			}

			// Open in browser if --web flag is set.
			if cmd.Bool("web") {
				url := ""
				if p.Links.HTML != nil {
					url = p.Links.HTML.Href
				}
				if url == "" {
					return fmt.Errorf("no web URL available for this pipeline")
				}
				fmt.Fprintln(f.IOOut, output.Muted.Render("Opening in browser: "+url))
				return exec.Command("open", url).Run()
			}

			format := cmdutil.GetFormat(ctx, cmd)

			if format == "json" {
				return output.RenderJSON(p)
			}

			// Rich detail view.
			status := p.StatusText()
			fmt.Fprintf(f.IOOut, "%s\n", output.Header.Render(fmt.Sprintf("Pipeline #%d", p.BuildNumber)))
			fmt.Fprintf(f.IOOut, "%s  %s\n", output.Bold.Render("State:"), output.StatusColor(status).Render(status))
			fmt.Fprintf(f.IOOut, "%s %s\n", output.Bold.Render("Branch:"), output.Cyan.Render(p.Target.RefName))

			if p.Target.Commit.Hash != "" {
				hash := p.Target.Commit.Hash
				if len(hash) > 12 {
					hash = hash[:12]
				}
				fmt.Fprintf(f.IOOut, "%s %s\n", output.Bold.Render("Commit:"), hash)
			}

			fmt.Fprintf(f.IOOut, "%s %s\n", output.Bold.Render("Creator:"), p.Creator.DisplayName)
			fmt.Fprintf(f.IOOut, "%s %s\n", output.Bold.Render("Duration:"), p.Duration())

			fmt.Fprintln(f.IOOut)
			fmt.Fprintf(f.IOOut, "%s  %s\n", output.Muted.Render("Created:"), models.FormatTime(p.CreatedOn))
			if p.CompletedOn != "" {
				fmt.Fprintf(f.IOOut, "%s %s\n", output.Muted.Render("Completed:"), models.FormatTime(p.CompletedOn))
			}

			// Fetch and display steps.
			stepsPath := fmt.Sprintf("/2.0/repositories/%s/%s/pipelines/%s/steps/", workspace, repo, pipelineUUID)
			steps, err := api.Paginate[models.PipelineStep](client, stepsPath, 0)
			if err != nil {
				return err
			}

			if len(steps) > 0 {
				fmt.Fprintln(f.IOOut)
				fmt.Fprintln(f.IOOut, output.Bold.Render("Steps:"))

				headers := []string{"Name", "Status", "Duration"}
				rows := make([][]string, len(steps))
				for i, s := range steps {
					stepStatus := stepStatusText(s)
					duration := stepDuration(s)
					rows[i] = []string{
						s.Name,
						output.StatusColor(stepStatus).Render(stepStatus),
						duration,
					}
				}
				output.RenderTable(headers, rows)
			}

			return nil
		},
	}
}

func stepStatusText(s models.PipelineStep) string {
	if s.State == nil {
		return "UNKNOWN"
	}
	if s.State.Result != nil {
		return s.State.Result.Name
	}
	if s.State.Stage != nil {
		return s.State.Stage.Name
	}
	return s.State.Name
}

func stepDuration(s models.PipelineStep) string {
	if s.StartedOn == "" {
		return ""
	}
	start, err := time.Parse(time.RFC3339Nano, s.StartedOn)
	if err != nil {
		return ""
	}

	var end time.Time
	if s.CompletedOn != "" {
		end, err = time.Parse(time.RFC3339Nano, s.CompletedOn)
		if err != nil {
			return ""
		}
	} else {
		end = time.Now()
	}

	d := end.Sub(start)
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}
