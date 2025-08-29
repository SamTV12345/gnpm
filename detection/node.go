package detection

import (
	"errors"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/samtv12345/gnpm/http"
	"github.com/samtv12345/gnpm/packageJson"
	"go.uber.org/zap"
)

func GetNodeVersion(proposedVersion *string, logger *zap.SugaredLogger, nodeVersions *[]http.NodeIndex) (*http.NodeIndex, error) {
	var versionToDownload *string
	packageManifest, err := packageJson.ReadPackageJson("package.json")
	if err != nil {
		return nil, err
	}

	// Check .nvmrc first
	nvmrcVersion, errNvmrc := packageJson.ReadNvmrc(".nvmrc")
	if errNvmrc == nil {
		versionToDownload = &nvmrcVersion
	}

	// Then check the package.json "engines" field
	if packageManifest != nil && packageManifest.Engines != nil && packageManifest.Engines.Node != nil {
		versionToDownload = packageManifest.Engines.Node
	}

	if proposedVersion != nil {
		versionToDownload = proposedVersion
	}

	if versionToDownload != nil {
		constraints, err := semver.NewConstraint(*versionToDownload)
		if err != nil {
			logger.Errorw("Error parsing version", "error", err)
		}

		var possibleVersions = make([]*semver.Version, 0)
		for _, nodeVersion := range *nodeVersions {
			v, err := semver.NewVersion(nodeVersion.Version[1:])
			if err != nil {
				logger.Errorw("Error parsing version", "error", err)
				continue
			}
			if constraints.Check(v) {
				possibleVersions = append(possibleVersions, v)
			}
		}

		// Sort possibleVersions in descending order
		sort.Sort(semver.Collection(possibleVersions))

		if len(possibleVersions) > 0 {
			latestVersion := possibleVersions[len(possibleVersions)-1].String()
			versionToDownload = &latestVersion
		}
	} else {
		// Default to the latest LTS version if no version is specified
		for _, nodeVersion := range *nodeVersions {
			if nodeVersion.LTS != false {
				versionToDownload = &nodeVersion.Version
				break
			}
		}
	}

	if versionToDownload == nil {
		return nil, errors.New("error finding a suitable Node.js version")
	}

	for _, nodeVersion := range *nodeVersions {
		if nodeVersion.Version == *versionToDownload || nodeVersion.Version == "v"+*versionToDownload {
			return &nodeVersion, nil
		}
	}
	return nil, errors.New("error finding a suitable Node.js version")
}
