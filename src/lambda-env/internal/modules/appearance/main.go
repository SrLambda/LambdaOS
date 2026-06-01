package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

var executor module.CLIExecutor = module.NewRealExecutor()

var themeMap = map[string]string{
	"dark":       "tokyonight",
	"light":      "tokyonight-light",
	"nord":       "nord",
	"catppuccin": "catppuccin-mocha",
}

var knownThemes = []string{"dark", "light", "nord", "catppuccin"}

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
	case "set-theme":
		theme := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				theme = v
			}
		}
		handleSetTheme(settingsPath, theme)
	case "set-wallpaper":
		wallpaperPath := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				wallpaperPath = v
			}
		}
		handleSetWallpaper(settingsPath, wallpaperPath)
	case "set-font-size":
		size := 0
		if params != nil {
			switch v := params["value"].(type) {
			case float64:
				size = int(v)
			case int:
				size = v
			case string:
				// ignore string values for font size
			}
		}
		handleSetFontSize(settingsPath, size)
	default:
		emitError(action, "unknown action", "use run, set-theme, set-wallpaper, or set-font-size")
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

	resp := module.Response{
		Status: "ok",
		Action: "run",
		Data: map[string]interface{}{
			"theme":     s.Appearance.Theme,
			"wallpaper": s.Appearance.Wallpaper,
			"font_size": s.Appearance.FontSize,
			"available_options": map[string]interface{}{
				"set-theme": knownThemes,
			},
			"current_value": map[string]interface{}{
				"set-theme":     s.Appearance.Theme,
				"set-wallpaper": s.Appearance.Wallpaper,
				"set-font-size": s.Appearance.FontSize,
			},
		},
		Message: "Appearance configuration loaded",
	}
	emit(resp)
}

func handleSetTheme(settingsPath, theme string) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("set-theme", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	valid := false
	for _, t := range knownThemes {
		if t == theme {
			valid = true
			break
		}
	}
	if !valid {
		emitError("set-theme", fmt.Sprintf("theme %q is not a known theme", theme), "")
		return
	}

	delta := map[string]interface{}{
		"appearance": map[string]interface{}{
			"theme": theme,
		},
	}

	if s.Neovim.UseGlobalTheme {
		if mapped, ok := themeMap[theme]; ok {
			delta["neovim"] = map[string]interface{}{
				"theme": mapped,
			}
		}
	}

	if s.Qtile.UseGlobalTheme {
		if mapped, ok := themeMap[theme]; ok {
			delta["qtile"] = map[string]interface{}{
				"color_scheme": mapped,
			}
		}
	}

	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-theme", fmt.Sprintf("save delta: %v", err), "")
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

func handleSetWallpaper(settingsPath, wallpaperPath string) {
	if wallpaperPath == "" {
		emitError("set-wallpaper", "wallpaper path cannot be empty", "")
		return
	}

	_, _, exitCode, err := executor.Run("feh", "--bg-scale", wallpaperPath)
	if exitCode != 0 || err != nil {
		emitError("set-wallpaper", fmt.Sprintf("feh failed: %v", err), "")
		return
	}

	delta := map[string]interface{}{
		"appearance": map[string]interface{}{
			"wallpaper": wallpaperPath,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-wallpaper", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-wallpaper",
		SettingsDelta: delta,
		Message:       fmt.Sprintf("Wallpaper set to %s", wallpaperPath),
	}
	emit(resp)
}

func handleSetFontSize(settingsPath string, size int) {
	if size < 1 {
		emitError("set-font-size", fmt.Sprintf("font size must be > 0, got %d", size), "")
		return
	}

	delta := map[string]interface{}{
		"appearance": map[string]interface{}{
			"font_size": size,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-font-size", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-font-size",
		SettingsDelta: delta,
		Message:       fmt.Sprintf("Font size set to %d", size),
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
