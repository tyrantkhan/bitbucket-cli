package config

import (
	"os"
	"path/filepath"
)

const appName = "bb"

// ConfigDir returns the configuration directory for bb.
// Uses XDG_CONFIG_HOME if set, otherwise ~/.config/bb/
func ConfigDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, appName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", appName)
}

// ConfigFilePath returns the full path to the config file.
func ConfigFilePath() string {
	return filepath.Join(ConfigDir(), "config.yml")
}

// CredentialsFilePath returns the full path to the credentials file.
func CredentialsFilePath() string {
	return filepath.Join(ConfigDir(), "credentials.json")
}

// EnsureConfigDir creates the config directory if it doesn't exist.
func EnsureConfigDir() error {
	return os.MkdirAll(ConfigDir(), 0700)
}
