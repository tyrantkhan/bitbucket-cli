package cmdutil

import (
	"context"
	"fmt"

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
