package packageJson

import (
	"encoding/json"
	"os"
)

func ReadPackageJson(path string) (*PackageManifest, error) {
	var packageFile, err = os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var packageJson PackageManifest

	err = json.Unmarshal(packageFile, &packageJson)
	if err != nil {
		return nil, err
	}

	return &packageJson, nil
}
