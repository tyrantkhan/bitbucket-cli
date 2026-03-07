package cmdutil

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/auth"
	"github.com/tyrantkhan/bb/internal/config"
)

type contextKey string

const factoryKey contextKey = "factory"

// Factory holds shared dependencies for commands.
type Factory struct {
	Config    *config.Config
	AuthStore auth.Store
	IOOut     io.Writer
	IOErr     io.Writer
}

// NewFactory creates a new Factory with default values.
func NewFactory() (*Factory, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	return &Factory{
		Config:    cfg,
		AuthStore: auth.NewFileStore(),
		IOOut:     os.Stdout,
		IOErr:     os.Stderr,
	}, nil
}

// APIClient creates an authenticated API client.
// For OAuth credentials, it auto-refreshes expired tokens before creating the client.
func (f *Factory) APIClient() (*api.Client, error) {
	creds, err := f.AuthStore.GetCredentials()
	if err != nil {
		return nil, err
	}

	if creds.IsOAuth() && creds.ExpiresAt > 0 && time.Now().Unix() >= creds.ExpiresAt {
		if creds.RefreshToken == "" || creds.ClientID == "" || creds.ClientSecret == "" {
			return nil, fmt.Errorf("OAuth token expired and cannot be refreshed. Run 'bb auth login --web' to re-authenticate")
		}

		tokenResp, err := auth.RefreshAccessToken(creds.ClientID, creds.ClientSecret, creds.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh OAuth token: %w. Run 'bb auth login --web' to re-authenticate", err)
		}

		creds.AccessToken = tokenResp.AccessToken
		if tokenResp.RefreshToken != "" {
			creds.RefreshToken = tokenResp.RefreshToken
		}
		creds.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Unix()

		if err := f.AuthStore.SetCredentials(creds); err != nil {
			fmt.Fprintf(f.IOErr, "Warning: failed to save refreshed token: %v\n", err)
		}
	}

	return api.NewClient(creds), nil
}

// WithFactory stores the factory in the context.
func WithFactory(ctx context.Context, f *Factory) context.Context {
	return context.WithValue(ctx, factoryKey, f)
}

// GetFactory retrieves the factory from the context.
func GetFactory(ctx context.Context) *Factory {
	f, _ := ctx.Value(factoryKey).(*Factory)
	return f
}
