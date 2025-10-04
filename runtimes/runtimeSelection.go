package runtimes

import (
	"os"
	"path/filepath"

	"github.com/samtv12345/gnpm/packageJson"
	"github.com/samtv12345/gnpm/runtimes/impl/bun"
	"github.com/samtv12345/gnpm/runtimes/impl/deno"
	"github.com/samtv12345/gnpm/runtimes/impl/node"
	"github.com/samtv12345/gnpm/runtimes/interfaces"
	"go.uber.org/zap"
)

func getSelectionFromPackageJson(path string, logger *zap.SugaredLogger) string {
	readPackageJson, err := packageJson.ReadPackageJson(path)
	if err != nil {
		logger.Warnf("Error reading package json from directory %s: %s", path, err)
		os.Exit(1)

		return "node"
	}
	if readPackageJson.Engines != nil {
		if readPackageJson.Engines.Node != nil {
			return "node"
		} else if readPackageJson.Engines.Deno != nil {
			return "deno"
		} else if readPackageJson.Engines.Bun != nil {
			return "bun"
		}
	} else {
		newPath := filepath.Join(filepath.Dir(path), "..")
		newPackageJsonPath := filepath.Join(newPath, "package.json")
		if _, err := os.Stat(newPackageJsonPath); err == nil {
			return getSelectionFromPackageJson(newPackageJsonPath, logger)
		}
	}
	return "node"
}

func GetRuntimeSelection(logger *zap.SugaredLogger) interfaces.IRuntime {
	var currentDir, err = os.Getwd()
	if err != nil {
		logger.Errorf("Error getting current working directory %s", err)
		os.Exit(1)
	}
	var selection = getSelectionFromPackageJson(filepath.Join(currentDir, "package.json"), logger)
	logger.Infof("Using %s as runtime", selection)
	if selection == "node" {
		return &node.Runtime{
			Logger: logger,
		}
	} else if selection == "deno" {
		return &deno.Runtime{
			Logger: logger,
		}
	} else if selection == "bun" {
		return &bun.Runtime{
			Logger: logger,
		}
	} else {
		panic("Unrecognized runtime selection " + selection)
	}
}
