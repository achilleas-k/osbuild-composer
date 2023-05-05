package distro

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/osbuild/osbuild-composer/internal/container"
	"github.com/osbuild/osbuild-composer/internal/crypt"
	"github.com/osbuild/osbuild-composer/internal/disk"
	"github.com/osbuild/osbuild-composer/internal/fsnode"
	"github.com/osbuild/osbuild-composer/internal/ignition"
	"github.com/osbuild/osbuild-composer/internal/ostree"
	"github.com/osbuild/osbuild-composer/internal/pathpolicy"
	"github.com/osbuild/osbuild-composer/internal/rhsm"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
	"github.com/osbuild/osbuild-composer/internal/subscription"
	"github.com/osbuild/osbuild-composer/internal/users"
)

// The ImageOptions specify options for a specific image build
type ImageOptions struct {
	// Size specifies the minimum file size for the image artifact.
	Size uint64
	// OSTree is the ostree commit to build or to pull and embed, depending on
	// the image type.
	OSTree *ostree.ImageOptions
	// Subscription configures the subscription manager.
	Subscription *subscription.ImageOptions
	// Facts embed an RHSM fact.
	Facts *rhsm.FactsImageOptions
	// Hostname for the image.
	Hostname *string
	// Kernel name and boot options.
	Kernel *Kernel
	// Public SSH keys to embed.
	SSHKey []SSHKey
	// Users to create during the build.
	Users []users.User
	// Groups to create during the build.
	Groups []users.Group
	// Timezone and NTP configuration.
	Timezone *Timezone
	// The default Locale.
	Locale *Locale
	// Firewall related configuration.
	Firewall *Firewall

	// Services to enable or disable.
	Services *Services
	// Custom partitions specified as separate mountpoints.
	Filesystem []disk.MountpointOption
	// Installation device for the coreos-installer.
	InstallationDevice string
	// FDO related options.
	FDO *FDO
	// OpenSCAP security policies (profiles).
	OpenSCAP *OpenSCAP
	// Ignition related options.
	Ignition *ignition.ImageOptions
	// Directories to create in the image.
	Directories []fsnode.Directory
	// Files with data to embed in the image.
	Files []fsnode.File

	Packages      []Package
	Modules       []Package
	PackageGroups []string
	// Containers is a list of source specifications for containers to fetch
	// and embed in the image.
	Containers []container.SourceSpec

	// Repositories is a list of repository configurations to write to the
	// image to be used for updates at runtime.
	Repositories ImageRepos
}

// A Package specifies an RPM package and optionally a specific version.
type Package struct {
	Name    string
	Version string
}

type FDO struct {
	ManufacturingServerURL string
	DiunPubKeyInsecure     string
	// This is the output of:
	// echo "sha256:$(openssl x509 -fingerprint -sha256 -noout -in diun_cert.pem | cut -d"=" -f2 | sed 's/://g')"
	DiunPubKeyHash      string
	DiunPubKeyRootCerts string
}

type Kernel struct {
	Name   string
	Append string
}

type SSHKey struct {
	User string
	Key  string
}

type Timezone struct {
	Timezone   *string
	NTPServers []string
}

type Locale struct {
	Languages []string
	Keyboard  *string
}

type Firewall struct {
	Ports    []string
	Services *FirewallServices
	Zones    []FirewallZone
}

type FirewallZone struct {
	Name    *string
	Sources []string
}

type FirewallServices struct {
	Enabled  []string
	Disabled []string
}

type Services struct {
	Enabled  []string
	Disabled []string
}

type OpenSCAP struct {
	DataStream string
	ProfileID  string
}

type RepoOptions struct {
	rpmmd.RepoConfig
	Filename string
}

type ImageRepos []RepoOptions

const repoFilenameRegex = "^[\\w.-]{1,250}\\.repo$"

func (repos ImageRepos) Validate() error {
	for _, repo := range repos {
		filenameRegex := regexp.MustCompile(repoFilenameRegex)
		if !filenameRegex.MatchString(repo.Filename) {
			return fmt.Errorf("Repository filename %q is invalid", repo.Filename)
		}

		return repo.RepoConfig.Validate()
	}
	return nil
}

func (options ImageOptions) GetHostname() *string {
	return options.Hostname
}

func (options ImageOptions) GetPrimaryLocale() (*string, *string) {
	if options.Locale == nil {
		return nil, nil
	}
	if len(options.Locale.Languages) == 0 {
		return nil, options.Locale.Keyboard
	}
	return &options.Locale.Languages[0], options.Locale.Keyboard
}

func (options ImageOptions) GetTimezoneSettings() (*string, []string) {
	if options.Timezone == nil {
		return nil, nil
	}
	return options.Timezone.Timezone, options.Timezone.NTPServers
}

func (options ImageOptions) GetUsers() []users.User {
	userOptions := []users.User{}

	// prepend sshkey for backwards compat (overridden by users)
	if len(options.SSHKey) > 0 {
		for _, options := range options.SSHKey {
			userOptions = append(userOptions, users.User{
				Name: options.User,
				Key:  &options.Key,
			})
		}
	}

	userOptions = append(userOptions, options.Users...)

	// sanitize user home directory: if it has a trailing slash,
	// it might lead to the directory not getting the correct selinux labels
	for idx := range userOptions {
		u := userOptions[idx]
		if u.Home != nil {
			homedir := strings.TrimRight(*u.Home, "/")
			u.Home = &homedir
			userOptions[idx] = u
		}
	}
	return userOptions
}

func (options ImageOptions) GetGroups() []users.Group {
	return options.Groups
}

func (options ImageOptions) GetKernel() *Kernel {
	var name string
	var append string
	if options.Kernel != nil {
		name = options.Kernel.Name
		append = options.Kernel.Append
	}

	if name == "" {
		name = "kernel"
	}

	return &Kernel{
		Name:   name,
		Append: append,
	}
}

func (options ImageOptions) GetFirewall() *Firewall {
	return options.Firewall
}

func (options ImageOptions) GetServices() *Services {
	return options.Services
}

func (options ImageOptions) GetFilesystems() []disk.MountpointOption {
	return options.Filesystem
}

func (options ImageOptions) GetFilesystemsMinSize() uint64 {
	var agg uint64
	for _, m := range options.Filesystem {
		agg += m.MinSize
	}
	// This ensures that file system customization `size` is a multiple of
	// sector size (512)
	if agg%512 != 0 {
		agg = (agg/512 + 1) * 512
	}
	return agg
}

func (options ImageOptions) GetInstallationDevice() string {
	return options.InstallationDevice
}

func (options ImageOptions) GetFDO() *FDO {
	return options.FDO
}

func (options ImageOptions) GetOpenSCAP() *OpenSCAP {
	return options.OpenSCAP
}

func (options ImageOptions) GetIgnition() *ignition.ImageOptions {
	return options.Ignition
}

func (options ImageOptions) GetDirectories() []fsnode.Directory {
	return options.Directories
}

func (options ImageOptions) GetFiles() []fsnode.File {
	return options.Files
}

// CheckAllowed returns an error if the options contain any customizations not
// specified in the arguments.
func (options ImageOptions) CheckAllowed(allowed ...string) error {

	allowMap := make(map[string]bool)

	for _, a := range allowed {
		allowMap[a] = true
	}

	t := reflect.TypeOf(options)
	v := reflect.ValueOf(options)

	for i := 0; i < t.NumField(); i++ {

		empty := false
		field := v.Field(i)

		switch field.Kind() {
		case reflect.String:
			if field.String() == "" {
				empty = true
			}
		case reflect.Array, reflect.Slice:
			if field.Len() == 0 {
				empty = true
			}
		case reflect.Ptr:
			if field.IsNil() {
				empty = true
			}
		default:
			panic(fmt.Sprintf("unhandled customization field type %s, %s", v.Kind(), t.Field(i).Name))

		}

		if !empty && !allowMap[t.Field(i).Name] {
			return fmt.Errorf("'%s' is not allowed", t.Field(i).Name)
		}
	}

	return nil
}

// CheckMountpointsPolicy checks if the mountpoints are allowed by the policy
func CheckMountpointsPolicy(mountpoints []disk.MountpointOption, mountpointAllowList *pathpolicy.PathPolicies) error {
	invalidMountpoints := []string{}
	for _, m := range mountpoints {
		err := mountpointAllowList.Check(m.Mountpoint)
		if err != nil {
			invalidMountpoints = append(invalidMountpoints, m.Mountpoint)
		}
	}

	if len(invalidMountpoints) > 0 {
		return fmt.Errorf("The following custom mountpoints are not supported %+q", invalidMountpoints)
	}

	return nil
}

// packages, modules, and groups all resolve to rpm packages right now. This
// function returns a combined list of "name-version" strings.
func (options ImageOptions) GetPackages() []string {
	return options.GetPackagesEx(true)
}

func (options ImageOptions) GetPackagesEx(bootable bool) []string {
	packages := []string{}
	for _, pkg := range options.Packages {
		packages = append(packages, pkg.ToNameVersion())
	}
	for _, pkg := range options.Modules {
		packages = append(packages, pkg.ToNameVersion())
	}
	for _, group := range options.PackageGroups {
		packages = append(packages, "@"+group)
	}

	if bootable {
		kc := options.GetKernel()
		kpkg := Package{Name: kc.Name}
		packages = append(packages, kpkg.ToNameVersion())
	}

	return packages
}

func (p Package) ToNameVersion() string {
	// Omit version to prevent all packages with prefix of name to be installed
	if p.Version == "*" || p.Version == "" {
		return p.Name
	}

	return p.Name + "-" + p.Version
}

// CryptPasswords ensures that all blueprint passwords are hashed
func (options ImageOptions) CryptPasswords() error {
	// Any passwords for users?
	for i := range options.Users {
		// Missing or empty password
		if options.Users[i].Password == nil {
			continue
		}

		// Prevent empty password from being hashed
		if len(*options.Users[i].Password) == 0 {
			options.Users[i].Password = nil
			continue
		}

		if !crypt.PasswordIsCrypted(*options.Users[i].Password) {
			pw, err := crypt.CryptSHA512(*options.Users[i].Password)
			if err != nil {
				return err
			}

			// Replace the password with the
			options.Users[i].Password = &pw
		}
	}

	return nil
}

func RepoCustomizationsToRepoConfigAndGPGKeyFiles(repos ImageRepos) (map[string][]rpmmd.RepoConfig, []*fsnode.File, error) {
	if len(repos) == 0 {
		return nil, nil, nil
	}

	repoMap := make(map[string][]rpmmd.RepoConfig, len(repos))
	var gpgKeyFiles []*fsnode.File
	for _, repo := range repos {
		filename := repo.Filename
		convertedRepo := repo.customRepoToRepoConfig()

		// convert any inline gpgkeys to fsnode.File and
		// replace the gpgkey with the file path
		for idx, gpgkey := range repo.GPGKeys {
			if _, ok := url.ParseRequestURI(gpgkey); ok != nil {
				// create the file path
				path := fmt.Sprintf("/etc/pki/rpm-gpg/RPM-GPG-KEY-%s-%d", repo.Id, idx)
				// replace the gpgkey with the file path
				convertedRepo.GPGKeys[idx] = fmt.Sprintf("file://%s", path)
				// create the fsnode for the gpgkey keyFile
				keyFile, err := fsnode.NewFile(path, nil, nil, nil, []byte(gpgkey))
				if err != nil {
					return nil, nil, err
				}
				gpgKeyFiles = append(gpgKeyFiles, keyFile)
			}
		}

		repoMap[filename] = append(repoMap[filename], convertedRepo)
	}

	return repoMap, gpgKeyFiles, nil
}

func (repo RepoOptions) customRepoToRepoConfig() rpmmd.RepoConfig {
	urls := make([]string, len(repo.BaseURLs))
	copy(urls, repo.BaseURLs)

	keys := make([]string, len(repo.GPGKeys))
	copy(keys, repo.GPGKeys)

	repoConfig := rpmmd.RepoConfig{
		Id:           repo.Id,
		BaseURLs:     urls,
		GPGKeys:      keys,
		Name:         repo.Name,
		Metalink:     repo.Metalink,
		MirrorList:   repo.MirrorList,
		CheckGPG:     repo.CheckGPG,
		CheckRepoGPG: repo.CheckRepoGPG,
		Priority:     repo.Priority,
		Enabled:      repo.Enabled,
		IgnoreSSL:    repo.IgnoreSSL,
	}

	return repoConfig
}
