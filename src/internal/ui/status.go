package ui

import (
	"fmt"
	"log"
	"strings"

	"excellgene.com/mirrorBox/internal/app"
	syncpkg "excellgene.com/mirrorBox/internal/sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// StatusWindow displays current sync status and job history.
type StatusWindow struct {
	app    fyne.App
	window fyne.Window

	state   *app.State
	content *fyne.Container
}

// NewStatusWindow creates a new status window.
func NewStatusWindow(app fyne.App, state *app.State) *StatusWindow {
	return &StatusWindow{
		app:   app,
		state: state,
	}
}

// Show displays the status window.
// The window is created lazily and reused.
func (w *StatusWindow) Show() {
	if w.window == nil {
		w.window = w.app.NewWindow("SambaSync - Status")
		w.window.Resize(fyne.NewSize(600, 500))

		w.content = container.NewVBox()
		scroll := container.NewVScroll(w.content)

		w.window.SetContent(scroll)

		w.window.SetOnClosed(func() {
			w.window = nil
			log.Println("Status window closed")
		})
	}

	w.refreshUI()
	w.window.Show()
	w.window.RequestFocus()

	log.Println("Status window opened")
}

// Hide closes the status window.
func (w *StatusWindow) Hide() {
	if w.window != nil {
		w.window.Hide()
	}
}

// Update refreshes the status display manually.
func (w *StatusWindow) Update() {
	w.refreshUI()
}

// OnJobEvent is called when a job status changes.
// This updates the UI safely.
func (w *StatusWindow) OnJobEvent(event app.JobEvent) {
	log.Printf("Job event: %s - %v", event.JobName, event.Status)
	w.refreshUI()
}

// refreshUI rebuilds the UI from current job state.
func (w *StatusWindow) refreshUI() {
	if w.content == nil {
		return
	}

	w.content.Objects = nil

	jobs := w.state.AllJobs()
	if len(jobs) == 0 {
		w.content.Add(widget.NewLabel("No jobs configured."))
		w.content.Refresh()
		return
	}

	for _, job := range jobs {
		w.content.Add(w.renderJob(job))
		w.content.Add(widget.NewSeparator())
	}

	w.content.Refresh()
}

// renderJob creates a UI block for a single job.
func (w *StatusWindow) renderJob(job *syncpkg.Job) fyne.CanvasObject {
	lines := []string{
		fmt.Sprintf("Job: %s", job.Name),
		fmt.Sprintf("Status: %v", job.Status()),
		fmt.Sprintf("Last Run: %v", job.LastRun()),
	}

	if result := job.LastResult(); result != nil {
		lines = append(lines,
			"Last Result:",
			fmt.Sprintf("  Created: %d", result.FilesCreated),
			fmt.Sprintf("  Updated: %d", result.FilesUpdated),
			fmt.Sprintf("  Deleted: %d", result.FilesDeleted),
			fmt.Sprintf("  Bytes Copied: %d", result.BytesCopied),
		)
	}

	if err := job.LastError(); err != nil {
		lines = append(lines, fmt.Sprintf("Error: %v", err))
	}

	text := strings.Join(lines, "\n")
	return widget.NewLabel(text)
}