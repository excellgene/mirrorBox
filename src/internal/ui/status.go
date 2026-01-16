package ui

import (
	"fmt"
	"log"

	"excellgene.com/symbaSync/internal/app"
)

// StatusWindow displays current sync status and job history.
// Responsibility:
//   - Show running/completed jobs
//   - Display sync statistics
//   - Show errors and logs
//
type StatusWindow struct {
	state *app.State
}

// NewStatusWindow creates a new status window.
func NewStatusWindow(state *app.State) *StatusWindow {
	return &StatusWindow{
		state: state,
	}
}

// Show displays the status window.
// Placeholder: In real implementation, this would open a GUI window.
func (w *StatusWindow) Show() {
	log.Println("Status window opened")
	w.logCurrentStatus()
	// TODO: Implement actual UI
}

// Hide closes the status window.
func (w *StatusWindow) Hide() {
	log.Println("Status window closed")
	// TODO: Implement actual UI
}

// Update refreshes the status display with current job states.
func (w *StatusWindow) Update() {
	w.logCurrentStatus()
	// TODO: Update actual UI
}

// OnJobEvent is called when a job status changes.
// Updates the display to reflect new status.
func (w *StatusWindow) OnJobEvent(event app.JobEvent) {
	log.Printf("Job event: %s - %v", event.JobName, event.Status)
	if event.Result != nil {
		log.Printf("  Created: %d, Updated: %d, Deleted: %d",
			event.Result.FilesCreated,
			event.Result.FilesUpdated,
			event.Result.FilesDeleted)
	}
	if event.Error != nil {
		log.Printf("  Error: %v", event.Error)
	}
	// TODO: Update actual UI
}

// logCurrentStatus prints current job status to console.
// Placeholder for actual UI rendering.
func (w *StatusWindow) logCurrentStatus() {
	jobs := w.state.AllJobs()

	fmt.Println("=== Current Status ===")
	for _, job := range jobs {
		fmt.Printf("Job: %s\n", job.Name)
		fmt.Printf("  Status: %v\n", job.Status())
		fmt.Printf("  Last Run: %s\n", job.LastRun())

		if result := job.LastResult(); result != nil {
			fmt.Printf("  Last Result:\n")
			fmt.Printf("    Created: %d\n", result.FilesCreated)
			fmt.Printf("    Updated: %d\n", result.FilesUpdated)
			fmt.Printf("    Deleted: %d\n", result.FilesDeleted)
			fmt.Printf("    Bytes: %d\n", result.BytesCopied)
		}

		if err := job.LastError(); err != nil {
			fmt.Printf("  Error: %v\n", err)
		}
		fmt.Println()
	}
	fmt.Println("=====================")
}
