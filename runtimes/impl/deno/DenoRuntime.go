package deno

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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

func (r Runtime) GetRcFilename() string {
	return ".deno-version"
}

func (r Runtime) GetEngine(engine *packageJson.Engines) *string {
	if engine != nil {
		return engine.Deno
	}
	return nil
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

func (r Runtime) GetAllVersionsOfRuntime(forceInstall *bool) (*[]interfaces.IRuntimeVersion, error) {
	dataDir, err := filemanagement.EnsureDataDir()
	if err != nil {
		return nil, err
	}
	nodeJSCacheFile := filepath.Join(*dataDir, ".cache", "deno_index.json")
	fsInfo, err := os.Stat(nodeJSCacheFile)
	if os.IsNotExist(err) || fsInfo.Size() == 0 || (forceInstall != nil && *forceInstall) {
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
