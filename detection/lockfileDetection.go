package detection

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/samtv12345/gnpm/packageJson"
	"go.uber.org/zap"
)

func pathExists(path string, isDirectory bool) bool {
	result, err := os.Stat(path)
	if err != nil {
		return false
	}

	if isDirectory {
		return result.IsDir()
	}

	return !result.IsDir()
}

func Lookup(cwd string) <-chan string {
	ch := make(chan string)
	go func() {
		directory := filepath.Clean(cwd)
		root := filepath.VolumeName(directory) + string(filepath.Separator)
		for directory != "" && directory != root {
			ch <- directory
			directory = filepath.Dir(directory)
		}
		close(ch)
	}()
	return ch
}

type PackageManagerDetectionResult struct {
	Name    string
	Agent   *string
	Version *string
}

func handleVer(version *string) *string {
	if version == nil {
		return nil
	}
	re := regexp.MustCompile(`\d+(\.\d+){0,2}`)
	match := re.FindString(*version)
	if match != "" {
		return &match
	}
	var asterisk = "*"
	return &asterisk
}

func getNameAndVer(pm packageJson.PackageManifest) *PackageManagerDetectionResult {

	if pm.PackageManager != nil {
		pmString, ok := (*pm.PackageManager).(string)
		if pm.PackageManager != nil && ok {
			var depRangeMarker = regexp.MustCompile("^\\^")
			replacedPackageManager := depRangeMarker.ReplaceAll([]byte(pmString), []byte(""))
			var resultingNameAndVersion = strings.Split(string(replacedPackageManager), "@")

			return &PackageManagerDetectionResult{
				Name:    resultingNameAndVersion[0],
				Version: handleVer(&resultingNameAndVersion[1]),
				Agent:   nil,
			}
		}
	}
	if pm.DevEngines != nil && pm.DevEngines.PackageManager != nil && pm.DevEngines.PackageManager.Name != "" {
		return &PackageManagerDetectionResult{
			Name:    pm.DevEngines.PackageManager.Name,
			Version: handleVer(&pm.DevEngines.PackageManager.Version),
			Agent:   nil,
		}
	}

	return nil
}

func handlePackageManager(filePath string, logger *zap.SugaredLogger) *PackageManagerDetectionResult {
	pm, err := packageJson.ReadPackageJson(filePath)
	var agent string
	if err != nil {
		logger.Error("Error reading package.json file")
		return nil
	}
	result := getNameAndVer(*pm)

	if result != nil {
		var name = result.Name
		var version = result.Version
		var res = FromStringToAgentName(result.Name)

		if res == AgentNameYarn {
			if parsedVersion, err := strconv.Atoi(*result.Version); err == nil && parsedVersion >= 2 {
				var agentBerry = AgentYarnBerry
				return &PackageManagerDetectionResult{
					Name:    result.Name + "@berry",
					Version: version,
					Agent:   &agentBerry,
				}
			} else {
				var agentClassic = AgentYarn
				return &PackageManagerDetectionResult{
					Name:    result.Name + "@classic",
					Version: version,
					Agent:   &agentClassic,
				}
			}
		} else if name == "pnpm" {
			if parsedVersion, err := strconv.Atoi(*result.Version); err == nil && parsedVersion < 7 {
				agent = "pnpm@6"
				return &PackageManagerDetectionResult{
					Name:    name,
					Version: version,
					Agent:   &agent,
				}
			} else if slices.Contains(Agents, name) {
				agent = name
				return &PackageManagerDetectionResult{
					Name:    name,
					Version: version,
					Agent:   &agent,
				}
			}
		} else if slices.Contains(Agents, name) {
			agent = name
			return &PackageManagerDetectionResult{
				Name:    name,
				Version: version,
				Agent:   &agent,
			}
		} else {
			logger.Warn("[gnpm] Unknown packageManager:" + name)
		}
	}

	return nil
}

func DetectLockFileTool(path string, logger *zap.SugaredLogger) *PackageManagerDetectionResult {
	var strategies = []string{
		"lockfile",
		"packageManager-field",
		"devEngines-field",
	}

	for dir := range Lookup(path) {
		for _, strategy := range strategies {
			logger.Debugf("[gnpm] Detecting lockfile strategy: %s", strategy)
			switch strategy {
			case "lockfile":
				for lock, name := range LOCKS {
					if pathExists(filepath.Join(dir, lock), false) {
						if pathExists(filepath.Join(dir, "package.json"), false) {
							// Package Manager field in package.json takes precedence over lockfile detection
							var result = handlePackageManager(filepath.Join(dir, "package.json"), logger)
							if result != nil {
								if result.Version != nil {
									logger.Infof("[gnpm] Detected package manager: %s in version %s", result.Name, *result.Version)
								} else {
									logger.Infof("[gnpm] Detected package manager: %s", result.Name)

								}
								return result
							} else {
								var asterisk = "*"
								return &PackageManagerDetectionResult{
									Name:    name,
									Version: &asterisk,
									Agent:   &name,
								}
							}
						}
					}
				}
				break
			case "packageManager-field", "devEngines-field":
				if pathExists(filepath.Join(dir, "package.json"), false) {
					var result = handlePackageManager(filepath.Join(dir, "package.json"), logger)
					if result != nil {
						return result
					}
				}
				break
			case "install-metadata":
				{
					for metadata, tool := range INSTALL_METADATA {
						var isMetadataDir = strings.HasPrefix(metadata, "/")
						if pathExists(filepath.Join(dir, metadata), isMetadataDir) {
							var result = PackageManagerDetectionResult{
								Name:  tool,
								Agent: &tool,
							}
							if tool == "yarn" {
								if isMetadataYarnClassic(metadata) {
									var agentYarn = AgentYarn
									result.Agent = &agentYarn
								} else {
									var agentYarn = AgentYarnBerry
									result.Agent = &agentYarn
								}
							}

							return &result
						}
					}
					break
				}
			}
		}
		if dir == path {
			break
		}
	}
	return nil
}
