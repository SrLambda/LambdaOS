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
		valid := false
		if params != nil {
			switch v := params["value"].(type) {
			case float64:
				volume = int(v)
				valid = true
			case int:
				volume = v
				valid = true
			case string:
				if parsed, err := strconv.Atoi(v); err == nil {
					volume = parsed
					valid = true
				}
			}
		}
		if !valid {
			emitError("set-volume", "invalid volume value", "expected a number between 0 and 100")
			return
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
	case "set-source":
		source := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				source = v
			}
		}
		handleSetSource(settingsPath, source)
	case "set-profile":
		profile := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				profile = v
			}
		}
		handleSetProfile(settingsPath, profile)
	case "set-app-volume":
		appVolume := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				appVolume = v
			}
		}
		handleSetAppVolume(settingsPath, appVolume)
	default:
		emitError(action, "unknown action", "use run, set-volume, set-mute, set-sink, set-source, set-profile, or set-app-volume")
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

func discoverSources(exe module.CLIExecutor) ([]string, error) {
	stdout, _, exitCode, err := exe.Run("pactl", "list", "short", "sources")
	if exitCode != 0 || err != nil {
		return nil, fmt.Errorf("pactl list short sources failed: %v", err)
	}
	var sources []string
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) >= 2 {
			sources = append(sources, fields[1])
		}
	}
	return sources, nil
}

// perAppVolume represents a single application stream volume entry.
type perAppVolume struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Volume float64 `json:"volume"`
}

// discoverPerAppVolumes attempts to list per-application volumes using wpctl.
// This is best-effort and only works when PipeWire is the backend.
func parseWpctlVolume(stdout string) float64 {
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Volume:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				v, err := strconv.ParseFloat(strings.TrimSuffix(parts[1], ","), 64)
				if err == nil {
					return v
				}
			}
		}
	}
	return 0
}

func discoverPerAppVolumes(exe module.CLIExecutor) ([]perAppVolume, string, error) {
	stdout, _, exitCode, err := exe.Run("wpctl", "status")
	if exitCode != 0 || err != nil {
		return nil, "wpctl not available or PipeWire not running", nil
	}

	var apps []perAppVolume
	inStreams := false
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Streams:") {
			inStreams = true
			continue
		}
		if inStreams && strings.HasSuffix(line, ":") {
			// Next section header ends the streams block.
			break
		}
		if inStreams && line != "" {
			// Try to parse lines like: "  55. firefox       [Stream/Output/Audio]"
			// We extract the ID and name.
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				idStr := strings.TrimSuffix(parts[0], ".")
				name := parts[1]
				if idStr != "" && name != "" {
					volume := 0.0
					inspectStdout, _, inspectExitCode, inspectErr := exe.Run("wpctl", "inspect", idStr)
					if inspectExitCode == 0 && inspectErr == nil {
						volume = parseWpctlVolume(inspectStdout)
					}
					apps = append(apps, perAppVolume{
						ID:     idStr,
						Name:   name,
						Volume: volume,
					})
				}
			}
		}
	}
	return apps, "", nil
}

func getProfileNames(profiles []settings.AudioProfile) []string {
	names := make([]string, 0, len(profiles))
	for _, p := range profiles {
		if p.Name != "" {
			names = append(names, p.Name)
		}
	}
	return names
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

	sources, err := discoverSources(executor)
	if err != nil {
		emitError("run", fmt.Sprintf("discover sources: %v", err), "")
		return
	}

	var perApp []perAppVolume
	var perAppWarning string
	if backend == "pipewire" {
		perApp, perAppWarning, _ = discoverPerAppVolumes(executor)
	} else {
		perAppWarning = "per-app volume requires PipeWire"
	}

	profileNames := getProfileNames(s.Audio.Profiles)

	data := map[string]interface{}{
		"volume":          s.Audio.Volume,
		"muted":           s.Audio.Muted,
		"default_sink":    s.Audio.DefaultSink,
		"default_source":  s.Audio.DefaultSource,
		"backend":         backend,
		"profile":         s.Audio.Profile,
		"profiles":        profileNames,
		"per_app_volumes": perApp,
		"available_options": map[string]interface{}{
			"set-sink":    sinks,
			"set-source":  sources,
			"set-profile": profileNames,
		},
		"current_value": map[string]interface{}{
			"set-volume":  s.Audio.Volume,
			"set-mute":    s.Audio.Muted,
			"set-sink":    s.Audio.DefaultSink,
			"set-source":  s.Audio.DefaultSource,
			"set-profile": s.Audio.Profile,
		},
	}
	if perAppWarning != "" {
		data["per_app_warning"] = perAppWarning
	}

	resp := module.Response{
		Status:  "ok",
		Action:  "run",
		Data:    data,
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

func handleSetSource(settingsPath, source string) {
	if _, err := settings.Load(settingsPath); err != nil {
		emitError("set-source", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	sources, err := discoverSources(executor)
	if err != nil {
		emitError("set-source", fmt.Sprintf("discover sources: %v", err), "")
		return
	}

	valid := false
	for _, src := range sources {
		if src == source {
			valid = true
			break
		}
	}
	if !valid {
		emitError("set-source", fmt.Sprintf("source %q is not available", source), "")
		return
	}

	_, _, exitCode, err := executor.Run("pactl", "set-default-source", source)
	if exitCode != 0 || err != nil {
		emitError("set-source", fmt.Sprintf("pactl failed: %v", err), "")
		return
	}

	delta := map[string]interface{}{
		"audio": map[string]interface{}{
			"default_source": source,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-source", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-source",
		SettingsDelta: delta,
		Message:       fmt.Sprintf("Default source set to %s", source),
	}
	emit(resp)
}

func handleSetProfile(settingsPath, profile string) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("set-profile", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	var matchedProfile *settings.AudioProfile
	if profile != "" {
		found := false
		for i, p := range s.Audio.Profiles {
			if p.Name == profile {
				found = true
				matchedProfile = &s.Audio.Profiles[i]
				break
			}
		}
		if !found {
			emitError("set-profile", fmt.Sprintf("profile %q does not exist", profile), "")
			return
		}
	}

	audioDelta := map[string]interface{}{
		"profile": profile,
	}

	if matchedProfile != nil {
		if matchedProfile.DefaultSink != "" {
			_, _, exitCode, err := executor.Run("pactl", "set-default-sink", matchedProfile.DefaultSink)
			if exitCode != 0 || err != nil {
				emitError("set-profile", fmt.Sprintf("pactl set-default-sink failed: %v", err), "")
				return
			}
			audioDelta["default_sink"] = matchedProfile.DefaultSink
		}
		if matchedProfile.DefaultSource != "" {
			_, _, exitCode, err := executor.Run("pactl", "set-default-source", matchedProfile.DefaultSource)
			if exitCode != 0 || err != nil {
				emitError("set-profile", fmt.Sprintf("pactl set-default-source failed: %v", err), "")
				return
			}
			audioDelta["default_source"] = matchedProfile.DefaultSource
		}
		if matchedProfile.Volume > 0 {
			sink := s.Audio.DefaultSink
			if matchedProfile.DefaultSink != "" {
				sink = matchedProfile.DefaultSink
			}
			if sink == "" {
				emitError("set-profile", "no default sink configured to apply volume", "")
				return
			}
			_, _, exitCode, err := executor.Run("pactl", "set-sink-volume", sink, fmt.Sprintf("%d%%", matchedProfile.Volume))
			if exitCode != 0 || err != nil {
				emitError("set-profile", fmt.Sprintf("pactl set-sink-volume failed: %v", err), "")
				return
			}
			audioDelta["volume"] = matchedProfile.Volume
		}
	}

	delta := map[string]interface{}{
		"audio": audioDelta,
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-profile", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-profile",
		SettingsDelta: delta,
		Message:       fmt.Sprintf("Profile set to %s", profile),
	}
	emit(resp)
}

func handleSetAppVolume(settingsPath, appVolume string) {
	if appVolume == "" {
		emitError("set-app-volume", "app volume value is required (format: id:volume)", "")
		return
	}

	parts := strings.SplitN(appVolume, ":", 2)
	if len(parts) != 2 {
		emitError("set-app-volume", "invalid app volume format", "expected id:volume")
		return
	}
	streamID := parts[0]
	volumeStr := parts[1]

	// Validate volume is a valid number.
	volume, err := strconv.ParseFloat(volumeStr, 64)
	if err != nil {
		emitError("set-app-volume", fmt.Sprintf("invalid volume: %v", err), "")
		return
	}
	if volume < 0 {
		volume = 0
	}
	if volume > 1.5 {
		volume = 1.5
	}

	_, _, exitCode, err := executor.Run("wpctl", "set-volume", streamID, volumeStr)
	if exitCode != 0 || err != nil {
		emitError("set-app-volume", fmt.Sprintf("wpctl failed: %v", err), "")
		return
	}

	// No settings delta for per-app volume (runtime-only).
	resp := module.Response{
		Status:  "ok",
		Action:  "set-app-volume",
		Message: fmt.Sprintf("App volume for %s set to %s", streamID, volumeStr),
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
