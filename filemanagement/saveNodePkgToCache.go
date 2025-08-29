package filemanagement

import (
	"os"
	"path/filepath"

	"github.com/samtv12345/gnpm/models"
	"go.uber.org/zap"
)

func SaveNodeJSToCacheDir(nodeData []byte, createNodeDat models.CreateNodeDownloadStruct, logger *zap.SugaredLogger) (*string, error) {
	dataDir, err := EnsureDataDir()
	if err != nil {
		return nil, err
	}
	cacheDir := filepath.Join(*dataDir, ".cache")
	err = os.Mkdir(cacheDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	locationToWriteNodeJSArchive := filepath.Join(cacheDir, createNodeDat.Filename)
	logger.Debugf("Saving NodeJS to cache at location: %s", locationToWriteNodeJSArchive)
	err = os.WriteFile(locationToWriteNodeJSArchive, nodeData, 0644)
	if err != nil {
		return nil, err
	}
	return &locationToWriteNodeJSArchive, nil
}

func HasNodeVersionInCache(downloadStruct *models.CreateNodeDownloadStruct, logger *zap.SugaredLogger) (*bool, error) {
	dataDir, err := EnsureDataDir()
	if err != nil {
		return nil, err
	}
	cacheFileOfNodeJs := filepath.Join(*dataDir, ".cache", downloadStruct.Filename)
	logger.Debugf("Checking if Node.js exists at: %s", cacheFileOfNodeJs)
	_, err = os.Stat(cacheFileOfNodeJs)

	if os.IsNotExist(err) {
		var falseVal = false
		return &falseVal, nil
	}
	if err != nil {
		return nil, err
	}
	return &[]bool{true}[0], nil
}
