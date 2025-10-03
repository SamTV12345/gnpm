package http

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

func DownloadFile(url string, sha256Sum *string, logger *zap.SugaredLogger, title string, sha512Sum *string) ([]byte, error) {
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

	if sha512Sum != nil && *sha512Sum != "" {
		s := sha512.New()
		s.Write(data)
		shaSumFromPackage := s.Sum(nil)
		hash := "sha512-" + base64.StdEncoding.EncodeToString(shaSumFromPackage)
		if hash != *sha512Sum {
			return nil, errors.New("SHA512 checksum does not match")
		}
	}

	return data, nil
}
