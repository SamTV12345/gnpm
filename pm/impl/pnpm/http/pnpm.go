package http

import (
	"encoding/json"
	"io"
	"net/http"

	http2 "github.com/samtv12345/gnpm/http"
)

type DownloadReleaseResult struct {
	Filename string
	Content  []byte
}

func GetSpecificReleaseOfPnpm(version string) (*http2.GitHubRelease, error) {
	data, err := http.Get("https://api.github.com/repos/pnpm/pnpm/releases/tags/v" + version)
	if err != nil {
		return nil, err
	}
	var release http2.GitHubRelease
	defer data.Body.Close()
	readBytes, err := io.ReadAll(data.Body)
	if err := json.Unmarshal(readBytes, &release); err != nil {
		return nil, err
	}
	return &release, nil
}
