package rhel84

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/osbuild/osbuild-composer/internal/crypt"
	"github.com/osbuild/osbuild-composer/internal/distro"
	osbuild "github.com/osbuild/osbuild-composer/internal/osbuild2"

	"github.com/osbuild/osbuild-composer/internal/blueprint"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

type ImageTypeS2 struct {
	arch                *architecture
	name                string
	filename            string
	mimeType            string
	packageSets         map[string][]string
	excludedPackageSets map[string][]string
	enabledServices     []string
	disabledServices    []string
	defaultTarget       string
	kernelOptions       string
	bootable            bool
	bootISO             bool
	rpmOstree           bool
	defaultSize         uint64
	depsolve            solver
	blueprint           *blueprint.Blueprint
}

func (t *ImageTypeS2) Arch() distro.Arch {
	return t.arch
}

func (t *ImageTypeS2) Name() string {
	return t.name
}

func (t *ImageTypeS2) Filename() string {
	return t.filename
}

func (t *ImageTypeS2) MIMEType() string {
	return t.mimeType
}

func (t *ImageTypeS2) OSTreeRef() string {
	if t.rpmOstree {
		return fmt.Sprintf(ostreeRef, t.arch.name)
	}
	return ""
}

func (t *ImageTypeS2) Size(size uint64) uint64 {
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

func (t *ImageTypeS2) Packages(bp blueprint.Blueprint) ([]string, []string) {
	// NOTE(akoutsou) 1to2t: rhel-edge-container image types perform their own
	// depsolving while creating the Manifest.  Returning empty string slices
	// to avoid unnecessary depsolving.
	return []string{}, []string{}
}

func (t *ImageTypeS2) BuildPackages() []string {
	// NOTE(akoutsou) 1to2t: rhel-edge-container image types perform their own
	// depsolving while creating the Manifest.  Returning empty string slice to
	// avoid unnecessary depsolving.
	return []string{}
}

func (t *ImageTypeS2) Exports() []string {
	return []string{"assembler"}
}

func (t *ImageTypeS2) DepsolvePackageSets() (map[string][]rpmmd.PackageSpec, map[string]string, error) {
	if t.depsolve == nil {
		return map[string][]rpmmd.PackageSpec{}, map[string]string{}, nil
	}

	pkgSpecs := make(map[string][]rpmmd.PackageSpec, len(t.packageSets)+1)
	var checksums map[string]string
	for name, pkgSet := range t.packageSets {
		if name == "commit" {
			// the main package set to be delivered
			// include blueprint and bootloader
			if bp := t.blueprint; bp != nil {
				pkgSet = append(pkgSet, bp.GetPackages()...)
				if timezone, _ := bp.Customizations.GetTimezoneSettings(); timezone != nil {
					pkgSet = append(pkgSet, "chrony")
				}
			}
			if t.bootable {
				pkgSet = append(pkgSet, t.arch.bootloaderPackages...)
			}
		}
		specs, csums, err := t.depsolve(pkgSet, t.excludedPackageSets[name])
		if err != nil {
			return nil, nil, err
		}
		if name == "commit" {
			checksums = csums
		}
		pkgSpecs[name] = specs
	}

	buildPackages := append(t.arch.distro.buildPackages, t.arch.buildPackages...)
	if t.rpmOstree {
		buildPackages = append(buildPackages, "rpm-ostree")
	}
	if t.bootISO {
		buildPackages = append(buildPackages, "lorax")
	}
	buildPackageSpecs, _, err := t.depsolve(buildPackages, nil)
	if err != nil {
		return nil, nil, err
	}
	pkgSpecs["build"] = buildPackageSpecs

	return pkgSpecs, checksums, nil
}

func (t *ImageTypeS2) Manifest(c *blueprint.Customizations,
	options distro.ImageOptions,
	repos []rpmmd.RepoConfig,
	packageSpecs,
	buildPackageSpecs []rpmmd.PackageSpec,
	seed int64) (distro.Manifest, error) {

	source := rand.NewSource(seed)
	rng := rand.New(source)

	// NOTE(akoutsou) 1to2t: package specs coming from the arguments should be
	// empty, so we depsolve them ourselves
	packageSetsSpecs, _, err := t.DepsolvePackageSets()
	if err != nil {
		return nil, err
	}

	pipelines, err := t.pipelines(c, options, repos, packageSetsSpecs, rng)
	if err != nil {
		return distro.Manifest{}, err
	}

	allPackageSpecs := make([]rpmmd.PackageSpec, 0)
	// flatten all package specs
	for _, pkgSpecs := range packageSetsSpecs {
		allPackageSpecs = append(allPackageSpecs, pkgSpecs...)
	}
	return json.Marshal(
		osbuild.Manifest{
			Version:   "2",
			Pipelines: pipelines,
			Sources:   t.sources(allPackageSpecs),
		},
	)
}

func (t *ImageTypeS2) sources(packages []rpmmd.PackageSpec) osbuild.Sources {
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

func (t *ImageTypeS2) pipelines(c *blueprint.Customizations, options distro.ImageOptions, repos []rpmmd.RepoConfig, packageSetsSpecs map[string][]rpmmd.PackageSpec, rng *rand.Rand) ([]osbuild.Pipeline, error) {

	if kernelOpts := c.GetKernel(); kernelOpts.Append != "" && t.rpmOstree {
		return nil, fmt.Errorf("kernel boot parameter customizations are not supported for ostree types")
	}

	pipelines := make([]osbuild.Pipeline, 0)

	pipelines = append(pipelines, *t.buildPipeline(repos, packageSetsSpecs["build"]))

	if t.rpmOstree {
		// NOTE(akoutsou) 1to2t: Currently all images of type imageTypeS2 are ostree
		treePipeline, err := t.ostreeTreePipeline(repos, packageSetsSpecs["commit"], c)
		if err != nil {
			return nil, err
		}
		pipelines = append(pipelines, *treePipeline)
		pipelines = append(pipelines, *t.ostreeCommitPipeline(options))
	}

	if t.bootISO {
		pipelines = append(pipelines, *t.anacondaTreePipeline(repos, packageSetsSpecs["installer"], options, c))
		pipelines = append(pipelines, *t.bootISOTreePipeline())
		pipelines = append(pipelines, *t.bootISOPipeline())
	} else {
		pipelines = append(pipelines, *t.containerTreePipeline(repos, packageSetsSpecs["container"], options, c))
		pipelines = append(pipelines, *t.containerPipeline())
	}

	return pipelines, nil
}

func (t *ImageTypeS2) buildPipeline(repos []rpmmd.RepoConfig, buildPackageSpecs []rpmmd.PackageSpec) *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	p.Name = "build"
	p.Runner = "org.osbuild.rhel84"
	p.AddStage(osbuild.NewRPMStage(t.rpmStageOptions(repos), t.rpmStageInputs(buildPackageSpecs)))
	p.AddStage(osbuild.NewSELinuxStage(t.selinuxStageOptions()))
	return p
}

func (t *ImageTypeS2) ostreeTreePipeline(repos []rpmmd.RepoConfig, packages []rpmmd.PackageSpec, c *blueprint.Customizations) (*osbuild.Pipeline, error) {
	p := new(osbuild.Pipeline)
	p.Name = "ostree-tree"
	p.Build = "name:build"
	p.AddStage(osbuild.NewRPMStage(t.rpmStageOptions(repos), t.rpmStageInputs(packages)))
	language, keyboard := c.GetPrimaryLocale()
	if language != nil {
		p.AddStage(osbuild.NewLocaleStage(&osbuild.LocaleStageOptions{Language: *language}))
	} else {
		p.AddStage(osbuild.NewLocaleStage(&osbuild.LocaleStageOptions{Language: "en_US.UTF-8"}))
	}
	if keyboard != nil {
		p.AddStage(osbuild.NewKeymapStage(&osbuild.KeymapStageOptions{Keymap: *keyboard}))
	}
	if hostname := c.GetHostname(); hostname != nil {
		p.AddStage(osbuild.NewHostnameStage(&osbuild.HostnameStageOptions{Hostname: *hostname}))
	}

	timezone, ntpServers := c.GetTimezoneSettings()
	if timezone != nil {
		p.AddStage(osbuild.NewTimezoneStage(&osbuild.TimezoneStageOptions{Zone: *timezone}))
	} else {
		p.AddStage(osbuild.NewTimezoneStage(&osbuild.TimezoneStageOptions{Zone: "America/New_York"}))
	}

	if len(ntpServers) > 0 {
		p.AddStage(osbuild.NewChronyStage(&osbuild.ChronyStageOptions{Timeservers: ntpServers}))
	}

	if groups := c.GetGroups(); len(groups) > 0 {
		p.AddStage(osbuild.NewGroupsStage(t.groupStageOptions(groups)))
	}

	if users := c.GetUsers(); len(users) > 0 {
		options, err := t.userStageOptions(users)
		if err != nil {
			return nil, err
		}
		p.AddStage(osbuild.NewUsersStage(options))
	}

	if services := c.GetServices(); services != nil || t.enabledServices != nil || t.disabledServices != nil || t.defaultTarget != "" {
		p.AddStage(osbuild.NewSystemdStage(t.systemdStageOptions(t.enabledServices, t.disabledServices, services, t.defaultTarget)))
	}

	if firewall := c.GetFirewall(); firewall != nil {
		p.AddStage(osbuild.NewFirewallStage(t.firewallStageOptions(firewall)))
	}

	if !t.bootISO {
		p.AddStage(osbuild.NewSELinuxStage(t.selinuxStageOptions()))
	}

	// These are the current defaults for the sysconfig stage. This can be changed to be image type exclusive if different configs are needed.
	p.AddStage(osbuild.NewSysconfigStage(&osbuild.SysconfigStageOptions{
		Kernel: osbuild.SysconfigKernelOptions{
			UpdateDefault: true,
			DefaultKernel: "kernel",
		},
		Network: osbuild.SysconfigNetworkOptions{
			Networking: true,
			NoZeroConf: true,
		},
	}))

	p.AddStage(osbuild.NewRPMOSTreePrepTreeStage(&osbuild.RPMOSTreePrepTreeStageOptions{
		EtcGroupMembers: []string{
			// NOTE: We may want to make this configurable.
			"wheel", "docker",
		},
	}))
	return p, nil
}

func (t *ImageTypeS2) ostreeCommitPipeline(options distro.ImageOptions) *osbuild.Pipeline {
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

func (t *ImageTypeS2) containerTreePipeline(repos []rpmmd.RepoConfig, packages []rpmmd.PackageSpec, options distro.ImageOptions, c *blueprint.Customizations) *osbuild.Pipeline {
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

func (t *ImageTypeS2) containerPipeline() *osbuild.Pipeline {
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

func (t *ImageTypeS2) anacondaTreePipeline(repos []rpmmd.RepoConfig, packages []rpmmd.PackageSpec, options distro.ImageOptions, c *blueprint.Customizations) *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	p.Name = "anaconda-tree"
	p.Build = "name:build"
	ostreePath := "/ostree/repo"
	p.AddStage(osbuild.NewRPMStage(t.rpmStageOptions(repos), t.rpmStageInputs(packages)))
	p.AddStage(osbuild.NewOSTreeInitStage(&osbuild.OSTreeInitStageOptions{Path: ostreePath}))
	p.AddStage(osbuild.NewOSTreePullStage(
		&osbuild.OSTreePullStageOptions{Repo: ostreePath},
		t.ostreePullStageInputs(options),
	))
	p.AddStage(osbuild.NewBuildstampStage(t.buildStampStageOptions()))
	language, _ := c.GetPrimaryLocale()
	if language != nil {
		p.AddStage(osbuild.NewLocaleStage(&osbuild.LocaleStageOptions{Language: *language}))
	} else {
		p.AddStage(osbuild.NewLocaleStage(&osbuild.LocaleStageOptions{Language: "en_US.UTF-8"}))
	}

	rootPassword := ""

	rootUser := osbuild.UsersStageOptionsUser{
		Password: &rootPassword,
	}

	installUID := 0
	installGID := 0
	installHome := "/root"
	installShell := "/usr/libexec/anaconda/run-anaconda"
	installPassword := ""
	installUser := osbuild.UsersStageOptionsUser{
		UID:      &installUID,
		GID:      &installGID,
		Home:     &installHome,
		Shell:    &installShell,
		Password: &installPassword,
	}
	usersStageOptions := &osbuild.UsersStageOptions{
		Users: map[string]osbuild.UsersStageOptionsUser{
			"root":    rootUser,
			"install": installUser,
		},
	}
	p.AddStage(osbuild.NewUsersStage(usersStageOptions))

	p.AddStage(osbuild.NewAnacondaStage(t.anacondaStageOptions()))

	p.AddStage(osbuild.NewLoraxScriptStage(t.loraxScriptStageOptions()))

	p.AddStage(osbuild.NewDracutStage(t.dracutStageOptions()))

	p.AddStage(osbuild.NewKickstartStage(t.kickstartStageOptions(ostreePath)))

	return p
}

func (t *ImageTypeS2) bootISOTreePipeline() *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	p.Name = "bootiso-tree"
	p.Build = "name:build"

	p.AddStage(osbuild.NewBootISOMonoStage(t.bootISOMonoStageOptions(), t.bootISOMonoStageInputs()))
	p.AddStage(osbuild.NewDiscinfoStage(t.discinfoStageOptions()))

	return p
}
func (t *ImageTypeS2) bootISOPipeline() *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	// NOTE(akoutsou) 1to2t: final pipeline should always be named "assembler"
	p.Name = "assembler"
	p.Build = "name:build"

	p.AddStage(osbuild.NewXorrisofsStage(t.xorrisofsStageOptions(), t.xorrisofsStageInputs()))
	p.AddStage(osbuild.NewImplantisomd5Stage(&osbuild.Implantisomd5StageOptions{Filename: t.Filename()}))

	return p
}

func (t *ImageTypeS2) rpmStageInputs(specs []rpmmd.PackageSpec) *osbuild.RPMStageInputs {
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

func (t *ImageTypeS2) ostreePullStageInputs(options distro.ImageOptions) *osbuild.OSTreePullStageInputs {
	pullStageInput := new(osbuild.OSTreePullStageInput)
	pullStageInput.Type = "org.osbuild.ostree"
	pullStageInput.Origin = "org.osbuild.pipeline"

	inputRefs := make(map[string]osbuild.OSTreePullStageReference)
	inputRefs["name:ostree-commit"] = osbuild.OSTreePullStageReference{Ref: t.OSTreeRef()}
	pullStageInput.References = inputRefs
	return &osbuild.OSTreePullStageInputs{Commits: pullStageInput}
}

func (t *ImageTypeS2) rpmStageOptions(repos []rpmmd.RepoConfig) *osbuild.RPMStageOptions {
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

func (t *ImageTypeS2) selinuxStageOptions() *osbuild.SELinuxStageOptions {
	return &osbuild.SELinuxStageOptions{
		FileContexts: "etc/selinux/targeted/contexts/files/file_contexts",
	}
}

func (t *ImageTypeS2) userStageOptions(users []blueprint.UserCustomization) (*osbuild.UsersStageOptions, error) {
	options := osbuild.UsersStageOptions{
		Users: make(map[string]osbuild.UsersStageOptionsUser),
	}

	for _, c := range users {
		if c.Password != nil && !crypt.PasswordIsCrypted(*c.Password) {
			cryptedPassword, err := crypt.CryptSHA512(*c.Password)
			if err != nil {
				return nil, err
			}

			c.Password = &cryptedPassword
		}

		user := osbuild.UsersStageOptionsUser{
			Groups:      c.Groups,
			Description: c.Description,
			Home:        c.Home,
			Shell:       c.Shell,
			Password:    c.Password,
			Key:         c.Key,
		}

		user.UID = c.UID
		user.GID = c.GID

		options.Users[c.Name] = user
	}

	return &options, nil
}

func (t *ImageTypeS2) groupStageOptions(groups []blueprint.GroupCustomization) *osbuild.GroupsStageOptions {
	options := osbuild.GroupsStageOptions{
		Groups: map[string]osbuild.GroupsStageOptionsGroup{},
	}

	for _, group := range groups {
		groupData := osbuild.GroupsStageOptionsGroup{
			Name: group.Name,
		}
		groupData.GID = group.GID

		options.Groups[group.Name] = groupData
	}

	return &options
}

func (t *ImageTypeS2) firewallStageOptions(firewall *blueprint.FirewallCustomization) *osbuild.FirewallStageOptions {
	options := osbuild.FirewallStageOptions{
		Ports: firewall.Ports,
	}

	if firewall.Services != nil {
		options.EnabledServices = firewall.Services.Enabled
		options.DisabledServices = firewall.Services.Disabled
	}

	return &options
}

func (t *ImageTypeS2) systemdStageOptions(enabledServices, disabledServices []string, s *blueprint.ServicesCustomization, target string) *osbuild.SystemdStageOptions {
	if s != nil {
		enabledServices = append(enabledServices, s.Enabled...)
		disabledServices = append(disabledServices, s.Disabled...)
	}
	return &osbuild.SystemdStageOptions{
		EnabledServices:  enabledServices,
		DisabledServices: disabledServices,
		DefaultTarget:    target,
	}
}

func (t *ImageTypeS2) buildStampStageOptions() *osbuild.BuildstampStageOptions {
	return &osbuild.BuildstampStageOptions{
		Arch:    t.Arch().Name(),
		Product: "Red Hat Enterprise Linux",
		Version: "8.4",
		Variant: "edge",
		Final:   true,
	}
}

func (t *ImageTypeS2) anacondaStageOptions() *osbuild.AnacondaStageOptions {
	return &osbuild.AnacondaStageOptions{
		KickstartModules: []string{
			"org.fedoraproject.Anaconda.Modules.Network",
			"org.fedoraproject.Anaconda.Modules.Payloads",
			"org.fedoraproject.Anaconda.Modules.Storage",
		},
	}
}

func (t *ImageTypeS2) loraxScriptStageOptions() *osbuild.LoraxScriptStageOptions {
	return &osbuild.LoraxScriptStageOptions{
		Path:     "99-generic/runtime-postinstall.tmpl",
		BaseArch: t.Arch().Name(),
	}
}

func (t *ImageTypeS2) dracutStageOptions() *osbuild.DracutStageOptions {
	kernel := []string{"4.18.0-293.el8.x86_64"}
	modules := []string{
		"bash",
		"systemd",
		"fips",
		"systemd-initrd",
		"modsign",
		"nss-softokn",
		"rdma",
		"rngd",
		"i18n",
		"convertfs",
		"network-manager",
		"network",
		"ifcfg",
		"url-lib",
		"drm",
		"plymouth",
		"prefixdevname",
		"prefixdevname-tools",
		"anaconda",
		"crypt",
		"dm",
		"dmsquash-live",
		"kernel-modules",
		"kernel-modules-extra",
		"kernel-network-modules",
		"livenet",
		"lvm",
		"mdraid",
		"multipath",
		"qemu",
		"qemu-net",
		"fcoe",
		"fcoe-uefi",
		"iscsi",
		"lunmask",
		"nfs",
		"resume",
		"rootfs-block",
		"terminfo",
		"udev-rules",
		"biosdevname",
		"dracut-systemd",
		"pollcdrom",
		"usrmount",
		"base",
		"fs-lib",
		"img-lib",
		"shutdown",
		"uefi-lib",
	}
	return &osbuild.DracutStageOptions{
		Kernel:  kernel,
		Modules: modules,
		Install: []string{"/.buildstamp"},
	}
}

func (t *ImageTypeS2) kickstartStageOptions(ostreePath string) *osbuild.KickstartStageOptions {
	return &osbuild.KickstartStageOptions{
		Path: "/usr/share/anaconda/interactive-defaults.ks",
		OSTree: osbuild.OSTreeOptions{
			OSName: "rhel",
			URL:    fmt.Sprintf("file://%s", ostreePath),
			Ref:    t.OSTreeRef(),
			GPG:    false,
		},
	}
}

func (t *ImageTypeS2) bootISOMonoStageOptions() *osbuild.BootISOMonoStageOptions {
	return &osbuild.BootISOMonoStageOptions{
		Product: osbuild.Product{
			Name:    "Red Hat Enterprise Linux",
			Version: "8.4",
		},
		ISOLabel: "RHEL-8-4-X86_64",
		// TODO: based on image arch
		Kernel: "4.18.0-293.el8.x86_64",
		EFI: osbuild.EFI{
			Architectures: []string{
				// TODO: based on image arch
				"IA32",
				"X64",
			},
			Vendor: "redhat",
		},
		RootFS: osbuild.RootFS{
			Size: 4096,
			Compression: osbuild.FSCompression{
				Method: "xz",
				Options: osbuild.FSCompressionOptions{
					// TODO: based on image arch
					BCJ: "x86",
				},
			},
		},
	}
}

func (t *ImageTypeS2) bootISOMonoStageInputs() *osbuild.BootISOMonoStageInputs {
	rootfsInput := new(osbuild.BootISOMonoStageInput)
	rootfsInput.Type = "org.osbuild.tree"
	rootfsInput.Origin = "org.osbuild.pipeline"
	rootfsInput.References = osbuild.BootISOStageReferences{"name:anaconda-tree"}
	return &osbuild.BootISOMonoStageInputs{
		RootFS: rootfsInput,
	}
}

func (t *ImageTypeS2) discinfoStageOptions() *osbuild.DiscinfoStageOptions {
	return &osbuild.DiscinfoStageOptions{
		BaseArch: t.Arch().Name(),
		Release:  "202010217.n.0",
	}
}

func (t *ImageTypeS2) xorrisofsStageOptions() *osbuild.XorrisofsStageOptions {
	return &osbuild.XorrisofsStageOptions{
		Filename: t.Filename(),
		VolID:    "RHEL-8-4-X86_64", // TODO: add to image type fields
		Boot: osbuild.XorrisofsBoot{
			Image:   "isolinux/isolinux.bin",
			Catalog: "isolinux/boot.cat",
		},
		EFI:          "images/efiboot.img",
		IsohybridMBR: "/usr/share/syslinux/isohdpfx.bin",
	}
}

func (t *ImageTypeS2) xorrisofsStageInputs() *osbuild.XorrisofsStageInputs {
	input := new(osbuild.XorrisofsStageInput)
	input.Type = "org.osbuild.tree"
	input.Origin = "org.osbuild.pipeline"
	input.References = osbuild.XorrisofsStageReferences{"name:bootiso-tree"}
	return &osbuild.XorrisofsStageInputs{Tree: input}
}

type solver func(specs []string, excludeSpecs []string) ([]rpmmd.PackageSpec, map[string]string, error)

func (t *ImageTypeS2) SetSolver(s solver, bp *blueprint.Blueprint) {
	t.depsolve = s
	t.blueprint = bp
}
