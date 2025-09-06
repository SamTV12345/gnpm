package interfaces

import (
	"github.com/samtv12345/gnpm/models"
)

type IRuntimeVersion interface {
	GetVersion() string
	IsLTS() bool
}

type IRuntime interface {
	GetRuntimeName() string
	GetAllVersionsOfRuntime() (*[]IRuntimeVersion, error)
	GetInformationFromPackageJSON(proposedVersion *string, path string, versions *[]IRuntimeVersion) (*IRuntimeVersion, error)
	ToDownloadUrl(filenamePrefix string, shaSumOFFiles []models.CreateFilenameStruct, version string) (*string, error)
	GetFilenamePrefix(version string) string
	GetShaSumsForRuntime(version string) (*[]models.CreateFilenameStruct, error)
}
