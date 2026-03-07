package auth

import (
	"context"
	"fmt"
	"time"

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
				if creds.IsOAuth() {
					fmt.Fprintf(f.IOOut, "  %s  %s\n", output.Muted.Render("Auth method:"), "OAuth")
					fmt.Fprintf(f.IOOut, "  %s  %s\n", output.Muted.Render("Access Token:"), maskSecret(creds.AccessToken))
					if creds.ExpiresAt > 0 {
						expiresAt := time.Unix(creds.ExpiresAt, 0)
						if time.Now().Before(expiresAt) {
							fmt.Fprintf(f.IOOut, "  %s  %s\n", output.Muted.Render("Token expires:"), expiresAt.Local().Format("2006-01-02 15:04"))
						} else {
							fmt.Fprintf(f.IOOut, "  %s  %s\n", output.Muted.Render("Token expires:"), output.Warning.Render("expired (will auto-refresh)"))
						}
					}
				} else {
					fmt.Fprintf(f.IOOut, "  %s  %s\n", output.Muted.Render("Auth method:"), "API token")
					fmt.Fprintf(f.IOOut, "  %s  %s\n", output.Muted.Render("API Token:"), maskSecret(creds.APIToken))
				}
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

// maskSecret masks a secret, showing only the last 4 characters.
func maskSecret(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return "****" + s[len(s)-4:]
}
