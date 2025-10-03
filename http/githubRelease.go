package http

type GithubIndex struct {
	Versions map[string]any `json:"versions"`
}

type Asset struct {
	BrowserDownloadURL string `json:"browser_download_url"`
	Name               string `json:"name"`
	Digest             string `json:"digest"`
}

type GitHubRelease struct {
	Assets []Asset `json:"assets"`
}
