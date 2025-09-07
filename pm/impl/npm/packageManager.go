package npm

import (
	"encoding/json"
	"io"
	http3 "net/http"
	"path/filepath"

	"github.com/samtv12345/gnpm/http"
	http4 "github.com/samtv12345/gnpm/pm/impl/npm/http"
	http2 "github.com/samtv12345/gnpm/pm/impl/pnpm/http"
	"github.com/samtv12345/gnpm/pm/interfaces"
	"go.uber.org/zap"
)

type Npm struct {
	Logger *zap.SugaredLogger
}

func (n Npm) GetVersionFileName() string {
	return "npm.json"
}

func (n Npm) GetName() string {
	return "npm"
}

func (n Npm) GetAllVersions() (*[]string, error) {
	data, err := http3.Get("https://registry.npmjs.org/npm?fields=versions")
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

func (n Npm) DownloadRelease(version string) (*http2.DownloadReleaseResult, error) {
	specificRelease, err := http4.GetNpmRelease(version)
	if err != nil {
		return nil, err
	}

	downloadUrl := specificRelease.Dist.Tarball
	if downloadUrl == "" {
		return nil, err
	}
	downloadedFile, err := http.DownloadFile(downloadUrl, nil, n.Logger, "Downloading npm version "+version, &specificRelease.Dist.Integrity)
	if err != nil {
		return nil, err
	}
	return &http2.DownloadReleaseResult{
		Filename: filepath.Base("npm-" + version + ".tgz"),
		Content:  downloadedFile,
	}, nil
}

var _ interfaces.IPackageManager = (*Npm)(nil)
