package caching

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/samtv12345/gnpm/filemanagement"
	"github.com/samtv12345/gnpm/pm/interfaces"
	"go.uber.org/zap"
)

func GetPnpmVersion(logger *zap.SugaredLogger, pm interfaces.IPackageManager) []string {
	cacheDir, err := filemanagement.GetCacheDir()
	if err != nil {
		logger.Warnf("Error getting cache dir: %v", err)
		return []string{}
	}
	pnpmJsonFile := filepath.Join(*cacheDir, pm.GetName()+".json")

	fsInfo, err := os.Stat(pnpmJsonFile)
	if os.IsNotExist(err) || fsInfo.Size() == 0 {
		pnpmVersion, err := pm.GetAllVersions()
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
	content, err := os.ReadFile(pnpmJsonFile)
	if err != nil {
		return []string{}
	}
	var pnpmVersions []string
	if err := json.Unmarshal(content, &pnpmVersions); err != nil {
		return []string{}
	}
	return pnpmVersions
}
