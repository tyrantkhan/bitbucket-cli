package pipeline

import (
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/urfave/cli/v3"
)

// NewCmdPipeline returns the parent pipeline command with its subcommands.
func NewCmdPipeline() *cli.Command {
	return &cli.Command{
		Name:            "pipeline",
		Usage:           "Manage pipelines",
		CommandNotFound: cmdutil.CommandNotFound,
		Commands: []*cli.Command{
			newCmdList(),
			newCmdView(),
			newCmdRun(),
			newCmdStop(),
			newCmdLogs(),
		},
	}
}
