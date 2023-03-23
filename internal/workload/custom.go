package workload

import "github.com/osbuild/osbuild-composer/internal/rpmmd"

type Custom struct {
	BaseWorkload
	UserPackages     []string
	UserRepos        []rpmmd.RepoConfig
	Services         []string
	DisabledServices []string
}

func (p *Custom) GetUserPackages() []string {
	return p.UserPackages
}

func (p *Custom) GetUserRepos() []rpmmd.RepoConfig {
	return p.UserRepos
}

func (p *Custom) GetServices() []string {
	return p.Services
}

// TODO: Does this belong here? What kind of workload requires
// services to be disabled?
func (p *Custom) GetDisabledServices() []string {
	return p.DisabledServices
}
