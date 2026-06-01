package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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

	dotfilesDir := os.Getenv("LAMBDA_ENV_DOTFILES_DIR")
	if dotfilesDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			emitError(action, "cannot determine home directory", "")
			return
		}
		dotfilesDir = filepath.Join(home, "dotfiles")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		emitError(action, "cannot determine home directory", "")
		return
	}

	switch action {
	case "run":
		handleRun(settingsPath)
	case "list":
		handleList(dotfilesDir)
	case "stow":
		handleStow(dotfilesDir)
	case "unstow":
		handleUnstow(dotfilesDir)
	case "check-conflicts":
		handleCheckConflicts(dotfilesDir, homeDir)
	case "backup":
		handleBackup(dotfilesDir, homeDir)
	default:
		emitError(action, "unknown action", "use run, list, stow, unstow, check-conflicts, or backup")
	}
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
			"dotfiles_enabled": s.Services.Enabled,
		},
		Message: "Dotfiles module ready",
	}
	emit(resp)
}

func handleList(dotfilesDir string) {
	modules, err := ListModules(dotfilesDir)
	if err != nil {
		emitError("list", fmt.Sprintf("list modules: %v", err), "")
		return
	}

	var data []map[string]interface{}
	for _, m := range modules {
		data = append(data, map[string]interface{}{
			"name":   m.Name,
			"stowed": m.Stowed,
		})
	}

	resp := module.Response{
		Status: "ok",
		Action: "list",
		Data:   map[string]interface{}{"modules": data},
	}
	emit(resp)
}

func handleStow(dotfilesDir string) {
	name := os.Getenv("LAMBDA_ENV_MODULE")
	if name == "" {
		emitError("stow", "LAMBDA_ENV_MODULE not set", "set the module name to stow")
		return
	}

	if err := Stow(dotfilesDir, name); err != nil {
		emitError("stow", fmt.Sprintf("stow %s: %v", name, err), "")
		return
	}

	resp := module.Response{
		Status: "ok",
		Action: "stow",
		Data:   map[string]interface{}{"module": name},
		Message: fmt.Sprintf("Stowed %s", name),
	}
	emit(resp)
}

func handleUnstow(dotfilesDir string) {
	name := os.Getenv("LAMBDA_ENV_MODULE")
	if name == "" {
		emitError("unstow", "LAMBDA_ENV_MODULE not set", "set the module name to unstow")
		return
	}

	if err := Unstow(dotfilesDir, name); err != nil {
		emitError("unstow", fmt.Sprintf("unstow %s: %v", name, err), "")
		return
	}

	resp := module.Response{
		Status: "ok",
		Action: "unstow",
		Data:   map[string]interface{}{"module": name},
		Message: fmt.Sprintf("Unstowed %s", name),
	}
	emit(resp)
}

func handleCheckConflicts(dotfilesDir, homeDir string) {
	name := os.Getenv("LAMBDA_ENV_MODULE")
	if name == "" {
		emitError("check-conflicts", "LAMBDA_ENV_MODULE not set", "set the module name to check")
		return
	}

	conflicts, err := DetectConflicts(dotfilesDir, name, homeDir)
	if err != nil {
		emitError("check-conflicts", fmt.Sprintf("detect conflicts: %v", err), "")
		return
	}

	data := map[string]interface{}{
		"module":    name,
		"conflicts": conflicts,
		"count":     len(conflicts),
	}

	resp := module.Response{
		Status: "ok",
		Action: "check-conflicts",
		Data:   data,
	}
	emit(resp)
}

func handleBackup(dotfilesDir, homeDir string) {
	name := os.Getenv("LAMBDA_ENV_MODULE")
	if name == "" {
		emitError("backup", "LAMBDA_ENV_MODULE not set", "set the module name to backup")
		return
	}

	count, err := BackupModule(dotfilesDir, name, homeDir)
	if err != nil {
		emitError("backup", fmt.Sprintf("backup %s: %v", name, err), "")
		return
	}

	resp := module.Response{
		Status: "ok",
		Action: "backup",
		Data: map[string]interface{}{
			"module": name,
			"count":  count,
		},
		Message: fmt.Sprintf("Backed up %d files", count),
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

func runStow(dotfilesDir string, args ...string) error {
	cmd := exec.Command("stow", args...)
	cmd.Dir = dotfilesDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("stow %v: %s (%w)", args, string(output), err)
	}
	return nil
}
