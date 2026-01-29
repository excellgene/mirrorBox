package ui

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"excellgene.com/mirrorBox/internal/app"
	"excellgene.com/mirrorBox/internal/config"
	syncpkg "excellgene.com/mirrorBox/internal/sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// SettingsWindow manages the settings/configuration window.
type SettingsWindow struct {
	app          fyne.App
	window       fyne.Window
	statusWidget *widget.Label

	config       *config.Config
	store        *config.Store
	state        *app.State
	statusWindow *StatusWindow
	jobFactory   *app.JobFactory
}

// NewSettingsWindow creates a new settings window.
func NewSettingsWindow(
	app fyne.App,
	cfg *config.Config,
	store *config.Store,
	state *app.State,
	statusWindow *StatusWindow,
	jobFactory *app.JobFactory,
) *SettingsWindow {
	return &SettingsWindow{
		app:          app,
		config:       cfg,
		store:        store,
		state:        state,
		statusWindow: statusWindow,
		jobFactory:   jobFactory,
	}
}

// This is called automatically after config changes.
func (w *SettingsWindow) reloadJobs() error {
	// Create new jobs from current config
	newJobs, err := w.jobFactory.CreateFromConfig(w.config)
	if err != nil {
		log.Printf("Failed to create jobs from config: %v", err)
		return fmt.Errorf("create jobs: %w", err)
	}

	// Reload jobs in the state
	w.state.ReloadJobs(newJobs)

	log.Printf("Reloaded %d jobs from config", len(newJobs))
	for _, job := range newJobs {
		log.Printf("  - %s", job.Name)
	}

	return nil
}

// saveConfigAndReload saves the config and automatically reloads jobs.
func (w *SettingsWindow) saveConfigAndReload() error {
	if err := w.store.Save(w.config); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	if err := w.reloadJobs(); err != nil {
		return fmt.Errorf("reload jobs: %w", err)
	}

	return nil
}

func deleteFolder(cfg *config.Config, store *config.Store, index int, refreshFunc func(), reloadJobsFunc func() error) {
	if index < 0 || index >= len(cfg.Folders) {
		log.Printf("Invalid folder index: %d", index)
		return
	}

	cfg.Folders = append(cfg.Folders[:index], cfg.Folders[index+1:]...)

	err := store.Save(cfg)
	if err != nil {
		log.Printf("Failed to save config after deleting folder: %v", err)
		return
	}

	if reloadJobsFunc != nil {
		if err := reloadJobsFunc(); err != nil {
			log.Printf("Failed to reload jobs after config change: %v", err)
		}
	}

	// Refresh the UI
	if refreshFunc != nil {
		refreshFunc()
	}
}

func addOrEditFolder(cfg *config.Config, store *config.Store, folderIndex int, refreshFunc func(), reloadJobsFunc func() error) {
	isEdit := folderIndex >= 0 && folderIndex < len(cfg.Folders)

	var title string
	var folder config.FolderToSync

	if isEdit {
		title = "Edit Sync Folder"
		folder = cfg.Folders[folderIndex]
	} else {
		title = "Add Sync Folder"
		folder = config.FolderToSync{Enabled: true}
	}

	modal := fyne.CurrentApp().NewWindow(title)

	// Source entry
	sourceEntry := widget.NewEntry()
	sourceEntry.SetPlaceHolder("Source Path")
	if isEdit {
		sourceEntry.SetText(folder.SourcePath)
	}

	sourceBtn := widget.NewButton("Browse…", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				sourceEntry.SetText(uri.Path())
			}
		}, modal).Show()
	})

	// Destination entry
	destinationEntry := widget.NewEntry()
	destinationEntry.SetPlaceHolder("Destination Path")
	if isEdit {
		destinationEntry.SetText(folder.DestinationPath)
	}

	destinationBtn := widget.NewButton("Browse…", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				destinationEntry.SetText(uri.Path())
			}
		}, modal).Show()
	})

	enabledCheck := widget.NewCheck("Enabled", func(checked bool) {})
	enabledCheck.SetChecked(folder.Enabled)

	saveButton := widget.NewButton("Save", func() {
		// Validate paths
		if sourceEntry.Text == "" || destinationEntry.Text == "" {
			dialog.ShowError(
				fmt.Errorf("both source and destination paths are required"),
				modal,
			)
			return
		}

		folder.SourcePath = sourceEntry.Text
		folder.DestinationPath = destinationEntry.Text
		folder.Enabled = enabledCheck.Checked

		if isEdit {
			cfg.Folders[folderIndex] = folder
			log.Printf("Updated folder at index %d: %s -> %s", folderIndex, folder.SourcePath, folder.DestinationPath)
		} else {
			cfg.Folders = append(cfg.Folders, folder)
			log.Printf("Adding new folder to sync: %s -> %s", folder.SourcePath, folder.DestinationPath)
		}

		if err := store.Save(cfg); err != nil {
			log.Printf("Failed to save config: %v", err)
			dialog.ShowError(
				fmt.Errorf("failed to save configuration: %w", err),
				modal,
			)
			return
		}

		// Reload jobs to reflect config changes
		if reloadJobsFunc != nil {
			if err := reloadJobsFunc(); err != nil {
				log.Printf("Failed to reload jobs after config change: %v", err)
			}
		}

		modal.Close()

		// Refresh the folder list
		if refreshFunc != nil {
			refreshFunc()
		}
	})

	cancelButton := widget.NewButton("Cancel", func() {
		modal.Close()
	})

	form := container.NewVBox(
		widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),

		widget.NewLabel("Source Path"),
		container.NewBorder(nil, nil, nil, sourceBtn, sourceEntry),

		widget.NewLabel("Destination Path"),
		container.NewBorder(nil, nil, nil, destinationBtn, destinationEntry),

		enabledCheck,
		widget.NewSeparator(),
		container.NewGridWithColumns(2, cancelButton, saveButton),
	)

	modal.SetContent(form)
	modal.Resize(fyne.NewSize(700, 400))
	modal.Show()
}

// NewFolderWindow creates a new folder window.
func NewFolderWindow(app fyne.App, cfg *config.Config, store *config.Store, reloadJobsFunc func() error) fyne.Window {
	modal := fyne.CurrentApp().NewWindow("Syncing Folders")

	// Container for the folder list
	var folderContainer *fyne.Container

	// Declare refreshFolders function variable to allow it to reference itself
	var refreshFolders func()

	// Function to refresh/rebuild the folder list
	refreshFolders = func() {
		widgets := []fyne.CanvasObject{}

		// Title
		widgets = append(widgets,
			widget.NewLabelWithStyle("Configured Sync Folders", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
		)

		// List existing folders
		if len(cfg.Folders) == 0 {
			widgets = append(widgets,
				widget.NewLabel("No folders configured yet."),
				widget.NewSeparator(),
			)
		} else {
			for i := range cfg.Folders {
				index := i // Capture index for closure

				folder := cfg.Folders[index]

				// Status indicator
				statusText := "✓ Enabled"
				if !folder.Enabled {
					statusText = "✗ Disabled"
				}

				// Folder card
				folderLabel := widget.NewLabel(folder.SourcePath + " → " + folder.DestinationPath)
				statusLabel := widget.NewLabel(statusText)

				editBtn := widget.NewButton("Edit", func() {
					addOrEditFolder(cfg, store, index, refreshFolders, reloadJobsFunc)
				})

				deleteBtn := widget.NewButton("Delete", func() {
					// Confirm deletion
					confirmDialog := dialog.NewConfirm(
						"Delete Folder",
						"Are you sure you want to delete this sync folder?\n\n"+
							folder.SourcePath+" → "+folder.DestinationPath,
						func(confirmed bool) {
							if confirmed {
								deleteFolder(cfg, store, index, refreshFolders, reloadJobsFunc)
							}
						},
						modal,
					)
					confirmDialog.Show()
				})

				buttonRow := container.NewGridWithColumns(2, editBtn, deleteBtn)

				folderCard := container.NewVBox(
					folderLabel,
					statusLabel,
					buttonRow,
					widget.NewSeparator(),
				)

				widgets = append(widgets, folderCard)
			}
		}

		// Add folder button at the bottom
		addButton := widget.NewButton("Add New Folder", func() {
			addOrEditFolder(cfg, store, -1, refreshFolders, reloadJobsFunc)
		})
		widgets = append(widgets, addButton)

		// Update the container
		folderContainer.Objects = widgets
		folderContainer.Refresh()
	}

	// Initialize the container
	folderContainer = container.NewVBox()

	// Build initial list
	refreshFolders()

	// Wrap in a scroll container
	scrollContainer := container.NewVScroll(folderContainer)

	modal.SetContent(scrollContainer)
	modal.Resize(fyne.NewSize(600, 500))

	return modal
}

// opens general Settings modal
func (w *SettingsWindow) changeGeneralSettings() {
	modal := w.app.NewWindow("General Settings")

	currentMinutes := int(w.config.CheckInterval.Minutes())

	// Entry for minutes
	minutesEntry := widget.NewEntry()
	minutesEntry.SetPlaceHolder("Minutes")
	minutesEntry.SetText(strconv.Itoa(currentMinutes))

	infoLabel := widget.NewLabel(fmt.Sprintf("Current interval: %v", w.config.CheckInterval))

	startAtBootCheck := widget.NewCheck("Start MirrorBox at login", nil)
	startAtBootCheck.SetChecked(w.config.StartAtBoot)

	saveButton := widget.NewButton("Save", func() {
		minutes, err := strconv.Atoi(minutesEntry.Text)
		if err != nil || minutes <= 0 {
			dialog.ShowError(
				fmt.Errorf("please enter a valid number of minutes (greater than 0)"),
				modal,
			)
			return
		}

		w.config.CheckInterval = time.Duration(minutes) * time.Minute
		w.config.StartAtBoot = startAtBootCheck.Checked

		if err := w.saveConfigAndReload(); err != nil {
			log.Printf("Failed to save config and reload jobs: %v", err)
			dialog.ShowError(
				fmt.Errorf("failed to save configuration: %w", err),
				modal,
			)
			return
		}

		infoLabel.SetText(fmt.Sprintf("Current interval: %v", w.config.CheckInterval))
		startAtBootCheck.SetChecked(w.config.StartAtBoot)

		log.Printf("Updated check interval to %v and reloaded jobs", w.config.CheckInterval)

		dialog.ShowInformation(
			"Success",
			fmt.Sprintf("Check interval updated to %v.\nRestart the application for changes to take effect.", w.config.CheckInterval),
			modal,
		)

		modal.Close()
	})

	cancelButton := widget.NewButton("Cancel", func() {
		modal.Close()
	})

	form := container.NewVBox(
		widget.NewLabelWithStyle("Change Check Interval", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		infoLabel,
		widget.NewLabel("Enter interval in minutes:"),
		minutesEntry,
		widget.NewLabel("Note: Application restart required for interval changes to take effect."),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Startup", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		startAtBootCheck,
		widget.NewLabel("Enable to launch MirrorBox automatically after you log in."),
		widget.NewSeparator(),
		container.NewGridWithColumns(2, cancelButton, saveButton),
	)

	modal.SetContent(form)
	modal.Resize(fyne.NewSize(500, 300))
	modal.Show()
}

// getLastJobStatus returns a formatted string with the last job execution status.
func (w *SettingsWindow) getLastJobStatus() string {
	jobs := w.state.AllJobs()
	if len(jobs) == 0 {
		return "No jobs configured"
	}

	var lastJob *syncpkg.Job
	var lastTime time.Time

	// Find the most recently run job
	for _, job := range jobs {
		if job.LastRun().After(lastTime) {
			lastTime = job.LastRun()
			lastJob = job
		}
	}

	if lastJob == nil || lastTime.IsZero() {
		return "No jobs have run yet"
	}

	statusText := fmt.Sprintf("Last Job: %s\n", lastJob.Name)
	statusText += fmt.Sprintf("Status: %v\n", formatJobStatus(lastJob.Status()))
	statusText += fmt.Sprintf("Run Time: %v\n", lastTime.Format("2006-01-02 15:04:05"))

	if result := lastJob.LastResult(); result != nil {
		statusText += fmt.Sprintf("Created: %d, Updated: %d, Deleted: %d",
			result.FilesCreated, result.FilesUpdated, result.FilesDeleted)
	}

	if err := lastJob.LastError(); err != nil {
		statusText += fmt.Sprintf("\nError: %v", err)
	}

	return statusText
}

// formatJobStatus converts JobStatus to a human-readable string.
func formatJobStatus(status syncpkg.JobStatus) string {
	switch status {
	case syncpkg.StatusIdle:
		return "Idle"
	case syncpkg.StatusRunning:
		return "Running"
	case syncpkg.StatusSuccess:
		return "Success"
	case syncpkg.StatusError:
		return "Error"
	default:
		return "Unknown"
	}
}

// refreshStatus updates the status widget with current job status.
// Must be called from the main UI thread.
func (w *SettingsWindow) refreshStatus() {
	if w.statusWidget != nil {
		w.statusWidget.SetText(w.getLastJobStatus())
	}
}

func (w *SettingsWindow) Show() {
	if w.window == nil {
		w.window = w.app.NewWindow("MirrorBox - Settings")
		w.window.Resize(fyne.NewSize(600, 550))

		w.statusWidget = widget.NewLabel(w.getLastJobStatus())
		w.statusWidget.Wrapping = fyne.TextWrapWord

		buttonsSection := container.NewVBox(
			widget.NewLabelWithStyle(
				"Configuration",
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true},
			),
			widget.NewButton(
				"Manage Sync Folders",
				func() {
					folderWindow := NewFolderWindow(w.app, w.config, w.store, w.reloadJobs)
					folderWindow.Show()
				},
			),
			widget.NewButton(
				"General Settings",
				func() {
					w.changeGeneralSettings()
				},
			),
			widget.NewButton(
				"View Job Status",
				func() {
					if w.statusWindow != nil {
						w.statusWindow.Show()
					}
				},
			),
		)

		statusSection := container.NewVBox(
			widget.NewSeparator(),
			widget.NewLabelWithStyle(
				"Last Job Execution",
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true},
			),
			w.statusWidget,
			widget.NewButton("Refresh Status", func() {
				w.refreshStatus()
			}),
		)

		form := container.NewVBox(
			buttonsSection,
			statusSection,
		)

		scrollContainer := container.NewVScroll(form)

		w.window.SetContent(scrollContainer)

		w.window.SetOnClosed(func() {
			w.window = nil
			w.statusWidget = nil
		})
	}

	// Refresh status when window is shown
	w.refreshStatus()
	w.window.Show()
	w.window.RequestFocus()
}

func (w *SettingsWindow) Hide() {
	if w.window != nil {
		w.window.Hide()
	}
}

func (w *SettingsWindow) OnSave(newConfig *config.Config) error {
	w.config = newConfig
	return w.saveConfigAndReload()
}

func (w *SettingsWindow) GetConfig() *config.Config {
	return w.config
}

func (w *SettingsWindow) UpdateJobStatus() {
	fyne.Do(func() {
		w.refreshStatus()
	})
}
