package icons

import (
	"os"
	"strings"
)

// DetectNerdFonts reads the LAMBDA_NERD_FONTS environment variable.
// It returns true when the variable is set to "1", "true", or "yes"
// (case-insensitive). All other values, including unset, default to false.
func DetectNerdFonts() bool {
	val := strings.ToLower(strings.TrimSpace(os.Getenv("LAMBDA_NERD_FONTS")))
	return val == "1" || val == "true" || val == "yes"
}

// ResolveNerdFonts returns the effective nerd-fonts setting.
// Priority: explicit CLI flag > LAMBDA_NERD_FONTS env var > default false.
func ResolveNerdFonts(cliSet bool, cliValue bool) bool {
	if cliSet {
		return cliValue
	}
	return DetectNerdFonts()
}
