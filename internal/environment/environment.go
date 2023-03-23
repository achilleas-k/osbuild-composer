package environment

import (
	"github.com/osbuild/osbuild-composer/internal/osbuild"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

type Environment interface {
	GetPackages() []string
	GetRepos() []rpmmd.RepoConfig
	GetServices() []string

	// TODO: replace type with internal (not osbuild)
	GetNTPConfig() ([]osbuild.ChronyConfigServer, *string)
}

type BaseEnvironment struct {
	Repos []rpmmd.RepoConfig
}

func (p BaseEnvironment) GetPackages() []string {
	return []string{}
}

func (p BaseEnvironment) GetRepos() []rpmmd.RepoConfig {
	return p.Repos
}

func (p BaseEnvironment) GetServices() []string {
	return []string{}
}

func (p BaseEnvironment) GetNTPConfig() ([]osbuild.ChronyConfigServer, *string) {
	return nil, nil
}
