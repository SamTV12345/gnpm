package commandRun

import (
	"os"
	"os/exec"

	"github.com/samtv12345/gnpm/detection"
	"go.uber.org/zap"
)

func RunCommand(detectionResult detection.PackageManagerDetectionResult, remainingArgs []string, logger *zap.SugaredLogger) {
	var cmd *exec.Cmd

	if remainingArgs[0] == "node" {
		logger.Info("Running node")
		cmdToRun := exec.Command("node", remainingArgs[1:]...)
		cmdToRun.Stdout = os.Stdout
		cmdToRun.Stderr = os.Stderr
		if err := cmdToRun.Start(); err != nil {
			logger.Errorw("Error running node command", "error", err)
		}
		if err := cmdToRun.Wait(); err != nil {
			logger.Errorw("Error waiting for node command to finish", "error", err)
		}
		return
	}

	if *detectionResult.Agent == detection.AgentNameNpm {
		cmd = exec.Command("npm", remainingArgs...)
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
	cmd.Start()
	cmd.Wait()
}
