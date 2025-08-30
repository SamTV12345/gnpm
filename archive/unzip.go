package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	sync2 "sync"

	"github.com/samtv12345/gnpm/filemanagement"
	"go.uber.org/zap"
)

func StripFirstDir(path string) string {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return path
}

func unzip(path string, logger *zap.SugaredLogger) (*string, error) {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer archive.Close()
	targetPath, err := filemanagement.CreateTargetDir(path)
	if err != nil {
		return nil, err
	}
	var sync sync2.WaitGroup
	for _, f := range archive.File {
		sync.Add(1)
		go func() {
			defer sync.Done()
			relPath := StripFirstDir(f.Name)

			if relPath == "" {
				return
			}

			filePath := filepath.Join(*targetPath, relPath)

			if !strings.HasPrefix(filePath, filepath.Clean(*targetPath)+string(os.PathSeparator)) {
				fmt.Println("invalid file path")
				logger.Error("invalid file path: ", filePath)
				return
			}
			if f.FileInfo().IsDir() {
				os.MkdirAll(filePath, os.ModePerm)
				return
			}

			if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				logger.Errorw("failed to create directory", "path", filePath, "error", err)
				return
			}

			dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				logger.Errorw("failed to open file", "path", filePath, "error", err)
				return
			}

			fileInArchive, err := f.Open()
			if err != nil {
				logger.Errorw("failed to open file", "path", filePath, "error", err)
				return
			}

			if _, err := io.Copy(dstFile, fileInArchive); err != nil {
				panic(err)
			}

			dstFile.Close()
			fileInArchive.Close()
		}()
	}
	sync.Wait()
	return targetPath, nil
}
