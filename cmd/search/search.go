package search

import (
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/urfave/cli/v3"
)

// NewCmdSearch returns the parent search command with its subcommands.
func NewCmdSearch() *cli.Command {
	return &cli.Command{
		Name:            "search",
		Usage:           "Search across repositories",
		CommandNotFound: cmdutil.CommandNotFound,
		Commands: []*cli.Command{
			newCmdCode(),
		},
	}
}
