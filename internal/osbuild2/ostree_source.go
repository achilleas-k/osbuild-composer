package osbuild2

type OSTreeSource struct {
	// URL of the repository.
	URL string `json:"url"`
	// GPG keys to verify the commits
	GPGKeys []string `json:"secrets,omitempty"`
}

func (OSTreeSource) isSource() {}

// The commits to fetch indexed their checksum
type OSTreeSoures struct {
	Items map[string]OSTreeSource
}
