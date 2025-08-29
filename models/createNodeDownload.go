package models

import "github.com/samtv12345/gnpm/http"

type CreateNodeDownloadStruct struct {
	NodeUrl string
	http.NodeShasumWithEncoding
}
