package http

import (
	"encoding/json"
	"net/http"

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
