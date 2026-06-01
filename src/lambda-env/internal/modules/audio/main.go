package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

var executor module.CLIExecutor = module.NewRealExecutor()

func main() {
	action := os.Getenv("LAMBDA_ENV_ACTION")
	settingsPath := os.Getenv("LAMBDA_ENV_SETTINGS")

	if settingsPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			emitError("run", "cannot determine home directory", "")
			return
		}
		settingsPath = filepath.Join(home, ".config", "lambdaos", "settings.json")
	}

	params := readParams()

	switch action {
	case "run":
		handleRun(settingsPath)
	case "set-volume":
		volume := 0
		if params != nil {
			switch v := params["value"].(type) {
			case float64:
				volume = int(v)
			case int:
				volume = v
			case string:
				if parsed, err := strconv.Atoi(v); err == nil {
					volume = parsed
				}
			}
		}
		handleSetVolume(settingsPath, volume)
	case "set-mute":
		muted := false
		if params != nil {
			if v, ok := params["value"].(bool); ok {
				muted = v
			}
		}
		handleSetMute(settingsPath, muted)
	case "set-sink":
		sink := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				sink = v
			}
		}
		handleSetSink(settingsPath, sink)
	default:
		emitError(action, "unknown action", "use run, set-volume, set-mute, or set-sink")
	}
}

func readParams() map[string]interface{} {
	p := os.Getenv("LAMBDA_ENV_PARAMS")
	if p == "" {
		return nil
	}
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(p), &params); err != nil {
		return nil
	}
	return params
}

func detectBackend(exe module.CLIExecutor) (string, error) {
	stdout, _, exitCode, err := exe.Run("pactl", "info")
	if exitCode != 0 || err != nil {
		return "", fmt.Errorf("pactl info failed: %v", err)
	}
	for _, line := range strings.Split(stdout, "\n") {
		if strings.Contains(strings.ToLower(line), "pipewire") {
			return "pipewire", nil
		}
	}
	return "pulseaudio", nil
}

func discoverSinks(exe module.CLIExecutor) ([]string, error) {
	stdout, _, exitCode, err := exe.Run("pactl", "list", "short", "sinks")
	if exitCode != 0 || err != nil {
		return nil, fmt.Errorf("pactl list short sinks failed: %v", err)
	}
	var sinks []string
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) >= 2 {
			sinks = append(sinks, fields[1])
		}
	}
	return sinks, nil
}

func handleRun(settingsPath string) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("run", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	backend, err := detectBackend(executor)
	if err != nil {
		emitError("run", fmt.Sprintf("detect backend: %v", err), "")
		return
	}

	sinks, err := discoverSinks(executor)
	if err != nil {
		emitError("run", fmt.Sprintf("discover sinks: %v", err), "")
		return
	}

	resp := module.Response{
		Status: "ok",
		Action: "run",
		Data: map[string]interface{}{
			"volume":       s.Audio.Volume,
			"muted":        s.Audio.Muted,
			"default_sink": s.Audio.DefaultSink,
			"backend":      backend,
			"available_options": map[string]interface{}{
				"set-sink": sinks,
			},
			"current_value": map[string]interface{}{
				"set-volume": s.Audio.Volume,
				"set-mute":   s.Audio.Muted,
				"set-sink":   s.Audio.DefaultSink,
			},
		},
		Message: "Audio configuration loaded",
	}
	emit(resp)
}

func handleSetVolume(settingsPath string, volume int) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("set-volume", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	if volume < 0 {
		volume = 0
	}
	if volume > 100 {
		volume = 100
	}

	sink := s.Audio.DefaultSink
	if sink == "" {
		emitError("set-volume", "no default sink configured", "")
		return
	}

	_, _, exitCode, err := executor.Run("pactl", "set-sink-volume", sink, fmt.Sprintf("%d%%", volume))
	if exitCode != 0 || err != nil {
		emitError("set-volume", fmt.Sprintf("pactl failed: %v", err), "")
		return
	}

	delta := map[string]interface{}{
		"audio": map[string]interface{}{
			"volume": volume,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-volume", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-volume",
		SettingsDelta: delta,
		Message:       fmt.Sprintf("Volume set to %d%%", volume),
	}
	emit(resp)
}

func handleSetMute(settingsPath string, muted bool) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("set-mute", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	sink := s.Audio.DefaultSink
	if sink == "" {
		emitError("set-mute", "no default sink configured", "")
		return
	}

	muteVal := "0"
	if muted {
		muteVal = "1"
	}

	_, _, exitCode, err := executor.Run("pactl", "set-sink-mute", sink, muteVal)
	if exitCode != 0 || err != nil {
		emitError("set-mute", fmt.Sprintf("pactl failed: %v", err), "")
		return
	}

	delta := map[string]interface{}{
		"audio": map[string]interface{}{
			"muted": muted,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-mute", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-mute",
		SettingsDelta: delta,
		Message:       fmt.Sprintf("Mute set to %v", muted),
	}
	emit(resp)
}

func handleSetSink(settingsPath, sink string) {
	if _, err := settings.Load(settingsPath); err != nil {
		emitError("set-sink", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	sinks, err := discoverSinks(executor)
	if err != nil {
		emitError("set-sink", fmt.Sprintf("discover sinks: %v", err), "")
		return
	}

	valid := false
	for _, sk := range sinks {
		if sk == sink {
			valid = true
			break
		}
	}
	if !valid {
		emitError("set-sink", fmt.Sprintf("sink %q is not available", sink), "")
		return
	}

	_, _, exitCode, err := executor.Run("pactl", "set-default-sink", sink)
	if exitCode != 0 || err != nil {
		emitError("set-sink", fmt.Sprintf("pactl failed: %v", err), "")
		return
	}

	delta := map[string]interface{}{
		"audio": map[string]interface{}{
			"default_sink": sink,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-sink", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-sink",
		SettingsDelta: delta,
		Message:       fmt.Sprintf("Default sink set to %s", sink),
	}
	emit(resp)
}

func emit(resp module.Response) {
	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal response: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func emitError(action, message, suggestion string) {
	resp := module.Response{
		Status:     "error",
		Action:     action,
		Message:    message,
		Suggestion: suggestion,
	}
	emit(resp)
}
