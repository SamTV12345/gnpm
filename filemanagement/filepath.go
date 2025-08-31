package filemanagement

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

const DataDir = "gnpm"

func EnsureDataDir() (*string, error) {
	dir := userDataDir()
	if dir == "" {
		return nil, errors.New(DataDir + " not found in $PATH")
	}
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(filepath.Join(dir, ".cache"), os.ModePerm)
	if err != nil {
		return nil, err
	}
	return &dir, nil
}

func EnsureModuleDir() (*string, error) {
	dataDir, err := EnsureDataDir()
	if err != nil {
		return nil, err
	}
	modulesDir := filepath.Join(*dataDir, "modules")
	err = os.MkdirAll(modulesDir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return &modulesDir, nil
}

func userDataDir() string {
	switch runtime.GOOS {
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, DataDir)
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
