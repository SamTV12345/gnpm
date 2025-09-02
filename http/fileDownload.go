package http

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

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
