package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func GetBunVersions(logger *zap.SugaredLogger) (*[]BunIndex, error) {
	data, err := http.Get("https://registry.npmjs.org/bun?fields=versions")
	if err != nil {
		logger.Errorf("Error downloading bun versions: %s", err)
		return nil, err
	}
	defer data.Body.Close()
	var versions []BunIndex
	readBytes, err := io.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}
	var bunIndex BunNpmResponse
	if err := json.Unmarshal(readBytes, &bunIndex); err != nil {
		return nil, err
	}
	for version := range bunIndex.Versions {
		versions = append(versions, BunIndex{
			Version: version,
			IsLts:   bunIndex.DistTags.Latest == version,
		})
	}

	return &versions, nil
}
