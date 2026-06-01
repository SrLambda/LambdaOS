package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

// Conflict represents a file that exists in both the repo and home with different content.
type Conflict struct {
	Path        string `json:"path"`
	RepoChecksum  string `json:"repo_checksum"`
	HomeChecksum  string `json:"home_checksum"`
}

// SHA256File reads a file and returns its hex-encoded SHA-256 checksum.
func SHA256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// DetectConflicts walks the module tree in dotfilesDir/<name>/, and for each file
// checks if the corresponding home file exists. Returns conflicts where both exist
// with different checksums.
//
// The target path for a file at ~/dotfiles/nvim/.config/nvim/init.lua is
// <homeDir>/.config/nvim/init.lua.
func DetectConflicts(dotfilesDir, name, homeDir string) ([]Conflict, error) {
	moduleDir := filepath.Join(dotfilesDir, name)

	info, err := os.Stat(moduleDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, nil
	}

	var conflicts []Conflict

	err = filepath.Walk(moduleDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(moduleDir, path)
		if err != nil {
			return nil
		}

		homePath := filepath.Join(homeDir, rel)

		homeInfo, err := os.Lstat(homePath)
		if err != nil {
			// Home file doesn't exist — not a conflict
			return nil
		}

		// Skip symlinks (they're managed by stow)
		if homeInfo.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		repoChecksum, err := SHA256File(path)
		if err != nil {
			return nil
		}

		homeChecksum, err := SHA256File(homePath)
		if err != nil {
			return nil
		}

		if repoChecksum != homeChecksum {
			conflicts = append(conflicts, Conflict{
				Path:         homePath,
				RepoChecksum: repoChecksum,
				HomeChecksum: homeChecksum,
			})
		}

		return nil
	})

	return conflicts, err
}
