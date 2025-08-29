package caching

import (
	"os"
	"path/filepath"

	"github.com/samtv12345/gnpm/filemanagement"
	"github.com/samtv12345/gnpm/http"
	"go.uber.org/zap"
)

func GetNodeJsVersion(logger *zap.SugaredLogger) (*[]http.NodeIndex, error) {
	dataDir, err := filemanagement.EnsureDataDir()
	if err != nil {
		return nil, err
	}
	nodeJSCacheFile := filepath.Join(*dataDir, ".cache", "nodejs_index.json")
	fsInfo, err := os.Stat(nodeJSCacheFile)
	if os.IsNotExist(err) || fsInfo.Size() == 0 {
		nodeVersions, err := http.GetNodeJsVersion(logger)
		if err != nil {
			return nil, err
		}
		err = filemanagement.SaveNodeInfoToFilesystem(*nodeVersions)
		if err != nil {
			return nil, err
		}
		return nodeVersions, nil
	}

	nodeVersions, err := filemanagement.ReadNodeInfoFromFilesystem()
	if err != nil {
		return nil, err
	}
	return nodeVersions, nil
}
