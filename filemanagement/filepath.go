package filemanagement

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

const DataDir = "gnpm"

func EncodePath(path string) (*string, error) {
	hash := sha256.Sum256([]byte(path))
	hashHex := hex.EncodeToString(hash[:8])
	symlinkName := "gnpm-" + hashHex
	moduledir, err := EnsureModuleDir()
	if err != nil {
		return nil, err
	}
	var moduleDir = filepath.Join(*moduledir, symlinkName)

	return &moduleDir, nil
}

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

func FindParentLinkModuleDir(startDir string) (*string, error) {
	potentialPackageJsonDir, err := EncodePath(startDir)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(*potentialPackageJsonDir)
	if os.IsNotExist(err) {
		parentDir := filepath.Dir(startDir)
		if parentDir == startDir {
			return nil, errors.New("could not find gnpm modules directory in any parent directories")
		}
		return FindParentLinkModuleDir(parentDir)
	}
	if err != nil {
		return nil, err
	}
	return potentialPackageJsonDir, nil
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

func PathExists(path string, isDirectory bool) bool {
	result, err := os.Stat(path)
	if err != nil {
		return false
	}

	if isDirectory {
		return result.IsDir()
	}

	return !result.IsDir()
}
