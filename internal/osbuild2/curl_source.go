package osbuild2

type CurlSource struct {
	Items map[string]CurlSourceItem `json:"items"`
}

func (CurlSource) isSource() {}

// CurlSourceItem can be either a URL string or a URL paired with a secrets
// provider
type CurlSourceItem interface {
	isCurlSourceItem()
}

type URL string

func (URL) isCurlSourceItem() {}

type URLWithSecrets struct {
	URL     string      `json:"url"`
	Secrets *URLSecrets `json:"secrets,omitempty"`
}

func (URLWithSecrets) isCurlSourceItem() {}

type URLSecrets struct {
	Name string `json:"name"`
}
