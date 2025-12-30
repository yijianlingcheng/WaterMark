package pkg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetStrMD5(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name:     "simple string",
			input:    "hello",
			expected: "5d41402abc4b2a76b9719d911017c592",
		},
		{
			name:     "string with spaces",
			input:    "hello world",
			expected: "5eb63bbbe01eeed093cb22bb8f5acdc3",
		},
		{
			name:     "string with special characters",
			input:    "hello@world!",
			expected: "854245bd59d2a3d09f4444e761e6cb5e",
		},
		{
			name:     "string with unicode",
			input:    "你好世界",
			expected: "65396ee4aad0b4f17aacd1c6112ee364",
		},
		{
			name:     "long string",
			input:    "This is a very long string that should still produce a valid MD5 hash",
			expected: "82a3a87a33054271935cfa3f0a463e72",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStrMD5(tt.input)
			if result != tt.expected {
				t.Errorf("GetStrMD5() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetFileMD5(t *testing.T) {
	t.Run("existing file", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		content := []byte("hello world")
		if err := os.WriteFile(testFile, content, 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		result, err := GetFileMD5(testFile)
		if err != NoError {
			t.Errorf("GetFileMD5() error = %v, want NoError", err)
		}
		expected := "5eb63bbbe01eeed093cb22bb8f5acdc3"
		if result != expected {
			t.Errorf("GetFileMD5() = %v, want %v", result, expected)
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		_, err := GetFileMD5("/nonexistent/file.txt")
		if err == NoError {
			t.Error("GetFileMD5() expected error for nonexistent file, got NoError")
		}
		if !HasError(err) {
			t.Errorf("GetFileMD5() error = %v, want error", err)
		}
	})

	t.Run("empty file", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "empty.txt")
		if err := os.WriteFile(testFile, []byte{}, 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		result, err := GetFileMD5(testFile)
		if err != NoError {
			t.Errorf("GetFileMD5() error = %v, want NoError", err)
		}
		expected := "d41d8cd98f00b204e9800998ecf8427e"
		if result != expected {
			t.Errorf("GetFileMD5() = %v, want %v", result, expected)
		}
	})

	t.Run("binary file", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "binary.bin")
		content := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD, 0xFC}
		if err := os.WriteFile(testFile, content, 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		result, err := GetFileMD5(testFile)
		if err != NoError {
			t.Errorf("GetFileMD5() error = %v, want NoError", err)
		}
		expected := "185781f9614c8100a813e0d8cbec0332"
		if result != expected {
			t.Errorf("GetFileMD5() = %v, want %v", result, expected)
		}
	})

	t.Run("large file", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "large.txt")
		content := make([]byte, 1024*1024)
		for i := range content {
			content[i] = byte(i % 256)
		}
		if err := os.WriteFile(testFile, content, 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		result, err := GetFileMD5(testFile)
		if err != NoError {
			t.Errorf("GetFileMD5() error = %v, want NoError", err)
		}
		if result == "" {
			t.Error("GetFileMD5() returned empty string for large file")
		}
	})
}
