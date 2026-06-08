package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

var executor module.CLIExecutor = module.NewRealExecutor()

func main() {
	action := os.Getenv("LAMBDA_ENV_ACTION")
	settingsPath := os.Getenv("LAMBDA_ENV_SETTINGS")

	if settingsPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			emitError("run", "cannot determine home directory", "")
			return
		}
		settingsPath = filepath.Join(home, ".config", "lambdaos", "settings.json")
	}

	params := readParams()

	switch action {
	case "run":
		handleRun(settingsPath)
	case "set-mode":
		value := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				value = v
			}
		}
		handleSetMode(settingsPath, value)
	case "set-position":
		value := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				value = v
			}
		}
		handleSetPosition(settingsPath, value)
	case "set-primary":
		var value interface{}
		if params != nil {
			if v, ok := params["value"]; ok {
				value = v
			}
		}
		handleSetPrimary(settingsPath, value)
	case "save-profile":
		name := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				name = v
			}
		}
		handleSaveProfile(settingsPath, name)
	case "load-profile":
		name := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				name = v
			}
		}
		handleLoadProfile(settingsPath, name)
	default:
		emitError(action, "unknown action", "use run, set-mode, set-position, set-primary, save-profile, or load-profile")
	}
}

func readParams() map[string]interface{} {
	p := os.Getenv("LAMBDA_ENV_PARAMS")
	if p == "" {
		return nil
	}
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(p), &params); err != nil {
		return nil
	}
	return params
}

// detectSessionType returns "x11" or "wayland" based on XDG_SESSION_TYPE and binary availability.
func detectSessionType(exe module.CLIExecutor) (string, error) {
	session := os.Getenv("XDG_SESSION_TYPE")
	session = strings.ToLower(session)

	if session == "wayland" {
		_, _, exitCode, _ := exe.Run("which", "wlr-randr")
		if exitCode == 0 {
			return "wayland", nil
		}
		return "", fmt.Errorf("wayland session detected but wlr-randr not found")
	}

	if session == "x11" {
		_, _, exitCode, _ := exe.Run("which", "xrandr")
		if exitCode == 0 {
			return "x11", nil
		}
		return "", fmt.Errorf("x11 session detected but xrandr not found")
	}

	// Fallback: probe binaries.
	_, _, exitCode, _ := exe.Run("which", "wlr-randr")
	if exitCode == 0 {
		return "wayland", nil
	}
	_, _, exitCode, _ = exe.Run("which", "xrandr")
	if exitCode == 0 {
		return "x11", nil
	}

	return "", fmt.Errorf("cannot determine display session type (XDG_SESSION_TYPE unset and neither xrandr nor wlr-randr found)")
}

// DisplayMode represents a single resolution+refresh option.
type DisplayMode struct {
	Resolution  string `json:"resolution"`
	RefreshRate string `json:"refresh_rate"`
	Preferred   bool   `json:"preferred"`
}

// DisplayOutput represents a detected monitor.
type DisplayOutput struct {
	Name        string        `json:"name"`
	Connected   bool          `json:"connected"`
	Primary     bool          `json:"primary"`
	CurrentMode string        `json:"current_mode"`
	Position    string        `json:"position"`
	Modes       []DisplayMode `json:"modes,omitempty"`
}

func parseXrandrOutput(stdout string) []DisplayOutput {
	var outputs []DisplayOutput

	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimRight(line, "\n")
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, " ") {
			// Mode line for the last output.
			if len(outputs) == 0 {
				continue
			}
			idx := len(outputs) - 1
			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}
			resolution := parts[0]
			lineHasPlus := strings.Contains(line, "+")
			for i := 1; i < len(parts); i++ {
				rateStr := strings.TrimSuffix(parts[i], "*")
				rateStr = strings.TrimSuffix(rateStr, "+")
				if rateStr == "" {
					continue
				}
				if _, err := strconv.ParseFloat(rateStr, 64); err != nil {
					continue
				}
				mode := DisplayMode{
					Resolution:  resolution,
					RefreshRate: rateStr,
					Preferred:   lineHasPlus,
				}
				outputs[idx].Modes = append(outputs[idx].Modes, mode)
				break // Only append once per line with the first valid rate.
			}
			continue
		}

		// Output line.
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		name := parts[0]
		if name == "Screen" {
			continue
		}
		status := parts[1]
		connected := status == "connected"

		primary := false
		currentMode := ""
		position := ""

		if connected {
			for i := 2; i < len(parts); i++ {
				if parts[i] == "primary" {
					primary = true
					continue
				}
				// Look for 1920x1080+0+0 pattern.
				if m := regexp.MustCompile(`^(\d+x\d+)\+(-?\d+)\+(-?\d+)$`).FindStringSubmatch(parts[i]); len(m) > 1 {
					currentMode = m[1]
					position = m[2] + "," + m[3]
				}
			}
		}

		outputs = append(outputs, DisplayOutput{
			Name:        name,
			Connected:   connected,
			Primary:     primary,
			CurrentMode: currentMode,
			Position:    position,
		})
	}

	return outputs
}

func parseWlrRandrOutput(stdout string) []DisplayOutput {
	var outputs []DisplayOutput

	for _, line := range strings.Split(stdout, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			// New output line: "name "description"" or just "name"
			parts := strings.Fields(line)
			if len(parts) == 0 {
				continue
			}
			name := parts[0]
			outputs = append(outputs, DisplayOutput{
				Name:      name,
				Connected: false,
			})
			continue
		}

		if len(outputs) == 0 {
			continue
		}

		idx := len(outputs) - 1
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Enabled:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "Enabled:"))
			outputs[idx].Connected = val == "yes"
			continue
		}
		if strings.HasPrefix(line, "Mode:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "Mode:"))
			// Format: 1920x1080@60.000000Hz
			if atIdx := strings.Index(val, "@"); atIdx > 0 {
				outputs[idx].CurrentMode = val[:atIdx]
				outputs[idx].Modes = append(outputs[idx].Modes, DisplayMode{
					Resolution:  val[:atIdx],
					RefreshRate: strings.TrimSuffix(val[atIdx+1:], "Hz"),
				})
			}
			continue
		}
		if strings.HasPrefix(line, "Position:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "Position:"))
			outputs[idx].Position = val
			continue
		}
	}

	return outputs
}

func discoverOutputs(exe module.CLIExecutor, sessionType string) ([]DisplayOutput, error) {
	if sessionType == "x11" {
		stdout, _, exitCode, err := exe.Run("xrandr", "--query")
		if exitCode != 0 || err != nil {
			return nil, fmt.Errorf("xrandr failed: %v", err)
		}
		return parseXrandrOutput(stdout), nil
	}
	if sessionType == "wayland" {
		stdout, _, exitCode, err := exe.Run("wlr-randr")
		if exitCode != 0 || err != nil {
			return nil, fmt.Errorf("wlr-randr failed: %v", err)
		}
		return parseWlrRandrOutput(stdout), nil
	}
	return nil, fmt.Errorf("unsupported session type: %s", sessionType)
}

func handleRun(settingsPath string) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("run", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	sessionType, err := detectSessionType(executor)
	if err != nil {
		emitError("run", fmt.Sprintf("detect session: %v", err), "")
		return
	}

	outputs, err := discoverOutputs(executor, sessionType)
	if err != nil {
		emitError("run", fmt.Sprintf("discover outputs: %v", err), "")
		return
	}

	// Build select options for set-mode: "output: resolution@refresh"
	var modeOptions []string
	for _, out := range outputs {
		if !out.Connected {
			continue
		}
		for _, m := range out.Modes {
			modeOptions = append(modeOptions, fmt.Sprintf("%s: %s@%s", out.Name, m.Resolution, m.RefreshRate))
		}
	}

	// Build profile names for load-profile.
	var profileNames []string
	for _, p := range s.Display.Profiles {
		profileNames = append(profileNames, p.Name)
	}

	data := map[string]interface{}{
		"session_type":   sessionType,
		"outputs":        outputs,
		"active_profile": s.Display.ActiveProfile,
		"available_options": map[string]interface{}{
			"set-mode":     modeOptions,
			"load-profile": profileNames,
		},
		"current_value": map[string]interface{}{
			"load-profile": s.Display.ActiveProfile,
		},
	}

	resp := module.Response{
		Status:  "ok",
		Action:  "run",
		Data:    data,
		Message: "Display configuration loaded",
	}
	emit(resp)
}

func handleSetMode(settingsPath, value string) {
	if value == "" {
		emitError("set-mode", "mode value is required (format: output: resolution@rate)", "")
		return
	}

	sessionType, err := detectSessionType(executor)
	if err != nil {
		emitError("set-mode", fmt.Sprintf("detect session: %v", err), "")
		return
	}

	// Parse value: "output: resolution@rate" or just "output:resolution@rate"
	value = strings.TrimSpace(value)
	colonIdx := strings.Index(value, ":")
	if colonIdx < 0 {
		emitError("set-mode", "invalid mode format", "expected output: resolution@rate")
		return
	}
	outputName := strings.TrimSpace(value[:colonIdx])
	modePart := strings.TrimSpace(value[colonIdx+1:])

	atIdx := strings.Index(modePart, "@")
	var resolution, rate string
	if atIdx < 0 {
		resolution = modePart
	} else {
		resolution = modePart[:atIdx]
		rate = modePart[atIdx+1:]
	}

	// Validate against discovered outputs.
	outputs, err := discoverOutputs(executor, sessionType)
	if err != nil {
		emitError("set-mode", fmt.Sprintf("discover outputs: %v", err), "")
		return
	}
	valid := false
	for _, out := range outputs {
		if out.Name == outputName {
			for _, m := range out.Modes {
				if m.Resolution == resolution && (rate == "" || m.RefreshRate == rate) {
					valid = true
					break
				}
			}
		}
	}
	if !valid {
		emitError("set-mode", fmt.Sprintf("mode %s@%s for output %s is not available", resolution, rate, outputName), "")
		return
	}

	// Apply.
	if sessionType == "x11" {
		args := []string{"--output", outputName, "--mode", resolution}
		if rate != "" {
			args = append(args, "--rate", rate)
		}
		_, _, exitCode, err := executor.Run("xrandr", args...)
		if exitCode != 0 || err != nil {
			emitError("set-mode", fmt.Sprintf("xrandr failed: %v", err), "")
			return
		}
	} else if sessionType == "wayland" {
		modeArg := resolution
		if rate != "" {
			modeArg = fmt.Sprintf("%s@%s", resolution, rate)
		}
		_, _, exitCode, err := executor.Run("wlr-randr", "--output", outputName, "--mode", modeArg)
		if exitCode != 0 || err != nil {
			emitError("set-mode", fmt.Sprintf("wlr-randr failed: %v", err), "")
			return
		}
	}

	// Persist the mode in a pseudo-field (not in schema, but we note it in settings).
	// Since DisplaySettings doesn't have a per-output runtime store, we save into
	// the active profile if it exists, otherwise just return ok without delta.
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("set-mode", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	// Find or create profile entry.
	profileIdx := -1
	for i, p := range s.Display.Profiles {
		if p.Name == s.Display.ActiveProfile {
			profileIdx = i
			break
		}
	}
	if profileIdx < 0 {
		// No active profile — just return ok without persisting mode.
		resp := module.Response{
			Status:  "ok",
			Action:  "set-mode",
			Message: fmt.Sprintf("Mode set to %s@%s on %s", resolution, rate, outputName),
		}
		emit(resp)
		return
	}

	// Update output config in the active profile.
	outputsCfg := s.Display.Profiles[profileIdx].Outputs
	found := false
	for i, o := range outputsCfg {
		if o.Name == outputName {
			outputsCfg[i].Mode = resolution
			if rate != "" {
				outputsCfg[i].Mode = fmt.Sprintf("%s@%s", resolution, rate)
			}
			found = true
			break
		}
	}
	if !found {
		outputsCfg = append(outputsCfg, settings.OutputConfig{
			Name: outputName,
			Mode: resolution,
		})
	}
	if rate != "" {
		for i := range outputsCfg {
			if outputsCfg[i].Name == outputName {
				outputsCfg[i].Mode = fmt.Sprintf("%s@%s", resolution, rate)
			}
		}
	}

	profiles := make([]map[string]interface{}, len(s.Display.Profiles))
	for i, p := range s.Display.Profiles {
		outputsMaps := make([]map[string]interface{}, len(p.Outputs))
		for j, o := range p.Outputs {
			outputsMaps[j] = map[string]interface{}{
				"name":     o.Name,
				"mode":     o.Mode,
				"position": o.Position,
				"primary":  o.Primary,
			}
		}
		profiles[i] = map[string]interface{}{
			"name":    p.Name,
			"outputs": outputsMaps,
		}
	}
	// Convert updated outputsCfg to []map[string]interface{} for the delta.
	updatedOutputsMaps := make([]map[string]interface{}, len(outputsCfg))
	for j, o := range outputsCfg {
		updatedOutputsMaps[j] = map[string]interface{}{
			"name":     o.Name,
			"mode":     o.Mode,
			"position": o.Position,
			"primary":  o.Primary,
		}
	}
	profiles[profileIdx]["outputs"] = updatedOutputsMaps

	delta := map[string]interface{}{
		"display": map[string]interface{}{
			"profiles": profiles,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-mode", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-mode",
		SettingsDelta: delta,
		Message:       fmt.Sprintf("Mode set to %s@%s on %s", resolution, rate, outputName),
	}
	emit(resp)
}

func handleSetPosition(settingsPath, value string) {
	if value == "" {
		emitError("set-position", "position value is required", "")
		return
	}

	sessionType, err := detectSessionType(executor)
	if err != nil {
		emitError("set-position", fmt.Sprintf("detect session: %v", err), "")
		return
	}

	outputs, err := discoverOutputs(executor, sessionType)
	if err != nil {
		emitError("set-position", fmt.Sprintf("discover outputs: %v", err), "")
		return
	}

	if len(outputs) == 0 {
		emitError("set-position", "no outputs available", "")
		return
	}

	// Pick the first connected output if not specified in value.
	outputName := ""
	for _, out := range outputs {
		if out.Connected {
			outputName = out.Name
			break
		}
	}
	if outputName == "" {
		emitError("set-position", "no connected outputs", "")
		return
	}

	if sessionType == "x11" {
		// value may be "--left-of HDMI-1" or "1920x0" or similar.
		parts := strings.Fields(value)
		args := []string{"--output", outputName}
		args = append(args, parts...)
		_, _, exitCode, err := executor.Run("xrandr", args...)
		if exitCode != 0 || err != nil {
			emitError("set-position", fmt.Sprintf("xrandr failed: %v", err), "")
			return
		}
	} else if sessionType == "wayland" {
		// wlr-randr expects --pos x,y
		args := []string{"--output", outputName, "--pos", value}
		_, _, exitCode, err := executor.Run("wlr-randr", args...)
		if exitCode != 0 || err != nil {
			emitError("set-position", fmt.Sprintf("wlr-randr failed: %v", err), "")
			return
		}
	}

	resp := module.Response{
		Status:  "ok",
		Action:  "set-position",
		Message: fmt.Sprintf("Position set to %s on %s", value, outputName),
	}
	emit(resp)
}

func handleSetPrimary(settingsPath string, value interface{}) {
	sessionType, err := detectSessionType(executor)
	if err != nil {
		emitError("set-primary", fmt.Sprintf("detect session: %v", err), "")
		return
	}

	if sessionType == "wayland" {
		emitError("set-primary", "primary output is not supported on Wayland", "")
		return
	}

	outputs, err := discoverOutputs(executor, sessionType)
	if err != nil {
		emitError("set-primary", fmt.Sprintf("discover outputs: %v", err), "")
		return
	}

	requestedOutput := ""
	switch v := value.(type) {
	case bool:
		if !v {
			emitError("set-primary", "Cannot unset primary display — at least one output must be primary. Use set-primary to assign a different output.", "")
			return
		}
	case string:
		requestedOutput = v
	}

	var target string
	if requestedOutput != "" {
		valid := false
		for _, out := range outputs {
			if out.Connected && out.Name == requestedOutput {
				valid = true
				if out.Primary {
					emitError("set-primary", fmt.Sprintf("%s is already primary", out.Name), "")
					return
				}
				target = out.Name
				break
			}
		}
		if !valid {
			emitError("set-primary", fmt.Sprintf("output %q is not connected", requestedOutput), "")
			return
		}
	} else {
		for _, out := range outputs {
			if out.Connected && !out.Primary {
				target = out.Name
				break
			}
		}
		if target == "" {
			emitError("set-primary", "no connected non-primary outputs found", "")
			return
		}
	}

	_, _, exitCode, err := executor.Run("xrandr", "--output", target, "--primary")
	if exitCode != 0 || err != nil {
		emitError("set-primary", fmt.Sprintf("xrandr failed: %v", err), "")
		return
	}

	resp := module.Response{
		Status:  "ok",
		Action:  "set-primary",
		Message: fmt.Sprintf("Primary output set to %s", target),
	}
	emit(resp)
}

func handleSaveProfile(settingsPath, name string) {
	if name == "" {
		emitError("save-profile", "profile name is required", "")
		return
	}

	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("save-profile", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	sessionType, err := detectSessionType(executor)
	if err != nil {
		emitError("save-profile", fmt.Sprintf("detect session: %v", err), "")
		return
	}

	outputs, err := discoverOutputs(executor, sessionType)
	if err != nil {
		emitError("save-profile", fmt.Sprintf("discover outputs: %v", err), "")
		return
	}

	var configs []settings.OutputConfig
	for _, out := range outputs {
		if !out.Connected {
			continue
		}
		mode := out.CurrentMode
		if mode == "" && len(out.Modes) > 0 {
			mode = out.Modes[0].Resolution
		}
		configs = append(configs, settings.OutputConfig{
			Name:     out.Name,
			Mode:     mode,
			Position: out.Position,
			Primary:  out.Primary,
		})
	}

	newProfile := settings.OutputProfile{
		Name:    name,
		Outputs: configs,
	}

	// Update or append profile.
	found := false
	for i, p := range s.Display.Profiles {
		if p.Name == name {
			s.Display.Profiles[i] = newProfile
			found = true
			break
		}
	}
	if !found {
		s.Display.Profiles = append(s.Display.Profiles, newProfile)
	}

	profiles := make([]map[string]interface{}, len(s.Display.Profiles))
	for i, p := range s.Display.Profiles {
		outputsMaps := make([]map[string]interface{}, len(p.Outputs))
		for j, o := range p.Outputs {
			outputsMaps[j] = map[string]interface{}{
				"name":     o.Name,
				"mode":     o.Mode,
				"position": o.Position,
				"primary":  o.Primary,
			}
		}
		profiles[i] = map[string]interface{}{
			"name":    p.Name,
			"outputs": outputsMaps,
		}
	}

	delta := map[string]interface{}{
		"display": map[string]interface{}{
			"active_profile": name,
			"profiles":       profiles,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("save-profile", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "save-profile",
		SettingsDelta: delta,
		Message:       fmt.Sprintf("Profile %s saved", name),
	}
	emit(resp)
}

func handleLoadProfile(settingsPath, name string) {
	if name == "" {
		emitError("load-profile", "profile name is required", "")
		return
	}

	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("load-profile", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	var profile *settings.OutputProfile
	for _, p := range s.Display.Profiles {
		if p.Name == name {
			profile = &p
			break
		}
	}
	if profile == nil {
		emitError("load-profile", fmt.Sprintf("profile %q not found", name), "")
		return
	}

	sessionType, err := detectSessionType(executor)
	if err != nil {
		emitError("load-profile", fmt.Sprintf("detect session: %v", err), "")
		return
	}

	// Discover currently connected outputs.
	currentOutputs, err := discoverOutputs(executor, sessionType)
	if err != nil {
		emitError("load-profile", fmt.Sprintf("discover outputs: %v", err), "")
		return
	}
	connectedMap := make(map[string]bool)
	for _, out := range currentOutputs {
		if out.Connected {
			connectedMap[out.Name] = true
		}
	}

	var warnings []string
	for _, cfg := range profile.Outputs {
		if !connectedMap[cfg.Name] {
			warnings = append(warnings, fmt.Sprintf("output %s is not connected", cfg.Name))
			continue
		}
		if sessionType == "x11" {
			args := []string{"--output", cfg.Name}
			if cfg.Mode != "" {
				args = append(args, "--mode", cfg.Mode)
			}
			if cfg.Position != "" {
				if strings.Contains(cfg.Position, ",") {
					args = append(args, "--pos", cfg.Position)
				} else {
					args = append(args, strings.Fields(cfg.Position)...)
				}
			}
			if cfg.Primary {
				args = append(args, "--primary")
			}
			_, _, exitCode, err := executor.Run("xrandr", args...)
			if exitCode != 0 || err != nil {
				warnings = append(warnings, fmt.Sprintf("xrandr failed for %s: %v", cfg.Name, err))
			}
		} else if sessionType == "wayland" {
			args := []string{"--output", cfg.Name}
			if cfg.Mode != "" {
				args = append(args, "--mode", cfg.Mode)
			}
			if cfg.Position != "" {
				args = append(args, "--pos", cfg.Position)
			}
			_, _, exitCode, err := executor.Run("wlr-randr", args...)
			if exitCode != 0 || err != nil {
				warnings = append(warnings, fmt.Sprintf("wlr-randr failed for %s: %v", cfg.Name, err))
			}
		}
	}

	delta := map[string]interface{}{
		"display": map[string]interface{}{
			"active_profile": name,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("load-profile", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	msg := fmt.Sprintf("Profile %s loaded", name)
	if len(warnings) > 0 {
		msg += "; warnings: " + strings.Join(warnings, ", ")
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "load-profile",
		SettingsDelta: delta,
		Message:       msg,
	}
	emit(resp)
}

func emit(resp module.Response) {
	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal response: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func emitError(action, message, suggestion string) {
	resp := module.Response{
		Status:     "error",
		Action:     action,
		Message:    message,
		Suggestion: suggestion,
	}
	emit(resp)
}
