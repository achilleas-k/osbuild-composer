package rhel85

import (
	"errors"
	"fmt"
	"sort"

	"github.com/osbuild/osbuild-composer/internal/blueprint"
	"github.com/osbuild/osbuild-composer/internal/distro"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

const name = "rhel-85"
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
	return nil
}

func (t *imageType) Exports() []string {
	if len(t.exports) > 0 {
		return t.exports
	}
	return []string{"assembler"}
}

func (t *imageType) Manifest(b *blueprint.Customizations, options distro.ImageOptions, repos []rpmmd.RepoConfig, packageSpecSets map[string][]rpmmd.PackageSpec, seed int64) (distro.Manifest, error) {
	return nil, nil
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
		exports:         []string{"ostree-commit"},
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
		exports:         []string{"ostree-commit"},
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
