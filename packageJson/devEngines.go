package packageJson

type PackageManager struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type DevEngines struct {
	PackageManager *PackageManager `json:"packageManager,omitempty"`
}
