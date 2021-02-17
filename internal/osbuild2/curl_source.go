package osbuild2

import "encoding/json"

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

// Unmarshal method for CurlSource for handling the CurlSourceItem interface:
// Tries each of the implementations until it finds the one that works.
func (cs *CurlSource) UnmarshalJSON(data []byte) (err error) {
	cs.Items = make(map[string]CurlSourceItem)
	type csSimple struct {
		Items map[string]URL `json:"items"`
	}
	simple := new(csSimple)
	if err = json.Unmarshal(data, simple); err == nil && len(simple.Items) != 0 {
		for k, v := range simple.Items {
			cs.Items[k] = v
		}
		return
	}

	type csWithSecrets struct {
		Items map[string]URLWithSecrets `json:"items"`
	}
	withSecrets := new(csWithSecrets)
	if err = json.Unmarshal(data, withSecrets); err == nil && len(withSecrets.Items) != 0 {
		for k, v := range withSecrets.Items {
			cs.Items[k] = v
		}
		return
	}

	return
}
