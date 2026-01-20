package ui

import (
	"log"

	"excellgene.com/symbaSync/internal/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
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

func deleteFolder(cfg *config.Config, store *config.Store, folder config.FolderToSync, modal fyne.Window) {
	// Remove folder from config
	var updatedFolders []config.FolderToSync
	for _, f := range cfg.Folders {
		if f != folder {
			updatedFolders = append(updatedFolders, f)
		}
	}
	cfg.Folders = updatedFolders

	// Save updated config
	err := store.Save(cfg)
	if err != nil {
		log.Printf("Failed to save config after deleting folder: %v", err)
		return
	}

	// Close the modal
	modal.Close()
}

func addFolder(cfg *config.Config, store *config.Store, folder config.FolderToSync) {
	modal := fyne.CurrentApp().NewWindow("Add Sync Folder")

	// Source entry
	sourceEntry := widget.NewEntry()
	sourceEntry.SetPlaceHolder("Source Path")

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

	destinationBtn := widget.NewButton("Browse…", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				destinationEntry.SetText(uri.Path())
			}
		}, modal).Show()
	})

	enabledCheck := widget.NewCheck("Enabled", func(checked bool) {
		folder.Enabled = checked
	})
	enabledCheck.SetChecked(true)

	saveButton := widget.NewButton("Save", func() {
		folder.SourcePath = sourceEntry.Text
		folder.DestinationPath = destinationEntry.Text
		folder.Enabled = enabledCheck.Checked

		cfg.Folders = append(cfg.Folders, folder)

		log.Printf("Adding new folder to sync: %s -> %s", folder.SourcePath, folder.DestinationPath)
		if err := store.Save(cfg); err != nil {
			log.Printf("Failed to save config after adding folder: %v", err)
			return
		}

		modal.Close()
	})

	form := container.NewVBox(
		widget.NewLabel("Add New Sync Folder"),

		widget.NewLabel("Source"),
		container.NewBorder(nil, nil, nil, sourceBtn, sourceEntry),

		widget.NewLabel("Destination"),
		container.NewBorder(nil, nil, nil, destinationBtn, destinationEntry),

		enabledCheck,
		saveButton,
	)

	modal.SetContent(form)
	modal.Resize(fyne.NewSize(700, 400))
	modal.Show()
}

// NewFolderWindow creates a new folder window.
func NewFolderWindow(app fyne.App, cfg *config.Config, store *config.Store) fyne.Window {
	modal := fyne.CurrentApp().NewWindow("Syncing Folders") 

	for _, folder := range cfg.Folders {
		label := widget.NewLabel("Sync: "+ folder.SourcePath + " <-> " + folder.DestinationPath)
		deleteButton := widget.NewButton("Delete", func() {
			deleteFolder(cfg, store, folder, modal)
		})

		modal.SetContent(container.NewVBox(label, deleteButton))
	}

	addButton := widget.NewButton("Add New Folder", func() {
		addFolder(cfg, store, config.FolderToSync{})
	})

	modal.SetContent(container.NewVBox(addButton))
	modal.Resize(fyne.NewSize(400, 300))
	
	return modal
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
			widget.NewButton(
				"Syncing Folders",
				func() {
					folderWindow := NewFolderWindow(w.app, w.config, w.store)
					folderWindow.Show()
				},
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