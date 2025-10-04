package yarnClassic

import (
	"encoding/json"
	"io"
	http3 "net/http"
	"os"
	"path/filepath"

	"github.com/samtv12345/gnpm/archive"
	"github.com/samtv12345/gnpm/http"
	http2 "github.com/samtv12345/gnpm/pm/impl/pnpm/http"
	http5 "github.com/samtv12345/gnpm/pm/impl/yarnClassic/http"
	"github.com/samtv12345/gnpm/pm/interfaces"
	"go.uber.org/zap"
)

type Yarn struct {
	Logger *zap.SugaredLogger
}

func (y Yarn) GetVersionFileName() string {
	return "yarn.json"
}

func (y Yarn) GetName() string {
	return "yarn"
}

func (y Yarn) GetAllVersions() (*[]string, error) {
	data, err := http3.Get("https://registry.npmjs.org/yarn?fields=versions")
	if err != nil {
		return nil, err
	}
	defer data.Body.Close()
	var versions []string
	readBytes, err := io.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}
	var npmIndex http.GithubIndex
	if err := json.Unmarshal(readBytes, &npmIndex); err != nil {
		return nil, err
	}
	for version := range npmIndex.Versions {
		versions = append(versions, version)
	}

	return &versions, nil
}

func (y Yarn) DownloadRelease(version string) (*http2.DownloadReleaseResult, error) {
	specificRelease, err := http5.GetYarnClassicRelease(version)
	if err != nil {
		return nil, err
	}

	downloadUrl := specificRelease.Dist.Tarball
	if downloadUrl == "" {
		return nil, err
	}
	downloadedFile, err := http.DownloadFile(downloadUrl, nil, y.Logger, "Downloading npm version "+version, &specificRelease.Dist.Integrity)
	if err != nil {
		return nil, err
	}
	return &http2.DownloadReleaseResult{
		Filename: filepath.Base("yarn-" + version + ".tgz"),
		Content:  downloadedFile,
	}, nil
}

func codeModCmd(cmdPath string) {
	var yarnJSFile = filepath.Join(cmdPath, "bin", "yarn.js")
	var cmdContent = "@echo off\nnode " + yarnJSFile + " %*"

	if err := os.WriteFile(filepath.Join(cmdPath, "bin", "yarn.cmd"), []byte(cmdContent), 0755); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func (y Yarn) GetAllPathsToLink(targetPath string) []string {
	var targetPathsToLink = make([]string, 0)
	targetPathsToLink = append(targetPathsToLink, filepath.Join(targetPath, "bin", "yarn.js"))
	targetPathsToLink = append(targetPathsToLink, filepath.Join(targetPath, "bin", "yarn.cmd"))
	targetPathsToLink = append(targetPathsToLink, filepath.Join(targetPath, "bin", "yarn"))
	targetPathsToLink = append(targetPathsToLink, filepath.Join(targetPath, "bin", "yarnpkg"))
	targetPathsToLink = append(targetPathsToLink, filepath.Join(targetPath, "bin", "yarnpkg.cmd"))

	codeModCmd(targetPath)
	return targetPathsToLink
}

func (y Yarn) ExtractToFilesystem(targetPath string) (*string, error) {
	return archive.UnarchiveFile(targetPath, y.Logger)

}

var _ interfaces.IPackageManager = (*Yarn)(nil)
