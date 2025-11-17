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
	http3 "github.com/samtv12345/gnpm/http"
	"github.com/samtv12345/gnpm/models"
	"github.com/samtv12345/gnpm/packageJson"
	"github.com/samtv12345/gnpm/runtimes/impl/node/http"
	"github.com/samtv12345/gnpm/runtimes/interfaces"
	"github.com/samtv12345/gnpm/utils"
	"go.uber.org/zap"
)

type Runtime struct {
	Logger *zap.SugaredLogger
}

func (n Runtime) GetRcFilename() string {
	return ".nvmrc"
}

func (n Runtime) GetEngine(engine *packageJson.Engines) *string {
	if engine != nil {
		return engine.Node
	}
	return nil
}

func (n Runtime) GetVersionedFilename(_ string, filename string) string {
	return filename
}

func (n Runtime) GetRuntimeName() string {
	return "node"
}

func (n Runtime) GetInformationFromPackageJSON(proposedVersion *string, path string, nodeVersions *[]interfaces.IRuntimeVersion) (*interfaces.IRuntimeVersion, error) {
	var versionToDownload *string

	if proposedVersion != nil {
		versionToDownload = proposedVersion
	} else {
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
	}

	if versionToDownload == nil {
		// Check parent directory for .nvmrc and package.json engines field
		parentPath := filepath.Join(path, "..")
		packageJsonPath := filepath.Join(parentPath, "package.json")
		if _, err := os.Stat(packageJsonPath); err == nil {
			return n.GetInformationFromPackageJSON(proposedVersion, parentPath, nodeVersions)
		}
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

func (n Runtime) GetAllVersionsOfRuntime(forceInstall *bool) (*[]interfaces.IRuntimeVersion, error) {
	cacheDir, err := filemanagement.GetCacheDir()
	if err != nil {
		return nil, err
	}
	nodeJSCacheFile := filepath.Join(*cacheDir, ".cache", "nodejs_index.json")
	fsInfo, err := os.Stat(nodeJSCacheFile)
	if os.IsNotExist(err) || fsInfo.Size() == 0 || (forceInstall != nil && *forceInstall) {
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

	var urlToNode string
	if utils.IsMusl() {
		urlToNode = "https://unofficial-builds.nodejs.org/download/release/" + version + "/" + filename.Filename
		n.Logger.Info("Musl detected, using unofficial builds:", urlToNode)
	} else {
		urlToNode = "https://nodejs.org/dist/" + version + "/" + filename.Filename

	}

	return &urlToNode, nil
}

func (n Runtime) GetShaSumsForRuntime(version string) (*[]models.CreateFilenameStruct, error) {
	var response *http2.Response
	var err error
	if utils.IsMusl() {
		response, err = http2.Get("https://unofficial-builds.nodejs.org/download/release/" + version + "/SHASUMS256.txt")
	} else {
		response, err = http2.Get("https://nodejs.org/dist/" + version + "/SHASUMS256.txt")
	}

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

	shasums := http3.DecodeShasumTxt(string(shasumData))
	return &shasums, nil
}

var _ interfaces.IRuntime = (*Runtime)(nil)
