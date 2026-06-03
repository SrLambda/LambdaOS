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

func capturePowerResponse(t *testing.T, action, params, settingsPath string) module.Response {
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

func TestHandleRunReturnsPowerSettingsAndBattery(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"upower -d": {
				Stdout:   "state: discharging\npercentage: 85%\ntime to empty: 3.5 hours\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Power.ScreenTimeout = 300
	s.Power.SleepTimeout = 600
	s.Power.LidCloseAction = "suspend"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := capturePowerResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	if resp.Data["screen_timeout"] != 300.0 {
		t.Errorf("screen_timeout = %v, want 300", resp.Data["screen_timeout"])
	}
	if resp.Data["sleep_timeout"] != 600.0 {
		t.Errorf("sleep_timeout = %v, want 600", resp.Data["sleep_timeout"])
	}
	if resp.Data["lid_close_action"] != "suspend" {
		t.Errorf("lid_close_action = %v, want suspend", resp.Data["lid_close_action"])
	}
	battery, ok := resp.Data["battery"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected battery map, got %T", resp.Data["battery"])
	}
	if battery["state"] != "discharging" {
		t.Errorf("battery state = %v, want discharging", battery["state"])
	}
	if battery["percentage"] != "85" {
		t.Errorf("battery percentage = %v, want 85", battery["percentage"])
	}
}

func TestHandleRunNoBattery(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"upower -d": {
				Stdout:   "",
				ExitCode: 1,
				Stderr:   "No battery",
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

	resp := capturePowerResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	if _, ok := resp.Data["battery"]; ok {
		t.Error("expected no battery data")
	}
	warning, ok := resp.Data["battery_warning"].(string)
	if !ok || warning == "" {
		t.Error("expected battery_warning to be set")
	}
}

func TestSetScreenTimeoutValid(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := capturePowerResponse(t, "set-screen-timeout", `{"value":"600"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["power"].(map[string]interface{})
	if !ok || delta["screen_timeout"] != 600.0 {
		t.Errorf("settings delta screen_timeout = %v, want 600", delta["screen_timeout"])
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Power.ScreenTimeout != 600 {
		t.Errorf("Power.ScreenTimeout = %d, want 600", loaded.Power.ScreenTimeout)
	}
}

func TestSetScreenTimeoutInvalid(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := capturePowerResponse(t, "set-screen-timeout", `{"value":"-10"}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error, got %s", resp.Status)
	}
}

func TestSetSleepTimeoutValid(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := capturePowerResponse(t, "set-sleep-timeout", `{"value":"1800"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["power"].(map[string]interface{})
	if !ok || delta["sleep_timeout"] != 1800.0 {
		t.Errorf("settings delta sleep_timeout = %v, want 1800", delta["sleep_timeout"])
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Power.SleepTimeout != 1800 {
		t.Errorf("Power.SleepTimeout = %d, want 1800", loaded.Power.SleepTimeout)
	}
}

func TestSetSleepTimeoutInvalid(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := capturePowerResponse(t, "set-sleep-timeout", `{"value":"abc"}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error, got %s", resp.Status)
	}
}

func TestSetLidCloseActionValid(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := capturePowerResponse(t, "set-lid-close-action", `{"value":"ignore"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["power"].(map[string]interface{})
	if !ok || delta["lid_close_action"] != "ignore" {
		t.Errorf("settings delta lid_close_action = %v, want ignore", delta["lid_close_action"])
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Power.LidCloseAction != "ignore" {
		t.Errorf("Power.LidCloseAction = %q, want ignore", loaded.Power.LidCloseAction)
	}
}

func TestSetLidCloseActionInvalid(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := capturePowerResponse(t, "set-lid-close-action", `{"value":"reboot"}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error, got %s", resp.Status)
	}
}

func TestParseLogindConf(t *testing.T) {
	data := `[Login]
IdleActionSec=450
HandleLidSwitch=hibernate
NAutoVTs=6
`
	screenTimeout, lidAction, err := parseLogindConf(data)
	if err != nil {
		t.Fatalf("parseLogindConf error: %v", err)
	}
	if screenTimeout != 450 {
		t.Errorf("screenTimeout = %d, want 450", screenTimeout)
	}
	if lidAction != "hibernate" {
		t.Errorf("lidAction = %q, want hibernate", lidAction)
	}
}

func TestParseLogindConfMissingKeys(t *testing.T) {
	data := `[Login]
NAutoVTs=6
`
	screenTimeout, lidAction, err := parseLogindConf(data)
	if err != nil {
		t.Fatalf("parseLogindConf error: %v", err)
	}
	if screenTimeout != 0 {
		t.Errorf("screenTimeout = %d, want 0", screenTimeout)
	}
	if lidAction != "" {
		t.Errorf("lidAction = %q, want empty", lidAction)
	}
}

func TestUpdateLogindConfKeyReplace(t *testing.T) {
	data := `[Login]
IdleActionSec=300
HandleLidSwitch=suspend
`
	updated := updateLogindConfKey(data, "IdleActionSec", "600")
	if !strings.Contains(updated, "IdleActionSec=600") {
		t.Errorf("expected IdleActionSec=600 in updated config")
	}
	if strings.Contains(updated, "IdleActionSec=300") {
		t.Error("old IdleActionSec value should have been replaced")
	}
}

func TestUpdateLogindConfKeyAppend(t *testing.T) {
	data := `[Login]
HandleLidSwitch=suspend
`
	updated := updateLogindConfKey(data, "IdleActionSec", "450")
	if !strings.Contains(updated, "IdleActionSec=450") {
		t.Errorf("expected IdleActionSec=450 to be appended")
	}
}

func TestUpdateLogindConfKeyNewSection(t *testing.T) {
	data := `SomeOtherKey=value
`
	updated := updateLogindConfKey(data, "HandleLidSwitch", "ignore")
	if !strings.Contains(updated, "[Login]") {
		t.Error("expected [Login] section to be created")
	}
	if !strings.Contains(updated, "HandleLidSwitch=ignore") {
		t.Errorf("expected HandleLidSwitch=ignore to be added")
	}
}

func TestParseUpowerOutput(t *testing.T) {
	stdout := `state: charging
percentage: 92%
time to empty: 2.1 hours
`
	battery := parseUpowerOutput(stdout)
	if battery == nil {
		t.Fatal("expected battery data")
	}
	if battery["state"] != "charging" {
		t.Errorf("state = %v, want charging", battery["state"])
	}
	if battery["percentage"] != "92" {
		t.Errorf("percentage = %v, want 92", battery["percentage"])
	}
	if battery["time_remaining"] != "2.1 hours" {
		t.Errorf("time_remaining = %v, want 2.1 hours", battery["time_remaining"])
	}
}

func TestParseUpowerOutputEmpty(t *testing.T) {
	battery := parseUpowerOutput("")
	if battery != nil {
		t.Error("expected nil for empty upower output")
	}
}

func TestReadSysBattery(t *testing.T) {
	batDir := t.TempDir()
	ueventPath := filepath.Join(batDir, "uevent")
	content := "POWER_SUPPLY_CAPACITY=77\nPOWER_SUPPLY_STATUS=Charging\n"
	if err := os.WriteFile(ueventPath, []byte(content), 0644); err != nil {
		t.Fatalf("write uevent: %v", err)
	}

	oldPath := sysBatteryPath
	sysBatteryPath = ueventPath
	defer func() { sysBatteryPath = oldPath }()

	battery, err := readSysBattery()
	if err != nil {
		t.Fatalf("readSysBattery error: %v", err)
	}
	if battery["percentage"] != "77" {
		t.Errorf("percentage = %v, want 77", battery["percentage"])
	}
	if battery["state"] != "charging" {
		t.Errorf("state = %v, want charging", battery["state"])
	}
}
