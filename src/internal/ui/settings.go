package ui

import (
	"log"

	"excellgene.com/symbaSync/internal/config"
)

// SettingsWindow manages the settings/configuration window.
// Responsibility:
//   - Display configuration UI
//   - Allow user to modify sync jobs, paths, schedules
//   - Save changes back to config
//
// NO BUSINESS LOGIC. Only UI presentation and user input handling.
//
// Implementation options:
//   - Wails (web-based UI with Go backend)
//   - Native GUI framework (fyne, gio, etc.)
//   - Web server + browser
type SettingsWindow struct {
	config *config.Config
	store  *config.Store
}

// NewSettingsWindow creates a new settings window.
func NewSettingsWindow(cfg *config.Config, store *config.Store) *SettingsWindow {
	return &SettingsWindow{
		config: cfg,
		store:  store,
	}
}

// Show displays the settings window.
// Placeholder: In real implementation, this would open a GUI window.
func (w *SettingsWindow) Show() {
	log.Println("Settings window opened")
	// TODO: Implement actual UI
	// Options:
	// - Wails: runtime.WindowShow()
	// - Native: create and show window
	// - Web: open browser to localhost:port/settings
}

// Hide closes the settings window.
func (w *SettingsWindow) Hide() {
	log.Println("Settings window closed")
	// TODO: Implement actual UI
}

// OnSave is called when user clicks Save in settings.
// Updates config and persists to disk.
func (w *SettingsWindow) OnSave(newConfig *config.Config) error {
	w.config = newConfig
	return w.store.Save(newConfig)
}

// GetConfig returns the current configuration.
func (w *SettingsWindow) GetConfig() *config.Config {
	return w.config
}
