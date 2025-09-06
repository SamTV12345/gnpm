package models

type CreateFilenameStruct struct {
	Filename string
	Sha256   string
}

type CreateDownloadStruct struct {
	RuntimeUrl string
	Sha256     string
	Filename   string
	Encoding   string
}
