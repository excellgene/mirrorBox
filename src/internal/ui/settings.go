package ui

import (
	"log"

	"excellgene.com/symbaSync/internal/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// SettingsWindow manages the settings/configuration window.
type SettingsWindow struct {
	app    fyne.App
	window fyne.Window

	config *config.Config
	store  *config.Store
}

// NewSettingsWindow creates a new settings window.
// IMPORTANT: the fyne.App MUST be injected (never create a new one here).
func NewSettingsWindow(
	app fyne.App,
	cfg *config.Config,
	store *config.Store,
) *SettingsWindow {
	return &SettingsWindow{
		app:    app,
		config: cfg,
		store:  store,
	}
}

// Show displays the settings window.
// The window is created lazily and reused.
func (w *SettingsWindow) Show() {
	if w.window == nil {
		w.window = w.app.NewWindow("SambaSync - Settings")
		w.window.Resize(fyne.NewSize(520, 420))


		form := container.NewVBox(
			widget.NewLabelWithStyle(
				"Configuration",
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true},
			),
		)

		w.window.SetContent(form)

		// When user closes the window via the window manager
		w.window.SetOnClosed(func() {
			w.window = nil
			log.Println("settings window closed")
		})
	}

	w.window.Show()
	w.window.RequestFocus()
	log.Println("settings window opened")
}

// Hide closes the settings window if open.
func (w *SettingsWindow) Hide() {
	if w.window != nil {
		w.window.Hide()
	}
}

// OnSave allows programmatic save if needed (optional use).
func (w *SettingsWindow) OnSave(newConfig *config.Config) error {
	w.config = newConfig
	return w.store.Save(newConfig)
}

// GetConfig returns the current configuration.
func (w *SettingsWindow) GetConfig() *config.Config {
	return w.config
}