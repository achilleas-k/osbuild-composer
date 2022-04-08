// Mock dnf-json
//
// The purpose of this program is to return fake but expected responses to
// dnf-json depsolve and dump queries.  Tests should initialise a
// dnfjson.Solver and configure it to run this program via the SetDNFJSONPath()
// method.  This utility accepts queries and returns responses with the same
// structure as the dnf-json Python script.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/osbuild/osbuild-composer/internal/dnfjson"
)

func maybeFail(err error) {
	if err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func readRequest(data []byte) dnfjson.Request {
	var req dnfjson.Request
	maybeFail(json.Unmarshal(data, &req))
	return req
}

func respond(result []map[string]interface{}) {
	resp, err := json.Marshal(result)
	maybeFail(err)
	fmt.Printf(string(resp))
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

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	maybeFail(err)

	req := readRequest(input)

	args := req.Arguments
	if len(args) == 0 {
		fail(errors.New("error: empty arguments"))
	}

	pkgs := createBaseDepsolveFixture()
	response := []map[string]interface{}{{
		"checksums": map[string]string{
			"0": "test:responsechecksum",
		},
		"dependencies": pkgs,
	}}
	respond(response)
}
