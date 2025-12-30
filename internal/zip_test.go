package internal

import (
	"WaterMark/pkg"
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestUnzip(t *testing.T) {
	tests := []struct {
		name        string
		zipContent  *bytes.Buffer
		unzipPath   string
		expectError bool
	}{
		{
			name:        "valid zip with single file",
			zipContent:  createTestZip([]string{"test.txt"}, []string{"test content"}),
			unzipPath:   filepath.Join(os.TempDir(), "test_unzip_single"),
			expectError: false,
		},
		{
			name:        "valid zip with multiple files",
			zipContent:  createTestZip([]string{"file1.txt", "file2.txt"}, []string{"content1", "content2"}),
			unzipPath:   filepath.Join(os.TempDir(), "test_unzip_multiple"),
			expectError: false,
		},
		{
			name:        "valid zip with directory",
			zipContent:  createTestZipWithDir([]string{"dir/file.txt"}, []string{"dir content"}),
			unzipPath:   filepath.Join(os.TempDir(), "test_unzip_dir"),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupTest(t, tt.unzipPath)
			defer cleanup()

			zipPath := filepath.Join(os.TempDir(), "test.zip")
			if err := os.WriteFile(zipPath, tt.zipContent.Bytes(), 0o644); err != nil {
				t.Fatalf("Failed to write test zip file: %v", err)
			}
			defer os.Remove(zipPath)

			Unzip(zipPath, tt.unzipPath)

			if tt.expectError {
				t.Error("Expected error but got none")
			} else {
				verifyUnzip(t, tt.unzipPath)
			}
		})
	}
}

func TestUnzipInvalidPath(t *testing.T) {
	unzipPath := filepath.Join(os.TempDir(), "test_unzip_invalid")
	cleanup := setupTest(t, unzipPath)
	defer cleanup()

	err := Unzip("nonexistent.zip", unzipPath)
	if !pkg.HasError(err) {
		t.Error("Expected error for invalid zip path, but got nil")
	}
}

func setupTest(t *testing.T, unzipPath string) func() {
	if err := os.RemoveAll(unzipPath); err != nil {
		t.Fatalf("Failed to cleanup test directory: %v", err)
	}

	return func() {
		if err := os.RemoveAll(unzipPath); err != nil {
			t.Logf("Failed to cleanup test directory: %v", err)
		}
	}
}

func createTestZip(filenames []string, contents []string) *bytes.Buffer {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for i, filename := range filenames {
		writer, err := zipWriter.Create(filename)
		if err != nil {
			continue
		}
		_, _ = writer.Write([]byte(contents[i]))
	}

	zipWriter.Close()
	return buf
}

func createTestZipWithDir(filenames []string, contents []string) *bytes.Buffer {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for i, filename := range filenames {
		writer, err := zipWriter.Create(filename)
		if err != nil {
			continue
		}
		_, _ = writer.Write([]byte(contents[i]))
	}

	zipWriter.Close()
	return buf
}

func verifyUnzip(t *testing.T, unzipPath string) {
	err := filepath.Walk(unzipPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("Failed to read file %s: %v", path, err)
			}

			if len(content) == 0 {
				t.Errorf("File %s is empty", path)
			}
		}

		return nil
	})
	if err != nil {
		t.Errorf("Failed to walk unzip directory: %v", err)
	}
}

func TestUnzipWithEmptyZip(t *testing.T) {
	unzipPath := filepath.Join(os.TempDir(), "test_unzip_empty")
	cleanup := setupTest(t, unzipPath)
	defer cleanup()

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	zipWriter.Close()

	zipPath := filepath.Join(os.TempDir(), "empty.zip")
	if err := os.WriteFile(zipPath, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("Failed to write empty zip file: %v", err)
	}
	defer os.Remove(zipPath)

	Unzip(zipPath, unzipPath)

	entries, err := os.ReadDir(unzipPath)
	if err != nil && !os.IsNotExist(err) {
		t.Errorf("Failed to read unzip directory: %v", err)
	}

	if len(entries) > 0 {
		t.Errorf("Empty zip should not create any files, got %d entries", len(entries))
	}
}

func TestUnzipWithNestedDirectories(t *testing.T) {
	unzipPath := filepath.Join(os.TempDir(), "test_unzip_nested")
	cleanup := setupTest(t, unzipPath)
	defer cleanup()

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	filenames := []string{
		"root.txt",
		"dir1/file1.txt",
		"dir1/subdir/file2.txt",
		"dir2/file3.txt",
	}
	contents := []string{
		"root content",
		"file1 content",
		"file2 content",
		"file3 content",
	}

	for i, filename := range filenames {
		writer, err := zipWriter.Create(filename)
		if err != nil {
			continue
		}
		_, _ = writer.Write([]byte(contents[i]))
	}

	zipWriter.Close()

	zipPath := filepath.Join(os.TempDir(), "nested.zip")
	if err := os.WriteFile(zipPath, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("Failed to write nested zip file: %v", err)
	}
	defer os.Remove(zipPath)

	Unzip(zipPath, unzipPath)

	for i, filename := range filenames {
		expectedPath := filepath.Join(unzipPath, filename)
		content, err := os.ReadFile(expectedPath)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", expectedPath, err)
		}

		if string(content) != contents[i] {
			t.Errorf("File %s content mismatch. Expected: %s, Got: %s", expectedPath, contents[i], string(content))
		}
	}
}

func TestUnzipWithLargeFile(t *testing.T) {
	unzipPath := filepath.Join(os.TempDir(), "test_unzip_large")
	cleanup := setupTest(t, unzipPath)
	defer cleanup()

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	writer, err := zipWriter.Create("large.txt")
	if err != nil {
		t.Fatalf("Failed to create file in zip: %v", err)
	}

	largeContent := make([]byte, 1024*1024)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	if _, err := writer.Write(largeContent); err != nil {
		t.Fatalf("Failed to write large content: %v", err)
	}

	zipWriter.Close()

	zipPath := filepath.Join(os.TempDir(), "large.zip")
	if err := os.WriteFile(zipPath, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("Failed to write large zip file: %v", err)
	}
	defer os.Remove(zipPath)

	Unzip(zipPath, unzipPath)

	expectedPath := filepath.Join(unzipPath, "large.txt")
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Errorf("Failed to read large file: %v", err)
	}

	if len(content) != len(largeContent) {
		t.Errorf("Large file size mismatch. Expected: %d, Got: %d", len(largeContent), len(content))
	}
}

func TestUnzipWithSpecialCharacters(t *testing.T) {
	unzipPath := filepath.Join(os.TempDir(), "test_unzip_special")
	cleanup := setupTest(t, unzipPath)
	defer cleanup()

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	filenames := []string{
		"file with spaces.txt",
		"文件.txt",
		"file-with-dashes.txt",
	}
	contents := []string{
		"spaces content",
		"中文内容",
		"dashes content",
	}

	for i, filename := range filenames {
		writer, err := zipWriter.Create(filename)
		if err != nil {
			continue
		}
		_, _ = writer.Write([]byte(contents[i]))
	}

	zipWriter.Close()

	zipPath := filepath.Join(os.TempDir(), "special.zip")
	if err := os.WriteFile(zipPath, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("Failed to write special zip file: %v", err)
	}
	defer os.Remove(zipPath)

	Unzip(zipPath, unzipPath)

	for i, filename := range filenames {
		expectedPath := filepath.Join(unzipPath, filename)
		content, err := os.ReadFile(expectedPath)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", expectedPath, err)
		}

		if string(content) != contents[i] {
			t.Errorf("File %s content mismatch. Expected: %s, Got: %s", expectedPath, contents[i], string(content))
		}
	}
}

func TestUnzipWithPermissions(t *testing.T) {
	unzipPath := filepath.Join(os.TempDir(), "test_unzip_permissions")
	cleanup := setupTest(t, unzipPath)
	defer cleanup()

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	writer, err := zipWriter.Create("executable.sh")
	if err != nil {
		t.Fatalf("Failed to create file in zip: %v", err)
	}

	_, _ = writer.Write([]byte("#!/bin/bash\necho 'hello'"))

	zipWriter.Close()

	zipPath := filepath.Join(os.TempDir(), "permissions.zip")
	if err := os.WriteFile(zipPath, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("Failed to write permissions zip file: %v", err)
	}
	defer os.Remove(zipPath)

	Unzip(zipPath, unzipPath)

	expectedPath := filepath.Join(unzipPath, "executable.sh")
	info, err := os.Stat(expectedPath)
	if err != nil {
		t.Errorf("Failed to stat file: %v", err)
	}

	if info.IsDir() {
		t.Error("Expected file, got directory")
	}

	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
	}

	if len(content) == 0 {
		t.Error("File is empty")
	}
}

func TestUnzipOverwriteExisting(t *testing.T) {
	unzipPath := filepath.Join(os.TempDir(), "test_unzip_overwrite")
	cleanup := setupTest(t, unzipPath)
	defer cleanup()

	if err := os.MkdirAll(unzipPath, os.ModePerm); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	existingFile := filepath.Join(unzipPath, "test.txt")
	if err := os.WriteFile(existingFile, []byte("old content"), 0o644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	writer, err := zipWriter.Create("test.txt")
	if err != nil {
		t.Fatalf("Failed to create file in zip: %v", err)
	}

	_, _ = writer.Write([]byte("new content"))

	zipWriter.Close()

	zipPath := filepath.Join(os.TempDir(), "overwrite.zip")
	if err := os.WriteFile(zipPath, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("Failed to write overwrite zip file: %v", err)
	}
	defer os.Remove(zipPath)

	Unzip(zipPath, unzipPath)

	content, err := os.ReadFile(existingFile)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
	}

	if string(content) != "new content" {
		t.Errorf("File content not overwritten. Expected: 'new content', Got: '%s'", string(content))
	}
}

func TestCreateTestZip(t *testing.T) {
	buf := createTestZip([]string{"test.txt"}, []string{"test content"})

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("Failed to create zip reader: %v", err)
	}

	if len(zipReader.File) != 1 {
		t.Errorf("Expected 1 file, got %d", len(zipReader.File))
	}

	if zipReader.File[0].Name != "test.txt" {
		t.Errorf("Expected file name 'test.txt', got '%s'", zipReader.File[0].Name)
	}

	rc, err := zipReader.File[0].Open()
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("Failed to read file content: %v", err)
	}

	if string(content) != "test content" {
		t.Errorf("Expected content 'test content', got '%s'", string(content))
	}
}
