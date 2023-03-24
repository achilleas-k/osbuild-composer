package environment

import (
	"github.com/osbuild/osbuild-composer/internal/common"
	"github.com/osbuild/osbuild-composer/internal/osbuild"
)

type EC2 struct {
	BaseEnvironment
}

func (p *EC2) GetPackages() []string {
	return []string{"cloud-init"}
}

func (p *EC2) GetServices() []string {
	return []string{
		"cloud-init",
		"cloud-init-local",
		"cloud-config",
		"cloud-final",
	}
}

func (p *EC2) GetNTPConfig() ([]osbuild.ChronyConfigServer, *string) {
	if p == nil {
		return nil, nil
	}
	ntp := []osbuild.ChronyConfigServer{
		{
			Hostname: "169.254.169.123",
			Prefer:   common.ToPtr(true),
			Iburst:   common.ToPtr(true),
			Minpoll:  common.ToPtr(4),
			Maxpoll:  common.ToPtr(4),
		},
	}
	return ntp, common.ToPtr("")
}

func NewEC2() *EC2 {
	return &EC2{}
}
