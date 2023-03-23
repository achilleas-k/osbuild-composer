package workload

import "github.com/osbuild/osbuild-composer/internal/rpmmd"

type Workload interface {
	GetPackages() []string
	GetOSPackages() []string
	GetOSExcludePackages() []string
	GetUserPackages() []string
	GetRepos() []rpmmd.RepoConfig
	GetOSRepos() []rpmmd.RepoConfig
	GetUserRepos() []rpmmd.RepoConfig
	GetServices() []string
	GetDisabledServices() []string
}

type BaseWorkload struct {
	OSPackages        []string
	OSExcludePackages []string
	OSRepos           []rpmmd.RepoConfig
}

func (p BaseWorkload) GetPackages() []string {
	// TODO: Remove in favour of GetUserPackages() and distinguish from OS packages
	return []string{}
}

func (p BaseWorkload) GetRepos() []rpmmd.RepoConfig {
	// TODO: Remove in favour of GetUserRepos() and distinguish from OS repos
	return p.GetUserRepos()
}

func (p BaseWorkload) GetServices() []string {
	return []string{}
}

func (p BaseWorkload) GetDisabledServices() []string {
	return []string{}
}

func (p BaseWorkload) GetOSPackages() []string {
	return p.OSPackages
}

func (p BaseWorkload) GetOSExcludePackages() []string {
	return p.OSExcludePackages
}

func (p BaseWorkload) GetUserPackages() []string {
	return nil
}

func (p BaseWorkload) GetOSRepos() []rpmmd.RepoConfig {
	return p.OSRepos
}

func (p BaseWorkload) GetUserRepos() []rpmmd.RepoConfig {
	return nil
}
