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

			for _, key := range validKeys {
				val, err := getConfigValue(f.Config, key)
				if err != nil {
					return err
				}
				fmt.Fprintf(f.IOOut, "%s=%s\n", key, val)
			}

			return nil
		}),
	}
}
