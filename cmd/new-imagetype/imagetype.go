package main

import (
	"math/rand"

	"github.com/osbuild/osbuild-composer/internal/disk"
	"github.com/osbuild/osbuild-composer/internal/distro"
	"github.com/osbuild/osbuild-composer/internal/environment"
	"github.com/osbuild/osbuild-composer/internal/image"
	"github.com/osbuild/osbuild-composer/internal/manifest"
	"github.com/osbuild/osbuild-composer/internal/platform"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
	"github.com/osbuild/osbuild-composer/internal/runner"
	"github.com/osbuild/osbuild-composer/internal/workload"
)

type imageType struct {
	name string

	platform    platform.Platform
	environment environment.Environment
	workload    workload.Workload

	basePartitionTables distro.BasePartitionTableMap
}

func (it *imageType) Manifest(depsolve solver) (*manifest.Manifest, error) {
	m := manifest.New()
	rng := rand.New(rand.NewSource(9))
	repos := getRepos("fedora-37", "x86_64")
	pkgSets := []rpmmd.PackageSet{
		{
			Include:      []string{"@core"},
			Repositories: repos,
		},
	}
	runner := &runner.Fedora{Version: 37}
	pt, err := disk.NewPartitionTable(&basePT, nil, 0, false, nil, rng)
	check(err)

	img := image.NewLiveImage()
	img.Platform = &platform.X86{
		BIOS:       true,
		UEFIVendor: "fedora",
		BasePlatform: platform.BasePlatform{
			ImageFormat: platform.FORMAT_QCOW2,
			QCOW2Compat: "1.1",
		},
	}
	img.OSCustomizations = manifest.OSCustomizations{
		Language: "en_GB",
		Timezone: "UTC",
	}
	img.Environment = nil
	img.Workload = &workload.Custom{
		BaseWorkload: workload.BaseWorkload{
			Repos: repos,
		},
		Packages: pkgSets[0].Include,
	}
	img.Compression = ""
	img.PartitionTable = pt
	img.Filename = "disk.qcow2"

	_, err = img.InstantiateManifest(&m, repos, runner, rng)
	check(err)

	solved := depsolve(m.GetPackageSetChains())

	m.AddPackages(solved)

	return &m, nil
}
