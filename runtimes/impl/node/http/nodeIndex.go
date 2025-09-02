package http

import (
	"github.com/samtv12345/gnpm/runtimes/interfaces"
)

type NodeIndex struct {
	Version  string      `json:"version"`
	Date     string      `json:"date"`
	Files    []string    `json:"files"`
	Npm      string      `json:"npm"`
	V8       string      `json:"v8"`
	UV       string      `json:"uv"`
	Zlib     string      `json:"zlib"`
	OpenSSL  string      `json:"openssl"`
	Modules  string      `json:"modules"`
	LTS      interface{} `json:"lts"`
	Security bool        `json:"security"`
}

func (n NodeIndex) GetVersion() string {
	return n.Version
}

func (n NodeIndex) IsLTS() bool {
	return n.LTS != false && n.LTS != nil
}

var _ interfaces.IRuntimeVersion = (*NodeIndex)(nil)
