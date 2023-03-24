package environment

import (
	"github.com/osbuild/osbuild-composer/internal/osbuild"
)

// TODO
type GCP struct {
	BaseEnvironment
}

func (p *GCP) GetNTPConfig() ([]osbuild.ChronyConfigServer, *string) {
	if p == nil {
		return nil, nil
	}
	ntp := []osbuild.ChronyConfigServer{{Hostname: "metadata.google.internal"}}
	return ntp, nil

}

func NewGCP() *GCP {
	return &GCP{}
}
