package tray

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// Menu holds references to system tray menu items.
type Menu struct {
	tray   *Tray
	menu   *fyne.Menu
	status *fyne.MenuItem
}

// NewMenu creates a new menu structure.
func NewMenu(tray *Tray) *Menu {
	return &Menu{tray: tray}
}

// Build creates the menu items and adds them to the tray.
func (m *Menu) Build() {
	desktopApp, ok := m.tray.app.(desktop.App)
	if !ok {
		return
	}

	syncNow := fyne.NewMenuItem("Sync Now", func() {
		m.tray.events <- EventSyncNow
	})

	m.status = fyne.NewMenuItem("Status: Idle", func() {
		m.tray.events <- EventStatus
	})

	settings := fyne.NewMenuItem("Settings", func() {
		m.tray.events <- EventSettings
	})

	quit := fyne.NewMenuItem("Quit", func() {
		m.tray.events <- EventQuit
		m.tray.Quit()
	})

	m.menu = fyne.NewMenu(
		"SambaSync",
		syncNow,
		m.status,
		fyne.NewMenuItemSeparator(),
		settings,
		fyne.NewMenuItemSeparator(),
		quit,
	)

	desktopApp.SetSystemTrayMenu(m.menu)
}

// SetStatusText updates the status menu item text.
func (m *Menu) SetStatusText(text string) {
	if m.status == nil || m.menu == nil {
		return
	}

	m.status.Label = "Status: " + text

	if desktopApp, ok := m.tray.app.(desktop.App); ok {
		desktopApp.SetSystemTrayMenu(m.menu)
	}
}