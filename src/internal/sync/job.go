package sync

import (
	"context"
	"fmt"
	"time"

	"excellgene.com/symbaSync/internal/infra/fs"
	"excellgene.com/symbaSync/internal/infra/smb"
)

// JobStatus represents the current state of a sync job.
type JobStatus int

const (
	StatusIdle    JobStatus = iota // Job not running
	StatusRunning                  // Job currently executing
	StatusSuccess                  // Last run completed successfully
	StatusError                    // Last run failed
)

// Job represents a complete sync job.
// Responsibility: Orchestrate the full sync workflow (walk, diff, sync).
// This is the main entry point for running a sync operation.
type Job struct {
	Name            string
	SourcePath      string
	DestinationPath string

	// Dependencies (injected)
	sourceWalker fs.Walker
	destWalker   fs.Walker
	differ       *Differ
	syncer       *Syncer

	// State
	status     JobStatus
	lastRun    time.Time
	lastResult *SyncResult
	lastError  error
}

// NewJob creates a new sync job.
// sourceWalker walks the local filesystem.
// smbClient provides access to remote filesystem.
func NewJob(name, sourcePath, destPath string, smbClient smb.Client) *Job {
	return &Job{
		Name:            name,
		SourcePath:      sourcePath,
		DestinationPath: destPath,
		sourceWalker:    fs.NewLocalWalker(sourcePath),
		// destWalker would be created from smbClient.Walk() when running
		differ: NewDiffer(),
		syncer: NewSyncer(smbClient),
		status: StatusIdle,
	}
}

// Run executes the sync job.
// Workflow:
//  1. Walk source filesystem
//  2. Walk destination filesystem
//  3. Compute diff
//  4. Apply sync operations
//
// Returns SyncResult with statistics and any errors encountered.
func (j *Job) Run(ctx context.Context) (*SyncResult, error) {
	j.status = StatusRunning
	j.lastRun = time.Now()

	// Step 1: Walk source
	var sourceFiles []fs.FileInfo
	err := j.sourceWalker.Walk(func(info fs.FileInfo) error {
		sourceFiles = append(sourceFiles, info)
		return nil
	})
	if err != nil {
		j.status = StatusError
		j.lastError = fmt.Errorf("walk source: %w", err)
		return nil, j.lastError
	}

	// Step 2: Walk destination
	// In real implementation, get walker from SMB client
	// For now, use empty list (first run will copy everything)
	var destFiles []fs.FileInfo
	if j.destWalker != nil {
		err = j.destWalker.Walk(func(info fs.FileInfo) error {
			destFiles = append(destFiles, info)
			return nil
		})
		if err != nil {
			j.status = StatusError
			j.lastError = fmt.Errorf("walk destination: %w", err)
			return nil, j.lastError
		}
	}

	// Step 3: Compute diff
	diffResult := j.differ.Diff(sourceFiles, destFiles)

	// Step 4: Apply sync
	syncResult, err := j.syncer.Sync(ctx, diffResult, j.SourcePath, j.DestinationPath)
	if err != nil {
		j.status = StatusError
		j.lastError = err
		j.lastResult = syncResult
		return syncResult, err
	}

	// Update status
	if len(syncResult.Errors) > 0 {
		j.status = StatusError
		j.lastError = fmt.Errorf("%d errors during sync", len(syncResult.Errors))
	} else {
		j.status = StatusSuccess
		j.lastError = nil
	}
	j.lastResult = syncResult

	return syncResult, nil
}

// Status returns the current job status.
func (j *Job) Status() JobStatus {
	return j.status
}

// LastResult returns the result of the last run.
func (j *Job) LastResult() *SyncResult {
	return j.lastResult
}

// LastError returns the error from the last run, if any.
func (j *Job) LastError() error {
	return j.lastError
}

// LastRun returns when the job was last executed.
func (j *Job) LastRun() time.Time {
	return j.lastRun
}
