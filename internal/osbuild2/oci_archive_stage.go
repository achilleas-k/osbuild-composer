package osbuild2

// OSTreeCommitStageOptions describes how to assemble a tree into an OSTree commit.
// OCIArchiveStageOptions Assemble an OCI image archive
type OCIArchiveStageOptions struct {
	Architecture string
}

func (OCIArchiveStageOptions) isStageOptions() {}

// The OCIArchiveStage describes how to assemble an OCI image archive.
func NewOCIArchiveStage(options *OCIArchiveStageOptions) *Stage {
	return &Stage{
		Type:    "org.osbuild.oci-archive",
		Options: options,
	}
}
