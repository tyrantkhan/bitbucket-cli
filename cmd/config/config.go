package config

import "github.com/urfave/cli/v3"

// NewCmdConfig returns the parent config command with its subcommands.
func NewCmdConfig() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Manage configuration",
		Commands: []*cli.Command{
			newCmdGet(),
			newCmdSet(),
			newCmdList(),
			newCmdClearCache(),
		},
	}
}
