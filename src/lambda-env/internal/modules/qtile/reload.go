package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"lambdaos.dev/lambda-env/internal/settings"
)

func ValidateConfigPy(path string) error {
	cmd := exec.Command("python3", "-m", "py_compile", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("py_compile failed: %s: %w", strings.TrimSpace(string(output)), err)
	}
	return nil
}

func ReloadQtile() error {
	cmd := exec.Command("qtile", "cmd-obj", "-o", "cmd", "-f", "reload_config")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "not found") || strings.Contains(string(output), "no such file") {
			return fmt.Errorf("qtile reload command not found (is Qtile running?): %w", err)
		}
		return fmt.Errorf("qtile reload failed: %s: %w", strings.TrimSpace(string(output)), err)
	}
	return nil
}

func SafeApply(settingsPath string) error {
	s, err := settings.Load(settingsPath)
	if err != nil {
		return fmt.Errorf("load settings: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}

	configPath := filepath.Join(home, ".config", "qtile", "config.py")
	backupPath := configPath + ".bak"
	keysPath := filepath.Join(home, ".config", "qtile", "keys.py")
	keysBackupPath := keysPath + ".bak"

	if _, err := os.Stat(configPath); err == nil {
		data, readErr := os.ReadFile(configPath)
		if readErr != nil {
			return fmt.Errorf("read existing config.py: %w", readErr)
		}
		if writeErr := os.WriteFile(backupPath, data, 0644); writeErr != nil {
			return fmt.Errorf("backup config.py: %w", writeErr)
		}
	}

	if _, err := os.Stat(keysPath); err == nil {
		data, readErr := os.ReadFile(keysPath)
		if readErr != nil {
			return fmt.Errorf("read existing keys.py: %w", readErr)
		}
		if writeErr := os.WriteFile(keysBackupPath, data, 0644); writeErr != nil {
			return fmt.Errorf("backup keys.py: %w", writeErr)
		}
	}

	generated, err := GenerateConfigPy(s.Qtile)
	if err != nil {
		return fmt.Errorf("generate config.py: %w", err)
	}

	if generated == "" {
		return fmt.Errorf("generated config.py content is empty")
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("create config.py directory: %w", err)
	}

	tmpPath := configPath + ".lambda-env.tmp"
	if err := os.WriteFile(tmpPath, []byte(generated), 0644); err != nil {
		return fmt.Errorf("write temp config.py: %w", err)
	}

	if err := ValidateConfigPy(tmpPath); err != nil {
		os.Remove(tmpPath)
		if _, statErr := os.Stat(backupPath); statErr == nil {
			restoreErr := os.Rename(backupPath, configPath)
			if restoreErr != nil {
				return fmt.Errorf("validation failed and restore also failed: %v: original error: %w", restoreErr, err)
			}
		}
		return fmt.Errorf("generated config.py failed validation: %w", err)
	}

	if err := os.Rename(tmpPath, configPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename temp config.py: %w", err)
	}

	if err := updateKeysPyTerminal(keysPath, s.Qtile.Terminal); err != nil {
		if _, statErr := os.Stat(keysBackupPath); statErr == nil {
			os.Rename(keysBackupPath, keysPath)
		}
		return fmt.Errorf("update keys.py terminal: %w", err)
	}

	if err := ReloadQtile(); err != nil {
		if _, statErr := os.Stat(backupPath); statErr == nil {
			restoreErr := os.Rename(backupPath, configPath)
			if restoreErr != nil {
				return fmt.Errorf("reload failed and restore also failed: %v: original error: %w", restoreErr, err)
			}
		}
		if _, statErr := os.Stat(keysBackupPath); statErr == nil {
			os.Rename(keysBackupPath, keysPath)
		}
		return fmt.Errorf("qtile reload failed: %w", err)
	}

	return nil
}
