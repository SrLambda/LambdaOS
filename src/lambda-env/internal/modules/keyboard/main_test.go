package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

func captureKeyboardResponse(t *testing.T, action, params, settingsPath string) module.Response {
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

func TestHandleRunReturnsCurrentLayoutAndAvailableLayouts(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"setxkbmap -layout": {
				Stdout:   "us\nes\nde\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Keyboard.Layout = "us"
	s.Keyboard.Variant = "intl"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureKeyboardResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	if resp.Data["layout"] != "us" {
		t.Errorf("layout = %v, want us", resp.Data["layout"])
	}
	if resp.Data["variant"] != "intl" {
		t.Errorf("variant = %v, want intl", resp.Data["variant"])
	}
	opts, ok := resp.Data["available_options"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected available_options map")
	}
	layouts, ok := opts["set-layout"].([]interface{})
	if !ok || len(layouts) != 3 {
		t.Fatalf("expected 3 layouts, got %v", opts["set-layout"])
	}
}

func TestSetLayoutValidCallsSetxkbmap(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"setxkbmap -layout": {
				Stdout:   "us\nes\n",
				ExitCode: 0,
			},
			"setxkbmap -layout es": {
				Stdout:   "",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Keyboard.Layout = "us"
	s.Keyboard.Variant = ""
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureKeyboardResponse(t, "set-layout", `{"value":"es"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["keyboard"].(map[string]interface{})
	if !ok || delta["layout"] != "es" {
		t.Errorf("settings delta layout = %v, want es", delta["layout"])
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Keyboard.Layout != "es" {
		t.Errorf("Keyboard.Layout = %q, want es", loaded.Keyboard.Layout)
	}
}

func TestSetLayoutInvalidReturnsError(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"setxkbmap -layout": {
				Stdout:   "us\nes\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureKeyboardResponse(t, "set-layout", `{"value":"invalid"}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error, got %s", resp.Status)
	}
}

func TestSetVariantCallsSetxkbmap(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"setxkbmap -layout us -variant intl": {
				Stdout:   "",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Keyboard.Layout = "us"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureKeyboardResponse(t, "set-variant", `{"value":"intl"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["keyboard"].(map[string]interface{})
	if !ok || delta["variant"] != "intl" {
		t.Errorf("settings delta variant = %v, want intl", delta["variant"])
	}
}

func TestLayoutDiscoveryParsing(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"setxkbmap -layout": {
				Stdout:   "  us  \n\tes\n\n  de\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	layouts, err := discoverLayouts(executor)
	if err != nil {
		t.Fatalf("discoverLayouts error: %v", err)
	}
	if len(layouts) != 3 {
		t.Fatalf("expected 3 layouts, got %d: %v", len(layouts), layouts)
	}
	expected := []string{"us", "es", "de"}
	for i, exp := range expected {
		if layouts[i] != exp {
			t.Errorf("layout[%d] = %q, want %q", i, layouts[i], exp)
		}
	}
}

func TestLayoutDiscoveryFallback(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"setxkbmap -layout": {
				Stdout:   "",
				ExitCode: 1,
				Stderr:   "Error: invalid",
			},
		},
	}
	defer func() { executor = oldExecutor }()

	layouts, err := discoverLayouts(executor)
	if err != nil {
		t.Fatalf("discoverLayouts error: %v", err)
	}
	if len(layouts) == 0 {
		t.Error("expected fallback layouts, got empty")
	}
}

func TestSetLayoutWithVariantPreservesVariant(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"setxkbmap -layout": {
				Stdout:   "us\nes\n",
				ExitCode: 0,
			},
			"setxkbmap -layout es -variant intl": {
				Stdout:   "",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Keyboard.Layout = "us"
	s.Keyboard.Variant = "intl"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureKeyboardResponse(t, "set-layout", `{"value":"es"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Keyboard.Layout != "es" {
		t.Errorf("Keyboard.Layout = %q, want es", loaded.Keyboard.Layout)
	}
	// Variant should remain unchanged in settings (setxkbmap call included it)
	if loaded.Keyboard.Variant != "intl" {
		t.Errorf("Keyboard.Variant = %q, want intl", loaded.Keyboard.Variant)
	}
}
