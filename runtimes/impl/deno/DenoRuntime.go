package deno

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/samtv12345/gnpm/archive"
	"github.com/samtv12345/gnpm/filemanagement"
	"github.com/samtv12345/gnpm/models"
	"github.com/samtv12345/gnpm/packageJson"
	http2 "github.com/samtv12345/gnpm/runtimes/impl/deno/http"
	"github.com/samtv12345/gnpm/runtimes/interfaces"
	"go.uber.org/zap"
)

type Runtime struct {
	Logger *zap.SugaredLogger
}

func (r Runtime) GetVersionedFilename(version string, filename string) string {
	filename = strings.Replace(filename, "deno", "deno-v"+version, 1)
	return filename
}

func (r Runtime) GetRuntimeName() string {
	return "deno"
}

func (r Runtime) ToDownloadUrl(filenamePrefix string, shaSumOFFiles []models.CreateFilenameStruct, version string) (*string, error) {
	filename := archive.FilterCorrectFilenameEnding(filenamePrefix, shaSumOFFiles)
	if filename == nil {
		r.Logger.Errorw("Error finding correct Node.js binary for your platform", "error")
		return nil, errors.New("error finding correct Node.js binary for your platform")
	}

	urlToNode := "https://github.com/denoland/deno/releases/download/v" + version + "/" + filename.Filename

	return &urlToNode, nil
}

func (r Runtime) GetFilenamePrefix(_ string) string {
	operatingSystem := runtime.GOOS
	architecture := runtime.GOARCH

	if strings.Contains(architecture, "amd64") {
		architecture = strings.Replace(architecture, "amd64", "x86_64", 1)
	}

	if operatingSystem == "windows" {
		operatingSystem = "pc-windows-msvc"
	} else if operatingSystem == "darwin" {
		operatingSystem = "apple-darwin"
	} else {
		operatingSystem = "unknown-linux-gnu"
	}

	return "deno-" + architecture + "-" + operatingSystem
}

func (r Runtime) GetShaSumsForRuntime(version string) (*[]models.CreateFilenameStruct, error) {
	releaseData, err := http2.GetDenoGitHubRelease(version, r.Logger)

	if err != nil {
		return nil, err
	}

	shaSums := make([]models.CreateFilenameStruct, 0)
	for _, asset := range releaseData.Assets {
		if strings.HasSuffix(asset.Name, ".zip") {
			shaSums = append(shaSums, models.CreateFilenameStruct{
				Filename: asset.Name,
				Sha256:   strings.TrimPrefix(asset.Digest, "sha256:"),
			})
		}
	}
	return &shaSums, nil
}

func (r Runtime) GetInformationFromPackageJSON(proposedVersion *string, path string, versions *[]interfaces.IRuntimeVersion) (*interfaces.IRuntimeVersion, error) {
	var versionToDownload *string
	packageManifest, err := packageJson.ReadPackageJson(filepath.Join(path, "package.json"))
	if err != nil {
		return nil, err
	}
	// Check .deno-version first
	bunvmRC, errNvmrc := packageJson.ReadRuntimeVersionFile(filepath.Join(path, ".deno-version"))
	if errNvmrc == nil {
		versionToDownload = &bunvmRC
	}

	// Then check the package.json "engines" field
	if packageManifest != nil && packageManifest.Engines != nil && packageManifest.Engines.Deno != nil {
		versionToDownload = packageManifest.Engines.Deno
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
			v, err := semver.NewVersion(nodeVersion.GetVersion())
			if err != nil {
				r.Logger.Errorw("Error parsing version", "error", err, nodeVersion.GetVersion())
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
		return nil, errors.New("error finding a suitable gnpm version")
	}

	for _, nodeVersion := range *versions {
		if nodeVersion.GetVersion() == *versionToDownload || nodeVersion.GetVersion() == "v"+*versionToDownload {
			return &nodeVersion, nil
		}
	}
	return nil, errors.New("error finding a suitable Bun version")
}

func (r Runtime) GetAllVersionsOfRuntime() (*[]interfaces.IRuntimeVersion, error) {
	dataDir, err := filemanagement.EnsureDataDir()
	if err != nil {
		return nil, err
	}
	nodeJSCacheFile := filepath.Join(*dataDir, ".cache", "deno_index.json")
	fsInfo, err := os.Stat(nodeJSCacheFile)
	if os.IsNotExist(err) || fsInfo.Size() == 0 {
		nodeVersions, err := http2.GetDenoVersions(r.Logger)
		if err != nil {
			return nil, err
		}
		err = filemanagement.SaveDenoInfoToFilesystem(*nodeVersions)
		if err != nil {
			return nil, err
		}

		converted := make([]interfaces.IRuntimeVersion, len(*nodeVersions))
		for i, v := range *nodeVersions {
			converted[i] = v
		}
		return &converted, nil
	}

	denoVersion, err := filemanagement.ReadDenoInfoFromFilesystem()
	if err != nil {
		return nil, err
	}
	return denoVersion, nil
}

var _ interfaces.IRuntime = (*Runtime)(nil)
