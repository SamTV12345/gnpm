package bun

import (
	"github.com/samtv12345/gnpm/models"
	"github.com/samtv12345/gnpm/runtimes/interfaces"
)

type Runtime struct {
}

func (r Runtime) ToDownloadUrl(filenamePrefix string, shaSumOFFiles []models.CreateFilenameStruct, version string) (*string, error) {
	//TODO implement me
	panic("implement me")
}

func (r Runtime) GetFilenamePrefix(version string) string {
	//TODO implement me
	panic("implement me")
}

func (r Runtime) GetShaSumsForRuntime(version string) (*[]models.CreateFilenameStruct, error) {
	//TODO implement me
	panic("implement me")
}

func (r Runtime) GetInformationFromPackageJSON(proposedVersion *string, path string, versions *[]interfaces.IRuntimeVersion) (*interfaces.IRuntimeVersion, error) {
	//TODO implement me
	panic("implement me")
}

func (r Runtime) GetAllVersionsOfRuntime() (*[]interfaces.IRuntimeVersion, error) {
	//TODO implement me
	panic("implement me")
}

var _ interfaces.IRuntime = (*Runtime)(nil)
