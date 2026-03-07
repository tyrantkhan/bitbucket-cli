package pipeline

import (
	"context"
	"fmt"

	"charm.land/huh/v2"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdStop() *cli.Command {
	return &cli.Command{
		Name:      "stop",
		Usage:     "Stop a running pipeline",
		ArgsUsage: "<uuid>",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
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

			// Confirmation prompt.
			var confirmed bool
			err = huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title(fmt.Sprintf("Stop pipeline %s?", pipelineUUID)).
						Value(&confirmed),
				),
			).Run()
			if err != nil {
				return err
			}
			if !confirmed {
				fmt.Fprintln(f.IOOut, output.Warning.Render("Stop cancelled."))
				return nil
			}

			path := fmt.Sprintf("/2.0/repositories/%s/%s/pipelines/%s/stopPipeline", workspace, repo, pipelineUUID)

			_, err = client.Post(path, nil)
			if err != nil {
				return err
			}

			fmt.Fprintln(f.IOOut, output.Success.Render("Pipeline stopped successfully."))

			return nil
		},
	}
}
