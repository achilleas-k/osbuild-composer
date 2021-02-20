package osbuild2

type RPMOSTreePrepTreeStageOptions struct {
	EtcGroupMembers []string `json:"etc_group_members,omitempty"`
}

func (RPMOSTreePrepTreeStageOptions) isStageOptions() {}

// The RPM OSTree PrepTree (org.osbuild.ostree.preptree) stage transforms the
// tree to an ostree layout.
func NewRPMOSTreePrepTreeStage(options *RPMOSTreePrepTreeStageOptions) *Stage {
	return &Stage{
		Type:    "org.osbuild.ostree.preptree",
		Options: options,
	}
}
