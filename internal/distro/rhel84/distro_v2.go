package rhel84

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/osbuild/osbuild-composer/internal/distro"
	osbuild "github.com/osbuild/osbuild-composer/internal/osbuild2"

	"github.com/osbuild/osbuild-composer/internal/blueprint"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

type imageTypeS2 struct {
	arch             *architecture
	name             string
	filename         string
	mimeType         string
	packages         []string
	excludedPackages []string
	enabledServices  []string
	disabledServices []string
	defaultTarget    string
	kernelOptions    string
	bootable         bool
	rpmOstree        bool
	defaultSize      uint64
}

func (t *imageTypeS2) Arch() distro.Arch {
	return t.arch
}

func (t *imageTypeS2) Name() string {
	return t.name
}

func (t *imageTypeS2) Filename() string {
	return t.filename
}

func (t *imageTypeS2) MIMEType() string {
	return t.mimeType
}

func (t *imageTypeS2) OSTreeRef() string {
	if t.rpmOstree {
		return fmt.Sprintf(ostreeRef, t.arch.name)
	}
	return ""
}

func (t *imageTypeS2) Size(size uint64) uint64 {
	const MegaByte = 1024 * 1024
	// Microsoft Azure requires vhd images to be rounded up to the nearest MB
	if t.name == "vhd" && size%MegaByte != 0 {
		size = (size/MegaByte + 1) * MegaByte
	}
	if size == 0 {
		size = t.defaultSize
	}
	return size
}

func (t *imageTypeS2) Packages(bp blueprint.Blueprint) ([]string, []string) {
	packages := append(t.packages, bp.GetPackages()...)
	timezone, _ := bp.Customizations.GetTimezoneSettings()
	if timezone != nil {
		packages = append(packages, "chrony")
	}
	if t.bootable {
		packages = append(packages, t.arch.bootloaderPackages...)
	}

	return packages, t.excludedPackages
}

func (t *imageTypeS2) BuildPackages() []string {
	packages := append(t.arch.distro.buildPackages, t.arch.buildPackages...)
	if t.rpmOstree {
		packages = append(packages, "rpm-ostree")
	}
	//  NOTE(akoutsou) 1to2t: for the container to serve the commit; this should be its
	//  own package set
	packages = append(packages, "httpd")
	return packages
}

func (t *imageTypeS2) Exports() []string {
	return []string{"assembler"}
}

func (t *imageTypeS2) Manifest(c *blueprint.Customizations,
	options distro.ImageOptions,
	repos []rpmmd.RepoConfig,
	packageSpecs,
	buildPackageSpecs []rpmmd.PackageSpec,
	seed int64) (distro.Manifest, error) {
	source := rand.NewSource(seed)
	rng := rand.New(source)
	pipelines, err := t.pipelines(c, options, repos, packageSpecs, buildPackageSpecs, rng)
	if err != nil {
		return distro.Manifest{}, err
	}

	return json.Marshal(
		osbuild.Manifest{
			Version:   "2",
			Pipelines: pipelines,
			Sources:   t.sources(append(packageSpecs, buildPackageSpecs...)),
		},
	)
}

func (t *imageTypeS2) sources(packages []rpmmd.PackageSpec) osbuild.Sources {
	source := &osbuild.CurlSource{
		Items: make(map[string]osbuild.CurlSourceItem),
	}
	for _, pkg := range packages {
		item := new(osbuild.URLWithSecrets)
		item.URL = pkg.RemoteLocation
		if pkg.Secrets == "org.osbuild.rhsm" {
			item.Secrets = &osbuild.URLSecrets{
				Name: "org.osbuild.rhsm",
			}
		}
		source.Items[pkg.Checksum] = item
	}
	return osbuild.Sources{
		"org.osbuild.curl": source,
	}
}

func (t *imageTypeS2) pipelines(c *blueprint.Customizations, options distro.ImageOptions, repos []rpmmd.RepoConfig, packageSpecs, buildPackageSpecs []rpmmd.PackageSpec, rng *rand.Rand) ([]osbuild.Pipeline, error) {
	pipelines := make([]osbuild.Pipeline, 0, 5)

	pipelines = append(pipelines, *t.buildPipeline(repos, buildPackageSpecs))

	if t.rpmOstree {
		// NOTE(akoutsou) 1to2t: Currently all images of type imageTypeS2 are ostree
		pipelines = append(pipelines, *t.ostreeTreePipeline(repos, packageSpecs, c))
		pipelines = append(pipelines, *t.ostreeCommitPipeline(options))
	}

	// NOTE(akoutsou) 1to2t: the container temporarily uses the build package
	// set but eventually it should get its own defined in the image type
	pipelines = append(pipelines, *t.containerTreePipeline(repos, buildPackageSpecs, options, c))
	pipelines = append(pipelines, *t.containerPipeline())

	return pipelines, nil
}

func (t *imageTypeS2) buildPipeline(repos []rpmmd.RepoConfig, buildPackageSpecs []rpmmd.PackageSpec) *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	p.Name = "build"
	p.Runner = "org.osbuild.rhel84"
	p.AddStage(osbuild.NewRPMStage(t.rpmStageOptions(repos), t.rpmStageInputs(buildPackageSpecs)))
	p.AddStage(osbuild.NewSELinuxStage(t.selinuxStageOptions()))
	return p
}

func (t *imageTypeS2) ostreeTreePipeline(repos []rpmmd.RepoConfig, packages []rpmmd.PackageSpec, c *blueprint.Customizations) *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	p.Name = "ostree-tree"
	p.Build = "name:build"
	p.AddStage(osbuild.NewRPMStage(t.rpmStageOptions(repos), t.rpmStageInputs(packages)))
	language, _ := c.GetPrimaryLocale()
	if language != nil {
		p.AddStage(osbuild.NewLocaleStage(&osbuild.LocaleStageOptions{Language: *language}))
	} else {
		p.AddStage(osbuild.NewLocaleStage(&osbuild.LocaleStageOptions{Language: "en_US"}))
	}
	p.AddStage(osbuild.NewSELinuxStage(t.selinuxStageOptions()))
	p.AddStage(osbuild.NewRPMOSTreePrepTreeStage(&osbuild.RPMOSTreePrepTreeStageOptions{
		EtcGroupMembers: []string{
			// NOTE: We may want to make this configurable.
			"wheel", "docker",
		},
	}))
	return p
}

func (t *imageTypeS2) ostreeCommitPipeline(options distro.ImageOptions) *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	p.Name = "ostree-commit"
	p.Build = "name:build"
	p.AddStage(osbuild.NewOSTreeInitStage(&osbuild.OSTreeInitStageOptions{Path: "/repo"}))

	commitStageInput := new(osbuild.OSTreeCommitStageInput)
	commitStageInput.Type = "org.osbuild.tree"
	commitStageInput.Origin = "org.osbuild.pipeline"
	commitStageInput.References = osbuild.OSTreeCommitStageReferences{"name:ostree-tree"}

	p.AddStage(osbuild.NewOSTreeCommitStage(
		&osbuild.OSTreeCommitStageOptions{
			Ref:       t.OSTreeRef(),
			OSVersion: "8.4", // NOTE: Set on image type?
			Parent:    options.OSTree.Parent,
		},
		&osbuild.OSTreeCommitStageInputs{Tree: commitStageInput}),
	)
	return p
}

func (t *imageTypeS2) containerTreePipeline(repos []rpmmd.RepoConfig, packages []rpmmd.PackageSpec, options distro.ImageOptions, c *blueprint.Customizations) *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	p.Name = "container-tree"
	p.Build = "name:build"
	p.AddStage(osbuild.NewRPMStage(t.rpmStageOptions(repos), t.rpmStageInputs(packages)))
	language, _ := c.GetPrimaryLocale()
	if language != nil {
		p.AddStage(osbuild.NewLocaleStage(&osbuild.LocaleStageOptions{Language: *language}))
	} else {
		p.AddStage(osbuild.NewLocaleStage(&osbuild.LocaleStageOptions{Language: "en_US"}))
	}
	p.AddStage(osbuild.NewOSTreeInitStage(&osbuild.OSTreeInitStageOptions{Path: "/var/www/html/repo"}))

	p.AddStage(osbuild.NewOSTreePullStage(
		&osbuild.OSTreePullStageOptions{Repo: "/var/www/html/repo"},
		t.ostreePullStageInputs(options),
	))
	return p
}

func (t *imageTypeS2) containerPipeline() *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	// NOTE(akoutsou) 1to2t: final pipeline should always be named "assembler"
	p.Name = "assembler"
	p.Build = "name:build"
	options := &osbuild.OCIArchiveStageOptions{
		Architecture: t.arch.Name(),
		Filename:     t.Filename(),
		Config: &osbuild.OCIArchiveConfig{
			Cmd:          []string{"httpd", "-D", "FOREGROUND"},
			ExposedPorts: []string{"80"},
		},
	}
	baseInput := new(osbuild.OCIArchiveStageInput)
	baseInput.Type = "org.osbuild.tree"
	baseInput.Origin = "org.osbuild.pipeline"
	baseInput.References = []string{"name:container-tree"}
	inputs := &osbuild.OCIArchiveStageInputs{Base: baseInput}
	p.AddStage(osbuild.NewOCIArchiveStage(options, inputs))
	return p
}

func (t *imageTypeS2) rpmStageInputs(specs []rpmmd.PackageSpec) *osbuild.RPMStageInputs {
	stageInput := new(osbuild.RPMStageInput)
	stageInput.Type = "org.osbuild.files"
	stageInput.Origin = "org.osbuild.source"
	stageInput.References = pkgRefs(specs)
	return &osbuild.RPMStageInputs{Packages: stageInput}
}

func pkgRefs(specs []rpmmd.PackageSpec) osbuild.RPMStageReferences {
	refs := make([]string, len(specs))
	for idx, pkg := range specs {
		refs[idx] = pkg.Checksum
	}
	return refs
}

func (t *imageTypeS2) ostreePullStageInputs(options distro.ImageOptions) *osbuild.OSTreePullStageInputs {
	pullStageInput := new(osbuild.OSTreePullStageInput)
	pullStageInput.Type = "org.osbuild.ostree"
	pullStageInput.Origin = "org.osbuild.pipeline"

	inputRefs := make(map[string]osbuild.OSTreePullStageReference)
	inputRefs["name:ostree-commit"] = osbuild.OSTreePullStageReference{Ref: t.OSTreeRef()}
	pullStageInput.References = inputRefs
	return &osbuild.OSTreePullStageInputs{Commits: pullStageInput}
}

func (t *imageTypeS2) rpmStageOptions(repos []rpmmd.RepoConfig) *osbuild.RPMStageOptions {
	var gpgKeys []string
	for _, repo := range repos {
		if repo.GPGKey == "" {
			continue
		}
		gpgKeys = append(gpgKeys, repo.GPGKey)
	}

	return &osbuild.RPMStageOptions{
		GPGKeys: gpgKeys,
		Exclude: &osbuild.Exclude{
			// NOTE: Make configurable?
			Docs: true,
		},
	}
}

func (t *imageTypeS2) selinuxStageOptions() *osbuild.SELinuxStageOptions {
	return &osbuild.SELinuxStageOptions{
		FileContexts: "etc/selinux/targeted/contexts/files/file_contexts",
	}
}
