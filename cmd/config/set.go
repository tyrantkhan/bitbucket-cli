package config

import (
	"context"
	"fmt"

	"github.com/tyrantkhan/bb/internal/cmdutil"
	internalConfig "github.com/tyrantkhan/bb/internal/config"
	"github.com/urfave/cli/v3"
)

func newCmdSet() *cli.Command {
	return &cli.Command{
		Name:      "set",
		Usage:     "Update configuration with a value for the given key",
		ArgsUsage: "<key> <value>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.NArg() != 2 {
				return fmt.Errorf("expected exactly 2 arguments, got %d", cmd.NArg())
			}

			f := cmdutil.GetFactory(ctx)
			key := cmd.Args().Get(0)
			value := cmd.Args().Get(1)

			if err := setConfigValue(f.Config, key, value); err != nil {
				return err
			}

			return f.Config.Save()
		},
	}
}

func setConfigValue(cfg *internalConfig.Config, key, value string) error {
	switch key {
	case "default_workspace":
		cfg.DefaultWorkspace = value
	case "default_format":
		if value != "table" && value != "json" {
			return fmt.Errorf("invalid value %q for default_format: must be \"table\" or \"json\"", value)
		}
		cfg.DefaultFormat = value
	case "editor":
		cfg.Editor = value
	default:
		return fmt.Errorf("unknown configuration key %q\nValid keys: %v", key, validKeys)
	}
	return nil
}
