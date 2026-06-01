package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"lambdaos.dev/lambda-env/internal/modules/qtile/templates"
	"lambdaos.dev/lambda-env/internal/settings"
)

var qtileThemeMap = map[string]string{
	"dark":       "tokyonight",
	"light":      "tokyonight-light",
	"nord":       "nord",
	"catppuccin": "catppuccin-mocha",
}

func resolveQtileColorScheme(s settings.Settings) string {
	if s.Qtile.UseGlobalTheme {
		if t, ok := qtileThemeMap[s.Appearance.Theme]; ok {
			return t
		}
	}
	return s.Qtile.ColorScheme
}

func GenerateConfigPy(s settings.QtileSettings) (string, error) {
	tmpl, err := template.New("config.py").Parse(templates.ConfigPyTemplate)
	if err != nil {
		return "", fmt.Errorf("parse config.py template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, s); err != nil {
		return "", fmt.Errorf("execute config.py template: %w", err)
	}

	return buf.String(), nil
}

func Apply(settingsPath string) error {
	s, err := settings.Load(settingsPath)
	if err != nil {
		return fmt.Errorf("load settings: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}

	configPath := filepath.Join(home, ".config", "qtile", "config.py")
	keysPath := filepath.Join(home, ".config", "qtile", "keys.py")
	backupPath := configPath + ".bak"
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

	s.Qtile.ColorScheme = resolveQtileColorScheme(*s)

	generated, err := GenerateConfigPy(s.Qtile)
	if err != nil {
		return err
	}

	if generated == "" {
		return fmt.Errorf("generated config.py content is empty")
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("create config.py directory: %w", err)
	}

	if err := os.WriteFile(configPath, []byte(generated), 0644); err != nil {
		return fmt.Errorf("write config.py: %w", err)
	}

	if err := updateKeysPyTerminal(keysPath, s.Qtile.Terminal); err != nil {
		return fmt.Errorf("update keys.py terminal: %w", err)
	}

	return nil
}

func updateKeysPyTerminal(keysPath, terminal string) error {
	if _, err := os.Stat(keysPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat keys.py: %w", err)
	}

	data, err := os.ReadFile(keysPath)
	if err != nil {
		return fmt.Errorf("read keys.py: %w", err)
	}

	content := string(data)
	oldLine := ""
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, `terminal = "`) || strings.HasPrefix(trimmed, "terminal = \"") {
			oldLine = line
			lines[i] = fmt.Sprintf("terminal = %q", terminal)
			break
		}
	}

	if oldLine == "" {
		return nil
	}

	newContent := strings.Join(lines, "\n")
	return os.WriteFile(keysPath, []byte(newContent), 0644)
}
