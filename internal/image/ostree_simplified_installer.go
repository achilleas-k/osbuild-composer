package image

import (
	"fmt"
	"math/rand"

	"github.com/osbuild/osbuild-composer/internal/artifact"
	"github.com/osbuild/osbuild-composer/internal/common"
	"github.com/osbuild/osbuild-composer/internal/disk"
	"github.com/osbuild/osbuild-composer/internal/environment"
	"github.com/osbuild/osbuild-composer/internal/manifest"
	"github.com/osbuild/osbuild-composer/internal/platform"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
	"github.com/osbuild/osbuild-composer/internal/runner"
	"github.com/osbuild/osbuild-composer/internal/workload"
)

type OSTreeSimplifiedInstaller struct {
	Base

	// Raw image that will be created and embedded
	rawImage *OSTreeRawImage

	Platform         platform.Platform
	OSCustomizations manifest.OSCustomizations
	Environment      environment.Environment
	Workload         workload.Workload

	ExtraBasePackages rpmmd.PackageSet

	// ISO label template (architecture-free)
	ISOLabelTempl string

	// Product string for ISO buildstamp
	Product string

	// OSVersion string for ISO buildstamp
	OSVersion string

	// Variant string for ISO buildstamp
	Variant string

	// OSName for ostree deployment
	OSName string

	installDevice string

	Filename string
}

func NewOSTreeSimplifiedInstaller(rawImage *OSTreeRawImage, installDevice string) *OSTreeSimplifiedInstaller {
	return &OSTreeSimplifiedInstaller{
		Base:          NewBase("ostree-simplified-installer"),
		rawImage:      rawImage,
		installDevice: installDevice,
	}
}

func (img *OSTreeSimplifiedInstaller) InstantiateManifest(m *manifest.Manifest,
	repos []rpmmd.RepoConfig,
	runner runner.Runner,
	rng *rand.Rand) (*artifact.Artifact, error) {
	buildPipeline := manifest.NewBuild(m, runner, repos)
	buildPipeline.Checkpoint()

	imageFilename := "image.raw.xz"

	// create the raw image
	img.rawImage.Filename = imageFilename
	rr := rawImagePipelines(img.rawImage, m, buildPipeline)

	coiPipeline := manifest.NewCOI(m,
		buildPipeline,
		img.Platform,
		repos,
		"kernel",
		img.Product,
		img.OSVersion,
		img.Variant)
	coiPipeline.ExtraPackages = img.ExtraBasePackages.Include
	coiPipeline.ExtraRepos = img.ExtraBasePackages.Repositories

	isoLabel := fmt.Sprintf(img.ISOLabelTempl, img.Platform.GetArch())

	// create boot ISO with raw image
	// rootfsImagePipeline := manifest.NewISORootfsImg(m, buildPipeline, coiPipeline)
	// rootfsImagePipeline.Size = 4 * common.GibiByte
	kernelOpts := []string{"rd.neednet=1",
		"coreos.inst.crypt_root=1",
		"coreos.inst.isoroot=" + isoLabel,
		"coreos.inst.install_dev=" + img.installDevice,
		"coreos.inst.image_file=/run/media/iso/disk.img.xz",
		"coreos.inst.insecure"}

	bootTreePipeline := manifest.NewEFIBootTree(m, buildPipeline, img.Product, img.OSVersion)
	bootTreePipeline.Platform = img.Platform
	bootTreePipeline.UEFIVendor = img.Platform.GetUEFIVendor()
	bootTreePipeline.ISOLabel = isoLabel
	bootTreePipeline.KernelOpts = kernelOpts

	rootfsPartitionTable := &disk.PartitionTable{
		Size: 20 * common.MebiByte,
		Partitions: []disk.Partition{
			{
				Start: 0,
				Size:  20 * common.MebiByte,
				Payload: &disk.Filesystem{
					Type:       "vfat",
					Mountpoint: "/",
					UUID:       disk.NewVolIDFromRand(rng),
				},
			},
		},
	}

	isoTreePipeline := manifest.NewCOIISOTree(m,
		buildPipeline,
		rr,
		coiPipeline,
		bootTreePipeline,
		isoLabel)
	isoTreePipeline.PartitionTable = rootfsPartitionTable
	isoTreePipeline.OSName = img.OSName

	isoPipeline := manifest.NewISO(m, buildPipeline, isoTreePipeline, isoLabel)
	isoPipeline.Filename = img.Filename
	isoPipeline.ISOLinux = true

	artifact := isoPipeline.Export()
	return artifact, nil
}
