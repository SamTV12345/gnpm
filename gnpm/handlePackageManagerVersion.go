package gnpm

import (
	"errors"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/samtv12345/gnpm/caching"
	"github.com/samtv12345/gnpm/commandRun"
	"github.com/samtv12345/gnpm/detection"
	"github.com/samtv12345/gnpm/filemanagement"
	"github.com/samtv12345/gnpm/pm"
	"go.uber.org/zap"
)

func HandlePackageManagerVersion(args commandRun.FlagArguments, logger *zap.SugaredLogger, result detection.PackageManagerDetectionResult) (*[]string, error) {
	pmManager := pm.GetPackageManagerSelection(result.Name, logger)
	// Get all available pm versions
	pmVersions := caching.GetPmVersion(logger, pmManager, nil)
	vs := make([]*semver.Version, len(pmVersions))
	for i, v := range pmVersions {
		version, err := semver.NewVersion(v)
		if err != nil {
			logger.Warnf("Error parsing %s version %s: %v", pmManager.GetName(), v, err)
			continue
		}
		vs[i] = version
	}
	constraint, err := semver.NewConstraint(*result.Version)
	if err != nil {
		logger.Warnf("Error parsing version constraint %s: %v", *result.Version, err)
		return nil, err
	}
	var matchedVersions []*semver.Version
	for _, v := range vs {
		if constraint.Check(v) {
			matchedVersions = append(matchedVersions, v)
		}
	}

	if len(matchedVersions) == 0 {
		refreshPmVersions := true
		pmVersions = caching.GetPmVersion(logger, pmManager, &refreshPmVersions)
		vs = make([]*semver.Version, len(pmVersions))
		for i, v := range pmVersions {
			version, err := semver.NewVersion(v)
			if err != nil {
				logger.Warnf("Error parsing %s version %s: %v", pmManager.GetName(), v, err)
				continue
			}
			vs[i] = version
		}
		constraint, err = semver.NewConstraint(*result.Version)
		if err != nil {
			logger.Warnf("Error parsing version constraint %s: %v", *result.Version, err)
			return nil, err
		}
		for _, v := range vs {
			if constraint.Check(v) {
				matchedVersions = append(matchedVersions, v)
			}
		}

		if len(matchedVersions) == 0 {
			logger.Warnf("No matching versions found for constraint %s", *result.Version)
			return nil, errors.New("no matching versions found")
		}
	}
	sort.Sort(semver.Collection(matchedVersions))
	selectedVersion := matchedVersions[len(matchedVersions)-1]
	isInstalled, targetPath, err := filemanagement.IsPackageManagerInstalled(selectedVersion.String(), pmManager)
	if err != nil {
		logger.Warnf("Error checking if %s version is installed: %v", pmManager.GetName(), err)
		return nil, err
	}
	if *isInstalled && targetPath != nil {
		logger.Infof("%s version %s is already installed in %s", pmManager.GetName(), *result.Version, *targetPath)
		var pmPaths = pmManager.GetAllPathsToLink(*targetPath)
		return &pmPaths, nil
	} else {
		logger.Infof("Selected %s version: %s", result.Name, selectedVersion.String())
		release, err := pmManager.DownloadRelease(selectedVersion.String())
		if err != nil {
			logger.Warnf("Error getting release of %s: %v", pmManager.GetName(), err)
			return nil, err
		}
		targetPath, err = filemanagement.SavePackageManager(release, logger, selectedVersion.String(), pmManager)
		if err != nil {
			logger.Warnf("Error saving %s to install dir: %v", pmManager.GetName(), err)
			return nil, err
		}
		logger.Infof("Installed %s version %s in %s", result.Name, selectedVersion.String(), *targetPath)
		targetPath, err = pmManager.ExtractToFilesystem(*targetPath)
		if err != nil {
			logger.Warnf("Error extracting %s to filesystem: %v", pmManager.GetName(), err)
			return nil, err
		}

	}
	var pmPaths = pmManager.GetAllPathsToLink(*targetPath)
	return &pmPaths, nil
}
