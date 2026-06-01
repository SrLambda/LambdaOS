package module

import (
	"errors"
	"strings"
	"testing"
)

func TestRealExecutorEcho(t *testing.T) {
	exec := NewRealExecutor()
	stdout, stderr, exitCode, err := exec.Run("echo", "hello", "world")
	if err != nil {
		t.Fatalf("echo failed: %v", err)
	}
	if exitCode != 0 {
		t.Errorf("exitCode = %d, want 0", exitCode)
	}
	if !strings.Contains(stdout, "hello world") {
		t.Errorf("stdout = %q, want to contain 'hello world'", stdout)
	}
	if stderr != "" {
		t.Errorf("stderr = %q, want empty", stderr)
	}
}

func TestRealExecutorExitCode(t *testing.T) {
	exec := NewRealExecutor()
	// sh -c 'exit 42' returns exit code 42.
	stdout, stderr, exitCode, err := exec.Run("sh", "-c", "exit 42")
	if err == nil {
		t.Fatal("expected error for non-zero exit code, got nil")
	}
	if exitCode != 42 {
		t.Errorf("exitCode = %d, want 42", exitCode)
	}
	if stdout != "" {
		t.Errorf("stdout = %q, want empty", stdout)
	}
	if stderr != "" {
		t.Errorf("stderr = %q, want empty", stderr)
	}
}

func TestRealExecutorStderr(t *testing.T) {
	exec := NewRealExecutor()
	stdout, stderr, exitCode, err := exec.Run("sh", "-c", "echo err >&2; exit 1")
	if err == nil {
		t.Fatal("expected error for non-zero exit code, got nil")
	}
	if exitCode != 1 {
		t.Errorf("exitCode = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr, "err") {
		t.Errorf("stderr = %q, want to contain 'err'", stderr)
	}
	if stdout != "" {
		t.Errorf("stdout = %q, want empty", stdout)
	}
}

func TestMockExecutor(t *testing.T) {
	mock := &MockExecutor{
		Responses: map[string]MockResponse{
			"echo hello": {
				Stdout:   "hello",
				Stderr:   "",
				ExitCode: 0,
				Err:      nil,
			},
			"false": {
				Stdout:   "",
				Stderr:   "error",
				ExitCode: 1,
				Err:      errors.New("command failed"),
			},
		},
	}

	stdout, stderr, exitCode, err := mock.Run("echo", "hello")
	if err != nil {
		t.Fatalf("mock echo hello: unexpected error: %v", err)
	}
	if stdout != "hello" {
		t.Errorf("stdout = %q, want %q", stdout, "hello")
	}
	if stderr != "" {
		t.Errorf("stderr = %q, want empty", stderr)
	}
	if exitCode != 0 {
		t.Errorf("exitCode = %d, want 0", exitCode)
	}

	stdout, stderr, exitCode, err = mock.Run("false")
	if err == nil {
		t.Fatal("mock false: expected error, got nil")
	}
	if exitCode != 1 {
		t.Errorf("exitCode = %d, want 1", exitCode)
	}
	if stderr != "error" {
		t.Errorf("stderr = %q, want %q", stderr, "error")
	}
	if stdout != "" {
		t.Errorf("stdout = %q, want empty", stdout)
	}
}

func TestMockExecutorFallback(t *testing.T) {
	mock := &MockExecutor{
		DefaultResponse: &MockResponse{
			Stdout:   "default",
			Stderr:   "",
			ExitCode: 0,
			Err:      nil,
		},
	}

	stdout, _, _, err := mock.Run("unknown", "cmd")
	if err != nil {
		t.Fatalf("mock unknown: unexpected error: %v", err)
	}
	if stdout != "default" {
		t.Errorf("stdout = %q, want %q", stdout, "default")
	}
}
