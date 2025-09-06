package http

type DenoRelease struct {
	Assets  []DenoAsset `json:"assets"`
	TagName string      `json:"tag_name"`
}

type DenoAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}
