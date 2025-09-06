package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	http2 "github.com/samtv12345/gnpm/http"
	"go.uber.org/zap"
)

func GetDenoGitHubRelease(version string, logger *zap.SugaredLogger) (*http2.PnpmRelease, error) {
	response, err := http.Get("https://api.github.com/repos/denoland/deno/releases/tags/v" + version)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	shasumData, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error("Error reading SHASUMS256.txt:", err)
		return nil, err
	}

	var denoRelease http2.PnpmRelease
	if err := json.Unmarshal(shasumData, &denoRelease); err != nil {
		logger.Error("Error unmarshalling SHASUMS256.txt:", err)
		return nil, err
	}

	return &denoRelease, nil
}

func GetDenoVersions(logger *zap.SugaredLogger) (*[]DenoIndex, error) {
	data, err := http.Get("https://api.github.com/repos/denoland/deno/git/refs/tags")
	if err != nil {
		logger.Errorf("Error downloading deno versions: %s", err)
		return nil, err
	}
	defer data.Body.Close()
	var versions []DenoIndex
	readBytes, err := io.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}
	var denoIndex []DenoNpmResponse
	if err := json.Unmarshal(readBytes, &denoIndex); err != nil {
		return nil, err
	}
	for _, version := range denoIndex {
		if strings.Contains(version.Ref, "refs/tags/std") {
			continue
		}

		versions = append(versions, DenoIndex{
			Version: strings.Replace(version.Ref, "refs/tags/v", "", 1),
			IsLts:   false,
		})
	}

	return &versions, nil
}
