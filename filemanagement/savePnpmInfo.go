package filemanagement

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func SavePnpmInfoToFilesystem(pnpmVersions []string) error {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return err
	}
	pnpmJsonFile := filepath.Join(*cacheDir, "pnpm.json")
	marshalledPnpmVersions, err := json.Marshal(pnpmVersions)
	if err != nil {
		return err
	}
	return os.WriteFile(pnpmJsonFile, marshalledPnpmVersions, os.ModePerm)
}

func ReadPnpmInfoFromFilesystem() (*[]string, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return nil, err
	}
	pnpmJsonFile := filepath.Join(*cacheDir, "pnpm.json")
	fileBytes, err := os.ReadFile(pnpmJsonFile)
	if err != nil {
		return nil, err
	}
	var pnpmVersions []string
	err = json.Unmarshal(fileBytes, &pnpmVersions)
	if err != nil {
		return nil, err
	}
	return &pnpmVersions, nil
}
