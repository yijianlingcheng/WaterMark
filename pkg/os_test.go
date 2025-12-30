package pkg

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

func TestIsWindows(t *testing.T) {
	result := IsWindows()
	expected := runtime.GOOS == Window

	if result != expected {
		t.Errorf("IsWindows() = %v, want %v (runtime.GOOS = %s)", result, expected, runtime.GOOS)
	}

	t.Logf("Running on %s, IsWindows() = %v", runtime.GOOS, result)
}

func TestGetDirFiles(t *testing.T) {
	t.Run("existing directory with files", func(t *testing.T) {
		tempDir := t.TempDir()

		testFiles := []string{"file1.txt", "file2.txt", "file3.jpg"}
		for _, filename := range testFiles {
			filePath := filepath.Join(tempDir, filename)
			if err := os.WriteFile(filePath, []byte("test content"), 0o644); err != nil {
				t.Fatalf("Failed to create test file %s: %v", filename, err)
			}
		}

		files, err := GetDirFiles(tempDir)
		if err != NoError {
			t.Errorf("GetDirFiles() error = %v, want NoError", err)
		}

		if len(files) != len(testFiles) {
			t.Errorf("GetDirFiles() returned %d files, want %d", len(files), len(testFiles))
		}

		for _, expectedFile := range testFiles {
			found := false
			for _, file := range files {
				if file == expectedFile {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("GetDirFiles() did not return expected file %s", expectedFile)
			}
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		tempDir := t.TempDir()

		files, err := GetDirFiles(tempDir)
		if err != NoError {
			t.Errorf("GetDirFiles() error = %v, want NoError", err)
		}

		if len(files) != 0 {
			t.Errorf("GetDirFiles() returned %d files, want 0", len(files))
		}
	})

	t.Run("directory with subdirectories", func(t *testing.T) {
		tempDir := t.TempDir()

		testFiles := []string{"file1.txt", "file2.txt"}
		for _, filename := range testFiles {
			filePath := filepath.Join(tempDir, filename)
			if err := os.WriteFile(filePath, []byte("test content"), 0o644); err != nil {
				t.Fatalf("Failed to create test file %s: %v", filename, err)
			}
		}

		subDir := filepath.Join(tempDir, "subdir")
		if err := os.Mkdir(subDir, 0o755); err != nil {
			t.Fatalf("Failed to create subdirectory: %v", err)
		}

		subFile := filepath.Join(subDir, "subfile.txt")
		if err := os.WriteFile(subFile, []byte("sub content"), 0o644); err != nil {
			t.Fatalf("Failed to create subfile: %v", err)
		}

		files, err := GetDirFiles(tempDir)
		if err != NoError {
			t.Errorf("GetDirFiles() error = %v, want NoError", err)
		}

		if len(files) != len(testFiles) {
			t.Errorf(
				"GetDirFiles() returned %d files, want %d (subdirectories should be excluded)",
				len(files),
				len(testFiles),
			)
		}

		for _, expectedFile := range testFiles {
			found := false
			for _, file := range files {
				if file == expectedFile {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("GetDirFiles() did not return expected file %s", expectedFile)
			}
		}
	})

	t.Run("nonexistent directory", func(t *testing.T) {
		nonexistentDir := filepath.Join(t.TempDir(), "nonexistent")

		files, err := GetDirFiles(nonexistentDir)
		if err == NoError {
			t.Error("GetDirFiles() expected error for nonexistent directory, got NoError")
		}

		if !HasError(err) {
			t.Errorf("GetDirFiles() error = %v, want error", err)
		}

		if len(files) != 0 {
			t.Errorf("GetDirFiles() returned %d files, want 0", len(files))
		}
	})

	t.Run("directory with hidden files", func(t *testing.T) {
		tempDir := t.TempDir()

		testFiles := []string{".hidden_file", "normal_file.txt"}
		for _, filename := range testFiles {
			filePath := filepath.Join(tempDir, filename)
			if err := os.WriteFile(filePath, []byte("test content"), 0o644); err != nil {
				t.Fatalf("Failed to create test file %s: %v", filename, err)
			}
		}

		files, err := GetDirFiles(tempDir)
		if err != NoError {
			t.Errorf("GetDirFiles() error = %v, want NoError", err)
		}

		if len(files) != len(testFiles) {
			t.Errorf("GetDirFiles() returned %d files, want %d", len(files), len(testFiles))
		}
	})

	t.Run("directory with many files", func(t *testing.T) {
		tempDir := t.TempDir()

		numFiles := 100
		for i := 0; i < numFiles; i++ {
			filePath := filepath.Join(tempDir, "file"+strconv.Itoa(i)+".txt")
			if err := os.WriteFile(filePath, []byte("test content"), 0o644); err != nil {
				t.Fatalf("Failed to create test file %d: %v", i, err)
			}
		}

		files, err := GetDirFiles(tempDir)
		if err != NoError {
			t.Errorf("GetDirFiles() error = %v, want NoError", err)
		}

		if len(files) != numFiles {
			t.Errorf("GetDirFiles() returned %d files, want %d", len(files), numFiles)
		}
	})
}

func TestOSConstants(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"Darwin", Darwin},
		{"Window", Window},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("Constant %s should not be empty", tt.name)
			}
		})
	}
}

func TestGetDirFilesWithSpecialCharacters(t *testing.T) {
	t.Run("directory with files containing special characters", func(t *testing.T) {
		tempDir := t.TempDir()

		testFiles := []string{"file with spaces.txt", "file-with-dashes.txt", "file_with_underscores.txt"}
		for _, filename := range testFiles {
			filePath := filepath.Join(tempDir, filename)
			if err := os.WriteFile(filePath, []byte("test content"), 0o644); err != nil {
				t.Fatalf("Failed to create test file %s: %v", filename, err)
			}
		}

		files, err := GetDirFiles(tempDir)
		if err != NoError {
			t.Errorf("GetDirFiles() error = %v, want NoError", err)
		}

		if len(files) != len(testFiles) {
			t.Errorf("GetDirFiles() returned %d files, want %d", len(files), len(testFiles))
		}
	})
}
