package app

import (
	"fmt"

	"excellgene.com/symbaSync/internal/config"
	"excellgene.com/symbaSync/internal/infra/smb"
	syncpkg "excellgene.com/symbaSync/internal/sync"
)

type JobFactory struct {
	smbClientFactory func(cfg smb.Config) smb.Client
}

func NewJobFactory(smbClientFactory func(cfg smb.Config) smb.Client) *JobFactory {
	return &JobFactory{
		smbClientFactory: smbClientFactory,
	}
}

// CreateFromConfig creates sync jobs from configuration.
// Only creates jobs that are enabled.
func (f *JobFactory) CreateFromConfig(cfg *config.Config) ([]*syncpkg.Job, error) {
	var jobs []*syncpkg.Job

	for _, jobCfg := range cfg.SyncJobs {
		if !jobCfg.Enabled {
			continue
		}

		job, err := f.createJob(jobCfg)
		if err != nil {
			return nil, fmt.Errorf("create job %s: %w", jobCfg.Name, err)
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

// createJob creates a single sync job from config.
func (f *JobFactory) createJob(cfg config.SyncJobConfig) (*syncpkg.Job, error) {
	// Parse SMB connection info from destination path
	// In real implementation, parse SMB URL (e.g., smb://server/share/path)
	// For now, use a mock config
	smbCfg := smb.Config{
		Host:     "localhost",
		Port:     445,
		Share:    "share",
		Username: "user",
		Password: "pass",
	}

	// Create SMB client
	smbClient := f.smbClientFactory(smbCfg)

	// Create sync job
	job := syncpkg.NewJob(
		cfg.Name,
		cfg.SourcePath,
		cfg.DestinationPath,
		smbClient,
	)

	return job, nil
}
