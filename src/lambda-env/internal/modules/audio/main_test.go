package main

import (
	"encoding/json"
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
