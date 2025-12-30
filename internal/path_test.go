package internal

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"WaterMark/pkg"
)

func TestInitRootPath(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "initialize root path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initRootPath()
			if GetRootPath() == "" {
				t.Error("GetRootPath should not return empty string after initRootPath")
			}
		})
	}
}

func TestGetRootPath(t *testing.T) {
	initRootPath()
	rootPath := GetRootPath()

	if rootPath == "" {
		t.Error("GetRootPath should not return empty string")
	}

	fileInfo, err := os.Stat(rootPath)
	if err != nil {
		t.Errorf("GetRootPath should return a valid directory: %v", err)
	}

	if !fileInfo.IsDir() {
		t.Error("GetRootPath should return a directory path")
	}
}

func TestGetPwdPath(t *testing.T) {
	initRootPath()
	testPath := "/test/path"
	result := GetPwdPath(testPath)

	if result == "" {
		t.Error("GetPwdPath should not return empty string")
	}

	if result[len(result)-len(testPath):] != testPath {
		t.Errorf("GetPwdPath should end with the provided path. Expected suffix: %s, Got: %s", testPath, result)
	}
}

func TestGetConfigPath(t *testing.T) {
	initRootPath()
	testPath := "test.yaml"
	result := GetConfigPath(testPath)

	if result == "" {
		t.Error("GetConfigPath should not return empty string")
	}

	if filepath.Base(result) != testPath {
		t.Errorf(
			"GetConfigPath should end with the provided filename. Expected: %s, Got: %s",
			testPath,
			filepath.Base(result),
		)
	}
}

func TestGetLogPath(t *testing.T) {
	initRootPath()
	testPath := "test.log"
	result := GetLogPath(testPath)

	if result == "" {
		t.Error("GetLogPath should not return empty string")
	}

	if filepath.Base(result) != testPath {
		t.Errorf(
			"GetLogPath should end with the provided filename. Expected: %s, Got: %s",
			testPath,
			filepath.Base(result),
		)
	}
}

func TestGetLogosPath(t *testing.T) {
	initRootPath()
	testPath := "test.png"
	result := GetLogosPath(testPath)

	if result == "" {
		t.Error("GetLogosPath should not return empty string")
	}

	if filepath.Base(result) != testPath {
		t.Errorf(
			"GetLogosPath should end with the provided filename. Expected: %s, Got: %s",
			testPath,
			filepath.Base(result),
		)
	}
}

func TestGetRuntimePath(t *testing.T) {
	initRootPath()
	testPath := "test.txt"
	result := GetRuntimePath(testPath)

	if result == "" {
		t.Error("GetRuntimePath should not return empty string")
	}

	if filepath.Base(result) != testPath {
		t.Errorf(
			"GetRuntimePath should end with the provided filename. Expected: %s, Got: %s",
			testPath,
			filepath.Base(result),
		)
	}
}

func TestGetUserDirectory(t *testing.T) {
	initRootPath()
	testPath := "user/test.txt"
	result := GetUserDirectory(testPath)

	if result == "" {
		t.Error("GetUserDirectory should not return empty string")
	}

	if filepath.Base(result) != "test.txt" {
		t.Errorf(
			"GetUserDirectory should end with the provided filename. Expected: test.txt, Got: %s",
			filepath.Base(result),
		)
	}
}

func TestGetFontFilePath(t *testing.T) {
	initRootPath()
	testPath := "test.ttf"
	result := GetFontFilePath(testPath)

	if result == "" {
		t.Error("GetFontFilePath should not return empty string")
	}

	if filepath.Base(result) != testPath {
		t.Errorf(
			"GetFontFilePath should end with the provided filename. Expected: %s, Got: %s",
			testPath,
			filepath.Base(result),
		)
	}
}

func TestGetAppExifCacheFilePath(t *testing.T) {
	initRootPath()
	result := GetAppExifCacheFilePath()

	if result == "" {
		t.Error("GetAppExifCacheFilePath should not return empty string")
	}

	if filepath.Base(result) != "exifCache.cache" {
		t.Errorf(
			"GetAppExifCacheFilePath should return path ending with 'exifCache.cache'. Got: %s",
			filepath.Base(result),
		)
	}
}

func TestGetExiftoolZipPath(t *testing.T) {
	initRootPath()
	result := GetExiftoolZipPath()

	if result == "" {
		t.Error("GetExiftoolZipPath should not return empty string")
	}

	if filepath.Base(result) != "exiftool.zip" {
		t.Errorf("GetExiftoolZipPath should return path ending with 'exiftool.zip'. Got: %s", filepath.Base(result))
	}
}

func TestGetExiftoolUnzipPath(t *testing.T) {
	initRootPath()
	result := GetExiftoolUnzipPath()

	if result == "" {
		t.Error("GetExiftoolUnzipPath should not return empty string")
	}

	if result[len(result)-1] != filepath.Separator && result[len(result)-1] != '/' {
		t.Errorf("GetExiftoolUnzipPath should end with path separator. Got: %s", result)
	}
}

func TestGetMainLayoutPath(t *testing.T) {
	initRootPath()
	result := GetMainLayoutPath()

	if result == "" {
		t.Error("GetMainLayoutPath should not return empty string")
	}

	if filepath.Base(result) != "layout.json" {
		t.Errorf("GetMainLayoutPath should return path ending with 'layout.json'. Got: %s", filepath.Base(result))
	}
}

func TestGetMagickPath(t *testing.T) {
	initRootPath()
	testPath := "test.exe"
	result := GetMagickPath(testPath)

	if result == "" {
		t.Error("GetMagickPath should not return empty string")
	}

	if filepath.Base(result) != testPath {
		t.Errorf(
			"GetMagickPath should end with the provided filename. Expected: %s, Got: %s",
			testPath,
			filepath.Base(result),
		)
	}
}

func TestGetWinMagick7zPath(t *testing.T) {
	initRootPath()
	result := GetWinMagick7zPath()

	if result == "" {
		t.Error("GetWinMagick7zPath should not return empty string")
	}

	if filepath.Base(result) != "ImageMagick.7z" {
		t.Errorf("GetWinMagick7zPath should return path ending with 'ImageMagick.7z'. Got: %s", filepath.Base(result))
	}
}

func TestGetMagickBinPath(t *testing.T) {
	initRootPath()
	result := GetMagickBinPath()

	if result == "" {
		t.Error("GetMagickBinPath should not return empty string")
	}

	if IsWindows() {
		if filepath.Base(result) != "magick.exe" {
			t.Errorf(
				"On Windows, GetMagickBinPath should return path ending with 'magick.exe'. Got: %s",
				filepath.Base(result),
			)
		}
	} else {
		if result != "magick" {
			t.Errorf("On non-Windows, GetMagickBinPath should return 'magick'. Got: %s", result)
		}
	}
}

func TestGetExiftoolPath(t *testing.T) {
	initRootPath()
	result := GetExiftoolPath()

	if result == "" {
		t.Error("GetExiftoolPath should not return empty string")
	}

	if IsWindows() {
		if filepath.Base(result) != "exiftool.exe" {
			t.Errorf(
				"On Windows, GetExiftoolPath should return path ending with 'exiftool.exe'. Got: %s",
				filepath.Base(result),
			)
		}
	} else {
		if filepath.Base(result) != "exiftool" {
			t.Errorf("On non-Windows, GetExiftoolPath should return path ending with 'exiftool'. Got: %s", filepath.Base(result))
		}
	}
}

func TestPathExists(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() (string, func())
		expected bool
	}{
		{
			name: "existing file",
			setup: func() (string, func()) {
				tmpDir := t.TempDir()
				testFile := filepath.Join(tmpDir, "test.txt")
				os.WriteFile(testFile, []byte("test"), 0o644)
				return testFile, func() {}
			},
			expected: true,
		},
		{
			name: "existing directory",
			setup: func() (string, func()) {
				tmpDir := t.TempDir()
				return tmpDir, func() {}
			},
			expected: true,
		},
		{
			name: "nonexistent path",
			setup: func() (string, func()) {
				tmpDir := t.TempDir()
				nonexistentPath := filepath.Join(tmpDir, "nonexistent.txt")
				return nonexistentPath, func() {}
			},
			expected: false,
		},
		{
			name: "empty path",
			setup: func() (string, func()) {
				return "", func() {}
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, cleanup := tt.setup()
			defer cleanup()

			result := PathExists(path)
			if result != tt.expected {
				t.Errorf("PathExists(%q) = %v, want %v", path, result, tt.expected)
			}
		})
	}
}

func TestGetAppBlurFilePath(t *testing.T) {
	initRootPath()
	testPath := "test.jpg"
	result := GetAppBlurFilePath(testPath)

	if result == "" {
		t.Error("GetAppBlurFilePath should not return empty string")
	}

	if filepath.Base(result) != testPath {
		t.Errorf(
			"GetAppBlurFilePath should end with the provided filename. Expected: %s, Got: %s",
			testPath,
			filepath.Base(result),
		)
	}
}

func TestCleanDir(t *testing.T) {
	initRootPath()
	blurPath := GetRootPath() + appBlurPath

	if err := os.MkdirAll(blurPath, os.ModePerm); err != nil {
		t.Fatalf("Failed to create blur directory: %v", err)
	}

	testFile := filepath.Join(blurPath, "test.jpg")
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	CleanDir()

	if PathExists(blurPath) {
		t.Error("CleanDir should remove the blur directory")
	}
}

func TestCreateAppDS(t *testing.T) {
	tests := []struct {
		name    string
		dirs    []string
		wantErr bool
	}{
		{
			name:    "create single directory",
			dirs:    []string{filepath.Join(t.TempDir(), "test1")},
			wantErr: false,
		},
		{
			name:    "create multiple directories",
			dirs:    []string{filepath.Join(t.TempDir(), "test1"), filepath.Join(t.TempDir(), "test2")},
			wantErr: false,
		},
		{
			name:    "create nested directories",
			dirs:    []string{filepath.Join(t.TempDir(), "test1", "sub1", "sub2")},
			wantErr: false,
		},
		{
			name:    "existing directory",
			dirs:    []string{t.TempDir()},
			wantErr: false,
		},
		{
			name:    "empty list",
			dirs:    []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := createAppDS(tt.dirs)
			if pkg.HasError(err) != tt.wantErr {
				t.Errorf("createAppDS() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, dir := range tt.dirs {
				if dir != "" && !PathExists(dir) {
					t.Errorf("createAppDS() should create directory %s", dir)
				}
			}
		})
	}
}

func TestIsWindows(t *testing.T) {
	expected := runtime.GOOS == "windows"
	result := IsWindows()

	if result != expected {
		t.Errorf("IsWindows() = %v, want %v", result, expected)
	}
}

func TestPathConsistency(t *testing.T) {
	initRootPath()
	rootPath := GetRootPath()

	if rootPath == "" {
		t.Fatal("GetRootPath returned empty string")
	}

	tests := []struct {
		name string
		fn   func(string) string
		arg  string
	}{
		{"GetPwdPath", GetPwdPath, "/test"},
		{"GetConfigPath", GetConfigPath, "config.yaml"},
		{"GetLogPath", GetLogPath, "app.log"},
		{"GetLogosPath", GetLogosPath, "logo.png"},
		{"GetRuntimePath", GetRuntimePath, "temp.txt"},
		{"GetUserDirectory", GetUserDirectory, "user/file.txt"},
		{"GetFontFilePath", GetFontFilePath, "font.ttf"},
		{"GetAppBlurFilePath", GetAppBlurFilePath, "blur.jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.arg)
			if result == "" {
				t.Errorf("%s returned empty string", tt.name)
			}

			if !filepath.IsAbs(result) {
				t.Errorf("%s should return absolute path, got: %s", tt.name, result)
			}
		})
	}
}
