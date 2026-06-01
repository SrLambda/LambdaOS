package hub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
	"lambdaos.dev/lambda-env/pkg/version"
)

// ExecuteModule runs the module's default "run" action.
// It is backward-compatible and delegates to ExecuteAction with action="run".
func (h *Hub) ExecuteModule(mod module.Manifest) (*module.Response, error) {
	return h.ExecuteAction(mod, "run", nil)
}

// ExecuteAction runs a specific module action with optional parameters.
// It sets LAMBDA_ENV_ACTION to the action name and passes params as JSON
// in the LAMBDA_ENV_PARAMS environment variable.
func (h *Hub) ExecuteAction(mod module.Manifest, actionName string, params map[string]interface{}) (*module.Response, error) {
	binPath := filepath.Join(mod.Path, "module")
	if _, err := os.Stat(binPath); err != nil {
		return nil, fmt.Errorf("module executable not found at %s: %w", binPath, err)
	}

	locale := os.Getenv("LANG")
	if locale == "" {
		locale = "en_US"
	}

	env := map[string]string{
		"LAMBDA_ENV_ACTION":      actionName,
		"LAMBDA_ENV_SETTINGS":    h.StorePath,
		"LAMBDA_ENV_HUB_VERSION": version.Version,
		"LAMBDA_ENV_LOCALE":      locale,
	}

	if params != nil && len(params) > 0 {
		paramsJSON, err := json.Marshal(params)
		if err == nil {
			env["LAMBDA_ENV_PARAMS"] = string(paramsJSON)
		}
	}

	timeout := mod.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath)
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	runErr := cmd.Run()
	exitCode := 0
	if runErr != nil {
		if exitErr, ok := runErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else if ctx.Err() == context.DeadlineExceeded {
			exitCode = 124 // standard timeout exit code
			return nil, fmt.Errorf("module timed out after %d seconds", timeout)
		} else {
			exitCode = 1
		}
	}

	stdoutStr := stdoutBuf.String()
	stderrStr := stderrBuf.String()

	// Log regardless of outcome.
	if h.Logger != nil {
		_ = h.Logger.Log(mod.Name, actionName, exitCode, stdoutStr, stderrStr, env)
	}

	// Parse JSON response from stdout.
	var resp module.Response
	if err := json.Unmarshal(stdoutBuf.Bytes(), &resp); err != nil {
		return nil, fmt.Errorf("failed to parse module JSON output: %w\nraw stdout: %s", err, stdoutStr)
	}

	// Merge settings delta if present and execution succeeded.
	if resp.Status == "ok" && len(resp.SettingsDelta) > 0 {
		if err := settings.SaveDelta(h.StorePath, resp.SettingsDelta); err != nil {
			return &resp, fmt.Errorf("settings delta merge failed: %w", err)
		}
	}

	// Treat non-zero exit codes that still emitted JSON as errors/warnings.
	if exitCode != 0 && resp.Status == "ok" {
		resp.Status = "error"
		if resp.Message == "" {
			resp.Message = fmt.Sprintf("module exited with code %d", exitCode)
		}
	}

	if resp.Status == "error" {
		return &resp, fmt.Errorf("module error: %s", resp.Message)
	}

	return &resp, nil
}
