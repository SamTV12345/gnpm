package gnpm

import (
	"errors"
	"path/filepath"
	"runtime"

	"github.com/samtv12345/gnpm/archive"
	"github.com/samtv12345/gnpm/filemanagement"
	"github.com/samtv12345/gnpm/http"
	"github.com/samtv12345/gnpm/models"
	"github.com/samtv12345/gnpm/runtimes"
	"github.com/samtv12345/gnpm/runtimes/interfaces"
	"go.uber.org/zap"
)

func createRelevantNodePaths(targetPath string) []string {
	if runtime.GOOS == "windows" {
		nodePath := filepath.Join(targetPath, "node.exe")
		return []string{nodePath}
	}

	if runtime.GOOS == "linux" {
		nodePath := filepath.Join(targetPath, "bin", "node")
		return []string{nodePath}
	}

	var nodePath = filepath.Join(targetPath, "node")
	return []string{nodePath}
}

func HandleRuntimeVersion(args []string, logger *zap.SugaredLogger) (relevantPathsToReturn *[]string, selectedRuntimeFor *interfaces.IRuntime, err error) {
	var selectedRuntime = runtimes.GetRuntimeSelection(logger)
	nodeVersions, err := selectedRuntime.GetAllVersionsOfRuntime()

	if err != nil {
		logger.Errorf("Error fetching %s versions with cause %s", selectedRuntime.GetRuntimeName(), err)
		return nil, nil, err
	}

	// Parse node version from .nvmrc or package.json
	nodeVersionToDownload, err := selectedRuntime.GetInformationFromPackageJSON(nil, ".", nodeVersions)
	if err != nil {
		logger.Errorf("Error determining  version %s with cause %s %s", selectedRuntime.GetRuntimeName(), "error", err)
		return nil, nil, err
	}
	logger.Infof("%s version to use: %s", selectedRuntime.GetRuntimeName(), (*nodeVersionToDownload).GetVersion())
	createNodeDownloadUrlInfo, err := createDownloadUrl(*nodeVersionToDownload, selectedRuntime, logger)
	if err != nil {
		logger.Errorf("Error creating %s download URL with %s %s", selectedRuntime.GetRuntimeName(), "error", err)
		return nil, nil, err
	}
	logger.Debugf("Node.js download URL: %s", createNodeDownloadUrlInfo.NodeUrl)
	exists, filename, err := filemanagement.HasNodeVersionInCache(createNodeDownloadUrlInfo, logger)
	if err != nil {
		logger.Errorw("Error checking Node.js cache", "error", err)
		return nil, nil, err
	}
	if *exists {
		logger.Infof("Node.js version %s already exists in cache", (*nodeVersionToDownload).GetVersion())
	} else {
		// Download and save to cache
		nodeJsData, err := http.DownloadFile(createNodeDownloadUrlInfo.NodeUrl, &createNodeDownloadUrlInfo.Sha256, logger, "Downloading Node.js")
		if err != nil {
			logger.Errorw("Error downloading Node.js", "error", err)
			return nil, nil, err
		}
		filename, err = filemanagement.SaveNodeJSToCacheDir(nodeJsData, *createNodeDownloadUrlInfo, logger)
		if err != nil {
			logger.Errorw("Error saving Node.js to cache", "error", err)
			return nil, nil, err
		}
		logger.Infof("Node.js saved to cache at: %s", *filename)
	}
	if filename == nil {
		logger.Errorw("Filename is nil after checking cache and downloading", "error", err)
		return nil, nil, errors.New("filename is nil after checking cache and downloading")
	}

	targetPath, err := filemanagement.DoesTargetDirExist(*filename)
	if err != nil {
		logger.Errorw("Error creating target directory for Node.js extraction", "error", err)
		return nil, nil, err
	}

	if filemanagement.HasArchiveBeenExtracted(*targetPath) {
		logger.Debugf("Node.js version %s already extracted at: %s", (*nodeVersionToDownload).GetVersion(), *targetPath)
		relevantPaths := createRelevantNodePaths(*targetPath)
		selectedRuntimeFor = &selectedRuntime
		return &relevantPaths, selectedRuntimeFor, nil
	} else {
		// Unpack the Node.js archive
		targetLocation, err := archive.UnarchiveFile(*filename, logger)
		if err != nil {
			logger.Errorw("Error extracting Node.js archive", "error", err)
			return nil, nil, err
		}
		logger.Debugf("Node.js extracted to: %s", *targetLocation)
	}
	relevantPaths := createRelevantNodePaths(*targetPath)
	return &relevantPaths, selectedRuntimeFor, nil
}

func createDownloadUrl(nodeVersionToDownload interfaces.IRuntimeVersion, nodeRuntime interfaces.IRuntime, logger *zap.SugaredLogger) (*models.CreateDownloadStruct, error) {
	shaSumsOFFiles, err := nodeRuntime.GetShaSumsForRuntime(nodeVersionToDownload.GetVersion())

	if err != nil {
		logger.Errorw("Error fetching SHASUMS256.txt", "error", err)
		return nil, err
	}

	filenamePrefix := nodeRuntime.GetFilenamePrefix(nodeVersionToDownload.GetVersion())
	urlToNode, err := nodeRuntime.ToDownloadUrl(filenamePrefix, *shaSumsOFFiles, nodeVersionToDownload.GetVersion())
	if err != nil {
		logger.Errorw("Error creating Node.js download URL", "error", err)
		return nil, err
	}
	downloadModel := archive.FilterCorrectFilenameEnding(filenamePrefix, *shaSumsOFFiles)
	if downloadModel == nil {
		logger.Errorw("No matching Node.js binary found for your platform")
		return nil, errors.New("no matching Node.js binary found for your platform")
	}
	downloadModel.NodeUrl = *urlToNode
	return downloadModel, err
}
