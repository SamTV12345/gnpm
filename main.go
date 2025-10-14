package main

import (
	"os"
	"sync"

	"github.com/samtv12345/gnpm/commandRun"
	"github.com/samtv12345/gnpm/detection"
	"github.com/samtv12345/gnpm/gnpm"
	"github.com/samtv12345/gnpm/logging"
	"github.com/samtv12345/gnpm/runtimes/interfaces"
	"github.com/samtv12345/gnpm/shell"
)

type LockableLinkPaths struct {
	mutex sync.Mutex
	paths []string
}

func main() {
	var logger = logging.CreateLogger()
	defer logger.Sync()
	cwd, err := os.Getwd()
	if err != nil {
		logger.Error("Failed to get current working directory", err)
		os.Exit(1)
	}
	var cmdFlags = commandRun.ParseFlags()
	var args = os.Args

	var remainingArgs = commandRun.FilterArgs(args[1:])

	if cmdFlags.Env {
		shell.ShowEnv(cwd)
		return
	}

	if len(remainingArgs) == 0 {
		logger.Warn("You need to specify a command to run")
		os.Exit(1)
	}

	var pathLinks = LockableLinkPaths{
		paths: make([]string, 0),
		mutex: sync.Mutex{},
	}

	var packageManagerDecision *detection.PackageManagerDetectionResult
	var selectedRuntime *interfaces.IRuntime

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		// Download and link all runtime and pm versions
		runtimeTargetPath, selectedRuntimeForGnpm, err := gnpm.HandleRuntimeVersion(cmdFlags, logger)

		if err != nil || selectedRuntimeForGnpm == nil {
			logger.Errorw("Error handling runtime version", "error", err)
			os.Exit(1)
		}
		selectedRuntime = selectedRuntimeForGnpm
		pathLinks.mutex.Lock()
		pathLinks.paths = append(pathLinks.paths, *runtimeTargetPath...)
		pathLinks.mutex.Unlock()
		waitGroup.Done()
		return
	}()

	go func() {
		// Handle package manager
		packageManagerDecision = detection.DetectLockFileTool(cwd, logger)
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
			pathLinks.mutex.Lock()
			pathLinks.paths = append(pathLinks.paths, *pmTargetPath...)
			pathLinks.mutex.Unlock()
			waitGroup.Done()
		}
	}()

	waitGroup.Wait()
	pathLinks.mutex.Lock()
	err = gnpm.LinkRequiredPaths(pathLinks.paths, logger)
	pathLinks.mutex.Unlock()
	if err != nil {
		logger.Errorf("Error linking package manager to %s", err)
		os.Exit(2)
	}
	commandRun.RunCommand(packageManagerDecision, *selectedRuntime, remainingArgs, logger)
}
