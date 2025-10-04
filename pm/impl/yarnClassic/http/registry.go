package http

import (
	"encoding/json"
	"net/http"
)

func GetYarnClassicRelease(version string) (*YarnClassicRelease, error) {
	data, err := http.Get("https://registry.npmjs.org/yarn/" + version)
	if err != nil {
		return nil, err
	}
	defer data.Body.Close()
	var release YarnClassicRelease
	if err := json.NewDecoder(data.Body).Decode(&release); err != nil {
		return nil, err
	}
	return &release, nil
}
