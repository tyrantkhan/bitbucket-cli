package pr

import "github.com/urfave/cli/v3"

// NewCmdPR returns the parent pr command with its subcommands.
func NewCmdPR() *cli.Command {
	return &cli.Command{
		Name:  "pr",
		Usage: "Manage pull requests",
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
		},
	}
}
