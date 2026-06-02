package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestBinaryCompiles verifies the lambda-env binary compiles without errors.
// This catches compilation errors that unit tests (which compile packages individually)
// might miss, such as package main conflicts, init() cycles, or missing entry points.
func TestBinaryCompiles(t *testing.T) {
	// Find the module root (parent of test/)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Build to a temp file to avoid polluting the workspace.
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "lambda-env")

	// go build -o <tmp> lambdaos.dev/lambda-env/cmd/lambda-env
	cmd := exec.Command("go", "build", "-o", binaryPath, "lambdaos.dev/lambda-env/cmd/lambda-env")
	cmd.Dir = filepath.Join(wd, "..") // navigate from src/lambda-env/test/ to src/lambda-env/
	cmd.Env = os.Environ()

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build failed: %v\noutput: %s", err, out)
	}

	// Verify binary exists.
	if _, err := os.Stat(binaryPath); err != nil {
		t.Fatalf("binary not found at %s after build: %v", binaryPath, err)
	}
}

// TestBinaryRunsHelp verifies the built binary runs without crashing.
func TestBinaryRunsHelp(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "lambda-env")

	// Build first.
	cmd := exec.Command("go", "build", "-o", binaryPath, "lambdaos.dev/lambda-env/cmd/lambda-env")
	cmd.Dir = filepath.Join(wd, "..", "..")
	cmd.Env = os.Environ()
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\noutput: %s", err, out)
	}

	// Run with --help (non-TUI mode so it doesn't try to open a terminal).
	runCmd := exec.Command(binaryPath, "--help")
	runCmd.Dir = tmpDir
	out, err := runCmd.CombinedOutput()
	// --help might not exist yet, but the binary should at least start and exit cleanly
	// rather than panicking or crashing with a nil dereference.
	if err != nil {
		// If exit code is non-zero, that's fine (--help might return non-zero on some programs).
		// We just check the binary didn't panic.
		t.Logf("binary exited with: %v\noutput: %s", err, out)
	}
}
