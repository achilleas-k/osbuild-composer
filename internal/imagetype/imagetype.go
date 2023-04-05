package imagetype

import (
	"fmt"
	"math/rand"

	"github.com/osbuild/osbuild-composer/internal/blueprint"
	"github.com/osbuild/osbuild-composer/internal/common"
	"github.com/osbuild/osbuild-composer/internal/container"
	"github.com/osbuild/osbuild-composer/internal/disk"
	"github.com/osbuild/osbuild-composer/internal/distro"
	"github.com/osbuild/osbuild-composer/internal/environment"
	"github.com/osbuild/osbuild-composer/internal/image"
	"github.com/osbuild/osbuild-composer/internal/manifest"
	"github.com/osbuild/osbuild-composer/internal/platform"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
	"github.com/osbuild/osbuild-composer/internal/workload"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

const (
	// package set names

	// build package set name
	buildPkgsKey = "build"

	// main/common os image package set name
	osPkgsKey = "os"

	// container package set name
	containerPkgsKey = "container"

	// installer package set name
	installerPkgsKey = "installer"

	// blueprint package set name
	blueprintPkgsKey = "blueprint"
)

type packageSetFunc func(t *ImageType) rpmmd.PackageSet

type imageFunc func(workload workload.Workload, t distro.ImageType, customizations *blueprint.Customizations, options distro.ImageOptions, packageSets map[string]rpmmd.PackageSet, containers []container.Spec, rng *rand.Rand) (image.ImageKind, error)

type ImageType struct {
	arch               distro.Arch
	platform           platform.Platform
	environment        environment.Environment
	workload           workload.Workload
	name               string
	nameAliases        []string
	filename           string
	compression        string // TODO: remove from image definition and make it a transport option
	mimeType           string
	packageSets        map[string]packageSetFunc
	defaultImageConfig *distro.ImageConfig
	kernelOptions      string
	defaultSize        uint64
	buildPipelines     []string
	payloadPipelines   []string
	exports            []string
	image              imageFunc

	// bootISO: installable ISO
	bootISO bool
	// rpmOstree: edge/ostree
	rpmOstree bool
	// bootable image
	bootable bool
	// If set to a value, it is preferred over the architecture value
	bootType distro.BootType
	// List of valid arches for the image type
	basePartitionTables distro.BasePartitionTableMap
}

func (t *ImageType) Name() string {
	return t.name
}

func (t *ImageType) Arch() distro.Arch {
	return t.arch
}

func (t *ImageType) Filename() string {
	return t.filename
}

func (t *ImageType) MIMEType() string {
	return t.mimeType
}

func (t *ImageType) OSTreeRef() string {
	return ""
}

func (t *ImageType) Size(size uint64) uint64 {
	// Microsoft Azure requires vhd images to be rounded up to the nearest MB
	if t.name == "vhd" && size%common.MebiByte != 0 {
		size = (size/common.MebiByte + 1) * common.MebiByte
	}
	if size == 0 {
		size = t.defaultSize
	}
	return size
}

func (t *ImageType) BuildPipelines() []string {
	return t.buildPipelines
}

func (t *ImageType) PayloadPipelines() []string {
	return t.payloadPipelines
}

func (t *ImageType) PayloadPackageSets() []string {
	return []string{blueprintPkgsKey}
}

func (t *ImageType) PackageSetsChains() map[string][]string {
	return nil
}

func (t *ImageType) Exports() []string {
	if len(t.exports) > 0 {
		return t.exports
	}
	return []string{"assembler"}
}

// getBootType returns the BootType which should be used for this particular
// combination of architecture and image type.
func (t *ImageType) getBootType() distro.BootType {
	return ""
}

func (t *ImageType) getPartitionTable(
	mountpoints []blueprint.FilesystemCustomization,
	options distro.ImageOptions,
	rng *rand.Rand,
) (*disk.PartitionTable, error) {
	archName := t.arch.Name()

	basePartitionTable, exists := t.basePartitionTables[archName]

	if !exists {
		return nil, fmt.Errorf("no partition table defined for architecture %q for image type %q", archName, t.Name())
	}

	imageSize := t.Size(options.Size)

	lvmify := !t.rpmOstree

	return disk.NewPartitionTable(&basePartitionTable, mountpoints, imageSize, lvmify, nil, rng)
}

func (t *ImageType) getDefaultImageConfig() *distro.ImageConfig {
	return nil
}

func (t *ImageType) PartitionType() string {
	archName := t.arch.Name()
	basePartitionTable, exists := t.basePartitionTables[archName]
	if !exists {
		return ""
	}

	return basePartitionTable.Type
}

func (t *ImageType) Manifest(customizations *blueprint.Customizations,
	options distro.ImageOptions,
	repos []rpmmd.RepoConfig,
	packageSpecs map[string][]rpmmd.PackageSpec,
	containers []container.Spec,
	seed int64) (distro.Manifest, []string, error) {

	bp := &blueprint.Blueprint{Name: "empty blueprint"}
	err := bp.Initialize()
	if err != nil {
		panic("could not initialize empty blueprint: " + err.Error())
	}
	bp.Customizations = customizations

	// the os pipeline filters repos based on the `osPkgsKey` package set, merge the repos which
	// contain a payload package set into the `osPkgsKey`, so those repos are included when
	// building the rpm stage in the os pipeline
	// TODO: roll this into workloads
	mergedRepos := make([]rpmmd.RepoConfig, 0, len(repos))
	for _, repo := range repos {
		for _, pkgsKey := range t.PayloadPackageSets() {
			// If the repo already contains the osPkgsKey, skip
			if slices.Contains(repo.PackageSets, osPkgsKey) {
				break
			}
			if slices.Contains(repo.PackageSets, pkgsKey) {
				repo.PackageSets = append(repo.PackageSets, osPkgsKey)
			}
		}
		mergedRepos = append(mergedRepos, repo)
	}

	repos = mergedRepos
	warnings, err := t.checkOptions(bp.Customizations, options, containers)
	if err != nil {
		return nil, nil, err
	}

	var packageSets map[string]rpmmd.PackageSet
	w := t.workload
	if w == nil {
		cw := &workload.Custom{
			BaseWorkload: workload.BaseWorkload{
				Repos: packageSets[blueprintPkgsKey].Repositories,
			},
			Packages: bp.GetPackagesEx(false),
		}
		if services := bp.Customizations.GetServices(); services != nil {
			cw.Services = services.Enabled
			cw.DisabledServices = services.Disabled
		}
		w = cw
	}

	source := rand.NewSource(seed)
	// math/rand is good enough in this case
	/* #nosec G404 */
	rng := rand.New(source)

	img, err := t.image(w, t, bp.Customizations, options, packageSets, containers, rng)
	if err != nil {
		return nil, nil, err
	}
	manifest := manifest.New()
	_, err = img.InstantiateManifest(&manifest, repos, t.arch.Distro().Runner(), rng)
	if err != nil {
		return nil, nil, err
	}

	ret, err := manifest.Serialize(packageSpecs)
	if err != nil {
		return ret, nil, err
	}
	return ret, warnings, err
}

func (t *ImageType) PackageSets(bp blueprint.Blueprint, options distro.ImageOptions, repos []rpmmd.RepoConfig) map[string][]rpmmd.PackageSet {
	// merge package sets that appear in the image type with the package sets
	// of the same name from the distro and arch
	packageSets := make(map[string]rpmmd.PackageSet)

	for name, getter := range t.packageSets {
		packageSets[name] = getter(t)
	}

	// amend with repository information
	for _, repo := range repos {
		if len(repo.PackageSets) > 0 {
			// only apply the repo to the listed package sets
			for _, psName := range repo.PackageSets {
				ps := packageSets[psName]
				ps.Repositories = append(ps.Repositories, repo)
				packageSets[psName] = ps
			}
		}
	}

	// In case of Cloud API, this method is called before the ostree commit
	// is resolved. Unfortunately, initializeManifest when called for
	// an ostree installer returns an error.
	//
	// Work around this by providing a dummy FetchChecksum to convince the
	// method that it's fine to initialize the manifest. Note that the ostree
	// content has no effect on the package sets, so this is fine.
	//
	// See: https://github.com/osbuild/osbuild-composer/issues/3125
	//
	// TODO: Remove me when it's possible the get the package set chain without
	//       resolving the ostree reference before. Also remove the test for
	//       this workaround
	if t.rpmOstree && t.bootISO && options.OSTree.FetchChecksum == "" {
		options.OSTree.FetchChecksum = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
		logrus.Warn("FIXME: Requesting package sets for iot-installer without a resolved ostree ref. Faking one.")
	}

	// Similar to above, for edge-commit and edge-container, we need to set an
	// ImageRef in order to properly initialize the manifest and package
	// selection.
	options.OSTree.ImageRef = t.OSTreeRef()

	// create a temporary container spec array with the info from the blueprint
	// to initialize the manifest
	containers := make([]container.Spec, len(bp.Containers))
	for idx := range bp.Containers {
		containers[idx] = container.Spec{
			Source:    bp.Containers[idx].Source,
			TLSVerify: bp.Containers[idx].TLSVerify,
			LocalName: bp.Containers[idx].Name,
		}
	}

	_, err := t.checkOptions(bp.Customizations, options, containers)
	if err != nil {
		logrus.Errorf("Initializing the manifest failed for %s (%s/%s): %v", t.Name(), t.arch.Distro().Name(), t.arch.Name(), err)
		return nil
	}

	w := t.workload
	if w == nil {
		cw := &workload.Custom{
			BaseWorkload: workload.BaseWorkload{
				Repos: packageSets[blueprintPkgsKey].Repositories,
			},
			Packages: bp.GetPackagesEx(false),
		}
		if services := bp.Customizations.GetServices(); services != nil {
			cw.Services = services.Enabled
			cw.DisabledServices = services.Disabled
		}
		w = cw
	}

	source := rand.NewSource(0)
	// math/rand is good enough in this case
	/* #nosec G404 */
	rng := rand.New(source)

	img, err := t.image(w, t, bp.Customizations, options, packageSets, containers, rng)
	if err != nil {
		logrus.Errorf("Initializing the manifest failed for %s (%s/%s): %v", t.Name(), t.arch.Distro().Name(), t.arch.Name(), err)
		return nil
	}
	manifest := manifest.New()
	_, err = img.InstantiateManifest(&manifest, repos, t.arch.Distro().Runner(), rng)
	if err != nil {
		logrus.Errorf("Initializing the manifest failed for %s (%s/%s): %v", t.Name(), t.arch.Distro().Name(), t.arch.Name(), err)
		return nil
	}
	return manifest.GetPackageSetChains()
}

// checkOptions checks the validity and compatibility of options and customizations for the image type.
// Returns ([]string, error) where []string, if non-nil, will hold any generated warnings (e.g. deprecation notices).
func (t *ImageType) checkOptions(customizations *blueprint.Customizations, options distro.ImageOptions, containers []container.Spec) ([]string, error) {
	panic("not implemented")
}
