package rhel85

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"path/filepath"
	"sort"

	"github.com/osbuild/osbuild-composer/internal/blueprint"
	"github.com/osbuild/osbuild-composer/internal/crypt"
	"github.com/osbuild/osbuild-composer/internal/distro"
	osbuild "github.com/osbuild/osbuild-composer/internal/osbuild2"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

const name = "rhel-85"
const osVersion = "8.5"
const modulePlatformID = "platform:el8"
const ostreeRef = "rhel/8/%s/edge"

type distribution struct {
	arches map[string]distro.Arch
}

func (d *distribution) Name() string {
	return name
}

func (d *distribution) ModulePlatformID() string {
	return modulePlatformID
}

func (d *distribution) ListArches() []string {
	archNames := make([]string, 0, len(d.arches))
	for name := range d.arches {
		archNames = append(archNames, name)
	}
	sort.Strings(archNames)
	return archNames
}

func (d *distribution) GetArch(name string) (distro.Arch, error) {
	arch, exists := d.arches[name]
	if !exists {
		return nil, errors.New("invalid architecture: " + name)
	}
	return arch, nil
}

func (d *distribution) addArches(arches ...architecture) {
	if d.arches == nil {
		d.arches = map[string]distro.Arch{}
	}

	for _, a := range arches {
		d.arches[a.name] = &architecture{
			distro:     d,
			name:       a.name,
			imageTypes: a.imageTypes,
		}
	}
}

type architecture struct {
	distro     *distribution
	name       string
	imageTypes map[string]distro.ImageType
}

func (a *architecture) Name() string {
	return a.name
}

func (a *architecture) ListImageTypes() []string {
	itNames := make([]string, 0, len(a.imageTypes))
	for name := range a.imageTypes {
		itNames = append(itNames, name)
	}
	sort.Strings(itNames)
	return itNames
}

func (a *architecture) GetImageType(name string) (distro.ImageType, error) {
	t, exists := a.imageTypes[name]
	if !exists {
		return nil, errors.New("invalid image type: " + name)
	}
	return t, nil
}

func (a *architecture) addImageTypes(imageTypes ...imageType) {
	if a.imageTypes == nil {
		a.imageTypes = map[string]distro.ImageType{}
	}
	for _, it := range imageTypes {
		a.imageTypes[it.name] = &imageType{
			arch:             a,
			name:             it.name,
			filename:         it.filename,
			mimeType:         it.mimeType,
			packageSets:      it.packageSets,
			enabledServices:  it.enabledServices,
			disabledServices: it.disabledServices,
			defaultTarget:    it.defaultTarget,
			kernelOptions:    it.kernelOptions,
			bootable:         it.bootable,
			bootISO:          it.bootISO,
			rpmOstree:        it.rpmOstree,
			defaultSize:      it.defaultSize,
			exports:          it.exports,
		}
	}
}

func (a *architecture) Distro() distro.Distro {
	return a.distro
}

type imageType struct {
	arch             *architecture
	name             string
	filename         string
	mimeType         string
	packageSets      map[string]rpmmd.PackageSet
	enabledServices  []string
	disabledServices []string
	defaultTarget    string
	bootISO          bool
	rpmOstree        bool
	defaultSize      uint64
	exports          []string
}

func (t *imageType) Name() string {
	return t.name
}

func (t *imageType) Arch() distro.Arch {
	return t.arch
}

func (t *imageType) Filename() string {
	return t.filename
}

func (t *imageType) MIMEType() string {
	return t.mimeType
}

func (t *imageType) OSTreeRef() string {
	if t.rpmOstree {
		return fmt.Sprintf(ostreeRef, t.arch.name)
	}
	return ""
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

func (t *imageType) PackageSets(bp blueprint.Blueprint) map[string]rpmmd.PackageSet {
	sets := make(map[string]rpmmd.PackageSet)
	for name, pkgSet := range t.packageSets {
		if name == "packages" {
			// combine image packages with blueprint
			pkgSet.Include = append(pkgSet.Include, bp.GetPackages()...)
			timezone, _ := bp.Customizations.GetTimezoneSettings()
			if timezone != nil {
				pkgSet.Include = append(pkgSet.Include, "chrony")
			}
		}
		sets[name] = pkgSet
	}
	return sets

}

func (t *imageType) Exports() []string {
	if len(t.exports) > 0 {
		return t.exports
	}
	return []string{"assembler"}
}

// local type for ostree commit metadata used to define commit sources
type ostreeCommit struct {
	Checksum string
	URL      string
}

func (t *imageType) Manifest(customizations *blueprint.Customizations,
	options distro.ImageOptions,
	repos []rpmmd.RepoConfig,
	packageSpecSets map[string][]rpmmd.PackageSpec,
	seed int64) (distro.Manifest, error) {

	source := rand.NewSource(seed)
	rng := rand.New(source)
	pipelines, err := t.pipelines(customizations, options, repos, packageSpecSets, rng)
	if err != nil {
		return distro.Manifest{}, err
	}

	// flatten spec sets for sources
	allPackageSpecs := make([]rpmmd.PackageSpec, 0)
	for _, specs := range packageSpecSets {
		allPackageSpecs = append(allPackageSpecs, specs...)
	}

	var commits []ostreeCommit
	if t.bootISO && options.OSTree.Parent != "" && options.OSTree.URL != "" {
		commits = []ostreeCommit{{Checksum: options.OSTree.Parent, URL: options.OSTree.URL}}
	}
	return json.Marshal(
		osbuild.Manifest{
			Version:   "2",
			Pipelines: pipelines,
			Sources:   t.sources(allPackageSpecs, commits),
		},
	)
}

// checkOptions checks the validity and compatibility of options and customizations for the image type.
func (t *imageType) checkOptions(customizations *blueprint.Customizations, options distro.ImageOptions) error {
	if t.bootISO {
		if options.OSTree.Parent == "" {
			return fmt.Errorf("boot ISO image type %q requires specifying a URL from which to retrieve the OSTree commit", t.name)
		}
		if customizations != nil {
			return fmt.Errorf("boot ISO image type %q does not support blueprint customizations", t.name)
		}
	}

	if kernelOpts := customizations.GetKernel(); kernelOpts.Append != "" && t.rpmOstree {
		return fmt.Errorf("kernel boot parameter customizations are not supported for ostree types")
	}

	return nil
}

func (t *imageType) installerPipelines(options distro.ImageOptions, repos []rpmmd.RepoConfig, installerPackages []rpmmd.PackageSpec) ([]osbuild.Pipeline, error) {
	kernelPkg := new(rpmmd.PackageSpec)
	for _, pkg := range installerPackages {
		if pkg.Name == "kernel" {
			kernelPkg = &pkg
			break
		}
	}
	if kernelPkg == nil {
		return nil, fmt.Errorf("kernel package not found in installer package set")
	}
	kernelVer := fmt.Sprintf("%s-%s.%s", kernelPkg.Version, kernelPkg.Release, kernelPkg.Arch)
	pipelines := make([]osbuild.Pipeline, 0)
	pipelines = append(pipelines, *t.anacondaTreePipeline(repos, installerPackages, options, kernelVer))
	pipelines = append(pipelines, *t.bootISOTreePipeline(kernelVer))
	pipelines = append(pipelines, *t.bootISOPipeline())
	return pipelines, nil

}

func (t *imageType) pipelines(customizations *blueprint.Customizations, options distro.ImageOptions, repos []rpmmd.RepoConfig, packageSetSpecs map[string][]rpmmd.PackageSpec, rng *rand.Rand) ([]osbuild.Pipeline, error) {

	if err := t.checkOptions(customizations, options); err != nil {
		return nil, err
	}

	pipelines := make([]osbuild.Pipeline, 0)
	pipelines = append(pipelines, *t.buildPipeline(repos, packageSetSpecs["build"]))
	if t.bootISO {
		ips, err := t.installerPipelines(options, repos, packageSetSpecs["installer"])
		if err != nil {
			return nil, err
		}
		return append(pipelines, ips...), nil
	}

	treePipeline, err := t.ostreeTreePipeline(repos, packageSetSpecs["packages"], customizations)
	if err != nil {
		return nil, err
	}

	pipelines = append(pipelines, *treePipeline)
	pipelines = append(pipelines, *t.ostreeCommitPipeline(options))

	if t.name == "edge-container" {
		pipelines = append(pipelines, *t.containerTreePipeline(repos, packageSetSpecs["container"], options, customizations))
		pipelines = append(pipelines, *t.containerPipeline())
	} else if t.name == "edge-commit" {
		pipelines = append(pipelines, *t.commitTarPipeline())
	}

	return pipelines, nil
}

func (t *imageType) sources(packages []rpmmd.PackageSpec, ostreeCommits []ostreeCommit) osbuild.Sources {
	sources := osbuild.Sources{}
	curl := &osbuild.CurlSource{
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
		curl.Items[pkg.Checksum] = item
	}
	if len(curl.Items) > 0 {
		sources["org.osbuild.curl"] = curl
	}

	ostree := &osbuild.OSTreeSource{
		Items: make(map[string]osbuild.OSTreeSourceItem),
	}
	for _, commit := range ostreeCommits {
		item := new(osbuild.OSTreeSourceItem)
		item.Remote.URL = commit.URL
		ostree.Items[commit.Checksum] = *item
	}
	if len(ostree.Items) > 0 {
		sources["org.osbuild.ostree"] = ostree
	}
	return sources
}

func (t *imageType) buildPipeline(repos []rpmmd.RepoConfig, buildPackageSpecs []rpmmd.PackageSpec) *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	p.Name = "build"
	p.Runner = "org.osbuild.rhel85"
	p.AddStage(osbuild.NewRPMStage(t.rpmStageOptions(repos), t.rpmStageInputs(buildPackageSpecs)))
	p.AddStage(osbuild.NewSELinuxStage(t.selinuxStageOptions()))
	return p
}

func (t *imageType) ostreeTreePipeline(repos []rpmmd.RepoConfig, packages []rpmmd.PackageSpec, c *blueprint.Customizations) (*osbuild.Pipeline, error) {
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
		p.AddStage(osbuild.NewFirstBootStage(t.usersFirstBootOptions(options)))
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

	p.AddStage(osbuild.NewOSTreePrepTreeStage(&osbuild.OSTreePrepTreeStageOptions{
		EtcGroupMembers: []string{
			// NOTE: We may want to make this configurable.
			"wheel", "docker",
		},
	}))
	return p, nil
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
			OSVersion: osVersion,
			Parent:    options.OSTree.Parent,
		},
		&osbuild.OSTreeCommitStageInputs{Tree: commitStageInput}),
	)
	return p
}

func (t *imageType) commitTarPipeline() *osbuild.Pipeline {
	options := osbuild.TarStageOptions{Filename: t.Filename()}
	commitTree := new(osbuild.TarStageInput)
	commitTree.Type = "org.osbuild.tree"
	commitTree.Origin = "org.osbuild.pipeline"
	commitTree.References = []string{"name:ostree-commit"}
	tarStage := osbuild.NewTarStage(&options, &osbuild.TarStageInputs{Tree: commitTree})
	p := new(osbuild.Pipeline)
	p.Name = "commit-archive"
	p.Build = "name:build"
	p.AddStage(tarStage)
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
		t.ostreePullStageInputs("org.osbuild.pipeline", "name:ostree-commit", options.OSTree.Ref),
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

func (t *imageType) anacondaTreePipeline(repos []rpmmd.RepoConfig, packages []rpmmd.PackageSpec, options distro.ImageOptions, kernelVer string) *osbuild.Pipeline {
	ostreeRepoPath := "/ostree/repo"
	p := new(osbuild.Pipeline)
	p.Name = "anaconda-tree"
	p.Build = "name:build"
	p.AddStage(osbuild.NewRPMStage(t.rpmStageOptions(repos), t.rpmStageInputs(packages)))
	p.AddStage(osbuild.NewOSTreeInitStage(&osbuild.OSTreeInitStageOptions{Path: ostreeRepoPath}))
	p.AddStage(osbuild.NewOSTreePullStage(
		&osbuild.OSTreePullStageOptions{Repo: ostreeRepoPath},
		t.ostreePullStageInputs("org.osbuild.source", options.OSTree.Parent, options.OSTree.Ref),
	))
	p.AddStage(osbuild.NewBuildstampStage(t.buildStampStageOptions()))
	p.AddStage(osbuild.NewLocaleStage(&osbuild.LocaleStageOptions{Language: "en_US.UTF-8"}))

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
	p.AddStage(osbuild.NewDracutStage(t.dracutStageOptions(kernelVer)))
	p.AddStage(osbuild.NewKickstartStage(t.kickstartStageOptions(fmt.Sprintf("file://%s", ostreeRepoPath), options.OSTree.Ref)))

	return p
}

func (t *imageType) bootISOTreePipeline(kernelVer string) *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	p.Name = "bootiso-tree"
	p.Build = "name:build"

	p.AddStage(osbuild.NewBootISOMonoStage(t.bootISOMonoStageOptions(kernelVer), t.bootISOMonoStageInputs()))
	p.AddStage(osbuild.NewDiscinfoStage(t.discinfoStageOptions()))

	return p
}
func (t *imageType) bootISOPipeline() *osbuild.Pipeline {
	p := new(osbuild.Pipeline)
	p.Name = "bootiso"
	p.Build = "name:build"

	p.AddStage(osbuild.NewXorrisofsStage(t.xorrisofsStageOptions(), t.xorrisofsStageInputs()))
	p.AddStage(osbuild.NewImplantisomd5Stage(&osbuild.Implantisomd5StageOptions{Filename: t.Filename()}))

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

func (t *imageType) ostreePullStageInputs(origin, source, commitRef string) *osbuild.OSTreePullStageInputs {
	pullStageInput := new(osbuild.OSTreePullStageInput)
	pullStageInput.Type = "org.osbuild.ostree"
	pullStageInput.Origin = origin

	inputRefs := make(map[string]osbuild.OSTreePullStageReference)
	inputRefs[source] = osbuild.OSTreePullStageReference{Ref: commitRef}
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

	options := &osbuild.SELinuxStageOptions{
		FileContexts: "etc/selinux/targeted/contexts/files/file_contexts",
	}
	if t.bootISO {
		options.Labels = map[string]string{
			"/usr/bin/cp":  "system_u:object_r:install_exec_t:s0",
			"/usr/bin/tar": "system_u:object_r:install_exec_t:s0",
		}
	}
	return options
}

func (t *imageType) userStageOptions(users []blueprint.UserCustomization) (*osbuild.UsersStageOptions, error) {
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

func (t *imageType) usersFirstBootOptions(usersStageOptions *osbuild.UsersStageOptions) *osbuild.FirstBootStageOptions {
	cmds := make([]string, 0, 3*len(usersStageOptions.Users)+1)
	// workaround for creating authorized_keys file for user
	varhome := filepath.Join("/var", "home")
	for name, user := range usersStageOptions.Users {
		if user.Key != nil {
			sshdir := filepath.Join(varhome, name, ".ssh")
			cmds = append(cmds, fmt.Sprintf("mkdir -p %s", sshdir))
			cmds = append(cmds, fmt.Sprintf("sh -c 'echo %q >> %q'", *user.Key, filepath.Join(sshdir, "authorized_keys")))
			cmds = append(cmds, fmt.Sprintf("chown %s:%s -Rc %s", name, name, sshdir))
		}
	}
	cmds = append(cmds, fmt.Sprintf("restorecon -rvF %s", varhome))
	options := &osbuild.FirstBootStageOptions{
		Commands:       cmds,
		WaitForNetwork: false,
	}

	return options
}

func (t *imageType) groupStageOptions(groups []blueprint.GroupCustomization) *osbuild.GroupsStageOptions {
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

func (t *imageType) firewallStageOptions(firewall *blueprint.FirewallCustomization) *osbuild.FirewallStageOptions {
	options := osbuild.FirewallStageOptions{
		Ports: firewall.Ports,
	}

	if firewall.Services != nil {
		options.EnabledServices = firewall.Services.Enabled
		options.DisabledServices = firewall.Services.Disabled
	}

	return &options
}

func (t *imageType) systemdStageOptions(enabledServices, disabledServices []string, s *blueprint.ServicesCustomization, target string) *osbuild.SystemdStageOptions {
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

func (t *imageType) buildStampStageOptions() *osbuild.BuildstampStageOptions {
	return &osbuild.BuildstampStageOptions{
		Arch:    t.Arch().Name(),
		Product: "Red Hat Enterprise Linux",
		Version: osVersion,
		Variant: "edge",
		Final:   true,
	}
}

func (t *imageType) anacondaStageOptions() *osbuild.AnacondaStageOptions {
	return &osbuild.AnacondaStageOptions{
		KickstartModules: []string{
			"org.fedoraproject.Anaconda.Modules.Network",
			"org.fedoraproject.Anaconda.Modules.Payloads",
			"org.fedoraproject.Anaconda.Modules.Storage",
		},
	}
}

func (t *imageType) loraxScriptStageOptions() *osbuild.LoraxScriptStageOptions {
	return &osbuild.LoraxScriptStageOptions{
		Path:     "99-generic/runtime-postinstall.tmpl",
		BaseArch: t.Arch().Name(),
	}
}

func (t *imageType) dracutStageOptions(kernelVer string) *osbuild.DracutStageOptions {
	kernel := []string{kernelVer}
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

func (t *imageType) kickstartStageOptions(ostreeURL, ostreeRef string) *osbuild.KickstartStageOptions {
	return &osbuild.KickstartStageOptions{
		Path: "/usr/share/anaconda/interactive-defaults.ks",
		OSTree: osbuild.OSTreeOptions{
			OSName: "rhel",
			URL:    ostreeURL,
			Ref:    ostreeRef,
			GPG:    false,
		},
	}
}

func (t *imageType) bootISOMonoStageOptions(kernelVer string) *osbuild.BootISOMonoStageOptions {
	comprOptions := new(osbuild.FSCompressionOptions)
	if bcj := osbuild.BCJOption(t.arch.Name()); bcj != "" {
		comprOptions.BCJ = bcj
	}
	return &osbuild.BootISOMonoStageOptions{
		Product: osbuild.Product{
			Name:    "Red Hat Enterprise Linux",
			Version: osVersion,
		},
		ISOLabel: fmt.Sprintf("RHEL-8-5-0-BaseOS-%s", t.Arch().Name()),
		Kernel:   kernelVer,
		EFI: osbuild.EFI{
			Architectures: []string{
				"IA32",
				"X64",
			},
			Vendor: "redhat",
		},
		ISOLinux: osbuild.ISOLinux{
			Enabled: true,
			Debug:   false,
		},
		Templates: "80-rhel",
		RootFS: osbuild.RootFS{
			Size: 4096,
			Compression: osbuild.FSCompression{
				Method:  "xz",
				Options: comprOptions,
			},
		},
	}
}

func (t *imageType) bootISOMonoStageInputs() *osbuild.BootISOMonoStageInputs {
	rootfsInput := new(osbuild.BootISOMonoStageInput)
	rootfsInput.Type = "org.osbuild.tree"
	rootfsInput.Origin = "org.osbuild.pipeline"
	rootfsInput.References = osbuild.BootISOMonoStageReferences{"name:anaconda-tree"}
	return &osbuild.BootISOMonoStageInputs{
		RootFS: rootfsInput,
	}
}

func (t *imageType) discinfoStageOptions() *osbuild.DiscinfoStageOptions {
	return &osbuild.DiscinfoStageOptions{
		BaseArch: t.Arch().Name(),
		Release:  "202010217.n.0",
	}
}

func (t *imageType) xorrisofsStageOptions() *osbuild.XorrisofsStageOptions {
	return &osbuild.XorrisofsStageOptions{
		Filename: t.Filename(),
		VolID:    fmt.Sprintf("RHEL-8-5-0-BaseOS-%s", t.Arch().Name()),
		SysID:    "LINUX",
		Boot: osbuild.XorrisofsBoot{
			Image:   "isolinux/isolinux.bin",
			Catalog: "isolinux/boot.cat",
		},
		EFI:          "images/efiboot.img",
		IsohybridMBR: "/usr/share/syslinux/isohdpfx.bin",
	}
}

func (t *imageType) xorrisofsStageInputs() *osbuild.XorrisofsStageInputs {
	input := new(osbuild.XorrisofsStageInput)
	input.Type = "org.osbuild.tree"
	input.Origin = "org.osbuild.pipeline"
	input.References = osbuild.XorrisofsStageReferences{"name:bootiso-tree"}
	return &osbuild.XorrisofsStageInputs{Tree: input}
}

// New creates a new distro object, defining the supported architectures and image types
func New() distro.Distro {
	rd := new(distribution)

	// Shared Package sets
	edgeCommitCommonPkgSet := rpmmd.PackageSet{
		Include: []string{
			"redhat-release",
			"glibc", "glibc-minimal-langpack", "nss-altfiles",
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
		},
		Exclude: []string{"rng-tools"},
	}
	edgeBuildPkgSet := rpmmd.PackageSet{
		Include: []string{
			"dnf", "dosfstools", "e2fsprogs", "efibootmgr", "genisoimage",
			"grub2-efi-ia32-cdboot", "grub2-efi-x64", "grub2-efi-x64-cdboot",
			"grub2-pc", "grub2-pc-modules", "grub2-tools", "grub2-tools-efi",
			"grub2-tools-extra", "grub2-tools-minimal", "isomd5sum",
			"lorax-templates-generic", "lorax-templates-rhel",
			"policycoreutils", "python36", "python3-iniparse", "qemu-img",
			"rpm-ostree", "selinux-policy-targeted", "shim-ia32", "shim-x64",
			"squashfs-tools", "syslinux", "syslinux-nonlinux", "systemd",
			"tar", "xfsprogs", "xorriso", "xz",
		},
		Exclude: nil,
	}
	edgeInstallerPkgSet := rpmmd.PackageSet{
		Include: []string{
			"aajohan-comfortaa-fonts", "abattis-cantarell-fonts",
			"alsa-firmware", "alsa-tools-firmware", "anaconda",
			"anaconda-dracut", "anaconda-install-env-deps", "anaconda-widgets",
			"audit", "bind-utils", "biosdevname", "bitmap-fangsongti-fonts",
			"bzip2", "cryptsetup", "curl", "dbus-x11", "dejavu-sans-fonts",
			"dejavu-sans-mono-fonts", "device-mapper-persistent-data",
			"dmidecode", "dnf", "dracut-config-generic", "dracut-network",
			"dump", "efibootmgr", "ethtool", "ftp", "gdb-gdbserver", "gdisk",
			"gfs2-utils", "glibc-all-langpacks",
			"google-noto-sans-cjk-ttc-fonts", "grub2-efi-ia32-cdboot",
			"grub2-efi-x64-cdboot", "grub2-tools", "grub2-tools-efi",
			"grub2-tools-extra", "grub2-tools-minimal", "grubby",
			"gsettings-desktop-schemas", "hdparm", "hexedit", "hostname",
			"initscripts", "ipmitool", "iwl1000-firmware", "iwl100-firmware",
			"iwl105-firmware", "iwl135-firmware", "iwl2000-firmware",
			"iwl2030-firmware", "iwl3160-firmware", "iwl3945-firmware",
			"iwl4965-firmware", "iwl5000-firmware", "iwl5150-firmware",
			"iwl6000-firmware", "iwl6000g2a-firmware", "iwl6000g2b-firmware",
			"iwl6050-firmware", "iwl7260-firmware", "jomolhari-fonts",
			"kacst-farsi-fonts", "kacst-qurn-fonts", "kbd", "kbd-misc",
			"kdump-anaconda-addon", "kernel", "khmeros-base-fonts", "less",
			"libblockdev-lvm-dbus", "libertas-sd8686-firmware",
			"libertas-sd8787-firmware", "libertas-usb8388-firmware",
			"libertas-usb8388-olpc-firmware", "libibverbs",
			"libreport-plugin-bugzilla", "libreport-plugin-reportuploader",
			"libreport-rhel-anaconda-bugzilla", "librsvg2", "linux-firmware",
			"lklug-fonts", "lohit-assamese-fonts", "lohit-bengali-fonts",
			"lohit-devanagari-fonts", "lohit-gujarati-fonts",
			"lohit-gurmukhi-fonts", "lohit-kannada-fonts", "lohit-odia-fonts",
			"lohit-tamil-fonts", "lohit-telugu-fonts", "lsof", "madan-fonts",
			"memtest86+", "metacity", "mtr", "mt-st", "net-tools", "nfs-utils",
			"nmap-ncat", "nm-connection-editor", "nss-tools",
			"openssh-clients", "openssh-server", "oscap-anaconda-addon",
			"ostree", "pciutils", "perl-interpreter", "pigz", "plymouth",
			"prefixdevname", "python3-pyatspi", "rdma-core",
			"redhat-release-eula", "rng-tools", "rpcbind", "rpm-ostree",
			"rsync", "rsyslog", "selinux-policy-targeted", "sg3_utils",
			"shim-ia32", "shim-x64", "sil-abyssinica-fonts",
			"sil-padauk-fonts", "sil-scheherazade-fonts", "smartmontools",
			"smc-meera-fonts", "spice-vdagent", "strace", "syslinux",
			"systemd", "system-storage-manager", "tar",
			"thai-scalable-waree-fonts", "tigervnc-server-minimal",
			"tigervnc-server-module", "udisks2", "udisks2-iscsi", "usbutils",
			"vim-minimal", "volume_key", "wget", "xfsdump", "xfsprogs",
			"xorg-x11-drivers", "xorg-x11-fonts-misc", "xorg-x11-server-utils",
			"xorg-x11-server-Xorg", "xorg-x11-xauth", "xz",
		},
		Exclude: nil,
	}
	edgeCommitX86PkgSet := rpmmd.PackageSet{
		Include: append(edgeCommitCommonPkgSet.Include,
			// x86 specific
			"grub2", "grub2-efi-x64", "efibootmgr", "shim-x64",
			"microcode_ctl", "iwl1000-firmware", "iwl100-firmware",
			"iwl105-firmware", "iwl135-firmware", "iwl2000-firmware",
			"iwl2030-firmware", "iwl3160-firmware", "iwl5000-firmware",
			"iwl5150-firmware", "iwl6000-firmware", "iwl6050-firmware",
			"iwl7260-firmware"),
		Exclude: edgeCommitCommonPkgSet.Exclude,
	}
	edgeCommitAarch64PkgSet := rpmmd.PackageSet{
		Include: append(edgeCommitCommonPkgSet.Include,
			// aarch64 specific
			"grub2-efi-aa64", "efibootmgr", "shim-aa64",
			"iwl7260-firmware"),
		Exclude: edgeCommitCommonPkgSet.Exclude,
	}

	// Shared Services
	edgeServices := []string{
		"NetworkManager.service", "firewalld.service", "sshd.service",
	}

	// Image Definitions
	edgeCommitImgTypeX86_64 := imageType{
		name:     "edge-commit",
		filename: "commit.tar",
		mimeType: "application/x-tar",
		packageSets: map[string]rpmmd.PackageSet{
			"build":    edgeBuildPkgSet,
			"packages": edgeCommitX86PkgSet,
		},
		enabledServices: edgeServices,
		rpmOstree:       true,
		exports:         []string{"commit-archive"},
	}
	edgeOCIImgTypeX86_64 := imageType{
		name:     "edge-container",
		filename: "container.tar",
		mimeType: "application/x-tar",
		packageSets: map[string]rpmmd.PackageSet{
			"build":     edgeBuildPkgSet,
			"packages":  edgeCommitX86PkgSet,
			"container": {Include: []string{"httpd"}},
		},
		enabledServices: edgeServices,
		rpmOstree:       true,
		bootISO:         false,
		exports:         []string{"container"},
	}
	edgeInstallerImgTypeX86_64 := imageType{
		name:     "edge-installer",
		filename: "installer.iso",
		mimeType: "application/x-iso9660-image",
		packageSets: map[string]rpmmd.PackageSet{
			"build":     edgeBuildPkgSet,
			"packages":  edgeCommitX86PkgSet,
			"installer": edgeInstallerPkgSet,
		},
		enabledServices: edgeServices,
		rpmOstree:       true,
		bootISO:         true,
		exports:         []string{"bootiso"},
	}

	x86_64 := architecture{
		name:   "x86_64",
		distro: rd,
	}
	x86_64.addImageTypes(edgeCommitImgTypeX86_64, edgeInstallerImgTypeX86_64, edgeOCIImgTypeX86_64)

	edgeCommitImgTypeAarch64 := imageType{
		name:     "edge-commit",
		filename: "commit.tar",
		mimeType: "application/x-tar",
		packageSets: map[string]rpmmd.PackageSet{
			"build":    edgeBuildPkgSet,
			"packages": edgeCommitAarch64PkgSet,
		},
		enabledServices: edgeServices,
		rpmOstree:       true,
		exports:         []string{"commit-archive"},
	}
	edgeOCIImgTypeAarch64 := imageType{
		name:     "edge-container",
		filename: "container.tar",
		mimeType: "application/x-tar",
		packageSets: map[string]rpmmd.PackageSet{
			"build":     edgeBuildPkgSet,
			"packages":  edgeCommitAarch64PkgSet,
			"container": {Include: []string{"httpd"}},
		},
		enabledServices: edgeServices,
		rpmOstree:       true,
		exports:         []string{"container"},
	}
	edgeInstallerImgTypeAarch64 := imageType{
		name:     "edge-installer",
		filename: "installer.iso",
		mimeType: "application/x-iso9660-image",
		packageSets: map[string]rpmmd.PackageSet{
			"build":     edgeBuildPkgSet,
			"packages":  edgeCommitX86PkgSet,
			"installer": edgeInstallerPkgSet,
		},
		enabledServices: edgeServices,
		rpmOstree:       true,
		bootISO:         true,
		exports:         []string{"bootiso"},
	}
	aarch64 := architecture{
		name:   "aarch64",
		distro: rd,
	}
	aarch64.addImageTypes(edgeCommitImgTypeAarch64, edgeOCIImgTypeAarch64, edgeInstallerImgTypeAarch64)
	rd.addArches(x86_64, aarch64)
	return rd
}
