package main

import (
	"log"
	"os"
	"path/filepath"

	"excellgene.com/symbaSync/internal/app"
	"excellgene.com/symbaSync/internal/config"
	"excellgene.com/symbaSync/internal/infra/smb"
	"excellgene.com/symbaSync/internal/tray"
	"excellgene.com/symbaSync/internal/ui"
)

func main() {
	log.Println("Starting SambaSync...")

	// Initialize config
	configPath := getConfigPath()
	configStore := config.NewStore(configPath)

	cfg, err := configStore.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize application state
	appState := app.NewState()

	// Create job factory with SMB client factory
	jobFactory := app.NewJobFactory(func(cfg smb.Config) smb.Client {
		// In production, replace with real SMB client
		return smb.NewMockClient(cfg)
	})

	// Create jobs from config
	jobs, err := jobFactory.CreateFromConfig(cfg)
	if err != nil {
		log.Fatalf("Failed to create jobs: %v", err)
	}

	// Register jobs with state
	for _, job := range jobs {
		appState.AddJob(job)
		log.Printf("Registered job: %s", job.Name)
	}

	// Initialize dispatcher
	dispatcher := app.NewDispatcher(appState)

	// Start scheduler if interval is configured
	if cfg.CheckInterval > 0 {
		log.Printf("Starting scheduler with interval: %v", cfg.CheckInterval)
		dispatcher.StartScheduler(cfg.CheckInterval)
	}

	// Initialize UI components
	settingsWindow := ui.NewSettingsWindow(cfg, configStore)
	statusWindow := ui.NewStatusWindow(appState)

	// Initialize system tray
	systemTray := tray.New()

	// Handle tray events
	go handleTrayEvents(systemTray, dispatcher, settingsWindow, statusWindow)

	// Listen to dispatcher events and update UI
	go handleDispatcherEvents(dispatcher, statusWindow, systemTray)

	log.Println("SambaSync started successfully")

	// Run system tray (blocking)
	systemTray.Run()

	// Cleanup on exit
	log.Println("Shutting down...")
	dispatcher.Stop()
	log.Println("Goodbye!")
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
			settings.Show()

		case tray.EventStatus:
			log.Println("User opened status")
			status.Show()

		case tray.EventQuit:
			log.Println("User quit application")
			systemTray.Quit()
			return
		}
	}
}

// handleDispatcherEvents processes job events and updates UI.
func handleDispatcherEvents(
	dispatcher *app.Dispatcher,
	status *ui.StatusWindow,
	systemTray *tray.Tray,
) {
	for event := range dispatcher.Events() {
		// Update status window
		status.OnJobEvent(event)

		// Update tray tooltip
		statusText := formatJobStatus(event)
		systemTray.UpdateStatus(statusText)
	}
}

// formatJobStatus creates a human-readable status string.
func formatJobStatus(event app.JobEvent) string {
	switch event.Status {
	case 0: // StatusIdle
		return "SambaSync - Idle"
	case 1: // StatusRunning
		return "SambaSync - Syncing..."
	case 2: // StatusSuccess
		return "SambaSync - Last sync successful"
	case 3: // StatusError
		return "SambaSync - Last sync failed"
	default:
		return "SambaSync"
	}
}

// getConfigPath returns the path to the config file.
// Uses platform-specific config directory.
func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	// Use ~/.config/sambasync/config.json on Unix-like systems
	// Use %APPDATA%/sambasync/config.json on Windows
	configDir := filepath.Join(homeDir, ".config", "sambasync")
	return filepath.Join(configDir, "config.json")
}
