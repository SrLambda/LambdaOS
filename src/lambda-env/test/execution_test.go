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

func TestExecuteActionSetsCorrectEnvVar(t *testing.T) {
	home := setupTestHome(t)
	modulesDir := filepath.Join(home, ".local", "share", "lambda-env", "modules")

	manifest := `{
  "name": "action-test",
  "version": "0.1.0",
  "description": "Test action env var",
  "description_es": "Test action env var",
  "category": "ops",
  "requires_root": false,
  "dependencies": [],
  "min_hub_version": "0.1.0"
}`
	modDir := filepath.Join(modulesDir, "action-test")
	// Script that writes env vars to files for verification
	script := "#!/usr/bin/env bash\n" +
		"echo \"$LAMBDA_ENV_ACTION\" > " + modDir + "/received_action.txt\n" +
		"echo \"$LAMBDA_ENV_PARAMS\" > " + modDir + "/received_params.txt\n" +
		"echo '{\"status\":\"ok\",\"action\":\"set-layout\"}'\n"
	writeModule(t, modulesDir, "action-test", manifest, script)

	settingsPath := filepath.Join(home, ".config", "lambdaos", "settings.json")
	d := settings.Defaults()
	if err := settings.Save(settingsPath, &d); err != nil {
		t.Fatalf("save settings: %v", err)
	}

	h, err := hub.New(settingsPath, false)
	if err != nil {
		t.Fatalf("hub.New: %v", err)
	}
	defer h.Logger.Close()

	var mod *module.Manifest
	for i := range h.Modules {
		if h.Modules[i].Name == "action-test" {
			mod = &h.Modules[i]
			break
		}
	}
	if mod == nil {
		t.Fatalf("action-test module not found")
	}

	resp, err := h.ExecuteAction(*mod, "set-layout", map[string]interface{}{"layout": "dvorak"})
	if err != nil {
		t.Fatalf("ExecuteAction: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("expected status ok, got %s", resp.Status)
	}

	// Verify the module received the correct action name
	actionBytes, err := os.ReadFile(filepath.Join(modDir, "received_action.txt"))
	if err != nil {
		t.Fatalf("read received_action: %v", err)
	}
	receivedAction := strings.TrimSpace(string(actionBytes))
	if receivedAction != "set-layout" {
		t.Errorf("expected action 'set-layout', got %q", receivedAction)
	}

	// Verify params were passed
	paramsBytes, err := os.ReadFile(filepath.Join(modDir, "received_params.txt"))
	if err != nil {
		t.Fatalf("read received_params: %v", err)
	}
	receivedParams := strings.TrimSpace(string(paramsBytes))
	if !strings.Contains(receivedParams, "dvorak") {
		t.Errorf("expected params to contain 'dvorak', got %q", receivedParams)
	}
}

func TestExecuteModuleBackwardCompatibility(t *testing.T) {
	home := setupTestHome(t)
	modulesDir := filepath.Join(home, ".local", "share", "lambda-env", "modules")

	manifest := `{
  "name": "compat-test",
  "version": "0.1.0",
  "description": "Test backward compatibility",
  "description_es": "Test backward compatibility",
  "category": "ops",
  "requires_root": false,
  "dependencies": [],
  "min_hub_version": "0.1.0"
}`
	script := "#!/usr/bin/env bash\n" +
		"action=$LAMBDA_ENV_ACTION\n" +
		"echo '{\"status\":\"ok\",\"action\":\"'\"$action\"'\"}'\n"
	writeModule(t, modulesDir, "compat-test", manifest, script)

	settingsPath := filepath.Join(home, ".config", "lambdaos", "settings.json")
	d := settings.Defaults()
	if err := settings.Save(settingsPath, &d); err != nil {
		t.Fatalf("save settings: %v", err)
	}

	h, err := hub.New(settingsPath, false)
	if err != nil {
		t.Fatalf("hub.New: %v", err)
	}
	defer h.Logger.Close()

	var mod *module.Manifest
	for i := range h.Modules {
		if h.Modules[i].Name == "compat-test" {
			mod = &h.Modules[i]
			break
		}
	}
	if mod == nil {
		t.Fatalf("compat-test module not found")
	}

	resp, err := h.ExecuteModule(*mod)
	if err != nil {
		t.Fatalf("ExecuteModule: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("expected status ok, got %s", resp.Status)
	}
	if resp.Action != "run" {
		t.Errorf("expected action 'run' for backward compatibility, got %s", resp.Action)
	}
}

func TestExecuteActionWithNilParams(t *testing.T) {
	home := setupTestHome(t)
	modulesDir := filepath.Join(home, ".local", "share", "lambda-env", "modules")

	manifest := `{
  "name": "nil-params-test",
  "version": "0.1.0",
  "description": "Test nil params",
  "description_es": "Test nil params",
  "category": "ops",
  "requires_root": false,
  "dependencies": [],
  "min_hub_version": "0.1.0"
}`
	script := "#!/usr/bin/env bash\n" +
		"params=$LAMBDA_ENV_PARAMS\n" +
		"echo '{\"status\":\"ok\",\"data\":{\"has_params\":\"'\"$params\"'\"}}'\n"
	writeModule(t, modulesDir, "nil-params-test", manifest, script)

	settingsPath := filepath.Join(home, ".config", "lambdaos", "settings.json")
	d := settings.Defaults()
	if err := settings.Save(settingsPath, &d); err != nil {
		t.Fatalf("save settings: %v", err)
	}

	h, err := hub.New(settingsPath, false)
	if err != nil {
		t.Fatalf("hub.New: %v", err)
	}
	defer h.Logger.Close()

	var mod *module.Manifest
	for i := range h.Modules {
		if h.Modules[i].Name == "nil-params-test" {
			mod = &h.Modules[i]
			break
		}
	}
	if mod == nil {
		t.Fatalf("nil-params-test module not found")
	}

	resp, err := h.ExecuteAction(*mod, "run", nil)
	if err != nil {
		t.Fatalf("ExecuteAction with nil params: %v", err)
	}

	// With nil params, LAMBDA_ENV_PARAMS may be empty or not set
	// The module should still execute successfully
	if resp.Status != "ok" {
		t.Errorf("expected status ok, got %s", resp.Status)
	}
}
