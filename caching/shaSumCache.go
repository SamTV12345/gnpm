package caching

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/samtv12345/gnpm/filemanagement"
	"github.com/samtv12345/gnpm/http"
	"github.com/samtv12345/gnpm/models"
	"github.com/samtv12345/gnpm/runtimes/interfaces"
)

func GetShaSumCacheInPath(runtime interfaces.IRuntime, version string) (*[]models.CreateFilenameStruct, error) {
	var dataDir, err = filemanagement.EnsureDataDir()
	if err != nil {
		return nil, err
	}
	var file = filepath.Join(*dataDir, ".cache", fmt.Sprintf("shaSumCache_%s_%s.json", runtime.GetRuntimeName(), version))
	if _, err := os.Stat(file); os.IsNotExist(err) {
		shaData, err := runtime.GetShaSumsForRuntime(version)
		if err != nil {
			return nil, err
		}
		err = filemanagement.SaveShaSumInfoToFilesystem(*shaData, file)
		if err != nil {
			return nil, err
		}
		return shaData, nil
	} else {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		createFilenames := http.DecodeShasumTxt(string(content))
		return &createFilenames, nil
	}
}
