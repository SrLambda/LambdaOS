package settings

import (
	"fmt"
)

// CurrentVersion is the schema version supported by this build.
const CurrentVersion = "1.0.0"

// Settings is the root settings struct.
type Settings struct {
	Version    string             `json:"version"`
	Appearance AppearanceSettings `json:"appearance"`
	Display    DisplaySettings    `json:"display"`
	Audio      AudioSettings      `json:"audio"`
	Network    NetworkSettings    `json:"network"`
	Bluetooth  BluetoothSettings  `json:"bluetooth"`
	Keyboard   KeyboardSettings   `json:"keyboard"`
	Neovim     NeovimSettings     `json:"neovim"`
	Qtile      QtileSettings      `json:"qtile"`
	Services   ServicesSettings   `json:"services"`
}

// AppearanceSettings defines look-and-feel options.
type AppearanceSettings struct {
	Theme    string `json:"theme"`
	FontSize int    `json:"font_size"`
	Opacity  int    `json:"opacity"`
	Wallpaper string `json:"wallpaper"`
}

// OutputConfig defines a single display output.
type OutputConfig struct {
	Name     string `json:"name"`
	Mode     string `json:"mode"`
	Position string `json:"position"`
	Primary  bool   `json:"primary,omitempty"`
}

// OutputProfile groups outputs under a named profile.
type OutputProfile struct {
	Name    string         `json:"name"`
	Outputs []OutputConfig `json:"outputs"`
}

// DisplaySettings defines display configuration.
type DisplaySettings struct {
	ActiveProfile string          `json:"active_profile"`
	Profiles      []OutputProfile `json:"profiles"`
}

// AudioSettings defines audio configuration.
type AudioSettings struct {
	DefaultSink string `json:"default_sink"`
	Volume      int    `json:"volume"`
	Muted       bool   `json:"muted"`
}

// NetworkSettings defines network configuration.
type NetworkSettings struct {
	WifiEnabled    bool     `json:"wifi_enabled"`
	KnownNetworks  []string `json:"known_networks"`
}

// BluetoothSettings defines bluetooth configuration.
type BluetoothSettings struct {
	Enabled       bool     `json:"enabled"`
	PairedDevices []string `json:"paired_devices"`
}

// KeyboardSettings defines keyboard configuration.
type KeyboardSettings struct {
	Layout string `json:"layout"`
	Variant string `json:"variant"`
	Options string `json:"options"`
}

// NeovimSettings defines Neovim configuration.
type NeovimSettings struct {
	Theme   string `json:"theme"`
	Font    string `json:"font"`
	Lines   int    `json:"lines"`
	Columns int    `json:"columns"`
}

// QtileSettings defines Qtile window manager configuration.
type QtileSettings struct {
	BarPosition string   `json:"bar_position"`
	BarSize     int      `json:"bar_size"`
	Layouts     []string `json:"layouts"`
}

// ServicesSettings defines enabled/disabled services.
type ServicesSettings struct {
	Enabled []string `json:"enabled"`
}

// Defaults returns a fully populated Settings with default values.
func Defaults() Settings {
	return Settings{
		Version: CurrentVersion,
		Appearance: AppearanceSettings{
			Theme:     "dark",
			FontSize:  14,
			Opacity:   100,
			Wallpaper: "",
		},
		Display: DisplaySettings{
			ActiveProfile: "default",
			Profiles:      []OutputProfile{},
		},
		Audio: AudioSettings{
			DefaultSink: "",
			Volume:      75,
			Muted:       false,
		},
		Network: NetworkSettings{
			WifiEnabled:   true,
			KnownNetworks: []string{},
		},
		Bluetooth: BluetoothSettings{
			Enabled:       true,
			PairedDevices: []string{},
		},
		Keyboard: KeyboardSettings{
			Layout:  "us",
			Variant: "",
			Options: "",
		},
		Neovim: NeovimSettings{
			Theme:   "tokyonight",
			Font:    "JetBrainsMono",
			Lines:   40,
			Columns: 120,
		},
		Qtile: QtileSettings{
			BarPosition: "top",
			BarSize:     24,
			Layouts:     []string{},
		},
		Services: ServicesSettings{
			Enabled: []string{},
		},
	}
}

// Validate checks that all fields have valid values.
func (s *Settings) Validate() error {
	if s.Version == "" {
		return fmt.Errorf("version is required")
	}

	if s.Audio.Volume < 0 || s.Audio.Volume > 100 {
		return fmt.Errorf("audio.volume must be between 0 and 100, got %d", s.Audio.Volume)
	}

	if s.Appearance.FontSize < 1 {
		return fmt.Errorf("appearance.font_size must be > 0, got %d", s.Appearance.FontSize)
	}

	if s.Appearance.Opacity < 0 || s.Appearance.Opacity > 100 {
		return fmt.Errorf("appearance.opacity must be between 0 and 100, got %d", s.Appearance.Opacity)
	}

	if s.Neovim.Lines < 1 {
		return fmt.Errorf("neovim.lines must be > 0, got %d", s.Neovim.Lines)
	}

	if s.Neovim.Columns < 1 {
		return fmt.Errorf("neovim.columns must be > 0, got %d", s.Neovim.Columns)
	}

	if s.Qtile.BarSize < 1 {
		return fmt.Errorf("qtile.bar_size must be > 0, got %d", s.Qtile.BarSize)
	}

	if len(s.Display.Profiles) > 0 {
		found := false
		for _, p := range s.Display.Profiles {
			if p.Name == s.Display.ActiveProfile {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("display.active_profile %q does not match any profile name", s.Display.ActiveProfile)
		}
	}

	return nil
}
