package osbuild2

// Options for the org.osbuild.ostree.pull stage.
type OSTreePullStageOptions struct {
	// Location of the ostree repo
	Repo string
}

func (OSTreePullStageOptions) isStageOptions() {}

// A new org.osbuild.ostree.init stage with given options and inputs.
func NewOSTreePullStage(options *OSTreePullStageOptions, inputs Inputs) *Stage {
	return &Stage{
		Type:    "org.osbuild.ostree.init",
		Inputs:  inputs,
		Options: options,
	}
}

type OSTreePullStageInput struct {
	Commits interface{}
}
