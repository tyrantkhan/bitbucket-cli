package cmdutil

import (
	"context"
	"io"
	"os"

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
func (f *Factory) APIClient() (*api.Client, error) {
	creds, err := f.AuthStore.GetCredentials()
	if err != nil {
		return nil, err
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
