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

	// Set the tray icon based on system theme
	icon := m.tray.getTrayIcon()
	desktopApp.SetSystemTrayIcon(icon)

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

	m.tray.app.Settings().AddListener(func(s fyne.Settings) {
		m.UpdateIcon()
	})

	desktopApp.SetSystemTrayMenu(m.menu)
}

func (m *Menu) SetStatusText(text string) {
	if m.status == nil || m.menu == nil {
		return
	}

	m.status.Label = "Status: " + text

	if desktopApp, ok := m.tray.app.(desktop.App); ok {
		desktopApp.SetSystemTrayMenu(m.menu)
	}
}

func (m *Menu) UpdateIcon() {
	if desktopApp, ok := m.tray.app.(desktop.App); ok {
		icon := m.tray.getTrayIcon()
		desktopApp.SetSystemTrayIcon(icon)
	}
}