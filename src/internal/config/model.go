package config

import "time"

// Config represents the application configuration.
// Responsibility: Define configuration schema only. No I/O, no validation beyond types.
type Config struct {
	// SyncJobs defines the configured sync jobs to run
	SyncJobs []SyncJobConfig `json:"sync_jobs"`

	// CheckInterval is how often to run scheduled syncs
	CheckInterval time.Duration `json:"check_interval"`

	// LogLevel controls logging verbosity
	LogLevel string `json:"log_level"`
}

// SyncJobConfig defines a single sync job configuration.
type SyncJobConfig struct {
	// Name is a human-readable identifier
	Name string `json:"name"`

	// SourcePath is the local directory to sync from
	SourcePath string `json:"source_path"`

	// DestinationPath is the remote/target path (e.g., SMB share)
	DestinationPath string `json:"destination_path"`

	// Enabled determines if this job should run
	Enabled bool `json:"enabled"`

	// Schedule defines when this job runs (e.g., "manual", "interval")
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
