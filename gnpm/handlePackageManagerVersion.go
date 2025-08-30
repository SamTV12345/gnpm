package gnpm

import (
	"github.com/samtv12345/gnpm/caching"
	"github.com/samtv12345/gnpm/detection"
	"go.uber.org/zap"
)

func HandlePackageManagerVersion(remainingArgs []string, logger *zap.SugaredLogger, result detection.PackageManagerDetectionResult) {
	caching.GetPnpmVersion(logger)
}
