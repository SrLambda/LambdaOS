package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListModules(t *testing.T) {
	tmp := t.TempDir()
	for _, d := range []string{"nvim", "tmux", "git"} {
		if err := os.MkdirAll(filepath.Join(tmp, d), 0755); err != nil {
			t.Fatal(err)
		}
	}

	modules, err := ListModules(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(modules) != 3 {
		t.Fatalf("expected 3 modules, got %d", len(modules))
	}

	names := make(map[string]bool)
	for _, m := range modules {
		names[m.Name] = true
	}
	for _, expected := range []string{"nvim", "tmux", "git"} {
		if !names[expected] {
			t.Errorf("expected module %q not found", expected)
		}
	}
}

func TestListModulesEmpty(t *testing.T) {
	tmp := t.TempDir()

	modules, err := ListModules(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(modules) != 0 {
		t.Errorf("expected 0 modules, got %d", len(modules))
	}
}

func TestDetectConflictsNoConflicts(t *testing.T) {
	tmp := t.TempDir()
	dotfilesDir := filepath.Join(tmp, "dotfiles")
	homeDir := filepath.Join(tmp, "home")

	// Create matching files with same content
	moduleDir := filepath.Join(dotfilesDir, "nvim", ".config", "nvim")
	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		t.Fatal(err)
	}
	repoFile := filepath.Join(moduleDir, "init.lua")
	if err := os.WriteFile(repoFile, []byte("same content"), 0644); err != nil {
		t.Fatal(err)
	}

	homeTarget := filepath.Join(homeDir, ".config", "nvim")
	if err := os.MkdirAll(homeTarget, 0755); err != nil {
		t.Fatal(err)
	}
	homeFile := filepath.Join(homeTarget, "init.lua")
	if err := os.WriteFile(homeFile, []byte("same content"), 0644); err != nil {
		t.Fatal(err)
	}

	conflicts, err := DetectConflicts(dotfilesDir, "nvim", homeDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(conflicts) != 0 {
		t.Errorf("expected 0 conflicts, got %d", len(conflicts))
	}
}

func TestDetectConflictsWithMismatch(t *testing.T) {
	tmp := t.TempDir()
	dotfilesDir := filepath.Join(tmp, "dotfiles")
	homeDir := filepath.Join(tmp, "home")

	// Create repo file
	moduleDir := filepath.Join(dotfilesDir, "nvim", ".config", "nvim")
	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		t.Fatal(err)
	}
	repoFile := filepath.Join(moduleDir, "init.lua")
	if err := os.WriteFile(repoFile, []byte("repo content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create home file with different content
	homeTarget := filepath.Join(homeDir, ".config", "nvim")
	if err := os.MkdirAll(homeTarget, 0755); err != nil {
		t.Fatal(err)
	}
	homeFile := filepath.Join(homeTarget, "init.lua")
	if err := os.WriteFile(homeFile, []byte("home content"), 0644); err != nil {
		t.Fatal(err)
	}

	conflicts, err := DetectConflicts(dotfilesDir, "nvim", homeDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(conflicts))
	}

	if conflicts[0].Path != homeFile {
		t.Errorf("expected conflict path %q, got %q", homeFile, conflicts[0].Path)
	}

	if conflicts[0].RepoChecksum == conflicts[0].HomeChecksum {
		t.Error("expected different checksums for conflicting files")
	}
}

func TestDetectConflictsHomeOnly(t *testing.T) {
	tmp := t.TempDir()
	dotfilesDir := filepath.Join(tmp, "dotfiles")
	homeDir := filepath.Join(tmp, "home")

	// Create module dir but no files
	if err := os.MkdirAll(filepath.Join(dotfilesDir, "nvim"), 0755); err != nil {
		t.Fatal(err)
	}

	// Create file only in home
	homeTarget := filepath.Join(homeDir, ".config", "nvim")
	if err := os.MkdirAll(homeTarget, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(homeTarget, "init.lua"), []byte("home only"), 0644); err != nil {
		t.Fatal(err)
	}

	conflicts, err := DetectConflicts(dotfilesDir, "nvim", homeDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(conflicts) != 0 {
		t.Errorf("expected 0 conflicts (home-only file), got %d", len(conflicts))
	}
}

func TestDetectConflictsRepoOnly(t *testing.T) {
	tmp := t.TempDir()
	dotfilesDir := filepath.Join(tmp, "dotfiles")
	homeDir := filepath.Join(tmp, "home")

	// Create repo file
	moduleDir := filepath.Join(dotfilesDir, "nvim", ".config", "nvim")
	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "init.lua"), []byte("repo only"), 0644); err != nil {
		t.Fatal(err)
	}

	// No corresponding home file

	conflicts, err := DetectConflicts(dotfilesDir, "nvim", homeDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(conflicts) != 0 {
		t.Errorf("expected 0 conflicts (repo-only file), got %d", len(conflicts))
	}
}

func TestBackupCopiesChangedFiles(t *testing.T) {
	tmp := t.TempDir()
	dotfilesDir := filepath.Join(tmp, "dotfiles")
	homeDir := filepath.Join(tmp, "home")

	// Create home file with content
	homeTarget := filepath.Join(homeDir, ".config", "nvim")
	if err := os.MkdirAll(homeTarget, 0755); err != nil {
		t.Fatal(err)
	}
	homeContent := []byte("home config content")
	if err := os.WriteFile(filepath.Join(homeTarget, "init.lua"), homeContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Create empty repo module dir
	repoModuleDir := filepath.Join(dotfilesDir, "nvim", ".config", "nvim")
	if err := os.MkdirAll(repoModuleDir, 0755); err != nil {
		t.Fatal(err)
	}

	count, err := BackupModule(dotfilesDir, "nvim", homeDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 file backed up, got %d", count)
	}

	// Verify the backed up file matches home
	repoFile := filepath.Join(repoModuleDir, "init.lua")
	got, err := os.ReadFile(repoFile)
	if err != nil {
		t.Fatalf("failed to read backed up file: %v", err)
	}

	if string(got) != string(homeContent) {
		t.Errorf("backed up content mismatch: got %q, want %q", string(got), string(homeContent))
	}
}

func TestBackupSkipsIdenticalFiles(t *testing.T) {
	tmp := t.TempDir()
	dotfilesDir := filepath.Join(tmp, "dotfiles")
	homeDir := filepath.Join(tmp, "home")

	content := []byte("identical content")

	// Create home file
	homeTarget := filepath.Join(homeDir, ".config", "nvim")
	if err := os.MkdirAll(homeTarget, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(homeTarget, "init.lua"), content, 0644); err != nil {
		t.Fatal(err)
	}

	// Create matching repo file
	repoModuleDir := filepath.Join(dotfilesDir, "nvim", ".config", "nvim")
	if err := os.MkdirAll(repoModuleDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repoModuleDir, "init.lua"), content, 0644); err != nil {
		t.Fatal(err)
	}

	count, err := BackupModule(dotfilesDir, "nvim", homeDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if count != 0 {
		t.Errorf("expected 0 files backed up (identical), got %d", count)
	}
}
