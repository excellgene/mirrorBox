package config

import (
	"log"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Store struct {
	configPath string
}

func NewStore(configPath string) *Store {
	return &Store{
		configPath: configPath,
	}
}

func (s *Store) Load() (*Config, error) {
	data, err := os.ReadFile(s.configPath)

	if err != nil {
		if os.IsNotExist(err) {
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

	log.Printf("Config saved to %s", s.configPath)
	log.Printf("Config content: %s", string(data))


	return nil
}
