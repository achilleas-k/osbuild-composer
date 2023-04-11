package main

import (
	"fmt"
	"os"
	"path/filepath"

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

func build(it imageType) {
	m, err := it.Manifest()
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
