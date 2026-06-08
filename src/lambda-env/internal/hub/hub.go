package hub

import (
	"fmt"
	"os"
	"os/exec"
	"sort"

	modlogger "lambdaos.dev/lambda-env/internal/module"
	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/internal/tui/icons"
	"lambdaos.dev/lambda-env/pkg/module"
)

// MenuCategory groups modules by their category for TUI rendering.
type MenuCategory struct {
	Name    string
	Modules []module.Manifest
	Count   int
}

// Hub is the central controller for module discovery and execution.
type Hub struct {
	Store     *settings.Settings
	StorePath string
	Modules   []module.Manifest
	Logger    *modlogger.Logger
}

// New creates a Hub, initializing the settings store and running module discovery.
func New(settingsPath string, nerdFonts bool) (*Hub, error) {
	store, err := settings.Load(settingsPath)
	if err != nil {
		return nil, fmt.Errorf("load settings: %w", err)
	}

	logger, err := modlogger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	iconProvider := icons.NewProvider(nerdFonts)
	modules, err := Scan(iconProvider)
	if err != nil {
		logger.Close()
		return nil, fmt.Errorf("module discovery: %w", err)
	}

	return &Hub{
		Store:     store,
		StorePath: settingsPath,
		Modules:   modules,
		Logger:    logger,
	}, nil
}

// BuildMenu groups discovered modules by category.
// Categories with no modules are omitted.
func (h *Hub) BuildMenu() []MenuCategory {
	groups := make(map[string][]module.Manifest)
	for _, m := range h.Modules {
		groups[m.Category] = append(groups[m.Category], m)
	}

	// Sort modules within each category.
	for _, mods := range groups {
		sort.Slice(mods, func(i, j int) bool {
			return mods[i].Name < mods[j].Name
		})
	}

	var result []MenuCategory
	for _, cat := range module.ValidCategories {
		if mods, ok := groups[cat]; ok && len(mods) > 0 {
			result = append(result, MenuCategory{
				Name:    cat,
				Modules: mods,
				Count:   len(mods),
			})
		}
	}

	return result
}

// CheckDeps verifies that all packages in deps are installed using pacman -Q.
// Returns true if all are installed, false otherwise.
func (h *Hub) CheckDeps(deps []string) bool {
	if len(deps) == 0 {
		return true
	}
	for _, dep := range deps {
		cmd := exec.Command("pacman", "-Q", dep)
		if err := cmd.Run(); err != nil {
			return false
		}
	}
	return true
}

// CheckRoot returns true if the current process is running as root.
func (h *Hub) CheckRoot() bool {
	return os.Geteuid() == 0
}
