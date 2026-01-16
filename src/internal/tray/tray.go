package tray

import (
	"log"

	"github.com/getlantern/systray"
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
// Responsibility:
//   - Display system tray icon
//   - Show menu items
//   - Emit events when user clicks menu items
//
type Tray struct {
	events chan Event
	menu   *Menu
}

// New creates a new system tray manager.
func New() *Tray {
	return &Tray{
		events: make(chan Event, 10),
	}
}

// Events returns a channel for receiving tray events.
func (t *Tray) Events() <-chan Event {
	return t.events
}

// Run starts the system tray.
// This is a blocking call that runs the tray event loop.
// Call this in a goroutine or as the last thing in main().
func (t *Tray) Run() {
	systray.Run(t.onReady, t.onExit)
}

// Quit stops the system tray and exits the application.
func (t *Tray) Quit() {
	systray.Quit()
}

// UpdateStatus updates the tray tooltip/status text.
func (t *Tray) UpdateStatus(status string) {
	systray.SetTooltip(status)
}

// onReady is called when systray is ready.
// Sets up the icon and menu.
func (t *Tray) onReady() {
	systray.SetTitle("SambaSync")
	systray.SetTooltip("SambaSync - Idle")

	// TODO: Set icon
	// systray.SetIcon(icon.Data)

	// Create menu
	t.menu = NewMenu()
	t.menu.Build()

	// Listen for menu events
	go t.handleMenuEvents()

	log.Println("System tray ready")
}

// onExit is called when systray is exiting.
func (t *Tray) onExit() {
	close(t.events)
	log.Println("System tray exited")
}

// handleMenuEvents listens to menu item clicks and emits events.
func (t *Tray) handleMenuEvents() {
	for {
		select {
		case <-t.menu.syncNow.ClickedCh:
			t.events <- EventSyncNow

		case <-t.menu.status.ClickedCh:
			t.events <- EventStatus

		case <-t.menu.settings.ClickedCh:
			t.events <- EventSettings

		case <-t.menu.quit.ClickedCh:
			t.events <- EventQuit
			return
		}
	}
}
