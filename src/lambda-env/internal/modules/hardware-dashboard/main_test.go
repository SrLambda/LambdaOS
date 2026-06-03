package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

func captureDashboardResponse(t *testing.T, action, params, settingsPath string) module.Response {
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

func TestHandleRunHappyPath(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"lscpu": {
				Stdout:   "Model name: Intel(R) Core(TM) i7-1165G7\nCPU(s): 8\nArchitecture: x86_64\n",
				ExitCode: 0,
			},
			"free -m": {
				Stdout:   "              total        used        free      shared  buff/cache   available\nMem:          16384        4096         256        1024       12032       11264\nSwap:          8192           0        8192\n",
				ExitCode: 0,
			},
			"df -h /": {
				Stdout:   "Filesystem      Size  Used Avail Use% Mounted on\n/dev/nvme0n1p1  250G   45G  205G  18% /\n",
				ExitCode: 0,
			},
			"sensors": {
				Stdout:   "coretemp-isa-0000\nCore 0: +52.0°C\nCore 1: +48.0°C\namdgpu-pci-0300\ntemp1: +45.0°C\n",
				ExitCode: 0,
			},
			"upower -d": {
				Stdout:   "state: discharging\npercentage: 85%\ntime to empty: 3.5 hours\n",
				ExitCode: 0,
			},
			"uptime -p": {
				Stdout:   "up 2 hours, 15 minutes\n",
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

	resp := captureDashboardResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}

	// Verify SettingsDelta is never emitted.
	if resp.SettingsDelta != nil {
		t.Fatalf("expected nil SettingsDelta for read-only module, got %v", resp.SettingsDelta)
	}

	cpu, ok := resp.Data["cpu"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected cpu map, got %T", resp.Data["cpu"])
	}
	if cpu["model"] != "Intel(R) Core(TM) i7-1165G7" {
		t.Errorf("cpu model = %v, want Intel(R) Core(TM) i7-1165G7", cpu["model"])
	}
	if cpu["cores"] != "8" {
		t.Errorf("cpu cores = %v, want 8", cpu["cores"])
	}

	ram, ok := resp.Data["ram"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected ram map, got %T", resp.Data["ram"])
	}
	if ram["total_mb"] != "16384" {
		t.Errorf("ram total = %v, want 16384", ram["total_mb"])
	}
	if ram["used_mb"] != "4096" {
		t.Errorf("ram used = %v, want 4096", ram["used_mb"])
	}

	disk, ok := resp.Data["disk"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected disk map, got %T", resp.Data["disk"])
	}
	if disk["filesystem"] != "/dev/nvme0n1p1" {
		t.Errorf("disk filesystem = %v, want /dev/nvme0n1p1", disk["filesystem"])
	}
	if disk["use_pct"] != "18%" {
		t.Errorf("disk use%% = %v, want 18%%", disk["use_pct"])
	}

	temps, ok := resp.Data["temperatures"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected temperatures map, got %T", resp.Data["temperatures"])
	}
	if temps["Core 0"] != "52.0" {
		t.Errorf("temp Core 0 = %v, want 52.0", temps["Core 0"])
	}

	battery, ok := resp.Data["battery"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected battery map, got %T", resp.Data["battery"])
	}
	if battery["percentage"] != "85" {
		t.Errorf("battery percentage = %v, want 85", battery["percentage"])
	}
	if battery["state"] != "discharging" {
		t.Errorf("battery state = %v, want discharging", battery["state"])
	}

	uptime, ok := resp.Data["uptime"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected uptime map, got %T", resp.Data["uptime"])
	}
	if uptime["human_readable"] != "up 2 hours, 15 minutes" {
		t.Errorf("uptime = %v, want 'up 2 hours, 15 minutes'", uptime["human_readable"])
	}
}

func TestHandleRunSensorsNotInstalled(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"lscpu": {
				Stdout:   "CPU(s): 4\n",
				ExitCode: 0,
			},
			"free -m": {
				Stdout:   "Mem: 8192 2048 6144\n",
				ExitCode: 0,
			},
			"df -h /": {
				Stdout:   "Filesystem Size Used Avail Use%\n/dev/sda1 100G 20G 80G 20% /\n",
				ExitCode: 0,
			},
			"sensors": {
				Stdout:   "",
				ExitCode: 127,
				Stderr:   "command not found",
			},
			"upower -d": {
				Stdout:   "state: charging\npercentage: 60%\n",
				ExitCode: 0,
			},
			"uptime -p": {
				Stdout:   "up 1 day, 3 hours\n",
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

	resp := captureDashboardResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}

	// Temperatures should fall back to thermal zones or show N/A.
	temps, ok := resp.Data["temperatures"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected temperatures map, got %T", resp.Data["temperatures"])
	}
	// If thermal zones are not present on this test runner, status should be N/A.
	if _, hasStatus := temps["status"]; hasStatus {
		if temps["status"] != "N/A" {
			t.Errorf("expected temps status N/A when sensors unavailable, got %v", temps["status"])
		}
	}
}

func TestHandleRunNoBattery(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"lscpu": {
				Stdout:   "CPU(s): 4\n",
				ExitCode: 0,
			},
			"free -m": {
				Stdout:   "Mem: 8192 2048 6144\n",
				ExitCode: 0,
			},
			"df -h /": {
				Stdout:   "Filesystem Size Used Avail Use%\n/dev/sda1 100G 20G 80G 20% /\n",
				ExitCode: 0,
			},
			"sensors": {
				Stdout:   "Core 0: +40.0°C\n",
				ExitCode: 0,
			},
			"upower -d": {
				Stdout:   "",
				ExitCode: 1,
				Stderr:   "No battery",
			},
			"uptime -p": {
				Stdout:   "up 5 minutes\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	// Ensure sysfs battery path points nowhere.
	oldPath := sysBatteryPath
	sysBatteryPath = filepath.Join(t.TempDir(), "nonexistent")
	defer func() { sysBatteryPath = oldPath }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureDashboardResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}

	battery, ok := resp.Data["battery"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected battery map, got %T", resp.Data["battery"])
	}
	if battery["status"] != "no battery detected" {
		t.Errorf("expected 'no battery detected', got %v", battery["status"])
	}
}

func TestHandleRunPartialFailure(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"lscpu": {
				Stdout:   "",
				ExitCode: 1,
			},
			"free -m": {
				Stdout:   "Mem: 8192 2048 6144\n",
				ExitCode: 0,
			},
			"df -h /": {
				Stdout:   "Filesystem Size Used Avail Use%\n/dev/sda1 100G 20G 80G 20% /\n",
				ExitCode: 0,
			},
			"sensors": {
				Stdout:   "",
				ExitCode: 127,
				Stderr:   "command not found",
			},
			"upower -d": {
				Stdout:   "",
				ExitCode: 1,
				Stderr:   "No battery",
			},
			"uptime -p": {
				Stdout:   "",
				ExitCode: 1,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	oldPath := sysBatteryPath
	sysBatteryPath = filepath.Join(t.TempDir(), "nonexistent")
	defer func() { sysBatteryPath = oldPath }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureDashboardResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}

	// Disk should still be present despite other failures.
	disk, ok := resp.Data["disk"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected disk map, got %T", resp.Data["disk"])
	}
	if disk["filesystem"] != "/dev/sda1" {
		t.Errorf("disk filesystem = %v, want /dev/sda1", disk["filesystem"])
	}

	// CPU should show N/A since lscpu failed and /proc/loadavg is host-dependent.
	cpu, ok := resp.Data["cpu"].(map[string]interface{})
	if ok {
		if cpu["model"] != "N/A" {
			t.Errorf("expected cpu model N/A on lscpu failure, got %v", cpu["model"])
		}
	}
}

func TestHandleRunNeverEmitsSettingsDelta(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"lscpu":      {Stdout: "CPU(s): 2\n", ExitCode: 0},
			"free -m":    {Stdout: "Mem: 4096 1024 3072\n", ExitCode: 0},
			"df -h /":    {Stdout: "Filesystem Size Used Avail Use%\n/dev/sda1 50G 10G 40G 20% /\n", ExitCode: 0},
			"sensors":    {Stdout: "", ExitCode: 127},
			"upower -d":  {Stdout: "", ExitCode: 1},
			"uptime -p":  {Stdout: "up 10 minutes\n", ExitCode: 0},
		},
	}
	defer func() { executor = oldExecutor }()

	oldPath := sysBatteryPath
	sysBatteryPath = filepath.Join(t.TempDir(), "nonexistent")
	defer func() { sysBatteryPath = oldPath }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureDashboardResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	if resp.SettingsDelta != nil {
		t.Fatalf("read-only module must never emit SettingsDelta, got %v", resp.SettingsDelta)
	}
}
