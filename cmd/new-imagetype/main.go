package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"

	"github.com/osbuild/osbuild-composer/internal/disk"
	"github.com/osbuild/osbuild-composer/internal/distro"
	"github.com/osbuild/osbuild-composer/internal/dnfjson"
	"github.com/osbuild/osbuild-composer/internal/environment"
	"github.com/osbuild/osbuild-composer/internal/image"
	"github.com/osbuild/osbuild-composer/internal/manifest"
	"github.com/osbuild/osbuild-composer/internal/osbuild"
	"github.com/osbuild/osbuild-composer/internal/platform"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
	"github.com/osbuild/osbuild-composer/internal/runner"
	"github.com/osbuild/osbuild-composer/internal/workload"
)

var basePT = disk.PartitionTable{
	UUID: "D209C89E-EA5E-4FBD-B161-B461CCE297E0",
	Type: "gpt",
	Partitions: []disk.Partition{
		{
			Size:     1048576, // 1MB
			Bootable: true,
			Type:     disk.BIOSBootPartitionGUID,
			UUID:     disk.BIOSBootPartitionUUID,
		},
		{
			Size: 209715200, // 200 MB
			Type: disk.EFISystemPartitionGUID,
			UUID: disk.EFISystemPartitionUUID,
			Payload: &disk.Filesystem{
				Type:         "vfat",
				UUID:         disk.EFIFilesystemUUID,
				Mountpoint:   "/boot/efi",
				Label:        "EFI-SYSTEM",
				FSTabOptions: "defaults,uid=0,gid=0,umask=077,shortname=winnt",
				FSTabFreq:    0,
				FSTabPassNo:  2,
			},
		},
		{
			Size: 524288000, // 500 MB
			Type: disk.FilesystemDataGUID,
			UUID: disk.FilesystemDataUUID,
			Payload: &disk.Filesystem{
				Type:         "ext4",
				Mountpoint:   "/boot",
				Label:        "boot",
				FSTabOptions: "defaults",
				FSTabFreq:    0,
				FSTabPassNo:  0,
			},
		},
		{
			Size: 2147483648, // 2GiB
			Type: disk.FilesystemDataGUID,
			UUID: disk.RootPartitionUUID,
			Payload: &disk.Filesystem{
				Type:         "ext4",
				Label:        "root",
				Mountpoint:   "/",
				FSTabOptions: "defaults",
				FSTabFreq:    0,
				FSTabPassNo:  0,
			},
		},
	},
}

const (
	store = "/media/scratch/osbuild-store"
)

type imageType struct {
	name string

	platform    platform.Platform
	environment environment.Environment
	workload    workload.Workload

	basePartitionTables distro.BasePartitionTableMap
}

func getRepos(distro, arch string) []rpmmd.RepoConfig {
	distroRepos, err := rpmmd.LoadRepositories([]string{"test/data/"}, distro)
	check(err)
	return distroRepos[arch]
}

func (it *imageType) Manifest() (*manifest.Manifest, []rpmmd.PackageSpec, error) {
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

	img.InstantiateManifest(&m, repos, runner, rng)
	solver := dnfjson.NewSolver("platform:f37", "37", "x86_64", "fedora-37", path.Join(store, "rpmmd"))
	solver.SetDNFJSONPath("./dnf-json")

	// Set cache size to 3 GiB
	solver.SetMaxCacheSize(1 * 1024 * 1024 * 1024)

	pkgs, err := solver.Depsolve(pkgSets)
	check(err)

	return &m, pkgs, nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func write_manifest(bytes []byte) {
	fname := "manifest.json"
	fp, err := os.Create(fname)
	check(err)
	_, err = fp.Write(bytes)
	check(err)
	fmt.Printf("Saved manifest to %s\n", fname)
}

func build(it imageType) {

	m, pkgs, err := it.Manifest()
	check(err)

	pkgSpecs := map[string][]rpmmd.PackageSpec{
		"build": pkgs,
		"os":    pkgs,
	}
	bytes, err := m.Serialize(pkgSpecs)
	check(err)

	write_manifest(bytes)

	outputDir := "./"
	extraEnv := []string{}
	jsonResult := false
	_, err = osbuild.RunOSBuild(bytes, store, outputDir, m.GetExports(), m.GetCheckpoints(), extraEnv, jsonResult, os.Stdout)
	check(err)

	fmt.Println("Done")
}

func main() {
	it := imageType{
		name: "qcow2",
	}

	build(it)
}
