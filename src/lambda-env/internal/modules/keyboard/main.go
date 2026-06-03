package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

var executor module.CLIExecutor = module.NewRealExecutor()

// curatedComposeOptions maps human labels to XKB option strings.
var curatedComposeOptions = []string{
	"compose:ralt",
	"compose:lwin",
	"compose:menu",
	"compose:caps",
}

// curatedXKBOptions maps human labels to XKB option strings.
var curatedXKBOptions = []string{
	"caps:swapescape",
	"ctrl:nocaps",
	"compose:ralt",
	"compose:lwin",
	"compose:menu",
	"compose:caps",
	"terminate:ctrl_alt_bksp",
	"altwin:swap_alt_win",
}

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
	case "set-compose":
		compose := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				compose = v
			}
		}
		handleSetCompose(settingsPath, compose)
	case "set-options":
		options := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				options = v
			}
		}
		handleSetOptions(settingsPath, options)
	default:
		emitError(action, "unknown action", "use run, set-layout, set-variant, set-compose, or set-options")
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
	stdout, _, exitCode, err := exe.Run("localectl", "list-x11-keymap-layouts")
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

func discoverVariants(exe module.CLIExecutor, layout string) ([]string, error) {
	stdout, _, exitCode, err := exe.Run("localectl", "list-x11-keymap-variants", layout)
	if exitCode != 0 || err != nil {
		return []string{""}, nil
	}
	var variants []string
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			variants = append(variants, line)
		}
	}
	// Always include empty variant for "no variant".
	return append([]string{""}, variants...), nil
}

// parseCurrentLayout queries setxkbmap and returns layout, variant, options.
func parseCurrentLayout(exe module.CLIExecutor) (layout, variant, options string, err error) {
	stdout, _, exitCode, execErr := exe.Run("setxkbmap", "-query")
	if exitCode != 0 || execErr != nil {
		return "", "", "", fmt.Errorf("setxkbmap -query failed: %v", execErr)
	}

	layoutRe := regexp.MustCompile(`(?m)^layout:\s*(\S+)`)
	variantRe := regexp.MustCompile(`(?m)^variant:\s*(\S*)`)
	optionsRe := regexp.MustCompile(`(?m)^options:\s*(\S*)`)

	if m := layoutRe.FindStringSubmatch(stdout); len(m) > 1 {
		layout = m[1]
	}
	if m := variantRe.FindStringSubmatch(stdout); len(m) > 1 {
		variant = m[1]
	}
	if m := optionsRe.FindStringSubmatch(stdout); len(m) > 1 {
		options = m[1]
	}
	return layout, variant, options, nil
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

	variants, err := discoverVariants(executor, s.Keyboard.Layout)
	if err != nil {
		emitError("run", fmt.Sprintf("discover variants: %v", err), "")
		return
	}

	_, _, currentOptions, err := parseCurrentLayout(executor)
	if err != nil {
		currentOptions = s.Keyboard.Options
	}

	resp := module.Response{
		Status: "ok",
		Action: "run",
		Data: map[string]interface{}{
			"layout":  s.Keyboard.Layout,
			"variant": s.Keyboard.Variant,
			"options": currentOptions,
			"available_options": map[string]interface{}{
				"set-layout":  layouts,
				"set-variant": variants,
				"set-compose": curatedComposeOptions,
				"set-options": curatedXKBOptions,
			},
			"current_value": map[string]interface{}{
				"set-layout":  s.Keyboard.Layout,
				"set-variant": s.Keyboard.Variant,
				"set-options": s.Keyboard.Options,
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
	if s.Keyboard.Options != "" {
		args = append(args, "-option", s.Keyboard.Options)
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
	if s.Keyboard.Options != "" {
		args = append(args, "-option", s.Keyboard.Options)
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

func handleSetCompose(settingsPath, compose string) {
	if compose == "" {
		emitError("set-compose", "compose option is required", "")
		return
	}

	valid := false
	for _, opt := range curatedComposeOptions {
		if opt == compose {
			valid = true
			break
		}
	}
	if !valid {
		emitError("set-compose", fmt.Sprintf("compose option %q is not available", compose), "")
		return
	}

	_, _, exitCode, err := executor.Run("setxkbmap", "-option", compose)
	if exitCode != 0 || err != nil {
		emitError("set-compose", fmt.Sprintf("setxkbmap failed: %v", err), "")
		return
	}

	delta := map[string]interface{}{
		"keyboard": map[string]interface{}{
			"options": compose,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-compose", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-compose",
		SettingsDelta: delta,
		Message:       fmt.Sprintf("Compose key set to %s", compose),
	}
	emit(resp)
}

func handleSetOptions(settingsPath, options string) {
	if options == "" {
		emitError("set-options", "options value is required", "")
		return
	}

	// Validate each option in the comma-separated list.
	selected := strings.Split(options, ",")
	for _, sel := range selected {
		sel = strings.TrimSpace(sel)
		if sel == "" {
			continue
		}
		valid := false
		for _, opt := range curatedXKBOptions {
			if opt == sel {
				valid = true
				break
			}
		}
		if !valid {
			emitError("set-options", fmt.Sprintf("option %q is not available", sel), "")
			return
		}
	}

	// Clear previous options first.
	_, _, exitCode, err := executor.Run("setxkbmap", "-option", "")
	if exitCode != 0 || err != nil {
		emitError("set-options", fmt.Sprintf("clearing options failed: %v", err), "")
		return
	}

	// Apply new options.
	_, _, exitCode, err = executor.Run("setxkbmap", "-option", options)
	if exitCode != 0 || err != nil {
		emitError("set-options", fmt.Sprintf("setxkbmap failed: %v", err), "")
		return
	}

	delta := map[string]interface{}{
		"keyboard": map[string]interface{}{
			"options": options,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-options", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-options",
		SettingsDelta: delta,
		Message:       fmt.Sprintf("Options set to %s", options),
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
