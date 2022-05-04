// dnfjson_mock provides data and methods for testing the dnfjson package.
package dnfjson_mock

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/osbuild/osbuild-composer/internal/dnfjson"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

func generatePackageList() rpmmd.PackageList {
	baseTime, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	if err != nil {
		panic(err)
	}

	var packageList rpmmd.PackageList

	for i := 0; i < 22; i++ {
		basePackage := rpmmd.Package{
			Name:        fmt.Sprintf("package%d", i),
			Summary:     fmt.Sprintf("pkg%d sum", i),
			Description: fmt.Sprintf("pkg%d desc", i),
			URL:         fmt.Sprintf("https://pkg%d.example.com", i),
			Epoch:       0,
			Version:     fmt.Sprintf("%d.0", i),
			Release:     fmt.Sprintf("%d.fc30", i),
			Arch:        "x86_64",
			BuildTime:   baseTime.AddDate(0, i, 0),
			License:     "MIT",
		}

		secondBuild := basePackage

		secondBuild.Version = fmt.Sprintf("%d.1", i)
		secondBuild.BuildTime = basePackage.BuildTime.AddDate(0, 0, 1)

		packageList = append(packageList, basePackage, secondBuild)
	}

	return packageList
}

func createBaseDepsolveFixture() []dnfjson.PackageSpec {
	return []dnfjson.PackageSpec{
		{
			Name:    "dep-package3",
			Epoch:   7,
			Version: "3.0.3",
			Release: "1.fc30",
			Arch:    "x86_64",
			RepoID:  "REPOID", // added by mock-dnf-json
		},
		{
			Name:    "dep-package1",
			Epoch:   0,
			Version: "1.33",
			Release: "2.fc30",
			Arch:    "x86_64",
			RepoID:  "REPOID", // added by mock-dnf-json
		},
		{
			Name:    "dep-package2",
			Epoch:   0,
			Version: "2.9",
			Release: "1.fc30",
			Arch:    "x86_64",
			RepoID:  "REPOID", // added by mock-dnf-json
		},
	}
}

// BaseDeps is the expected list of dependencies (as rpmmd.PackageSpec) from
// the Base ResponseGenerator
func BaseDeps() []rpmmd.PackageSpec {
	return []rpmmd.PackageSpec{
		{
			Name:     "dep-package3",
			Epoch:    7,
			Version:  "3.0.3",
			Release:  "1.fc30",
			Arch:     "x86_64",
			CheckGPG: true,
		},
		{
			Name:     "dep-package1",
			Epoch:    0,
			Version:  "1.33",
			Release:  "2.fc30",
			Arch:     "x86_64",
			CheckGPG: true,
		},
		{
			Name:     "dep-package2",
			Epoch:    0,
			Version:  "2.9",
			Release:  "1.fc30",
			Arch:     "x86_64",
			CheckGPG: true,
		},
	}
}

type ResponseGenerator func(string) string

func Base(tmpdir string) string {
	deps := map[string]interface{}{
		"checksums": map[string]string{
			"REPOID": "test:responsechecksum",
		},
		"dependencies": createBaseDepsolveFixture(),
	}

	pkgs := map[string]interface{}{
		"checksums": map[string]string{
			"REPOID": "test:responsechecksum",
		},
		"packages": generatePackageList(),
	}

	data := map[string]interface{}{
		"depsolve": deps,
		"dump":     pkgs,
	}
	path := filepath.Join(tmpdir, "base.json")
	write(data, path)
	return path
}

func NonExistingPackage(tmpdir string) string {
	deps := dnfjson.Error{
		Kind:   "MarkingErrors",
		Reason: "Error occurred when marking packages for installation: Problems in request:\nmissing packages: fash",
	}
	data := map[string]interface{}{
		"depsolve": deps,
	}
	path := filepath.Join(tmpdir, "notexist.json")
	write(data, path)
	return path
}

func BadDepsolve(tmpdir string) string {
	deps := dnfjson.Error{
		Kind:   "DepsolveError",
		Reason: "There was a problem depsolving ['go2rpm']: \n Problem: conflicting requests\n  - nothing provides askalono-cli needed by go2rpm-1-4.fc31.noarch",
	}
	pkgs := map[string]interface{}{
		"checksums": map[string]string{
			"REPOID": "test:responsechecksum",
		},
		"packages": generatePackageList(),
	}

	data := map[string]interface{}{
		"depsolve": deps,
		"dump":     pkgs,
	}
	path := filepath.Join(tmpdir, "baddepsolve.json")
	write(data, path)
	return path
}

func BadFetch(tmpdir string) string {
	deps := dnfjson.Error{
		Kind:   "DepsolveError",
		Reason: "There was a problem depsolving ['go2rpm']: \n Problem: conflicting requests\n  - nothing provides askalono-cli needed by go2rpm-1-4.fc31.noarch",
	}
	pkgs := dnfjson.Error{
		Kind:   "FetchError",
		Reason: "There was a problem when fetching packages.",
	}
	data := map[string]interface{}{
		"depsolve": deps,
		"dump":     pkgs,
	}
	path := filepath.Join(tmpdir, "badfetch.json")
	write(data, path)
	return path
}

func marshal(data interface{}) []byte {
	jdata, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return jdata
}

func write(data interface{}, path string) {
	fp, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	if _, err := fp.Write(marshal(data)); err != nil {
		panic(err)
	}
}
