package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"lambdaos.dev/lambda-env/internal/modules/neovim/templates"
	"lambdaos.dev/lambda-env/internal/settings"
)

func GenerateLazyLua(s settings.NeovimSettings) (string, error) {
	tmpl, err := template.New("lazy.lua").Parse(templates.LazyLuaTemplate)
	if err != nil {
		return "", fmt.Errorf("parse lazy.lua template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, s); err != nil {
		return "", fmt.Errorf("execute lazy.lua template: %w", err)
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

	lazyPath := filepath.Join(home, ".config", "nvim", "lua", "core", "lazy.lua")
	backupPath := lazyPath + ".bak"

	if _, err := os.Stat(lazyPath); err == nil {
		data, readErr := os.ReadFile(lazyPath)
		if readErr != nil {
			return fmt.Errorf("read existing lazy.lua: %w", readErr)
		}
		if writeErr := os.WriteFile(backupPath, data, 0644); writeErr != nil {
			return fmt.Errorf("backup lazy.lua: %w", writeErr)
		}
	}

	generated, err := GenerateLazyLua(s.Neovim)
	if err != nil {
		return err
	}

	if generated == "" {
		return fmt.Errorf("generated lazy.lua content is empty")
	}

	lazyDir := filepath.Dir(lazyPath)
	if err := os.MkdirAll(lazyDir, 0755); err != nil {
		return fmt.Errorf("create lazy.lua directory: %w", err)
	}

	if err := os.WriteFile(lazyPath, []byte(generated), 0644); err != nil {
		return fmt.Errorf("write lazy.lua: %w", err)
	}

	return nil
}
