package pkg

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadImageWithDecode(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() (string, func())
		expectError bool
		errorCode   int
	}{
		{
			name: "load valid jpeg image",
			setupFunc: func() (string, func()) {
				buf := new(bytes.Buffer)
				img := image.NewRGBA(image.Rect(0, 0, 100, 100))
				jpeg.Encode(buf, img, nil)

				tmpFile := filepath.Join(os.TempDir(), "test.jpg")
				if err := os.WriteFile(tmpFile, buf.Bytes(), 0o644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				cleanup := func() { os.Remove(tmpFile) }
				return tmpFile, cleanup
			},
			expectError: false,
		},
		{
			name: "load valid png image",
			setupFunc: func() (string, func()) {
				buf := new(bytes.Buffer)
				img := image.NewRGBA(image.Rect(0, 0, 100, 100))
				png.Encode(buf, img)

				tmpFile := filepath.Join(os.TempDir(), "test.png")
				if err := os.WriteFile(tmpFile, buf.Bytes(), 0o644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				cleanup := func() { os.Remove(tmpFile) }
				return tmpFile, cleanup
			},
			expectError: false,
		},
		{
			name: "load nonexistent file",
			setupFunc: func() (string, func()) {
				tmpFile := filepath.Join(os.TempDir(), "nonexistent.jpg")
				cleanup := func() {}
				return tmpFile, cleanup
			},
			expectError: true,
			errorCode:   FILE_NOT_OPEN_ERROR,
		},
		{
			name: "load unsupported image format",
			setupFunc: func() (string, func()) {
				tmpFile := filepath.Join(os.TempDir(), "test.txt")
				if err := os.WriteFile(tmpFile, []byte("This is not an image file"), 0o644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				cleanup := func() { os.Remove(tmpFile) }
				return tmpFile, cleanup
			},
			expectError: true,
			errorCode:   IMAGE_NO_SUPPORT_ERROR,
		},
		{
			name: "load corrupted jpeg file",
			setupFunc: func() (string, func()) {
				corruptedData := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01}

				tmpFile := filepath.Join(os.TempDir(), "corrupted.jpg")
				if err := os.WriteFile(tmpFile, corruptedData, 0o644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				cleanup := func() { os.Remove(tmpFile) }
				return tmpFile, cleanup
			},
			expectError: true,
			errorCode:   IMAGE_DECODE_ERROR,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, cleanup := tt.setupFunc()
			defer cleanup()

			img, err := LoadImageWithDecode(path)

			if tt.expectError {
				if !HasError(err) {
					t.Errorf("Expected error but got none")
				}
				if tt.errorCode != 0 && err.Code != tt.errorCode {
					t.Errorf("Expected error code %d, got %d", tt.errorCode, err.Code)
				}
			} else {
				if HasError(err) {
					t.Errorf("Unexpected error: %v", err)
				}
				if img == nil {
					t.Errorf("Expected image but got nil")
				}
				bounds := img.Bounds()
				if bounds.Dx() != 100 || bounds.Dy() != 100 {
					t.Errorf("Expected image size 100x100, got %dx%d", bounds.Dx(), bounds.Dy())
				}
			}
		})
	}
}

func TestGetFileType(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() (*os.File, func())
		expectType  string
		expectError bool
		errorCode   int
	}{
		{
			name: "detect jpeg file type",
			setupFunc: func() (*os.File, func()) {
				buf := new(bytes.Buffer)
				img := image.NewRGBA(image.Rect(0, 0, 100, 100))
				jpeg.Encode(buf, img, nil)

				tmpFile := filepath.Join(os.TempDir(), "test.jpg")
				if err := os.WriteFile(tmpFile, buf.Bytes(), 0o644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				file, err := os.Open(tmpFile)
				if err != nil {
					t.Fatalf("Failed to open test file: %v", err)
				}
				cleanup := func() {
					file.Close()
					os.Remove(tmpFile)
				}
				return file, cleanup
			},
			expectType:  "image/jpeg",
			expectError: false,
		},
		{
			name: "detect png file type",
			setupFunc: func() (*os.File, func()) {
				buf := new(bytes.Buffer)
				img := image.NewRGBA(image.Rect(0, 0, 100, 100))
				png.Encode(buf, img)

				tmpFile := filepath.Join(os.TempDir(), "test.png")
				if err := os.WriteFile(tmpFile, buf.Bytes(), 0o644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				file, err := os.Open(tmpFile)
				if err != nil {
					t.Fatalf("Failed to open test file: %v", err)
				}
				cleanup := func() {
					file.Close()
					os.Remove(tmpFile)
				}
				return file, cleanup
			},
			expectType:  "image/png",
			expectError: false,
		},
		{
			name: "detect empty file",
			setupFunc: func() (*os.File, func()) {
				tmpFile := filepath.Join(os.TempDir(), "empty.txt")
				if err := os.WriteFile(tmpFile, []byte{}, 0o644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				file, err := os.Open(tmpFile)
				if err != nil {
					t.Fatalf("Failed to open test file: %v", err)
				}
				cleanup := func() {
					file.Close()
					os.Remove(tmpFile)
				}
				return file, cleanup
			},
			expectError: true,
			errorCode:   FILE_NOT_READ_ERROR,
		},
		{
			name: "detect text file",
			setupFunc: func() (*os.File, func()) {
				tmpFile := filepath.Join(os.TempDir(), "test.txt")
				if err := os.WriteFile(tmpFile, []byte("Hello, World!"), 0o644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				file, err := os.Open(tmpFile)
				if err != nil {
					t.Fatalf("Failed to open test file: %v", err)
				}
				cleanup := func() {
					file.Close()
					os.Remove(tmpFile)
				}
				return file, cleanup
			},
			expectType:  "application/octet-stream",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, cleanup := tt.setupFunc()
			defer cleanup()

			fileType, err := GetFileType(file)

			if tt.expectError {
				if !HasError(err) {
					t.Errorf("Expected error but got none")
				}
				if tt.errorCode != 0 && err.Code != tt.errorCode {
					t.Errorf("Expected error code %d, got %d", tt.errorCode, err.Code)
				}
			} else {
				if HasError(err) {
					t.Errorf("Unexpected error: %v", err)
				}
				if fileType != tt.expectType {
					t.Errorf("Expected file type %s, got %s", tt.expectType, fileType)
				}
			}
		})
	}
}

func TestGenerateImageByWidthHeight(t *testing.T) {
	tests := []struct {
		name         string
		originalSize image.Rectangle
		targetWidth  int
		targetHeight int
	}{
		{
			name:         "resize to smaller dimensions",
			originalSize: image.Rect(0, 0, 200, 200),
			targetWidth:  100,
			targetHeight: 100,
		},
		{
			name:         "resize to larger dimensions",
			originalSize: image.Rect(0, 0, 50, 50),
			targetWidth:  100,
			targetHeight: 100,
		},
		{
			name:         "resize to same dimensions",
			originalSize: image.Rect(0, 0, 100, 100),
			targetWidth:  100,
			targetHeight: 100,
		},
		{
			name:         "resize to different aspect ratio",
			originalSize: image.Rect(0, 0, 200, 100),
			targetWidth:  100,
			targetHeight: 100,
		},
		{
			name:         "resize small image to large",
			originalSize: image.Rect(0, 0, 10, 10),
			targetWidth:  500,
			targetHeight: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalImg := image.NewRGBA(tt.originalSize)

			resizedImg := GenerateImageByWidthHeight(originalImg, tt.targetWidth, tt.targetHeight)

			if resizedImg == nil {
				t.Fatalf("Expected resized image but got nil")
			}

			bounds := resizedImg.Bounds()
			if bounds.Dx() != tt.targetWidth {
				t.Errorf("Expected width %d, got %d", tt.targetWidth, bounds.Dx())
			}
			if bounds.Dy() != tt.targetHeight {
				t.Errorf("Expected height %d, got %d", tt.targetHeight, bounds.Dy())
			}
		})
	}
}

func TestImageRotate(t *testing.T) {
	tests := []struct {
		name         string
		orientation  int
		originalSize image.Rectangle
		expectedSize image.Rectangle
	}{
		{
			name:         "rotate 90 degrees",
			orientation:  90,
			originalSize: image.Rect(0, 0, 100, 50),
			expectedSize: image.Rect(0, 0, 50, 100),
		},
		{
			name:         "rotate 180 degrees",
			orientation:  180,
			originalSize: image.Rect(0, 0, 100, 50),
			expectedSize: image.Rect(0, 0, 100, 50),
		},
		{
			name:         "rotate 270 degrees",
			orientation:  270,
			originalSize: image.Rect(0, 0, 100, 50),
			expectedSize: image.Rect(0, 0, 50, 100),
		},
		{
			name:         "no rotation (invalid orientation)",
			orientation:  0,
			originalSize: image.Rect(0, 0, 100, 50),
			expectedSize: image.Rect(0, 0, 100, 50),
		},
		{
			name:         "invalid orientation 45",
			orientation:  45,
			originalSize: image.Rect(0, 0, 100, 50),
			expectedSize: image.Rect(0, 0, 100, 50),
		},
		{
			name:         "square image rotate 90",
			orientation:  90,
			originalSize: image.Rect(0, 0, 100, 100),
			expectedSize: image.Rect(0, 0, 100, 100),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalImg := image.NewRGBA(tt.originalSize)

			rotatedImg := ImageRotate(tt.orientation, originalImg)

			if rotatedImg == nil {
				t.Fatalf("Expected rotated image but got nil")
			}

			bounds := rotatedImg.Bounds()
			if bounds.Dx() != tt.expectedSize.Dx() {
				t.Errorf("Expected width %d, got %d", tt.expectedSize.Dx(), bounds.Dx())
			}
			if bounds.Dy() != tt.expectedSize.Dy() {
				t.Errorf("Expected height %d, got %d", tt.expectedSize.Dy(), bounds.Dy())
			}
		})
	}
}

func TestImageRotateWithColors(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 3))

	img.Set(0, 0, color.RGBA{255, 0, 0, 255})
	img.Set(1, 0, color.RGBA{0, 255, 0, 255})
	img.Set(0, 1, color.RGBA{0, 0, 255, 255})
	img.Set(1, 1, color.RGBA{255, 255, 0, 255})
	img.Set(0, 2, color.RGBA{255, 0, 255, 255})
	img.Set(1, 2, color.RGBA{0, 255, 255, 255})

	tests := []struct {
		name        string
		orientation int
	}{
		{
			name:        "rotate 90 degrees with colors",
			orientation: 90,
		},
		{
			name:        "rotate 180 degrees with colors",
			orientation: 180,
		},
		{
			name:        "rotate 270 degrees with colors",
			orientation: 270,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rotatedImg := ImageRotate(tt.orientation, img)

			if rotatedImg == nil {
				t.Fatalf("Expected rotated image but got nil")
			}

			bounds := rotatedImg.Bounds()
			if bounds.Dx() == 0 || bounds.Dy() == 0 {
				t.Errorf("Rotated image has invalid dimensions: %dx%d", bounds.Dx(), bounds.Dy())
			}

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					c := rotatedImg.At(x, y)
					if c == nil {
						t.Errorf("Pixel at (%d, %d) is nil", x, y)
					}
				}
			}
		})
	}
}

func TestGenerateImageByWidthHeightWithNilImage(t *testing.T) {
	t.Run("resize nil image", func(t *testing.T) {
		result := GenerateImageByWidthHeight(nil, 100, 100)
		if result != nil {
			t.Errorf("Expected nil result for nil input image, got non-nil")
		}
	})
}

func TestImageRotateWithNilImage(t *testing.T) {
	tests := []struct {
		name        string
		orientation int
	}{
		{
			name:        "rotate nil image 90 degrees",
			orientation: 90,
		},
		{
			name:        "rotate nil image 180 degrees",
			orientation: 180,
		},
		{
			name:        "rotate nil image 270 degrees",
			orientation: 270,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ImageRotate(tt.orientation, nil)
			if result != nil {
				t.Errorf("Expected nil result for nil input image, got non-nil")
			}
		})
	}
}

func TestLoadImageWithDecodeWithLargeImage(t *testing.T) {
	t.Run("load large jpeg image", func(t *testing.T) {
		buf := new(bytes.Buffer)
		img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
		jpeg.Encode(buf, img, nil)

		tmpFile := filepath.Join(os.TempDir(), "large.jpg")
		if err := os.WriteFile(tmpFile, buf.Bytes(), 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		defer os.Remove(tmpFile)

		loadedImg, err := LoadImageWithDecode(tmpFile)

		if HasError(err) {
			t.Errorf("Unexpected error: %v", err)
		}
		if loadedImg == nil {
			t.Errorf("Expected image but got nil")
		}
		bounds := loadedImg.Bounds()
		if bounds.Dx() != 1000 || bounds.Dy() != 1000 {
			t.Errorf("Expected image size 1000x1000, got %dx%d", bounds.Dx(), bounds.Dy())
		}
	})
}
