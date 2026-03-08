package config

import (
	"context"
	"fmt"

	"github.com/tyrantkhan/bb/internal/cmdutil"
	internalConfig "github.com/tyrantkhan/bb/internal/config"
	"github.com/urfave/cli/v3"
)

func newCmdGet() *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "Print the value of a given configuration key",
		ArgsUsage: "<key>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.NArg() != 1 {
				return fmt.Errorf("expected exactly 1 argument, got %d", cmd.NArg())
			}

			f := cmdutil.GetFactory(ctx)
			key := cmd.Args().First()

			val, err := getConfigValue(f.Config, key)
			if err != nil {
				return err
			}

			fmt.Fprintln(f.IOOut, val)
			return nil
		},
	}
}

// validKeys lists the known configuration keys.
var validKeys = []string{"default_workspace", "default_format", "editor"}

func getConfigValue(cfg *internalConfig.Config, key string) (string, error) {
	switch key {
	case "default_workspace":
		return cfg.DefaultWorkspace, nil
	case "default_format":
		return cfg.DefaultFormat, nil
	case "editor":
		return cfg.Editor, nil
	default:
		return "", fmt.Errorf("unknown configuration key %q\nValid keys: %v", key, validKeys)
	}
}
