package main

import (
	"os"
	"path/filepath"
)

// StowModule represents a dotfiles module discovered in the repo.
type StowModule struct {
	Name   string `json:"name"`
	Stowed bool   `json:"stowed"`
}

// ListModules scans dotfilesDir for subdirectories and returns a list of modules
// with their stow status.
func ListModules(dotfilesDir string) ([]StowModule, error) {
	entries, err := os.ReadDir(dotfilesDir)
	if err != nil {
		return nil, err
	}

	var modules []StowModule
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if name == "." || name == ".." {
			continue
		}
		modules = append(modules, StowModule{
			Name:   name,
			Stowed: IsStowed(dotfilesDir, name),
		})
	}

	return modules, nil
}

// Stow runs `stow <name>` from dotfilesDir.
func Stow(dotfilesDir, name string) error {
	return runStow(dotfilesDir, name)
}

// Unstow runs `stow -D <name>` from dotfilesDir.
func Unstow(dotfilesDir, name string) error {
	return runStow(dotfilesDir, "-D", name)
}

// IsStowed checks if any symlinks exist in target directories for this module.
// It walks the module tree and checks if the corresponding home paths are symlinks.
func IsStowed(dotfilesDir, name string) bool {
	moduleDir := filepath.Join(dotfilesDir, name)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	found := false
	filepath.Walk(moduleDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(moduleDir, path)
		if err != nil {
			return nil
		}

		targetPath := filepath.Join(homeDir, rel)
		if linfo, err := os.Lstat(targetPath); err == nil {
			if linfo.Mode()&os.ModeSymlink != 0 {
				found = true
			}
		}

		return nil
	})

	return found
}
