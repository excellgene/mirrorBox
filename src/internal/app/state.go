package app

import (
	"sync"

	syncpkg "excellgene.com/mirrorBox/internal/sync"
)

// State holds the runtime state of the application.
type State struct {
	mu   sync.RWMutex
	jobs map[string]*syncpkg.Job 
}

func NewState() *State {
	return &State{
		jobs: make(map[string]*syncpkg.Job),
	}
}

func (s *State) AddJob(job *syncpkg.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.Name] = job
}

func (s *State) GetJob(name string) *syncpkg.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.jobs[name]
}

func (s *State) AllJobs() []*syncpkg.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make([]*syncpkg.Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

func (s *State) RemoveJob(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, name)
}
