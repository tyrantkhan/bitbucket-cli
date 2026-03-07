package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdStatus() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Show authentication status",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			f := cmdutil.GetFactory(ctx)

			client, err := f.APIClient()
			if err != nil {
				fmt.Fprintln(f.IOErr, output.Error.Render("Not logged in."))
				fmt.Fprintln(f.IOErr, output.Muted.Render("Run 'bb auth login' to authenticate."))
				return err
			}

			resp, err := client.Get("/2.0/user")
			if err != nil {
				fmt.Fprintln(f.IOErr, output.Error.Render("Authentication failed: "+err.Error()))
				return err
			}

			var user models.User
			if err := api.DecodeJSON(resp, &user); err != nil {
				return fmt.Errorf("failed to decode user response: %w", err)
			}

			// Retrieve credentials to show masked app password.
			creds, _ := f.AuthStore.GetCredentials()

			fmt.Fprintln(f.IOOut, output.Bold.Render("Authentication Status"))
			fmt.Fprintln(f.IOOut)
			fmt.Fprintf(f.IOOut, "  %s  %s\n", output.Muted.Render("User:"), user.DisplayName)
			fmt.Fprintf(f.IOOut, "  %s  %s\n", output.Muted.Render("Nickname:"), user.Nickname)

			if creds != nil {
				fmt.Fprintf(f.IOOut, "  %s  %s\n", output.Muted.Render("App Password:"), maskPassword(creds.AppPassword))
			}

			workspace := f.Config.DefaultWorkspace
			if workspace != "" {
				fmt.Fprintf(f.IOOut, "  %s  %s\n", output.Muted.Render("Workspace:"), workspace)
			} else {
				fmt.Fprintf(f.IOOut, "  %s  %s\n", output.Muted.Render("Workspace:"), output.Warning.Render("not set"))
			}

			fmt.Fprintln(f.IOOut)
			fmt.Fprintln(f.IOOut, output.Success.Render("Authenticated"))

			return nil
		},
	}
}

// maskPassword masks all but the last 4 characters of a password.
func maskPassword(password string) string {
	if len(password) <= 4 {
		return strings.Repeat("*", len(password))
	}
	return strings.Repeat("*", len(password)-4) + password[len(password)-4:]
}
