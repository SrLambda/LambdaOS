package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"lambdaos.dev/lambda-env/internal/hub"
	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

// setupTestHome creates a temporary home directory with module and config paths.
func setupTestHome(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	logDir := filepath.Join(tmpDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("create log dir: %v", err)
	}
	t.Setenv("LAMBDA_ENV_LOG_DIR", logDir)

	modulesDir := filepath.Join(tmpDir, ".local", "share", "lambda-env", "modules")
	if err := os.MkdirAll(modulesDir, 0755); err != nil {
		t.Fatalf("create modules dir: %v", err)
	}

	configDir := filepath.Join(tmpDir, ".config", "lambdaos")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("create config dir: %v", err)
	}

	return tmpDir
}

// writeModule creates a module directory with manifest and executable script.
func writeModule(t *testing.T, baseDir, name, manifest, script string) {
	t.Helper()
	modDir := filepath.Join(baseDir, name)
	if err := os.MkdirAll(modDir, 0755); err != nil {
		t.Fatalf("create module dir: %v", err)
	}
	manifestPath := filepath.Join(modDir, "manifest.json")
	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	scriptPath := filepath.Join(modDir, "module")
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		t.Fatalf("write script: %v", err)
	}
}

func TestDiscoveryWithFixtures(t *testing.T) {
	home := setupTestHome(t)
	modulesDir := filepath.Join(home, ".local", "share", "lambda-env", "modules")

	validManifest := `{
  "name": "screen",
  "version": "0.1.0",
  "description": "Manage display",
  "description_es": "Gestionar pantalla",
  "category": "system",
  "requires_root": false,
  "dependencies": [],
  "min_hub_version": "0.1.0"
}`
	screenScript := "#!/usr/bin/env bash\necho '{\"status\":\"ok\",\"action\":\"run\",\"data\":{}}'\n"
	writeModule(t, modulesDir, "screen", validManifest, screenScript)

	invalidManifest := `{"version": "0.1.0", "description": "broken"}`
	if err := os.MkdirAll(filepath.Join(modulesDir, "broken"), 0755); err != nil {
		t.Fatalf("create broken dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(modulesDir, "broken", "manifest.json"), []byte(invalidManifest), 0644); err != nil {
		t.Fatalf("write broken manifest: %v", err)
	}

	settingsPath := filepath.Join(home, ".config", "lambdaos", "settings.json")
	d := settings.Defaults()
	if err := settings.Save(settingsPath, &d); err != nil {
		t.Fatalf("save settings: %v", err)
	}

	h, err := hub.New(settingsPath)
	if err != nil {
		t.Fatalf("hub.New: %v", err)
	}
	defer h.Logger.Close()

	foundScreen := false
	for _, m := range h.Modules {
		if m.Name == "screen" {
			foundScreen = true
		}
		if m.Name == "broken" {
			t.Errorf("broken module should have been skipped")
		}
	}
	if !foundScreen {
		t.Errorf("expected screen module to be discovered")
	}
}

func TestModuleExecutionAndJSONParse(t *testing.T) {
	home := setupTestHome(t)
	modulesDir := filepath.Join(home, ".local", "share", "lambda-env", "modules")

	manifest := `{
  "name": "audio",
  "version": "0.1.0",
  "description": "Manage audio",
  "description_es": "Gestionar audio",
  "category": "system",
  "requires_root": false,
  "dependencies": [],
  "min_hub_version": "0.1.0",
  "timeout": 5
}`
	script := "#!/usr/bin/env bash\necho '{\"status\":\"ok\",\"action\":\"run\",\"data\":{\"sinks\":[]}}'\n"
	writeModule(t, modulesDir, "audio", manifest, script)

	settingsPath := filepath.Join(home, ".config", "lambdaos", "settings.json")
	d := settings.Defaults()
	if err := settings.Save(settingsPath, &d); err != nil {
		t.Fatalf("save settings: %v", err)
	}

	h, err := hub.New(settingsPath)
	if err != nil {
		t.Fatalf("hub.New: %v", err)
	}
	defer h.Logger.Close()

	var audioMod *module.Manifest
	for i := range h.Modules {
		if h.Modules[i].Name == "audio" {
			audioMod = &h.Modules[i]
			break
		}
	}
	if audioMod == nil {
		t.Fatalf("audio module not found")
	}

	resp, err := h.ExecuteModule(*audioMod)
	if err != nil {
		t.Fatalf("ExecuteModule: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("expected status ok, got %s", resp.Status)
	}
	if resp.Action != "run" {
		t.Errorf("expected action run, got %s", resp.Action)
	}
}

func TestSettingsDeltaMerge(t *testing.T) {
	home := setupTestHome(t)
	modulesDir := filepath.Join(home, ".local", "share", "lambda-env", "modules")

	manifest := `{
  "name": "screen",
  "version": "0.1.0",
  "description": "Manage display",
  "description_es": "Gestionar pantalla",
  "category": "system",
  "requires_root": false,
  "dependencies": [],
  "min_hub_version": "0.1.0"
}`
	script := "#!/usr/bin/env bash\necho '{\"status\":\"ok\",\"action\":\"run\",\"data\":{},\"settings_delta\":{\"display\":{\"active_profile\":\"home\"}}}'\n"
	writeModule(t, modulesDir, "screen", manifest, script)

	settingsPath := filepath.Join(home, ".config", "lambdaos", "settings.json")
	defaults := settings.Defaults()
	defaults.Display.ActiveProfile = "default"
	if err := settings.Save(settingsPath, &defaults); err != nil {
		t.Fatalf("save settings: %v", err)
	}

	h, err := hub.New(settingsPath)
	if err != nil {
		t.Fatalf("hub.New: %v", err)
	}
	defer h.Logger.Close()

	var screenMod *module.Manifest
	for i := range h.Modules {
		if h.Modules[i].Name == "screen" {
			screenMod = &h.Modules[i]
			break
		}
	}
	if screenMod == nil {
		t.Fatalf("screen module not found")
	}

	_, err = h.ExecuteModule(*screenMod)
	if err != nil {
		t.Fatalf("ExecuteModule: %v", err)
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("load settings after delta: %v", err)
	}
	if loaded.Display.ActiveProfile != "home" {
		t.Errorf("expected active_profile=home after delta merge, got %s", loaded.Display.ActiveProfile)
	}
}

func TestExecutionLog(t *testing.T) {
	home := setupTestHome(t)
	modulesDir := filepath.Join(home, ".local", "share", "lambda-env", "modules")

	manifest := `{
  "name": "logger-test",
  "version": "0.1.0",
  "description": "Test logging",
  "description_es": "Test logging",
  "category": "ops",
  "requires_root": false,
  "dependencies": [],
  "min_hub_version": "0.1.0"
}`
	script := "#!/usr/bin/env bash\n>&2 echo 'stderr msg'\necho '{\"status\":\"ok\",\"action\":\"run\"}'\n"
	writeModule(t, modulesDir, "logger-test", manifest, script)

	settingsPath := filepath.Join(home, ".config", "lambdaos", "settings.json")
	d := settings.Defaults()
	if err := settings.Save(settingsPath, &d); err != nil {
		t.Fatalf("save settings: %v", err)
	}

	h, err := hub.New(settingsPath)
	if err != nil {
		t.Fatalf("hub.New: %v", err)
	}
	defer h.Logger.Close()

	var mod *module.Manifest
	for i := range h.Modules {
		if h.Modules[i].Name == "logger-test" {
			mod = &h.Modules[i]
			break
		}
	}
	if mod == nil {
		t.Fatalf("logger-test module not found")
	}

	resp, err := h.ExecuteModule(*mod)
	if err != nil {
		t.Fatalf("ExecuteModule: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("expected ok, got %s", resp.Status)
	}

	// Verify log file contains execution record.
	logPath := filepath.Join(home, "logs", "modules.log")
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	logContent := string(data)
	if !strings.Contains(logContent, "module=logger-test") {
		t.Errorf("log missing module name")
	}
	if !strings.Contains(logContent, "stderr msg") {
		t.Errorf("log missing stderr content")
	}
	if !strings.Contains(logContent, "exit_code=0") {
		t.Errorf("log missing exit code")
	}
}

// TestFixtureFilesExist ensures the committed fixture files are present.
func TestFixtureFilesExist(t *testing.T) {
	fixtures := []string{
		"fixtures/modules/screen/manifest.json",
		"fixtures/modules/screen/module",
		"fixtures/modules/audio/manifest.json",
		"fixtures/modules/audio/module",
		"fixtures/modules/broken/manifest.json",
	}
	for _, f := range fixtures {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("fixture file missing: %s", f)
		}
	}
}
