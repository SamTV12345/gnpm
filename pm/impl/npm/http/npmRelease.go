package http

type NpmRelease struct {
	Dist struct {
		Tarball    string `json:"tarball"`
		Shasum     string `json:"shasum"`
		Integrity  string `json:"integrity"`
		Signatures []struct {
			KeyID     string `json:"keyid"`
			Signature string `json:"signature"`
		}
	}
}
