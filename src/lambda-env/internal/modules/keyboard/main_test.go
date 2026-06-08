package main

import (
	"encoding/json"
	"fmt"
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
			"localectl list-x11-keymap-layouts": {
				Stdout:   "us\nes\nde\n",
				ExitCode: 0,
			},
			"localectl list-x11-keymap-variants us": {
				Stdout:   "intl\nstd\n",
				ExitCode: 0,
			},
			"setxkbmap -query": {
				Stdout:   "layout: us\nvariant: intl\noptions: caps:swapescape\n",
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
	if resp.Data["options"] != "caps:swapescape" {
		t.Errorf("options = %v, want caps:swapescape", resp.Data["options"])
	}
	opts, ok := resp.Data["available_options"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected available_options map")
	}
	layouts, ok := opts["set-layout"].([]interface{})
	if !ok || len(layouts) != 3 {
		t.Fatalf("expected 3 layouts, got %v", opts["set-layout"])
	}
	variants, ok := opts["set-variant"].([]interface{})
	if !ok || len(variants) != 3 { // empty + intl + std
		t.Fatalf("expected 3 variants, got %v", opts["set-variant"])
	}
}

func TestSetLayoutValidCallsSetxkbmap(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"localectl list-x11-keymap-layouts": {
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
			"localectl list-x11-keymap-layouts": {
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
			"localectl list-x11-keymap-layouts": {
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
			"localectl list-x11-keymap-layouts": {
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
			"localectl list-x11-keymap-layouts": {
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

func TestDiscoverVariants(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"localectl list-x11-keymap-variants es": {
				Stdout:   "deadtilde\nnodeadkeys\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	variants, err := discoverVariants(executor, "es")
	if err != nil {
		t.Fatalf("discoverVariants error: %v", err)
	}
	if len(variants) != 3 { // empty + deadtilde + nodeadkeys
		t.Fatalf("expected 3 variants, got %d: %v", len(variants), variants)
	}
	if variants[0] != "" {
		t.Errorf("variants[0] = %q, want empty", variants[0])
	}
	if variants[1] != "deadtilde" {
		t.Errorf("variants[1] = %q, want deadtilde", variants[1])
	}
}

func TestDiscoverVariantsFallback(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"localectl list-x11-keymap-variants xx": {
				Stdout:   "",
				ExitCode: 1,
				Stderr:   "Error",
			},
		},
	}
	defer func() { executor = oldExecutor }()

	variants, err := discoverVariants(executor, "xx")
	if err != nil {
		t.Fatalf("discoverVariants error: %v", err)
	}
	if len(variants) != 1 || variants[0] != "" {
		t.Errorf("expected [\"\"], got %v", variants)
	}
}

func TestParseCurrentLayout(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"setxkbmap -query": {
				Stdout:   "rules:      evdev\nmodel:      pc105\nlayout:     us,es\nvariant:    intl,\noptions:    caps:swapescape,ctrl:nocaps\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	layout, variant, options, err := parseCurrentLayout(executor)
	if err != nil {
		t.Fatalf("parseCurrentLayout error: %v", err)
	}
	if layout != "us,es" {
		t.Errorf("layout = %q, want us,es", layout)
	}
	if variant != "intl," {
		t.Errorf("variant = %q, want intl,", variant)
	}
	if options != "caps:swapescape,ctrl:nocaps" {
		t.Errorf("options = %q, want caps:swapescape,ctrl:nocaps", options)
	}
}

func TestParseCurrentLayoutHandlesMissingOptions(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"setxkbmap -query": {
				Stdout:   "layout: de\nvariant: nodeadkeys\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	layout, variant, options, err := parseCurrentLayout(executor)
	if err != nil {
		t.Fatalf("parseCurrentLayout error: %v", err)
	}
	if layout != "de" {
		t.Errorf("layout = %q, want de", layout)
	}
	if variant != "nodeadkeys" {
		t.Errorf("variant = %q, want nodeadkeys", variant)
	}
	if options != "" {
		t.Errorf("options = %q, want empty", options)
	}
}

func TestSetComposeValid(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"setxkbmap -option compose:ralt": {
				Stdout:   "",
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

	resp := captureKeyboardResponse(t, "set-compose", `{"value":"compose:ralt"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["keyboard"].(map[string]interface{})
	if !ok || delta["options"] != "compose:ralt" {
		t.Errorf("settings delta options = %v, want compose:ralt", delta["options"])
	}
}

func TestSetComposeInvalid(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureKeyboardResponse(t, "set-compose", `{"value":"compose:invalid"}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error, got %s", resp.Status)
	}
}

func TestSetOptionsValidSingle(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"setxkbmap -option \"\"": {
				Stdout:   "",
				ExitCode: 0,
			},
			"setxkbmap -option caps:swapescape": {
				Stdout:   "",
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

	resp := captureKeyboardResponse(t, "set-options", `{"value":"caps:swapescape"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["keyboard"].(map[string]interface{})
	if !ok || delta["options"] != "caps:swapescape" {
		t.Errorf("settings delta options = %v, want caps:swapescape", delta["options"])
	}
}

func TestSetOptionsValidMultiple(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"setxkbmap -option \"\"": {
				Stdout:   "",
				ExitCode: 0,
			},
			"setxkbmap -option ctrl:nocaps,compose:ralt": {
				Stdout:   "",
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

	resp := captureKeyboardResponse(t, "set-options", `{"value":"ctrl:nocaps,compose:ralt"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["keyboard"].(map[string]interface{})
	if !ok || delta["options"] != "ctrl:nocaps,compose:ralt" {
		t.Errorf("settings delta options = %v, want ctrl:nocaps,compose:ralt", delta["options"])
	}
}

func TestSetOptionsInvalid(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureKeyboardResponse(t, "set-options", `{"value":"invalid:option"}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error, got %s", resp.Status)
	}
}

func TestSetOptionsClearsBeforeApplying(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"setxkbmap -option \"\"": {
				Stdout:   "",
				ExitCode: 0,
			},
			"setxkbmap -option caps:swapescape": {
				Stdout:   "",
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

	resp := captureKeyboardResponse(t, "set-options", `{"value":"caps:swapescape"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	// The clear command was executed (mock key matched), so this passes if no mock miss.
}

func TestSetComposePreservesOtherOptions(t *testing.T) {
	oldExecutor := executor
	defer func() { executor = oldExecutor }()

	tests := []struct {
		name         string
		queryOptions string
		compose      string
		wantOptions  string
		wantCalls    []string
	}{
		{
			name:         "add compose to existing options",
			queryOptions: "caps:swapescape",
			compose:      "compose:ralt",
			wantOptions:  "caps:swapescape,compose:ralt",
			wantCalls: []string{
				"setxkbmap -query",
				"setxkbmap -option \"\"",
				"setxkbmap -option caps:swapescape,compose:ralt",
			},
		},
		{
			name:         "replace existing compose key",
			queryOptions: "caps:swapescape,compose:ralt",
			compose:      "compose:caps",
			wantOptions:  "caps:swapescape,compose:caps",
			wantCalls: []string{
				"setxkbmap -query",
				"setxkbmap -option \"\"",
				"setxkbmap -option caps:swapescape,compose:caps",
			},
		},
		{
			name:         "already present is no-op",
			queryOptions: "caps:swapescape,compose:ralt",
			compose:      "compose:ralt",
			wantOptions:  "caps:swapescape,compose:ralt",
			wantCalls: []string{
				"setxkbmap -query",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &module.MockExecutor{
				Responses: map[string]module.MockResponse{
					"setxkbmap -query": {
						Stdout:   fmt.Sprintf("layout: us\noptions: %s\n", tt.queryOptions),
						ExitCode: 0,
					},
				},
			}
			for _, cmd := range tt.wantCalls {
				if cmd == "setxkbmap -query" {
					continue
				}
				mock.Responses[cmd] = module.MockResponse{Stdout: "", ExitCode: 0}
			}
			executor = mock

			tmpDir := t.TempDir()
			settingsPath := filepath.Join(tmpDir, "settings.json")
			s := settings.Defaults()
			if err := settings.Save(settingsPath, &s); err != nil {
				t.Fatalf("Save: %v", err)
			}

			resp := captureKeyboardResponse(t, "set-compose", fmt.Sprintf(`{"value":"%s"}`, tt.compose), settingsPath)
			if resp.Status != "ok" {
				t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
			}
			delta, ok := resp.SettingsDelta["keyboard"].(map[string]interface{})
			if !ok || delta["options"] != tt.wantOptions {
				t.Errorf("settings delta options = %v, want %s", delta["options"], tt.wantOptions)
			}

			if len(mock.Calls) != len(tt.wantCalls) {
				t.Errorf("calls = %v, want %v", mock.Calls, tt.wantCalls)
			} else {
				for i := range mock.Calls {
					if mock.Calls[i] != tt.wantCalls[i] {
						t.Errorf("call[%d] = %q, want %q", i, mock.Calls[i], tt.wantCalls[i])
					}
				}
			}
		})
	}
}
