package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"lambdaos.dev/lambda-env/internal/hub"
	"lambdaos.dev/lambda-env/internal/tui"
	"lambdaos.dev/lambda-env/internal/tui/icons"
	"lambdaos.dev/lambda-env/pkg/version"
)

// nerdFontsFlag is a custom flag.Value that tracks whether the flag was
// explicitly provided on the command line.
type nerdFontsFlag struct {
	value bool
	set   bool
}

func (n *nerdFontsFlag) Set(val string) error {
	n.set = true
	v, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	n.value = v
	return nil
}

func (n *nerdFontsFlag) String() string {
	return strconv.FormatBool(n.value)
}

func (n *nerdFontsFlag) IsBoolFlag() bool { return true }

func main() {
	var (
		showHelp    = flag.Bool("help", false, "Show usage information")
		showVersion = flag.Bool("version", false, "Show version information")
		nerdFonts   = &nerdFontsFlag{}
	)
	flag.Var(nerdFonts, "nerd-fonts", "Enable Nerd Font icons")
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
	enabled := icons.ResolveNerdFonts(nerdFonts.set, nerdFonts.value)

	h, err := hub.New(settingsPath, enabled)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize hub: %v\n", err)
		os.Exit(1)
	}
	defer h.Logger.Close()

	m := tui.NewModel(h, enabled)
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
