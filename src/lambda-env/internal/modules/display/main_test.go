package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

func captureDisplayResponse(t *testing.T, action, params, settingsPath string) module.Response {
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

func TestDetectSessionTypeX11(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"which xrandr": {
				Stdout:   "/usr/bin/xrandr\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	t.Setenv("XDG_SESSION_TYPE", "x11")
	session, err := detectSessionType(executor)
	if err != nil {
		t.Fatalf("detectSessionType error: %v", err)
	}
	if session != "x11" {
		t.Errorf("session = %q, want x11", session)
	}
}

func TestDetectSessionTypeWayland(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"which wlr-randr": {
				Stdout:   "/usr/bin/wlr-randr\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	t.Setenv("XDG_SESSION_TYPE", "wayland")
	session, err := detectSessionType(executor)
	if err != nil {
		t.Fatalf("detectSessionType error: %v", err)
	}
	if session != "wayland" {
		t.Errorf("session = %q, want wayland", session)
	}
}

func TestDetectSessionTypeFallback(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"which wlr-randr": {
				Stdout:   "",
				ExitCode: 1,
			},
			"which xrandr": {
				Stdout:   "/usr/bin/xrandr\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	os.Unsetenv("XDG_SESSION_TYPE")
	session, err := detectSessionType(executor)
	if err != nil {
		t.Fatalf("detectSessionType error: %v", err)
	}
	if session != "x11" {
		t.Errorf("session = %q, want x11", session)
	}
}

func TestParseXrandrOutput(t *testing.T) {
	stdout := `Screen 0: minimum 320 x 200, current 1920 x 1080, maximum 16384 x 16384
eDP-1 connected primary 1920x1080+0+0 (normal left inverted right x axis y axis) 309mm x 174mm
   1920x1080     60.03 +  59.97    59.96    59.93
   1680x1050     59.95    59.88
HDMI-1 disconnected (normal left inverted right x axis y axis)
HDMI-2 connected 1920x1080+1920+0 (normal left inverted right x axis y axis) 509mm x 286mm
   1920x1080     60.00 +  50.00    59.94
`
	outputs := parseXrandrOutput(stdout)
	if len(outputs) != 3 {
		t.Fatalf("expected 3 outputs, got %d", len(outputs))
	}

	edp := outputs[0]
	if edp.Name != "eDP-1" {
		t.Errorf("output[0].Name = %q, want eDP-1", edp.Name)
	}
	if !edp.Connected {
		t.Error("expected eDP-1 connected")
	}
	if !edp.Primary {
		t.Error("expected eDP-1 primary")
	}
	if edp.CurrentMode != "1920x1080" {
		t.Errorf("eDP-1 current mode = %q, want 1920x1080", edp.CurrentMode)
	}
	if len(edp.Modes) != 2 {
		t.Errorf("expected 2 modes for eDP-1, got %d", len(edp.Modes))
	}
	if edp.Modes[0].Resolution != "1920x1080" {
		t.Errorf("eDP-1 mode[0].Resolution = %q", edp.Modes[0].Resolution)
	}
	if !edp.Modes[0].Preferred {
		t.Error("expected eDP-1 1920x1080 to be preferred")
	}

	hdmi1 := outputs[1]
	if hdmi1.Connected {
		t.Error("expected HDMI-1 disconnected")
	}

	hdmi2 := outputs[2]
	if !hdmi2.Connected {
		t.Error("expected HDMI-2 connected")
	}
	if hdmi2.CurrentMode != "1920x1080" {
		t.Errorf("HDMI-2 current mode = %q", hdmi2.CurrentMode)
	}
}

func TestParseWlrRandrOutput(t *testing.T) {
	stdout := `HDMI-A-1 "ASUS 27" (HDMI-A-1)"
  Enabled: yes
  Mode: 1920x1080@60.000000Hz
  Position: 0,0
  Scale: 1.000000
  Transform: normal
DP-1 "Dell 24" (DP-1)"
  Enabled: no
  Mode: 2560x1440@144.000000Hz
  Position: 1920,0
`
	outputs := parseWlrRandrOutput(stdout)
	if len(outputs) != 2 {
		t.Fatalf("expected 2 outputs, got %d", len(outputs))
	}

	hdmi := outputs[0]
	if hdmi.Name != "HDMI-A-1" {
		t.Errorf("output[0].Name = %q", hdmi.Name)
	}
	if !hdmi.Connected {
		t.Error("expected HDMI-A-1 connected")
	}
	if hdmi.CurrentMode != "1920x1080" {
		t.Errorf("HDMI-A-1 current mode = %q", hdmi.CurrentMode)
	}
	if hdmi.Position != "0,0" {
		t.Errorf("HDMI-A-1 position = %q", hdmi.Position)
	}

	dp := outputs[1]
	if dp.Connected {
		t.Error("expected DP-1 disconnected")
	}
	if dp.CurrentMode != "2560x1440" {
		t.Errorf("DP-1 current mode = %q", dp.CurrentMode)
	}
}

func TestHandleRunWithXrandr(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"which xrandr": {
				Stdout:   "/usr/bin/xrandr\n",
				ExitCode: 0,
			},
			"xrandr --query": {
				Stdout: `eDP-1 connected primary 1920x1080+0+0
   1920x1080     60.03 +  59.97
HDMI-1 connected 1920x1080+1920+0
   1920x1080     60.00 +
`,
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

	t.Setenv("XDG_SESSION_TYPE", "x11")
	resp := captureDisplayResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	if resp.Data["session_type"] != "x11" {
		t.Errorf("session_type = %v, want x11", resp.Data["session_type"])
	}
	outputs, ok := resp.Data["outputs"].([]interface{})
	if !ok || len(outputs) != 2 {
		t.Fatalf("expected 2 outputs, got %v", resp.Data["outputs"])
	}
}

func TestHandleSetModeXrandr(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"which xrandr": {
				Stdout:   "/usr/bin/xrandr\n",
				ExitCode: 0,
			},
			"xrandr --query": {
				Stdout: `eDP-1 connected primary 1920x1080+0+0
   1920x1080     60.03 +  59.97
   1680x1050     59.95
`,
				ExitCode: 0,
			},
			"xrandr --output eDP-1 --mode 1920x1080 --rate 60.03": {
				Stdout:   "",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Display.ActiveProfile = "default"
	s.Display.Profiles = []settings.OutputProfile{
		{Name: "default", Outputs: []settings.OutputConfig{}},
	}
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	t.Setenv("XDG_SESSION_TYPE", "x11")
	resp := captureDisplayResponse(t, "set-mode", `{"value":"eDP-1:1920x1080@60.03"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
}

func TestHandleSetModeInvalid(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"which xrandr": {
				Stdout:   "/usr/bin/xrandr\n",
				ExitCode: 0,
			},
			"xrandr --query": {
				Stdout: `eDP-1 connected primary 1920x1080+0+0
   1920x1080     60.03 +
`,
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

	t.Setenv("XDG_SESSION_TYPE", "x11")
	resp := captureDisplayResponse(t, "set-mode", `{"value":"eDP-1:2560x1440@60"}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error, got %s", resp.Status)
	}
}

func TestHandleSetPosition(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"which xrandr": {
				Stdout:   "/usr/bin/xrandr\n",
				ExitCode: 0,
			},
			"xrandr --query": {
				Stdout: `eDP-1 connected primary 1920x1080+0+0
   1920x1080     60.03 +
HDMI-1 connected 1920x1080+1920+0
   1920x1080     60.00 +
`,
				ExitCode: 0,
			},
			"xrandr --output eDP-1 --left-of HDMI-1": {
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

	t.Setenv("XDG_SESSION_TYPE", "x11")
	resp := captureDisplayResponse(t, "set-position", `{"value":"--left-of HDMI-1"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
}

func TestHandleSetPrimary(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"which xrandr": {
				Stdout:   "/usr/bin/xrandr\n",
				ExitCode: 0,
			},
			"xrandr --query": {
				Stdout: `eDP-1 connected 1920x1080+0+0
   1920x1080     60.03 +
HDMI-1 connected 1920x1080+1920+0
   1920x1080     60.00 +
`,
				ExitCode: 0,
			},
			"xrandr --output eDP-1 --primary": {
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

	t.Setenv("XDG_SESSION_TYPE", "x11")
	resp := captureDisplayResponse(t, "set-primary", `{"value":true}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
}

func TestHandleSaveProfile(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"which xrandr": {
				Stdout:   "/usr/bin/xrandr\n",
				ExitCode: 0,
			},
			"xrandr --query": {
				Stdout: `eDP-1 connected primary 1920x1080+0+0
   1920x1080     60.03 +
HDMI-1 connected 1920x1080+1920+0
   1920x1080     60.00 +
`,
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

	t.Setenv("XDG_SESSION_TYPE", "x11")
	resp := captureDisplayResponse(t, "save-profile", `{"value":"Work Setup"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	found := false
	for _, p := range loaded.Display.Profiles {
		if p.Name == "Work Setup" {
			found = true
			if len(p.Outputs) != 2 {
				t.Errorf("expected 2 outputs in profile, got %d", len(p.Outputs))
			}
			break
		}
	}
	if !found {
		t.Error("profile 'Work Setup' not found in settings")
	}
	if loaded.Display.ActiveProfile != "Work Setup" {
		t.Errorf("ActiveProfile = %q, want Work Setup", loaded.Display.ActiveProfile)
	}
}

func TestHandleLoadProfile(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"which xrandr": {
				Stdout:   "/usr/bin/xrandr\n",
				ExitCode: 0,
			},
			"xrandr --query": {
				Stdout: `eDP-1 connected primary 1920x1080+0+0
   1920x1080     60.03 +
HDMI-1 connected 1920x1080+1920+0
   1920x1080     60.00 +
`,
				ExitCode: 0,
			},
			"xrandr --output eDP-1 --mode 1920x1080 --pos 0,0 --primary": {
				Stdout:   "",
				ExitCode: 0,
			},
			"xrandr --output HDMI-1 --mode 1920x1080 --pos 1920,0": {
				Stdout:   "",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Display.ActiveProfile = "Work Setup"
	s.Display.Profiles = []settings.OutputProfile{
		{
			Name: "Work Setup",
			Outputs: []settings.OutputConfig{
				{Name: "eDP-1", Mode: "1920x1080", Position: "0,0", Primary: true},
				{Name: "HDMI-1", Mode: "1920x1080", Position: "1920,0"},
			},
		},
	}
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	t.Setenv("XDG_SESSION_TYPE", "x11")
	resp := captureDisplayResponse(t, "load-profile", `{"value":"Work Setup"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Display.ActiveProfile != "Work Setup" {
		t.Errorf("ActiveProfile = %q, want Work Setup", loaded.Display.ActiveProfile)
	}
}

func TestHandleLoadProfileMissingOutput(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"which xrandr": {
				Stdout:   "/usr/bin/xrandr\n",
				ExitCode: 0,
			},
			"xrandr --query": {
				Stdout: `eDP-1 connected primary 1920x1080+0+0
   1920x1080     60.03 +
`,
				ExitCode: 0,
			},
			"xrandr --output eDP-1 --mode 1920x1080 --pos 0,0 --primary": {
				Stdout:   "",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Display.ActiveProfile = "Work Setup"
	s.Display.Profiles = []settings.OutputProfile{
		{
			Name: "Work Setup",
			Outputs: []settings.OutputConfig{
				{Name: "eDP-1", Mode: "1920x1080", Position: "0,0", Primary: true},
				{Name: "HDMI-1", Mode: "1920x1080", Position: "1920,0"},
			},
		},
	}
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	t.Setenv("XDG_SESSION_TYPE", "x11")
	resp := captureDisplayResponse(t, "load-profile", `{"value":"Work Setup"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	if !strings.Contains(resp.Message, "HDMI-1 is not connected") {
		t.Errorf("expected warning about missing HDMI-1, got: %s", resp.Message)
	}
}
