package http

import (
	"strings"

	"github.com/samtv12345/gnpm/models"
)

func DecodeShasumTxt(shasumTxt string) []models.CreateFilenameStruct {
	var shaSumToFileMappingArr = make([]models.CreateFilenameStruct, 0)
	splittedShaSumData := strings.Split(shasumTxt, "\n")
	for _, line := range splittedShaSumData {
		shaSumToFileMapping := strings.Split(line, "  ")
		if len(shaSumToFileMapping) == 2 {
			shaSumToFileMappingArr = append(shaSumToFileMappingArr, models.CreateFilenameStruct{
				Sha256:   shaSumToFileMapping[0],
				Filename: shaSumToFileMapping[1],
			})
		}
	}
	return shaSumToFileMappingArr
}
