package config

import "time"

type Config struct {
	SyncJobs []SyncJobConfig `json:"sync_jobs"`
	CheckInterval time.Duration `json:"check_interval"`
	LogLevel string `json:"log_level"`
}

// SyncJobConfig defines a single sync job configuration.
type SyncJobConfig struct {
	Name string `json:"name"`
	SourcePath string `json:"source_path"`
	DestinationPath string `json:"destination_path"`
	Enabled bool `json:"enabled"`
	Schedule string `json:"schedule"`
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() *Config {
	return &Config{
		SyncJobs:      []SyncJobConfig{},
		CheckInterval: 5 * time.Minute,
		LogLevel:      "info",
	}
}
