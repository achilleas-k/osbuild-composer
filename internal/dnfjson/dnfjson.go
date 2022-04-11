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
	ModulePlatformID string

	// System architecture
	Arch string

	// Cache directory for the DNF metadata
	CacheDir string

	// Release version of the distro. This is used in repo files on the host
	// system and required for subscription support.
	releaseVer string

	// Path to the dnf-json binary and optional args (default: "osbuild-dnf-json", assumed in $PATH)
	dnfJsonCmd []string
}

// Create a new unconfigured Solver (without platform information). It should
// be configured using the SetConfig() function before use.
func NewBaseSolver(cacheDir string) *Solver {
	return &Solver{
		CacheDir:   cacheDir,
		dnfJsonCmd: []string{"osbuild-dnf-json"},
	}
}

// Create a new Solver with the given configuration
func NewSolver(modulePlatformID string, releaseVer string, arch string, cacheDir string) *Solver {
	return &Solver{
		ModulePlatformID: modulePlatformID,
		Arch:             arch,
		CacheDir:         cacheDir,
		releaseVer:       releaseVer,
		dnfJsonCmd:       []string{"osbuild-dnf-json"},
	}
}

// SetConfig sets the platform configuration values.
func (s *Solver) SetConfig(modulePlatformID string, releaseVer string, arch string) {
	s.ModulePlatformID = modulePlatformID
	s.Arch = arch
	s.releaseVer = releaseVer
}

func makeDepsolveRequest(s *Solver, pkgSets []rpmmd.PackageSet, repoSets [][]rpmmd.RepoConfig) (*Request, error) {
	args := make([]Arguments, len(pkgSets))
	for idx := range pkgSets {
		repos, err := ReposFromRPMMD(repoSets[idx], s.Arch, s.releaseVer)
		if err != nil {
			return nil, err
		}
		args[idx] = Arguments{
			PackageSpecs: pkgSets[idx].Include,
			ExcludSpecs:  pkgSets[idx].Exclude,
			Repos:        repos,
		}
	}
	req := Request{
		Command:          "depsolve",
		ModulePlatformID: s.ModulePlatformID,
		Arch:             s.Arch,
		CacheDir:         s.CacheDir,
		Arguments:        args,
	}
	return &req, nil
}

// Depsolve the given packages with explicit excludes using the solver configuration and provided repos
func (s *Solver) Depsolve(pkgSets []rpmmd.PackageSet, repoSets [][]rpmmd.RepoConfig) ([]DepsolveResult, error) {
	if len(pkgSets) != len(repoSets) {
		return nil, fmt.Errorf("error: different number of package sets and repositories: %d != %d", len(pkgSets), len(repoSets))
	}

	req, err := makeDepsolveRequest(s, pkgSets, repoSets)
	if err != nil {
		return nil, err
	}

	output, err := run(s.dnfJsonCmd, req)
	if err != nil {
		return nil, err
	}
	var result []depsolveResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	return resultsToPublic(result, repoSets), nil
}

func makeDumpRequest(s *Solver, repos []rpmmd.RepoConfig) (*Request, error) {
	dnfRepos, err := ReposFromRPMMD(repos, s.Arch, "")
	if err != nil {
		return nil, err
	}
	req := Request{
		Command:          "dump",
		ModulePlatformID: s.ModulePlatformID,
		Arch:             s.Arch,
		CacheDir:         s.CacheDir,
		Arguments: []Arguments{
			{
				Repos: dnfRepos,
			},
		},
	}
	return &req, nil
}

func (s *Solver) FetchMetadata(repos []rpmmd.RepoConfig) (*Metadata, error) {
	req, err := makeDumpRequest(s, repos)
	if err != nil {
		return nil, err
	}
	result, err := run(s.dnfJsonCmd, req)
	if err != nil {
		return nil, err
	}

	var metadata *Metadata
	if err := json.Unmarshal(result, metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}

func (s *Solver) SetDNFJSONPath(path ...string) {
	s.dnfJsonCmd = path
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

// convert a slice of internal depsolveResult to a slice of public DepsolveResult
func resultsToPublic(results []depsolveResult, repoSets [][]rpmmd.RepoConfig) []DepsolveResult {
	pubRes := make([]DepsolveResult, len(results))
	for idx := range results {
		pubRes[idx] = resultToPublic(results[idx], repoSets[idx])
	}
	return pubRes
}

// convert an internal depsolveResult to a public DepsolveResult
func resultToPublic(result depsolveResult, repos []rpmmd.RepoConfig) DepsolveResult {
	return DepsolveResult{
		Checksums:    result.Checksums,
		Dependencies: depsToRPMMD(result.Dependencies, repos),
	}
}

func depsToRPMMD(dependencies []PackageSpec, repos []rpmmd.RepoConfig) []rpmmd.PackageSpec {
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
	// Command should be either "depsolve" or "dump"
	Command string `json:"command"`

	// Platform ID, e.g., "platform:el8"
	ModulePlatformID string `json:"module_platform_id"`

	// System architecture
	Arch string `json:"arch"`

	// Cache directory for the DNF metadata
	CacheDir string `json:"cachedir"`

	// One or more arguments for the action defined by Command
	Arguments []Arguments `json:"arguments"`
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

// Private version of the depsolve result.  Uses a slightly different
// PackageSpec than the public one that uses the rpmmd type.
type depsolveResult struct {
	// Repository checksums
	Checksums map[string]string `json:"checksums"`

	// Resolved package dependencies
	Dependencies []PackageSpec `json:"dependencies"`
}

// DepsolveResult is the result returned from a Depsolve call.
type DepsolveResult struct {
	// Repository checksums
	Checksums map[string]string

	// Resolved package dependencies
	Dependencies []rpmmd.PackageSpec
}

// Metadata is the result returned from a FetchMetadata call.
type Metadata struct {
	Checksums map[string]string `json:"checksums"`
	Packages  rpmmd.PackageList `json:"packages"`
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

func ParseError(data []byte) Error {
	var e Error
	if err := json.Unmarshal(data, &e); err != nil {
		// dumping the error into the Reason can get noisy, but it's good for troubleshooting
		return Error{
			Kind:   "InternalError",
			Reason: fmt.Sprintf("Failed to unmarshal dnf-json error output %q: %s", string(data), err.Error()),
		}
	}
	return e
}

// Depsolve the given packages with explicit excludes using the given configuration and repos
func Depsolve(pkgSets []rpmmd.PackageSet, repoSets [][]rpmmd.RepoConfig, modulePlatformID string, releaseVer string, arch string, cacheDir string) ([]DepsolveResult, error) {
	return NewSolver(modulePlatformID, releaseVer, arch, cacheDir).Depsolve(pkgSets, repoSets)
}

func FetchMetadata(repos []rpmmd.RepoConfig, modulePlatformID string, releaseVer string, arch string, cacheDir string) (*Metadata, error) {
	return NewSolver(modulePlatformID, releaseVer, arch, cacheDir).FetchMetadata(repos)
}

func run(dnfJsonCmd []string, req *Request) ([]byte, error) {
	if len(dnfJsonCmd) == 0 {
		return nil, fmt.Errorf("dnf-json command undefined")
	}
	ex := dnfJsonCmd[0]
	args := make([]string, len(dnfJsonCmd)-1)
	if len(dnfJsonCmd) > 1 {
		args = dnfJsonCmd[1:]
	}
	cmd := exec.Command(ex, args...)
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
		return nil, ParseError(output)
	}

	return output, nil
}
