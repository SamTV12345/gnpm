package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/samtv12345/gnpm/filemanagement"
)

func StripFirstDir(path string) string {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return path
}

func unzip(path string) (*string, error) {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer archive.Close()
	targetPath, err := filemanagement.CreateTargetDir(path)
	if err != nil {
		return nil, err
	}
	for _, f := range archive.File {
		relPath := StripFirstDir(f.Name)

		if relPath == "" {
			continue
		}

		filePath := filepath.Join(*targetPath, relPath)

		if !strings.HasPrefix(filePath, filepath.Clean(*targetPath)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return nil, fmt.Errorf("invalid file path: %s", filePath)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return nil, err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return nil, err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return nil, err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return targetPath, nil
}
