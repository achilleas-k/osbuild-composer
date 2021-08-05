package osbuild2

type OSTreeReferences interface {
	isOSTreeReferences()
	isReferences()
}

type CommitReferences map[string]CommitReference

type CommitReference struct {
	// OSTree reference to create for this commit
	Ref string `json:"ref"`
}

func (CommitReferences) isOSTreeReferences() {}

// Alias ReferenceList to match OSTReeReferences interface
type OSTreeListReferences ReferenceList

func (OSTreeListReferences) isOSTreeReferences() {}

// Returns a new OSTree input object. References can be a list of strings or
// CommitReferences map.
func NewOSTreeInput(origin InputOriginType, references OSTreeReferences) *Input {
	return &Input{
		Type:       "org.osbuild.ostree",
		Origin:     origin,
		References: references,
	}
}
