package main

import (
	"os"

	"github.com/samtv12345/gnpm/commandRun"
	"github.com/samtv12345/gnpm/detection"
	"github.com/samtv12345/gnpm/gnpm"
	"github.com/samtv12345/gnpm/logging"
)

func main() {
	var logger = logging.CreateLogger()
	defer logger.Sync()
	var cwd, err = os.Getwd()
	if err != nil {
		logger.Error("Failed to get current working directory", err)
		return
	}
	var args = os.Args
	if len(args) == 1 {
		logger.Warn("You need to specify a command to run")
		return
	}
	var remainingArgs = args[1:]

	if remainingArgs[0] == "use" {
		gnpm.HandleNodeVersion(remainingArgs[1:], logger)
		var packageManagerDecision = detection.DetectLockFileTool(cwd, logger)
		if packageManagerDecision == nil {
			logger.Info("No package manager detected")
			return
		}
		logger.Infof("Package Manager detected: %s", packageManagerDecision.Name)
		gnpm.HandlePackageManagerVersion(remainingArgs[1:], logger, *packageManagerDecision)
	} else {
		var packageManagerDecision = detection.DetectLockFileTool(cwd, logger)
		if packageManagerDecision == nil {
			logger.Info("No package manager detected")
			return
		}
		logger.Infof("Package Manager detected: %s", packageManagerDecision.Name)
		commandRun.RunCommand(*packageManagerDecision, remainingArgs, logger)
	}
}
