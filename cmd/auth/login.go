package auth

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"charm.land/huh/v2"
	"charm.land/huh/v2/spinner"
	"charm.land/lipgloss/v2"
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
			&cli.BoolFlag{
				Name:  "web",
				Usage: "Authenticate via OAuth browser flow",
			},
			&cli.BoolFlag{
				Name:  "api-token",
				Usage: "Authenticate via API token",
			},
			&cli.StringFlag{
				Name:  "username",
				Usage: "Atlassian account email (for API token auth)",
			},
			&cli.StringFlag{
				Name:  "token",
				Usage: "Bitbucket API token (for API token auth)",
			},
			&cli.StringFlag{
				Name:  "client-id",
				Usage: "OAuth consumer key (overrides default)",
			},
			&cli.StringFlag{
				Name:  "client-secret",
				Usage: "OAuth consumer secret (overrides default)",
			},
		},
		Action: cmdutil.NoArgs(func(ctx context.Context, cmd *cli.Command) error {
			f := cmdutil.GetFactory(ctx)

			useWeb := cmd.Bool("web")
			useAPIToken := cmd.Bool("api-token")

			// If neither flag is set, show interactive picker.
			if !useWeb && !useAPIToken {
				box := lipgloss.NewStyle().
					Width(50).
					BorderStyle(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("#0052CC")).
					Padding(1, 2).
					Render("bb has no backend or servers — all requests go directly to the Bitbucket API and your credentials never leave your device.")
				fmt.Fprintln(f.IOOut, box)

				var method string
				err := huh.NewForm(
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("How would you like to authenticate?").
							Options(
								huh.NewOption("Web browser (OAuth) — recommended", "web"),
								huh.NewOption("API token", "api_token"),
							).
							Value(&method),
					),
				).Run()
				if err != nil {
					return err
				}

				useWeb = method == "web"
			}

			var creds *auth.Credentials
			var err error

			if useWeb {
				creds, err = loginOAuth(f, cmd)
			} else {
				creds, err = loginAPIToken(f, cmd)
			}
			if err != nil {
				return err
			}

			// Validate credentials by fetching the authenticated user.
			client := api.NewClient(creds)
			var user models.User
			var validationErr error

			if err := spinner.New().
				Title("Verifying credentials...").
				Action(func() {
					resp, reqErr := client.Get("/2.0/user")
					if reqErr != nil {
						validationErr = reqErr
						return
					}
					if decErr := api.DecodeJSON(resp, &user); decErr != nil {
						validationErr = fmt.Errorf("failed to decode user response: %w", decErr)
					}
				}).
				Run(); err != nil {
				return err
			}

			if validationErr != nil {
				return fmt.Errorf("authentication failed: %w", validationErr)
			}

			// For OAuth, set the username from the API response.
			if creds.IsOAuth() {
				creds.Username = user.Nickname
			}

			fmt.Fprintln(f.IOOut, output.Success.Render(
				fmt.Sprintf("Authenticated as %s (%s)", user.DisplayName, user.Nickname),
			))

			// Fetch workspaces so the user can select a default.
			var workspaces []models.Workspace
			var fetchErr error

			if err := spinner.New().
				Title("Fetching workspaces...").
				Action(func() {
					workspaces, fetchErr = api.Paginate[models.Workspace](client, "/2.0/workspaces", 100)
				}).
				Run(); err != nil {
				return err
			}

			if fetchErr != nil {
				fmt.Fprintln(f.IOErr, output.Warning.Render("Could not fetch workspaces: "+fetchErr.Error()))
			}

			var defaultWorkspace string

			if len(workspaces) == 1 {
				defaultWorkspace = workspaces[0].Slug
				fmt.Fprintln(f.IOOut, output.Muted.Render(
					fmt.Sprintf("Default workspace set to %s (only workspace available)", defaultWorkspace),
				))
			} else if len(workspaces) > 1 {
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
		}),
	}
}

func loginOAuth(f *cmdutil.Factory, cmd *cli.Command) (*auth.Credentials, error) {
	clientID := cmd.String("client-id")
	clientSecret := cmd.String("client-secret")

	// Fall back to environment variables.
	if clientID == "" {
		clientID = os.Getenv("BB_CLIENT_ID")
	}
	if clientSecret == "" {
		clientSecret = os.Getenv("BB_CLIENT_SECRET")
	}

	// Fall back to built-in defaults.
	if clientID == "" {
		clientID = auth.DefaultClientID
	}
	if clientSecret == "" {
		clientSecret = auth.DefaultClientSecret
	}

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("OAuth client credentials not configured.\n" +
			"Set BB_CLIENT_ID and BB_CLIENT_SECRET environment variables,\n" +
			"or use --client-id and --client-secret flags.\n\n" +
			"Create an OAuth consumer at: https://bitbucket.org/account/settings/ → OAuth consumers")
	}

	// Start callback server on a random port.
	port, codeCh, errCh, err := auth.StartCallbackServer()
	if err != nil {
		return nil, err
	}

	// Build authorization URL and open browser.
	authURL := auth.AuthorizeURL(clientID, port)
	fmt.Fprintln(f.IOOut, output.Muted.Render("Opening browser to authorize..."))
	fmt.Fprintln(f.IOOut, output.Muted.Render("If the browser doesn't open, visit:"))
	fmt.Fprintln(f.IOOut, authURL)

	if err := openBrowser(authURL); err != nil {
		fmt.Fprintln(f.IOErr, output.Warning.Render("Could not open browser automatically."))
	}

	fmt.Fprintln(f.IOOut, output.Muted.Render("Waiting for authorization..."))

	// Wait for the callback.
	var code string
	select {
	case code = <-codeCh:
	case err := <-errCh:
		return nil, err
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("authentication timed out after 5 minutes")
	}

	// Exchange code for tokens.
	var tokenResp *auth.TokenResponse
	var exchangeErr error

	if err := spinner.New().
		Title("Exchanging authorization code for tokens...").
		Action(func() {
			tokenResp, exchangeErr = auth.ExchangeCode(clientID, clientSecret, code, port)
		}).
		Run(); err != nil {
		return nil, err
	}

	if exchangeErr != nil {
		return nil, exchangeErr
	}

	return &auth.Credentials{
		AuthMethod:   "oauth",
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Unix(),
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}, nil
}

func loginAPIToken(f *cmdutil.Factory, cmd *cli.Command) (*auth.Credentials, error) {
	email := cmd.String("username")
	apiToken := cmd.String("token")

	// Interactive mode: prompt for credentials if not provided via flags.
	if email == "" || apiToken == "" {
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
					Title("API Token Authentication").
					Description(
						"Generate a token at https://id.atlassian.com/manage-profile/security/api-tokens\n\n"+
							"Required scopes: Repositories, Pull Requests, Pipelines (read/write)\n\n"+
							"All requests go directly to the Bitbucket API.\n"+
							"Your credentials are stored locally and never leave your device.",
					),
				huh.NewInput().
					Title("Email").
					Description("Your Atlassian account email").
					Value(&email),
				huh.NewInput().
					Title("API Token").
					Description("Paste your Atlassian API token").
					EchoMode(huh.EchoModePassword).
					Value(&apiToken),
			),
		).Run()
		if err != nil {
			return nil, err
		}
	}

	if email == "" || apiToken == "" {
		return nil, fmt.Errorf("email and API token are required")
	}

	return &auth.Credentials{
		AuthMethod: "api_token",
		Username:   email,
		APIToken:   apiToken,
	}, nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Start()
}
