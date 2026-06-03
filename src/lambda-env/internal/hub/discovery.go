package hub

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"lambdaos.dev/lambda-env/pkg/module"
)

const (
	// SystemModulesPath is the system-wide module directory.
	SystemModulesPath = "/usr/share/lambda-env/modules"
	// UserModulesPath is the user-specific module directory.
	UserModulesPath = "~/.local/share/lambda-env/modules"
)

// Scan discovers modules from system and user paths.
// User modules override system modules with the same name.
// Results are sorted by category, then name.
func Scan() ([]module.Manifest, error) {
	system, err := scanPath(SystemModulesPath)
	if err != nil {
		// Non-fatal: system path may not exist.
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	userPath := expandHome(UserModulesPath)
	user, err := scanPath(userPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	// Merge: user overrides system.
	merged := make(map[string]module.Manifest)
	for _, m := range system {
		merged[m.Name] = m
	}
	for _, m := range user {
		merged[m.Name] = m
	}

	result := make([]module.Manifest, 0, len(merged))
	for _, m := range merged {
		result = append(result, m)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Category != result[j].Category {
			return result[i].Category < result[j].Category
		}
		return result[i].Name < result[j].Name
	})

	return result, nil
}

func scanPath(base string) ([]module.Manifest, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}

	var manifests []module.Manifest
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(base, entry.Name(), "manifest.json")
		data, err := os.ReadFile(manifestPath)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "warning: cannot read %s: %v\n", manifestPath, err)
			continue
		}

		var m module.Manifest
		if err := json.Unmarshal(data, &m); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "warning: invalid manifest %s: %v\n", manifestPath, err)
			continue
		}

		if err := m.Validate(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "warning: invalid manifest %s: %v\n", manifestPath, err)
			continue
		}

		m.Path = filepath.Join(base, entry.Name())
		manifests = append(manifests, m)
	}

	return manifests, nil
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
