package filemanagement

import (
	"os"
	"path/filepath"
	"runtime"

	http2 "github.com/samtv12345/gnpm/pm/impl/pnpm/http"
	"github.com/samtv12345/gnpm/pm/interfaces"
	"go.uber.org/zap"
)

func SavePackageManager(result *http2.DownloadReleaseResult, logger *zap.SugaredLogger, version string, pm interfaces.IPackageManager) (*string, error) {
	dataDir, err := EnsureDataDir()
	if err != nil {
		return nil, err
	}
	gnpmDir := filepath.Join(*dataDir, ".cache")
	_, err = os.Stat(gnpmDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(gnpmDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	locationToWritePnpm := filepath.Join(gnpmDir, pm.GetName()+"-"+version+filepath.Ext(result.Filename))
	_, err = os.Stat(locationToWritePnpm)
	if os.IsNotExist(err) {
	} else if err != nil {
		return nil, err
	}

	logger.Debugf("Saving pnpm to cache at location: %s with name %s", locationToWritePnpm, result.Filename)
	err = os.WriteFile(locationToWritePnpm, result.Content, 0644)
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

func IsPackageManagerInstalled(version string, pmManager interfaces.IPackageManager) (*bool, *string, error) {
	if version == "*" {
		return &[]bool{false}[0], nil, nil
	}
	dataDir, err := EnsureDataDir()
	if err != nil {
		return nil, nil, err
	}
	locationToCheck := filepath.Join(*dataDir, "_gnpm", pmManager.GetName()+"-"+version)
	locationToCheck2 := filepath.Join(*dataDir, "_gnpm", pmManager.GetName()+"-"+version+".exe")
	_, err = os.Stat(locationToCheck)
	_, err2 := os.Stat(locationToCheck2)

	if os.IsNotExist(err) && os.IsNotExist(err2) {
		var falseVal = false
		return &falseVal, nil, nil
	}
	if err != nil && err2 != nil {
		return nil, nil, err
	}

	if err != nil {
		locationToCheck = locationToCheck2
	}

	return &[]bool{true}[0], &locationToCheck, nil
}
