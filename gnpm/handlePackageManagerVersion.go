package gnpm

import (
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/samtv12345/gnpm/caching"
	"github.com/samtv12345/gnpm/detection"
	"github.com/samtv12345/gnpm/filemanagement"
	"github.com/samtv12345/gnpm/http"
	"go.uber.org/zap"
)

func HandlePackageManagerVersion(remainingArgs []string, logger *zap.SugaredLogger, result detection.PackageManagerDetectionResult) {
	isInstalled, targetPath, err := filemanagement.IsPnpmVersionInInstallDir(*result.Version)
	if err != nil {
		logger.Warnf("Error checking if pnpm version is installed: %v", err)
		return
	}
	if *isInstalled && targetPath != nil {
		logger.Infof("pnpm version %s is already installed in %s", *result.Version, *targetPath)
		return
	}
	// Get all available pnpm versions
	pnpmVersions := caching.GetPnpmVersion(logger)
	vs := make([]*semver.Version, len(pnpmVersions))
	for i, v := range pnpmVersions {
		version, err := semver.NewVersion(v)
		if err != nil {
			logger.Warnf("Error parsing pnpm version %s: %v", v, err)
			continue
		}
		vs[i] = version
	}
	constraint, err := semver.NewConstraint(*result.Version)
	if err != nil {
		logger.Warnf("Error parsing version constraint %s: %v", *result.Version, err)
		return
	}
	var matchedVersions []*semver.Version
	for _, v := range vs {
		if constraint.Check(v) {
			matchedVersions = append(matchedVersions, v)
		}
	}
	if len(matchedVersions) == 0 {
		logger.Warnf("No matching versions found for constraint %s", *result.Version)
		return
	}
	sort.Sort(semver.Collection(matchedVersions))
	selectedVersion := matchedVersions[len(matchedVersions)-1]
	logger.Infof("Selected pnpm version: %s", selectedVersion.String())
	release, err := http.DownloadPnpmRelease(selectedVersion.String(), logger)
	if err != nil {
		logger.Warnf("Error getting release of pnpm: %v", err)
		return
	}
	targetPath, err = filemanagement.SavePnpmToInstallDir(release, logger, selectedVersion.String())

}
