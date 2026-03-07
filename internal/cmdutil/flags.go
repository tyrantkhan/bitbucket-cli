package cmdutil

import (
	"github.com/urfave/cli/v3"
)

// WorkspaceFlag is the --workspace / -w flag.
var WorkspaceFlag = &cli.StringFlag{
	Name:    "workspace",
	Aliases: []string{"w"},
	Usage:   "Bitbucket workspace slug",
}

// RepoFlag is the --repo / -r flag.
var RepoFlag = &cli.StringFlag{
	Name:    "repo",
	Aliases: []string{"r"},
	Usage:   "Repository slug",
}

// FormatFlag is the --format flag.
var FormatFlag = &cli.StringFlag{
	Name:  "format",
	Usage: "Output format: table, json",
	Value: "table",
}

// LimitFlag is the --limit flag.
var LimitFlag = &cli.IntFlag{
	Name:  "limit",
	Usage: "Maximum number of results",
	Value: 30,
}

// WebFlag is the --web flag.
var WebFlag = &cli.BoolFlag{
	Name:  "web",
	Usage: "Open in browser",
}
