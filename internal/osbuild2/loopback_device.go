package osbuild2

// Expose a file (or part of it) as a device node

type LoopbackDeviceOptions struct {
	// File to associate with the loopback device
	Filename string `json:"filename"`

	// Start of the data segment
	Start uint64 `json:"start,omitempty"`

	// Size limit of the data segment (in sectors)
	Size uint64 `json:"size,omitempty"`

	// Sector size (in bytes)
	SectorSize uint64 `json:"sector-size,omitempty"`
}

func (LoopbackDeviceOptions) isDeviceOptions() {}
