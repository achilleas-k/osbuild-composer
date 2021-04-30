package osbuild2

type ChownStageOptions struct {
	// Paths to operate on for changing mode bits
	Paths []string `json:"paths"`

	User  string `json:"user,omitempty"`
	Group string `json:"group,omitempty"`

	Recursive bool `json:"recursive,omitempty"`
}

func (ChownStageOptions) isStageOptions() {}

// NewChownStage creates a new org.osbuild.chown stage
func NewChownStage(options *ChownStageOptions) *Stage {
	return &Stage{
		Type:    "org.osbuild.chown",
		Options: options,
	}
}
