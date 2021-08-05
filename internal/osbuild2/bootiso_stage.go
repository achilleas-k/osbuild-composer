package osbuild2

type BootISOMonoStageOptions struct {
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
	Compression FSCompression `json:"compression"`

	// Size in MiB
	Size int `json:"size"`
}

type FSCompression struct {
	Method  string                `json:"method"`
	Options *FSCompressionOptions `json:"options,omitempty"`
}

type FSCompressionOptions struct {
	BCJ string `json:"bcj"`
}

// BCJOption returns the appropriate xz branch/call/jump (BCJ) filter for the
// given architecture
func BCJOption(arch string) string {
	switch arch {
	case "x86_64":
		return "x86"
	case "aarch64":
		return "arm"
	case "ppc64le":
		return "powerpc"
	}
	return ""
}

func (BootISOMonoStageOptions) isStageOptions() {}

// Assemble a file system tree for a bootable ISO
func NewBootISOMonoStage(options *BootISOMonoStageOptions, rootPipeline string) *Stage {
	// NOTE: The bootiso.mono stage is deprecated and we never used the
	// 'kernel' property, so keeping support for it is not necessary
	inputs := Inputs{
		"rootfs": *NewTreeInput(rootPipeline),
	}
	return &Stage{
		Type:    "org.osbuild.bootiso.mono",
		Options: options,
		Inputs:  inputs,
	}
}
