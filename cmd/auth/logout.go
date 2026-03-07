package auth

import (
	"context"
	"fmt"

	"charm.land/huh/v2"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdLogout() *cli.Command {
	return &cli.Command{
		Name:  "logout",
		Usage: "Remove stored authentication credentials",
		Action: cmdutil.NoArgs(func(ctx context.Context, cmd *cli.Command) error {
			f := cmdutil.GetFactory(ctx)

			var confirmed bool
			err := huh.NewConfirm().
				Title("Are you sure you want to log out?").
				Description("This will remove your stored credentials.").
				Value(&confirmed).
				Run()
			if err != nil {
				return err
			}

			if !confirmed {
				fmt.Fprintln(f.IOOut, output.Muted.Render("Logout cancelled."))
				return nil
			}

			if err := f.AuthStore.DeleteCredentials(); err != nil {
				return fmt.Errorf("failed to delete credentials: %w", err)
			}

			fmt.Fprintln(f.IOOut, output.Success.Render("Logged out successfully."))
			return nil
		}),
	}
}
