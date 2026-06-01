package main

import (
	"io"
	"os"
	"path/filepath"
)

// BackupModule copies files from homeDir config dirs to the dotfiles repo,
// preserving directory structure. Returns count of files actually copied.
//
// For each file in <homeDir>/.config/<module>/, copies to
// <dotfilesDir>/<name>/.config/<module>/.
// Skips files that already match (compare checksums before copy).
func BackupModule(dotfilesDir, name, homeDir string) (int, error) {
	homeModuleDir := filepath.Join(homeDir, ".config", name)
	repoModuleDir := filepath.Join(dotfilesDir, name, ".config", name)

	info, err := os.Stat(homeModuleDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	if !info.IsDir() {
		return 0, nil
	}

	if err := os.MkdirAll(repoModuleDir, 0755); err != nil {
		return 0, err
	}

	count := 0

	err = filepath.Walk(homeModuleDir, func(homePath string, homeInfo os.FileInfo, err error) error {
		if err != nil || homeInfo.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(homeModuleDir, homePath)
		if err != nil {
			return nil
		}

		repoPath := filepath.Join(repoModuleDir, rel)

		// Check if repo file exists and matches
		if _, err := os.Stat(repoPath); err == nil {
			repoChecksum, err := SHA256File(repoPath)
			if err != nil {
				return nil
			}

			homeChecksum, err := SHA256File(homePath)
			if err != nil {
				return nil
			}

			if repoChecksum == homeChecksum {
				return nil
			}
		}

		// Create target directory
		targetDir := filepath.Dir(repoPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return err
		}

		// Copy home file to repo
		if err := copyFile(homePath, repoPath); err != nil {
			return err
		}

		count++
		return nil
	})

	return count, err
}

// copyFile copies src to dst using io.Copy.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
