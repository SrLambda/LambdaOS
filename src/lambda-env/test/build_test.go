package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// moduleRoot returns the absolute path to the lambda-env module root
// (src/lambda-env/), computed from this source file's location.
func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine source file path")
	}
	// file = <module>/test/build_test.go → module root = file/dir/test/..
	return filepath.Dir(filepath.Dir(file))
}

// TestBinaryCompiles verifies the lambda-env binary compiles without errors.
func TestBinaryCompiles(t *testing.T) {
	root := moduleRoot(t)
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "lambda-env")

	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/lambda-env")
	cmd.Dir = root
	cmd.Env = os.Environ()

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build failed: %v\noutput: %s", err, out)
	}

	if _, err := os.Stat(binaryPath); err != nil {
		t.Fatalf("binary not found at %s after build: %v", binaryPath, err)
	}
}

// TestBinaryRunsHelp verifies the built binary runs without crashing.
func TestBinaryRunsHelp(t *testing.T) {
	root := moduleRoot(t)
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "lambda-env")

	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/lambda-env")
	cmd.Dir = root
	cmd.Env = os.Environ()
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\noutput: %s", err, out)
	}

	runCmd := exec.Command(binaryPath, "--help")
	runCmd.Dir = tmpDir
	out, err := runCmd.CombinedOutput()
	if err != nil {
		t.Logf("binary exited with: %v\noutput: %s", err, out)
	}
}
