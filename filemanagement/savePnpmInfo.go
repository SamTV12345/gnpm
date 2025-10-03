package filemanagement

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/samtv12345/gnpm/pm/interfaces"
)

func SavePnpmInfoToFilesystem(pnpmVersions []string, pm interfaces.IPackageManager) error {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return err
	}
	pnpmJsonFile := filepath.Join(*cacheDir, pm.GetName()+".json")
	marshalledPnpmVersions, err := json.Marshal(pnpmVersions)
	if err != nil {
		return err
	}
	return os.WriteFile(pnpmJsonFile, marshalledPnpmVersions, os.ModePerm)
}
