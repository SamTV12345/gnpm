package gnpm

import (
	"errors"
	"path/filepath"
	"runtime"

	"github.com/samtv12345/gnpm/archive"
	"github.com/samtv12345/gnpm/caching"
	"github.com/samtv12345/gnpm/commandRun"
	"github.com/samtv12345/gnpm/filemanagement"
	"github.com/samtv12345/gnpm/http"
	"github.com/samtv12345/gnpm/models"
	"github.com/samtv12345/gnpm/runtimes"
	"github.com/samtv12345/gnpm/runtimes/interfaces"
	"go.uber.org/zap"
)

func createRelevantRuntimePaths(targetPath string, selectedRuntime interfaces.IRuntime) []string {
	if runtime.GOOS == "windows" {
		runtimePath := filepath.Join(targetPath, selectedRuntime.GetRuntimeName()+".exe")
		return []string{runtimePath}
	}

	if runtime.GOOS == "linux" {
		runtimePath := filepath.Join(targetPath, "bin", selectedRuntime.GetRuntimeName())
		return []string{runtimePath}
	}

	var runtimePath = filepath.Join(targetPath, selectedRuntime.GetRuntimeName())
	return []string{runtimePath}
}

func HandleRuntimeVersion(args commandRun.FlagArguments, logger *zap.SugaredLogger) (relevantPathsToReturn *[]string, selectedRuntimeFor *interfaces.IRuntime, err error) {
	var selectedRuntime = runtimes.GetRuntimeSelection(logger)
	runtimeVersions, err := selectedRuntime.GetAllVersionsOfRuntime()

	if err != nil {
		logger.Errorf("Error fetching %s versions with cause %s", selectedRuntime.GetRuntimeName(), err)
		return nil, nil, err
	}

	// Parse runtime version from e.g. .nvmrc or package.json
	runtimeVersionToDownload, err := selectedRuntime.GetInformationFromPackageJSON(args.RuntimeVersion, ".", runtimeVersions)
	if err != nil {
		logger.Errorf("Error determining  version %s with cause %s %s", selectedRuntime.GetRuntimeName(), "error", err)
		return nil, nil, err
	}
	logger.Infof("%s version to use: %s", selectedRuntime.GetRuntimeName(), (*runtimeVersionToDownload).GetVersion())
	createRuntimeDownloadUrlInfo, err := createDownloadUrl(*runtimeVersionToDownload, selectedRuntime, logger)
	if err != nil {
		logger.Errorf("Error creating %s download URL with %s %s", selectedRuntime.GetRuntimeName(), "error", err)
		return nil, nil, err
	}
	logger.Debugf("%s download URL: %s", selectedRuntime.GetRuntimeName(), createRuntimeDownloadUrlInfo.RuntimeUrl)
	exists, filename, err := filemanagement.HasRuntimeVersionInCache(createRuntimeDownloadUrlInfo, logger, &selectedRuntime, *runtimeVersionToDownload)
	if err != nil {
		logger.Errorf("Error checking %s cache %s", selectedRuntime.GetRuntimeName(), err)
		return nil, nil, err
	}
	if *exists {
		logger.Infof("%s version %s already exists in cache", selectedRuntime.GetRuntimeName(), (*runtimeVersionToDownload).GetVersion())
	} else {
		// Download and save to cache
		runtimeData, err := http.DownloadFile(createRuntimeDownloadUrlInfo.RuntimeUrl, &createRuntimeDownloadUrlInfo.Sha256, logger, "Downloading "+selectedRuntime.GetRuntimeName(), &createRuntimeDownloadUrlInfo.Sha512)
		if err != nil {
			logger.Errorf("Error downloading %s with cause %s", selectedRuntime.GetRuntimeName(), err)
			return nil, nil, err
		}
		filename, err = filemanagement.SaveRuntimeToCacheDir(runtimeData, *createRuntimeDownloadUrlInfo, logger)
		if err != nil {
			logger.Errorf("Error saving %s to cache %s", selectedRuntime.GetRuntimeName(), err)
			return nil, nil, err
		}
		logger.Infof("%s saved to cache at: %s", selectedRuntime.GetRuntimeName(), *filename)
	}
	if filename == nil {
		logger.Errorw("Filename is nil after checking cache and downloading", "error", err)
		return nil, nil, errors.New("filename is nil after checking cache and downloading")
	}

	targetPath, err := filemanagement.DoesTargetDirExist(*filename)
	if err != nil {
		logger.Errorf("Error creating target directory for %s with cause %s", selectedRuntime.GetRuntimeName(), err)
		return nil, nil, err
	}

	if filemanagement.HasArchiveBeenExtracted(*targetPath) {
		logger.Debugf("%s version %s already extracted at: %s", selectedRuntime.GetRuntimeName(), (*runtimeVersionToDownload).GetVersion(), *targetPath)
		relevantPaths := createRelevantRuntimePaths(*targetPath, selectedRuntime)
		selectedRuntimeFor = &selectedRuntime
		return &relevantPaths, selectedRuntimeFor, nil
	} else {
		// Unpack the runtime archive
		targetLocation, err := archive.UnarchiveFile(*filename, logger)
		if err != nil {
			logger.Errorf("Error extracting %s archive", selectedRuntime.GetRuntimeName())
			return nil, nil, err
		}
		logger.Debugf("%s extracted to: %s", selectedRuntime.GetRuntimeName(), *targetLocation)
		selectedRuntimeFor = &selectedRuntime
	}
	relevantPaths := createRelevantRuntimePaths(*targetPath, selectedRuntime)
	return &relevantPaths, selectedRuntimeFor, nil
}

func createDownloadUrl(runtimeVersionToDownload interfaces.IRuntimeVersion, runtime interfaces.IRuntime, logger *zap.SugaredLogger) (*models.CreateDownloadStruct, error) {
	shaSumsOFFiles, err := caching.GetShaSumCacheInPath(runtime, runtimeVersionToDownload.GetVersion())

	if err != nil {
		logger.Errorw("Error fetching SHASUMS256.txt", "error", err)
		return nil, err
	}

	filenamePrefix := runtime.GetFilenamePrefix(runtimeVersionToDownload.GetVersion())
	urlToruntime, err := runtime.ToDownloadUrl(filenamePrefix, *shaSumsOFFiles, runtimeVersionToDownload.GetVersion())
	if err != nil {
		logger.Errorw("Error creating runtime.js download URL", "error", err)
		return nil, err
	}
	downloadModel := archive.FilterCorrectFilenameEnding(filenamePrefix, *shaSumsOFFiles)
	if downloadModel == nil {
		logger.Errorw("No matching runtime.js binary found for your platform")
		return nil, errors.New("no matching runtime.js binary found for your platform")
	}
	downloadModel.RuntimeUrl = *urlToruntime
	return downloadModel, err
}
