package engine

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/internal"
	"WaterMark/pkg"
)

func TestCacheGetImageExif(t *testing.T) {
	internal.InitAppConfigsAndRes()
	InitAllTools()

	tests := []struct {
		name string
		path string
	}{
		{
			name: "Cache get image exif with non-existent file",
			path: "non_existent_file.jpg",
		},
		{
			name: "Cache get image exif with empty path",
			path: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := CacheGetImageExif(tt.path)

			if pkg.HasError(err) {
				t.Logf("CacheGetImageExif() returned expected error: %v", err)
			} else {
				t.Logf("CacheGetImageExif() returned info for %s", tt.path)
				if info.File != tt.path {
					t.Errorf("CacheGetImageExif() File = %v, want %v", info.File, tt.path)
				}
			}
		})
	}
}

func TestCacheGetImageExifWithCache(t *testing.T) {
	internal.InitAppConfigsAndRes()
	InitAllTools()

	tests := []struct {
		name string
		path string
	}{
		{
			name: "Cache get image exif with cache hit",
			path: "test_cached_file.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md5 := pkg.GetStrMD5(tt.path)
			testExif := exiftool.FileMetadata{
				File: tt.path,
				Fields: map[string]any{
					"Make":  "Canon",
					"Model": "EOS 5D",
				},
				Err: nil,
			}

			exiftoolCache.Store(md5, testExif)

			info, err := CacheGetImageExif(tt.path)

			if pkg.HasError(err) {
				t.Errorf("CacheGetImageExif() returned error: %v", err)
			}

			if info.File != tt.path {
				t.Errorf("CacheGetImageExif() File = %v, want %v", info.File, tt.path)
			}

			if info.Fields["Make"] != "Canon" {
				t.Errorf("CacheGetImageExif() Fields[Make] = %v, want Canon", info.Fields["Make"])
			}

			exiftoolCache.Delete(md5)
		})
	}
}

func TestExiftoolCacheOperations(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test exiftool cache operations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPath := "test_cache_operations.jpg"
			md5 := pkg.GetStrMD5(testPath)
			testExif := exiftool.FileMetadata{
				File: testPath,
				Fields: map[string]any{
					"Make":  "Nikon",
					"Model": "D850",
				},
				Err: nil,
			}

			exiftoolCache.Store(md5, testExif)

			loaded, ok := exiftoolCache.Load(md5)
			if !ok {
				t.Error("exiftoolCache.Load() failed to load stored value")
			}

			exifData, ok := loaded.(exiftool.FileMetadata)
			if !ok {
				t.Error("exiftoolCache.Load() returned wrong type")
			}

			if exifData.Fields["Make"] != "Nikon" {
				t.Errorf("exiftoolCache.Load() Fields[Make] = %v, want Nikon", exifData.Fields["Make"])
			}

			exiftoolCache.Delete(md5)

			_, ok = exiftoolCache.Load(md5)
			if ok {
				t.Error("exiftoolCache.Delete() failed to delete value")
			}
		})
	}
}

func TestExiftoolCacheRange(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test exiftool cache range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPaths := []string{"test1.jpg", "test2.jpg", "test3.jpg"}

			for _, path := range testPaths {
				md5 := pkg.GetStrMD5(path)
				exiftoolCache.Store(md5, exiftool.FileMetadata{
					File:   path,
					Fields: map[string]any{},
					Err:    nil,
				})
			}

			count := 0
			exiftoolCache.Range(func(key, value any) bool {
				count++
				return true
			})

			if count < len(testPaths) {
				t.Errorf("exiftoolCache.Range() found %d items, want at least %d", count, len(testPaths))
			}

			for _, path := range testPaths {
				md5 := pkg.GetStrMD5(path)
				exiftoolCache.Delete(md5)
			}
		})
	}
}

func TestCacheGetImageExifWithRealImage(t *testing.T) {
	internal.InitAppConfigsAndRes()
	InitAllTools()

	testDir := "test_images"
	testFileName := "test_real_image.jpg"
	testFilePath := filepath.Join(testDir, testFileName)

	absPath, err := filepath.Abs(testFilePath)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x * 255 / 100), G: uint8(y * 255 / 100), B: 128, A: 255})
		}
	}

	file, err := os.Create(testFilePath)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}
	defer file.Close()

	if err := jpeg.Encode(file, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}

	file.Close()

	info, eErr := CacheGetImageExif(absPath)
	if pkg.HasError(eErr) {
		t.Errorf("CacheGetImageExif() returned error: %v", eErr)
	}

	if info.File != absPath {
		t.Errorf("CacheGetImageExif() File = %v, want %v", info.File, absPath)
	}

	info2, eErr := CacheGetImageExif(absPath)
	if pkg.HasError(eErr) {
		t.Errorf("CacheGetImageExif() second call returned error: %v", eErr)
	}

	if info2.File != absPath {
		t.Errorf("CacheGetImageExif() second call File = %v, want %v", info2.File, absPath)
	}

	if info.File != info2.File {
		t.Errorf("CacheGetImageExif() returned different File values: %v vs %v", info.File, info2.File)
	}
}
