package settings

import (
	"fmt"
)

// CurrentVersion is the schema version supported by this build.
const CurrentVersion = "1.1.0"

// Settings is the root settings struct.
type Settings struct {
	Version       string               `json:"version"`
	Appearance    AppearanceSettings   `json:"appearance"`
	Display       DisplaySettings      `json:"display"`
	Audio         AudioSettings        `json:"audio"`
	Network       NetworkSettings      `json:"network"`
	Bluetooth     BluetoothSettings    `json:"bluetooth"`
	Keyboard      KeyboardSettings     `json:"keyboard"`
	Neovim        NeovimSettings       `json:"neovim"`
	Qtile         QtileSettings        `json:"qtile"`
	Services      ServicesSettings     `json:"services"`
	Power         PowerSettings        `json:"power"`
	Defaults      DefaultsSettings     `json:"defaults"`
	Autostart     AutostartSettings    `json:"autostart"`
	Updates       UpdatesSettings      `json:"updates"`
	Security      SecuritySettings     `json:"security"`
	Fonts         FontsSettings        `json:"fonts"`
	Notifications NotificationsSettings `json:"notifications"`
}

// AppearanceSettings defines look-and-feel options.
type AppearanceSettings struct {
	Theme     string `json:"theme"`
	FontSize  int    `json:"font_size"`
	Opacity   int    `json:"opacity"`
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
	WifiEnabled   bool     `json:"wifi_enabled"`
	KnownNetworks []string `json:"known_networks"`
}

// BluetoothSettings defines bluetooth configuration.
type BluetoothSettings struct {
	Enabled       bool     `json:"enabled"`
	PairedDevices []string `json:"paired_devices"`
}

// KeyboardSettings defines keyboard configuration.
type KeyboardSettings struct {
	Layout  string `json:"layout"`
	Variant string `json:"variant"`
	Options string `json:"options"`
}

// GroupConfig defines a Qtile workspace group.
type GroupConfig struct {
	Name string `json:"name"`
	Icon string `json:"icon,omitempty"`
}

// NeovimSettings defines Neovim configuration.
type NeovimSettings struct {
	Theme         string   `json:"theme"`
	Font          string   `json:"font"`
	Lines         int      `json:"lines"`
	Columns       int      `json:"columns"`
	EnableLSP     bool     `json:"enable_lsp"`
	EnableCopilot bool     `json:"enable_copilot"`
	EnableNeotree bool     `json:"enable_neotree"`
	LspServers    []string `json:"lsp_servers"`
	UseGlobalTheme bool    `json:"use_global_theme"`
}

// QtileSettings defines Qtile window manager configuration.
type QtileSettings struct {
	BarPosition        string        `json:"bar_position"`
	BarSize            int           `json:"bar_size"`
	Layouts            []string      `json:"layouts"`
	Terminal           string        `json:"terminal"`
	Browser            string        `json:"browser"`
	DefaultFileManager string        `json:"default_file_manager"`
	Groups             []GroupConfig `json:"groups"`
	UseGlobalTheme     bool          `json:"use_global_theme"`
}

// ServicesSettings defines enabled/disabled services.
type ServicesSettings struct {
	Enabled []string `json:"enabled"`
}

// PowerSettings defines power management configuration.
type PowerSettings struct {
	ScreenTimeout  int    `json:"screen_timeout"`
	SleepTimeout   int    `json:"sleep_timeout"`
	LidCloseAction string `json:"lid_close_action"`
}

// DefaultsSettings defines default application assignments.
type DefaultsSettings struct {
	Browser     string `json:"browser"`
	Terminal    string `json:"terminal"`
	Editor      string `json:"editor"`
	FileManager string `json:"file_manager"`
}

// AutostartSettings defines applications that start on login.
type AutostartSettings struct {
	Enabled []string `json:"enabled"`
}

// UpdatesSettings defines system update configuration.
type UpdatesSettings struct {
	AutoUpdate      bool     `json:"auto_update"`
	CheckInterval   int      `json:"check_interval"`
	ExcludePackages []string `json:"exclude_packages"`
}

// SecuritySettings defines security-related configuration.
type SecuritySettings struct {
	FirewallEnabled   bool `json:"firewall_enabled"`
	SudoTimeout       int  `json:"sudo_timeout"`
	ScreenLockTimeout int  `json:"screen_lock_timeout"`
}

// FontsSettings defines font configuration.
type FontsSettings struct {
	Monospace  string `json:"monospace"`
	SansSerif  string `json:"sans_serif"`
	Serif      string `json:"serif"`
	FontSize   int    `json:"font_size"`
}

// NotificationsSettings defines notification daemon configuration.
type NotificationsSettings struct {
	Enabled         bool `json:"enabled"`
	DoNotDisturb    bool `json:"do_not_disturb"`
	TimeoutSeconds  int  `json:"timeout_seconds"`
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
			Theme:          "tokyonight",
			Font:           "JetBrainsMono",
			Lines:          40,
			Columns:        120,
			EnableLSP:      true,
			EnableCopilot:  true,
			EnableNeotree:  true,
			LspServers:     []string{"gopls", "pyright"},
			UseGlobalTheme: false,
		},
		Qtile: QtileSettings{
			BarPosition:        "top",
			BarSize:            24,
			Layouts:            []string{},
			Terminal:           "kitty",
			Browser:            "firefox",
			DefaultFileManager: "thunar",
			Groups: []GroupConfig{
				{Name: "1"}, {Name: "2"}, {Name: "3"},
				{Name: "4"}, {Name: "5"}, {Name: "6"},
				{Name: "7"}, {Name: "8"}, {Name: "9"},
			},
			UseGlobalTheme: false,
		},
		Services: ServicesSettings{
			Enabled: []string{},
		},
		Power: PowerSettings{
			ScreenTimeout:  300,
			SleepTimeout:   600,
			LidCloseAction: "suspend",
		},
		Defaults: DefaultsSettings{
			Browser:     "",
			Terminal:    "",
			Editor:      "",
			FileManager: "",
		},
		Autostart: AutostartSettings{
			Enabled: []string{},
		},
		Updates: UpdatesSettings{
			AutoUpdate:      false,
			CheckInterval:   86400,
			ExcludePackages: []string{},
		},
		Security: SecuritySettings{
			FirewallEnabled:   true,
			SudoTimeout:       5,
			ScreenLockTimeout: 300,
		},
		Fonts: FontsSettings{
			Monospace: "JetBrainsMono",
			SansSerif: "Noto Sans",
			Serif:     "Noto Serif",
			FontSize:  14,
		},
		Notifications: NotificationsSettings{
			Enabled:        true,
			DoNotDisturb:   false,
			TimeoutSeconds: 5,
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

	if s.Qtile.Terminal != "" {
		validTerminals := map[string]bool{
			"kitty":     true,
			"foot":      true,
			"alacritty": true,
			"st":        true,
			"xterm":     true,
		}
		if !validTerminals[s.Qtile.Terminal] {
			return fmt.Errorf("qtile.terminal %q is not allowed, must be one of: kitty, foot, alacritty, st, xterm", s.Qtile.Terminal)
		}
	}

	if s.Qtile.Browser != "" {
		validBrowsers := map[string]bool{
			"firefox":  true,
			"chromium": true,
			"brave":    true,
			"chrome":   true,
		}
		if !validBrowsers[s.Qtile.Browser] {
			return fmt.Errorf("qtile.browser %q is not allowed, must be one of: firefox, chromium, brave, chrome", s.Qtile.Browser)
		}
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

	if s.Power.ScreenTimeout < 0 {
		return fmt.Errorf("power.screen_timeout must be >= 0, got %d", s.Power.ScreenTimeout)
	}

	if s.Power.SleepTimeout < 0 {
		return fmt.Errorf("power.sleep_timeout must be >= 0, got %d", s.Power.SleepTimeout)
	}

	if s.Security.SudoTimeout < 0 {
		return fmt.Errorf("security.sudo_timeout must be >= 0, got %d", s.Security.SudoTimeout)
	}

	if s.Security.ScreenLockTimeout < 0 {
		return fmt.Errorf("security.screen_lock_timeout must be >= 0, got %d", s.Security.ScreenLockTimeout)
	}

	if s.Fonts.FontSize < 1 {
		return fmt.Errorf("fonts.font_size must be > 0, got %d", s.Fonts.FontSize)
	}

	if s.Updates.CheckInterval < 0 {
		return fmt.Errorf("updates.check_interval must be >= 0, got %d", s.Updates.CheckInterval)
	}

	if s.Notifications.TimeoutSeconds < 0 {
		return fmt.Errorf("notifications.timeout_seconds must be >= 0, got %d", s.Notifications.TimeoutSeconds)
	}

	return nil
}
