package tray

import (
	"github.com/getlantern/systray"
)

// Menu holds references to system tray menu items.
// Responsibility: Define menu structure and items.
type Menu struct {
	syncNow  *systray.MenuItem
	status   *systray.MenuItem
	settings *systray.MenuItem
	quit     *systray.MenuItem
}

// NewMenu creates a new menu structure.
func NewMenu() *Menu {
	return &Menu{}
}

// Build creates the menu items and adds them to the tray.
func (m *Menu) Build() {
	// Main actions
	m.syncNow = systray.AddMenuItem("Sync Now", "Run sync immediately")
	m.status = systray.AddMenuItem("Status", "Show sync status")

	systray.AddSeparator()

	// Configuration
	m.settings = systray.AddMenuItem("Settings", "Open settings")

	systray.AddSeparator()

	// Exit
	m.quit = systray.AddMenuItem("Quit", "Exit SambaSync")
}

// SetSyncEnabled enables or disables the sync now button.
// Useful when a sync is already running.
func (m *Menu) SetSyncEnabled(enabled bool) {
	if enabled {
		m.syncNow.Enable()
	} else {
		m.syncNow.Disable()
	}
}

// SetStatusText updates the status menu item text.
func (m *Menu) SetStatusText(text string) {
	m.status.SetTitle(text)
}
