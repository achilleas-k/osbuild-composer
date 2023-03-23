package workload

import (
	"fmt"

	"github.com/osbuild/osbuild-composer/internal/blueprint"
	"github.com/osbuild/osbuild-composer/internal/distro"
	"github.com/osbuild/osbuild-composer/internal/osbuild"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

type Custom struct {
	BaseWorkload
	UserPackages     []string
	UserRepos        []rpmmd.RepoConfig
	Services         []string
	DisabledServices []string

	// KernelName indicates that a kernel is installed, and names the kernel
	// package.
	KernelName string

	// TODO: replace type with internal (not osbuild)
	NTPServers     []osbuild.ChronyConfigServer
	LeapSecTZ      *string
	OpenSCAPConfig *osbuild.OscapRemediationStageOptions
	SElinux        string
	Subscription   *distro.SubscriptionImageOptions
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

func NewCustomWorkload(customizations *blueprint.Customizations, extraBasePackages, userPackages rpmmd.PackageSet) *Custom {
	_, ntpCustomizations := customizations.GetTimezoneSettings()
	ntpServers := make([]osbuild.ChronyConfigServer, len(ntpCustomizations))
	for idx, server := range ntpCustomizations {
		ntpServers[idx] = osbuild.ChronyConfigServer{Hostname: server}
	}
	var openSCAP *osbuild.OscapRemediationStageOptions
	if oscapConfig := customizations.GetOpenSCAP(); oscapConfig != nil {
		openSCAP = osbuild.NewOscapRemediationStageOptions(
			osbuild.OscapConfig{
				Datastream: oscapConfig.DataStream,
				ProfileID:  oscapConfig.ProfileID,
			},
		)
	}

	w := &Custom{
		// user packages should be depsolved with both user repos and base repos
		UserRepos:      append(extraBasePackages.Repositories, userPackages.Repositories...),
		UserPackages:   userPackages.Include,
		KernelName:     customizations.GetKernel().Name,
		NTPServers:     ntpServers,
		OpenSCAPConfig: openSCAP,
		SElinux:        "targeted",
		BaseWorkload: BaseWorkload{
			OSPackages:        extraBasePackages.Include,
			OSExcludePackages: extraBasePackages.Exclude,
			OSRepos:           extraBasePackages.Repositories,
		},
	}
	if services := customizations.GetServices(); services != nil {
		w.Services = services.Enabled
		w.DisabledServices = services.Disabled
	}
	return w
}

func (p Custom) GetOSPackages() []string {
	packages := p.OSPackages
	if p.KernelName != "" {
		packages = append(packages, p.KernelName)
	}
	if len(p.NTPServers) > 0 {
		packages = append(packages, "chrony")
	}
	if p.OpenSCAPConfig != nil {
		packages = append(packages, "openscap-scanner", "scap-security-guide")
	}
	if p.SElinux != "" {
		packages = append(packages, fmt.Sprintf("selinux-policy-%s", p.SElinux))
	}

	// Make sure the right packages are included for subscriptions
	// rhc always uses insights, and depends on subscription-manager
	// non-rhc uses subscription-manager and optionally includes Insights
	if p.Subscription != nil {
		packages = append(packages, "subscription-manager")
		if p.Subscription.Rhc {
			packages = append(packages, "rhc", "insights-client")
		} else if p.Subscription.Insights {
			packages = append(packages, "insights-client")
		}
	}
	return packages
}

func (p Custom) GetKernelName() string {
	return p.KernelName
}

func (p Custom) GetNTPConfig() ([]osbuild.ChronyConfigServer, *string) {
	return p.NTPServers, p.LeapSecTZ
}

func (p Custom) GetSubscription() *distro.SubscriptionImageOptions {
	return p.Subscription
}

func (p Custom) GetOSCAPConfig() *osbuild.OscapRemediationStageOptions {
	return p.OpenSCAPConfig
}
