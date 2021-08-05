package osbuild2

// Collection of Inputs for a Stage
type Inputs map[string]Input

// Single Input for a Stage
type Input struct {
	// Input type
	Type string `json:"type"`

	// Origin should be either 'org.osbuild.source' or 'org.osbuild.pipeline'
	Origin InputOriginType `json:"origin"`

	References References `json:"references"`
}

// TODO: define these using type aliases
type InputOriginType string

const (
	InputOriginSource   InputOriginType = "org.osbuild.source"
	InputOriginPipeline                 = "org.osbuild.pipeline"
)

type inputCommon struct {
	Type string `json:"type"`
	// Origin should be either 'org.osbuild.source' or 'org.osbuild.pipeline'
	Origin string `json:"origin"`
}
type References interface {
	isReferences()
}

type ReferenceList []string

func (ReferenceList) isReferences() {}
