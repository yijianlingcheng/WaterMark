package internal

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/pkg"
)

func TestCacheLoadImageWithDecode(t *testing.T) {
	originalImagesCache := imagesCache
	defer func() { imagesCache = originalImagesCache }()

	imagesCache = make(map[string]image.Image)

	tests := []struct {
		name      string
		setup     func() (string, func())
		wantCache bool
		wantErr   bool
	}{
		{
			name: "成功加载并缓存图片",
			setup: func() (string, func()) {
				tmpDir := t.TempDir()
				imgPath := filepath.Join(tmpDir, "test.jpg")

				img := image.NewRGBA(image.Rect(0, 0, 100, 100))
				for y := 0; y < 100; y++ {
					for x := 0; x < 100; x++ {
						img.Set(x, y, color.RGBA{255, 0, 0, 255})
					}
				}

				file, err := os.Create(imgPath)
				if err != nil {
					t.Fatalf("Failed to create test image: %v", err)
				}
				err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
				if err != nil {
					t.Fatalf("Failed to encode test image: %v", err)
				}
				file.Close()

				return imgPath, func() {}
			},
			wantCache: false,
			wantErr:   false,
		},
		{
			name: "从缓存中获取图片",
			setup: func() (string, func()) {
				tmpDir := t.TempDir()
				imgPath := filepath.Join(tmpDir, "test.jpg")

				img := image.NewRGBA(image.Rect(0, 0, 100, 100))
				for y := 0; y < 100; y++ {
					for x := 0; x < 100; x++ {
						img.Set(x, y, color.RGBA{255, 0, 0, 255})
					}
				}

				file, err := os.Create(imgPath)
				if err != nil {
					t.Fatalf("Failed to create test image: %v", err)
				}
				err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
				if err != nil {
					t.Fatalf("Failed to encode test image: %v", err)
				}
				file.Close()

				md5 := pkg.GetStrMD5(imgPath)
				imagesCache[md5] = img

				return imgPath, func() {}
			},
			wantCache: true,
			wantErr:   false,
		},
		{
			name: "文件不存在",
			setup: func() (string, func()) {
				imgPath := "nonexistent.jpg"
				return imgPath, func() {}
			},
			wantCache: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imgPath, cleanup := tt.setup()
			defer cleanup()

			img, err := CacheLoadImageWithDecode(imgPath)

			if tt.wantErr {
				if !pkg.HasError(err) {
					t.Errorf("CacheLoadImageWithDecode() expected error, got NoError")
				}
				if img != nil {
					t.Errorf("CacheLoadImageWithDecode() expected nil image on error, got %v", img)
				}
			} else {
				if pkg.HasError(err) {
					t.Errorf("CacheLoadImageWithDecode() unexpected error = %v", err)
				}
				if img == nil {
					t.Errorf("CacheLoadImageWithDecode() expected non-nil image, got nil")
				}
				if tt.wantCache {
					md5 := pkg.GetStrMD5(imgPath)
					if _, ok := imagesCache[md5]; !ok {
						t.Errorf("CacheLoadImageWithDecode() expected image to be cached")
					}
				}
			}
		})
	}
}

func TestImportImageFiles(t *testing.T) {
	originalImagesCache := imagesCache
	defer func() { imagesCache = originalImagesCache }()

	imagesCache = make(map[string]image.Image)

	tests := []struct {
		name       string
		setup      func() ([]string, []exiftool.FileMetadata, func())
		wantCached int
	}{
		{
			name: "成功导入多张图片",
			setup: func() ([]string, []exiftool.FileMetadata, func()) {
				tmpDir := t.TempDir()
				paths := make([]string, 3)
				exifInfos := make([]exiftool.FileMetadata, 3)

				for i := 0; i < 3; i++ {
					imgPath := filepath.Join(tmpDir, fmt.Sprintf("test%d.jpg", i))
					paths[i] = imgPath

					img := image.NewRGBA(image.Rect(0, 0, 100, 100))
					for y := 0; y < 100; y++ {
						for x := 0; x < 100; x++ {
							img.Set(x, y, color.RGBA{byte(i * 80), 0, 0, 255})
						}
					}

					file, err := os.Create(imgPath)
					if err != nil {
						t.Fatalf("Failed to create test image: %v", err)
					}
					err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
					if err != nil {
						t.Fatalf("Failed to encode test image: %v", err)
					}
					file.Close()

					exifInfos[i] = exiftool.FileMetadata{
						Fields: map[string]any{
							"Orientation": "Horizontal (normal)",
						},
					}
				}

				return paths, exifInfos, func() {}
			},
			wantCached: 3,
		},
		{
			name: "导入带旋转信息的图片",
			setup: func() ([]string, []exiftool.FileMetadata, func()) {
				tmpDir := t.TempDir()
				paths := make([]string, 2)
				exifInfos := make([]exiftool.FileMetadata, 2)

				for i := 0; i < 2; i++ {
					imgPath := filepath.Join(tmpDir, fmt.Sprintf("test%d.jpg", i))
					paths[i] = imgPath

					img := image.NewRGBA(image.Rect(0, 0, 100, 100))
					for y := 0; y < 100; y++ {
						for x := 0; x < 100; x++ {
							img.Set(x, y, color.RGBA{0, 0, 255, 255})
						}
					}

					file, err := os.Create(imgPath)
					if err != nil {
						t.Fatalf("Failed to create test image: %v", err)
					}
					err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
					if err != nil {
						t.Fatalf("Failed to encode test image: %v", err)
					}
					file.Close()

					orientations := []string{"Rotate 90 CW", "Rotate 180 CW"}
					exifInfos[i] = exiftool.FileMetadata{
						Fields: map[string]any{
							"Orientation": orientations[i],
						},
					}
				}

				return paths, exifInfos, func() {}
			},
			wantCached: 2,
		},
		{
			name: "部分文件不存在",
			setup: func() ([]string, []exiftool.FileMetadata, func()) {
				tmpDir := t.TempDir()
				paths := make([]string, 2)
				exifInfos := make([]exiftool.FileMetadata, 2)

				imgPath := filepath.Join(tmpDir, "test.jpg")
				paths[0] = imgPath

				img := image.NewRGBA(image.Rect(0, 0, 100, 100))
				for y := 0; y < 100; y++ {
					for x := 0; x < 100; x++ {
						img.Set(x, y, color.RGBA{0, 0, 255, 255})
					}
				}

				file, err := os.Create(imgPath)
				if err != nil {
					t.Fatalf("Failed to create test image: %v", err)
				}
				err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
				if err != nil {
					t.Fatalf("Failed to encode test image: %v", err)
				}
				file.Close()

				paths[1] = "nonexistent.jpg"

				exifInfos[0] = exiftool.FileMetadata{
					Fields: map[string]any{
						"Orientation": "Horizontal (normal)",
					},
				}
				exifInfos[1] = exiftool.FileMetadata{
					Fields: map[string]any{
						"Orientation": "Horizontal (normal)",
					},
				}

				return paths, exifInfos, func() {}
			},
			wantCached: 1,
		},
		{
			name: "空路径列表",
			setup: func() ([]string, []exiftool.FileMetadata, func()) {
				return []string{}, []exiftool.FileMetadata{}, func() {}
			},
			wantCached: 0,
		},
		{
			name: "从缓存中导入图片",
			setup: func() ([]string, []exiftool.FileMetadata, func()) {
				tmpDir := t.TempDir()
				paths := make([]string, 2)
				exifInfos := make([]exiftool.FileMetadata, 2)

				for i := 0; i < 2; i++ {
					imgPath := filepath.Join(tmpDir, fmt.Sprintf("test%d.jpg", i))
					paths[i] = imgPath

					img := image.NewRGBA(image.Rect(0, 0, 100, 100))
					for y := 0; y < 100; y++ {
						for x := 0; x < 100; x++ {
							img.Set(x, y, color.RGBA{byte(i * 100), byte(i * 100), byte(i * 100), 255})
						}
					}

					file, err := os.Create(imgPath)
					if err != nil {
						t.Fatalf("Failed to create test image: %v", err)
					}
					err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
					if err != nil {
						t.Fatalf("Failed to encode test image: %v", err)
					}
					file.Close()

					md5 := pkg.GetStrMD5(imgPath)
					imagesCache[md5] = img

					exifInfos[i] = exiftool.FileMetadata{
						Fields: map[string]any{
							"Orientation": "Horizontal (normal)",
						},
					}
				}

				return paths, exifInfos, func() {}
			},
			wantCached: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths, exifInfos, cleanup := tt.setup()
			defer cleanup()

			ImportImageFiles(paths, exifInfos)

			if len(imagesCache) != tt.wantCached {
				t.Errorf("ImportImageFiles() cached %d images, want %d", len(imagesCache), tt.wantCached)
			}

			for _, path := range paths {
				md5 := pkg.GetStrMD5(path)
				if _, ok := imagesCache[md5]; ok {
					img := imagesCache[md5]
					if img == nil {
						t.Errorf("ImportImageFiles() cached nil image for path %s", path)
					}
				}
			}
		})
	}
}

func TestImportImageFilesConcurrency(t *testing.T) {
	originalImagesCache := imagesCache
	defer func() { imagesCache = originalImagesCache }()

	imagesCache = make(map[string]image.Image)

	tmpDir := t.TempDir()
	paths := make([]string, 20)
	exifInfos := make([]exiftool.FileMetadata, 20)

	for i := 0; i < 20; i++ {
		imgPath := filepath.Join(tmpDir, fmt.Sprintf("test%d.jpg", i))
		paths[i] = imgPath

		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		for y := 0; y < 100; y++ {
			for x := 0; x < 100; x++ {
				img.Set(x, y, color.RGBA{byte(i * 10), byte(i * 10), byte(i * 10), 255})
			}
		}

		file, err := os.Create(imgPath)
		if err != nil {
			t.Fatalf("Failed to create test image: %v", err)
		}
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
		if err != nil {
			t.Fatalf("Failed to encode test image: %v", err)
		}
		file.Close()

		exifInfos[i] = exiftool.FileMetadata{
			Fields: map[string]any{
				"Orientation": "Horizontal (normal)",
			},
		}
	}

	ImportImageFiles(paths, exifInfos)

	if len(imagesCache) != 20 {
		t.Errorf("ImportImageFiles() cached %d images, want 20", len(imagesCache))
	}

	for _, path := range paths {
		md5 := pkg.GetStrMD5(path)
		if _, ok := imagesCache[md5]; !ok {
			t.Errorf("ImportImageFiles() did not cache image for path %s", path)
		}
	}
}
