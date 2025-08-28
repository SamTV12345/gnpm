package main

import (
	"os"

	"github.com/samtv12345/gnpm/detection"
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

	var packageManagerDecision = detection.DetectLockFileTool(cwd, logger)
	if packageManagerDecision == nil {
		logger.Info("No package manager detected")
		return
	}
	println(packageManagerDecision.Name)
}
