package module

import (
	"os/exec"
	"strings"
)

// CLIExecutor abstracts command execution for testability.
type CLIExecutor interface {
	Run(cmd string, args ...string) (stdout string, stderr string, exitCode int, err error)
}

// RealExecutor runs commands using os/exec.
type RealExecutor struct{}

// NewRealExecutor returns a production CLIExecutor.
func NewRealExecutor() CLIExecutor {
	return &RealExecutor{}
}

// Run executes a command and returns stdout, stderr, exit code, and any error.
func (r *RealExecutor) Run(cmd string, args ...string) (stdout string, stderr string, exitCode int, err error) {
	c := exec.Command(cmd, args...)
	outBuf := new(strings.Builder)
	errBuf := new(strings.Builder)
	c.Stdout = outBuf
	c.Stderr = errBuf

	runErr := c.Run()
	if runErr != nil {
		if exitError, ok := runErr.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = -1
		}
		return outBuf.String(), errBuf.String(), exitCode, runErr
	}
	return outBuf.String(), errBuf.String(), 0, nil
}

// MockResponse defines the canned response for a mock command.
type MockResponse struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Err      error
}

// MockExecutor is a test double that returns fixed responses.
type MockExecutor struct {
	Responses       map[string]MockResponse
	DefaultResponse *MockResponse
}

// Run looks up the command by key "cmd arg1 arg2 ..." and returns the canned response.
func (m *MockExecutor) Run(cmd string, args ...string) (stdout string, stderr string, exitCode int, err error) {
	key := cmd
	if len(args) > 0 {
		key = cmd + " " + strings.Join(args, " ")
	}

	if resp, ok := m.Responses[key]; ok {
		return resp.Stdout, resp.Stderr, resp.ExitCode, resp.Err
	}

	if m.DefaultResponse != nil {
		return m.DefaultResponse.Stdout, m.DefaultResponse.Stderr, m.DefaultResponse.ExitCode, m.DefaultResponse.Err
	}

	return "", "", 0, nil
}
