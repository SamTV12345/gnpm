package filemanagement

import (
	"os"
	"path/filepath"
	"strings"
)

func stripSuffix(filename string) string {
	filename = strings.Replace(filename, ".zip", "", 1)
	filename = strings.Replace(filename, ".tar.gz", "", 1)
	filename = strings.Replace(filename, ".tgz", "", 1)
	return filename
}

func CreateTargetDir(sourcefile string) (*string, error) {
	filename := filepath.Base(sourcefile)
	targetDir := filepath.Join(filepath.Dir(sourcefile), "..", "_gnpm", stripSuffix(filename))
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	return &targetDir, nil
}

func DoesTargetDirExist(sourcefile string) (*string, error) {
	filename := filepath.Base(sourcefile)
	println(stripSuffix(filename))
	targetDir := filepath.Join(filepath.Dir(sourcefile), "..", "_gnpm", stripSuffix(filename))
	return &targetDir, nil
}
