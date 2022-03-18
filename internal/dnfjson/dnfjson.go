package dnfjson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// Solver is configured with a set of repositories and system information in
// order to resolve dependencies for RPM packages using DNF.
type Solver struct {
	// Repositories to use for depsolving
	Repos []RepoConfig `json:"repos"`

	// Platform ID, e.g., "platform:el8"
	ModulePlatformID string `json:"module_platform_id"`

	// System architecture
	Arch string `json:"arch"`

	// Cache directory for the DNF metadata
	CacheDir string `json:"cachedir"`
}

// Create a new Solver with the given configuration
func NewSolver(repos []RepoConfig, modulePlatformID string, arch string, cacheDir string) *Solver {
	return &Solver{
		Repos:            repos,
		ModulePlatformID: modulePlatformID,
		Arch:             arch,
		CacheDir:         cacheDir,
	}
}

// Depsolve the given packages with explicit excludes using the solver configuration
func (s *Solver) Depsolve(includes []string, excludes []string) (*Result, error) {
	req := Request{
		Command: "depsolve",
		Arguments: Arguments{
			PackageSpecs: includes,
			ExcludSpecs:  excludes,
			Solver:       *s,
		},
	}
	return run(req)
}

func (s *Solver) FetchMetadata() (*Result, error) {
	req := Request{
		Command: "dump",
		Arguments: Arguments{
			Solver: *s,
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

// Request command and arguments for dnf-json
type Request struct {
	Command   string    `json:"command"`
	Arguments Arguments `json:"arguments"`
}

// Arguments for a dnf-json request
type Arguments struct {
	// Packages to depsolve
	PackageSpecs []string `json:"package-specs"`

	// Packages to exclude from results
	ExcludSpecs []string `json:"exclude-specs"`

	// Solver configuration
	Solver
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

// Depsolve the given packages with explicit excludes using the given configuration
func Depsolve(packages []string, excludes []string, repos []RepoConfig, modulePlatformID string, arch string, cacheDir string) (*Result, error) {
	req := Request{
		Command: "depsolve",
		Arguments: Arguments{
			PackageSpecs: packages,
			ExcludSpecs:  excludes,
			Solver: Solver{
				Repos:            repos,
				ModulePlatformID: modulePlatformID,
				Arch:             arch,
				CacheDir:         cacheDir,
			},
		},
	}
	return run(req)
}

func FetchMetadata(repos []RepoConfig, modulePlatformID string, arch string, cacheDir string) (*Result, error) {
	req := Request{
		Command: "dump",
		Arguments: Arguments{
			Solver: Solver{
				Repos:            repos,
				ModulePlatformID: modulePlatformID,
				Arch:             arch,
				CacheDir:         cacheDir,
			},
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
