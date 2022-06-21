package distro

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"

	"github.com/osbuild/osbuild-composer/internal/blueprint"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

type externalDistroInfo struct {
	Name             string `json:"name"`
	Releasever       string `json:"releasever"`
	ModulePlatformID string `json:"module_platform_id"`
	OSTreeRef        string `json:"ostree_ref"`

	Arches map[string]externalArchInfo `json:"arches"`
}

type externalArchInfo struct {
	ImageTypes map[string]externalImageTypeInfo `json:"image_types"`
}

type externalImageTypeInfo struct {
	Filename      string `json:"filename"`
	MIMEType      string `json:"mimetype"`
	OSTreeRef     string `json:"ostree_ref"`
	Size          uint64 `json:"size"`
	PartitionType string `json:"partition_type"`

	PackageSets        map[string]externalPackageSet `json:"package_sets"`
	PayloadPackageSets []string                      `json:"payload_package_sets"`

	BuildPipelines   []string `json:"build_pipelines"`
	PayloadPipelines []string `json:"payload_pipelines"`
	Exports          []string `json:"exports"`
}

type externalPackageSet struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

// externalDistro distro definition.  Implements the Distro interface.  Each
// method calls out to a single binary with appropriate command line arguments.
// The externalDistro distro binary should implement the osbuild-composer
// distro definition specification.
type externalDistro struct {
	cmd              string
	name             string
	releasever       string
	modulePlatformID string
	ostreeRef        string

	arches map[string]externalArch
}

func NewExternal(cmd string) (Distro, error) {
	infoRaw, err := run(cmd, "get-info")
	if err != nil {
		return nil, err
	}

	var edi *externalDistroInfo
	decoder := json.NewDecoder(bytes.NewBuffer(infoRaw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(edi); err != nil {
		return nil, err
	}

	distro := &externalDistro{
		cmd:              cmd,
		name:             edi.Name,
		releasever:       edi.Releasever,
		modulePlatformID: edi.ModulePlatformID,
		ostreeRef:        edi.OSTreeRef,
	}

	for archName, archInfo := range edi.Arches {
		arch := newExternalArch(archName, archInfo)
		arch.distro = distro
		distro.arches[archName] = arch

	}
	return distro, nil
}

// Returns the name of the distro.
func (d externalDistro) Name() string {
	return d.name
}

// Returns the release version of the distro. This is used in repo
// files on the host system and required for the subscription support.
func (d externalDistro) Releasever() string {
	return d.releasever
}

// Returns the module platform id of the distro. This is used by DNF
// for modularity support.
func (d externalDistro) ModulePlatformID() string {
	return d.modulePlatformID
}

// Returns the ostree reference template
func (d externalDistro) OSTreeRef() string {
	return d.ostreeRef
}

// Returns a sorted list of the names of the architectures this distro
// supports.
func (d externalDistro) ListArches() []string {
	names := make([]string, len(d.arches))
	var idx int
	for name := range d.arches {
		names[idx] = name
		idx++
	}
	return names
}

// Returns an object representing the given architecture as support
// by this distro.
func (d externalDistro) GetArch(name string) (Arch, error) {
	arch, valid := d.arches[name]
	if !valid {
		return nil, errors.New("invalid architecture: " + name)
	}
	return arch, nil
}

func run(distroCmd string, args ...string) ([]byte, error) {
	cmd := exec.Command(distroCmd, args...)
	cmd.Stderr = os.Stderr
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return stdout.Bytes(), nil
}

// An Arch represents a given distribution's support for a given architecture.
type externalArch struct {
	name string

	imageTypes map[string]externalImageType

	distro Distro
}

func newExternalArch(name string, archInfo externalArchInfo) externalArch {
	arch := externalArch{name: name}
	for imgTypeName, imgTypeInfo := range archInfo.ImageTypes {
		imgType := newExternalImageType(imgTypeName, imgTypeInfo)
		imgType.arch = arch
		arch.imageTypes[imgTypeName] = imgType
	}
	return arch
}

// Returns the name of the architecture.
func (a externalArch) Name() string {
	return a.name
}

// Returns a sorted list of the names of the image types this architecture
// supports.
func (a externalArch) ListImageTypes() []string {
	imgTypes := make([]string, len(a.imageTypes))
	var idx int
	for name := range a.imageTypes {
		imgTypes[idx] = name
		idx++
	}
	return imgTypes
}

// Returns an object representing a given image format for this architecture,
// on this distro.
func (a externalArch) GetImageType(name string) (ImageType, error) {
	imgType, valid := a.imageTypes[name]
	if !valid {
		return nil, errors.New("invalid image type: " + name)
	}
	return imgType, nil
}

// Returns the parent distro
func (a externalArch) Distro() Distro {
	return a.distro
}

type externalImageType struct {
	name          string
	filename      string
	mimeType      string
	osTreeRef     string
	size          uint64
	partitionType string

	payloadPackageSets []string

	buildPipelines   []string
	payloadPipelines []string
	exports          []string

	arch Arch
}

func newExternalImageType(name string, imgTypeInfo externalImageTypeInfo) externalImageType {
	imgType := externalImageType{
		name:               name,
		filename:           imgTypeInfo.Filename,
		mimeType:           imgTypeInfo.MIMEType,
		osTreeRef:          imgTypeInfo.OSTreeRef,
		size:               imgTypeInfo.Size,
		partitionType:      imgTypeInfo.PartitionType,
		payloadPackageSets: imgTypeInfo.PayloadPackageSets,
		buildPipelines:     imgTypeInfo.BuildPipelines,
		payloadPipelines:   imgTypeInfo.PayloadPipelines,
		exports:            imgTypeInfo.Exports,
	}
	return imgType
}

// Returns the name of the image type.
func (it externalImageType) Name() string {
	return it.name
}

// Returns the parent architecture
func (it externalImageType) Arch() Arch {
	return it.arch
}

// Returns the canonical filename for the image type.
func (it externalImageType) Filename() string {
	return it.filename
}

// Retrns the MIME-type for the image type.
func (it externalImageType) MIMEType() string {
	return it.mimeType
}

// Returns the default OSTree ref for the image type.
func (it externalImageType) OSTreeRef() string {
	return it.osTreeRef
}

// Returns the proper image size for a given output format. If the input size
// is 0 the default value for the format will be returned.
func (it externalImageType) Size(size uint64) uint64 {
	return it.size
}

// Returns the corresponding partion type ("gpt", "dos") or "" the image type
// has no partition table. Only support for RHEL 8.5+
func (it externalImageType) PartitionType() string {
	return it.partitionType
}

// Returns the names of the pipelines that set up the build environment (buildroot).
func (it externalImageType) BuildPipelines() []string {
	return it.buildPipelines
}

// Returns the names of the pipelines that create the image.
func (it externalImageType) PayloadPipelines() []string {
	return it.payloadPipelines
}

// Returns the package set names safe to install custom packages via custom repositories.
func (it externalImageType) PayloadPackageSets() []string {
	return it.payloadPackageSets
}

// Returns named arrays of package set names which should be depsolved in a chain.
func (it externalImageType) PackageSetsChains() map[string][]string {
	return nil
}

// Returns the names of the stages that will produce the build output.
func (it externalImageType) Exports() []string {
	return it.exports
}

// Returns the sets of packages to include and exclude when building the image.
// Indexed by a string label. How each set is labeled and used depends on the
// image type.
func (it externalImageType) PackageSets(bp blueprint.Blueprint, repos []rpmmd.RepoConfig) map[string][]rpmmd.PackageSet {
	return nil
}

// Returns an osbuild manifest, containing the sources and pipeline necessary
// to build an image, given output format with all packages and customizations
// specified in the given blueprint. The packageSpecSets must be labelled in
// the same way as the originating PackageSets.
func (it externalImageType) Manifest(b *blueprint.Customizations, options ImageOptions, repos []rpmmd.RepoConfig, packageSpecSets map[string][]rpmmd.PackageSpec, seed int64) (Manifest, error) {
	return nil, nil
}
