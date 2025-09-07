package models

type CreateFilenameStruct struct {
	Filename string
	Sha256   string
}

type CreateDownloadStruct struct {
	RuntimeUrl string
	Sha256     string
	Sha512     string
	Filename   string
	Encoding   string
}
