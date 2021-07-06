package osbuild2

// Install the grub2 boot loader for non-UEFI systems or hybrid boot

type Grub2InstStageOptions struct {
	// Filename of the disk image
	Filename string `json:"filename"`

	// Platform of the target system
	Platform string `json:"platform"`

	Location uint64 `json:"location,omitempty"`

	// How to obtain the GRUB core image
	Core CoreMkImage `json:"core"`

	// Location of grub config
	Prefix PrefixPartition `json:"prefix"`

	// Sector size (in bytes)
	SectorSize uint64 `json:"sector-size,omitempty"`
}

func (Grub2InstStageOptions) isStageOptions() {}

// Generate the core image via grub-mkimage
type CoreMkImage struct {
	// TODO: verify that it's "mkimage"
	Type string `json:"type"`

	// TODO: "gpt" or "dos"
	PartLabel string `json:"partlabel"`

	// TODO: "ext4", "xfs", or "btrfs"
	Filesystem string `json:"filesystem"`
}

// Grub2 config on a specific partition, e.g. (,gpt3)/boot
type PrefixPartition struct {
	// TODO: must be "partition"
	Type string `json:"type"`

	// TODO: "gpt" or "dos"
	PartLabel string `json:"partlabel"`

	// The partition number, starting at zero
	Number uint `json:"number"`

	// Location of the grub config inside the partition
	Path string `json:"path"`
}

func NewGrub2InstStage(options *Grub2InstStageOptions) *Stage {
	return &Stage{
		Type:    "org.osbuild.grub2.inst",
		Options: options,
	}
}
