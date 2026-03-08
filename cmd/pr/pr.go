package pr

import (
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/urfave/cli/v3"
)

// NewCmdPR returns the parent pr command with its subcommands.
func NewCmdPR() *cli.Command {
	return &cli.Command{
		Name:            "pr",
		Usage:           "Manage pull requests",
		CommandNotFound: cmdutil.CommandNotFound,
		Commands: []*cli.Command{
			newCmdList(),
			newCmdView(),
			newCmdCreate(),
			newCmdMerge(),
			newCmdApprove(),
			newCmdDecline(),
			newCmdComment(),
			newCmdDiff(),
			newCmdActivity(),
			newCmdStatus(),
			newCmdReady(),
			newCmdDraft(),
			newCmdEdit(),
		},
	}
}
