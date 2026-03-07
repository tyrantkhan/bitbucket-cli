package pipeline

import "github.com/urfave/cli/v3"

// NewCmdPipeline returns the parent pipeline command with its subcommands.
func NewCmdPipeline() *cli.Command {
	return &cli.Command{
		Name:  "pipeline",
		Usage: "Manage pipelines",
		Commands: []*cli.Command{
			newCmdList(),
			newCmdView(),
			newCmdRun(),
			newCmdStop(),
			newCmdLogs(),
		},
	}
}
