package caching

import (
	"os"
	"path/filepath"

	"github.com/samtv12345/gnpm/filemanagement"
	"github.com/samtv12345/gnpm/http"
	"go.uber.org/zap"
)

func GetPnpmVersion(logger *zap.SugaredLogger) []string {
	cacheDir, err := filemanagement.GetCacheDir()
	if err != nil {
		logger.Warnf("Error getting cache dir: %v", err)
		return []string{}
	}
	pnpmJsonFile := filepath.Join(*cacheDir, "pnpm.json")

	fsInfo, err := os.Stat(pnpmJsonFile)
	if os.IsNotExist(err) || fsInfo.Size() == 0 {
		pnpmVersion, err := http.GetAllVersionsOfPnpm()
		if err != nil {
			logger.Warnf("Error fetching pnpm versions: %v", err)
			return []string{}
		}
		err = filemanagement.SavePnpmInfoToFilesystem(*pnpmVersion)
		if err != nil {
			logger.Warnf("Error saving pnpm versions to filesystem: %v", err)
			return []string{}
		}
		return *pnpmVersion
	}
	pnpmVersions, err := filemanagement.ReadPnpmInfoFromFilesystem()
	if err != nil {
		logger.Warnf("Error reading pnpm versions from filesystem: %v", err)
		return []string{}
	}
	return *pnpmVersions
}
