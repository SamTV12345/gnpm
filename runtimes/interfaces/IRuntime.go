package interfaces

import (
	"github.com/samtv12345/gnpm/models"
	"github.com/samtv12345/gnpm/packageJson"
)

type IRuntimeVersion interface {
	GetVersion() string
	IsLTS() bool
}

type IRuntime interface {
	GetRuntimeName() string
	GetRcFilename() string
	GetEngine(engine *packageJson.Engines) *string
	GetVersionedFilename(version string, filename string) string
	GetAllVersionsOfRuntime(forceDownload *bool) (*[]IRuntimeVersion, error)
	ToDownloadUrl(filenamePrefix string, shaSumOFFiles []models.CreateFilenameStruct, version string) (*string, error)
	GetFilenamePrefix(version string) string
	GetShaSumsForRuntime(version string) (*[]models.CreateFilenameStruct, error)
}
