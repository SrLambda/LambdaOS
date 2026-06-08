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

func captureAudioResponse(t *testing.T, action, params, settingsPath string) module.Response {
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

func TestHandleRunReturnsCurrentVolumeAndSinks(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl info": {
				Stdout:   "Server Name: pulseaudio\nDefault Sink: alsa_output.pci-0000_00_1f.3.analog-stereo\n",
				ExitCode: 0,
			},
			"pactl list short sinks": {
				Stdout:   "0\talsa_output.pci-0000_00_1f.3.analog-stereo\tmodule-alsa-card.c\ts16le 2ch 44100Hz\tRUNNING\n",
				ExitCode: 0,
			},
			"pactl list short sources": {
				Stdout:   "0\talsa_input.pci-0000_00_1f.3.analog-stereo\tmodule-alsa-card.c\ts16le 2ch 44100Hz\tRUNNING\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Audio.Volume = 60
	s.Audio.Muted = false
	s.Audio.DefaultSink = "alsa_output.pci-0000_00_1f.3.analog-stereo"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAudioResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	if resp.Data["volume"] != float64(60) {
		t.Errorf("volume = %v, want 60", resp.Data["volume"])
	}
	if resp.Data["muted"] != false {
		t.Errorf("muted = %v, want false", resp.Data["muted"])
	}
	if resp.Data["default_sink"] != "alsa_output.pci-0000_00_1f.3.analog-stereo" {
		t.Errorf("default_sink = %v, want alsa_output.pci-0000_00_1f.3.analog-stereo", resp.Data["default_sink"])
	}
	if resp.Data["backend"] != "pulseaudio" {
		t.Errorf("backend = %v, want pulseaudio", resp.Data["backend"])
	}
	opts, ok := resp.Data["available_options"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected available_options map")
	}
	sinks, ok := opts["set-sink"].([]interface{})
	if !ok || len(sinks) != 1 {
		t.Fatalf("expected 1 sink, got %v", opts["set-sink"])
	}
	if sinks[0] != "alsa_output.pci-0000_00_1f.3.analog-stereo" {
		t.Errorf("sink[0] = %v, want alsa_output.pci-0000_00_1f.3.analog-stereo", sinks[0])
	}
	sources, ok := opts["set-source"].([]interface{})
	if !ok || len(sources) != 1 {
		t.Fatalf("expected 1 source, got %v", opts["set-source"])
	}
	if sources[0] != "alsa_input.pci-0000_00_1f.3.analog-stereo" {
		t.Errorf("source[0] = %v, want alsa_input.pci-0000_00_1f.3.analog-stereo", sources[0])
	}
	// Per-app should show warning for pulseaudio
	if resp.Data["per_app_warning"] != "per-app volume requires PipeWire" {
		t.Errorf("per_app_warning = %v", resp.Data["per_app_warning"])
	}
}

func TestSetVolumeValidCallsPactl(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl set-sink-volume alsa_output.pci-0000_00_1f.3.analog-stereo 45%": {
				Stdout:   "",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Audio.DefaultSink = "alsa_output.pci-0000_00_1f.3.analog-stereo"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAudioResponse(t, "set-volume", `{"value":45}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["audio"].(map[string]interface{})
	if !ok || delta["volume"] != float64(45) {
		t.Errorf("settings delta volume = %v, want 45", delta["volume"])
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Audio.Volume != 45 {
		t.Errorf("Audio.Volume = %d, want 45", loaded.Audio.Volume)
	}
}

func TestSetVolumeClamping(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl set-sink-volume alsa_output.pci-0000_00_1f.3.analog-stereo 0%": {
				Stdout:   "",
				ExitCode: 0,
			},
			"pactl set-sink-volume alsa_output.pci-0000_00_1f.3.analog-stereo 100%": {
				Stdout:   "",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Audio.DefaultSink = "alsa_output.pci-0000_00_1f.3.analog-stereo"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Clamp below 0 to 0
	resp := captureAudioResponse(t, "set-volume", `{"value":-10}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok for negative volume, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["audio"].(map[string]interface{})
	if !ok || delta["volume"] != float64(0) {
		t.Errorf("settings delta volume = %v, want 0", delta["volume"])
	}

	// Reset settings for second test
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Clamp above 100 to 100
	resp = captureAudioResponse(t, "set-volume", `{"value":150}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok for over-max volume, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok = resp.SettingsDelta["audio"].(map[string]interface{})
	if !ok || delta["volume"] != float64(100) {
		t.Errorf("settings delta volume = %v, want 100", delta["volume"])
	}
}

func TestSetVolumeInvalidStringReturnsError(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Audio.DefaultSink = "alsa_output.pci-0000_00_1f.3.analog-stereo"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAudioResponse(t, "set-volume", `{"value":"not-a-number"}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error for invalid volume string, got %s", resp.Status)
	}
}

func TestSetMuteCallsPactl(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl set-sink-mute alsa_output.pci-0000_00_1f.3.analog-stereo 1": {
				Stdout:   "",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Audio.DefaultSink = "alsa_output.pci-0000_00_1f.3.analog-stereo"
	s.Audio.Muted = false
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAudioResponse(t, "set-mute", `{"value":true}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["audio"].(map[string]interface{})
	if !ok || delta["muted"] != true {
		t.Errorf("settings delta muted = %v, want true", delta["muted"])
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !loaded.Audio.Muted {
		t.Errorf("Audio.Muted = %v, want true", loaded.Audio.Muted)
	}
}

func TestSetSinkValidCallsPactl(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl list short sinks": {
				Stdout:   "0\talsa_output.pci-0000_00_1f.3.analog-stereo\tmodule-alsa-card.c\ts16le 2ch 44100Hz\tRUNNING\n1\tbluez_sink.00_00_00_00_00_00.a2dp-sink\tmodule-bluez5-device.c\ts16le 2ch 44100Hz\tRUNNING\n",
				ExitCode: 0,
			},
			"pactl set-default-sink bluez_sink.00_00_00_00_00_00.a2dp-sink": {
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

	resp := captureAudioResponse(t, "set-sink", `{"value":"bluez_sink.00_00_00_00_00_00.a2dp-sink"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["audio"].(map[string]interface{})
	if !ok || delta["default_sink"] != "bluez_sink.00_00_00_00_00_00.a2dp-sink" {
		t.Errorf("settings delta default_sink = %v", delta["default_sink"])
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Audio.DefaultSink != "bluez_sink.00_00_00_00_00_00.a2dp-sink" {
		t.Errorf("Audio.DefaultSink = %q", loaded.Audio.DefaultSink)
	}
}

func TestSetSinkInvalidReturnsError(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl list short sinks": {
				Stdout:   "0\talsa_output.pci-0000_00_1f.3.analog-stereo\tmodule-alsa-card.c\ts16le 2ch 44100Hz\tRUNNING\n",
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

	resp := captureAudioResponse(t, "set-sink", `{"value":"nonexistent"}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error, got %s", resp.Status)
	}
}

func TestBackendDetectionPipewire(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl info": {
				Stdout:   "Server Name: PulseAudio (on PipeWire 0.3.48)\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	backend, err := detectBackend(executor)
	if err != nil {
		t.Fatalf("detectBackend error: %v", err)
	}
	if backend != "pipewire" {
		t.Errorf("backend = %q, want pipewire", backend)
	}
}

func TestBackendDetectionPulseaudio(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl info": {
				Stdout:   "Server Name: pulseaudio\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	backend, err := detectBackend(executor)
	if err != nil {
		t.Fatalf("detectBackend error: %v", err)
	}
	if backend != "pulseaudio" {
		t.Errorf("backend = %q, want pulseaudio", backend)
	}
}

func TestDynamicSinkDiscovery(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl list short sinks": {
				Stdout: "0\tsink1\tmodule-alsa-card.c\ts16le 2ch 44100Hz\tRUNNING\n" +
					"1\tsink2\tmodule-bluez5-device.c\ts16le 2ch 44100Hz\tIDLE\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	sinks, err := discoverSinks(executor)
	if err != nil {
		t.Fatalf("discoverSinks error: %v", err)
	}
	if len(sinks) != 2 {
		t.Fatalf("expected 2 sinks, got %d: %v", len(sinks), sinks)
	}
	if sinks[0] != "sink1" {
		t.Errorf("sinks[0] = %q, want sink1", sinks[0])
	}
	if sinks[1] != "sink2" {
		t.Errorf("sinks[1] = %q, want sink2", sinks[1])
	}
}

func TestBackendDetectionFailure(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl info": {
				Stdout:   "",
				ExitCode: 1,
				Stderr:   "Connection refused",
			},
		},
	}
	defer func() { executor = oldExecutor }()

	backend, err := detectBackend(executor)
	if err == nil {
		t.Fatal("expected error for failed pactl info, got nil")
	}
	if backend != "" {
		t.Errorf("backend = %q, want empty string on error", backend)
	}
}

func TestDiscoverSources(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl list short sources": {
				Stdout: "0\talsa_input.usb\tmodule-alsa-card.c\ts16le 2ch 44100Hz\tRUNNING\n" +
					"1\talsa_input.pci\tmodule-alsa-card.c\ts16le 2ch 44100Hz\tIDLE\n",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	sources, err := discoverSources(executor)
	if err != nil {
		t.Fatalf("discoverSources error: %v", err)
	}
	if len(sources) != 2 {
		t.Fatalf("expected 2 sources, got %d: %v", len(sources), sources)
	}
	if sources[0] != "alsa_input.usb" {
		t.Errorf("sources[0] = %q, want alsa_input.usb", sources[0])
	}
	if sources[1] != "alsa_input.pci" {
		t.Errorf("sources[1] = %q, want alsa_input.pci", sources[1])
	}
}

func TestSetSourceValid(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl list short sources": {
				Stdout:   "0\talsa_input.usb\tmodule-alsa-card.c\ts16le 2ch 44100Hz\tRUNNING\n",
				ExitCode: 0,
			},
			"pactl set-default-source alsa_input.usb": {
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

	resp := captureAudioResponse(t, "set-source", `{"value":"alsa_input.usb"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["audio"].(map[string]interface{})
	if !ok || delta["default_source"] != "alsa_input.usb" {
		t.Errorf("settings delta default_source = %v", delta["default_source"])
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Audio.DefaultSource != "alsa_input.usb" {
		t.Errorf("Audio.DefaultSource = %q", loaded.Audio.DefaultSource)
	}
}

func TestSetSourceInvalid(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl list short sources": {
				Stdout:   "0\talsa_input.usb\tmodule-alsa-card.c\ts16le 2ch 44100Hz\tRUNNING\n",
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

	resp := captureAudioResponse(t, "set-source", `{"value":"nonexistent"}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error, got %s", resp.Status)
	}
}

func TestSetProfileValid(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Audio.Profiles = []settings.AudioProfile{
		{Name: "Headphones", DefaultSink: "alsa_output.usb", DefaultSource: "alsa_input.usb", Volume: 80},
	}
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAudioResponse(t, "set-profile", `{"value":"Headphones"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["audio"].(map[string]interface{})
	if !ok || delta["profile"] != "Headphones" {
		t.Errorf("settings delta profile = %v", delta["profile"])
	}
}

func TestSetProfileInvalid(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAudioResponse(t, "set-profile", `{"value":"NonExistent"}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error, got %s", resp.Status)
	}
}

func TestSetProfileClear(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Audio.Profile = "Headphones"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAudioResponse(t, "set-profile", `{"value":""}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["audio"].(map[string]interface{})
	if !ok || delta["profile"] != "" {
		t.Errorf("settings delta profile = %v, want empty", delta["profile"])
	}
}

func TestPerAppVolumesPipewire(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl info": {
				Stdout:   "Server Name: PulseAudio (on PipeWire 0.3.48)\n",
				ExitCode: 0,
			},
			"pactl list short sinks": {
				Stdout:   "0\talsa_output.pci\tmodule-alsa-card.c\ts16le 2ch 44100Hz\tRUNNING\n",
				ExitCode: 0,
			},
			"pactl list short sources": {
				Stdout:   "0\talsa_input.pci\tmodule-alsa-card.c\ts16le 2ch 44100Hz\tRUNNING\n",
				ExitCode: 0,
			},
			"wpctl status": {
				Stdout: "PipeWire 'pipewire-0' [0.3.48]\n" +
					" └─ Clients:\n" +
					" └─ Sinks:\n" +
					" └─ Sources:\n" +
					" └─ Streams:\n" +
					"        55. firefox       [Stream/Output/Audio]\n" +
					"        56. spotify       [Stream/Output/Audio]\n",
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

	resp := captureAudioResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	perApp, ok := resp.Data["per_app_volumes"].([]interface{})
	if !ok || len(perApp) != 2 {
		t.Fatalf("expected 2 per-app volumes, got %v", resp.Data["per_app_volumes"])
	}
	if _, hasWarning := resp.Data["per_app_warning"]; hasWarning {
		t.Errorf("unexpected per_app_warning in response")
	}
}

func TestPerAppVolumesGracefulDegradation(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl info": {
				Stdout:   "Server Name: PulseAudio (on PipeWire 0.3.48)\n",
				ExitCode: 0,
			},
			"pactl list short sinks": {
				Stdout:   "0\talsa_output.pci\tmodule-alsa-card.c\ts16le 2ch 44100Hz\tRUNNING\n",
				ExitCode: 0,
			},
			"pactl list short sources": {
				Stdout:   "0\talsa_input.pci\tmodule-alsa-card.c\ts16le 2ch 44100Hz\tRUNNING\n",
				ExitCode: 0,
			},
			"wpctl status": {
				Stdout:   "",
				ExitCode: 1,
				Stderr:   "wpctl: command not found",
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

	resp := captureAudioResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	if resp.Data["per_app_warning"] != "wpctl not available or PipeWire not running" {
		t.Errorf("per_app_warning = %v", resp.Data["per_app_warning"])
	}
}

func TestSetAppVolumeValid(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"wpctl set-volume 55 0.80": {
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

	resp := captureAudioResponse(t, "set-app-volume", `{"value":"55:0.80"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	// No settings delta for per-app volume.
	if resp.SettingsDelta != nil {
		t.Errorf("expected nil SettingsDelta, got %v", resp.SettingsDelta)
	}
}

func TestSetAppVolumeInvalidFormat(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAudioResponse(t, "set-app-volume", `{"value":"invalid"}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error, got %s", resp.Status)
	}
}

func TestDiscoverPerAppVolumesReal(t *testing.T) {
	oldExecutor := executor
	defer func() { executor = oldExecutor }()

	mock := &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"wpctl status": {
				Stdout: "PipeWire 'pipewire-0' [0.3.48]\n" +
					" \u2514\u2500 Clients:\n" +
					" \u2514\u2500 Sinks:\n" +
					" \u2514\u2500 Sources:\n" +
					" \u2514\u2500 Streams:\n" +
					"        55. firefox       [Stream/Output/Audio]\n" +
					"        56. spotify       [Stream/Output/Audio]\n",
				ExitCode: 0,
			},
			"wpctl inspect 55": {
				Stdout: "id: 55\n" +
					"  Volume: 0.40\n" +
					"  Mute: 0\n",
				ExitCode: 0,
			},
			"wpctl inspect 56": {
				Stdout: "id: 56\n" +
					"  Volume: 0.75\n" +
					"  Mute: 0\n",
				ExitCode: 0,
			},
		},
	}
	executor = mock

	apps, warning, err := discoverPerAppVolumes(executor)
	if err != nil {
		t.Fatalf("discoverPerAppVolumes error: %v", err)
	}
	if warning != "" {
		t.Errorf("unexpected warning: %s", warning)
	}
	if len(apps) != 2 {
		t.Fatalf("expected 2 apps, got %d", len(apps))
	}
	if apps[0].ID != "55" || apps[0].Name != "firefox" {
		t.Errorf("app[0] = %+v", apps[0])
	}
	if apps[0].Volume != 0.40 {
		t.Errorf("app[0].Volume = %v, want 0.40", apps[0].Volume)
	}
	if apps[1].ID != "56" || apps[1].Name != "spotify" {
		t.Errorf("app[1] = %+v", apps[1])
	}
	if apps[1].Volume != 0.75 {
		t.Errorf("app[1].Volume = %v, want 0.75", apps[1].Volume)
	}

	wantCalls := []string{
		"wpctl status",
		"wpctl inspect 55",
		"wpctl inspect 56",
	}
	if len(mock.Calls) != len(wantCalls) {
		t.Errorf("calls = %v, want %v", mock.Calls, wantCalls)
	} else {
		for i := range mock.Calls {
			if mock.Calls[i] != wantCalls[i] {
				t.Errorf("call[%d] = %q, want %q", i, mock.Calls[i], wantCalls[i])
			}
		}
	}
}

func TestSetProfileAppliesHardware(t *testing.T) {
	oldExecutor := executor
	defer func() { executor = oldExecutor }()

	tests := []struct {
		name      string
		profile   settings.AudioProfile
		wantCalls []string
		wantDelta map[string]interface{}
	}{
		{
			name: "profile with sink source volume",
			profile: settings.AudioProfile{
				Name:          "Headphones",
				DefaultSink:   "alsa_output.usb",
				DefaultSource: "alsa_input.usb",
				Volume:        80,
			},
			wantCalls: []string{
				"pactl set-default-sink alsa_output.usb",
				"pactl set-default-source alsa_input.usb",
				"pactl set-sink-volume alsa_output.usb 80%",
			},
			wantDelta: map[string]interface{}{
				"profile":        "Headphones",
				"default_sink":   "alsa_output.usb",
				"default_source": "alsa_input.usb",
				"volume":         float64(80),
			},
		},
		{
			name: "profile with no fields",
			profile: settings.AudioProfile{
				Name: "Empty",
			},
			wantCalls: []string{},
			wantDelta: map[string]interface{}{
				"profile": "Empty",
			},
		},
		{
			name: "clear profile",
			profile: settings.AudioProfile{
				Name: "",
			},
			wantCalls: []string{},
			wantDelta: map[string]interface{}{
				"profile": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &module.MockExecutor{
				Responses: map[string]module.MockResponse{},
			}
			for _, cmd := range tt.wantCalls {
				mock.Responses[cmd] = module.MockResponse{Stdout: "", ExitCode: 0}
			}
			executor = mock

			tmpDir := t.TempDir()
			settingsPath := filepath.Join(tmpDir, "settings.json")
			s := settings.Defaults()
			if tt.profile.Name != "" {
				s.Audio.Profiles = []settings.AudioProfile{tt.profile}
			}
			if err := settings.Save(settingsPath, &s); err != nil {
				t.Fatalf("Save: %v", err)
			}

			value := tt.profile.Name
			resp := captureAudioResponse(t, "set-profile", fmt.Sprintf(`{"value":"%s"}`, value), settingsPath)
			if resp.Status != "ok" {
				t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
			}
			delta, ok := resp.SettingsDelta["audio"].(map[string]interface{})
			if !ok {
				t.Fatalf("expected audio delta")
			}
			for k, v := range tt.wantDelta {
				if delta[k] != v {
					t.Errorf("delta[%s] = %v, want %v", k, delta[k], v)
				}
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
