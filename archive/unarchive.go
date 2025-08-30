package archive

import (
	"path/filepath"

	"go.uber.org/zap"
)

func UnarchiveFile(path string, logger *zap.SugaredLogger) (*string, error) {
	var extension = filepath.Ext(path)
	if extension == ".zip" {
		return unzip(path, logger)
	} else if extension == ".gz" {
		return untar(path, logger)
	} else {
		return nil, nil
	}
}
