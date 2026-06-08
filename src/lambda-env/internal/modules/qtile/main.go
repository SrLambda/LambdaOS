package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

var (
	terminalAllowlist = map[string]bool{
		"kitty":     true,
		"foot":      true,
		"alacritty": true,
		"st":        true,
		"xterm":     true,
	}

	browserAllowlist = map[string]bool{
		"firefox":  true,
		"chromium": true,
		"brave":    true,
		"chrome":   true,
	}

	fileManagerAllowlist = map[string]bool{
		"thunar":   true,
		"yazi":     true,
		"nemo":     true,
		"nautilus": true,
		"ranger":   true,
	}
)

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

	switch action {
	case "run":
		handleRun(settingsPath)
	case "set-terminal":
		handleSet(settingsPath, "set-terminal", "terminal", terminalAllowlist)
	case "set-browser":
		handleSet(settingsPath, "set-browser", "browser", browserAllowlist)
	case "set-file-manager":
		handleSet(settingsPath, "set-file-manager", "default_file_manager", fileManagerAllowlist)
	case "reload":
		handleReload(settingsPath)
	default:
		emitError(action, "unknown action", "use run, set-terminal, set-browser, set-file-manager, or reload")
	}
}

func handleRun(settingsPath string) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("run", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	if err := Apply(settingsPath); err != nil {
		emitError("run", fmt.Sprintf("apply qtile config: %v", err), "")
		return
	}

	resp := module.Response{
		Status: "ok",
		Action: "run",
		Data: map[string]interface{}{
			"terminal":             s.Qtile.Terminal,
			"browser":              s.Qtile.Browser,
			"default_file_manager": s.Qtile.DefaultFileManager,
		},
		Message: "Qtile configuration applied",
	}
	emit(resp)
}

func handleSet(settingsPath, actionName, field string, allowlist map[string]bool) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError(actionName, fmt.Sprintf("load settings: %v", err), "")
		return
	}

	var newValue string
	switch field {
	case "terminal":
		newValue = s.Qtile.Terminal
	case "browser":
		newValue = s.Qtile.Browser
	case "default_file_manager":
		newValue = s.Qtile.DefaultFileManager
	}

	if !allowlist[newValue] {
		emitError(actionName, fmt.Sprintf("value %q is not in allowlist", newValue), "")
		return
	}

	delta := map[string]interface{}{
		"qtile": map[string]interface{}{
			field: newValue,
		},
	}

	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError(actionName, fmt.Sprintf("save delta: %v", err), "")
		return
	}

	if err := Apply(settingsPath); err != nil {
		emitError(actionName, fmt.Sprintf("apply qtile config: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        actionName,
		SettingsDelta: delta,
		Message:       fmt.Sprintf("%s set to %s", field, newValue),
	}
	emit(resp)
}

func handleReload(settingsPath string) {
	if err := SafeApply(settingsPath); err != nil {
		emitError("reload", fmt.Sprintf("safe apply: %v", err), "")
		return
	}

	resp := module.Response{
		Status:  "ok",
		Action:  "reload",
		Message: "Qtile configuration reloaded",
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
