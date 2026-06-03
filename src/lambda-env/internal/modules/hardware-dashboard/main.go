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

// sysBatteryPath is overridable for tests.
var sysBatteryPath = "/sys/class/power_supply/BAT0/uevent"

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

	switch action {
	case "run":
		handleRun(settingsPath)
	default:
		emitError(action, "unknown action", "use run")
	}
}

func handleRun(settingsPath string) {
	_, err := settings.Load(settingsPath)
	if err != nil {
		emitError("run", fmt.Sprintf("load settings: %v", err), "")
		return
	}

	data := make(map[string]interface{})

	// CPU info (independent from other sections)
	cpuData := collectCPU(executor)
	if cpuData != nil {
		data["cpu"] = cpuData
	}

	// RAM info
	ramData := collectRAM(executor)
	if ramData != nil {
		data["ram"] = ramData
	}

	// Disk info
	diskData := collectDisk(executor)
	if diskData != nil {
		data["disk"] = diskData
	}

	// Temperatures
	tempData := collectTemps(executor)
	if tempData != nil {
		data["temperatures"] = tempData
	}

	// Battery
	batteryData := collectBattery(executor)
	if batteryData != nil {
		data["battery"] = batteryData
	}

	// Uptime
	uptimeData := collectUptime(executor)
	if uptimeData != nil {
		data["uptime"] = uptimeData
	}

	resp := module.Response{
		Status:  "ok",
		Action:  "run",
		Data:    data,
		Message: "Hardware dashboard refreshed",
	}
	emit(resp)
}

// collectCPU gathers CPU model, cores, and load. Failures are isolated.
func collectCPU(exe module.CLIExecutor) map[string]interface{} {
	result := map[string]interface{}{
		"model":     "N/A",
		"cores":     "N/A",
		"load_1min": "N/A",
	}

	// Model name and cores via lscpu
	stdout, _, exitCode, _ := exe.Run("lscpu")
	if exitCode == 0 && stdout != "" {
		for _, line := range strings.Split(stdout, "\n") {
			if strings.Contains(line, "Model name") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					result["model"] = strings.TrimSpace(parts[1])
				}
			}
			if strings.HasPrefix(line, "CPU(s):") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					result["cores"] = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	// Cores fallback via nproc
	if result["cores"] == "N/A" {
		stdout, _, exitCode, _ = exe.Run("nproc")
		if exitCode == 0 {
			result["cores"] = strings.TrimSpace(stdout)
		}
	}

	// Load via /proc/loadavg
	loadData, err := os.ReadFile("/proc/loadavg")
	if err == nil {
		fields := strings.Fields(string(loadData))
		if len(fields) > 0 {
			result["load_1min"] = fields[0]
		}
		if len(fields) > 1 {
			result["load_5min"] = fields[1]
		}
		if len(fields) > 2 {
			result["load_15min"] = fields[2]
		}
	}

	// Load fallback via top
	if result["load_1min"] == "N/A" {
		stdout, _, exitCode, _ = exe.Run("top", "-bn1")
		if exitCode == 0 {
			re := regexp.MustCompile(`load average[s]?:\s*([\d.]+)`)
			if m := re.FindStringSubmatch(stdout); len(m) > 1 {
				result["load_1min"] = m[1]
			}
		}
	}

	return result
}

// collectRAM gathers RAM total, used, available. Failures are isolated.
func collectRAM(exe module.CLIExecutor) map[string]interface{} {
	result := map[string]interface{}{
		"total_mb":     "N/A",
		"used_mb":      "N/A",
		"available_mb": "N/A",
		"percentage":   "N/A",
	}

	stdout, _, exitCode, _ := exe.Run("free", "-m")
	if exitCode == 0 && stdout != "" {
		for _, line := range strings.Split(stdout, "\n") {
			fields := strings.Fields(line)
			if len(fields) >= 4 && fields[0] == "Mem:" {
				total := fields[1]
				used := fields[2]
				available := fields[3]
				result["total_mb"] = total
				result["used_mb"] = used
				result["available_mb"] = available

				totalNum, tErr := strconv.ParseFloat(total, 64)
				usedNum, uErr := strconv.ParseFloat(used, 64)
				if tErr == nil && uErr != nil && totalNum > 0 {
					pct := (usedNum / totalNum) * 100
					result["percentage"] = fmt.Sprintf("%.1f", pct)
				}
				break
			}
		}
	}

	// Fallback to /proc/meminfo
	if result["total_mb"] == "N/A" {
		memData, err := os.ReadFile("/proc/meminfo")
		if err == nil {
			var totalKB, availableKB int
			for _, line := range strings.Split(string(memData), "\n") {
				if strings.HasPrefix(line, "MemTotal:") {
					fields := strings.Fields(line)
					if len(fields) >= 2 {
						totalKB, _ = strconv.Atoi(fields[1])
					}
				}
				if strings.HasPrefix(line, "MemAvailable:") {
					fields := strings.Fields(line)
					if len(fields) >= 2 {
						availableKB, _ = strconv.Atoi(fields[1])
					}
				}
			}
			if totalKB > 0 {
				result["total_mb"] = strconv.Itoa(totalKB / 1024)
				usedKB := totalKB - availableKB
				result["used_mb"] = strconv.Itoa(usedKB / 1024)
				result["available_mb"] = strconv.Itoa(availableKB / 1024)
				pct := float64(usedKB) / float64(totalKB) * 100
				result["percentage"] = fmt.Sprintf("%.1f", pct)
			}
		}
	}

	return result
}

// collectDisk gathers disk usage. Failures are isolated.
func collectDisk(exe module.CLIExecutor) map[string]interface{} {
	result := map[string]interface{}{
		"filesystem": "N/A",
		"size":       "N/A",
		"used":       "N/A",
		"available":  "N/A",
		"use_pct":    "N/A",
	}

	stdout, _, exitCode, _ := exe.Run("df", "-h", "/")
	if exitCode == 0 && stdout != "" {
		lines := strings.Split(stdout, "\n")
		if len(lines) >= 2 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 6 {
				result["filesystem"] = fields[0]
				result["size"] = fields[1]
				result["used"] = fields[2]
				result["available"] = fields[3]
				result["use_pct"] = fields[4]
			}
		}
	}

	return result
}

// collectTemps gathers CPU/GPU temperatures. Failures are isolated.
func collectTemps(exe module.CLIExecutor) map[string]interface{} {
	result := make(map[string]interface{})

	// Try sensors first.
	stdout, _, exitCode, _ := exe.Run("sensors")
	if exitCode == 0 && stdout != "" {
		tempRe := regexp.MustCompile(`([\w\s]+):\s*\+?([\d.]+)\s*°C`)
		for _, line := range strings.Split(stdout, "\n") {
			if m := tempRe.FindStringSubmatch(line); len(m) > 2 {
				label := strings.TrimSpace(m[1])
				val := m[2]
				result[label] = val
			}
		}
	}

	// Fallback to thermal zones.
	if len(result) == 0 {
		zones, err := filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
		if err == nil && len(zones) > 0 {
			for _, zonePath := range zones {
				data, err := os.ReadFile(zonePath)
				if err != nil {
					continue
				}
				tempStr := strings.TrimSpace(string(data))
				tempFloat, err := strconv.ParseFloat(tempStr, 64)
				if err != nil {
					continue
				}
				// temp is usually in millidegrees
				tempC := tempFloat / 1000.0

				// Try to read type for labeling
				dir := filepath.Dir(zonePath)
				typeData, _ := os.ReadFile(filepath.Join(dir, "type"))
				label := strings.TrimSpace(string(typeData))
				if label == "" {
					label = filepath.Base(dir)
				}
				result[label] = fmt.Sprintf("%.1f", tempC)
			}
		}
	}

	if len(result) == 0 {
		result["status"] = "N/A"
	}

	return result
}

// collectBattery gathers battery status. Failures are isolated.
func collectBattery(exe module.CLIExecutor) map[string]interface{} {
	result := make(map[string]interface{})

	// Try upower first.
	stdout, _, exitCode, _ := exe.Run("upower", "-d")
	if exitCode == 0 && stdout != "" {
		stateRe := regexp.MustCompile(`state:\s*(\w+)`)
		pctRe := regexp.MustCompile(`percentage:\s*(\d+)%`)
		timeRe := regexp.MustCompile(`time to empty:\s*([\d.]+\s*\w+)`)

		for _, line := range strings.Split(stdout, "\n") {
			if m := stateRe.FindStringSubmatch(line); len(m) > 1 {
				result["state"] = m[1]
			}
			if m := pctRe.FindStringSubmatch(line); len(m) > 1 {
				result["percentage"] = m[1]
			}
			if m := timeRe.FindStringSubmatch(line); len(m) > 1 {
				result["time_remaining"] = m[1]
			}
		}
	}

	// If upower gave nothing useful, fallback to sysfs.
	if len(result) == 0 {
		data, err := os.ReadFile(sysBatteryPath)
		if err == nil {
			var capacity, status string
			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(line, "POWER_SUPPLY_CAPACITY=") {
					capacity = strings.TrimPrefix(line, "POWER_SUPPLY_CAPACITY=")
				}
				if strings.HasPrefix(line, "POWER_SUPPLY_STATUS=") {
					status = strings.TrimPrefix(line, "POWER_SUPPLY_STATUS=")
				}
			}
			if capacity != "" {
				result["percentage"] = capacity
			}
			if status != "" {
				result["state"] = strings.ToLower(status)
			}
		}
	}

	if len(result) == 0 {
		result["status"] = "no battery detected"
	}

	return result
}

// collectUptime gathers system uptime. Failures are isolated.
func collectUptime(exe module.CLIExecutor) map[string]interface{} {
	result := map[string]interface{}{
		"human_readable": "N/A",
		"seconds":        "N/A",
	}

	// Try uptime -p first.
	stdout, _, exitCode, _ := exe.Run("uptime", "-p")
	if exitCode == 0 {
		result["human_readable"] = strings.TrimSpace(stdout)
	}

	// Read /proc/uptime for seconds.
	uptimeData, err := os.ReadFile("/proc/uptime")
	if err == nil {
		fields := strings.Fields(string(uptimeData))
		if len(fields) > 0 {
			result["seconds"] = fields[0]
		}
	}

	return result
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
