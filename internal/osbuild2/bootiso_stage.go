package osbuild2

type BootISOStageOptions struct {
	Product Product `json:"product"`

	Kernel string `json:"kernel"`

	ISOLabel string `json:"isolabel"`

	EFI EFI `json:"efi,omitempty"`

	ISOLinux ISOLinux `json:"isolinux,omitempty"`

	// Additional kernel boot options
	KernelOpts string `json:"kernel_opts,omitempty"`

	Templates string `json:"templates,omitempty"`

	RootFS RootFS `json:"rootfs,omitempty"`
}

type EFI struct {
	Architectures []string `json:"architectures"`
	Vendor        string   `json:"vendor"`
}

type ISOLinux struct {
	Enabled bool `json:"enabled"`
	Debug   bool `json:"debug,omitempty"`
}

type RootFS struct {
	Compression struct {
		Method  string `json:"method"`
		Options struct {
			BCJ string `json:"bcj"`
		} `json:"options,omitempty"`
	} `json:"compression"`

	// Size in MiB
	Size int `json:"size"`
}

func (BootISOStageOptions) isStageOptions() {}

type BootISOStageInputs struct {
	RootFS *BootISOStageInput `json:"rootfs"`
	Kernel *BootISOStageInput `json:"kernel"`
}

func (BootISOStageInputs) isStageInputs() {}

type BootISOStageInput struct {
	inputCommon
	References BootISOStageReferences `json:"references"`
}

func (BootISOStageInput) isStageInput() {}

type BootISOStageReferences []string

func (BootISOStageReferences) isReferences() {}

// Assemble a file system tree for a bootable ISO
func NewBootISOStage(options *BootISOStageOptions, inputs *BootISOStageInputs) *Stage {
	return &Stage{
		Type:    "org.osbuild.bootiso",
		Options: options,
		Inputs:  inputs,
	}
}
