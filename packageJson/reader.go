package packageJson

import (
	"encoding/json"
	"os"
	"strings"
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

func ReadNvmrc(path string) (string, error) {
	var nvmrcFile, err = os.ReadFile(path)
	if err != nil {
		return "", err
	}
	contentOfNvmrc := string(nvmrcFile)
	contentOfNvmrc = strings.TrimSpace(contentOfNvmrc)
	return contentOfNvmrc, nil
}
