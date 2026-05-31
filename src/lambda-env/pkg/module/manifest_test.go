package module

import (
	"strings"
	"testing"
)

func TestManifestValidate(t *testing.T) {
	tests := []struct {
		name     string
		manifest Manifest
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid manifest",
			manifest: Manifest{
				Name:          "screen",
				Version:       "0.1.0",
				Description:   "Manage display configuration",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "system",
				RequiresRoot:  false,
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			manifest: Manifest{
				Version:       "0.1.0",
				Description:   "Manage display configuration",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "system",
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
			},
			wantErr: true,
			errMsg:  "manifest name is required",
		},
		{
			name: "missing version",
			manifest: Manifest{
				Name:          "screen",
				Description:   "Manage display configuration",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "system",
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
			},
			wantErr: true,
			errMsg:  "manifest version is required",
		},
		{
			name: "invalid category",
			manifest: Manifest{
				Name:          "screen",
				Version:       "0.1.0",
				Description:   "Manage display configuration",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "invalid",
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
			},
			wantErr: true,
			errMsg:  "manifest category",
		},
		{
			name: "empty dependencies and zero timeout",
			manifest: Manifest{
				Name:          "screen",
				Version:       "0.1.0",
				Description:   "Manage display configuration",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "system",
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
				Timeout:       0,
			},
			wantErr: false,
		},
		{
			name: "name with spaces",
			manifest: Manifest{
				Name:          "My Module",
				Version:       "0.1.0",
				Description:   "Manage display configuration",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "system",
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
			},
			wantErr: true,
			errMsg:  "must be lowercase with hyphens only",
		},
		{
			name: "missing description",
			manifest: Manifest{
				Name:          "screen",
				Version:       "0.1.0",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "system",
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
			},
			wantErr: true,
			errMsg:  "manifest description is required",
		},
		{
			name: "missing description_es",
			manifest: Manifest{
				Name:          "screen",
				Version:       "0.1.0",
				Description:   "Manage display configuration",
				Category:      "system",
				Dependencies:  []string{},
				MinHubVersion: "0.1.0",
			},
			wantErr: true,
			errMsg:  "manifest description_es is required",
		},
		{
			name: "missing min_hub_version",
			manifest: Manifest{
				Name:          "screen",
				Version:       "0.1.0",
				Description:   "Manage display configuration",
				DescriptionES: "Gestionar configuración de pantalla",
				Category:      "system",
				Dependencies:  []string{},
			},
			wantErr: true,
			errMsg:  "manifest min_hub_version is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.manifest.Validate()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errMsg)
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Fatalf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}
