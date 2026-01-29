//go:build darwin || cgo

package app

import (
	"github.com/emersion/go-autostart"

	"excellgene.com/mirrorBox/internal/config"
)

func EnableAutoStart(appName, appPath string, cfg *config.Config) error {
	autostartApp := &autostart.App{
		Name:        appName,
		DisplayName: appName,
		Exec:        []string{appPath},
	}

	if cfg.StartAtBoot && !autostartApp.IsEnabled() {
		err := autostartApp.Enable()
		if err != nil {
			return err
		}
	}

	return nil
}

func DisableAutoStart(appName string) error {
	autostartApp := &autostart.App{
		Name:        appName,
		DisplayName: appName,
	}

	if autostartApp.IsEnabled() {
		err := autostartApp.Disable()
		if err != nil {
			return err
		}
	}

	return nil
}
