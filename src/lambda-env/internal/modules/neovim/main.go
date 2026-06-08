package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
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

	params := readParams()

	switch action {
	case "run":
		handleRun(settingsPath)
	case "toggle-lsp":
		handleToggle(settingsPath, "toggle-lsp", "enable_lsp", "LSP")
	case "toggle-copilot":
		handleToggle(settingsPath, "toggle-copilot", "enable_copilot", "Copilot")
	case "toggle-neotree":
		handleToggle(settingsPath, "toggle-neotree", "enable_neotree", "Neo-tree")
	case "set-theme":
		theme := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				theme = v
			}
		}
		handleSetTheme(settingsPath, theme)
	case "apply":
		handleApply(settingsPath)
	default:
		emitError(action, "unknown action", "use run, toggle-lsp, toggle-copilot, toggle-neotree, set-theme, or apply")
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

func handleRun(settingsPath string) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("run", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	if err := Apply(settingsPath); err != nil {
		emitError("run", fmt.Sprintf("apply neovim config: %v", err), "")
		return
	}

	resp := module.Response{
		Status: "ok",
		Action: "run",
		Data: map[string]interface{}{
			"enable_lsp":     s.Neovim.EnableLSP,
			"enable_copilot": s.Neovim.EnableCopilot,
			"enable_neotree": s.Neovim.EnableNeotree,
		},
		Message: "Neovim configuration applied",
	}
	emit(resp)
}

func handleSetTheme(settingsPath, theme string) {
	if theme == "" {
		emitError("set-theme", "theme value is required", "")
		return
	}

	_, err := settings.Load(settingsPath)
	if err != nil {
		emitError("set-theme", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	delta := map[string]interface{}{
		"neovim": map[string]interface{}{
			"theme": theme,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-theme", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	if err := Apply(settingsPath); err != nil {
		emitError("set-theme", fmt.Sprintf("apply neovim config: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-theme",
		SettingsDelta: delta,
		Message:       fmt.Sprintf("Theme set to %s", theme),
	}
	emit(resp)
}

func handleApply(settingsPath string) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("apply", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	if err := Apply(settingsPath); err != nil {
		emitError("apply", fmt.Sprintf("apply neovim config: %v", err), "")
		return
	}

	resp := module.Response{
		Status:  "ok",
		Action:  "apply",
		Message: "Neovim configuration applied",
		Data: map[string]interface{}{
			"theme":          s.Neovim.Theme,
			"enable_lsp":     s.Neovim.EnableLSP,
			"enable_copilot": s.Neovim.EnableCopilot,
			"enable_neotree": s.Neovim.EnableNeotree,
		},
	}
	emit(resp)
}

func handleToggle(settingsPath, actionName, field, label string) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError(actionName, fmt.Sprintf("load settings: %v", err), "")
		return
	}

	var currentValue bool
	switch field {
	case "enable_lsp":
		currentValue = s.Neovim.EnableLSP
	case "enable_copilot":
		currentValue = s.Neovim.EnableCopilot
	case "enable_neotree":
		currentValue = s.Neovim.EnableNeotree
	}

	newValue := !currentValue

	delta := map[string]interface{}{
		"neovim": map[string]interface{}{
			field: newValue,
		},
	}

	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError(actionName, fmt.Sprintf("save delta: %v", err), "")
		return
	}

	if err := Apply(settingsPath); err != nil {
		emitError(actionName, fmt.Sprintf("apply neovim config: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        actionName,
		SettingsDelta: delta,
		Message:       fmt.Sprintf("%s %s", label, toggleLabel(newValue)),
	}
	emit(resp)
}

func toggleLabel(on bool) string {
	if on {
		return "enabled"
	}
	return "disabled"
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
