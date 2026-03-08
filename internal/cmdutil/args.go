package cmdutil

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

// NoArgs wraps an action to reject unknown positional arguments on leaf commands.
func NoArgs(action cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Present() {
			return fmt.Errorf("unknown command %q for %q", cmd.Args().First(), cmd.FullName())
		}
		return action(ctx, cmd)
	}
}

// CommandNotFound is a shared handler for unknown subcommands. It prints an
// error message, suggests a similar command, and lists available commands.
func CommandNotFound(_ context.Context, cmd *cli.Command, command string) {
	fmt.Fprintf(os.Stderr, "unknown command %q for %q\n", command, cmd.FullName())

	var visible []*cli.Command
	for _, c := range cmd.Commands {
		if !c.Hidden {
			visible = append(visible, c)
		}
	}

	if suggestion := cli.SuggestCommand(visible, command); suggestion != "" {
		fmt.Fprintf(os.Stderr, "\nDid you mean this?\n\t%s\n", suggestion)
	}

	fmt.Fprintf(os.Stderr, "\nUsage:  %s <command> [flags]\n", cmd.FullName())
	fmt.Fprintf(os.Stderr, "\nAvailable commands:\n")
	for _, c := range visible {
		fmt.Fprintf(os.Stderr, "  %s\n", c.Name)
	}

	cli.OsExiter(1)
}
