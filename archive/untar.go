package archive

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/samtv12345/gnpm/filemanagement"
	"go.uber.org/zap"
)

func splitAfterNthOccurenceOfSign(s string, sep rune, n int) (string, string) {
	for i, sep2 := range s {
		if sep2 == sep {
			n--
			if n == 0 {
				return s[:i], s[i+1:]
			}
		}
	}
	return s, ""
}

func IsGzipTar(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	magic := make([]byte, 2)
	_, err = f.Read(magic)
	if err != nil {
		return false, err
	}
	return magic[0] == 0x1F && magic[1] == 0x8B, nil
}

func untar(path string, logger *zap.SugaredLogger) (*string, error) {
	tarballArchive, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var readerToUse io.Reader

	readerToUse = bytes.NewReader(tarballArchive)

	isGzipTar, err := IsGzipTar(path)
	if err != nil {
		return nil, err
	}

	if isGzipTar {
		logger.Debugf("Extracting tarball with gzip from %s", path)
		gzipReader, err := gzip.NewReader(readerToUse)
		defer gzipReader.Close()
		if err != nil {
			return nil, err
		}
		readerToUse = gzipReader
	}

	logger.Debugf("Extracting tarball from %s", path)

	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(readerToUse)
	targetPath, err := filemanagement.CreateTargetDir(path)
	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return targetPath, nil

		// return any other error
		case err != nil:
			return nil, err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		_, strippedHeader := splitAfterNthOccurenceOfSign(header.Name, '/', 1)

		// skip because tar always contains a folder with the same name as the tarball
		if strippedHeader == "" {
			continue
		}

		target := filepath.Join(*targetPath, strippedHeader)
		// If it is a file create dir to file
		if !header.FileInfo().IsDir() {
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
					return nil, err
				}
			}
		}

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()
		// check the file type
		switch header.Typeflag {

		// if it's a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return nil, err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return nil, err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return nil, err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}
