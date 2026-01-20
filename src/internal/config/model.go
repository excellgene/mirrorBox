package config

import "time"

type FolderToSync struct {
	SourcePath      string `json:"SourcePath"`
	DestinationPath string `json:"DestinationPath"`
	Enabled         bool   `json:"Enabled"`
}

// Config defines the application configuration structure.
type Config struct {
	CheckInterval time.Duration  `json:"check_interval"`
	StartAtBoot   bool           `json:"start_at_boot"`
	Folders       []FolderToSync `json:"folders"`
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() *Config {
	return &Config{
		CheckInterval: 5 * time.Minute,
		StartAtBoot:   false,
		Folders:       []FolderToSync{},
	}
}
