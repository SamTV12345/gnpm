package node

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	http2 "net/http"

	"github.com/Masterminds/semver/v3"
	"github.com/samtv12345/gnpm/archive"
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

func (n Runtime) GetInformationFromPackageJSON(proposedVersion *string, path string, nodeVersions *[]interfaces.IRuntimeVersion) (*interfaces.IRuntimeVersion, error) {
	var versionToDownload *string
	packageManifest, err := packageJson.ReadPackageJson(filepath.Join(path, "package.json"))
	if err != nil {
		return nil, err
	}

	// Check .nvmrc first
	nvmrcVersion, errNvmrc := packageJson.ReadRuntimeVersionFile(filepath.Join(path, ".nvmrc"))
	if errNvmrc == nil {
		versionToDownload = &nvmrcVersion
	}

	// Then check the package.json "engines" field
	if packageManifest != nil && packageManifest.Engines != nil && packageManifest.Engines.Node != nil {
		versionToDownload = packageManifest.Engines.Node
	}

	if proposedVersion != nil {
		versionToDownload = proposedVersion
	}

	if versionToDownload != nil {
		constraints, err := semver.NewConstraint(*versionToDownload)
		if err != nil {
			n.Logger.Errorw("Error parsing version", "error", err)
		}

		var possibleVersions = make([]*semver.Version, 0)
		for _, nodeVersion := range *nodeVersions {
			v, err := semver.NewVersion(nodeVersion.GetVersion()[1:])
			if err != nil {
				n.Logger.Errorw("Error parsing version", "error", err)
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
		for _, nodeVersion := range *nodeVersions {
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

	for _, nodeVersion := range *nodeVersions {
		if nodeVersion.GetVersion() == *versionToDownload || nodeVersion.GetVersion() == "v"+*versionToDownload {
			return &nodeVersion, nil
		}
	}
	return nil, errors.New("error finding a suitable Node.js version")
}

func (n Runtime) GetAllVersionsOfRuntime() (*[]interfaces.IRuntimeVersion, error) {
	dataDir, err := filemanagement.EnsureDataDir()
	if err != nil {
		return nil, err
	}
	nodeJSCacheFile := filepath.Join(*dataDir, ".cache", "nodejs_index.json")
	fsInfo, err := os.Stat(nodeJSCacheFile)
	if os.IsNotExist(err) || fsInfo.Size() == 0 {
		nodeVersions, err := http.GetNodeJsVersion(n.Logger)
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

func (n Runtime) GetFilenamePrefix(version string) string {
	operatingSystem := runtime.GOOS
	architecture := runtime.GOARCH

	if operatingSystem == "windows" {
		operatingSystem = "win"
	}

	if strings.Contains(architecture, "amd") {
		architecture = strings.Replace(architecture, "amd", "x", 1)
	}

	var filenamePrefix = "node-" + version + "-" + operatingSystem + "-" + architecture
	return filenamePrefix
}

func (n Runtime) ToDownloadUrl(filenamePrefix string, shaSumOFFiles []models.CreateFilenameStruct, version string) (*string, error) {
	filename := archive.FilterCorrectFilenameEnding(filenamePrefix, shaSumOFFiles)
	if filename == nil {
		n.Logger.Errorw("Error finding correct Node.js binary for your platform", "error")
		return nil, errors.New("error finding correct Node.js binary for your platform")
	}

	urlToNode := "https://nodejs.org/dist/" + version + "/" + filename.Filename

	return &urlToNode, nil
}

func (n Runtime) GetShaSumsForRuntime(version string) (*[]models.CreateFilenameStruct, error) {
	response, err := http2.Get("https://nodejs.org/dist/" + version + "/SHASUMS256.txt")
	if err != nil {
		n.Logger.Error("Error fetching SHASUMS256.txt:", err)
		return nil, err
	}
	defer response.Body.Close()
	shasumData, err := io.ReadAll(response.Body)
	if err != nil {
		n.Logger.Error("Error reading SHASUMS256.txt:", err)
		return nil, err
	}

	var shaSumData = string(shasumData)
	var shaSumToFileMappingArr = make([]models.CreateFilenameStruct, 0)
	splittedShaSumData := strings.Split(shaSumData, "\n")
	for _, line := range splittedShaSumData {
		shaSumToFileMapping := strings.Split(line, "  ")
		if len(shaSumToFileMapping) == 2 {
			shaSumToFileMappingArr = append(shaSumToFileMappingArr, models.CreateFilenameStruct{
				Sha256:   shaSumToFileMapping[0],
				Filename: shaSumToFileMapping[1],
			})
		}
	}
	return &shaSumToFileMappingArr, nil
}

var _ interfaces.IRuntime = (*Runtime)(nil)
