package main

import (
	"log"
	"os"
	"path/filepath"

	"excellgene.com/mirrorBox/internal/app"
	"excellgene.com/mirrorBox/internal/config"
	"excellgene.com/mirrorBox/internal/tray"
	"excellgene.com/mirrorBox/internal/ui"
	"fyne.io/fyne/v2"
)

func main() {
	configPath := getConfigPath()
	configStore := config.NewStore(configPath)

	cfg, err := configStore.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appState := app.NewState()

	jobFactory := app.NewJobFactory()

	jobs, err := jobFactory.CreateFromConfig(cfg)
	if err != nil {
		log.Fatalf("Failed to create jobs: %v", err)
	}

	for _, job := range jobs {
		appState.AddJob(job)
		log.Printf("Registered job: %s", job.Name)
	}

	dispatcher := app.NewDispatcher(appState)

	if cfg.CheckInterval > 0 {
		log.Printf("Starting scheduler with interval: %v", cfg.CheckInterval)
		dispatcher.StartScheduler(cfg.CheckInterval)
	}
	systemTray := tray.New()

	statusWindow := ui.NewStatusWindow(
		systemTray.App(),
		appState,
	)

	settingsWindow := ui.NewSettingsWindow(
		systemTray.App(),
		cfg,
		configStore,
		appState,
		statusWindow,
		jobFactory,
	)

	go handleTrayEvents(systemTray, dispatcher, settingsWindow, statusWindow)
	go handleDispatcherEvents(dispatcher, statusWindow, settingsWindow, systemTray)

	systemTray.Run()
	dispatcher.Stop()
}

// handleTrayEvents processes user actions from system tray.
func handleTrayEvents(
	systemTray *tray.Tray,
	dispatcher *app.Dispatcher,
	settings *ui.SettingsWindow,
	status *ui.StatusWindow,
) {
	for event := range systemTray.Events() {
		switch event {
		case tray.EventSyncNow:
			log.Println("User triggered sync")
			dispatcher.RunAll()

		case tray.EventSettings:
			log.Println("User opened settings")
			fyne.Do(func() {
				settings.Show()
			})

		case tray.EventStatus:
			log.Println("User opened status")
			fyne.Do(func() {
				status.Show()
			})

		case tray.EventQuit:
			log.Println("User quit application")
			fyne.Do(func() {
				systemTray.Quit()
			})
			return
		}
	}
}

// handleDispatcherEvents processes job events and updates UI.
func handleDispatcherEvents(
	dispatcher *app.Dispatcher,
	status *ui.StatusWindow,
	settings *ui.SettingsWindow,
	systemTray *tray.Tray,
) {
	for event := range dispatcher.Events() {
		fyne.Do(func() {
			status.OnJobEvent(event)
			settings.UpdateJobStatus()
			systemTray.UpdateStatus(formatJobStatus(event))
		})
	}
}

// formatJobStatus creates a human-readable status string.
func formatJobStatus(event app.JobEvent) string {
	switch event.Status {
	case 0:
		return "MirrorBox - Idle"
	case 1:
		return "MirrorBox - Syncing..."
	case 2:
		return "MirrorBox - Last sync successful"
	case 3:
		return "MirrorBox - Last sync failed"
	default:
		return "MirrorBox"
	}
}

// getConfigPath returns the path to the config file.
func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".config", "mirrorbox")
	return filepath.Join(configDir, "config.json")
}
