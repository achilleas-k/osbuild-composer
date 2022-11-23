package image

import (
	"fmt"
	"math/rand"

	"github.com/osbuild/osbuild-composer/internal/artifact"
	"github.com/osbuild/osbuild-composer/internal/common"
	"github.com/osbuild/osbuild-composer/internal/disk"
	"github.com/osbuild/osbuild-composer/internal/manifest"
	"github.com/osbuild/osbuild-composer/internal/ostree"
	"github.com/osbuild/osbuild-composer/internal/platform"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
	"github.com/osbuild/osbuild-composer/internal/runner"
	"github.com/osbuild/osbuild-composer/internal/users"
	"github.com/osbuild/osbuild-composer/internal/workload"
)

type OSTreeSimplifiedInstaller struct {
	Base

	Platform       platform.Platform
	Workload       workload.Workload
	PartitionTable *disk.PartitionTable

	Users  []users.User
	Groups []users.Group

	Commit ostree.CommitSpec

	SysrootReadOnly bool

	ISOLabelTempl string
	Product       string
	OSVersion     string
	Variant       string
	OSName        string

	KernelOptionsAppend []string
	Keyboard            string
	Locale              string

	Filename string
}

func NewOSTreeSimplifiedInstaller(commit ostree.CommitSpec) *OSTreeSimplifiedInstaller {
	return &OSTreeSimplifiedInstaller{
		Base:   NewBase("ostree-simplified-installer"),
		Commit: commit,
	}
}

func (img *OSTreeSimplifiedInstaller) InstantiateManifest(m *manifest.Manifest,
	repos []rpmmd.RepoConfig,
	runner runner.Runner,
	rng *rand.Rand) (*artifact.Artifact, error) {
	buildPipeline := manifest.NewBuild(m, runner, repos)
	buildPipeline.Checkpoint()

	// create the raw image
	osPipeline := manifest.NewOSTreeDeployment(m, buildPipeline, img.Commit, img.OSName, img.Platform)
	osPipeline.PartitionTable = img.PartitionTable
	osPipeline.KernelOptionsAppend = img.KernelOptionsAppend
	osPipeline.Keyboard = img.Keyboard
	osPipeline.Locale = img.Locale
	osPipeline.Users = img.Users
	osPipeline.Groups = img.Groups
	osPipeline.SysrootReadOnly = img.SysrootReadOnly

	imagePipeline := manifest.NewRawOStreeImage(m, buildPipeline, img.Platform, osPipeline)

	xzPipeline := manifest.NewXZ(m, buildPipeline, imagePipeline)
	xzPipeline.Filename = img.Filename

	coiPipeline := manifest.NewCOI(m,
		buildPipeline,
		img.Platform,
		repos,
		"kernel",
		img.Product,
		img.OSVersion)

	// create boot ISO with raw image
	rootfsImagePipeline := manifest.NewISORootfsImg(m, buildPipeline, coiPipeline)
	rootfsImagePipeline.Size = 4 * common.GibiByte

	bootTreePipeline := manifest.NewEFIBootTree(m, buildPipeline, img.Product, img.OSVersion)
	bootTreePipeline.Platform = img.Platform
	bootTreePipeline.UEFIVendor = img.Platform.GetUEFIVendor()
	bootTreePipeline.ISOLabel = "" // TODO

	isoLabel := fmt.Sprintf(img.ISOLabelTempl, img.Platform.GetArch())
	isoTreePipeline := manifest.NewISOTree(m,
		buildPipeline,
		coiPipeline,
		rootfsImagePipeline,
		bootTreePipeline,
		isoLabel)
	isoTreePipeline.PartitionTable = rootfsPartitionTable
	isoTreePipeline.Release = img.Release
	isoTreePipeline.OSName = img.OSName
	isoTreePipeline.Users = img.Users
	isoTreePipeline.Groups = img.Groups
	isoTreePipeline.PayloadPath = tarPath

	isoTreePipeline.SquashfsCompression = img.SquashfsCompression

	isoTreePipeline.OSPipeline = osPipeline
	isoTreePipeline.KernelOpts = img.AdditionalKernelOpts

	isoPipeline := manifest.NewISO(m, buildPipeline, isoTreePipeline)
	isoPipeline.Filename = img.Filename
	isoPipeline.ISOLinux = true

	return art, nil
}
