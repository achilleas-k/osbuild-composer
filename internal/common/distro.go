package common

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

func GetHostDistroName() (string, bool, bool, error) {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return "", false, false, err
	}
	defer f.Close()
	osrelease, err := readOSRelease(f)
	if err != nil {
		return "", false, false, err
	}

	isStream := osrelease["NAME"] == "CentOS Stream"

	version := strings.Split(osrelease["VERSION_ID"], ".")
	name := osrelease["ID"] + "-" + strings.Join(version, "")

	// TODO: We should probably index these things by the full CPE
	beta := strings.Contains(osrelease["CPE_NAME"], "beta")
	return name, beta, isStream, nil
}

func readOSRelease(r io.Reader) (map[string]string, error) {
	osrelease := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("readOSRelease: invalid input")
		}

		key := strings.TrimSpace(parts[0])
		// drop all surrounding whitespace and double-quotes
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"")
		osrelease[key] = value
	}

	return osrelease, nil
}

// Returns true if the version represented by the first argument is
// semantically older than the second.
// Meant to be used for comparing distro versions for differences between minor
// releases.
// Evaluates to false if a and b are equal.
// Assumes any missing components are 0, so 8 < 8.1.
func VersionLessThan(a, b string) bool {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	// pad shortest argument with zeroes
	for len(aParts) < len(bParts) {
		aParts = append(aParts, "0")
	}
	for len(bParts) < len(aParts) {
		bParts = append(bParts, "0")
	}

	for idx := 0; idx < len(aParts); idx++ {
		if aParts[idx] < bParts[idx] {
			return true
		} else if aParts[idx] > bParts[idx] {
			return false
		}
	}

	// equal
	return false
}
