package packageJson

type Dependencies = map[string]string
type PeerDependenciesMeta = map[string]string
type PackageScripts = map[string]string

type Engines struct {
	Node *string `json:"node,omitempty"`
	Npm  *string `json:"npm,omitempty"`
	Yarn *string `json:"yarn,omitempty"`
	Bun  *string `json:"bun,omitempty"`
	Deno *string `json:"deno,omitempty"`
}

type PublishConfig struct {
	directory       *string
	linkDirectory   *string
	executableFiles *[]string
	registry        *string
}

type Version = string
type Pattern = map[string][]string

type Versions = map[Version]Pattern

type PackageManifest struct {
	Name        string  `json:"name"`
	Version     string  `json:"version"`
	PType       *string `json:"type"`
	Bin         *any    `json:"bin,omitempty"`
	Description *string `json:"description,omitempty"`
	Directories *struct {
		Bin *string `json:"bin,omitempty"`
	} `json:"directories,omitempty"`
	Files                *[]string            `json:"files,omitempty"`
	Dependencies         *Dependencies        `json:"dependencies,omitempty"`
	DevDependencies      *Dependencies        `json:"devDependencies,omitempty"`
	OptionalDependencies *Dependencies        `json:"optionalDependencies,omitempty"`
	PeerDependencies     *Dependencies        `json:"peerDependencies,omitempty"`
	PeerDependenciesMeta PeerDependenciesMeta `json:"peerDependenciesMeta,omitempty"`
	DependenciesMeta     PeerDependenciesMeta `json:"dependenciesMeta,omitempty"`
	BundleDependencies   *interface{}         `json:"bundleDependencies,omitempty"`
	Homepage             *string              `json:"homepage,omitempty"`
	Repository           *interface{}         `json:"repository,omitempty"`
	Scripts              *PackageScripts      `json:"scripts,omitempty"`
	Config               *any                 `json:"config,omitempty"`
	Engines              *Engines             `json:"engines,omitempty"`
	Cpu                  []string             `json:"cpu,omitempty"`
	Os                   []string             `json:"os,omitempty"`
	Libc                 []string             `json:"libc,omitempty"`
	Main                 *string              `json:"main,omitempty"`
	Module               *string              `json:"module,omitempty"`
	Typings              *string              `json:"typings,omitempty"`
	Types                *string              `json:"types,omitempty"`
	PublishConfig        *PublishConfig       `json:"publishConfig,omitempty"`
	TypesVersions        *Versions            `json:"typesVersions,omitempty"`
	Readme               *string              `json:"readme,omitempty"`
	Author               *string              `json:"author,omitempty"`
	License              *string              `json:"license,omitempty"`
	Exports              *map[string]string   `json:"exports,omitempty"`
	Deprecated           *string              `json:"deprecated,omitempty"`
	PackageManager       *interface{}         `json:"packageManager,omitempty"`
	DevEngines           *DevEngines          `json:"devEngines,omitempty"`
}
