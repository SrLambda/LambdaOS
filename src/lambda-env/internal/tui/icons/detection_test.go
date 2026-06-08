package icons

import (
	"testing"
)

func TestResolveNerdFonts(t *testing.T) {
	tests := []struct {
		name      string
		cliSet    bool
		cliValue  bool
		env       map[string]string
		want      bool
	}{
		{
			name:   "CLI flag true overrides env false",
			cliSet: true, cliValue: true,
			env:  map[string]string{"LAMBDA_NERD_FONTS": "false"},
			want: true,
		},
		{
			name:   "CLI flag false overrides env true",
			cliSet: true, cliValue: false,
			env:  map[string]string{"LAMBDA_NERD_FONTS": "true"},
			want: false,
		},
		{
			name:   "no CLI flag falls back to env true",
			cliSet: false, cliValue: false,
			env:  map[string]string{"LAMBDA_NERD_FONTS": "true"},
			want: true,
		},
		{
			name:   "no CLI flag falls back to env false",
			cliSet: false, cliValue: false,
			env:  map[string]string{"LAMBDA_NERD_FONTS": "false"},
			want: false,
		},
		{
			name:   "no CLI flag no env defaults to false",
			cliSet: false, cliValue: false,
			env:  map[string]string{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			got := ResolveNerdFonts(tt.cliSet, tt.cliValue)
			if got != tt.want {
				t.Errorf("ResolveNerdFonts(%v, %v) = %v, want %v", tt.cliSet, tt.cliValue, got, tt.want)
			}
		})
	}
}

func TestDetectNerdFonts(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want bool
	}{
		{
			name: "unset defaults to false",
			env:  map[string]string{},
			want: false,
		},
		{
			name: "1 enables nerd fonts",
			env:  map[string]string{"LAMBDA_NERD_FONTS": "1"},
			want: true,
		},
		{
			name: "true enables nerd fonts",
			env:  map[string]string{"LAMBDA_NERD_FONTS": "true"},
			want: true,
		},
		{
			name: "yes enables nerd fonts",
			env:  map[string]string{"LAMBDA_NERD_FONTS": "yes"},
			want: true,
		},
		{
			name: "0 disables nerd fonts",
			env:  map[string]string{"LAMBDA_NERD_FONTS": "0"},
			want: false,
		},
		{
			name: "false disables nerd fonts",
			env:  map[string]string{"LAMBDA_NERD_FONTS": "false"},
			want: false,
		},
		{
			name: "no disables nerd fonts",
			env:  map[string]string{"LAMBDA_NERD_FONTS": "no"},
			want: false,
		},
		{
			name: "random string disables nerd fonts",
			env:  map[string]string{"LAMBDA_NERD_FONTS": "maybe"},
			want: false,
		},
		{
			name: "case insensitive true",
			env:  map[string]string{"LAMBDA_NERD_FONTS": "TRUE"},
			want: true,
		},
		{
			name: "case insensitive yes",
			env:  map[string]string{"LAMBDA_NERD_FONTS": "YES"},
			want: true,
		},
		{
			name: "whitespace trimmed",
			env:  map[string]string{"LAMBDA_NERD_FONTS": "  true  "},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			got := DetectNerdFonts()
			if got != tt.want {
				t.Errorf("DetectNerdFonts() = %v, want %v", got, tt.want)
			}
		})
	}
}
