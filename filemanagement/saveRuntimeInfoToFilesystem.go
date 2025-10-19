package filemanagement

import (
	"encoding/json"
	"os"
	"path/filepath"

	http2 "github.com/samtv12345/gnpm/runtimes/impl/bun/http"
	http3 "github.com/samtv12345/gnpm/runtimes/impl/deno/http"
	"github.com/samtv12345/gnpm/runtimes/impl/node/http"
	"github.com/samtv12345/gnpm/runtimes/interfaces"
)

func SaveNodeInfoToFilesystem[T interfaces.IRuntimeVersion](nodeIndices []T) error {
	jsonBytes, err := json.Marshal(nodeIndices)
	if err != nil {
		return err
	}
	cacheDir, err := GetCacheDir()
	if err != nil {
		return err
	}
	pathToSaveTo := filepath.Join(*cacheDir, "nodejs_index.json")
	return os.WriteFile(pathToSaveTo, jsonBytes, os.ModePerm)
}

func SaveBunInfoToFilesystem[T interfaces.IRuntimeVersion](nodeIndices []T) error {
	jsonBytes, err := json.Marshal(nodeIndices)
	if err != nil {
		return err
	}
	cacheDir, err := GetCacheDir()
	if err != nil {
		return err
	}
	pathToSaveTo := filepath.Join(*cacheDir, "bun_index.json")
	return os.WriteFile(pathToSaveTo, jsonBytes, os.ModePerm)
}

func SaveDenoInfoToFilesystem[T interfaces.IRuntimeVersion](nodeIndices []T) error {
	jsonBytes, err := json.Marshal(nodeIndices)
	if err != nil {
		return err
	}
	cacheDir, err := GetCacheDir()
	if err != nil {
		return err
	}
	pathToSaveTo := filepath.Join(*cacheDir, "deno_index.json")
	return os.WriteFile(pathToSaveTo, jsonBytes, os.ModePerm)
}

func ReadNodeInfoFromFilesystem() (*[]interfaces.IRuntimeVersion, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return nil, err
	}
	pathToReadFrom := filepath.Join(*cacheDir, "nodejs_index.json")
	fileBytes, err := os.ReadFile(pathToReadFrom)
	if err != nil {
		return nil, err
	}
	var nodeIndices []http.NodeIndex
	err = json.Unmarshal(fileBytes, &nodeIndices)
	if err != nil {
		return nil, err
	}
	converted := make([]interfaces.IRuntimeVersion, len(nodeIndices))
	for i, v := range nodeIndices {
		converted[i] = &v
	}

	return &converted, nil
}

func ReadBunInfoFromFilesystem() (*[]interfaces.IRuntimeVersion, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return nil, err
	}
	pathToReadFrom := filepath.Join(*cacheDir, "bun_index.json")
	fileBytes, err := os.ReadFile(pathToReadFrom)
	if err != nil {
		return nil, err
	}
	var nodeIndices []http2.BunIndex
	err = json.Unmarshal(fileBytes, &nodeIndices)
	if err != nil {
		return nil, err
	}
	converted := make([]interfaces.IRuntimeVersion, len(nodeIndices))
	for i, v := range nodeIndices {
		converted[i] = &v
	}

	return &converted, nil
}

func ReadDenoInfoFromFilesystem() (*[]interfaces.IRuntimeVersion, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return nil, err
	}
	pathToReadFrom := filepath.Join(*cacheDir, "deno_index.json")
	fileBytes, err := os.ReadFile(pathToReadFrom)
	if err != nil {
		return nil, err
	}
	var nodeIndices []http3.DenoIndex
	err = json.Unmarshal(fileBytes, &nodeIndices)
	if err != nil {
		return nil, err
	}
	converted := make([]interfaces.IRuntimeVersion, len(nodeIndices))
	for i, v := range nodeIndices {
		converted[i] = &v
	}

	return &converted, nil
}
