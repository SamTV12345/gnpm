package runtimes

import (
	"github.com/samtv12345/gnpm/runtimes/impl/bun"
	"github.com/samtv12345/gnpm/runtimes/impl/deno"
	"github.com/samtv12345/gnpm/runtimes/impl/node"
	"github.com/samtv12345/gnpm/runtimes/interfaces"
	"go.uber.org/zap"
)

func GetRuntimeSelection(selection string, logger *zap.SugaredLogger) interfaces.IRuntime {
	if selection == "node" {
		return &node.Runtime{
			Logger: logger,
		}
	} else if selection == "deno" {
		return &deno.Runtime{}
	} else if selection == "bun" {
		return &bun.Runtime{}
	} else {
		panic("Unrecognized runtime selection " + selection)
	}
}
