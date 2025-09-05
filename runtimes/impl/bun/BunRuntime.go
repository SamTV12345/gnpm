package bun

import (
	"errors"
	"os"
	"path/filepath"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/samtv12345/gnpm/filemanagement"
	"github.com/samtv12345/gnpm/models"
	"github.com/samtv12345/gnpm/packageJson"
	"github.com/samtv12345/gnpm/runtimes/impl/node/http"
	"github.com/samtv12345/gnpm/runtimes/interfaces"
	"go.uber.org/zap"
)

type Runtime struct {
	Logger *zap.SugaredLogger
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
	var versionToDownload *string
	packageManifest, err := packageJson.ReadPackageJson(filepath.Join(path, "package.json"))
	if err != nil {
		return nil, err
	}
	// Check .nvmrc first
	nvmrcVersion, errNvmrc := packageJson.ReadRuntimeVersionFile(filepath.Join(path, ".bun-version"))
	if errNvmrc == nil {
		versionToDownload = &nvmrcVersion
	}

	// Then check the package.json "engines" field
	if packageManifest != nil && packageManifest.Engines != nil && packageManifest.Engines.Bun != nil {
		versionToDownload = packageManifest.Engines.Bun
	}

	if proposedVersion != nil {
		versionToDownload = proposedVersion
	}

	if versionToDownload != nil {
		constraints, err := semver.NewConstraint(*versionToDownload)
		if err != nil {
			r.Logger.Errorw("Error parsing version", "error", err)
		}

		var possibleVersions = make([]*semver.Version, 0)
		for _, nodeVersion := range *versions {
			v, err := semver.NewVersion(nodeVersion.GetVersion()[1:])
			if err != nil {
				r.Logger.Errorw("Error parsing version", "error", err)
				continue
			}
			if constraints.Check(v) {
				possibleVersions = append(possibleVersions, v)
			}
		}

		// Sort possibleVersions in descending order
		sort.Sort(semver.Collection(possibleVersions))

		if len(possibleVersions) > 0 {
			latestVersion := possibleVersions[len(possibleVersions)-1].String()
			versionToDownload = &latestVersion
		}
	} else {
		// Default to the latest LTS version if no version is specified
		for _, nodeVersion := range *versions {
			if nodeVersion.IsLTS() != false {
				var nodeVersion = nodeVersion.GetVersion()
				versionToDownload = &nodeVersion
				break
			}
		}
	}

	if versionToDownload == nil {
		return nil, errors.New("error finding a suitable Node.js version")
	}

	for _, nodeVersion := range *versions {
		if nodeVersion.GetVersion() == *versionToDownload || nodeVersion.GetVersion() == "v"+*versionToDownload {
			return &nodeVersion, nil
		}
	}
	return nil, errors.New("error finding a suitable Node.js version")
}

func (r Runtime) GetAllVersionsOfRuntime() (*[]interfaces.IRuntimeVersion, error) {
	dataDir, err := filemanagement.EnsureDataDir()
	if err != nil {
		return nil, err
	}
	nodeJSCacheFile := filepath.Join(*dataDir, ".cache", "nodejs_index.json")
	fsInfo, err := os.Stat(nodeJSCacheFile)
	if os.IsNotExist(err) || fsInfo.Size() == 0 {
		nodeVersions, err := http.GetNodeJsVersion(r.Logger)
		if err != nil {
			return nil, err
		}
		err = filemanagement.SaveNodeInfoToFilesystem(*nodeVersions)
		if err != nil {
			return nil, err
		}

		converted := make([]interfaces.IRuntimeVersion, len(*nodeVersions))
		for i, v := range *nodeVersions {
			converted[i] = v
		}
		return &converted, nil
	}

	nodeVersions, err := filemanagement.ReadNodeInfoFromFilesystem()
	if err != nil {
		return nil, err
	}
	return nodeVersions, nil
}

var _ interfaces.IRuntime = (*Runtime)(nil)
