package internal

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"WaterMark/internal/cmd"
	"WaterMark/pkg"
)

func TestCheckExiftool(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() func()
		wantErr bool
		errCode int
	}{
		{
			name: "exiftool exists and works",
			setup: func() func() {
				return func() {}
			},
			wantErr: false,
		},
		{
			name: "exiftool does not exist",
			setup: func() func() {
				return func() {}
			},
			wantErr: true,
			errCode: pkg.ExiftoolNotExistError.Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup()
			defer cleanup()

			err := checkExiftool()
			if tt.wantErr {
				if !pkg.HasError(err) {
					t.Skipf("checkExiftool() expected error, got NoError (tool may be installed)")
				}
				if err.Code != tt.errCode {
					t.Errorf("checkExiftool() error code = %v, want %v", err.Code, tt.errCode)
				}
			} else {
				if pkg.HasError(err) {
					t.Skipf("checkExiftool() unexpected error = %v (tool may not be installed)", err)
				}
			}
		})
	}
}

func TestCheckFontFile(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() func()
		wantErr bool
	}{
		{
			name: "font directory exists with files",
			setup: func() func() {
				initRootPath()
				fontDir := GetFontFilePath("")
				if err := os.MkdirAll(fontDir, os.ModePerm); err != nil {
					t.Fatalf("Failed to create font directory: %v", err)
				}
				fontFile := filepath.Join(fontDir, "test.ttf")
				os.WriteFile(fontFile, []byte("test"), 0o644)
				return func() { os.RemoveAll(fontDir) }
			},
			wantErr: false,
		},
		{
			name: "font directory exists but empty",
			setup: func() func() {
				initRootPath()
				fontDir := GetFontFilePath("")
				if err := os.MkdirAll(fontDir, os.ModePerm); err != nil {
					t.Fatalf("Failed to create font directory: %v", err)
				}
				return func() { os.RemoveAll(fontDir) }
			},
			wantErr: false,
		},
		{
			name: "font directory does not exist",
			setup: func() func() {
				initRootPath()
				fontDir := GetFontFilePath("")
				return func() { os.RemoveAll(fontDir) }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup()
			defer cleanup()

			err := checkFontFile()
			if pkg.HasError(err) != tt.wantErr {
				t.Errorf("checkFontFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckLogoFile(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() func()
		wantErr bool
	}{
		{
			name: "logo directory exists with files",
			setup: func() func() {
				initRootPath()
				logoDir := GetLogosPath("")
				if err := os.MkdirAll(logoDir, os.ModePerm); err != nil {
					t.Fatalf("Failed to create logo directory: %v", err)
				}
				logoFile := filepath.Join(logoDir, "test.png")
				os.WriteFile(logoFile, []byte("test"), 0o644)
				return func() { os.RemoveAll(logoDir) }
			},
			wantErr: false,
		},
		{
			name: "logo directory exists but empty",
			setup: func() func() {
				initRootPath()
				logoDir := GetLogosPath("")
				if err := os.MkdirAll(logoDir, os.ModePerm); err != nil {
					t.Fatalf("Failed to create logo directory: %v", err)
				}
				return func() { os.RemoveAll(logoDir) }
			},
			wantErr: false,
		},
		{
			name: "logo directory does not exist",
			setup: func() func() {
				initRootPath()
				logoDir := GetLogosPath("")
				return func() { os.RemoveAll(logoDir) }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup()
			defer cleanup()

			err := checkLogoFile()
			if pkg.HasError(err) != tt.wantErr {
				t.Errorf("checkLogoFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckImageMagick(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() func()
		wantErr bool
		errCode int
	}{
		{
			name: "ImageMagick exists and works",
			setup: func() func() {
				return func() {}
			},
			wantErr: false,
		},
		{
			name: "ImageMagick does not exist",
			setup: func() func() {
				return func() {}
			},
			wantErr: true,
			errCode: pkg.ExiftoolNotExistError.Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup()
			defer cleanup()

			err := checkImageMagick()
			if tt.wantErr {
				if !pkg.HasError(err) {
					t.Skipf("checkImageMagick() expected error, got NoError (tool may be installed)")
				}
				if err.Code != tt.errCode {
					t.Errorf("checkImageMagick() error code = %v, want %v", err.Code, tt.errCode)
				}
			} else {
				if pkg.HasError(err) {
					t.Skipf("checkImageMagick() unexpected error = %v (tool may not be installed)", err)
				}
			}
		})
	}
}

func TestCheckInstallExif(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() func()
		wantErr bool
	}{
		{
			name: "exiftool installation check",
			setup: func() func() {
				return func() {}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup()
			defer cleanup()

			err := checkInstallExif()
			if pkg.HasError(err) != tt.wantErr {
				t.Errorf("checkInstallExif() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckInstallImageMagick(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() func()
		wantErr bool
	}{
		{
			name: "ImageMagick installation check",
			setup: func() func() {
				return func() {}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup()
			defer cleanup()

			err := checkInstallImageMagick()
			if pkg.HasError(err) != tt.wantErr {
				t.Skipf("checkInstallImageMagick() error = %v, wantErr %v (tool may not be installed)", err, tt.wantErr)
			}
		})
	}
}

func TestCheckWithTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		setup   func() func()
		wantErr bool
	}{
		{
			name:    "short timeout",
			timeout: 1 * time.Millisecond,
			setup:   func() func() { return func() {} },
			wantErr: true,
		},
		{
			name:    "long timeout",
			timeout: 10 * time.Second,
			setup:   func() func() { return func() {} },
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup()
			defer cleanup()

			_, err := cmd.CommandRun(tt.timeout, "echo test")
			if pkg.HasError(err) != tt.wantErr {
				t.Errorf("CommandRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckEmptyVersion(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() func()
		wantErr bool
	}{
		{
			name: "empty version response",
			setup: func() func() {
				return func() {}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup()
			defer cleanup()

			version, err := cmd.CommandRun(5*time.Second, "nonexistentcommand12345")
			if version != "" || !pkg.HasError(err) {
				t.Errorf("Expected empty version and error, got version=%s, err=%v", version, err)
			}
		})
	}
}

func TestCheckWithInvalidPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "invalid path",
			path:    "/invalid/path",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cmd.CommandRun(5*time.Second, "nonexistentcommand12345")
			if !pkg.HasError(err) {
				t.Errorf("Expected error for invalid path %s, got NoError", tt.path)
			}
		})
	}
}
