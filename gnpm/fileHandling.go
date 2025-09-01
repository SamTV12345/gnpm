package gnpm

import (
	"path/filepath"
	"runtime"
	"strings"
)

func rewriteFileTargetName(path string) string {
	var filename = filepath.Base(path)
	operatingSystem := runtime.GOOS
	architecture := runtime.GOARCH

	filePrefix := ""

	if architecture == "amd64" {
		architecture = "x64"
	}

	if operatingSystem == "darwin" {
		filePrefix += "-win-" + architecture
	}

	if operatingSystem == "windows" {
		filePrefix += "-win-" + architecture
	} else if operatingSystem == "darwin" {
		filePrefix += "-macos-" + architecture
	} else {
		filePrefix += "-linux-" + architecture
	}

	return strings.Replace(filename, filePrefix, "", -1)
}
