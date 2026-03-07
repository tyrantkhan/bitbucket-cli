package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration.
type Config struct {
	DefaultWorkspace string `yaml:"default_workspace,omitempty"`
	DefaultFormat    string `yaml:"default_format,omitempty"` // "table" or "json"
	Editor           string `yaml:"editor,omitempty"`
}

// Load reads the config from disk. Returns default config if file doesn't exist.
func Load() (*Config, error) {
	cfg := &Config{
		DefaultFormat: "table",
	}

	data, err := os.ReadFile(ConfigFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.DefaultFormat == "" {
		cfg.DefaultFormat = "table"
	}

	return cfg, nil
}

// Save writes the config to disk.
func (c *Config) Save() error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigFilePath(), data, 0644)
}
