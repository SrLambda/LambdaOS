package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	case "set-layout":
		layout := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				layout = v
			}
		}
		handleSetLayout(settingsPath, layout)
	case "set-variant":
		variant := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				variant = v
			}
		}
		handleSetVariant(settingsPath, variant)
	default:
		emitError(action, "unknown action", "use run, set-layout, or set-variant")
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

func discoverLayouts(exe module.CLIExecutor) ([]string, error) {
	stdout, _, exitCode, err := exe.Run("setxkbmap", "-layout")
	if exitCode != 0 || err != nil {
		return []string{"us", "es", "de", "fr", "gb", "it", "jp", "ru"}, nil
	}
	var layouts []string
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			layouts = append(layouts, line)
		}
	}
	return layouts, nil
}

func handleRun(settingsPath string) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("run", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	layouts, err := discoverLayouts(executor)
	if err != nil {
		emitError("run", fmt.Sprintf("discover layouts: %v", err), "")
		return
	}

	resp := module.Response{
		Status: "ok",
		Action: "run",
		Data: map[string]interface{}{
			"layout":            s.Keyboard.Layout,
			"variant":           s.Keyboard.Variant,
			"available_options": map[string]interface{}{"set-layout": layouts},
			"current_value": map[string]interface{}{
				"set-layout":  s.Keyboard.Layout,
				"set-variant": s.Keyboard.Variant,
			},
		},
		Message: "Keyboard configuration loaded",
	}
	emit(resp)
}

func handleSetLayout(settingsPath, layout string) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("set-layout", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	layouts, err := discoverLayouts(executor)
	if err != nil {
		emitError("set-layout", fmt.Sprintf("discover layouts: %v", err), "")
		return
	}

	valid := false
	for _, l := range layouts {
		if l == layout {
			valid = true
			break
		}
	}
	if !valid {
		emitError("set-layout", fmt.Sprintf("layout %q is not available", layout), "")
		return
	}

	args := []string{"-layout", layout}
	if s.Keyboard.Variant != "" {
		args = append(args, "-variant", s.Keyboard.Variant)
	}
	_, _, exitCode, err := executor.Run("setxkbmap", args...)
	if exitCode != 0 || err != nil {
		emitError("set-layout", fmt.Sprintf("setxkbmap failed: %v", err), "")
		return
	}

	delta := map[string]interface{}{
		"keyboard": map[string]interface{}{
			"layout": layout,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-layout", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-layout",
		SettingsDelta: delta,
		Message:       fmt.Sprintf("Layout set to %s", layout),
	}
	emit(resp)
}

func handleSetVariant(settingsPath, variant string) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("set-variant", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	args := []string{"-layout", s.Keyboard.Layout}
	if variant != "" {
		args = append(args, "-variant", variant)
	}
	_, _, exitCode, err := executor.Run("setxkbmap", args...)
	if exitCode != 0 || err != nil {
		emitError("set-variant", fmt.Sprintf("setxkbmap failed: %v", err), "")
		return
	}

	delta := map[string]interface{}{
		"keyboard": map[string]interface{}{
			"variant": variant,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-variant", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-variant",
		SettingsDelta: delta,
		Message:       fmt.Sprintf("Variant set to %s", variant),
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
