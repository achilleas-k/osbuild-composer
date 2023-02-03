package main

import (
	"encoding/json"
	"fmt"

	"github.com/osbuild/osbuild-composer/internal/blueprint"
	"github.com/osbuild/osbuild-composer/internal/container"
	"github.com/osbuild/osbuild-composer/internal/distro"
	"github.com/osbuild/osbuild-composer/internal/distro/nu"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	d := nu.New()
	x86, err := d.GetArch("x86_64")
	check(err)
	qcow2, err := x86.GetImageType("qcow2")
	check(err)
	manifest, err := nuManifest(qcow2)
	check(err)

	mj, err := json.MarshalIndent(manifest, "", "  ")
	check(err)
	fmt.Println(string(mj))
}

func nuManifest(it *nu.ImageType) (distro.Manifest, error) {
	var customizations *blueprint.Customizations
	var options distro.ImageOptions
	var repos []rpmmd.RepoConfig
	packageSets := make(map[string][]rpmmd.PackageSpec)
	var containers []container.Spec
	var seed int64

	bob := rpmmd.PackageSpec{
		Name:           "bobby",
		Epoch:          0,
		Version:        "42",
		Release:        "",
		Arch:           "",
		RemoteLocation: "https://example.com/repo/packages/bobby-42.rpm",
		Checksum:       "ffffff",
		Secrets:        "",
		CheckGPG:       false,
		IgnoreSSL:      false,
	}
	packageSets["build"] = []rpmmd.PackageSpec{bob}
	packageSets["os"] = []rpmmd.PackageSpec{bob}

	manifest, err := it.Manifest(customizations, options, repos, packageSets, containers, seed)

	return manifest, err
}
