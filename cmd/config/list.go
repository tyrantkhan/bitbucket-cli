package config

import (
	"context"
	"fmt"

	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/urfave/cli/v3"
)

func newCmdList() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "Print a list of configuration keys and values",
		Action: cmdutil.NoArgs(func(ctx context.Context, cmd *cli.Command) error {
			f := cmdutil.GetFactory(ctx)

			fmt.Fprintf(f.IOOut, "default_workspace=%s\n", f.Config.DefaultWorkspace)
			fmt.Fprintf(f.IOOut, "default_format=%s\n", f.Config.DefaultFormat)
			fmt.Fprintf(f.IOOut, "editor=%s\n", f.Config.Editor)

			return nil
		}),
	}
}
