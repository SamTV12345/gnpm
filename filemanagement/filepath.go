package filemanagement

import (
	"os"
	"path/filepath"
	"runtime"
)

const DataDir = "gnpm"

func EnsureDataDir() (string, error) {
	dir := userDataDir()
	if dir == "" {
		return "", nil
	}
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func userDataDir() string {
	switch runtime.GOOS {
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "Roaming", DataDir)
		}
	case "darwin":
		home := os.Getenv("HOME")
		if home != "" {
			return filepath.Join(home, "Library", "Application Support", DataDir)
		}
	default: // Linux und andere Unix
		if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
			return filepath.Join(xdg, DataDir)
		}
		home := os.Getenv("HOME")
		if home != "" {
			return filepath.Join(home, ".local", "share", DataDir)
		}
	}
	return ""
}
