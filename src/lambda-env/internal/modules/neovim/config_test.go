package main

import (
	"strings"
	"testing"

	"lambdaos.dev/lambda-env/internal/settings"
)

func TestGenerateLazyLuaAllTogglesOn(t *testing.T) {
	s := settings.NeovimSettings{
		EnableLSP:     true,
		EnableCopilot: true,
		EnableNeotree: true,
	}

	out, err := GenerateLazyLua(s)
	if err != nil {
		t.Fatalf("GenerateLazyLua() error = %v", err)
	}

	if !strings.Contains(out, `{ import = "plugins.lsp" }`) {
		t.Error("expected output to contain plugins.lsp import when EnableLSP=true")
	}
	if !strings.Contains(out, `{ import = "plugins.ai" }`) {
		t.Error("expected output to contain plugins.ai import when EnableCopilot=true")
	}
}

func TestGenerateLazyLuaLspOff(t *testing.T) {
	s := settings.NeovimSettings{
		EnableLSP:     false,
		EnableCopilot: true,
		EnableNeotree: true,
	}

	out, err := GenerateLazyLua(s)
	if err != nil {
		t.Fatalf("GenerateLazyLua() error = %v", err)
	}

	if strings.Contains(out, `{ import = "plugins.lsp" }`) {
		t.Error("expected output NOT to contain plugins.lsp import when EnableLSP=false")
	}
	if !strings.Contains(out, `{ import = "plugins.ai" }`) {
		t.Error("expected output to contain plugins.ai import when EnableCopilot=true")
	}
}

func TestGenerateLazyLuaAllOff(t *testing.T) {
	s := settings.NeovimSettings{
		EnableLSP:     false,
		EnableCopilot: false,
		EnableNeotree: false,
	}

	out, err := GenerateLazyLua(s)
	if err != nil {
		t.Fatalf("GenerateLazyLua() error = %v", err)
	}

	if strings.Contains(out, `{ import = "plugins.lsp" }`) {
		t.Error("expected output NOT to contain plugins.lsp import when all toggles off")
	}
	if strings.Contains(out, `{ import = "plugins.ai" }`) {
		t.Error("expected output NOT to contain plugins.ai import when all toggles off")
	}

	if !strings.Contains(out, `{ import = "plugins" }`) {
		t.Error("expected output to contain base plugins import")
	}
}

func TestGenerateLazyLuaNonEmpty(t *testing.T) {
	s := settings.NeovimSettings{
		EnableLSP:     true,
		EnableCopilot: true,
		EnableNeotree: true,
	}

	out, err := GenerateLazyLua(s)
	if err != nil {
		t.Fatalf("GenerateLazyLua() error = %v", err)
	}

	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestGenerateLazyLuaValidLua(t *testing.T) {
	s := settings.NeovimSettings{
		EnableLSP:     true,
		EnableCopilot: true,
		EnableNeotree: true,
	}

	out, err := GenerateLazyLua(s)
	if err != nil {
		t.Fatalf("GenerateLazyLua() error = %v", err)
	}

	if !strings.Contains(out, `require("lazy").setup`) {
		t.Error("expected output to contain require(\"lazy\").setup for valid Lua")
	}
}

func TestGenerateLazyLuaNeotreeOff(t *testing.T) {
	s := settings.NeovimSettings{
		EnableLSP:     true,
		EnableCopilot: true,
		EnableNeotree: false,
	}

	out, err := GenerateLazyLua(s)
	if err != nil {
		t.Fatalf("GenerateLazyLua() error = %v", err)
	}

	if strings.Contains(out, "neo-tree") {
		t.Error("expected output NOT to contain neo-tree when EnableNeotree=false")
	}
}

func TestGenerateLazyLuaNeotreeOn(t *testing.T) {
	s := settings.NeovimSettings{
		EnableLSP:     true,
		EnableCopilot: true,
		EnableNeotree: true,
	}

	out, err := GenerateLazyLua(s)
	if err != nil {
		t.Fatalf("GenerateLazyLua() error = %v", err)
	}

	if !strings.Contains(out, "neo-tree") {
		t.Error("expected output to contain neo-tree when EnableNeotree=true")
	}
}

func TestResolveNeovimThemeGlobalEnabled(t *testing.T) {
	s := settings.Settings{
		Appearance: settings.AppearanceSettings{Theme: "dark"},
		Neovim:     settings.NeovimSettings{UseGlobalTheme: true, Theme: "gruvbox"},
	}
	got := resolveNeovimTheme(s)
	if got != "tokyonight" {
		t.Errorf("resolveNeovimTheme = %q, want tokyonight", got)
	}
}

func TestResolveNeovimThemeGlobalDisabled(t *testing.T) {
	s := settings.Settings{
		Appearance: settings.AppearanceSettings{Theme: "dark"},
		Neovim:     settings.NeovimSettings{UseGlobalTheme: false, Theme: "gruvbox"},
	}
	got := resolveNeovimTheme(s)
	if got != "gruvbox" {
		t.Errorf("resolveNeovimTheme = %q, want gruvbox", got)
	}
}
