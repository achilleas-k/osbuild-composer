package osbuild2

// Returns an input that represents the tree produced by the named pipeline.
func NewTreeInput(pipeline string) *Input {
	return &Input{
		Type:       "org.osbuild.tree",
		Origin:     "org.osbuild.pipeline",
		References: ReferenceList([]string{"name:" + pipeline}),
	}
}
