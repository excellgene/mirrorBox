package app

import (
	"sync"

	syncpkg "excellgene.com/symbaSync/internal/sync"
)

// State holds the runtime state of the application.
// Responsibility: Thread-safe access to application state (jobs, status, etc.)
// This is the single source of truth for runtime data.
type State struct {
	mu   sync.RWMutex
	jobs map[string]*syncpkg.Job // Job name -> Job
}

// NewState creates a new application state.
func NewState() *State {
	return &State{
		jobs: make(map[string]*syncpkg.Job),
	}
}

// AddJob registers a new job.
func (s *State) AddJob(job *syncpkg.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.Name] = job
}

// GetJob retrieves a job by name.
// Returns nil if job doesn't exist.
func (s *State) GetJob(name string) *syncpkg.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.jobs[name]
}

// AllJobs returns a copy of all registered jobs.
func (s *State) AllJobs() []*syncpkg.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make([]*syncpkg.Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// RemoveJob unregisters a job.
func (s *State) RemoveJob(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, name)
}
