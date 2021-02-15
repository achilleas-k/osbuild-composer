package osbuild2

type OSTreeCommitStageOptions struct {
	// OStree ref to create for the commit
	Ref string `json:"ref"`
	// Set the version of the OS as commit metadata
	OSVersion string `json:"os_version,omitempty"`
	// Commit ID of the parent commit
	Parent string `json:"parent,omitempty"`
}

func (OSTreeCommitStageOptions) isStageOptions() {}

// The OSTreeCommitStage (org.osbuild.ostree.commit) describes how to assemble
// a tree into an OSTree commit.
func NewOSTreeCommitStage(options *OSTreeCommitStageOptions) *Stage {
	return &Stage{
		Type:    "org.osbuild.ostree.commit",
		Options: options,
	}
}
