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
	cwd, err := os.Getwd()
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

	// Download and link all node and pnpm versions
	nodeTargetPath, err := gnpm.HandleRuntimeVersion(remainingArgs[1:], logger)
	if err != nil {
		logger.Errorw("Error handling node version", "error", err)
		return
	}
	var packageManagerDecision = detection.DetectLockFileTool(cwd, logger)
	if packageManagerDecision == nil {
		logger.Info("No package manager detected")
	} else {
		logger.Infof("Package Manager detected: %s", packageManagerDecision.Name)
		pmTargetPath, err := gnpm.HandlePackageManagerVersion(remainingArgs[1:], logger, *packageManagerDecision)
		if err != nil {
			logger.Errorw("Error handling package manager version", "error", err)
			return
		}
		logger.Infof("Package manager %s installed at %s", packageManagerDecision.Name, *pmTargetPath)

		// Link
		*nodeTargetPath = append(*nodeTargetPath, *pmTargetPath)
		logger.Infof("Package Manager detected: %s", packageManagerDecision.Name)
	}

	err = gnpm.LinkRequiredPaths(*nodeTargetPath, logger, packageManagerDecision)
	if err != nil {
		logger.Errorf("Error linking package manager to %s", err)
		return
	}
	commandRun.RunCommand(packageManagerDecision, remainingArgs, logger)
}
