package http

type PnpmIndex struct {
	Versions map[string]any `json:"versions"`
}
