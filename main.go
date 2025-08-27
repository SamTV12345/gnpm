package main

import (
	"github.com/samtv12345/gnpm/detection"
	"github.com/samtv12345/gnpm/logging"
)

func main() {
	var logger = logging.CreateLogger()
	defer logger.Sync()
	detection.HandlePackageManager(".", logger)
}
