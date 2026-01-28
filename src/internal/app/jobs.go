package app

import (
	"fmt"

	"excellgene.com/mirrorBox/internal/config"
	syncpkg "excellgene.com/mirrorBox/internal/sync"
)

type JobFactory struct {
}

func NewJobFactory() *JobFactory {
	return &JobFactory{}
}

// CreateFromConfig creates sync jobs from configuration.
// Only creates jobs that are enabled.
func (f *JobFactory) CreateFromConfig(cfg *config.Config) ([]*syncpkg.Job, error) {
	var jobs []*syncpkg.Job

	for _, jobCfg := range cfg.Folders {
		if !jobCfg.Enabled {
			continue
		}

		name := "Sync " + jobCfg.SourcePath + " to " + jobCfg.DestinationPath
		job, err := f.createJob(jobCfg)
		if err != nil {
			return nil, fmt.Errorf("create job %s: %w", name, err)
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

// createJob creates a single sync job from config.
func (f *JobFactory) createJob(cfg config.FolderToSync) (*syncpkg.Job, error) {
	job := syncpkg.NewJob(
		"Sync " + cfg.SourcePath + " to " + cfg.DestinationPath,
		cfg.SourcePath,
		cfg.DestinationPath,
	)

	return job, nil
}
