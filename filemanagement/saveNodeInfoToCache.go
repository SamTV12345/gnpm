package filemanagement

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/samtv12345/gnpm/http"
)

func SaveNodeInfoToFilesystem(nodeIndices []http.NodeIndex) error {
	jsonBytes, err := json.Marshal(nodeIndices)
	if err != nil {
		return err
	}
	cacheDir, err := getCacheDir()
	if err != nil {
		return err
	}
	pathToSaveTo := filepath.Join(*cacheDir, "nodejs_index.json")
	return os.WriteFile(pathToSaveTo, jsonBytes, os.ModePerm)
}

func getCacheDir() (*string, error) {
	dataDir, err := EnsureDataDir()
	if err != nil {
		return nil, err
	}
	cacheDir := filepath.Join(*dataDir, ".cache")
	err = os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	return &cacheDir, nil
}

func ReadNodeInfoFromFilesystem() (*[]http.NodeIndex, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return nil, err
	}
	pathToReadFrom := filepath.Join(*cacheDir, "nodejs_index.json")
	fileBytes, err := os.ReadFile(pathToReadFrom)
	if err != nil {
		return nil, err
	}
	var nodeIndices []http.NodeIndex
	err = json.Unmarshal(fileBytes, &nodeIndices)
	if err != nil {
		return nil, err
	}
	return &nodeIndices, nil
}
