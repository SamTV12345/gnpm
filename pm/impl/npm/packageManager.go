package npm

import (
	"encoding/json"
	"io"
	http3 "net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/samtv12345/gnpm/archive"
	"github.com/samtv12345/gnpm/http"
	http4 "github.com/samtv12345/gnpm/pm/impl/npm/http"
	http2 "github.com/samtv12345/gnpm/pm/impl/pnpm/http"
	"github.com/samtv12345/gnpm/pm/interfaces"
	"go.uber.org/zap"
)

type Npm struct {
	Logger *zap.SugaredLogger
}

func (n Npm) ExtractToFilesystem(targetPath string) (*string, error) {
	return archive.UnarchiveFile(targetPath, n.Logger)
}

func (n Npm) GetAllPathsToLink(targetPath string) []string {
	var targetPathsToLink = make([]string, 0)

	cmdPath := filepath.Join(targetPath, "bin", "npm.cmd")
	ps1Path := filepath.Join(targetPath, "bin", "npm.ps1")
	bashPath := filepath.Join(targetPath, "bin", "npm")
	targetPathsToLink = append(targetPathsToLink, bashPath)
	targetPathsToLink = append(targetPathsToLink, cmdPath)
	targetPathsToLink = append(targetPathsToLink, ps1Path)
	targetPathsToLink = append(targetPathsToLink, filepath.Join(targetPath, "bin", "npx.cmd"))
	targetPathsToLink = append(targetPathsToLink, filepath.Join(targetPath, "bin", "npx.ps1"))
	targetPathsToLink = append(targetPathsToLink, filepath.Join(targetPath, "bin", "npm-cli.js"))
	targetPathsToLink = append(targetPathsToLink, filepath.Join(targetPath, "bin", "npm-prefix.js"))
	targetPathsToLink = append(targetPathsToLink, filepath.Join(targetPath, "node_modules"))

	// Do code mods
	codeModCmd(cmdPath, n.Logger, targetPath)
	powershellMod(ps1Path, n.Logger, targetPath)
	bashMod(bashPath, n.Logger, targetPath)
	return targetPathsToLink
}

func bashMod(ps1Path string, logger *zap.SugaredLogger, targetPath string) {
	content, err := os.ReadFile(ps1Path)
	if err != nil {
		logger.Warnf("Error reading file %s: %v", ps1Path, err)
		return
	}
	newContent := string(content)
	newContent = strings.ReplaceAll(newContent, "$NPM_PREFIX/node_modules/npm/bin/", targetPath+"/")
	err = os.WriteFile(ps1Path, []byte(newContent), 0644)
	if err != nil {
		logger.Warnf("Error writing file %s: %v", ps1Path, err)
		return
	}
}

func powershellMod(ps1Path string, logger *zap.SugaredLogger, targetPath string) {
	content, err := os.ReadFile(ps1Path)
	if err != nil {
		logger.Warnf("Error reading file %s: %v", ps1Path, err)
		return
	}
	newContent := string(content)
	newContent = strings.ReplaceAll(newContent, "$PSScriptRoot/node_modules/npm/bin/", filepath.Join(targetPath, "bin")+"\\")
	err = os.WriteFile(ps1Path, []byte(newContent), 0644)
	if err != nil {
		logger.Warnf("Error writing file %s: %v", ps1Path, err)
		return
	}
}

func codeModCmd(cmdPath string, logger *zap.SugaredLogger, targetPath string) {
	content, err := os.ReadFile(cmdPath)
	if err != nil {
		logger.Warnf("Error reading file %s: %v", cmdPath, err)
		return
	}
	newContent := string(content)
	newContent = strings.ReplaceAll(newContent, "%~dp0\\node_modules\\npm\\bin\\", filepath.Join(targetPath, "bin")+"\\")
	newContent = strings.ReplaceAll(newContent, "%%F\\node_modules\\npm\\bin\\", filepath.Join(targetPath, "bin")+"\\")
	err = os.WriteFile(cmdPath, []byte(newContent), 0644)
	if err != nil {
		logger.Warnf("Error writing file %s: %v", cmdPath, err)
		return
	}
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
