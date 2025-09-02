package models

type CreateFilenameStruct struct {
	Filename string
	Sha256   string
}

type CreateDownloadStruct struct {
	NodeUrl  string
	Sha256   string
	Filename string
	Encoding string
}
