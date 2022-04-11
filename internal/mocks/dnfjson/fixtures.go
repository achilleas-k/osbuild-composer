// dnfjson_mock provides data and methods for testing the dnfjson package.
package dnfjson_mock

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/osbuild/osbuild-composer/internal/dnfjson"
)

func generatePackageList() []dnfjson.PackageSpec {
	var packageList []dnfjson.PackageSpec

	for i := 0; i < 22; i++ {
		basePackage := dnfjson.PackageSpec{
			Name:           fmt.Sprintf("package%d", i),
			Version:        fmt.Sprintf("%d.0", i),
			Release:        fmt.Sprintf("%d.fc30", i),
			Arch:           "x86_64",
			RepoID:         "0",
			RemoteLocation: fmt.Sprintf("https://pkg%d.example.com", i),
			Checksum:       fmt.Sprintf("notachecksum-%d", i),
		}

		secondBuild := basePackage

		secondBuild.Version = fmt.Sprintf("%d.1", i)

		packageList = append(packageList, basePackage, secondBuild)
	}

	sort.Slice(packageList, func(i, j int) bool {
		return packageList[i].Name < packageList[j].Name
	})

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
			RepoID:  "0",
		},
		{
			Name:    "dep-package1",
			Epoch:   0,
			Version: "1.33",
			Release: "2.fc30",
			Arch:    "x86_64",
			RepoID:  "0",
		},
		{
			Name:    "dep-package2",
			Epoch:   0,
			Version: "2.9",
			Release: "1.fc30",
			Arch:    "x86_64",
			RepoID:  "0",
		},
	}
}

func BaseDepsolve(tmpdir string) string {
	pkgs := createBaseDepsolveFixture()
	data := []map[string]interface{}{{
		"checksums": map[string]string{
			"0": "test:responsechecksum",
		},
		"dependencies": pkgs,
	}}
	path := filepath.Join(tmpdir, "base.json")
	write(data, path)
	return path
}

func PackageList(tmpdir string) string {
	pkgs := generatePackageList()
	data := map[string]interface{}{
		"checksums": map[string]string{
			"0": "test:responsechecksum",
		},
		"dependencies": pkgs,
	}
	path := filepath.Join(tmpdir, "pkgs.json")
	write(data, path)
	return path
}

func NonExistingPackage(tmpdir string) string {
	data := dnfjson.Error{
		Kind:   "MarkingErrors",
		Reason: "Error occurred when marking packages for installation: Problems in request:\nmissing packages: fash",
	}
	path := filepath.Join(tmpdir, "notexist.err.json")
	write(data, path)
	return path
}

func BadDepsolve(tmpdir string) string {
	data := dnfjson.Error{
		Kind:   "DepsolveError",
		Reason: "There was a problem depsolving ['go2rpm']: \n Problem: conflicting requests\n  - nothing provides askalono-cli needed by go2rpm-1-4.fc31.noarch",
	}
	path := filepath.Join(tmpdir, "baddepsolve.err.json")
	write(data, path)
	return path
}

func BadFetch(tmpdir string) string {
	data := dnfjson.Error{
		Kind:   "FetchError",
		Reason: "There was a problem when fetching packages.",
	}
	path := filepath.Join(tmpdir, "badfetch.err.json")
	write(data, path)
	return path
}

func write(data interface{}, path string) {
	fp, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	jdata, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	fp.Write(jdata)
}
