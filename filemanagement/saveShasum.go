package filemanagement

import (
	"os"

	"github.com/samtv12345/gnpm/models"
)

func SaveShaSumInfoToFilesystem(shasums []models.CreateFilenameStruct, path string) error {
	contentToWrite := ""
	for _, shasum := range shasums {
		contentToWrite += shasum.Sha256 + "  " + shasum.Filename + "\n"
	}
	return os.WriteFile(path, []byte(contentToWrite), 0644)
}
