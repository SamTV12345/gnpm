package detection

import "strings"

func isMetadataYarnClassic(metadataPath string) bool {
	return strings.HasPrefix(metadataPath, ".yarn_integrity")
}
