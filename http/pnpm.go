package http

import (
	"encoding/json"
	"io"
	"net/http"
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
