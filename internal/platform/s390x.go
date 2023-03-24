package platform

type S390X struct {
	BasePlatform
	BIOS bool
}

func (p *S390X) GetArch() Arch {
	return ARCH_S390X
}

func (p *S390X) GetPackages() []string {
	return []string{
		"dracut-config-generic",
		"s390utils-base",
		"s390utils-core",
	}
}

func (p *S390X) GetBuildPackages() []string {
	return []string{
		"s390utils-base",
	}
}

type S390XUnbootable struct {
	BasePlatform
}

func (p *S390XUnbootable) Bootable() bool {
	return false
}

func (p *S390XUnbootable) GetArch() Arch {
	return ARCH_AARCH64
}

func (p *S390XUnbootable) GetBuildPackages() []string {
	// TODO: remove these for unbootable?
	return []string{
		"s390utils-base",
	}
}

func (p *S390XUnbootable) GetPackages() []string {
	// TODO: remove these for unbootable?
	return []string{
		"dracut-config-generic",
		"s390utils-base",
		"s390utils-core",
	}
}
