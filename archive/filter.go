package archive

import (
	"path/filepath"
	"strings"

	"github.com/samtv12345/gnpm/models"
)

func FilterCorrectFilenameEnding(filenamePrefix string, shaSumsOFFiles []models.CreateFilenameStruct) *models.CreateDownloadStruct {
	for _, file := range shaSumsOFFiles {
		if strings.HasPrefix(file.Filename, filenamePrefix) && (strings.HasSuffix(file.Filename, ".tar.gz") || strings.HasSuffix(file.Filename, ".zip")) {
			fileExtension := filepath.Ext(file.Filename)
			return &models.CreateDownloadStruct{
				Sha256:   file.Sha256,
				Encoding: fileExtension,
			}
		}
	}
	return nil
}
