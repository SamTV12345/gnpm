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

func LinkPackageManager(targetNodePath, targetPackageManagerPath string, logger *zap.SugaredLogger, detection *detection.PackageManagerDetectionResult) error {
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

	if _, err := os.Stat(filepath.Join(moduleDir, "node")); err == nil {
		err = os.Remove(filepath.Join(moduleDir, "node"))
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(filepath.Join(moduleDir, detection.Name)); err == nil {
		err = os.Remove(filepath.Join(moduleDir, detection.Name))
		if err != nil {
			return err
		}
	}

	// Create symlink
	err = os.Symlink(targetNodePath, filepath.Join(moduleDir, "node"))
	if err != nil && !os.IsExist(err) {
		return err
	}

	logger.Debugf("Symlink %s to %s", targetNodePath, moduleDir+"/node")
	err = os.Symlink(targetPackageManagerPath, filepath.Join(moduleDir, detection.Name))
	logger.Debugf("Symlink %s to %s", targetPackageManagerPath, filepath.Join(moduleDir, detection.Name))
	if err != nil && !os.IsExist(err) {
		return err
	}
	return shell.PropagateChangesToCurrentShell(moduleDir)
}
