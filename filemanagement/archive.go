package filemanagement

import (
	"os"
	"path/filepath"
)

func HasArchiveBeenExtracted(filename string) bool {
	extractedFlagFile := filepath.Join(filename)
	_, err := os.Stat(extractedFlagFile)
	if err != nil {
		return false
	}
	return true
}
