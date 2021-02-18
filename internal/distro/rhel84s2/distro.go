package rhel84s2

import (
	"encoding/json"
	"errors"
	"math/rand"
	"sort"

	"github.com/osbuild/osbuild-composer/internal/distro"
	osbuild "github.com/osbuild/osbuild-composer/internal/osbuild2"

	"github.com/osbuild/osbuild-composer/internal/blueprint"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

const name = "rhel-84"
const modulePlatformID = "platform:el8"

type distribution struct {
	arches        map[string]architecture
	imageTypes    map[string]imageType
	buildPackages []string
}

type architecture struct {
	distro             *distribution
	name               string
	bootloaderPackages []string
	buildPackages      []string
	legacy             string
	uefi               bool
	imageTypes         map[string]imageType
}

type imageType struct {
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

func (a *architecture) Distro() distro.Distro {
	return a.distro
}

func (t *imageType) Arch() distro.Arch {
	return t.arch
}

func (d *distribution) ListArches() []string {
	archs := make([]string, 0, len(d.arches))
	for name := range d.arches {
		archs = append(archs, name)
	}
	sort.Strings(archs)
	return archs
}

func (d *distribution) GetArch(arch string) (distro.Arch, error) {
	a, exists := d.arches[arch]
	if !exists {
		return nil, errors.New("invalid architecture: " + arch)
	}

	return &a, nil
}

func (d *distribution) setArches(arches ...architecture) {
	d.arches = map[string]architecture{}
	for _, a := range arches {
		d.arches[a.name] = architecture{
			distro:             d,
			name:               a.name,
			bootloaderPackages: a.bootloaderPackages,
			buildPackages:      a.buildPackages,
			uefi:               a.uefi,
			imageTypes:         a.imageTypes,
		}
	}
}

func (a *architecture) Name() string {
	return a.name
}

func (a *architecture) ListImageTypes() []string {
	formats := make([]string, 0, len(a.imageTypes))
	for name := range a.imageTypes {
		formats = append(formats, name)
	}
	sort.Strings(formats)
	return formats
}

func (a *architecture) GetImageType(imageType string) (distro.ImageType, error) {
	t, exists := a.imageTypes[imageType]
	if !exists {
		return nil, errors.New("invalid image type: " + imageType)
	}

	return &t, nil
}

func (a *architecture) setImageTypes(imageTypes ...imageType) {
	a.imageTypes = map[string]imageType{}
	for _, it := range imageTypes {
		a.imageTypes[it.name] = imageType{
			arch:             a,
			name:             it.name,
			filename:         it.filename,
			mimeType:         it.mimeType,
			packages:         it.packages,
			excludedPackages: it.excludedPackages,
			enabledServices:  it.enabledServices,
			disabledServices: it.disabledServices,
			defaultTarget:    it.defaultTarget,
			kernelOptions:    it.kernelOptions,
			bootable:         it.bootable,
			rpmOstree:        it.rpmOstree,
			defaultSize:      it.defaultSize,
		}
	}
}

func (t *imageType) Name() string {
	return t.name
}

func (t *imageType) Filename() string {
	return t.filename
}

func (t *imageType) MIMEType() string {
	return t.mimeType
}

func (t *imageType) Size(size uint64) uint64 {
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

func (t *imageType) Packages(bp blueprint.Blueprint) ([]string, []string) {
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

func (t *imageType) BuildPackages() []string {
	packages := append(t.arch.distro.buildPackages, t.arch.buildPackages...)
	if t.rpmOstree {
		packages = append(packages, "rpm-ostree")
	}
	return packages
}

func (t *imageType) Manifest(c *blueprint.Customizations,
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
			Sources:   sources(append(packageSpecs, buildPackageSpecs...)),
		},
	)
}

func (d *distribution) Name() string {
	return name
}

func (d *distribution) ModulePlatformID() string {
	return modulePlatformID
}

func sources(packages []rpmmd.PackageSpec) osbuild.Sources {
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

func (t *imageType) pipelines(c *blueprint.Customizations, options distro.ImageOptions, repos []rpmmd.RepoConfig, packageSpecs, buildPackageSpecs []rpmmd.PackageSpec, rng *rand.Rand) ([]osbuild.Pipeline, error) {
	pipelines := make([]osbuild.Pipeline, 0, 5)

	pipelines = append(pipelines, *t.buildPipeline(repos, buildPackageSpecs))

	if t.rpmOstree {
		// NOTE: Currently all image types in this distro are ostree
		pipelines = append(pipelines, *t.ostreeTreePipeline(repos, packageSpecs, c))
		pipelines = append(pipelines, *t.ostreeCommitPipeline(options))
	}

	pipelines = append(pipelines, *t.containerTreePipeline(repos, buildPackageSpecs, options, c))
	pipelines = append(pipelines, *t.containerPipeline())

	return pipelines, nil
}

func (t *imageType) buildPipeline(repos []rpmmd.RepoConfig, buildPackageSpecs []rpmmd.PackageSpec) *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	p.Name = "build"
	p.Runner = "org.osbuild.rhel84"
	p.AddStage(osbuild.NewRPMStage(t.rpmStageOptions(repos), t.rpmStageInputs(buildPackageSpecs)))
	p.AddStage(osbuild.NewSELinuxStage(t.selinuxStageOptions()))
	return p
}

func (t *imageType) ostreeTreePipeline(repos []rpmmd.RepoConfig, packages []rpmmd.PackageSpec, c *blueprint.Customizations) *osbuild.Pipeline {
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

func (t *imageType) ostreeCommitPipeline(options distro.ImageOptions) *osbuild.Pipeline {
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
			Ref:       options.OSTree.Ref,
			OSVersion: "8.4", // NOTE: Set on image type?
			Parent:    options.OSTree.Parent,
		},
		&osbuild.OSTreeCommitStageInputs{Tree: commitStageInput}),
	)
	return p
}

func (t *imageType) containerTreePipeline(repos []rpmmd.RepoConfig, packages []rpmmd.PackageSpec, options distro.ImageOptions, c *blueprint.Customizations) *osbuild.Pipeline {
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

func (t *imageType) containerPipeline() *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	p.Name = "container"
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

func (t *imageType) rpmStageInputs(specs []rpmmd.PackageSpec) *osbuild.RPMStageInputs {
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

func (t *imageType) ostreePullStageInputs(options distro.ImageOptions) *osbuild.OSTreePullStageInputs {
	pullStageInput := new(osbuild.OSTreePullStageInput)
	pullStageInput.Type = "org.osbuild.tree"
	pullStageInput.Origin = "org.osbuild.pipeline"

	inputRefs := make(map[string]osbuild.OSTreePullStageReference)
	inputRefs["name:ostree-commit"] = osbuild.OSTreePullStageReference{Ref: options.OSTree.Ref}
	pullStageInput.References = inputRefs
	return &osbuild.OSTreePullStageInputs{Commits: pullStageInput}
}

func (t *imageType) rpmStageOptions(repos []rpmmd.RepoConfig) *osbuild.RPMStageOptions {
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

func (t *imageType) selinuxStageOptions() *osbuild.SELinuxStageOptions {
	return &osbuild.SELinuxStageOptions{
		FileContexts: "etc/selinux/targeted/contexts/files/file_contexts",
	}
}

// New creates a new distro object, defining the supported architectures and image types
func New() distro.Distro {
	const GigaByte = 1024 * 1024 * 1024

	edgeOCIImgTypeX86_64 := imageType{
		name:     "rhel-edge-container",
		filename: "rhel84-container.tar",
		mimeType: "application/x-tar",
		packages: []string{
			"redhat-release", // TODO: is this correct for Edge?
			"glibc", "glibc-minimal-langpack", "nss-altfiles",
			"kernel",
			"dracut-config-generic", "dracut-network",
			"basesystem", "bash", "platform-python",
			"shadow-utils", "chrony", "setup", "shadow-utils",
			"sudo", "systemd", "coreutils", "util-linux",
			"curl", "vim-minimal",
			"rpm", "rpm-ostree", "polkit",
			"lvm2", "cryptsetup", "pinentry",
			"e2fsprogs", "dosfstools",
			"keyutils", "gnupg2",
			"attr", "xz", "gzip",
			"firewalld", "iptables",
			"NetworkManager", "NetworkManager-wifi", "NetworkManager-wwan",
			"wpa_supplicant",
			"dnsmasq", "traceroute",
			"hostname", "iproute", "iputils",
			"openssh-clients", "procps-ng", "rootfiles",
			"openssh-server", "passwd",
			"policycoreutils", "policycoreutils-python-utils",
			"selinux-policy-targeted", "setools-console",
			"less", "tar", "rsync",
			"fwupd", "usbguard",
			"bash-completion", "tmux",
			"ima-evm-utils",
			"audit",
			"podman", "container-selinux", "skopeo", "criu",
			"slirp4netns", "fuse-overlayfs",
			"clevis", "clevis-dracut", "clevis-luks",
			"greenboot", "greenboot-grub2", "greenboot-rpm-ostree-grub2", "greenboot-reboot", "greenboot-status",
			// x86 specific
			"grub2", "grub2-efi-x64", "efibootmgr", "shim-x64", "microcode_ctl",
			"iwl1000-firmware", "iwl100-firmware", "iwl105-firmware", "iwl135-firmware",
			"iwl2000-firmware", "iwl2030-firmware", "iwl3160-firmware", "iwl5000-firmware",
			"iwl5150-firmware", "iwl6000-firmware", "iwl6050-firmware", "iwl7260-firmware",
		},
		excludedPackages: []string{
			"rng-tools",
			"subscription-manager",
		},
		enabledServices: []string{
			"NetworkManager.service", "firewalld.service", "sshd.service",
			"greenboot-grub2-set-counter", "greenboot-grub2-set-success", "greenboot-healthcheck",
			"greenboot-rpm-ostree-grub2-check-fallback", "greenboot-status", "greenboot-task-runner",
			"redboot-auto-reboot", "redboot-task-runner",
		},
		rpmOstree: true,
	}
	edgeOCIImgTypeArch64 := imageType{
		name:     "rhel-edge-container",
		filename: "rhel84-container.tar",
		mimeType: "application/x-tar",
		packages: []string{
			"redhat-release", // TODO: is this correct for Edge?
			"glibc", "glibc-minimal-langpack", "nss-altfiles",
			"kernel",
			"dracut-config-generic", "dracut-network",
			"basesystem", "bash", "platform-python",
			"shadow-utils", "chrony", "setup", "shadow-utils",
			"sudo", "systemd", "coreutils", "util-linux",
			"curl", "vim-minimal",
			"rpm", "rpm-ostree", "polkit",
			"lvm2", "cryptsetup", "pinentry",
			"e2fsprogs", "dosfstools",
			"keyutils", "gnupg2",
			"attr", "xz", "gzip",
			"firewalld", "iptables",
			"NetworkManager", "NetworkManager-wifi", "NetworkManager-wwan",
			"wpa_supplicant",
			"dnsmasq", "traceroute",
			"hostname", "iproute", "iputils",
			"openssh-clients", "procps-ng", "rootfiles",
			"openssh-server", "passwd",
			"policycoreutils", "policycoreutils-python-utils",
			"selinux-policy-targeted", "setools-console",
			"less", "tar", "rsync",
			"fwupd", "usbguard",
			"bash-completion", "tmux",
			"ima-evm-utils",
			"audit",
			"podman", "container-selinux", "skopeo", "criu",
			"slirp4netns", "fuse-overlayfs",
			"clevis", "clevis-dracut", "clevis-luks",
			"greenboot", "greenboot-grub2", "greenboot-rpm-ostree-grub2", "greenboot-reboot", "greenboot-status",
			// aarch64 specific
			"grub2-efi-aa64", "efibootmgr", "shim-aa64",
			"iwl7260-firmware",
		},
		excludedPackages: []string{
			"rng-tools",
			"subscription-manager",
		},
		enabledServices: []string{
			"NetworkManager.service", "firewalld.service", "sshd.service",
			"greenboot-grub2-set-counter", "greenboot-grub2-set-success", "greenboot-healthcheck",
			"greenboot-rpm-ostree-grub2-check-fallback", "greenboot-status", "greenboot-task-runner",
			"redboot-auto-reboot", "redboot-task-runner",
		},
		rpmOstree: true,
	}

	r := distribution{
		imageTypes: map[string]imageType{},
		buildPackages: []string{
			"dnf",
			"dosfstools",
			"e2fsprogs",
			"glibc",
			"policycoreutils",
			"python36",
			"python3-iniparse", // dependency of org.osbuild.rhsm stage
			"qemu-img",
			"selinux-policy-targeted",
			"systemd",
			"tar",
			"xfsprogs",
			"xz",

			// for the container
			"httpd",
		},
	}
	x8664 := architecture{
		distro: &r,
		name:   "x86_64",
		bootloaderPackages: []string{
			"dracut-config-generic",
			"grub2-pc",
			"grub2-efi-x64",
			"shim-x64",
		},
		buildPackages: []string{
			"grub2-pc",
		},
		legacy: "i386-pc",
		uefi:   true,
	}
	x8664.setImageTypes(
		edgeOCIImgTypeX86_64,
	)

	aarch64 := architecture{
		distro: &r,
		name:   "aarch64",
		bootloaderPackages: []string{
			"dracut-config-generic",
			"efibootmgr",
			"grub2-efi-aa64",
			"grub2-tools",
			"shim-aa64",
		},
		uefi: true,
	}
	aarch64.setImageTypes(
		edgeOCIImgTypeArch64,
	)

	r.setArches(x8664, aarch64)

	return &r
}
