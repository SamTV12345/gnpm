package gnpm

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/samtv12345/gnpm/detection"
	"github.com/samtv12345/gnpm/filemanagement"
	"github.com/samtv12345/gnpm/shell"
	"go.uber.org/zap"
)

func LinkPackageManager(targetPaths []string, logger *zap.SugaredLogger, detection *detection.PackageManagerDetectionResult) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	hash := sha256.Sum256([]byte(cwd))
	hashHex := hex.EncodeToString(hash[:8])
	symlinkName := "gnpm-" + hashHex
	moduledir, err := filemanagement.EnsureModuleDir()
	if err != nil {
		return err
	}
	var moduleDir = filepath.Join(*moduledir, symlinkName)
	if err = os.Mkdir(moduleDir, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}

	for _, path := range targetPaths {
		filePathToCheck := filepath.Join(moduleDir, rewriteFileTargetName(path))
		logger.Debugf("Checking path: %s", filePathToCheck)
		if _, err := os.Lstat(filePathToCheck); err == nil {
			if err := os.Remove(filePathToCheck); err != nil {
				return err
			}
		}
	}

	for _, path := range targetPaths {
		logger.Debugf("Creating path: %s", path)
		filePathToCheck := filepath.Join(moduleDir, rewriteFileTargetName(path))
		err := os.Symlink(path, filePathToCheck)
		if err != nil {
			return err
		}
	}

	return shell.PropagateChangesToCurrentShell(moduleDir, logger)
}
