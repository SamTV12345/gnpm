package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

	"go.uber.org/zap"
)

func GetAllVersionsOfPnpm() (*[]string, error) {
	data, err := http.Get("https://registry.npmjs.org/pnpm?fields=versions")
	if err != nil {
		return nil, err
	}
	defer data.Body.Close()
	var versions []string
	readBytes, err := io.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}
	var pnpmIndex PnpmIndex
	if err := json.Unmarshal(readBytes, &pnpmIndex); err != nil {
		return nil, err
	}
	for version := range pnpmIndex.Versions {
		versions = append(versions, version)
	}

	return &versions, nil
}

type DownloadPnpmReleaseResult struct {
	Filename string
	Content  []byte
}

func DownloadPnpmRelease(version string, logger *zap.SugaredLogger) (*DownloadPnpmReleaseResult, error) {
	specificRelease, err := getSpecificReleaseOfPnpm(version)
	if err != nil {
		return nil, err
	}

	operatingSystem := runtime.GOOS
	architecture := runtime.GOARCH
	if operatingSystem == "windows" {
		operatingSystem = "win"
	}

	if architecture == "amd64" {
		architecture = "x64"
	}

	var url string
	var shasum string
	for _, asset := range specificRelease.Assets {
		if strings.Contains(asset.BrowserDownloadURL, "pnpm-"+operatingSystem+"-"+architecture) {
			url = asset.BrowserDownloadURL
			shasum = strings.TrimPrefix(asset.Digest, "sha256:")
			break
		}
	}

	if url == "" {
		return nil, errors.New("no compatible pnpm binary found for your platform")
	}

	downloadedPnpmRelease, err := DownloadFile(url, &shasum, logger, "Downloading pnpm")
	if err != nil {
		return nil, err
	}
	return &DownloadPnpmReleaseResult{
		Filename: filepath.Base(url),
		Content:  downloadedPnpmRelease,
	}, nil
}

func getSpecificReleaseOfPnpm(version string) (*PnpmRelease, error) {
	data, err := http.Get("https://api.github.com/repos/pnpm/pnpm/releases/tags/v" + version)
	if err != nil {
		return nil, err
	}
	var release PnpmRelease
	defer data.Body.Close()
	readBytes, err := io.ReadAll(data.Body)
	if err := json.Unmarshal(readBytes, &release); err != nil {
		return nil, err
	}
	return &release, nil
}
