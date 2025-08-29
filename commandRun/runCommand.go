package commandRun

import (
	"os"
	"os/exec"

	"github.com/samtv12345/gnpm/detection"
	"go.uber.org/zap"
)

func RunCommand(detectionResult detection.PackageManagerDetectionResult, remainingArgs []string, logger *zap.SugaredLogger) {
	var cmd *exec.Cmd
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
