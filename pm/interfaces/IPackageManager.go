package interfaces

import (
	http2 "github.com/samtv12345/gnpm/pm/impl/pnpm/http"
)

type IPackageManager interface {
	GetVersionFileName() string
	GetName() string
	GetAllVersions() (*[]string, error)
	DownloadRelease(version string) (*http2.DownloadReleaseResult, error)
	GetAllPathsToLink(targetPath string) []string
	ExtractToFilesystem(targetPath string) (*string, error)
}
