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

var directoriesInGnpm = []string{
	"modules",
	".cache",
	"_gnpm",
}

var files = map[string]string{}

func init() {
	var userDir = dataDir()
	if err := os.MkdirAll(userDir, os.ModePerm); err != nil {
		panic("Could not create gnpm data directory: " + err.Error())
	}
	for _, dir := range directoriesInGnpm {
		err := os.MkdirAll(filepath.Join(userDir, dir), os.ModePerm)
		if err != nil {
			panic("Could not create gnpm subdirectory: " + dir + err.Error())
		}
		files[dir] = filepath.Join(userDir, dir)
	}
}

func EncodePath(path string) (*string, error) {
	hash := sha256.Sum256([]byte(path))
	hashHex := hex.EncodeToString(hash[:8])
	symlinkName := "gnpm-" + hashHex
	moduledir, err := GetModuleDir()
	if err != nil {
		return nil, err
	}
	var moduleDir = filepath.Join(*moduledir, symlinkName)

	return &moduleDir, nil
}

func GetModuleDir() (*string, error) {
	moduleDir, exists := files["modules"]
	if !exists {
		return nil, errors.New("could not find gnpm modules directory")
	}
	return &moduleDir, nil
}

func GetCacheDir() (*string, error) {
	cacheDir, exists := files[".cache"]
	if !exists {
		return nil, errors.New("could not find gnpm cache directory")
	}
	return &cacheDir, nil
}

func GetGnpmDir() (*string, error) {
	gnpmDir, exists := files["_gnpm"]
	if !exists {
		return nil, errors.New("could not find gnpm _gnpm directory")
	}
	return &gnpmDir, nil
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

func dataDir() string {
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
