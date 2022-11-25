package manifest

import (
	"fmt"

	"github.com/osbuild/osbuild-composer/internal/disk"
	"github.com/osbuild/osbuild-composer/internal/osbuild"
	"github.com/osbuild/osbuild-composer/internal/users"
)

type COIISOTree struct {
	Base

	// TODO: review optional and mandatory fields and their meaning
	OSName  string
	Release string
	Users   []users.User
	Groups  []users.Group

	PartitionTable *disk.PartitionTable

	coiPipeline      *CoreOSInstaller
	bootTreePipeline *EFIBootTree

	// The path where the payload (tarball or ostree repo) will be stored.
	PayloadPath string

	isoLabel string

	KernelOpts []string
}

func NewCOIISOTree(m *Manifest,
	buildPipeline *Build,
	coiPipeline *CoreOSInstaller,
	bootTreePipeline *EFIBootTree,
	isoLabel string) *COIISOTree {

	p := &COIISOTree{
		Base:             NewBase(m, "bootiso-tree", buildPipeline),
		coiPipeline:      coiPipeline,
		bootTreePipeline: bootTreePipeline,
		isoLabel:         isoLabel,
	}
	buildPipeline.addDependent(p)
	if coiPipeline.Base.manifest != m {
		panic("anaconda pipeline from different manifest")
	}
	m.addPipeline(p)
	return p
}

func (p *COIISOTree) serialize() osbuild.Pipeline {
	pipeline := p.Base.serialize()

	kernelOpts := []string{}

	if len(p.KernelOpts) > 0 {
		kernelOpts = append(kernelOpts, p.KernelOpts...)
	}

	pipeline.AddStage(osbuild.NewMkdirStage(&osbuild.MkdirStageOptions{
		Paths: []osbuild.Path{
			{
				Path: "images",
			},
			{
				Path: "images/pxeboot",
			},
		},
	}))

	inputName := "tree"
	copyStageOptions := &osbuild.CopyStageOptions{
		Paths: []osbuild.CopyStagePath{
			{
				From: fmt.Sprintf("input://%s/boot/vmlinuz-%s", inputName, p.coiPipeline.kernelVer),
				To:   "tree:///images/pxeboot/vmlinuz",
			},
			{
				From: fmt.Sprintf("input://%s/boot/initramfs-%s.img", inputName, p.coiPipeline.kernelVer),
				To:   "tree:///images/pxeboot/initrd.img",
			},
		},
	}
	copyStageInputs := osbuild.NewPipelineTreeInputs(inputName, p.coiPipeline.Name())
	copyStage := osbuild.NewCopyStageSimple(copyStageOptions, copyStageInputs)
	pipeline.AddStage(copyStage)

	isoLinuxOptions := &osbuild.ISOLinuxStageOptions{
		Product: osbuild.ISOLinuxProduct{
			Name:    p.coiPipeline.product,
			Version: p.coiPipeline.version,
		},
		Kernel: osbuild.ISOLinuxKernel{
			Dir:  "/images/pxeboot",
			Opts: kernelOpts,
		},
	}
	isoLinuxStage := osbuild.NewISOLinuxStage(isoLinuxOptions, p.coiPipeline.Name())
	pipeline.AddStage(isoLinuxStage)

	filename := "images/efiboot.img"
	pipeline.AddStage(osbuild.NewTruncateStage(&osbuild.TruncateStageOptions{
		Filename: filename,
		Size:     fmt.Sprintf("%d", p.PartitionTable.Size),
	}))

	efibootDevice := osbuild.NewLoopbackDevice(&osbuild.LoopbackDeviceOptions{Filename: filename})
	for _, stage := range osbuild.GenMkfsStages(p.PartitionTable, efibootDevice) {
		pipeline.AddStage(stage)
	}

	inputName = "root-tree"
	copyInputs := osbuild.NewPipelineTreeInputs(inputName, p.bootTreePipeline.Name())
	copyOptions, copyDevices, copyMounts := osbuild.GenCopyFSTreeOptions(inputName, p.bootTreePipeline.Name(), filename, p.PartitionTable)
	pipeline.AddStage(osbuild.NewCopyStage(copyOptions, copyInputs, copyDevices, copyMounts))

	copyInputs = osbuild.NewPipelineTreeInputs(inputName, p.bootTreePipeline.Name())
	pipeline.AddStage(osbuild.NewCopyStageSimple(
		&osbuild.CopyStageOptions{
			Paths: []osbuild.CopyStagePath{
				{
					From: fmt.Sprintf("input://%s/EFI", inputName),
					To:   "tree:///",
				},
			},
		},
		copyInputs,
	))

	return pipeline
}
