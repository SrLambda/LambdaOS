package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Load reads settings from path, applying defaults and migration.
// If the file does not exist, it returns defaults with no error.
func Load(path string) (*Settings, error) {
	s := Defaults()

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &s, nil
	}
	if err != nil {
		return nil, err
	}

	var loaded map[string]interface{}
	if err := json.Unmarshal(data, &loaded); err != nil {
		return nil, err
	}

	// Merge loaded values on top of defaults.
	defaultsJSON, _ := json.Marshal(s)
	var defaultsMap map[string]interface{}
	json.Unmarshal(defaultsJSON, &defaultsMap)

	merged := deepMerge(defaultsMap, loaded)

	mergedJSON, _ := json.Marshal(merged)
	var result Settings
	if err := json.Unmarshal(mergedJSON, &result); err != nil {
		return nil, err
	}

	if err := Migrate(&result); err != nil {
		return nil, err
	}

	if err := result.Validate(); err != nil {
		return nil, err
	}

	return &result, nil
}

// Save writes settings to path atomically using a temp file + rename.
func Save(path string, s *Settings) error {
	if err := s.Validate(); err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, "settings-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()

	var writeErr error
	defer func() {
		if writeErr != nil {
			os.Remove(tmpPath)
		}
	}()

	enc := json.NewEncoder(tmp)
	enc.SetIndent("", "  ")
	if writeErr = enc.Encode(s); writeErr != nil {
		tmp.Close()
		return writeErr
	}

	if writeErr = tmp.Close(); writeErr != nil {
		return writeErr
	}

	return os.Rename(tmpPath, path)
}

// SaveDelta loads current settings, deep-merges delta, and saves atomically.
func SaveDelta(path string, delta map[string]interface{}) error {
	s, err := Load(path)
	if err != nil {
		return err
	}

	currentJSON, _ := json.Marshal(s)
	var currentMap map[string]interface{}
	json.Unmarshal(currentJSON, &currentMap)

	merged := deepMerge(currentMap, delta)

	mergedJSON, _ := json.Marshal(merged)
	var result Settings
	if err := json.Unmarshal(mergedJSON, &result); err != nil {
		return err
	}

	if err := result.Validate(); err != nil {
		return err
	}

	return Save(path, &result)
}

// Migrate checks the version in s and applies migrations if needed.
// It rejects downgrades.
func Migrate(s *Settings) error {
	if s.Version == "" {
		s.Version = "0.0.0"
	}

	cmp, err := compareVersions(s.Version, CurrentVersion)
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", s.Version, err)
	}

	if cmp > 0 {
		return fmt.Errorf("cannot downgrade: file version %s > current version %s", s.Version, CurrentVersion)
	}

	if cmp < 0 {
		// Fill any missing fields with defaults.
		defaults := Defaults()

		currentJSON, _ := json.Marshal(s)
		var currentMap map[string]interface{}
		json.Unmarshal(currentJSON, &currentMap)

		defaultsJSON, _ := json.Marshal(defaults)
		var defaultsMap map[string]interface{}
		json.Unmarshal(defaultsJSON, &defaultsMap)

		merged := deepMerge(defaultsMap, currentMap)
		mergedJSON, _ := json.Marshal(merged)

		if err := json.Unmarshal(mergedJSON, s); err != nil {
			return err
		}

		s.Version = CurrentVersion
	}

	return nil
}

// deepMerge recursively merges override into base. override values win.
func deepMerge(base, override map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range base {
		result[k] = v
	}
	for k, v := range override {
		if vm, ok := v.(map[string]interface{}); ok {
			if bm, ok := result[k].(map[string]interface{}); ok {
				result[k] = deepMerge(bm, vm)
			} else {
				result[k] = vm
			}
		} else {
			result[k] = v
		}
	}
	return result
}

// compareVersions compares two semver strings.
// Returns -1 if a < b, 0 if a == b, 1 if a > b.
func compareVersions(a, b string) (int, error) {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	for i := 0; i < 3; i++ {
		var av, bv int
		if i < len(aParts) {
			v, err := strconv.Atoi(aParts[i])
			if err != nil {
				return 0, err
			}
			av = v
		}
		if i < len(bParts) {
			v, err := strconv.Atoi(bParts[i])
			if err != nil {
				return 0, err
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
