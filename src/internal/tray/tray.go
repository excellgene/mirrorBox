package tray

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"excellgene.com/symbaSync/internal/ui/icons"
)

// Event represents a user action from the system tray.
type Event int

const (
	EventSyncNow Event = iota
	EventSettings
	EventStatus
	EventQuit
)

// Tray manages the system tray icon and menu.
type Tray struct {
	app                fyne.App
	events             chan Event
	menu               *Menu
	unregisterCallback func()
}

// New creates a new tray manager.
func New() *Tray {
	return &Tray{
		app:    app.NewWithID("com.excellgene.sambasync"),
		events: make(chan Event, 10),
	}
}

// getTrayIcon returns the appropriate icon based on the system theme.
func (t *Tray) getTrayIcon() fyne.Resource {
	if t.app.Settings().ThemeVariant() == theme.VariantDark {
		return icons.TrayWhite()
	}
	return icons.TrayBlack()
}

// Events returns a channel for receiving tray events.
func (t *Tray) Events() <-chan Event {
	return t.events
}

// Run starts the Fyne application.
func (t *Tray) Run() {
	t.menu = NewMenu(t)
	t.menu.Build()

	log.Println("System tray ready (Fyne)")
	t.app.Run()

	close(t.events)
	log.Println("System tray exited")
}

// Quit cleanly exits the application.
func (t *Tray) Quit() {
	t.app.Quit()
}

// UpdateStatus updates the tray tooltip/status text.
func (t *Tray) UpdateStatus(status string) {
	t.menu.SetStatusText(status)
}

func (t *Tray) App() fyne.App {
	return t.app
}
