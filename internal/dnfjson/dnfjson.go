package dnfjson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/osbuild/osbuild-composer/internal/rhsm"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

// Solver is configured with system information in order to resolve
// dependencies for RPM packages using DNF.
type Solver struct {
	// Platform ID, e.g., "platform:el8"
	ModulePlatformID string `json:"module_platform_id"`

	// System architecture
	Arch string `json:"arch"`

	// Cache directory for the DNF metadata
	CacheDir string `json:"cachedir"`
}

// Create a new Solver with the given configuration
func NewSolver(modulePlatformID string, arch string, cacheDir string) *Solver {
	return &Solver{
		ModulePlatformID: modulePlatformID,
		Arch:             arch,
		CacheDir:         cacheDir,
	}
}

// Depsolve the given packages with explicit excludes using the solver configuration and provided repos
func (s *Solver) Depsolve(includes []string, excludes []string, repos []RepoConfig) (*Result, error) {
	req := Request{
		Command: "depsolve",
		Solver:  s,
		Arguments: Arguments{
			PackageSpecs: includes,
			ExcludSpecs:  excludes,
			Repos:        repos,
		},
	}
	return run(req)
}

func (s *Solver) FetchMetadata(repos []RepoConfig) (*Result, error) {
	req := Request{
		Command: "dump",
		Solver:  s,
		Arguments: Arguments{
			Repos: repos,
		},
	}
	return run(req)
}

// Repository configuration for resolving dependencies for a set of packages. A
// Solver needs at least one RPM repository configured to be able to depsolve.
type RepoConfig struct {
	ID             string `json:"id"`
	Name           string `json:"name,omitempty"`
	BaseURL        string `json:"baseurl,omitempty"`
	Metalink       string `json:"metalink,omitempty"`
	MirrorList     string `json:"mirrorlist,omitempty"`
	GPGKey         string `json:"gpgkey,omitempty"`
	IgnoreSSL      bool   `json:"ignoressl"`
	SSLCACert      string `json:"sslcacert,omitempty"`
	SSLClientKey   string `json:"sslclientkey,omitempty"`
	SSLClientCert  string `json:"sslclientcert,omitempty"`
	MetadataExpire string `json:"metadata_expire,omitempty"`
}

// ReposFromRPMMD converts an rpmmd.RepoConfig to a RepoConfig. If the
// repository requires a subscription, the system subscriptions are loaded and
// included in the new configs.
func ReposFromRPMMD(rpmRepos []rpmmd.RepoConfig, arch string, releaseVer string) ([]RepoConfig, error) {
	subscriptions, _ := rhsm.LoadSystemSubscriptions()
	dnfRepos := make([]RepoConfig, len(rpmRepos))
	for idx, rr := range rpmRepos {
		id := strconv.Itoa(idx)
		dr := RepoConfig{
			ID:             id,
			Name:           rr.Name,
			BaseURL:        rr.BaseURL,
			Metalink:       rr.Metalink,
			MirrorList:     rr.MirrorList,
			GPGKey:         rr.GPGKey,
			IgnoreSSL:      rr.IgnoreSSL,
			MetadataExpire: rr.MetadataExpire,
		}
		if rr.RHSM {
			if subscriptions == nil {
				return nil, fmt.Errorf("This system does not have any valid subscriptions. Subscribe it before specifying rhsm: true in sources.")
			}
			secrets, err := subscriptions.GetSecretsForBaseurl(rr.BaseURL, arch, releaseVer)
			if err != nil {
				return nil, fmt.Errorf("RHSM secrets not found on the host for this baseurl: %s", rr.BaseURL)
			}
			dr.SSLCACert = secrets.SSLCACert
			dr.SSLClientKey = secrets.SSLClientKey
			dr.SSLClientCert = secrets.SSLClientCert

		}
		dnfRepos[idx] = dr
	}
	return dnfRepos, nil
}

func DepsToRPMMD(dependencies []PackageSpec, repos []rpmmd.RepoConfig) []rpmmd.PackageSpec {
	rpmDependencies := make([]rpmmd.PackageSpec, len(dependencies))
	for i, dep := range dependencies {
		id, err := strconv.Atoi(dep.RepoID)
		if err != nil {
			panic(err)
		}
		repo := repos[id]
		dep := dependencies[i]
		rpmDependencies[i].Name = dep.Name
		rpmDependencies[i].Epoch = dep.Epoch
		rpmDependencies[i].Version = dep.Version
		rpmDependencies[i].Release = dep.Release
		rpmDependencies[i].Arch = dep.Arch
		rpmDependencies[i].RemoteLocation = dep.RemoteLocation
		rpmDependencies[i].Checksum = dep.Checksum
		rpmDependencies[i].CheckGPG = repo.CheckGPG
		if repo.RHSM {
			rpmDependencies[i].Secrets = "org.osbuild.rhsm"
		}
	}
	return rpmDependencies
}

// Request command and arguments for dnf-json
type Request struct {
	Command string `json:"command"`
	*Solver
	Arguments Arguments `json:"arguments"`
}

// Arguments for a dnf-json request
type Arguments struct {
	// Repositories to use for depsolving
	Repos []RepoConfig `json:"repos"`

	// Packages to depsolve
	PackageSpecs []string `json:"package-specs"`

	// Packages to exclude from results
	ExcludSpecs []string `json:"exclude-specs"`
}

// Result of a dnf-json depsolve run
type Result struct {
	// Repository checksums
	Checksums map[string]string `json:"checksums"`

	// Resolved package dependencies
	Dependencies []PackageSpec `json:"dependencies"`
}

// Package specification
type PackageSpec struct {
	Name           string `json:"name"`
	Epoch          uint   `json:"epoch"`
	Version        string `json:"version,omitempty"`
	Release        string `json:"release,omitempty"`
	Arch           string `json:"arch,omitempty"`
	RepoID         string `json:"repo_id,omitempty"`
	Path           string `json:"path,omitempty"`
	RemoteLocation string `json:"remote_location,omitempty"`
	Checksum       string `json:"checksum,omitempty"`
	Secrets        string `json:"secrets,omitempty"`
}

// dnf-json error structure
type Error struct {
	Kind   string `json:"kind"`
	Reason string `json:"reason"`
}

func (err Error) Error() string {
	return fmt.Sprintf("DNF error occurred: %s: %s", err.Kind, err.Reason)
}

// Depsolve the given packages with explicit excludes using the given configuration and repos
func Depsolve(packages []string, excludes []string, repos []RepoConfig, modulePlatformID string, arch string, cacheDir string) (*Result, error) {
	req := Request{
		Command: "depsolve",
		Solver: &Solver{
			ModulePlatformID: modulePlatformID,
			Arch:             arch,
			CacheDir:         cacheDir,
		},
		Arguments: Arguments{
			PackageSpecs: packages,
			ExcludSpecs:  excludes,
			Repos:        repos,
		},
	}
	return run(req)
}

func FetchMetadata(repos []RepoConfig, modulePlatformID string, arch string, cacheDir string) (*Result, error) {
	req := Request{
		Command: "dump",
		Solver: &Solver{
			ModulePlatformID: modulePlatformID,
			Arch:             arch,
			CacheDir:         cacheDir,
		},
		Arguments: Arguments{
			Repos: repos,
		},
	}
	return run(req)
}

func run(req Request) (*Result, error) {
	cmd := exec.Command("osbuild-dnf-json")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	err = json.NewEncoder(stdin).Encode(req)
	if err != nil {
		return nil, err
	}
	stdin.Close()

	output, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if runError, ok := err.(*exec.ExitError); ok && runError.ExitCode() != 0 {
		return nil, err
	}

	res := new(Result)
	if err := json.Unmarshal(output, res); err != nil {
		return nil, err
	}

	return res, nil
}
