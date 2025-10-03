package commandRun

import (
	"os"
	"os/exec"

	"github.com/samtv12345/gnpm/detection"
	"github.com/samtv12345/gnpm/runtimes/interfaces"
	"go.uber.org/zap"
)

func RunCommand(detectionResult *detection.PackageManagerDetectionResult, selectedRuntime interfaces.IRuntime, remainingArgs []string, logger *zap.SugaredLogger) {
	var cmd *exec.Cmd

	if remainingArgs[0] == selectedRuntime.GetRuntimeName() {
		logger.Infof("Running %s", selectedRuntime.GetRuntimeName())
		cmdToRun := exec.Command(selectedRuntime.GetRuntimeName(), remainingArgs[1:]...)
		cmdToRun.Stdout = os.Stdout
		cmdToRun.Stderr = os.Stderr
		if err := cmdToRun.Run(); err != nil {
			logger.Errorf("Error running %s command with %s %s", selectedRuntime.GetRuntimeName(), "error", err)
		}
		return
	}

	if detectionResult == nil || detectionResult.Agent == nil {
		logger.Warn("No package manager detected, cannot run command")
		return
	}

	logger.Infof("Running command with package manager: %s", *detectionResult.Agent)

	// Prepare the command based on the detected package manager
	// The first argument in remainingArgs is the command to run (e.g., install, add, etc.)
	// We need to prepend the package manager executable to the command

	if *detectionResult.Agent == detection.AgentNameNpm {
		cmd = exec.Command("npm.cmd", remainingArgs...)
	} else if *detectionResult.Agent == detection.AgentNameYarn {
		cmd = exec.Command("yarn", remainingArgs...)
	} else if *detectionResult.Agent == detection.AgentNamePnpm {
		cmd = exec.Command("pnpm", remainingArgs...)
	} else if *detectionResult.Agent == detection.AgentNameBun {
		cmd = exec.Command("bun", remainingArgs...)
	} else if *detectionResult.Agent == detection.AgentNameDeno {
		cmd = exec.Command("deno", remainingArgs...)
	} else {
		logger.Warnf("Unsupported package manager: %s", *detectionResult.Agent)
		return
	}

	if cmd == nil {
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Errorw("Error running command", "error", err)
	}
}
