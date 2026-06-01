package main

import (
	"path/filepath"
	"testing"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

func TestIntegrationAudioFullFlow(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"pactl info": {
				Stdout:   "Server Name: pulseaudio\nDefault Sink: sink1\n",
				ExitCode: 0,
			},
			"pactl list short sinks": {
				Stdout:   "0\tsink1\tmodule-alsa-card.c\ts16le 2ch 44100Hz\tRUNNING\n1\tsink2\tmodule-bluez5-device.c\ts16le 2ch 44100Hz\tIDLE\n",
				ExitCode: 0,
			},
			"pactl set-sink-volume sink1 30%": {
				Stdout:   "",
				ExitCode: 0,
			},
			"pactl set-sink-mute sink1 1": {
				Stdout:   "",
				ExitCode: 0,
			},
			"pactl set-default-sink sink2": {
				Stdout:   "",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Audio.DefaultSink = "sink1"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Run action returns dynamic sink options
	resp := captureAudioResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("run expected ok, got %s", resp.Status)
	}
	opts, ok := resp.Data["available_options"].(map[string]interface{})
	if !ok || len(opts["set-sink"].([]interface{})) != 2 {
		t.Fatalf("expected 2 sinks in run response")
	}

	// Set volume
	resp = captureAudioResponse(t, "set-volume", `{"value":30}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("set-volume expected ok, got %s", resp.Status)
	}
	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Audio.Volume != 30 {
		t.Errorf("volume = %d, want 30", loaded.Audio.Volume)
	}

	// Set mute
	resp = captureAudioResponse(t, "set-mute", `{"value":true}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("set-mute expected ok, got %s", resp.Status)
	}
	loaded, err = settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !loaded.Audio.Muted {
		t.Errorf("muted = %v, want true", loaded.Audio.Muted)
	}

	// Set sink
	resp = captureAudioResponse(t, "set-sink", `{"value":"sink2"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("set-sink expected ok, got %s", resp.Status)
	}
	loaded, err = settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Audio.DefaultSink != "sink2" {
		t.Errorf("default_sink = %q, want sink2", loaded.Audio.DefaultSink)
	}
}
