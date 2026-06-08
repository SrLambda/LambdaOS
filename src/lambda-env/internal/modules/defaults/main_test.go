package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

func captureDefaultsResponse(t *testing.T, action, params, settingsPath string) module.Response {
	t.Helper()
	t.Setenv("LAMBDA_ENV_ACTION", action)
	if params != "" {
		t.Setenv("LAMBDA_ENV_PARAMS", params)
	} else {
		os.Unsetenv("LAMBDA_ENV_PARAMS")
	}
	t.Setenv("LAMBDA_ENV_SETTINGS", settingsPath)

	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = oldStdout

	var resp module.Response
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return resp
}

func TestHandleRunReturnsCurrentDefaultsAndApps(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{}
	defer func() { executor = oldExecutor }()

	// Create fake app directory
	tmpDir := t.TempDir()
	appsDir = tmpDir
	for _, app := range []string{"firefox", "chromium", "kitty", "thunar", "nvim"} {
		if err := os.WriteFile(filepath.Join(tmpDir, app+".desktop"), []byte("[Desktop Entry]"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	settingsDir := t.TempDir()
	settingsPath := filepath.Join(settingsDir, "settings.json")
	s := settings.Defaults()
	s.Defaults.Browser = "firefox"
	s.Defaults.Terminal = "kitty"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureDefaultsResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	if resp.Data["browser"] != "firefox" {
		t.Errorf("browser = %v, want firefox", resp.Data["browser"])
	}
	if resp.Data["terminal"] != "kitty" {
		t.Errorf("terminal = %v, want kitty", resp.Data["terminal"])
	}
	apps, ok := resp.Data["available_apps"].([]interface{})
	if !ok || len(apps) != 5 {
		t.Fatalf("expected 5 apps, got %v", resp.Data["available_apps"])
	}
}

func TestSetBrowserValidSavesDelta(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	appsDir = tmpDir
	if err := os.WriteFile(filepath.Join(tmpDir, "firefox.desktop"), []byte("[Desktop Entry]"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "chromium.desktop"), []byte("[Desktop Entry]"), 0644); err != nil {
		t.Fatal(err)
	}

	settingsDir := t.TempDir()
	settingsPath := filepath.Join(settingsDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureDefaultsResponse(t, "set-browser", `{"value":"chromium"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["defaults"].(map[string]interface{})
	if !ok || delta["browser"] != "chromium" {
		t.Errorf("settings delta browser = %v, want chromium", delta["browser"])
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Defaults.Browser != "chromium" {
		t.Errorf("Defaults.Browser = %q, want chromium", loaded.Defaults.Browser)
	}
}

func TestSetBrowserInvalidDesktopFileReturnsError(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	appsDir = tmpDir
	// Only firefox.desktop exists
	if err := os.WriteFile(filepath.Join(tmpDir, "firefox.desktop"), []byte("[Desktop Entry]"), 0644); err != nil {
		t.Fatal(err)
	}

	settingsDir := t.TempDir()
	settingsPath := filepath.Join(settingsDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureDefaultsResponse(t, "set-browser", `{"value":"nonexistent"}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error, got %s", resp.Status)
	}
}

func TestSetTerminalSavesDelta(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	appsDir = tmpDir
	if err := os.WriteFile(filepath.Join(tmpDir, "kitty.desktop"), []byte("[Desktop Entry]"), 0644); err != nil {
		t.Fatal(err)
	}

	settingsDir := t.TempDir()
	settingsPath := filepath.Join(settingsDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureDefaultsResponse(t, "set-terminal", `{"value":"kitty"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["defaults"].(map[string]interface{})
	if !ok || delta["terminal"] != "kitty" {
		t.Errorf("settings delta terminal = %v, want kitty", delta["terminal"])
	}
}

func TestSetEditorSavesDelta(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	appsDir = tmpDir
	if err := os.WriteFile(filepath.Join(tmpDir, "nvim.desktop"), []byte("[Desktop Entry]"), 0644); err != nil {
		t.Fatal(err)
	}

	settingsDir := t.TempDir()
	settingsPath := filepath.Join(settingsDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureDefaultsResponse(t, "set-editor", `{"value":"nvim"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["defaults"].(map[string]interface{})
	if !ok || delta["editor"] != "nvim" {
		t.Errorf("settings delta editor = %v, want nvim", delta["editor"])
	}
}

func TestSetFileManagerSavesDelta(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	appsDir = tmpDir
	if err := os.WriteFile(filepath.Join(tmpDir, "thunar.desktop"), []byte("[Desktop Entry]"), 0644); err != nil {
		t.Fatal(err)
	}

	settingsDir := t.TempDir()
	settingsPath := filepath.Join(settingsDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureDefaultsResponse(t, "set-file-manager", `{"value":"thunar"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["defaults"].(map[string]interface{})
	if !ok || delta["file_manager"] != "thunar" {
		t.Errorf("settings delta file_manager = %v, want thunar", delta["file_manager"])
	}
}

func TestApplyConstructsXdgCommands(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"xdg-settings set default-web-browser firefox.desktop": {
				Stdout:   "",
				ExitCode: 0,
			},
			"xdg-mime default kitty.desktop x-scheme-handler/terminal": {
				Stdout:   "",
				ExitCode: 0,
			},
			"xdg-mime default nvim.desktop text/plain": {
				Stdout:   "",
				ExitCode: 0,
			},
			"xdg-mime default thunar.desktop inode/directory": {
				Stdout:   "",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	settingsDir := t.TempDir()
	settingsPath := filepath.Join(settingsDir, "settings.json")
	s := settings.Defaults()
	s.Defaults.Browser = "firefox"
	s.Defaults.Terminal = "kitty"
	s.Defaults.Editor = "nvim"
	s.Defaults.FileManager = "thunar"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureDefaultsResponse(t, "apply", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	results, ok := resp.Data["results"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected results map in data")
	}
	if len(results) != 4 {
		t.Errorf("expected 4 results, got %d", len(results))
	}
}

func TestApplyPartialFailureHandling(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"xdg-settings set default-web-browser firefox.desktop": {
				Stdout:   "",
				ExitCode: 0,
			},
			"xdg-mime default kitty.desktop x-scheme-handler/terminal": {
				Stdout:   "",
				ExitCode: 1,
				Stderr:   "xdg-mime: command failed",
			},
			"xdg-mime default nvim.desktop text/plain": {
				Stdout:   "",
				ExitCode: 0,
			},
			"xdg-mime default thunar.desktop inode/directory": {
				Stdout:   "",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	settingsDir := t.TempDir()
	settingsPath := filepath.Join(settingsDir, "settings.json")
	s := settings.Defaults()
	s.Defaults.Browser = "firefox"
	s.Defaults.Terminal = "kitty"
	s.Defaults.Editor = "nvim"
	s.Defaults.FileManager = "thunar"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureDefaultsResponse(t, "apply", "", settingsPath)
	// Partial failure returns warning status with per-item results
	if resp.Status != "warning" {
		t.Fatalf("expected status warning for partial failure, got %s: %s", resp.Status, resp.Message)
	}
	results, ok := resp.Data["results"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected results map in data")
	}
	termResult, ok := results["terminal"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected terminal result map")
	}
	if termResult["ok"] != false {
		t.Errorf("terminal ok = %v, want false", termResult["ok"])
	}
}

func TestDynamicAppDiscovery(t *testing.T) {
	tmpDir := t.TempDir()
	for _, app := range []string{"firefox.desktop", "chromium.desktop", "kitty.desktop"} {
		if err := os.WriteFile(filepath.Join(tmpDir, app), []byte("[Desktop Entry]"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	apps, err := discoverApps(tmpDir)
	if err != nil {
		t.Fatalf("discoverApps error: %v", err)
	}
	if len(apps) != 3 {
		t.Fatalf("expected 3 apps, got %d: %v", len(apps), apps)
	}
	expected := map[string]bool{"firefox": true, "chromium": true, "kitty": true}
	for _, app := range apps {
		if !expected[app] {
			t.Errorf("unexpected app %q", app)
		}
	}
}

func TestDesktopFileValidation(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "firefox.desktop"), []byte("[Desktop Entry]"), 0644); err != nil {
		t.Fatal(err)
	}

	if !validateDesktopFile(tmpDir, "firefox") {
		t.Error("expected firefox.desktop to be valid")
	}
	if validateDesktopFile(tmpDir, "nonexistent") {
		t.Error("expected nonexistent to be invalid")
	}
}

func TestApplyEmptyDefaultsReturnsOk(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{}
	defer func() { executor = oldExecutor }()

	settingsDir := t.TempDir()
	settingsPath := filepath.Join(settingsDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureDefaultsResponse(t, "apply", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok for empty defaults, got %s: %s", resp.Status, resp.Message)
	}
	results, ok := resp.Data["results"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected results map in data")
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty defaults, got %d", len(results))
	}
}
