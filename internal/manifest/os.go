package manifest

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/osbuild/osbuild-composer/internal/common"
	"github.com/osbuild/osbuild-composer/internal/container"
	"github.com/osbuild/osbuild-composer/internal/disk"
	"github.com/osbuild/osbuild-composer/internal/distro"
	"github.com/osbuild/osbuild-composer/internal/environment"
	"github.com/osbuild/osbuild-composer/internal/fsnode"
	"github.com/osbuild/osbuild-composer/internal/osbuild"
	"github.com/osbuild/osbuild-composer/internal/ostree"
	"github.com/osbuild/osbuild-composer/internal/platform"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
	"github.com/osbuild/osbuild-composer/internal/users"
	"github.com/osbuild/osbuild-composer/internal/workload"
)

// OSCustomizations encapsulates all configuration applied to the base
// operating system independently of where and how it is integrated and what
// workload it is running.
// TODO: move out kernel/bootloader/cloud-init/... to other
//
//	abstractions, this should ideally only contain things that
//	can always be applied.
type OSCustomizations struct {

	// Packages to install in addition to the ones required by the
	// pipeline.
	// WORKLOAD or REMOVE: what would require adding arbitrary packages other than a workload (whether built-in or custom)?
	ExtraBasePackages []string

	// Packages to exclude from the base package set. This is useful in
	// case of weak dependencies, comps groups, or where multiple packages
	// can satisfy a dependency. Must not conflict with the included base
	// package set.
	// WORKLOAD or REMOVE: what would require excluding arbitrary packages other than a workload (whether built-in or custom)?
	ExcludeBasePackages []string

	// Additional repos to install the base packages from.
	ExtraBaseRepos []rpmmd.RepoConfig

	// Containers to embed in the image
	// WORKLOAD: currently, only user-defined (custom) workloads include containers.
	Containers []container.Spec

	// KernelName indicates that a kernel is installed, and names the kernel
	// package.
	// WORKLOAD: this is modified by the user, usually to install the rt kernel
	KernelName string

	// KernelOptionsAppend are appended to the kernel commandline
	KernelOptionsAppend []string

	// KernelOptionsBootloader controls whether kernel command line options
	// should be specified in the bootloader grubenv configuration. Otherwise
	// they are specified in /etc/kernel/cmdline (default).
	//
	// NB: The kernel options need to be still specified in /etc/default/grub
	// under the GRUB_CMDLINE_LINUX variable. The reason is that it is used by
	// the 10_linux script executed by grub2-mkconfig to override the kernel
	// options in /etc/kernel/cmdline if the file has older timestamp than
	// /etc/default/grub.
	//
	// This should only be used for RHEL 8 and CentOS 8 images that use grub
	// (non s390x).  Newer releases (9+) should keep this disabled.
	// PLATFORM
	KernelOptionsBootloader bool

	GPGKeyFiles []string

	// WORKLOAD: built-in or custom
	Language         string
	Keyboard         *string
	X11KeymapLayouts []string
	Hostname         string
	Timezone         string
	EnabledServices  []string
	DisabledServices []string
	DefaultTarget    string

	// SELinux policy, when set it enables the labeling of the tree with the
	// selected profile
	// ENVIRONMENT
	SElinux string

	// Only used on RHEL 7
	SELinuxForceRelabel *bool

	// Do not install documentation
	// WORKLOAD
	ExcludeDocs bool

	// WORKLOAD or ENVIRONMENT depending on purpose.  For example, the
	// "ec2-user" for AWS should be defined by the environment.
	Groups []users.Group
	Users  []users.User

	// WORKLOAD
	Firewall *osbuild.FirewallStageOptions

	// PLATFORM: bootloader configuration; hardware dependent
	Grub2Config *osbuild.GRUB2Config

	// PLATFORM: hardware (network) configs
	Sysconfig []*osbuild.SysconfigStageOptions

	// PLATFORM or ENVIRONMENT: used for AMIs to set NAutoVTs
	SystemdLogind []*osbuild.SystemdLogindStageOptions

	// ENVIRONMENT: depends on the deployment environment of the image
	CloudInit []*osbuild.CloudInitStageOptions

	// PLATFORM: hardware
	Modprobe []*osbuild.ModprobeStageOptions

	// PLATFORM: hardware
	DracutConf []*osbuild.DracutConfStageOptions

	// ENVIRONMENT
	// Only supports overriding the environment of the nm-cloud-setup.service
	// for automatically setting up network manager in cloud environments
	SystemdUnit []*osbuild.SystemdUnitStageOptions

	// ENVIRONMENT or WORKLOAD
	// Select identity and authentication sources
	Authselect    *osbuild.AuthselectStageOptions
	SELinuxConfig *osbuild.SELinuxConfigStageOptions

	// PLATFORM: hardware related
	Tuned *osbuild.TunedStageOptions

	// WORKLOAD
	Tmpfilesd []*osbuild.TmpfilesdStageOptions

	// WORKLOAD
	PamLimitsConf []*osbuild.PamLimitsConfStageOptions
	Sysctld       []*osbuild.SysctldStageOptions

	// WORKLOAD: usually release ver, but could be other things
	DNFConfig []*osbuild.DNFConfigStageOptions

	// WORKLOAD
	DNFAutomaticConfig *osbuild.DNFAutomaticConfigStageOptions

	// Same as DNFConfig
	YUMConfig *osbuild.YumConfigStageOptions
	YUMRepos  []*osbuild.YumReposStageOptions

	// WORKLOAD
	SshdConfig *osbuild.SshdConfigStageOptions

	// ENVIRONMENT: GCP stuff
	GCPGuestAgentConfig *osbuild.GcpGuestAgentConfigOptions

	// WORKLOAD
	AuthConfig *osbuild.AuthconfigStageOptions

	// ENVIRONMENT: used in azure
	PwQuality *osbuild.PwqualityConfStageOptions

	// WORKLOAD
	OpenSCAPConfig *osbuild.OscapRemediationStageOptions

	// ENVIRONMENT (default for image type) and WORKLOAD (overrides)
	NTPServers []osbuild.ChronyConfigServer

	// ENVIRONMENT: Azure agent
	WAAgentConfig *osbuild.WAAgentConfStageOptions

	// PLATFORM: hardware thingies
	UdevRules *osbuild.UdevRulesStageOptions

	// ENVIRONMENT (default for image type) and WORKLOAD (overrides)
	LeapSecTZ *string

	FactAPIType string

	Subscription *distro.SubscriptionImageOptions
	RHSMConfig   map[distro.RHSMSubscriptionStatus]*osbuild.RHSMStageOptions

	// VERY WORKLOAD
	// Custom directories and files to create in the image
	Directories []*fsnode.Directory
	Files       []*fsnode.File
}

// OS represents the filesystem tree of the target image. This roughly
// corresponds to the root filesystem once an instance of the image is running.
type OS struct {
	Base
	// Customizations to apply to the base OS
	OSCustomizations
	// Environment the system will run in
	Environment environment.Environment
	// Workload to install on top of the base system
	Workload workload.Workload
	// Ref of ostree commit, if empty the tree cannot be in an ostree commit
	OSTreeRef string
	// OSTree parent spec, if nil the new commit (if applicable) will have no parent
	OSTreeParent *ostree.CommitSpec
	// Partition table, if nil the tree cannot be put on a partitioned disk
	PartitionTable *disk.PartitionTable

	repos        []rpmmd.RepoConfig
	packageSpecs []rpmmd.PackageSpec
	platform     platform.Platform
	kernelVer    string

	// NoBLS configures the image bootloader with traditional menu entries
	// instead of BLS. Required for legacy systems like RHEL 7.
	NoBLS     bool
	OSProduct string
	OSVersion string
	OSNick    string
}

// NewOS creates a new OS pipeline. build is the build pipeline to use for
// building the OS pipeline. platform is the target platform for the final
// image. repos are the repositories to install RPMs from.
func NewOS(m *Manifest,
	buildPipeline *Build,
	platform platform.Platform,
	repos []rpmmd.RepoConfig) *OS {
	name := "os"
	p := &OS{
		Base:     NewBase(m, name, buildPipeline),
		repos:    filterRepos(repos, name),
		platform: platform,
	}
	buildPipeline.addDependent(p)
	m.addPipeline(p)
	return p
}

func (p *OS) getPackageSetChain() []rpmmd.PackageSet {
	packages := p.platform.GetPackages()

	// If we have a logical volume we need to include the lvm2 package.
	// OSTree-based images (commit and container) aren't bootable images and
	// don't have partition tables.
	if p.PartitionTable != nil && p.OSTreeRef == "" {
		packages = append(packages, p.PartitionTable.GetBuildPackages()...)
	}

	if p.Environment != nil {
		packages = append(packages, p.Environment.GetPackages()...)
	}

	chain := []rpmmd.PackageSet{
		{
			Include:      packages,
			Repositories: p.repos,
		},
	}

	if p.Workload != nil {
		osPackages := p.Workload.GetOSPackages()
		if len(osPackages) > 0 {
			chain[0].Include = append(chain[0].Include, osPackages...)
			chain[0].Exclude = p.Workload.GetOSExcludePackages()
			chain[0].Repositories = append(chain[0].Repositories, p.repos...)
		}
		userPackages := p.Workload.GetUserPackages()
		if len(userPackages) > 0 {
			chain = append(chain, rpmmd.PackageSet{
				Include:      userPackages,
				Repositories: append(chain[0].Repositories, p.Workload.GetUserRepos()...),
			})
		}
	}

	return chain
}

func (p *OS) getBuildPackages() []string {
	packages := p.platform.GetBuildPackages()
	if p.PartitionTable != nil {
		packages = append(packages, p.PartitionTable.GetBuildPackages()...)
	}
	packages = append(packages, "rpm")
	if p.OSTreeRef != "" {
		packages = append(packages, "rpm-ostree")
	}
	if p.SElinux != "" {
		packages = append(packages, "policycoreutils", fmt.Sprintf("selinux-policy-%s", p.SElinux))
	}
	if len(p.CloudInit) > 0 {
		packages = append(packages, "python3-pyyaml")
	}
	if len(p.DNFConfig) > 0 || len(p.RHSMConfig) > 0 {
		packages = append(packages, "python3-iniparse")
	}

	if len(p.OSCustomizations.Containers) > 0 {
		if p.OSTreeRef != "" {
			packages = append(packages, "python3-toml")
		}
		packages = append(packages, "skopeo")
	}

	return packages
}

func (p *OS) getOSTreeCommits() []ostree.CommitSpec {
	if p.OSTreeParent == nil {
		return nil
	}
	return []ostree.CommitSpec{*p.OSTreeParent}
}

func (p *OS) getPackageSpecs() []rpmmd.PackageSpec {
	return p.packageSpecs
}

func (p *OS) getContainerSpecs() []container.Spec {
	return p.Containers
}

func (p *OS) serializeStart(packages []rpmmd.PackageSpec) {
	if len(p.packageSpecs) > 0 {
		panic("double call to serializeStart()")
	}
	p.packageSpecs = packages
	if p.KernelName != "" {
		p.kernelVer = rpmmd.GetVerStrFromPackageSpecListPanic(p.packageSpecs, p.KernelName)
	}
}

func (p *OS) serializeEnd() {
	if len(p.packageSpecs) == 0 {
		panic("serializeEnd() call when serialization not in progress")
	}
	p.kernelVer = ""
	p.packageSpecs = nil
}

func (p *OS) serialize() osbuild.Pipeline {
	if len(p.packageSpecs) == 0 {
		panic("serialization not started")
	}

	pipeline := p.Base.serialize()

	if p.OSTreeRef != "" && p.OSTreeParent != nil {
		pipeline.AddStage(osbuild.NewOSTreePasswdStage("org.osbuild.source", p.OSTreeParent.Checksum))
	}

	rpmOptions := osbuild.NewRPMStageOptions(p.repos)
	if p.ExcludeDocs {
		if rpmOptions.Exclude == nil {
			rpmOptions.Exclude = &osbuild.Exclude{}
		}
		rpmOptions.Exclude.Docs = true
	}
	rpmOptions.GPGKeysFromTree = p.GPGKeyFiles
	if p.OSTreeRef != "" {
		rpmOptions.OSTreeBooted = common.ToPtr(true)
		rpmOptions.DBPath = "/usr/share/rpm"
	}
	pipeline.AddStage(osbuild.NewRPMStage(rpmOptions, osbuild.NewRpmStageSourceFilesInputs(p.packageSpecs)))

	if !p.NoBLS {
		// If the /boot is on a separate partition, the prefix for the BLS stage must be ""
		if p.PartitionTable == nil || p.PartitionTable.FindMountable("/boot") == nil {
			pipeline.AddStage(osbuild.NewFixBLSStage(&osbuild.FixBLSStageOptions{}))
		} else {
			pipeline.AddStage(osbuild.NewFixBLSStage(&osbuild.FixBLSStageOptions{Prefix: common.ToPtr("")}))
		}
	}

	if len(p.Containers) > 0 {
		images := osbuild.NewContainersInputForSources(p.Containers)

		var storagePath string

		// OSTree commits do not include data in `/var` since that is tied to the
		// deployment, rather than the commit. Therefore the containers need to be
		// stored in a different location, like `/usr/share`, and the container
		// storage engine configured accordingly.
		if p.OSTreeRef != "" {
			storagePath = "/usr/share/containers/storage"
			storageConf := "/etc/containers/storage.conf"

			containerStoreOpts := osbuild.NewContainerStorageOptions(storageConf, storagePath)
			pipeline.AddStage(osbuild.NewContainersStorageConfStage(containerStoreOpts))
		}

		skopeo := osbuild.NewSkopeoStage(images, storagePath)
		pipeline.AddStage(skopeo)
	}

	pipeline.AddStage(osbuild.NewLocaleStage(&osbuild.LocaleStageOptions{Language: p.Language}))

	if p.Keyboard != nil {
		keymapOptions := &osbuild.KeymapStageOptions{Keymap: *p.Keyboard}
		if len(p.X11KeymapLayouts) > 0 {
			keymapOptions.X11Keymap = &osbuild.X11KeymapOptions{Layouts: p.X11KeymapLayouts}
		}
		pipeline.AddStage(osbuild.NewKeymapStage(keymapOptions))
	}

	if p.Hostname != "" {
		pipeline.AddStage(osbuild.NewHostnameStage(&osbuild.HostnameStageOptions{Hostname: p.Hostname}))
	}
	pipeline.AddStage(osbuild.NewTimezoneStage(&osbuild.TimezoneStageOptions{Zone: p.Timezone}))

	if len(p.NTPServers) > 0 {
		chronyOptions := &osbuild.ChronyStageOptions{Servers: p.NTPServers}
		if p.LeapSecTZ != nil {
			chronyOptions.LeapsecTz = p.LeapSecTZ
		}
		pipeline.AddStage(osbuild.NewChronyStage(chronyOptions))
	}

	if len(p.Groups) > 0 {
		pipeline.AddStage(osbuild.GenGroupsStage(p.Groups))
	}

	if len(p.Users) > 0 {
		if p.OSTreeRef != "" {
			// for ostree, writing the key during user creation is
			// redundant and can cause issues so create users without keys
			// and write them on first boot
			usersStageSansKeys, err := osbuild.GenUsersStage(p.Users, true)
			if err != nil {
				// TODO: move encryption into weldr
				panic("password encryption failed")
			}
			pipeline.AddStage(usersStageSansKeys)
			pipeline.AddStage(osbuild.NewFirstBootStage(usersFirstBootOptions(p.Users)))
		} else {
			usersStage, err := osbuild.GenUsersStage(p.Users, false)
			if err != nil {
				// TODO: move encryption into weldr
				panic("password encryption failed")
			}
			pipeline.AddStage(usersStage)
		}
	}

	if p.Firewall != nil {
		pipeline.AddStage(osbuild.NewFirewallStage(p.Firewall))
	}

	for _, sysconfigConfig := range p.Sysconfig {
		pipeline.AddStage(osbuild.NewSysconfigStage(sysconfigConfig))
	}

	for _, systemdLogindConfig := range p.SystemdLogind {
		pipeline.AddStage(osbuild.NewSystemdLogindStage(systemdLogindConfig))
	}

	for _, cloudInitConfig := range p.CloudInit {
		pipeline.AddStage(osbuild.NewCloudInitStage(cloudInitConfig))
	}

	for _, modprobeConfig := range p.Modprobe {
		pipeline.AddStage(osbuild.NewModprobeStage(modprobeConfig))
	}

	for _, dracutConfConfig := range p.DracutConf {
		pipeline.AddStage(osbuild.NewDracutConfStage(dracutConfConfig))
	}

	for _, systemdUnitConfig := range p.SystemdUnit {
		pipeline.AddStage(osbuild.NewSystemdUnitStage(systemdUnitConfig))
	}

	if p.Authselect != nil {
		pipeline.AddStage(osbuild.NewAuthselectStage(p.Authselect))
	}

	if p.SELinuxConfig != nil {
		pipeline.AddStage(osbuild.NewSELinuxConfigStage(p.SELinuxConfig))
	}

	if p.Tuned != nil {
		pipeline.AddStage(osbuild.NewTunedStage(p.Tuned))
	}

	for _, tmpfilesdConfig := range p.Tmpfilesd {
		pipeline.AddStage(osbuild.NewTmpfilesdStage(tmpfilesdConfig))
	}

	for _, pamLimitsConfConfig := range p.PamLimitsConf {
		pipeline.AddStage(osbuild.NewPamLimitsConfStage(pamLimitsConfConfig))
	}

	for _, sysctldConfig := range p.Sysctld {
		pipeline.AddStage(osbuild.NewSysctldStage(sysctldConfig))
	}

	for _, dnfConfig := range p.DNFConfig {
		pipeline.AddStage(osbuild.NewDNFConfigStage(dnfConfig))
	}

	if p.DNFAutomaticConfig != nil {
		pipeline.AddStage(osbuild.NewDNFAutomaticConfigStage(p.DNFAutomaticConfig))
	}

	for _, yumRepo := range p.YUMRepos {
		pipeline.AddStage(osbuild.NewYumReposStage(yumRepo))
	}

	if p.YUMConfig != nil {
		pipeline.AddStage(osbuild.NewYumConfigStage(p.YUMConfig))
	}

	if p.GCPGuestAgentConfig != nil {
		pipeline.AddStage(osbuild.NewGcpGuestAgentConfigStage(p.GCPGuestAgentConfig))
	}

	if p.SshdConfig != nil {
		pipeline.AddStage((osbuild.NewSshdConfigStage(p.SshdConfig)))
	}

	if p.AuthConfig != nil {
		pipeline.AddStage(osbuild.NewAuthconfigStage(p.AuthConfig))
	}

	if p.PwQuality != nil {
		pipeline.AddStage(osbuild.NewPwqualityConfStage(p.PwQuality))
	}

	// If subscription settings are included there are 3 possible setups:
	// - Register the system with rhc and enable Insights
	// - Register with subscription-manager, no Insights or rhc
	// - Register with subscription-manager and enable Insights, no rhc
	if p.Subscription != nil {
		var commands []string
		if p.Subscription.Rhc {
			// Use rhc for registration instead of subscription manager
			commands = []string{fmt.Sprintf("/usr/bin/rhc connect -o=%s -a=%s --server %s", p.Subscription.Organization, p.Subscription.ActivationKey, p.Subscription.ServerUrl)}

			// Always enable Insights when using rhc
			commands = append(commands, "/usr/bin/insights-client --register")
		} else {
			commands = []string{fmt.Sprintf("/usr/sbin/subscription-manager register --org=%s --activationkey=%s --serverurl %s --baseurl %s", p.Subscription.Organization, p.Subscription.ActivationKey, p.Subscription.ServerUrl, p.Subscription.BaseUrl)}

			// Insights is optional when using subscription-manager
			if p.Subscription.Insights {
				commands = append(commands, "/usr/bin/insights-client --register")
			}
		}

		pipeline.AddStage(osbuild.NewFirstBootStage(&osbuild.FirstBootStageOptions{
			Commands:       commands,
			WaitForNetwork: true,
		}))

		if rhsmConfig, exists := p.RHSMConfig[distro.RHSMConfigWithSubscription]; exists {
			pipeline.AddStage(osbuild.NewRHSMStage(rhsmConfig))
		}
	} else {
		if rhsmConfig, exists := p.RHSMConfig[distro.RHSMConfigNoSubscription]; exists {
			pipeline.AddStage(osbuild.NewRHSMStage(rhsmConfig))
		}
	}

	if waConfig := p.WAAgentConfig; waConfig != nil {
		pipeline.AddStage(osbuild.NewWAAgentConfStage(waConfig))
	}

	if udevRules := p.UdevRules; udevRules != nil {
		pipeline.AddStage(osbuild.NewUdevRulesStage(udevRules))
	}

	if pt := p.PartitionTable; pt != nil {
		kernelOptions := osbuild.GenImageKernelOptions(p.PartitionTable)
		kernelOptions = append(kernelOptions, p.KernelOptionsAppend...)
		if !p.KernelOptionsBootloader {
			pipeline = prependKernelCmdlineStage(pipeline, strings.Join(kernelOptions, " "), pt)
		}

		pipeline.AddStage(osbuild.NewFSTabStage(osbuild.NewFSTabStageOptions(pt)))

		var bootloader *osbuild.Stage
		switch p.platform.GetArch() {
		case platform.ARCH_S390X:
			bootloader = osbuild.NewZiplStage(new(osbuild.ZiplStageOptions))
		default:
			if p.NoBLS {
				// BLS entries not supported: use grub2.legacy
				id := "76a22bf4-f153-4541-b6c7-0332c0dfaeac"
				product := osbuild.GRUB2Product{
					Name:    p.OSProduct,
					Version: p.OSVersion,
					Nick:    p.OSNick,
				}

				_, err := rpmmd.GetVerStrFromPackageSpecList(p.packageSpecs, "dracut-config-rescue")
				hasRescue := err == nil
				bootloader = osbuild.NewGrub2LegacyStage(
					osbuild.NewGrub2LegacyStageOptions(
						p.Grub2Config,
						p.PartitionTable,
						kernelOptions,
						p.platform.GetBIOSPlatform(),
						p.platform.GetUEFIVendor(),
						osbuild.MakeGrub2MenuEntries(id, p.kernelVer, product, hasRescue),
					),
				)
			} else {
				options := osbuild.NewGrub2StageOptionsUnified(pt,
					strings.Join(kernelOptions, " "),
					p.kernelVer,
					p.platform.GetUEFIVendor() != "",
					p.platform.GetBIOSPlatform(),
					p.platform.GetUEFIVendor(), false)
				if cfg := p.Grub2Config; cfg != nil {
					// TODO: don't store Grub2Config in OSPipeline, making the overrides unnecessary
					// grub2.Config.Default is owned and set by `NewGrub2StageOptionsUnified`
					// and thus we need to preserve it
					if options.Config != nil {
						cfg.Default = options.Config.Default
					}

					options.Config = cfg
				}
				if p.KernelOptionsBootloader {
					options.WriteCmdLine = nil
					if options.UEFI != nil {
						options.UEFI.Unified = false
					}
				}
				bootloader = osbuild.NewGRUB2Stage(options)
			}
		}

		pipeline.AddStage(bootloader)
	}

	if p.OpenSCAPConfig != nil {
		pipeline.AddStage(osbuild.NewOscapRemediationStage(p.OpenSCAPConfig))
	}

	if p.FactAPIType != "" {
		pipeline.AddStage(osbuild.NewRHSMFactsStage(&osbuild.RHSMFactsStageOptions{
			Facts: osbuild.RHSMFacts{
				ApiType: p.FactAPIType,
			},
		}))
	}

	if p.OSTreeRef != "" {
		pipeline.AddStage(osbuild.NewSystemdJournaldStage(
			&osbuild.SystemdJournaldStageOptions{
				Filename: "10-persistent.conf",
				Config: osbuild.SystemdJournaldConfigDropin{
					Journal: osbuild.SystemdJournaldConfigJournalSection{
						Storage: osbuild.StoragePresistent,
					},
				},
			}))
	}

	// First create custom directories, because some of the custom files may depend on them
	if len(p.Directories) > 0 {
		pipeline.AddStages(osbuild.GenDirectoryNodesStages(p.Directories)...)
	}

	if len(p.Files) > 0 {
		pipeline.AddStages(osbuild.GenFileNodesStages(p.Files)...)
	}

	enabledServices := []string{}
	disabledServices := []string{}
	enabledServices = append(enabledServices, p.EnabledServices...)
	disabledServices = append(disabledServices, p.DisabledServices...)
	if p.Environment != nil {
		enabledServices = append(enabledServices, p.Environment.GetServices()...)
	}
	if p.Workload != nil {
		enabledServices = append(enabledServices, p.Workload.GetServices()...)
		disabledServices = append(disabledServices, p.Workload.GetDisabledServices()...)
	}
	if len(enabledServices) != 0 ||
		len(disabledServices) != 0 || p.DefaultTarget != "" {
		pipeline.AddStage(osbuild.NewSystemdStage(&osbuild.SystemdStageOptions{
			EnabledServices:  enabledServices,
			DisabledServices: disabledServices,
			DefaultTarget:    p.DefaultTarget,
		}))
	}

	if p.SElinux != "" {
		pipeline.AddStage(osbuild.NewSELinuxStage(&osbuild.SELinuxStageOptions{
			FileContexts:     fmt.Sprintf("etc/selinux/%s/contexts/files/file_contexts", p.SElinux),
			ForceAutorelabel: p.SELinuxForceRelabel,
		}))
	}

	if p.OSTreeRef != "" {
		pipeline.AddStage(osbuild.NewOSTreePrepTreeStage(&osbuild.OSTreePrepTreeStageOptions{
			EtcGroupMembers: []string{
				// NOTE: We may want to make this configurable.
				"wheel", "docker",
			},
		}))
	}

	return pipeline
}

func prependKernelCmdlineStage(pipeline osbuild.Pipeline, kernelOptions string, pt *disk.PartitionTable) osbuild.Pipeline {
	rootFs := pt.FindMountable("/")
	if rootFs == nil {
		panic("root filesystem must be defined for kernel-cmdline stage, this is a programming error")
	}
	rootFsUUID := rootFs.GetFSSpec().UUID
	kernelStage := osbuild.NewKernelCmdlineStage(osbuild.NewKernelCmdlineStageOptions(rootFsUUID, kernelOptions))
	pipeline.Stages = append([]*osbuild.Stage{kernelStage}, pipeline.Stages...)
	return pipeline
}

func usersFirstBootOptions(users []users.User) *osbuild.FirstBootStageOptions {
	cmds := make([]string, 0, 3*len(users)+2)
	// workaround for creating authorized_keys file for user
	// need to special case the root user, which has its home in a different place
	varhome := filepath.Join("/var", "home")
	roothome := filepath.Join("/var", "roothome")

	for _, user := range users {
		if user.Key != nil {
			var home string

			if user.Name == "root" {
				home = roothome
			} else {
				home = filepath.Join(varhome, user.Name)
			}

			sshdir := filepath.Join(home, ".ssh")

			cmds = append(cmds, fmt.Sprintf("mkdir -p %s", sshdir))
			cmds = append(cmds, fmt.Sprintf("sh -c 'echo %q >> %q'", *user.Key, filepath.Join(sshdir, "authorized_keys")))
			cmds = append(cmds, fmt.Sprintf("chown %s:%s -Rc %s", user.Name, user.Name, sshdir))
		}
	}
	cmds = append(cmds, fmt.Sprintf("restorecon -rvF %s", varhome))
	cmds = append(cmds, fmt.Sprintf("restorecon -rvF %s", roothome))

	options := &osbuild.FirstBootStageOptions{
		Commands:       cmds,
		WaitForNetwork: false,
	}

	return options
}

func (p *OS) GetPlatform() platform.Platform {
	return p.platform
}

func (p *OS) getInline() []string {
	inlineData := []string{}

	// inline data for custom files
	for _, file := range p.Files {
		inlineData = append(inlineData, string(file.Data()))
	}

	return inlineData
}
