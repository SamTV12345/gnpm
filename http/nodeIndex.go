package http

type NodeIndex struct {
	Version  string      `json:"version"`
	Date     string      `json:"date"`
	Files    []string    `json:"files"`
	Npm      string      `json:"npm"`
	V8       string      `json:"v8"`
	UV       string      `json:"uv"`
	Zlib     string      `json:"zlib"`
	OpenSSL  string      `json:"openssl"`
	Modules  string      `json:"modules"`
	LTS      interface{} `json:"lts"`
	Security bool        `json:"security"`
}
