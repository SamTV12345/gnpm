package filemanagement

import (
	"os"
	"path/filepath"

	http2 "github.com/samtv12345/gnpm/pm/impl/pnpm/http"
	"github.com/samtv12345/gnpm/pm/interfaces"
	"go.uber.org/zap"
)

func SavePackageManager(result *http2.DownloadReleaseResult, logger *zap.SugaredLogger, version string, pm interfaces.IPackageManager) (*string, error) {
	gnpmDir, err := GetCacheDir()
	if err != nil {
		return nil, err
	}

	locationToWritePnpm := filepath.Join(*gnpmDir, pm.GetName()+"-"+version+filepath.Ext(result.Filename))
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

func IsPackageManagerInstalled(version string, pmManager interfaces.IPackageManager) (*bool, *string, error) {
	if version == "*" {
		return &[]bool{false}[0], nil, nil
	}
	gnpmDir, err := GetGnpmDir()
	if err != nil {
		return nil, nil, err
	}
	locationToCheck := filepath.Join(*gnpmDir, pmManager.GetName()+"-"+version)
	locationToCheck2 := filepath.Join(*gnpmDir, pmManager.GetName()+"-"+version+".exe")
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
