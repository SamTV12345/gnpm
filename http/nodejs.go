package http

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

func GetNodeJsVersion(logger *zap.SugaredLogger) (*[]NodeIndex, error) {
	response, err := http.Get("https://nodejs.org/dist/index.json")
	if err != nil {
		logger.Error("Error fetching Node.js versions:", err)
		return nil, err
	}
	defer response.Body.Close()
	var nodeIndexes []NodeIndex
	err = json.NewDecoder(response.Body).Decode(&nodeIndexes)
	if err != nil {
		logger.Error("Error decoding Node.js versions:", err)
		return nil, err
	}
	return &nodeIndexes, nil
}

func GetShasumForNodeJSVersion(version string, logger *zap.SugaredLogger) (*[]NodeShasum, error) {
	response, err := http.Get("https://nodejs.org/dist/" + version + "/SHASUMS256.txt")
	if err != nil {
		logger.Error("Error fetching SHASUMS256.txt:", err)
		return nil, err
	}
	defer response.Body.Close()
	shasumData, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error("Error reading SHASUMS256.txt:", err)
		return nil, err
	}

	var shaSumData = string(shasumData)
	var shaSumToFileMappingArr = make([]NodeShasum, 0)
	splittedShaSumData := strings.Split(shaSumData, "\n")
	for _, line := range splittedShaSumData {
		shaSumToFileMapping := strings.Split(line, "  ")
		if len(shaSumToFileMapping) == 2 {
			shaSumToFileMappingArr = append(shaSumToFileMappingArr, NodeShasum{

				Sha256:   shaSumToFileMapping[0],
				Filename: shaSumToFileMapping[1],
			})
		}
	}
	return &shaSumToFileMappingArr, nil
}

func DownloadFile(url string, sha256Sum *string, logger *zap.SugaredLogger, title string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		logger.Error("Error downloading Node.js:", err)
		return nil, err
	}
	defer response.Body.Close()
	buf := bytes.Buffer{}
	bar := progressbar.NewOptions64(response.ContentLength, progressbar.OptionSetDescription(title), progressbar.OptionShowTotalBytes(true))
	multiWriter := io.MultiWriter(&buf, bar)
	_, err = io.Copy(multiWriter, response.Body)
	println()
	if err != nil {
		logger.Error("Error reading Node.js download:", err)
		return nil, err
	}
	data := buf.Bytes()
	if sha256Sum != nil {
		hash := sha256.Sum256(data)
		hashHex := hex.EncodeToString(hash[:])
		if *sha256Sum != "" && *sha256Sum != strings.ToLower(hashHex) {
			logger.Error("SHA256 checksum does not match")
			return nil, errors.New("SHA256 checksum does not match")
		}
	}
	return data, nil
}
