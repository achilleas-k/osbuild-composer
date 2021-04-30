package osbuild2

type ChmodStageOptions struct {
	// Paths to operate on for changing mode bits
	Paths []string `json:"paths"`

	Mode string `json:"mode"`

	Recursive bool `json:"recursive,omitempty"`
}

func (ChmodStageOptions) isStageOptions() {}

// NewChmodStage creates a new org.osbuild.chmod stage
func NewChmodStage(options *ChmodStageOptions) *Stage {
	return &Stage{
		Type:    "org.osbuild.chmod",
		Options: options,
	}
}
