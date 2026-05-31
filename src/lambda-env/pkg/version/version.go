package version

import (
	"fmt"
	"strconv"
	"strings"
)

// Version is the current hub version.
const Version = "0.1.0"

// CheckCompatibility returns true if the hub version satisfies the minimum
// version requirement (hubVersion >= minVersion).
func CheckCompatibility(minVersion string) bool {
	cmp, err := compareVersions(Version, minVersion)
	if err != nil {
		return false
	}
	return cmp >= 0
}

// compareVersions compares two simple semver strings.
// Returns -1 if a < b, 0 if a == b, 1 if a > b.
func compareVersions(a, b string) (int, error) {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	for i := 0; i < 3; i++ {
		var av, bv int
		if i < len(aParts) {
			v, err := strconv.Atoi(aParts[i])
			if err != nil {
				return 0, fmt.Errorf("invalid version component %q in %q: %w", aParts[i], a, err)
			}
			av = v
		}
		if i < len(bParts) {
			v, err := strconv.Atoi(bParts[i])
			if err != nil {
				return 0, fmt.Errorf("invalid version component %q in %q: %w", bParts[i], b, err)
			}
			bv = v
		}
		if av < bv {
			return -1, nil
		}
		if av > bv {
			return 1, nil
		}
	}
	return 0, nil
}
