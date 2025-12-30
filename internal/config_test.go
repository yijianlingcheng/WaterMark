package internal

import (
	"os"
	"path/filepath"
	"testing"

	"WaterMark/pkg"
)

func TestSetAppMode(t *testing.T) {
	tests := []struct {
		name string
		mode string
	}{
		{
			name: "set debug mode",
			mode: APP_DEV,
		},
		{
			name: "set api debug mode",
			mode: APP_API_DEV,
		},
		{
			name: "set release mode",
			mode: APP_RELEASE,
		},
		{
			name: "set custom mode",
			mode: "custom",
		},
		{
			name: "set empty mode",
			mode: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetAppMode(tt.mode)
			if appMode != tt.mode {
				t.Errorf("SetAppMode(%q) = %q, want %q", tt.mode, appMode, tt.mode)
			}
		})
	}
}

func TestISRelease(t *testing.T) {
	tests := []struct {
		name     string
		mode     string
		expected bool
	}{
		{
			name:     "release mode",
			mode:     APP_RELEASE,
			expected: true,
		},
		{
			name:     "debug mode",
			mode:     APP_DEV,
			expected: false,
		},
		{
			name:     "api debug mode",
			mode:     APP_API_DEV,
			expected: false,
		},
		{
			name:     "custom mode",
			mode:     "custom",
			expected: false,
		},
		{
			name:     "empty mode",
			mode:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetAppMode(tt.mode)
			result := ISRelease()
			if result != tt.expected {
				t.Errorf("ISRelease() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestISApiDebug(t *testing.T) {
	tests := []struct {
		name     string
		mode     string
		expected bool
	}{
		{
			name:     "api debug mode",
			mode:     APP_API_DEV,
			expected: true,
		},
		{
			name:     "release mode",
			mode:     APP_RELEASE,
			expected: false,
		},
		{
			name:     "debug mode",
			mode:     APP_DEV,
			expected: false,
		},
		{
			name:     "custom mode",
			mode:     "custom",
			expected: false,
		},
		{
			name:     "empty mode",
			mode:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetAppMode(tt.mode)
			result := ISApiDebug()
			if result != tt.expected {
				t.Errorf("ISApiDebug() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRestoreAppConfigFile(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
	}{
		{
			name: "file exists",
			setup: func() string {
				tmpDir := t.TempDir()
				testFile := filepath.Join(tmpDir, "app.yaml")
				os.WriteFile(testFile, []byte("test"), 0o644)
				return testFile
			},
			wantErr: false,
		},
		{
			name: "file does not exist",
			setup: func() string {
				tmpDir := t.TempDir()
				testFile := filepath.Join(tmpDir, "app.yaml")
				return testFile
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPath := tt.setup()

			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("restoreAppConfigFile() panicked unexpectedly: %v", r)
				}
			}()

			restoreAppConfigFile(testPath)
		})
	}
}

func TestGetPlugin(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		expected string
	}{
		{
			name:     "default plugin",
			setup:    func() {},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result := GetPlugin()
			if result != tt.expected {
				t.Logf("GetPlugin() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestAppModeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "APP_DEV constant",
			constant: APP_DEV,
			expected: "debug",
		},
		{
			name:     "APP_API_DEV constant",
			constant: APP_API_DEV,
			expected: "api_debug",
		},
		{
			name:     "APP_RELEASE constant",
			constant: APP_RELEASE,
			expected: "release",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Constant %q = %q, want %q", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

func TestModeTransitions(t *testing.T) {
	tests := []struct {
		name     string
		modes    []string
		expected []bool
		check    func() bool
	}{
		{
			name:  "release to debug",
			modes: []string{APP_RELEASE, APP_DEV},
			check: func() bool { return !ISRelease() },
		},
		{
			name:  "debug to api debug",
			modes: []string{APP_DEV, APP_API_DEV},
			check: func() bool { return ISApiDebug() },
		},
		{
			name:  "api debug to release",
			modes: []string{APP_API_DEV, APP_RELEASE},
			check: func() bool { return ISRelease() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, mode := range tt.modes {
				SetAppMode(mode)
			}

			if !tt.check() {
				t.Errorf("Mode transition failed")
			}
		})
	}
}

func TestInitAppConfig(t *testing.T) {
	InitAppConfigsAndRes()

	tests := []struct {
		name    string
		setup   func() func()
		wantErr bool
	}{
		{
			name: "valid config file",
			setup: func() func() {
				tmpDir := t.TempDir()
				configPath := filepath.Join(tmpDir, "app.yaml")
				os.WriteFile(configPath, []byte("test: value\n"), 0o644)
				return func() { os.RemoveAll(tmpDir) }
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup()
			defer cleanup()

			err := initAppConfig()
			if pkg.HasError(err) != tt.wantErr {
				t.Errorf("initAppConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
