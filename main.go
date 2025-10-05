package main

import (
	"os"

	"github.com/samtv12345/gnpm/commandRun"
	"github.com/samtv12345/gnpm/detection"
	"github.com/samtv12345/gnpm/gnpm"
	"github.com/samtv12345/gnpm/logging"
	"github.com/samtv12345/gnpm/shell"
)

func main() {
	var logger = logging.CreateLogger()
	defer logger.Sync()
	cwd, err := os.Getwd()
	if err != nil {
		logger.Error("Failed to get current working directory", err)
		return
	}
	var cmdFlags = commandRun.ParseFlags()
	var args = os.Args

	var remainingArgs = commandRun.FilterArgs(args[1:])

	if cmdFlags.Env {
		print(shell.ShowEnv(cwd))
		os.Exit(0)
	}

	if len(remainingArgs) == 0 {
		logger.Warn("You need to specify a command to run")
		return
	}
	// Download and link all runtime and pm versions
	runtimeTargetPath, selectedRuntime, err := gnpm.HandleRuntimeVersion(cmdFlags, logger)
	if err != nil || selectedRuntime == nil {
		logger.Errorw("Error handling runtime version", "error", err)
		return
	}
	var packageManagerDecision = detection.DetectLockFileTool(cwd, logger)
	if packageManagerDecision == nil {
		logger.Info("No package manager detected")
		os.Exit(1)
	} else {

		if cmdFlags.PackageManagerVersion != nil {
			packageManagerDecision.Version = cmdFlags.PackageManagerVersion
		}

		logger.Infof("Package Manager detected: %s", packageManagerDecision.Name)
		pmTargetPath, err := gnpm.HandlePackageManagerVersion(cmdFlags, logger, *packageManagerDecision)
		if err != nil {
			logger.Errorw("Error handling package manager version", "error", err)
			os.Exit(1)
		}
		logger.Infof("Package manager %s installed at %s", packageManagerDecision.Name, *pmTargetPath)

		// Link
		*runtimeTargetPath = append(*runtimeTargetPath, *pmTargetPath...)
	}

	err = gnpm.LinkRequiredPaths(*runtimeTargetPath, logger, packageManagerDecision)
	if err != nil {
		logger.Errorf("Error linking package manager to %s", err)
		os.Exit(2)
	}
	commandRun.RunCommand(packageManagerDecision, *selectedRuntime, remainingArgs, logger)
}
