package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Store handles persistence of configuration.
// Responsibility: Load/save config from/to disk. No business logic.
type Store struct {
	configPath string
}

// NewStore creates a new config store.
// configPath should be an absolute path to the config file.
func NewStore(configPath string) *Store {
	return &Store{
		configPath: configPath,
	}
}

// Load reads configuration from disk.
// Returns default config if file doesn't exist.
func (s *Store) Load() (*Config, error) {
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}

// Save writes configuration to disk.
func (s *Store) Save(cfg *Config) error {
	// Ensure directory exists
	dir := filepath.Dir(s.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(s.configPath, data, 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}
