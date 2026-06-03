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

var allowedLidActions = []string{"suspend", "hibernate", "ignore", "poweroff"}

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
	case "set-screen-timeout":
		timeoutStr := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				timeoutStr = v
			}
		}
		handleSetScreenTimeout(settingsPath, timeoutStr)
	case "set-sleep-timeout":
		timeoutStr := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				timeoutStr = v
			}
		}
		handleSetSleepTimeout(settingsPath, timeoutStr)
	case "set-lid-close-action":
		actionVal := ""
		if params != nil {
			if v, ok := params["value"].(string); ok {
				actionVal = v
			}
		}
		handleSetLidCloseAction(settingsPath, actionVal)
	default:
		emitError(action, "unknown action", "use run, set-screen-timeout, set-sleep-timeout, or set-lid-close-action")
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

// readLogindConf reads /etc/systemd/logind.conf and extracts IdleActionSec and HandleLidSwitch.
func readLogindConf() (screenTimeout int, lidAction string, err error) {
	data, err := os.ReadFile("/etc/systemd/logind.conf")
	if err != nil {
		return 0, "", err
	}
	return parseLogindConf(string(data))
}

func parseLogindConf(data string) (screenTimeout int, lidAction string, err error) {
	screenRe := regexp.MustCompile(`(?m)^IdleActionSec\s*=\s*(\d+)`)
	lidRe := regexp.MustCompile(`(?m)^HandleLidSwitch\s*=\s*(\w+)`)

	if m := screenRe.FindStringSubmatch(data); len(m) > 1 {
		if v, parseErr := strconv.Atoi(m[1]); parseErr == nil {
			screenTimeout = v
		}
	}
	if m := lidRe.FindStringSubmatch(data); len(m) > 1 {
		lidAction = m[1]
	}
	return screenTimeout, lidAction, nil
}

// updateLogindConfKey replaces or appends a key in logind.conf data.
func updateLogindConfKey(data, key, value string) string {
	re := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(key) + `\s*=.*$`)
	line := key + "=" + value
	if re.MatchString(data) {
		return re.ReplaceAllString(data, line)
	}
	// Append under [Login] or at end.
	if strings.Contains(data, "[Login]") {
		return strings.Replace(data, "[Login]", "[Login]\n"+line, 1)
	}
	return data + "\n[Login]\n" + line + "\n"
}

func readBatteryStatus(exe module.CLIExecutor) (map[string]interface{}, string) {
	// Try upower first.
	stdout, _, exitCode, execErr := exe.Run("upower", "-d")
	if exitCode == 0 && execErr == nil {
		battery := parseUpowerOutput(stdout)
		if battery != nil {
			return battery, ""
		}
	}

	// Fallback to /sys/class/power_supply/BAT0/uevent.
	battery, err := readSysBattery()
	if err == nil && battery != nil {
		return battery, ""
	}

	if exitCode != 0 || execErr != nil {
		return nil, "battery info unavailable (upower not installed or no battery)"
	}
	return nil, "no battery detected"
}

func parseUpowerOutput(stdout string) map[string]interface{} {
	var state, percentage, time string
	stateRe := regexp.MustCompile(`state:\s*(\w+)`)
	pctRe := regexp.MustCompile(`percentage:\s*(\d+)%`)
	timeRe := regexp.MustCompile(`time to empty:\s*([\d.]+\s*\w+)`)

	for _, line := range strings.Split(stdout, "\n") {
		if m := stateRe.FindStringSubmatch(line); len(m) > 1 {
			state = m[1]
		}
		if m := pctRe.FindStringSubmatch(line); len(m) > 1 {
			percentage = m[1]
		}
		if m := timeRe.FindStringSubmatch(line); len(m) > 1 {
			time = m[1]
		}
	}

	if state == "" && percentage == "" {
		return nil
	}

	result := map[string]interface{}{}
	if state != "" {
		result["state"] = state
	}
	if percentage != "" {
		result["percentage"] = percentage
	}
	if time != "" {
		result["time_remaining"] = time
	}
	return result
}

var sysBatteryPath = "/sys/class/power_supply/BAT0/uevent"

func readSysBattery() (map[string]interface{}, error) {
	data, err := os.ReadFile(sysBatteryPath)
	if err != nil {
		return nil, err
	}

	var capacity, status string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "POWER_SUPPLY_CAPACITY=") {
			capacity = strings.TrimPrefix(line, "POWER_SUPPLY_CAPACITY=")
		}
		if strings.HasPrefix(line, "POWER_SUPPLY_STATUS=") {
			status = strings.TrimPrefix(line, "POWER_SUPPLY_STATUS=")
		}
	}

	if capacity == "" && status == "" {
		return nil, fmt.Errorf("no battery data found")
	}

	result := map[string]interface{}{}
	if capacity != "" {
		result["percentage"] = capacity
	}
	if status != "" {
		result["state"] = strings.ToLower(status)
	}
	return result, nil
}

func handleRun(settingsPath string) {
	s, err := settings.Load(settingsPath)
	if err != nil {
		emitError("run", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	// Read logind.conf for current system values (best-effort).
	logindScreenTimeout, logindLidAction, _ := readLogindConf()
	screenTimeout := s.Power.ScreenTimeout
	sleepTimeout := s.Power.SleepTimeout
	lidAction := s.Power.LidCloseAction

	if logindScreenTimeout > 0 {
		screenTimeout = logindScreenTimeout
	}
	if logindLidAction != "" {
		lidAction = logindLidAction
	}

	battery, batteryWarning := readBatteryStatus(executor)

	data := map[string]interface{}{
		"screen_timeout":  screenTimeout,
		"sleep_timeout":   sleepTimeout,
		"lid_close_action": lidAction,
		"current_value": map[string]interface{}{
			"set-screen-timeout":    screenTimeout,
			"set-sleep-timeout":     sleepTimeout,
			"set-lid-close-action":  lidAction,
		},
		"available_options": map[string]interface{}{
			"set-lid-close-action": allowedLidActions,
		},
	}
	if battery != nil {
		data["battery"] = battery
	}
	if batteryWarning != "" {
		data["battery_warning"] = batteryWarning
	}

	resp := module.Response{
		Status:  "ok",
		Action:  "run",
		Data:    data,
		Message: "Power configuration loaded",
	}
	emit(resp)
}

func handleSetScreenTimeout(settingsPath, timeoutStr string) {
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil || timeout < 0 {
		emitError("set-screen-timeout", fmt.Sprintf("invalid timeout %q", timeoutStr), "expected a non-negative integer in seconds")
		return
	}

	// Best-effort update of logind.conf (requires root).
	data, readErr := os.ReadFile("/etc/systemd/logind.conf")
	if readErr == nil {
		updated := updateLogindConfKey(string(data), "IdleActionSec", strconv.Itoa(timeout))
		_ = os.WriteFile("/etc/systemd/logind.conf", []byte(updated), 0644)
	}

	delta := map[string]interface{}{
		"power": map[string]interface{}{
			"screen_timeout": timeout,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-screen-timeout", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	msg := fmt.Sprintf("Screen timeout set to %d seconds", timeout)
	if readErr != nil {
		msg += " (requires root to apply to system)"
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-screen-timeout",
		SettingsDelta: delta,
		Message:       msg,
	}
	emit(resp)
}

func handleSetSleepTimeout(settingsPath, timeoutStr string) {
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil || timeout < 0 {
		emitError("set-sleep-timeout", fmt.Sprintf("invalid timeout %q", timeoutStr), "expected a non-negative integer in seconds")
		return
	}

	// Note: there is no standard standalone systemd key for sleep timeout separate from
	// idle action. We persist the value in settings and document that applying it to
	// the system may require root and manual configuration.

	delta := map[string]interface{}{
		"power": map[string]interface{}{
			"sleep_timeout": timeout,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-sleep-timeout", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-sleep-timeout",
		SettingsDelta: delta,
		Message:       fmt.Sprintf("Sleep timeout set to %d seconds", timeout),
	}
	emit(resp)
}

func handleSetLidCloseAction(settingsPath, actionVal string) {
	valid := false
	for _, a := range allowedLidActions {
		if a == actionVal {
			valid = true
			break
		}
	}
	if !valid {
		emitError("set-lid-close-action", fmt.Sprintf("invalid lid action %q", actionVal), fmt.Sprintf("must be one of: %v", allowedLidActions))
		return
	}

	// Best-effort update of logind.conf (requires root).
	data, readErr := os.ReadFile("/etc/systemd/logind.conf")
	if readErr == nil {
		updated := updateLogindConfKey(string(data), "HandleLidSwitch", actionVal)
		_ = os.WriteFile("/etc/systemd/logind.conf", []byte(updated), 0644)
	}

	delta := map[string]interface{}{
		"power": map[string]interface{}{
			"lid_close_action": actionVal,
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		emitError("set-lid-close-action", fmt.Sprintf("save delta: %v", err), "")
		return
	}

	msg := fmt.Sprintf("Lid close action set to %s", actionVal)
	if readErr != nil {
		msg += " (requires root to apply to system)"
	}

	resp := module.Response{
		Status:        "ok",
		Action:        "set-lid-close-action",
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
