package bun

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/samtv12345/gnpm/archive"
	"github.com/samtv12345/gnpm/filemanagement"
	http3 "github.com/samtv12345/gnpm/http"
	"github.com/samtv12345/gnpm/models"
	"github.com/samtv12345/gnpm/packageJson"
	http2 "github.com/samtv12345/gnpm/runtimes/impl/bun/http"
	"github.com/samtv12345/gnpm/runtimes/interfaces"
	"go.uber.org/zap"
)

type Runtime struct {
	Logger *zap.SugaredLogger
}

func (r Runtime) GetVersionedFilename(version string, filename string) string {
	filename = strings.Replace(filename, "bun", "bun-v"+version, 1)
	return filename
}

func (r Runtime) GetRuntimeName() string {
	return "bun"
}

func (r Runtime) ToDownloadUrl(filenamePrefix string, shaSumOFFiles []models.CreateFilenameStruct, version string) (*string, error) {
	filename := archive.FilterCorrectFilenameEnding(filenamePrefix, shaSumOFFiles)
	if filename == nil {
		r.Logger.Errorw("Error finding correct Node.js binary for your platform", "error")
		return nil, errors.New("error finding correct Node.js binary for your platform")
	}

	urlToNode := "https://github.com/oven-sh/bun/releases/download/bun-v" + version + "/" + filename.Filename

	return &urlToNode, nil
}

// GetFilenamePrefix bun does not use the version in the filename
func (r Runtime) GetFilenamePrefix(_ string) string {
	operatingSystem := runtime.GOOS
	architecture := runtime.GOARCH

	if strings.Contains(architecture, "amd") {
		architecture = strings.Replace(architecture, "amd", "x", 1)
	}

	return "bun-" + operatingSystem + "-" + architecture

}

func (r Runtime) GetShaSumsForRuntime(version string) (*[]models.CreateFilenameStruct, error) {
	response, err := http.Get("https://github.com/oven-sh/bun/releases/download/bun-v" + version + "/SHASUMS256.txt")
	if err != nil {
		r.Logger.Error("Error fetching SHASUMS256.txt:", err)
		return nil, err
	}
	defer response.Body.Close()
	shasumData, err := io.ReadAll(response.Body)
	if err != nil {
		r.Logger.Error("Error reading SHASUMS256.txt:", err)
		return nil, err
	}

	shasums := http3.DecodeShasumTxt(string(shasumData))
	for i, shasum := range shasums {
		if strings.Contains(shasum.Filename, "profile") {
			shasums = append(shasums[:i], shasums[i+1:]...)
			continue
		}
	}
	return &shasums, nil
}

func (r Runtime) GetInformationFromPackageJSON(proposedVersion *string, path string, versions *[]interfaces.IRuntimeVersion) (*interfaces.IRuntimeVersion, error) {
	var versionToDownload *string
	packageManifest, err := packageJson.ReadPackageJson(filepath.Join(path, "package.json"))
	if err != nil {
		return nil, err
	}
	// Check .bun-version first
	bunvmRC, errNvmrc := packageJson.ReadRuntimeVersionFile(filepath.Join(path, ".bun-version"))
	if errNvmrc == nil {
		versionToDownload = &bunvmRC
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
		return nil, errors.New("error finding a suitable bun version")
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
	nodeJSCacheFile := filepath.Join(*dataDir, ".cache", "bun_index.json")
	fsInfo, err := os.Stat(nodeJSCacheFile)
	if os.IsNotExist(err) || fsInfo.Size() == 0 {
		nodeVersions, err := http2.GetBunVersions(r.Logger)
		if err != nil {
			return nil, err
		}
		err = filemanagement.SaveBunInfoToFilesystem(*nodeVersions)
		if err != nil {
			return nil, err
		}

		converted := make([]interfaces.IRuntimeVersion, len(*nodeVersions))
		for i, v := range *nodeVersions {
			converted[i] = v
		}
		return &converted, nil
	}

	bunVersions, err := filemanagement.ReadBunInfoFromFilesystem()
	if err != nil {
		return nil, err
	}
	return bunVersions, nil
}

var _ interfaces.IRuntime = (*Runtime)(nil)
