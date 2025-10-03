package http

import (
	"encoding/json"
	"net/http"
)

func GetNpmRelease(version string) (*NpmRelease, error) {
	data, err := http.Get("https://registry.npmjs.org/npm/" + version)
	if err != nil {
		return nil, err
	}
	defer data.Body.Close()
	var release NpmRelease
	if err := json.NewDecoder(data.Body).Decode(&release); err != nil {
		return nil, err
	}
	return &release, nil
}
