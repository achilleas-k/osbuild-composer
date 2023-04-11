package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/osbuild/osbuild-composer/internal/dnfjson"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

const (
	store  = "/media/scratch/osbuild-store"
	source = "/home/achilleas/projects/osbuild/osbuild-composer"
)

func getRepos(distro, arch string) []rpmmd.RepoConfig {
	distroRepos, err := rpmmd.LoadRepositories([]string{filepath.Join(source, "test/data/")}, distro)
	check(err)
	return distroRepos[arch]
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

func depsolveJob(chains map[string][]rpmmd.PackageSet) map[string][]rpmmd.PackageSpec {
	solver := dnfjson.NewSolver("platform:f37", "37", "x86_64", "fedora-37", path.Join(store, "rpmmd"))
	solver.SetDNFJSONPath(filepath.Join(source, "./dnf-json"))

	// Set cache size to 3 GiB
	solver.SetMaxCacheSize(1 * 1024 * 1024 * 1024)

	solved := make(map[string][]rpmmd.PackageSpec, len(chains))
	for name, chain := range chains {
		pkgs, err := solver.Depsolve(chain)
		check(err)
		solved[name] = pkgs
	}
	return solved
}

func build(it imageType) {
	m, err := it.Manifest(depsolveJob)
	check(err)

	bytes, err := m.Serialize(nil)
	check(err)

	write_manifest(bytes)

	// outputDir := "./"
	// extraEnv := []string{}
	// jsonResult := false
	// _, err = osbuild.RunOSBuild(bytes, store, outputDir, m.GetExports(), m.GetCheckpoints(), extraEnv, jsonResult, os.Stdout)
	// check(err)

	fmt.Println("Done")
}

func main() {
	it := imageType{
		name: "qcow2",
	}

	build(it)
}
