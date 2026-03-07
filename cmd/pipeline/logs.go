package pipeline

import (
	"context"
	"fmt"
	"io"
	"time"

	"charm.land/huh/v2"
	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/urfave/cli/v3"
)

func newCmdLogs() *cli.Command {
	return &cli.Command{
		Name:      "logs",
		Usage:     "View pipeline step logs",
		ArgsUsage: "<uuid>",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
			&cli.StringFlag{
				Name:  "step",
				Usage: "Step UUID (interactive picker if not provided)",
			},
			&cli.BoolFlag{
				Name:  "follow",
				Usage: "Follow log output for in-progress steps",
			},
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

			// Fetch steps.
			stepsPath := fmt.Sprintf("/2.0/repositories/%s/%s/pipelines/%s/steps/", workspace, repo, pipelineUUID)
			steps, err := api.Paginate[models.PipelineStep](client, stepsPath, 0)
			if err != nil {
				return err
			}

			if len(steps) == 0 {
				return fmt.Errorf("no steps found for pipeline %s", pipelineUUID)
			}

			stepUUID := cmd.String("step")

			// Interactive step picker if --step not provided.
			if stepUUID == "" {
				var options []huh.Option[string]
				for _, s := range steps {
					label := fmt.Sprintf("%s (%s)", s.Name, stepStatusText(s))
					options = append(options, huh.NewOption(label, s.UUID))
				}

				err = huh.NewSelect[string]().
					Title("Select a step").
					Options(options...).
					Value(&stepUUID).
					Run()
				if err != nil {
					return err
				}
			}

			if stepUUID == "" {
				return fmt.Errorf("step UUID is required")
			}

			follow := cmd.Bool("follow")

			logPath := fmt.Sprintf("/2.0/repositories/%s/%s/pipelines/%s/steps/%s/log", workspace, repo, pipelineUUID, stepUUID)

			if !follow {
				// Single fetch of the log.
				resp, err := client.Get(logPath)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return fmt.Errorf("failed to read log: %w", err)
				}

				fmt.Fprint(f.IOOut, string(body))
				return nil
			}

			// Follow mode: poll for new log content.
			var printed int

			for {
				resp, err := client.Get(logPath)
				if err != nil {
					return err
				}

				body, err := io.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil {
					return fmt.Errorf("failed to read log: %w", err)
				}

				if len(body) > printed {
					fmt.Fprint(f.IOOut, string(body[printed:]))
					printed = len(body)
				}

				// Check if the step is still in progress.
				stepPath := fmt.Sprintf("/2.0/repositories/%s/%s/pipelines/%s/steps/%s", workspace, repo, pipelineUUID, stepUUID)
				stepResp, err := client.Get(stepPath)
				if err != nil {
					return err
				}

				var step models.PipelineStep
				if err := api.DecodeJSON(stepResp, &step); err != nil {
					return fmt.Errorf("failed to decode step: %w", err)
				}

				status := stepStatusText(step)
				if status != "IN_PROGRESS" && status != "PENDING" && status != "RUNNING" {
					// Step is no longer in progress; fetch final log content.
					resp, err := client.Get(logPath)
					if err != nil {
						return err
					}

					body, err := io.ReadAll(resp.Body)
					resp.Body.Close()
					if err != nil {
						return fmt.Errorf("failed to read log: %w", err)
					}

					if len(body) > printed {
						fmt.Fprint(f.IOOut, string(body[printed:]))
					}
					break
				}

				time.Sleep(5 * time.Second)
			}

			return nil
		},
	}
}
