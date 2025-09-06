package http

import (
	"strings"

	"github.com/samtv12345/gnpm/runtimes/interfaces"
)

type BunNpmResponse struct {
	Versions map[string]any `json:"versions"`
	DistTags DistTags       `json:"dist-tags"`
}

type DistTags struct {
	Latest string `json:"latest"`
	Canary string `json:"canary"`
}

type BunIndex struct {
	Version string
	IsLts   bool
}

func (b BunIndex) GetVersion() string {
	if strings.Contains(b.Version, "canary") {
		return strings.Split(b.Version, "-")[0]
	}
	return b.Version
}

func (b BunIndex) IsLTS() bool {
	return b.IsLts
}

var _ interfaces.IRuntimeVersion = (*BunIndex)(nil)
