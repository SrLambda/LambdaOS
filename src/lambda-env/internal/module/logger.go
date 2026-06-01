package module

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SystemLogDir is the default log directory when running as root.
const SystemLogDir = "/var/log/lambda-env"

// Logger writes structured module execution logs.
type Logger struct {
	file *os.File
}

// logDir returns the effective log directory.
// - Root: /var/log/lambda-env (system-wide logs)
// - Non-root: ~/.local/share/lambda-env/logs (user-scoped logs)
// - Override: LAMBDA_ENV_LOG_DIR env var takes precedence.
func logDir() string {
	if d := os.Getenv("LAMBDA_ENV_LOG_DIR"); d != "" {
		return d
	}
	if os.Geteuid() == 0 {
		return "/var/log/lambda-env"
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "/var/log/lambda-env" // fallback to system path
	}
	return filepath.Join(home, ".local", "share", "lambda-env", "logs")
}

// logFile returns the effective log file path.
func logFile() string {
	return logDir() + "/modules.log"
}

// NewLogger creates a logger, ensuring the log directory exists.
func NewLogger() (*Logger, error) {
	dir := logDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}

	f, err := os.OpenFile(logFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	return &Logger{file: f}, nil
}

// Log writes a structured log entry for a module execution.
func (l *Logger) Log(module, action string, exitCode int, stdout, stderr string, env map[string]string) error {
	level := "INFO"
	if exitCode != 0 {
		level = "ERROR"
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)

	var envParts []string
	for k, v := range env {
		envParts = append(envParts, fmt.Sprintf("%s=%s", k, v))
	}
	envStr := strings.Join(envParts, ", ")

	entry := fmt.Sprintf(
		"%s [%s] module=%s action=%s exit_code=%d\n%s\n%s\n  env: %s\n",
		timestamp,
		level,
		module,
		action,
		exitCode,
		indentLines("  stdout: ", stdout),
		indentLines("  stderr: ", stderr),
		envStr,
	)

	_, err := l.file.WriteString(entry)
	return err
}

// Close closes the log file handle.
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func indentLines(prefix, text string) string {
	if text == "" {
		return prefix
	}
	lines := strings.Split(text, "\n")
	result := prefix + lines[0]
	indent := strings.Repeat(" ", len(prefix))
	for i := 1; i < len(lines); i++ {
		result += "\n" + indent + lines[i]
	}
	return result
}
