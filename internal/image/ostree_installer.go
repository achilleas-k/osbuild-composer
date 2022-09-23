package image

import (
	"math/rand"

	"github.com/osbuild/osbuild-composer/internal/artifact"
	"github.com/osbuild/osbuild-composer/internal/manifest"
	"github.com/osbuild/osbuild-composer/internal/platform"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
	"github.com/osbuild/osbuild-composer/internal/runner"
	"github.com/osbuild/osbuild-composer/internal/users"
)

type OSTreeInstaller struct {
	Base
	Platform          platform.Platform
	ExtraBasePackages rpmmd.PackageSet
	Users             []users.User
	Groups            []users.Group

	ISOLabelTempl string
	Product       string
	Variant       string
	OSName        string
	OSVersion     string
	Release       string

	OSTreeURL    string
	OSTreeRef    string
	OSTreeCommit string

	Filename string
}

func NewOSTreeInstaller() *OSTreeInstaller {
	return &OSTreeInstaller{
		Base: NewBase("ostree-installer"),
	}
}

func (img *OSTreeInstaller) InstantiateManifest(m *manifest.Manifest,
	repos []rpmmd.RepoConfig,
	runner runner.Runner,
	rng *rand.Rand) (*artifact.Artifact, error) {
	buildPipeline := manifest.NewBuild(m, runner, repos)
	buildPipeline.Checkpoint()

	anacondaPipeline := manifest.NewAnaconda(m,
		buildPipeline,
		img.Platform,
		repos, "kernel",
		img.Product,
		img.OSVersion)
	anacondaPipeline.ExtraPackages = img.ExtraBasePackages.Include
	anacondaPipeline.ExtraRepos = img.ExtraBasePackages.Repositories
	anacondaPipeline.Users = len(img.Users)+len(img.Groups) > 0
	anacondaPipeline.Variant = img.Variant
	anacondaPipeline.Biosdevname = (img.Platform.GetArch() == platform.ARCH_X86_64)
	anacondaPipeline.Checkpoint()

	isoTreePipeline := manifest.NewISOTree(m,
		buildPipeline,
		anacondaPipeline,
		nil,
		nil,
		img.OSTreeCommit,
		img.OSTreeURL,
		img.OSTreeRef,
		"")
	isoTreePipeline.Release = img.Release
	isoTreePipeline.OSName = img.OSName
	isoTreePipeline.Users = img.Users
	isoTreePipeline.Groups = img.Groups
	isoTreePipeline.KSPath = "/ostree.ks"

	isoPipeline := manifest.NewISO(m, buildPipeline, isoTreePipeline)
	isoPipeline.Filename = img.Filename
	isoPipeline.ISOLinux = true
	artifact := isoPipeline.Export()

	return artifact, nil
}
