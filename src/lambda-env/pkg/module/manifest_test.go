package module

import (
	"strings"
	"testing"
)

func TestManifestValidate(t *testing.T) {
	tests := []struct {
		name     string
		manifest Manifest
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid manifest",
			manifest: Manifest{
				Name:          "screen",
				Version:       "0.1.0",
				Description:   "Manage display configuration",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "system",
				RequiresRoot:  false,
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			manifest: Manifest{
				Version:       "0.1.0",
				Description:   "Manage display configuration",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "system",
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
			},
			wantErr: true,
			errMsg:  "manifest name is required",
		},
		{
			name: "missing version",
			manifest: Manifest{
				Name:          "screen",
				Description:   "Manage display configuration",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "system",
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
			},
			wantErr: true,
			errMsg:  "manifest version is required",
		},
		{
			name: "invalid category",
			manifest: Manifest{
				Name:          "screen",
				Version:       "0.1.0",
				Description:   "Manage display configuration",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "invalid",
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
			},
			wantErr: true,
			errMsg:  "manifest category",
		},
		{
			name: "empty dependencies and zero timeout",
			manifest: Manifest{
				Name:          "screen",
				Version:       "0.1.0",
				Description:   "Manage display configuration",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "system",
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
				Timeout:       0,
			},
			wantErr: false,
		},
		{
			name: "name with spaces",
			manifest: Manifest{
				Name:          "My Module",
				Version:       "0.1.0",
				Description:   "Manage display configuration",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "system",
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
			},
			wantErr: true,
			errMsg:  "must be lowercase with hyphens only",
		},
		{
			name: "missing description",
			manifest: Manifest{
				Name:          "screen",
				Version:       "0.1.0",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "system",
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
			},
			wantErr: true,
			errMsg:  "manifest description is required",
		},
		{
			name: "missing description_es",
			manifest: Manifest{
				Name:          "screen",
				Version:       "0.1.0",
				Description:   "Manage display configuration",
				Category:      "system",
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
			},
			wantErr: true,
			errMsg:  "manifest description_es is required",
		},
		{
			name: "missing min_hub_version",
			manifest: Manifest{
				Name:          "screen",
				Version:       "0.1.0",
				Description:   "Manage display configuration",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "system",
				Dependencies:  []string{},
			},
			wantErr: true,
			errMsg:  "manifest min_hub_version is required",
		},
		{
			name: "valid manifest with actions",
			manifest: Manifest{
				Name:          "neovim",
				Version:       "0.1.0",
				Description:   "Configure Neovim",
				DescriptionES: "Configurar Neovim",
				Category:      "apps",
				Dependencies:  []string{},
				MinHubVersion: "1.0.0",
				Actions: []ActionConfig{
					{Name: "toggle-lsp", Label: "Enable LSP", Type: "toggle", Field: "neovim.enable_lsp"},
					{Name: "set-theme", Label: "Theme", Type: "select", Field: "neovim.theme", Options: []string{"tokyonight", "gruvbox"}},
				},
			},
			wantErr: false,
		},
		{
			name: "manifest with invalid action type",
			manifest: Manifest{
				Name:          "neovim",
				Version:       "0.1.0",
				Description:   "Configure Neovim",
				DescriptionES: "Configurar Neovim",
				Category:      "apps",
				Dependencies:  []string{},
				MinHubVersion: "1.0.0",
				Actions: []ActionConfig{
					{Name: "bad-action", Label: "Bad", Type: "invalid", Field: "neovim.bad"},
				},
			},
			wantErr: true,
			errMsg:  "invalid action type",
		},
		{
			name: "manifest with select action missing options",
			manifest: Manifest{
				Name:          "neovim",
				Version:       "0.1.0",
				Description:   "Configure Neovim",
				DescriptionES: "Configurar Neovim",
				Category:      "apps",
				Dependencies:  []string{},
				MinHubVersion: "1.0.0",
				Actions: []ActionConfig{
					{Name: "set-theme", Label: "Theme", Type: "select", Field: "neovim.theme"},
				},
			},
			wantErr: true,
			errMsg:  "select action must have options",
		},
		{
			name: "manifest with list action missing options",
			manifest: Manifest{
				Name:          "neovim",
				Version:       "0.1.0",
				Description:   "Configure Neovim",
				DescriptionES: "Configurar Neovim",
				Category:      "apps",
				Dependencies:  []string{},
				MinHubVersion: "1.0.0",
				Actions: []ActionConfig{
					{Name: "set-servers", Label: "Servers", Type: "list", Field: "neovim.lsp_servers"},
				},
			},
			wantErr: true,
			errMsg:  "list action must have options",
		},
		{
			name: "manifest with empty actions is valid",
			manifest: Manifest{
				Name:          "dotfiles",
				Version:       "0.1.0",
				Description:   "Manage dotfiles",
				DescriptionES: "Gestionar dotfiles",
				Category:      "ops",
				Dependencies:  []string{},
				MinHubVersion: "1.0.0",
				Actions:       []ActionConfig{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.manifest.Validate()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errMsg)
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Fatalf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestActionConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		action  ActionConfig
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid toggle action",
			action:  ActionConfig{Name: "toggle-mute", Label: "Mute", Type: "toggle", Field: "audio.muted"},
			wantErr: false,
		},
		{
			name:    "valid text action",
			action:  ActionConfig{Name: "set-layout", Label: "Layout", Type: "text", Field: "keyboard.layout"},
			wantErr: false,
		},
		{
			name:    "valid confirm action",
			action:  ActionConfig{Name: "reset-config", Label: "Reset", Type: "confirm", Field: ""},
			wantErr: false,
		},
		{
			name:    "valid execute action",
			action:  ActionConfig{Name: "apply-changes", Label: "Apply", Type: "execute", Field: ""},
			wantErr: false,
		},
		{
			name:    "valid select with options",
			action:  ActionConfig{Name: "set-theme", Label: "Theme", Type: "select", Field: "appearance.theme", Options: []string{"dark", "light"}},
			wantErr: false,
		},
		{
			name:    "valid list with options",
			action:  ActionConfig{Name: "set-servers", Label: "Servers", Type: "list", Field: "neovim.lsp_servers", Options: []string{"gopls", "pyright"}},
			wantErr: false,
		},
		{
			name:    "invalid type",
			action:  ActionConfig{Name: "bad", Label: "Bad", Type: "unknown", Field: ""},
			wantErr: true,
			errMsg:  "invalid action type",
		},
		{
			name:    "empty name",
			action:  ActionConfig{Name: "", Label: "Bad", Type: "toggle", Field: ""},
			wantErr: true,
			errMsg:  "action name is required",
		},
		{
			name:    "empty label",
			action:  ActionConfig{Name: "bad", Label: "", Type: "toggle", Field: ""},
			wantErr: true,
			errMsg:  "action label is required",
		},
		{
			name:    "select without options",
			action:  ActionConfig{Name: "set-theme", Label: "Theme", Type: "select", Field: "appearance.theme"},
			wantErr: true,
			errMsg:  "select action must have options",
		},
		{
			name:    "list without options",
			action:  ActionConfig{Name: "set-servers", Label: "Servers", Type: "list", Field: "neovim.lsp_servers"},
			wantErr: true,
			errMsg:  "list action must have options",
		},
		{
			name:    "toggle with options is ok",
			action:  ActionConfig{Name: "toggle", Label: "Toggle", Type: "toggle", Field: "x.y", Options: []string{"a", "b"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.action.Validate()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errMsg)
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Fatalf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}
