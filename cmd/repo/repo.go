package repo

import (
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/urfave/cli/v3"
)

// NewCmdRepo returns the parent repo command with its subcommands.
func NewCmdRepo() *cli.Command {
	return &cli.Command{
		Name:            "repo",
		Usage:           "Manage repositories",
		CommandNotFound: cmdutil.CommandNotFound,
		Commands: []*cli.Command{
			newCmdList(),
			newCmdView(),
			newCmdClone(),
			newCmdCreate(),
			newCmdMove(),
		},
	}
}
