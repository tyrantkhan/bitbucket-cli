package auth

import (
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/urfave/cli/v3"
)

// NewCmdAuth returns the parent auth command with its subcommands.
func NewCmdAuth() *cli.Command {
	return &cli.Command{
		Name:            "auth",
		Usage:           "Manage authentication",
		CommandNotFound: cmdutil.CommandNotFound,
		Commands: []*cli.Command{
			newCmdLogin(),
			newCmdLogout(),
			newCmdStatus(),
		},
	}
}
