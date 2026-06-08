package theme

import (
	"math"
	"testing"
)

// relativeLuminance computes the WCAG 2.1 relative luminance of a hex color.
func relativeLuminance(hex string) float64 {
	// Parse hex #RRGGBB
	r := float64(hexCharValue(hex[1])<<4|hexCharValue(hex[2])) / 255.0
	g := float64(hexCharValue(hex[3])<<4|hexCharValue(hex[4])) / 255.0
	b := float64(hexCharValue(hex[5])<<4|hexCharValue(hex[6])) / 255.0

	r = gammaCorrect(r)
	g = gammaCorrect(g)
	b = gammaCorrect(b)

	return 0.2126*r + 0.7152*g + 0.0722*b
}

func hexCharValue(c byte) int {
	if c >= '0' && c <= '9' {
		return int(c - '0')
	}
	if c >= 'a' && c <= 'f' {
		return int(c-'a') + 10
	}
	if c >= 'A' && c <= 'F' {
		return int(c-'A') + 10
	}
	return 0
}

func gammaCorrect(c float64) float64 {
	if c <= 0.03928 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

// contrastRatio computes the WCAG contrast ratio between two hex colors.
func contrastRatio(fg, bg string) float64 {
	l1 := relativeLuminance(fg)
	l2 := relativeLuminance(bg)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

func TestWCAGContrast(t *testing.T) {
	backgrounds := []struct {
		name string
		hex  string
	}{
		{"Bg", Bg},
		{"StatusBg", StatusBg},
	}

	// Colors used for normal body text must meet the stricter 4.5:1 threshold.
	normalTextColors := []struct {
		name string
		hex  string
	}{
		{"Success", Success},
		{"Error", Error},
		{"Warn", Warn},
		{"Dimmed", Dimmed},
		{"Accent", Accent},
	}

	// Bold/large text colors that only need >= 3:1.
	boldColors := []struct {
		name string
		hex  string
	}{}

	for _, bg := range backgrounds {
		for _, col := range normalTextColors {
			t.Run(col.name+"_on_"+bg.name, func(t *testing.T) {
				ratio := contrastRatio(col.hex, bg.hex)

				// WCAG AA normal text requires >= 4.5:1
				if ratio < 4.5 {
					t.Errorf(
						"contrast ratio for %s on %s = %.2f:1, want >= 4.5:1",
						col.name, bg.name, ratio,
					)
				}
			})
		}

		for _, col := range boldColors {
			t.Run(col.name+"_on_"+bg.name, func(t *testing.T) {
				ratio := contrastRatio(col.hex, bg.hex)

				// WCAG AA large/bold text requires >= 3:1
				if ratio < 3.0 {
					t.Errorf(
						"contrast ratio for %s on %s = %.2f:1, want >= 3.0:1 for large/bold text",
						col.name, bg.name, ratio,
					)
				}
			})
		}
	}
}

func TestDimmedOnStatusBg(t *testing.T) {
	// Explicitly verify the fixed dimmed color passes AA on StatusBg.
	ratio := contrastRatio(Dimmed, StatusBg)
	if ratio < 4.5 {
		t.Errorf("Dimmed on StatusBg = %.2f:1, want >= 4.5:1", ratio)
	}
}
