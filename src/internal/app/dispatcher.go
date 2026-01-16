package app

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	syncpkg "excellgene.com/symbaSync/internal/sync"
)

// JobEvent represents an event from job execution.
type JobEvent struct {
	JobName string
	Status  syncpkg.JobStatus
	Result  *syncpkg.SyncResult
	Error   error
}

// Dispatcher manages and executes sync jobs.
// Responsibility:
//   - Schedule and run jobs
//   - Manage goroutines and cancellation
//   - Emit events for UI updates
//
// Dispatcher is the coordination layer between jobs and the rest of the app.
type Dispatcher struct {
	state  *State
	events chan JobEvent

	// Cancellation
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewDispatcher creates a new job dispatcher.
func NewDispatcher(state *State) *Dispatcher {
	ctx, cancel := context.WithCancel(context.Background())

	return &Dispatcher{
		state:  state,
		events: make(chan JobEvent, 100),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Events returns a channel for receiving job events.
// UI can listen to this channel for updates.
func (d *Dispatcher) Events() <-chan JobEvent {
	return d.events
}

// RunNow executes a job immediately.
// Non-blocking: runs job in background goroutine.
func (d *Dispatcher) RunNow(jobName string) error {
	job := d.state.GetJob(jobName)
	if job == nil {
		return fmt.Errorf("job not found: %s", jobName)
	}

	// Run in background
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		d.runJob(job)
	}()

	return nil
}

// RunAll executes all registered jobs.
// Non-blocking: each job runs in its own goroutine.
func (d *Dispatcher) RunAll() {
	jobs := d.state.AllJobs()

	for _, job := range jobs {
		job := job // Capture for goroutine
		d.wg.Add(1)
		go func() {
			defer d.wg.Done()
			d.runJob(job)
		}()
	}
}

// StartScheduler starts the periodic job scheduler.
// Runs jobs at configured intervals until Stop() is called.
func (d *Dispatcher) StartScheduler(interval time.Duration) {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-d.ctx.Done():
				return
			case <-ticker.C:
				log.Println("Scheduler tick: running all jobs")
				d.RunAll()
			}
		}
	}()
}

// Stop gracefully shuts down the dispatcher.
// Waits for all running jobs to complete.
func (d *Dispatcher) Stop() {
	log.Println("Stopping dispatcher...")
	d.cancel()
	d.wg.Wait()
	close(d.events)
	log.Println("Dispatcher stopped")
}

// runJob executes a single job and emits events.
// This is the internal method that actually runs the job.
func (d *Dispatcher) runJob(job *syncpkg.Job) {
	log.Printf("Running job: %s", job.Name)

	// Create job-specific context with timeout
	ctx, cancel := context.WithTimeout(d.ctx, 30*time.Minute)
	defer cancel()

	// Run the job
	result, err := job.Run(ctx)

	// Emit event
	event := JobEvent{
		JobName: job.Name,
		Status:  job.Status(),
		Result:  result,
		Error:   err,
	}

	select {
	case d.events <- event:
	case <-d.ctx.Done():
		// Dispatcher is shutting down
	}

	if err != nil {
		log.Printf("Job %s failed: %v", job.Name, err)
	} else {
		log.Printf("Job %s completed: %d created, %d updated, %d deleted",
			job.Name, result.FilesCreated, result.FilesUpdated, result.FilesDeleted)
	}
}
