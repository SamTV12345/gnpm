package http

type PnpmIndex struct {
	Versions map[string]any `json:"versions"`
}

type Asset struct {
	BrowserDownloadURL string `json:"browser_download_url"`
	Digest             string `json:"digest"`
}

type PnpmRelease struct {
	Assets []Asset `json:"assets"`
}
