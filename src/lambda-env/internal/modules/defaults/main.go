package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

var executor module.CLIExecutor = module.NewRealExecutor()
var appsDir = "/usr/share/applications/"

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
	case "set-browser":
		app := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				app = v
			}
		}
		handleSetDefault(settingsPath, "browser", app)
	case "set-terminal":
		app := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				app = v
			}
		}
		handleSetDefault(settingsPath, "terminal", app)
	case "set-editor":
		app := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				app = v
			}
		}
		handleSetDefault(settingsPath, "editor", app)
	case "set-file-manager":
		app := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				app = v
			}
		}
		handleSetDefault(settingsPath, "file_manager", app)
	case "apply":
		handleApply(settingsPath)
	default:
		emitError(action, "unknown action", "use run, set-browser, set-terminal, set-editor, set-file-manager, or apply")
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

func discoverApps(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var apps []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".desktop") {
			apps = append(apps, strings.TrimSuffix(name, ".desktop"))
		}
	}
	sort.Strings(apps)
	return apps, nil
}

func validateDesktopFile(dir, app string) bool {
	path := filepath.Join(dir, app+".desktop")
	_, err := os.Stat(path)
	return err == nil
}

func handleRun(settingsPath string) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("run", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	apps, err := discoverApps(appsDir)
	if err != nil {
		emitError("run", fmt.Sprintf("discover apps: %v", err), "")
		return
	}

	resp := module.Response{
		Status: "ok",
		Action: "run",
		Data: map[string]interface{}{
			"browser":        s.Defaults.Browser,
			"terminal":       s.Defaults.Terminal,
			"editor":         s.Defaults.Editor,
			"file_manager":   s.Defaults.FileManager,
			"available_apps": apps,
			"available_options": map[string]interface{}{
				"set-browser":      apps,
				"set-terminal":     apps,
				"set-editor":       apps,
				"set-file-manager": apps,
			},
			"current_value": map[string]interface{}{
				"set-browser":      s.Defaults.Browser,
				"set-terminal":     s.Defaults.Terminal,
				"set-editor":       s.Defaults.Editor,
				"set-file-manager": s.Defaults.FileManager,
			},
		},
		Message: "Defaults configuration loaded",
	}
	emit(resp)
}

func handleSetDefault(settingsPath, field, app string) {
	if app == "" {
		emitError("set-"+field, fmt.Sprintf("%s cannot be empty", field), "")
		return
	}

	if !validateDesktopFile(appsDir, app) {
		emitError("set-"+field, fmt.Sprintf("%q is not an available application", app), "")
		return
	}

	delta := map[string]interface{}{
		"defaults": map[string]interface{}{
			field: app,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-"+field, fmt.Sprintf("save delta: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-" + field,
		SettingsDelta: delta,
		Message:       fmt.Sprintf("%s set to %s", field, app),
	}
	emit(resp)
}

func handleApply(settingsPath string) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("apply", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	type defaultCmd struct {
		name string
		app  string
		cmd  string
		args []string
	}

	var cmds []defaultCmd
	if s.Defaults.Browser != "" {
		cmds = append(cmds, defaultCmd{
			name: "browser",
			app:  s.Defaults.Browser,
			cmd:  "xdg-settings",
			args: []string{"set", "default-web-browser", s.Defaults.Browser + ".desktop"},
		})
	}
	if s.Defaults.Terminal != "" {
		cmds = append(cmds, defaultCmd{
			name: "terminal",
			app:  s.Defaults.Terminal,
			cmd:  "xdg-mime",
			args: []string{"default", s.Defaults.Terminal + ".desktop", "x-scheme-handler/terminal"},
		})
	}
	if s.Defaults.Editor != "" {
		cmds = append(cmds, defaultCmd{
			name: "editor",
			app:  s.Defaults.Editor,
			cmd:  "xdg-mime",
			args: []string{"default", s.Defaults.Editor + ".desktop", "text/plain"},
		})
	}
	if s.Defaults.FileManager != "" {
		cmds = append(cmds, defaultCmd{
			name: "file_manager",
			app:  s.Defaults.FileManager,
			cmd:  "xdg-mime",
			args: []string{"default", s.Defaults.FileManager + ".desktop", "inode/directory"},
		})
	}

	results := make(map[string]interface{})
	hasError := false
	for _, c := range cmds {
		_, _, exitCode, err := executor.Run(c.cmd, c.args...)
		if exitCode != 0 || err != nil {
			results[c.name] = map[string]interface{}{
				"ok":    false,
				"error": fmt.Sprintf("%v", err),
			}
			hasError = true
		} else {
			results[c.name] = map[string]interface{}{
				"ok": true,
			}
		}
	}

	status := "ok"
	message := "All defaults applied"
	if hasError {
		status = "warning"
		message = "Some defaults could not be applied"
	}

	resp := module.Response{
		Status:  status,
		Action:  "apply",
		Data:    map[string]interface{}{"results": results},
		Message: message,
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
