package archive

import "path/filepath"

func UnarchiveFile(path string) (*string, error) {
	var extension = filepath.Ext(path)
	if extension == ".zip" {
		return unzip(path)
	} else if extension == ".gz" {
		return untar(path)
	} else {
		return nil, nil
	}
}
