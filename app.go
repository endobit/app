// Package app provides utilities for applications.
package app

import (
	"os"
	"path/filepath"
	"runtime"
)

// UserConfigFilePath returns the path to the directory containing an
// application's config file. The path is determined based on the operating
// system:
//   - On Windows, it uses the `AppData` environment variable.
//   - On Unix-like systems, it checks the `XDG_CONFIG_HOME` environment variable.
//     If `XDG_CONFIG_HOME` is not set, it defaults to `~/.config/<appName>`.
func UserConfigFilePath(appName string) (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		baseDir = os.Getenv("AppData")
	default:
		xdg := os.Getenv("XDG_CONFIG_HOME")
		if xdg != "" {
			baseDir = xdg
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}

			baseDir = filepath.Join(home, ".config")
		}
	}

	return filepath.Join(baseDir, appName), nil
}
