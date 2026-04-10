package workspace

import (
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/urfave/cli/v3"
)

// NewCmdWorkspace returns the parent workspace command with its subcommands.
func NewCmdWorkspace() *cli.Command {
	return &cli.Command{
		Name:            "workspace",
		Usage:           "Manage workspaces",
		CommandNotFound: cmdutil.CommandNotFound,
		Commands: []*cli.Command{
			newCmdMembers(),
		},
	}
}
