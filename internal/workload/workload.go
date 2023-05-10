package workload

import (
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

type Workload interface {
	GetPackages() []string
	GetRepos() []rpmmd.RepoConfig
	GetServices() []string
	GetDisabledServices() []string

	SetRepos(repos []rpmmd.RepoConfig)
}

type BaseWorkload struct {
	Repos []rpmmd.RepoConfig
}

func (p BaseWorkload) GetPackages() []string {
	return []string{}
}

func (p BaseWorkload) GetRepos() []rpmmd.RepoConfig {
	return p.Repos
}

func (p BaseWorkload) GetServices() []string {
	return []string{}
}

func (p BaseWorkload) GetDisabledServices() []string {
	return []string{}
}

func (p *BaseWorkload) SetRepos(repos []rpmmd.RepoConfig) {
	p.Repos = repos
}
