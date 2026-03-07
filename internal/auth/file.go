package auth

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/tyrantkhan/bb/internal/config"
)

// FileStore implements Store using a JSON file on disk.
type FileStore struct{}

// NewFileStore creates a new file-based credential store.
func NewFileStore() *FileStore {
	return &FileStore{}
}

// GetCredentials reads credentials from the credentials file.
func (f *FileStore) GetCredentials() (*Credentials, error) {
	path := config.CredentialsFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("not logged in. Run 'bb auth login' to authenticate")
		}
		return nil, err
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}

	if creds.Username == "" || creds.AppPassword == "" {
		return nil, errors.New("invalid credentials. Run 'bb auth login' to re-authenticate")
	}

	return &creds, nil
}

// SetCredentials writes credentials to the credentials file with 0600 permissions.
func (f *FileStore) SetCredentials(creds *Credentials) error {
	if err := config.EnsureConfigDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(config.CredentialsFilePath(), data, 0600)
}

// DeleteCredentials removes the credentials file.
func (f *FileStore) DeleteCredentials() error {
	err := os.Remove(config.CredentialsFilePath())
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
