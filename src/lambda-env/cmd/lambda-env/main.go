package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"lambdaos.dev/lambda-env/internal/hub"
	"lambdaos.dev/lambda-env/internal/tui"
	"lambdaos.dev/lambda-env/pkg/version"
)

func main() {
	var (
		showHelp    = flag.Bool("help", false, "Show usage information")
		showVersion = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	if *showHelp {
		fmt.Println("lambda-env — LambdaOS Settings TUI")
		fmt.Println()
		fmt.Println("Usage: lambda-env [options]")
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("lambda-env version %s\n", version.Version)
		os.Exit(0)
	}

	settingsPath := defaultSettingsPath()

	h, err := hub.New(settingsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize hub: %v\n", err)
		os.Exit(1)
	}
	defer h.Logger.Close()

	m := tui.NewModel(h)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func defaultSettingsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/lambdaos-settings.json"
	}
	return filepath.Join(home, ".config", "lambdaos", "settings.json")
}
