package filemanagement

import (
	"os"
	"path/filepath"
)

func CreateTargetDir(sourcefile string) (*string, error) {
	targetDir := filepath.Join(filepath.Dir(sourcefile), "..", "_gnpm")
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	return &targetDir, nil
}
