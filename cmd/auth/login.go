package auth

import (
	"context"
	"fmt"

	"charm.land/huh/v2"
	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/auth"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdLogin() *cli.Command {
	return &cli.Command{
		Name:  "login",
		Usage: "Authenticate with Bitbucket Cloud",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "username",
				Usage: "Bitbucket username",
			},
			&cli.StringFlag{
				Name:  "app-password",
				Usage: "Bitbucket app password",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			f := cmdutil.GetFactory(ctx)

			username := cmd.String("username")
			appPassword := cmd.String("app-password")

			// Interactive mode: prompt for credentials if not provided via flags.
			if username == "" || appPassword == "" {
				err := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Username").
							Description("Your Bitbucket username").
							Value(&username),
						huh.NewInput().
							Title("App Password").
							Description("Generate one at Bitbucket > Personal settings > App passwords").
							EchoMode(huh.EchoModePassword).
							Value(&appPassword),
					),
				).Run()
				if err != nil {
					return err
				}
			}

			if username == "" || appPassword == "" {
				return fmt.Errorf("username and app password are required")
			}

			// Validate credentials by fetching the authenticated user.
			creds := &auth.Credentials{
				Username:    username,
				AppPassword: appPassword,
			}
			client := api.NewClient(creds)

			resp, err := client.Get("/2.0/user")
			if err != nil {
				fmt.Fprintln(f.IOErr, output.Error.Render("Authentication failed: "+err.Error()))
				return err
			}

			var user models.User
			if err := api.DecodeJSON(resp, &user); err != nil {
				return fmt.Errorf("failed to decode user response: %w", err)
			}

			fmt.Fprintln(f.IOOut, output.Success.Render(
				fmt.Sprintf("Authenticated as %s (%s)", user.DisplayName, user.Nickname),
			))

			// Fetch workspaces so the user can select a default.
			workspaces, err := api.Paginate[models.Workspace](client, "/2.0/workspaces", 100)
			if err != nil {
				fmt.Fprintln(f.IOErr, output.Warning.Render("Could not fetch workspaces: "+err.Error()))
			}

			var defaultWorkspace string

			if len(workspaces) == 1 {
				// Only one workspace, use it automatically.
				defaultWorkspace = workspaces[0].Slug
				fmt.Fprintln(f.IOOut, output.Muted.Render(
					fmt.Sprintf("Default workspace set to %s (only workspace available)", defaultWorkspace),
				))
			} else if len(workspaces) > 1 {
				// Build options for the select prompt.
				options := make([]huh.Option[string], len(workspaces))
				for i, ws := range workspaces {
					options[i] = huh.NewOption(fmt.Sprintf("%s (%s)", ws.Name, ws.Slug), ws.Slug)
				}

				err := huh.NewForm(
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("Select default workspace").
							Options(options...).
							Value(&defaultWorkspace),
					),
				).Run()
				if err != nil {
					return err
				}
			}

			// Store credentials.
			if err := f.AuthStore.SetCredentials(creds); err != nil {
				return fmt.Errorf("failed to save credentials: %w", err)
			}

			// Save default workspace to config.
			if defaultWorkspace != "" {
				f.Config.DefaultWorkspace = defaultWorkspace
				if err := f.Config.Save(); err != nil {
					fmt.Fprintln(f.IOErr, output.Warning.Render("Failed to save config: "+err.Error()))
				}
			}

			fmt.Fprintln(f.IOOut, output.Success.Render("Login complete!"))
			if defaultWorkspace != "" {
				fmt.Fprintln(f.IOOut, output.Muted.Render("Default workspace: "+output.Bold.Render(defaultWorkspace)))
			}

			return nil
		},
	}
}
