package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bodgit/sevenzip"

	"WaterMark/pkg"
)

func TestUnzip7z(t *testing.T) {
	tests := []struct {
		name        string
		zipPath     string
		unzipPath   string
		expectError bool
	}{
		{
			name:        "valid 7z with single file",
			zipPath:     filepath.Join("testdata", "single.7z"),
			unzipPath:   filepath.Join(os.TempDir(), "test_unzip7z_single"),
			expectError: false,
		},
		{
			name:        "valid 7z with multiple files",
			zipPath:     filepath.Join("testdata", "multiple.7z"),
			unzipPath:   filepath.Join(os.TempDir(), "test_unzip7z_multiple"),
			expectError: false,
		},
		{
			name:        "valid 7z with directory",
			zipPath:     filepath.Join("testdata", "withdir.7z"),
			unzipPath:   filepath.Join(os.TempDir(), "test_unzip7z_dir"),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupTest(t, tt.unzipPath)
			defer cleanup()

			if _, err := os.Stat(tt.zipPath); os.IsNotExist(err) {
				t.Skip("Test file not found: " + tt.zipPath)
			}

			err := Unzip7z(tt.zipPath, tt.unzipPath)

			if tt.expectError {
				if !pkg.HasError(err) {
					t.Error("Expected error, but got nil")
				}
			} else {
				if pkg.HasError(err) {
					t.Errorf("Unexpected error: %v", err)
				}
				verifyUnzip(t, tt.unzipPath)
			}
		})
	}
}

func TestUnzip7zInvalidPath(t *testing.T) {
	unzipPath := filepath.Join(os.TempDir(), "test_unzip7z_invalid")
	cleanup := setupTest(t, unzipPath)
	defer cleanup()

	err := Unzip7z("nonexistent.7z", unzipPath)
	if !pkg.HasError(err) {
		t.Error("Expected error for invalid 7z path, but got nil")
	}
}

func TestUnzip7zWithNestedDirectories(t *testing.T) {
	unzipPath := filepath.Join(os.TempDir(), "test_unzip7z_nested")
	zipPath := filepath.Join("testdata", "nested.7z")
	cleanup := setupTest(t, unzipPath)
	defer cleanup()

	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		t.Skip("Test file not found: " + zipPath)
	}

	err := Unzip7z(zipPath, unzipPath)
	if pkg.HasError(err) {
		t.Fatalf("Failed to unzip: %v", err)
	}

	filenames := []string{
		"root.txt",
		"dir1/file1.txt",
		"dir1/subdir/file2.txt",
		"dir2/file3.txt",
	}

	for _, filename := range filenames {
		expectedPath := filepath.Join(unzipPath, filename)
		if _, statErr := os.Stat(expectedPath); os.IsNotExist(statErr) {
			t.Errorf("File %s was not created", expectedPath)
		}
	}
}

func TestUnzip7zWithLargeFile(t *testing.T) {
	unzipPath := filepath.Join(os.TempDir(), "test_unzip7z_large")
	zipPath := filepath.Join("testdata", "large.7z")
	cleanup := setupTest(t, unzipPath)
	defer cleanup()

	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		t.Skip("Test file not found: " + zipPath)
	}

	err := Unzip7z(zipPath, unzipPath)
	if pkg.HasError(err) {
		t.Fatalf("Failed to unzip: %v", err)
	}

	expectedPath := filepath.Join(unzipPath, "large.txt")
	_, statErr := os.Stat(expectedPath)
	if statErr != nil {
		t.Errorf("Failed to stat large file: %v", statErr)
	}

	info, _ := os.Stat(expectedPath)
	if info.Size() == 0 {
		t.Error("Large file is empty")
	}
}

func TestUnzip7zWithSpecialCharacters(t *testing.T) {
	unzipPath := filepath.Join(os.TempDir(), "test_unzip7z_special")
	zipPath := filepath.Join("testdata", "special.7z")
	cleanup := setupTest(t, unzipPath)
	defer cleanup()

	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		t.Skip("Test file not found: " + zipPath)
	}

	err := Unzip7z(zipPath, unzipPath)
	if pkg.HasError(err) {
		t.Fatalf("Failed to unzip: %v", err)
	}

	filenames := []string{
		"file with spaces.txt",
		"文件.txt",
		"file-with-dashes.txt",
	}

	for _, filename := range filenames {
		expectedPath := filepath.Join(unzipPath, filename)
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("File %s was not created", expectedPath)
		}
	}
}

func TestUnzip7zWithPermissions(t *testing.T) {
	unzipPath := filepath.Join(os.TempDir(), "test_unzip7z_permissions")
	zipPath := filepath.Join("testdata", "permissions.7z")
	cleanup := setupTest(t, unzipPath)
	defer cleanup()

	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		t.Skip("Test file not found: " + zipPath)
	}

	err := Unzip7z(zipPath, unzipPath)
	if pkg.HasError(err) {
		t.Fatalf("Failed to unzip: %v", err)
	}

	expectedPath := filepath.Join(unzipPath, "executable.sh")
	info, statErr := os.Stat(expectedPath)
	if statErr != nil {
		t.Errorf("Failed to stat file: %v", statErr)
	}

	if info.IsDir() {
		t.Error("Expected file, got directory")
	}

	if info.Size() == 0 {
		t.Error("File is empty")
	}
}

func TestUnzip7zOverwriteExisting(t *testing.T) {
	unzipPath := filepath.Join(os.TempDir(), "test_unzip7z_overwrite")
	zipPath := filepath.Join("testdata", "overwrite.7z")
	cleanup := setupTest(t, unzipPath)
	defer cleanup()

	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		t.Skip("Test file not found: " + zipPath)
	}

	if err := os.MkdirAll(unzipPath, os.ModePerm); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	existingFile := filepath.Join(unzipPath, "test.txt")
	if err := os.WriteFile(existingFile, []byte("old content"), 0o644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	err := Unzip7z(zipPath, unzipPath)
	if pkg.HasError(err) {
		t.Fatalf("Failed to unzip: %v", err)
	}

	content, readErr := os.ReadFile(existingFile)
	if readErr != nil {
		t.Errorf("Failed to read file: %v", readErr)
	}

	if string(content) == "old content" {
		t.Error("File content was not overwritten")
	}
}

func TestUnzip7zEmpty(t *testing.T) {
	unzipPath := filepath.Join(os.TempDir(), "test_unzip7z_empty")
	zipPath := filepath.Join("testdata", "empty.7z")
	cleanup := setupTest(t, unzipPath)
	defer cleanup()

	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		t.Skip("Test file not found: " + zipPath)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unzip with empty 7z should not panic, got: %v", r)
		}
	}()

	err := Unzip7z(zipPath, unzipPath)
	if pkg.HasError(err) {
		t.Fatalf("Failed to unzip: %v", err)
	}

	entries, readErr := os.ReadDir(unzipPath)
	if readErr != nil && !os.IsNotExist(readErr) {
		t.Errorf("Failed to read unzip directory: %v", readErr)
	}

	if len(entries) > 0 {
		t.Errorf("Empty 7z should not create any files, got %d entries", len(entries))
	}
}

func TestExtract7zFile(t *testing.T) {
	unzipPath := filepath.Join(os.TempDir(), "test_extract7z_file")
	zipPath := filepath.Join("testdata", "single.7z")
	cleanup := setupTest(t, unzipPath)
	defer cleanup()

	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		t.Skip("Test file not found: " + zipPath)
	}

	reader, err := sevenzip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("Failed to open 7z reader: %v", err)
	}
	defer reader.Close()

	if len(reader.File) == 0 {
		t.Fatal("No files in test 7z archive")
	}

	for _, f := range reader.File {
		if !f.FileInfo().IsDir() {
			err := extract7zFile(unzipPath, f)
			if err != nil {
				t.Errorf("Failed to extract file: %v", err)
			}

			expectedPath := filepath.Join(unzipPath, f.Name)
			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				t.Errorf("File %s was not created", expectedPath)
			}
		}
	}
}
