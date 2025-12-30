package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetRootPath(t *testing.T) {
	rootPath := getRootPath()
	if rootPath == "" {
		t.Error("getRootPath() returned empty string")
	}

	_, err := os.Stat(rootPath)
	if err != nil {
		t.Errorf("getRootPath() returned invalid path: %v", err)
	}
}

func TestGetFilePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple path",
			input:    "test.txt",
			expected: "test.txt",
		},
		{
			name:     "path with directory",
			input:    "dir/file.txt",
			expected: "dir/file.txt",
		},
		{
			name:     "empty path",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFilePath(tt.input)
			if !strings.HasSuffix(result, tt.expected) {
				t.Errorf("getFilePath(%q) = %q, want suffix %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestReplaceMode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		mode     string
		expected string
	}{
		{
			name:     "dev mode",
			code:     "internal.SetAppMode(internal.APP_RELEASE)",
			mode:     "dev",
			expected: "internal.SetAppMode(internal.APP_DEV)",
		},
		{
			name:     "api mode",
			code:     "internal.SetAppMode(internal.APP_RELEASE)",
			mode:     "api-dev",
			expected: "internal.SetAppMode(internal.APP_API_DEV)",
		},
		{
			name:     "release mode",
			code:     "internal.SetAppMode(internal.APP_DEV)",
			mode:     "release",
			expected: "internal.SetAppMode(internal.APP_RELEASE)",
		},
		{
			name: "multiple occurrences",
			code: "internal.SetAppMode(internal.APP_DEV)\n" +
				"internal.SetAppMode(internal.APP_RELEASE)",
			mode: "api-dev",
			expected: "internal.SetAppMode(internal.APP_API_DEV)\n" +
				"internal.SetAppMode(internal.APP_API_DEV)",
		},
		{
			name:     "code without mode",
			code:     "some other code",
			mode:     "dev",
			expected: "some other code",
		},
		{
			name:     "invalid mode keeps original",
			code:     "internal.SetAppMode(internal.APP_DEV)",
			mode:     "invalid",
			expected: "internal.SetAppMode(internal.APP_RELEASE)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceMode(tt.code, tt.mode)
			if result != tt.expected {
				t.Errorf("replaceMode() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetAppVersion(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		versionFile string
		expected    string
		wantErr     bool
	}{
		{
			name:        "valid version",
			versionFile: "APP_VERSION:1.0.0",
			expected:    "1.0.0",
			wantErr:     false,
		},
		{
			name:        "version with spaces",
			versionFile: "APP_VERSION: 2.0.0 ",
			expected:    " 2.0.0 ",
			wantErr:     false,
		},
		{
			name:        "version with newline",
			versionFile: "APP_VERSION:3.0.0\n",
			expected:    "3.0.0",
			wantErr:     false,
		},
		{
			name:        "version with CRLF",
			versionFile: "APP_VERSION:4.0.0\r\n",
			expected:    "4.0.0",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			versionFile := filepath.Join(tmpDir, "version")
			err := os.WriteFile(versionFile, []byte(tt.versionFile), 0o600)
			if err != nil {
				t.Fatalf("failed to create version file: %v", err)
			}

			originalGetFilePathFunc := getFilePathFunc
			getFilePathFunc = func(p string) string {
				return filepath.Join(tmpDir, p)
			}
			defer func() { getFilePathFunc = originalGetFilePathFunc }()

			result := getAppVersion()
			if result != tt.expected {
				t.Errorf("getAppVersion() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestChangeAppModeToDev(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		mode      string
		inputCode string
		expected  string
		wantErr   bool
	}{
		{
			name: "dev mode",
			mode: "dev",
			inputCode: `package main
import "internal"
func main() {
	internal.SetAppMode(internal.APP_RELEASE)
}`,
			expected: `package main
import "internal"
func main() {
	internal.SetAppMode(internal.APP_DEV)
}`,
			wantErr: false,
		},
		{
			name: "api mode",
			mode: "api-dev",
			inputCode: `package main
import "internal"
func main() {
	internal.SetAppMode(internal.APP_DEV)
}`,
			expected: `package main
import "internal"
func main() {
	internal.SetAppMode(internal.APP_API_DEV)
}`,
			wantErr: false,
		},
		{
			name: "release mode",
			mode: "release",
			inputCode: `package main
import "internal"
func main() {
	internal.SetAppMode(internal.APP_DEV)
}`,
			expected: `package main
import "internal"
func main() {
	internal.SetAppMode(internal.APP_RELEASE)
}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mainFile := filepath.Join(tmpDir, "main.go")
			err := os.WriteFile(mainFile, []byte(tt.inputCode), 0o600)
			if err != nil {
				t.Fatalf("failed to create main.go: %v", err)
			}

			originalGetFilePathFunc := getFilePathFunc
			getFilePathFunc = func(p string) string {
				return filepath.Join(tmpDir, p)
			}
			defer func() { getFilePathFunc = originalGetFilePathFunc }()

			originalPrintErrorFunc := printErrorFunc
			printErrorFunc = func(err error) {
				if !tt.wantErr {
					t.Errorf("unexpected error: %v", err)
				}
			}
			defer func() { printErrorFunc = originalPrintErrorFunc }()

			changeAppModeToDev(tt.mode)

			result, err := os.ReadFile(mainFile)
			if err != nil {
				t.Fatalf("failed to read main.go: %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("changeAppModeToDev() result = %q, want %q", string(result), tt.expected)
			}
		})
	}
}

func TestReplaceAppVersion(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		versionFile   string
		htmlContent   string
		expectedMatch bool
	}{
		{
			name:          "replace version",
			versionFile:   "APP_VERSION:1.2.3",
			htmlContent:   `<html><body><p>版本: 1.0.0</p></body></html>`,
			expectedMatch: true,
		},
		{
			name:          "version with case insensitive",
			versionFile:   "APP_VERSION:2.0.0",
			htmlContent:   `<html><body><P>版本: old</P></body></html>`,
			expectedMatch: true,
		},
		{
			name:          "version with spaces",
			versionFile:   "APP_VERSION:3.0.0",
			htmlContent:   `<html><body><p> 版本 : old </p></body></html>`,
			expectedMatch: true,
		},
		{
			name:          "no version tag",
			versionFile:   "APP_VERSION:1.0.0",
			htmlContent:   `<html><body><p>Other content</p></body></html>`,
			expectedMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			versionFile := filepath.Join(tmpDir, "version")
			err := os.WriteFile(versionFile, []byte(tt.versionFile), 0o600)
			if err != nil {
				t.Fatalf("failed to create version file: %v", err)
			}

			frontendDir := filepath.Join(tmpDir, "frontend", "src", "A")
			err = os.MkdirAll(frontendDir, 0o755)
			if err != nil {
				t.Fatalf("failed to create frontend directory: %v", err)
			}

			htmlFile := filepath.Join(frontendDir, "aboutVersionView.html")
			err = os.WriteFile(htmlFile, []byte(tt.htmlContent), 0o600)
			if err != nil {
				t.Fatalf("failed to create html file: %v", err)
			}

			originalGetRootPathFunc := getRootPathFunc
			getRootPathFunc = func() string {
				return tmpDir
			}
			defer func() { getRootPathFunc = originalGetRootPathFunc }()

			originalPrintErrorFunc := printErrorFunc
			printErrorFunc = func(err error) {
				t.Errorf("unexpected error: %v", err)
			}
			defer func() { printErrorFunc = originalPrintErrorFunc }()

			replaceAppVersion()

			result, err := os.ReadFile(htmlFile)
			if err != nil {
				t.Fatalf("failed to read html file: %v", err)
			}

			resultStr := string(result)
			if tt.expectedMatch {
				if !strings.Contains(resultStr, strings.TrimSpace(strings.ReplaceAll(tt.versionFile, "APP_VERSION:", ""))) {
					t.Errorf("replaceAppVersion() result does not contain expected version: %q", resultStr)
				}
			} else {
				if resultStr != tt.htmlContent {
					t.Errorf("replaceAppVersion() modified content when it shouldn't: %q", resultStr)
				}
			}
		})
	}
}

func TestChangeAppModeToDevInvalidMode(t *testing.T) {
	tmpDir := t.TempDir()

	mainFile := filepath.Join(tmpDir, "main.go")
	err := os.WriteFile(mainFile, []byte("test"), 0o600)
	if err != nil {
		t.Fatalf("failed to create main.go: %v", err)
	}

	originalGetFilePathFunc := getFilePathFunc
	getFilePathFunc = func(p string) string {
		return filepath.Join(tmpDir, p)
	}
	defer func() { getFilePathFunc = originalGetFilePathFunc }()

	errorCalled := false
	originalPrintErrorFunc := printErrorFunc
	printErrorFunc = func(err error) {
		errorCalled = true
	}
	defer func() { printErrorFunc = originalPrintErrorFunc }()

	changeAppModeToDev("invalid_mode")

	if !errorCalled {
		t.Error("printError was not called for invalid mode")
	}
}
