package runtimes

import (
	"github.com/samtv12345/gnpm/packageJson"
	"github.com/samtv12345/gnpm/runtimes/impl/bun"
	"github.com/samtv12345/gnpm/runtimes/impl/deno"
	"github.com/samtv12345/gnpm/runtimes/impl/node"
	"github.com/samtv12345/gnpm/runtimes/interfaces"
	"go.uber.org/zap"
)

func getSelectionFromPackageJson() string {
	readPackageJson, err := packageJson.ReadPackageJson("package.json")
	if err != nil {
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
	}
	return "node"
}

func GetRuntimeSelection(logger *zap.SugaredLogger) interfaces.IRuntime {
	var selection = getSelectionFromPackageJson()
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
