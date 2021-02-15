package osbuild2

type Secret struct {
	Name string `json:"name,omitempty"`
}

type CurlSource struct {
	URL     string  `json:"url"`
	Secrets *Secret `json:"secrets,omitempty"`
}

type CurlSources struct {
	URLs map[string]CurlSource `json:"urls"`
}

func (CurlSource) isSource() {}
