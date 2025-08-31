package filemanagement

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/samtv12345/gnpm/http"
	"go.uber.org/zap"
)

func SavePnpmToInstallDir(result *http.DownloadPnpmReleaseResult, logger *zap.SugaredLogger, version string) (*string, error) {
	dataDir, err := EnsureDataDir()
	if err != nil {
		return nil, err
	}
	gnpmDir := filepath.Join(*dataDir, "_gnpm")
	_, err = os.Stat(gnpmDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(gnpmDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	locationToWritePnpm := filepath.Join(gnpmDir, "pnpm-"+version)
	_, err = os.Stat(locationToWritePnpm)
	if os.IsNotExist(err) {
		err = os.Mkdir(locationToWritePnpm, os.ModePerm)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	logger.Debugf("Saving pnpm to cache at location: %s with name %s", locationToWritePnpm, result.Filename)
	err = os.WriteFile(filepath.Join(locationToWritePnpm, result.Filename), result.Content, 0644)
	if err != nil {
		return nil, err
	}
	return &locationToWritePnpm, nil
}

func buildPnpmFilename() string {
	operatingSystem := runtime.GOOS
	architecture := runtime.GOARCH

	fileSuffix := ""
	filePrefix := "pnpm-"

	if architecture == "amd64" {
		architecture = "x64"
	}

	if operatingSystem == "darwin" {
		filePrefix += "win-" + architecture
	}

	if operatingSystem == "windows" {
		fileSuffix = ".exe"
		filePrefix += "win-" + architecture
	} else if operatingSystem == "darwin" {
		filePrefix += "macos-" + architecture
	} else {
		filePrefix += "linux-" + architecture
	}

	return filePrefix + fileSuffix
}

func IsPnpmVersionInInstallDir(version string) (*bool, *string, error) {
	dataDir, err := EnsureDataDir()
	if err != nil {
		return nil, nil, err
	}
	locationToCheck := filepath.Join(*dataDir, "_gnpm", "pnpm-"+version)
	_, err = os.Stat(locationToCheck)

	if os.IsNotExist(err) {
		var falseVal = false
		return &falseVal, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	filenameInPnpmDir := buildPnpmFilename()
	locationToCheck = filepath.Join(*dataDir, "_gnpm", "pnpm-"+version, filenameInPnpmDir)
	_, err = os.Stat(locationToCheck)

	if os.IsNotExist(err) {
		var falseVal = false
		return &falseVal, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	return &[]bool{true}[0], &locationToCheck, nil
}
