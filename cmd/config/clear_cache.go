package config

import (
	"context"
	"fmt"
	"os"

	"github.com/tyrantkhan/bb/internal/cmdutil"
	internalConfig "github.com/tyrantkhan/bb/internal/config"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdClearCache() *cli.Command {
	return &cli.Command{
		Name:  "clear-cache",
		Usage: "Delete the local cache used for update checks",
		Action: cmdutil.NoArgs(func(ctx context.Context, cmd *cli.Command) error {
			f := cmdutil.GetFactory(ctx)

			path := internalConfig.StateFilePath()
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove cache: %w", err)
			}

			fmt.Fprintln(f.IOErr, output.Success.Render("Cache cleared."))
			return nil
		}),
	}
}
