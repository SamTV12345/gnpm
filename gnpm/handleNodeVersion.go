package gnpm

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/samtv12345/gnpm/archive"
	"github.com/samtv12345/gnpm/caching"
	"github.com/samtv12345/gnpm/detection"
	"github.com/samtv12345/gnpm/filemanagement"
	"github.com/samtv12345/gnpm/http"
	"github.com/samtv12345/gnpm/models"
	"go.uber.org/zap"
)

func HandleNodeVersion(args []string, logger *zap.SugaredLogger) {
	nodeVersions, err := caching.GetNodeJsVersion(logger)

	if err != nil {
		logger.Errorw("Error fetching Node.js versions", "error", err)
		return
	}

	if len(args) == 0 {
		// Parse node version from .nvmrc or package.json
		nodeVersionToDownload, err := detection.GetNodeVersion(nil, logger, nodeVersions)
		if err != nil {
			logger.Errorw("Error determining Node.js version", "error", err)
			return
		}
		logger.Infof("Node.js version to use: %s", nodeVersionToDownload.Version)
		createNodeDownloadUrlInfo, err := createNodeDownloadUrl(*nodeVersionToDownload, logger)
		if err != nil {
			logger.Errorw("Error creating Node.js download URL", "error", err)
			return
		}
		logger.Debugf("Node.js download URL: %s", createNodeDownloadUrlInfo.NodeUrl)
		exists, filename, err := filemanagement.HasNodeVersionInCache(createNodeDownloadUrlInfo, logger)
		if err != nil {
			logger.Errorw("Error checking Node.js cache", "error", err)
			return
		}
		if *exists {
			logger.Infof("Node.js version %s already exists in cache", nodeVersionToDownload.Version)
		} else {
			// Download and save to cache
			nodeJsData, err := http.DownloadNodeJS(createNodeDownloadUrlInfo.NodeUrl, createNodeDownloadUrlInfo.Sha256, logger)
			if err != nil {
				logger.Errorw("Error downloading Node.js", "error", err)
				return
			}
			filename, err = filemanagement.SaveNodeJSToCacheDir(nodeJsData, *createNodeDownloadUrlInfo, logger)
			if err != nil {
				logger.Errorw("Error saving Node.js to cache", "error", err)
				return
			}
			logger.Infof("Node.js saved to cache at: %s", *filename)
		}
		if filename == nil {
			logger.Errorw("Filename is nil after checking cache and downloading", "error", err)
			return
		}

		targetPath, err := filemanagement.DoesTargetDirExist(*filename)
		if err != nil {
			logger.Errorw("Error creating target directory for Node.js extraction", "error", err)
			return
		}

		if filemanagement.HasArchiveBeenExtracted(*targetPath) {
			logger.Debugf("Node.js version %s already extracted at: %s", nodeVersionToDownload.Version, *targetPath)
			return
		} else {
			// Unpack the Node.js archive
			targetLocation, err := archive.UnarchiveFile(*filename, logger)
			if err != nil {
				logger.Errorw("Error extracting Node.js archive", "error", err)
				return
			}
			logger.Debugf("Node.js extracted to: %s", *targetLocation)
		}
	}
}

func filterCorrectFilenameEnding(filenamePrefix string, shaSumsOFFiles []http.NodeShasum) *http.NodeShasumWithEncoding {
	for _, file := range shaSumsOFFiles {
		if strings.HasPrefix(file.Filename, filenamePrefix) && (strings.HasSuffix(file.Filename, ".tar.xz") || strings.HasSuffix(file.Filename, ".zip")) {
			fileExtension := filepath.Ext(file.Filename)
			return &http.NodeShasumWithEncoding{
				NodeShasum: file,
				Encoding:   fileExtension,
			}
		}
	}
	return nil
}

func createNodeDownloadUrl(nodeVersionToDownload http.NodeIndex, logger *zap.SugaredLogger) (*models.CreateNodeDownloadStruct, error) {
	shaSumsOFFiles, err := http.GetShasumForNodeJSVersion(nodeVersionToDownload.Version, logger)

	if err != nil {
		logger.Errorw("Error fetching SHASUMS256.txt", "error", err)
		return nil, err
	}

	operatingSystem := runtime.GOOS
	architecture := runtime.GOARCH

	if operatingSystem == "windows" {
		operatingSystem = "win"
	}

	if strings.Contains(architecture, "amd") {
		architecture = strings.Replace(architecture, "amd", "x", 1)
	}

	var filenamePrefix = "node-" + nodeVersionToDownload.Version + "-" + operatingSystem + "-" + architecture

	filename := filterCorrectFilenameEnding(filenamePrefix, *shaSumsOFFiles)
	if filename == nil {
		logger.Errorw("Error finding correct Node.js binary for your platform", "error", err)
		return nil, errors.New("error finding correct Node.js binary for your platform")
	}

	urlToNode := "https://nodejs.org/dist/" + nodeVersionToDownload.Version + "/" + filename.Filename
	return &models.CreateNodeDownloadStruct{
		NodeUrl:                urlToNode,
		NodeShasumWithEncoding: *filename,
	}, nil
}
