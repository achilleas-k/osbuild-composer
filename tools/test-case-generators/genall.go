package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/osbuild/osbuild-composer/internal/blueprint"
	"github.com/osbuild/osbuild-composer/internal/distro"
	"github.com/osbuild/osbuild-composer/internal/distroregistry"
	"github.com/osbuild/osbuild-composer/internal/dnfjson"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

func makeManifestJob(imgType distro.ImageType, options distro.ImageOptions, repos []rpmmd.RepoConfig, distribution distro.Distro, archName string, seedArg int64) (func(), string) {
	distroName := distribution.Name()
	imgTypeName := imgType.Name()
	u := func(s string) string {
		return strings.Replace(s, "-", "_", -1)
	}
	filename := fmt.Sprintf("%s-%s-%s-boot.json", u(distroName), u(archName), u(imgTypeName))
	workerName := archName + distribution.Name()
	cacheDir := path.Join("/tmp", "rpmmd", workerName)
	job := func() {
		fmt.Printf("Starting job %s\n", filename)
		packageSets := imgType.PackageSets(blueprint.Blueprint{})
		if len(repos) == 0 {
			fmt.Println("No repos")
			return
		}
		packageSpecs, err := depsolve(cacheDir, packageSets, repos, distribution, archName)
		if err != nil {
			fmt.Printf("depsolve failed: %s\n", err.Error())
			return
		}
		if packageSpecs == nil {
			fmt.Println("nil package specs...")
			return
		}
		manifest, err := imgType.Manifest(nil, options, repos, packageSpecs, seedArg)
		if err != nil {
			fmt.Printf("%q failed: %s\n", filename, err)
			return
		}
		saveManifest(manifest, filename)
		fmt.Printf("Finished job %s\n", filename)
	}
	return job, workerName
}

type DistroArchRepoMap map[string]map[string][]rpmmd.RepoConfig

func readRepos() DistroArchRepoMap {
	file := "./tools/test-case-generators/repos.json"
	darm := new(DistroArchRepoMap)
	fp, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	data, err := ioutil.ReadAll(fp)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, darm); err != nil {
		panic(err)
	}
	return *darm
}

func depsolve(cacheDir string, packageSets map[string]rpmmd.PackageSet, repos []rpmmd.RepoConfig, d distro.Distro, arch string) (map[string][]rpmmd.PackageSpec, error) {
	// convert distro repos
	dnfRepos, err := dnfjson.ReposFromRPMMD(repos, arch, d.Releasever())
	if err != nil {
		return nil, err
	}
	solver := dnfjson.NewSolver(d.ModulePlatformID(), arch, cacheDir)
	packageSpecSets := make(map[string][]rpmmd.PackageSpec)
	for name, packages := range packageSets {
		results, err := solver.Depsolve(packages, dnfRepos)
		if err != nil {
			fmt.Printf("Could not depsolve: %s", err.Error())
			return nil, err
		}
		// NOTE: we only depsolve one package set at a time for now
		result := results[0]
		packageSpecSets[name] = dnfjson.DepsToRPMMD(result.Dependencies, repos)
	}
	return packageSpecSets, nil
}

func saveManifest(manifest distro.Manifest, filename string) {
	b, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		panic(err)
	}
	fp, err := os.Create("test/data/manifests.noinfo/" + filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	fp.Write(b)
}

type Worker struct {
	name string
	done chan bool
}

func (w *Worker) runJobs(q chan func()) {
	c := 0
	go func() {
		for j := range q {
			j()
			c++
			fmt.Printf("Worker %s finished %d jobs\n", w.name, c)
		}
		w.done <- true
	}()
}

func (w *Worker) addJob(j func(), q chan func()) {
	q <- j
}

func initWorkers(names []string) map[string]*Worker {
	workers := make(map[string]*Worker)
	for _, name := range names {
		// run through once to add one entry per name
		workers[name] = nil
	}
	for name := range workers {
		w := new(Worker)
		w.name = name
		w.done = make(chan bool, 1)
		workers[name] = w
	}
	return workers
}
func initWorkersN(n int) map[int]*Worker {
	workers := make(map[int]*Worker)
	for idx := 0; idx < n; idx++ {
		w := new(Worker)
		name := fmt.Sprint(idx)
		w.name = name
		w.done = make(chan bool, 1)
		workers[idx] = w
	}
	return workers
}

func main() {
	seedArg := int64(0)
	darm := readRepos()
	distros := distroregistry.NewDefault()
	wNames := make([]string, 0)
	jobs := make([]func(), 0)
	fmt.Println("Collecting jobs")
	for _, distroName := range distros.List() {
		distribution := distros.GetDistro(distroName)
		for _, archName := range distribution.ListArches() {
			if archName == "s390x" {
				fmt.Println("Skipping s390x")
				continue
			}
			arch, err := distribution.GetArch(archName)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			repos := darm[distroName][archName]
			if len(repos) == 0 {
				fmt.Printf("No repos for %s %s. Skipping\n", distroName, archName)
				continue
			}

			for _, imgTypeName := range arch.ListImageTypes() {
				imgType, err := arch.GetImageType(imgTypeName)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				if imgType.Name() == "edge-installer" || imgType.Name() == "edge-simplified-installer" || imgType.Name() == "edge-raw-image" {
					continue
				}
				options := distro.ImageOptions{Size: imgType.Size(0)}
				job, wName := makeManifestJob(imgType, options, repos, distribution, archName, seedArg)
				jobs = append(jobs, job)
				wNames = append(wNames, wName)
			}
		}
	}

	fmt.Printf("Collected %d jobs\n", len(jobs))
	n := 32
	workers := initWorkersN(n)
	fmt.Printf("Initialised %d workers\n", len(workers))
	fmt.Println("Enqueueing jobs")
	q := make(chan func(), len(jobs))
	for idx := range jobs {
		j := jobs[idx]
		workers[idx%n].addJob(j, q)
	}

	fmt.Println("Starting workers")
	for name := range workers {
		workers[name].runJobs(q)
	}
	fmt.Println("Finalizing workers")
	close(q)
	for name := range workers {
		<-workers[name].done
	}
	fmt.Printf("ALL DONE (%d)\n", len(jobs))
}
