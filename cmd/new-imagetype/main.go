package main

import (
	"fmt"
	"os"

	"github.com/osbuild/osbuild-composer/internal/distro"
	"github.com/osbuild/osbuild-composer/internal/environment"
	"github.com/osbuild/osbuild-composer/internal/manifest"
	"github.com/osbuild/osbuild-composer/internal/osbuild"
	"github.com/osbuild/osbuild-composer/internal/platform"
	"github.com/osbuild/osbuild-composer/internal/workload"
)

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

func (it *imageType) Manifest() (*manifest.Manifest, error) {
	m := manifest.New()
	return &m, nil
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
