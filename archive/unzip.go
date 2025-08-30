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

func Unzip(path string) (*string, error) {
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
		filePath := filepath.Join(*targetPath, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(*targetPath)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return nil, fmt.Errorf("invalid file path: %s", filePath)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return targetPath, nil
}
