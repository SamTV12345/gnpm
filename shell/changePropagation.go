package shell

import (
	"os"

	"go.uber.org/zap"
)

func PropagateChangesToCurrentShell(path string, logger *zap.SugaredLogger) error {
	logger.Debugf("Propagating changes to current shell")
	logger.Infof("Using path: %s", path)
	return os.Setenv("PATH", path+string(os.PathListSeparator)+os.Getenv("PATH"))
}
