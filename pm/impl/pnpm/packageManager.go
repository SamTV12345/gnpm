package pnpm

import (
	"encoding/json"
	"errors"
	"io"
	http2 "net/http"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/samtv12345/gnpm/http"
	http3 "github.com/samtv12345/gnpm/pm/impl/pnpm/http"
	"github.com/samtv12345/gnpm/pm/interfaces"
	"go.uber.org/zap"
)

type Pnpm struct {
	Logger *zap.SugaredLogger
}

func (p Pnpm) ExtractToFilesystem(targetPath string) (*string, error) {
	return &targetPath, nil
}

func (p Pnpm) GetAllPathsToLink(targetPath string) []string {
	return []string{targetPath}
}

func (p Pnpm) GetVersionFileName() string {
	return "pnpm.json"
}

func (p Pnpm) GetName() string {
	return "pnpm"
}

func (p Pnpm) GetAllVersions() (*[]string, error) {
	data, err := http2.Get("https://registry.npmjs.org/pnpm?fields=versions")
	if err != nil {
		return nil, err
	}
	defer data.Body.Close()
	var versions []string
	readBytes, err := io.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}
	var pnpmIndex http.GithubIndex
	if err := json.Unmarshal(readBytes, &pnpmIndex); err != nil {
		return nil, err
	}
	for version := range pnpmIndex.Versions {
		versions = append(versions, version)
	}

	return &versions, nil
}

func (p Pnpm) DownloadRelease(version string) (*http3.DownloadReleaseResult, error) {
	specificRelease, err := http3.GetSpecificReleaseOfPnpm(version)
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

	downloadedPnpmRelease, err := http.DownloadFile(url, &shasum, p.Logger, "Downloading pnpm", nil)
	if err != nil {
		return nil, err
	}
	return &http3.DownloadReleaseResult{
		Filename: filepath.Base(url),
		Content:  downloadedPnpmRelease,
	}, nil
}

var _ interfaces.IPackageManager = (*Pnpm)(nil)
