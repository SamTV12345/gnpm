package http

type NodeShasum struct {
	Sha256   string
	Filename string
}

type NodeShasumWithEncoding struct {
	NodeShasum
	Encoding string
}
