package http

type DenoIndex struct {
	Version string
	IsLts   bool
}

func (d DenoIndex) GetVersion() string {
	return d.Version
}

func (d DenoIndex) IsLTS() bool {
	return d.IsLts
}

type DenoNpmResponse struct {
	Ref string `json:"ref"` // e.g. refs/tags/v1.34.3
}
