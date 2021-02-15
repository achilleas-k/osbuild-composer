package osbuild2

type InitMode string

const (
	ModeBare         InitMode = "bare"
	ModeBareUser     InitMode = "bare-user"
	ModeBareUserOnly InitMode = "bare-user-only"
	ModeArchvie      InitMode = "archive"
)

// Options for the org.osbuild.ostree.init stage.
type OSTreeInitStageOptions struct {
	// The Mode in which to initialise the repo
	Mode InitMode
	// Location in which to create the repo
	Path string
}

func (OSTreeInitStageOptions) isStageOptions() {}

// A new org.osbuild.ostree.init stage with given options and inputs.
func NewOSTreeInitStage(options *OSTreeInitStageOptions, inputs Inputs) *Stage {
	return &Stage{
		Type:    "org.osbuild.ostree.init",
		Inputs:  inputs,
		Options: options,
	}
}
