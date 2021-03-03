package osbuild2

type KickstartStageOptions struct {
	// Where to place the kickstart file
	Path string `json:"path"`

	OSTree OSTreeOptions `json:"ostree,omitempty"`

	LiveIMG struct {
		URL string `json:"url"`
	} `json:"liveimg"`
}

type OSTreeOptions struct {
	OSName string `json:"osname"`
	URL    string `json:"url"`
	Ref    string `json:"ref"`
	GPG    bool   `json:"gpg"`
}

func (KickstartStageOptions) isStageOptions() {}

// Creates an Anaconda kickstart file
func NewKickstartStageOption(options *KickstartStageOptions) *Stage {
	return &Stage{
		Type:    "org.osbuild.kickstart",
		Options: options,
	}
}
