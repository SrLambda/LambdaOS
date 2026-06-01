package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"lambdaos.dev/lambda-env/internal/settings"
)

func TestGenerateConfigPyDefaults(t *testing.T) {
	s := settings.QtileSettings{
		Terminal:           "kitty",
		Browser:            "firefox",
		DefaultFileManager: "thunar",
	}

	out, err := GenerateConfigPy(s)
	if err != nil {
		t.Fatalf("GenerateConfigPy() error = %v", err)
	}

	if !strings.Contains(out, `terminal = "kitty"`) {
		t.Error("expected output to contain terminal = \"kitty\"")
	}
	if !strings.Contains(out, `browser = "firefox"`) {
		t.Error("expected output to contain browser = \"firefox\"")
	}
	if !strings.Contains(out, `file_manager = "thunar"`) {
		t.Error("expected output to contain file_manager = \"thunar\"")
	}
}

func TestGenerateConfigPyCustomTerminal(t *testing.T) {
	s := settings.QtileSettings{
		Terminal:           "foot",
		Browser:            "firefox",
		DefaultFileManager: "thunar",
	}

	out, err := GenerateConfigPy(s)
	if err != nil {
		t.Fatalf("GenerateConfigPy() error = %v", err)
	}

	if !strings.Contains(out, `terminal = "foot"`) {
		t.Error("expected output to contain terminal = \"foot\"")
	}
}

func TestGenerateConfigPyCustomBrowser(t *testing.T) {
	s := settings.QtileSettings{
		Terminal:           "kitty",
		Browser:            "brave",
		DefaultFileManager: "thunar",
	}

	out, err := GenerateConfigPy(s)
	if err != nil {
		t.Fatalf("GenerateConfigPy() error = %v", err)
	}

	if !strings.Contains(out, `browser = "brave"`) {
		t.Error("expected output to contain browser = \"brave\"")
	}
}

func TestGenerateConfigPyNonEmpty(t *testing.T) {
	s := settings.QtileSettings{
		Terminal:           "kitty",
		Browser:            "firefox",
		DefaultFileManager: "thunar",
	}

	out, err := GenerateConfigPy(s)
	if err != nil {
		t.Fatalf("GenerateConfigPy() error = %v", err)
	}

	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestGenerateConfigPyValidPython(t *testing.T) {
	s := settings.QtileSettings{
		Terminal:           "kitty",
		Browser:            "firefox",
		DefaultFileManager: "thunar",
	}

	out, err := GenerateConfigPy(s)
	if err != nil {
		t.Fatalf("GenerateConfigPy() error = %v", err)
	}

	if !strings.Contains(out, "import os") {
		t.Error("expected output to contain 'import os' for valid Python")
	}
	if !strings.Contains(out, "from libqtile") {
		t.Error("expected output to contain 'from libqtile' for valid Python")
	}
}

func TestValidateConfigPyValid(t *testing.T) {
	tmpDir := t.TempDir()
	validPy := filepath.Join(tmpDir, "valid.py")
	content := `import os
print("hello")
`
	if err := os.WriteFile(validPy, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	err := ValidateConfigPy(validPy)
	if err != nil {
		t.Errorf("ValidateConfigPy() expected no error for valid Python, got: %v", err)
	}
}

func TestGenerateConfigPyWithGroups(t *testing.T) {
	s := settings.QtileSettings{
		Terminal:           "kitty",
		Browser:            "firefox",
		DefaultFileManager: "thunar",
		Groups: []settings.GroupConfig{
			{Name: "1"},
			{Name: "2"},
			{Name: "3"},
		},
	}

	out, err := GenerateConfigPy(s)
	if err != nil {
		t.Fatalf("GenerateConfigPy() error = %v", err)
	}

	if !strings.Contains(out, `Group("1")`) {
		t.Error("expected output to contain Group(\"1\") when Groups is populated")
	}
	if !strings.Contains(out, `Group("2")`) {
		t.Error("expected output to contain Group(\"2\") when Groups is populated")
	}
	if !strings.Contains(out, `Group("3")`) {
		t.Error("expected output to contain Group(\"3\") when Groups is populated")
	}
}

func TestUpdateKeysPyTerminal(t *testing.T) {
	tmpDir := t.TempDir()
	keysPath := filepath.Join(tmpDir, "keys.py")

	content := `import os
terminal = "kitty"
keys = []
`
	if err := os.WriteFile(keysPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write keys.py: %v", err)
	}

	if err := updateKeysPyTerminal(keysPath, "foot"); err != nil {
		t.Fatalf("updateKeysPyTerminal() error = %v", err)
	}

	data, err := os.ReadFile(keysPath)
	if err != nil {
		t.Fatalf("failed to read keys.py: %v", err)
	}

	if !strings.Contains(string(data), `terminal = "foot"`) {
		t.Errorf("expected terminal = \"foot\", got:\n%s", string(data))
	}
}

func TestUpdateKeysPyTerminalNoExistingTerminal(t *testing.T) {
	tmpDir := t.TempDir()
	keysPath := filepath.Join(tmpDir, "keys.py")

	content := `import os
keys = []
`
	if err := os.WriteFile(keysPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write keys.py: %v", err)
	}

	if err := updateKeysPyTerminal(keysPath, "foot"); err != nil {
		t.Fatalf("updateKeysPyTerminal() error = %v", err)
	}

	data, err := os.ReadFile(keysPath)
	if err != nil {
		t.Fatalf("failed to read keys.py: %v", err)
	}

	if string(data) != content {
		t.Errorf("expected file to be unchanged when no terminal line exists")
	}
}

func TestValidateConfigPyInvalid(t *testing.T) {
	tmpDir := t.TempDir()
	invalidPy := filepath.Join(tmpDir, "invalid.py")
	content := `def broken(
`
	if err := os.WriteFile(invalidPy, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	err := ValidateConfigPy(invalidPy)
	if err == nil {
		t.Error("ValidateConfigPy() expected error for invalid Python, got nil")
	}
}
