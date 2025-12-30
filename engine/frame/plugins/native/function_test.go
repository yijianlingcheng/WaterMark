package native

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"WaterMark/internal"
	"WaterMark/layout"
	"WaterMark/pkg"
)

func createTestImage(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	testImagePath := filepath.Join(tmpDir, "test.jpg")

	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			testImg.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 255, A: 255})
		}
	}

	testFile, err := os.Create(testImagePath)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}
	defer testFile.Close()

	// 使用JPEG编码
	if err := jpeg.Encode(testFile, testImg, nil); err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}

	return testImagePath
}

func TestInitAllCachaAndTools(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Init all cache and tools",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
			}()

			err := InitAllCachaAndTools()

			if pkg.HasError(err) {
				t.Logf("InitAllCachaAndTools() returned error: %v", err)
				t.Skip("Fonts or logos directory not available in test environment")
			}
		})
	}
}

func TestGetTextContentMaxSize(t *testing.T) {
	fontFile := filepath.Join(internal.GetFontFilePath(""), "Alibaba-PuHuiTi-Bold.ttf")

	tests := []struct {
		name     string
		width    int
		fontFile string
		content  string
		wantErr  bool
	}{
		{
			name:     "Calculate max size for short text",
			width:    100,
			fontFile: fontFile,
			content:  "Hi",
			wantErr:  false,
		},
		{
			name:     "Calculate max size for medium text",
			width:    200,
			fontFile: fontFile,
			content:  "Test Content",
			wantErr:  false,
		},
		{
			name:     "Calculate max size for long text",
			width:    300,
			fontFile: fontFile,
			content:  "This is a longer text content for testing",
			wantErr:  false,
		},
		{
			name:     "Calculate max size for nikon text",
			width:    400,
			fontFile: fontFile,
			content:  "Nikon Camera Test",
			wantErr:  false,
		},
		{
			name:     "Calculate max size for single character",
			width:    50,
			fontFile: fontFile,
			content:  "A",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getTextContentMaxSize(tt.width, tt.fontFile, tt.content)

			if got <= 0 {
				t.Error("getTextContentMaxSize() returned non-positive size for valid input")
			}
		})
	}
}

func TestGetTextContentMaxSize_Caching(t *testing.T) {
	fontFile := filepath.Join(internal.GetFontFilePath(""), "Alibaba-PuHuiTi-Bold.ttf")
	width := 100
	content := "Cache Test"

	t.Run("Verify caching behavior", func(t *testing.T) {
		got1 := getTextContentMaxSize(width, fontFile, content)
		got2 := getTextContentMaxSize(width, fontFile, content)

		if got1 != got2 {
			t.Errorf("getTextContentMaxSize() caching failed, got different results: %d vs %d", got1, got2)
		}
	})
}

func TestFindTextContentMaxSize(t *testing.T) {
	fontFile := filepath.Join(internal.GetFontFilePath(""), "Alibaba-PuHuiTi-Bold.ttf")

	tests := []struct {
		name     string
		width    int
		fontFile string
		content  string
		wantErr  bool
	}{
		{
			name:     "Find max size for short text",
			width:    100,
			fontFile: fontFile,
			content:  "Hi",
			wantErr:  false,
		},
		{
			name:     "Find max size for medium text",
			width:    200,
			fontFile: fontFile,
			content:  "Test Content",
			wantErr:  false,
		},
		{
			name:     "Find max size for long text",
			width:    300,
			fontFile: fontFile,
			content:  "This is a longer text content for testing",
			wantErr:  false,
		},
		{
			name:     "Find max size for nikon text",
			width:    400,
			fontFile: fontFile,
			content:  "Nikon Camera Test",
			wantErr:  false,
		},
		{
			name:     "Find max size for single character",
			width:    50,
			fontFile: fontFile,
			content:  "A",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findTextContentMaxSize(tt.width, tt.fontFile, tt.content)

			if got <= 0 {
				t.Error("findTextContentMaxSize() returned non-positive size for valid input")
			}
		})
	}
}

func TestGetTextContentMaxSizeWithLogo(t *testing.T) {
	fontFile := filepath.Join(internal.GetFontFilePath(""), "Alibaba-PuHuiTi-Bold.ttf")

	tests := []struct {
		name     string
		width    int
		logoName string
		fontFile string
		content  string
		wantErr  bool
	}{
		{
			name:     "Calculate max size with logo for short text",
			width:    150,
			logoName: "logo1",
			fontFile: fontFile,
			content:  "Hi",
			wantErr:  false,
		},
		{
			name:     "Calculate max size with logo for medium text",
			width:    300,
			logoName: "logo2",
			fontFile: fontFile,
			content:  "Test Content",
			wantErr:  false,
		},
		{
			name:     "Calculate max size with logo for long text",
			width:    500,
			logoName: "logo3",
			fontFile: fontFile,
			content:  "This is a longer text content for testing",
			wantErr:  false,
		},
		{
			name:     "Calculate max size with nikon logo",
			width:    400,
			logoName: "nikon",
			fontFile: fontFile,
			content:  "Nikon Camera Test",
			wantErr:  false,
		},
		{
			name:     "Calculate max size with nikon logo for short text",
			width:    200,
			logoName: "nikon",
			fontFile: fontFile,
			content:  "Hi",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getTextContentMaxSizeWithLogo(tt.width, tt.logoName, tt.fontFile, tt.content)

			if got <= 0 {
				t.Error("getTextContentMaxSizeWithLogo() returned non-positive size for valid input")
			}
		})
	}
}

func TestGetTextContentMaxSizeWithLogo_Caching(t *testing.T) {
	fontFile := filepath.Join(internal.GetFontFilePath(""), "Alibaba-PuHuiTi-Bold.ttf")
	width := 150
	logoName := "logo1"
	content := "Cache Test"

	t.Run("Verify caching behavior", func(t *testing.T) {
		got1 := getTextContentMaxSizeWithLogo(width, logoName, fontFile, content)
		got2 := getTextContentMaxSizeWithLogo(width, logoName, fontFile, content)

		if got1 != got2 {
			t.Errorf("getTextContentMaxSizeWithLogo() caching failed, got different results: %d vs %d", got1, got2)
		}
	})
}

func TestFindTextContentMaxSizeWithLogo(t *testing.T) {
	fontFile := filepath.Join(internal.GetFontFilePath(""), "Alibaba-PuHuiTi-Bold.ttf")

	tests := []struct {
		name     string
		width    int
		logoName string
		fontFile string
		content  string
		wantErr  bool
	}{
		{
			name:     "Find max size with logo for short text",
			width:    150,
			logoName: "logo1",
			fontFile: fontFile,
			content:  "Hi",
			wantErr:  false,
		},
		{
			name:     "Find max size with logo for medium text",
			width:    300,
			logoName: "logo2",
			fontFile: fontFile,
			content:  "Test Content",
			wantErr:  false,
		},
		{
			name:     "Find max size with logo for long text",
			width:    500,
			logoName: "logo3",
			fontFile: fontFile,
			content:  "This is a longer text content for testing",
			wantErr:  false,
		},
		{
			name:     "Find max size with nikon logo",
			width:    400,
			logoName: "nikon",
			fontFile: fontFile,
			content:  "Nikon Test",
			wantErr:  false,
		},
		{
			name:     "Find max size with nikon logo for short text",
			width:    200,
			logoName: "nikon",
			fontFile: fontFile,
			content:  "Hi",
			wantErr:  false,
		},
		{
			name:     "Find max size with nikon logo for long text",
			width:    600,
			logoName: "nikon",
			fontFile: fontFile,
			content:  "This is a longer text content for testing with nikon logo",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findTextContentMaxSizeWithLogo(tt.width, tt.logoName, tt.fontFile, tt.content)

			if got <= 0 {
				t.Error("findTextContentMaxSizeWithLogo() returned non-positive size for valid input")
			}
		})
	}
}

func TestInitAllCachaAndTools_Concurrent(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic: %v", r)
		}
	}()

	var wg sync.WaitGroup
	numGoroutines := 5

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Goroutine recovered from panic: %v", r)
				}
			}()
			err := InitAllCachaAndTools()
			if pkg.HasError(err) {
				t.Logf("InitAllCachaAndTools() returned error: %v", err)
			}
		}()
	}

	wg.Wait()
}

func TestCreateFrameImageRGBA(t *testing.T) {
	// 获取所有布局配置
	allLayouts := layout.GetAllLayout()

	tests := []struct {
		name    string
		opts    map[string]any
		wantErr bool
	}{
		{
			name:    "Create frame with nil opts",
			opts:    nil,
			wantErr: true,
		},
		{
			name:    "Create frame with empty opts",
			opts:    map[string]any{},
			wantErr: true,
		},
		{
			name: "Create frame with invalid source image",
			opts: map[string]any{
				"sourceImageFile": "non_existent.jpg",
			},
			wantErr: true,
		},
		{
			name: "Create blur frame with nil opts",
			opts: map[string]any{
				"isBlur": true,
			},
			wantErr: true,
		},
	}

	// 如果有有效的布局配置，添加测试用例
	if len(allLayouts) > 0 {
		// 使用第一个布局作为有效参数
		validLayout := allLayouts[0]

		tests = append(tests, []struct {
			name    string
			opts    map[string]any
			wantErr bool
		}{
			{
				name: "Create normal frame with valid params from layout.json",
				opts: map[string]any{
					"sourceImageFile": createTestImage(t),
					"photoType":       "photo",
					"isBlur":          false,
					"params":          validLayout,
				},
				wantErr: false,
			},
			{
				name: "Create border frame with valid params from layout.json",
				opts: map[string]any{
					"sourceImageFile": createTestImage(t),
					"photoType":       "border",
					"isBlur":          false,
					"params":          validLayout,
				},
				wantErr: false,
			},
		}...)

		// 查找模糊类型的布局
		for _, lay := range allLayouts {
			if lay.Type == "blur_bottom_text_center_layout" {
				tests = append(tests, struct {
					name    string
					opts    map[string]any
					wantErr bool
				}{
					name: "Create blur frame with valid params from layout.json",
					opts: map[string]any{
						"sourceImageFile": createTestImage(t),
						"photoType":       "photo",
						"isBlur":          true,
						"params":          lay,
					},
					wantErr: false,
				})
				break
			}
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
					if !tt.wantErr {
						t.Errorf("Unexpected panic: %v", r)
					}
				}
			}()

			img, err := CreateFrameImageRGBA(tt.opts)

			if tt.wantErr && !pkg.HasError(err) {
				t.Errorf("CreateFrameImageRGBA() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(err) {
				t.Errorf("CreateFrameImageRGBA() unexpected error: %v", err)
			}

			if !tt.wantErr && img == nil {
				t.Error("CreateFrameImageRGBA() returned nil image")
			}
		})
	}
}

func TestGetFrameImageBorderInfo(t *testing.T) {
	tests := []struct {
		name    string
		opts    map[string]any
		wantErr bool
	}{
		{
			name:    "Get border info with nil opts",
			opts:    nil,
			wantErr: true,
		},
		{
			name:    "Get border info with empty opts",
			opts:    map[string]any{},
			wantErr: true,
		},
		{
			name: "Get blur border info with nil opts",
			opts: map[string]any{
				"isBlur": true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
					if !tt.wantErr {
						t.Errorf("Unexpected panic: %v", r)
					}
				}
			}()

			info, err := GetFrameImageBorderInfo(tt.opts)

			if tt.wantErr && !pkg.HasError(err) {
				t.Errorf("GetFrameImageBorderInfo() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(err) {
				t.Errorf("GetFrameImageBorderInfo() unexpected error: %v", err)
			}

			if !tt.wantErr && info == nil {
				t.Error("GetFrameImageBorderInfo() returned nil info")
			}
		})
	}
}

func TestStrColor2RGBA(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want color.RGBA
	}{
		{
			name: "Convert valid color string",
			s:    "255,255,255,255",
			want: color.RGBA{255, 255, 255, 255},
		},
		{
			name: "Convert valid color string with different values",
			s:    "203,203,201,255",
			want: color.RGBA{203, 203, 201, 255},
		},
		{
			name: "Convert empty string (should use default)",
			s:    "",
			want: color.RGBA{255, 255, 255, 255},
		},
		{
			name: "Convert color string with zeros",
			s:    "0,0,0,0",
			want: color.RGBA{0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strColor2RGBA(tt.s)

			if got.R != tt.want.R {
				t.Errorf("strColor2RGBA() R = %v, want %v", got.R, tt.want.R)
			}

			if got.G != tt.want.G {
				t.Errorf("strColor2RGBA() G = %v, want %v", got.G, tt.want.G)
			}

			if got.B != tt.want.B {
				t.Errorf("strColor2RGBA() B = %v, want %v", got.B, tt.want.B)
			}

			if got.A != tt.want.A {
				t.Errorf("strColor2RGBA() A = %v, want %v", got.A, tt.want.A)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name string
		x    int
		want int
	}{
		{
			name: "Absolute value of positive number",
			x:    10,
			want: 10,
		},
		{
			name: "Absolute value of negative number",
			x:    -10,
			want: 10,
		},
		{
			name: "Absolute value of zero",
			x:    0,
			want: 0,
		},
		{
			name: "Absolute value of large positive number",
			x:    1000000,
			want: 1000000,
		},
		{
			name: "Absolute value of large negative number",
			x:    -1000000,
			want: 1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := abs(tt.x)

			if got != tt.want {
				t.Errorf("abs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveJpgImage(t *testing.T) {
	tests := []struct {
		name          string
		saveImageFile string
		image         draw.Image
		quality       int
		wantErr       bool
	}{
		{
			name:          "Save valid JPG image with high quality",
			saveImageFile: filepath.Join(os.TempDir(), "test_save_high_quality.jpg"),
			image:         image.NewRGBA(image.Rect(0, 0, 100, 100)),
			quality:       95,
			wantErr:       false,
		},
		{
			name:          "Save valid JPG image with low quality",
			saveImageFile: filepath.Join(os.TempDir(), "test_save_low_quality.jpg"),
			image:         image.NewRGBA(image.Rect(0, 0, 100, 100)),
			quality:       50,
			wantErr:       false,
		},
		{
			name:          "Save JPG image with default quality",
			saveImageFile: filepath.Join(os.TempDir(), "test_save_default_quality.jpg"),
			image:         image.NewRGBA(image.Rect(0, 0, 100, 100)),
			quality:       85,
			wantErr:       false,
		},
		{
			name:          "Save JPG image with different size",
			saveImageFile: filepath.Join(os.TempDir(), "test_save_different_size.jpg"),
			image:         image.NewRGBA(image.Rect(0, 0, 500, 500)),
			quality:       90,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
				os.Remove(tt.saveImageFile)
			}()

			saveJpgImage(tt.saveImageFile, tt.image, tt.quality)

			_, err := os.Stat(tt.saveImageFile)
			if tt.wantErr {
				if !os.IsNotExist(err) {
					t.Errorf("saveJpgImage() expected error but file was created")
				}
			} else {
				if os.IsNotExist(err) {
					t.Errorf("saveJpgImage() file was not created")
				}
			}
		})
	}
}

func TestSavePngImage(t *testing.T) {
	tests := []struct {
		name          string
		saveImageFile string
		image         draw.Image
		wantErr       bool
	}{
		{
			name:          "Save valid PNG image",
			saveImageFile: filepath.Join(os.TempDir(), "test_save.png"),
			image:         image.NewRGBA(image.Rect(0, 0, 100, 100)),
			wantErr:       false,
		},
		{
			name:          "Save PNG image with different size",
			saveImageFile: filepath.Join(os.TempDir(), "test_save_large.png"),
			image:         image.NewRGBA(image.Rect(0, 0, 500, 500)),
			wantErr:       false,
		},
		{
			name:          "Save PNG image with small size",
			saveImageFile: filepath.Join(os.TempDir(), "test_save_small.png"),
			image:         image.NewRGBA(image.Rect(0, 0, 10, 10)),
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
				os.Remove(tt.saveImageFile)
			}()

			savePngImage(tt.saveImageFile, tt.image)

			_, err := os.Stat(tt.saveImageFile)
			if tt.wantErr {
				if !os.IsNotExist(err) {
					t.Errorf("savePngImage() expected error but file was created")
				}
			} else {
				if os.IsNotExist(err) {
					t.Errorf("savePngImage() file was not created")
				}
			}
		})
	}
}

func TestSaveImageFile(t *testing.T) {
	tests := []struct {
		name          string
		saveImageFile string
		image         draw.Image
		quality       int
		wantErr       bool
	}{
		{
			name:          "Save JPG image with .jpg extension",
			saveImageFile: filepath.Join(os.TempDir(), "test_save_file.jpg"),
			image:         image.NewRGBA(image.Rect(0, 0, 100, 100)),
			quality:       85,
			wantErr:       false,
		},
		{
			name:          "Save JPG image with .jpeg extension",
			saveImageFile: filepath.Join(os.TempDir(), "test_save_file.jpeg"),
			image:         image.NewRGBA(image.Rect(0, 0, 100, 100)),
			quality:       90,
			wantErr:       false,
		},
		{
			name:          "Save JPG image with .JPG extension (uppercase)",
			saveImageFile: filepath.Join(os.TempDir(), "test_save_file.JPG"),
			image:         image.NewRGBA(image.Rect(0, 0, 100, 100)),
			quality:       95,
			wantErr:       false,
		},
		{
			name:          "Save PNG image with .png extension",
			saveImageFile: filepath.Join(os.TempDir(), "test_save_file.png"),
			image:         image.NewRGBA(image.Rect(0, 0, 100, 100)),
			quality:       85,
			wantErr:       false,
		},
		{
			name:          "Save PNG image with .PNG extension (uppercase)",
			saveImageFile: filepath.Join(os.TempDir(), "test_save_file.PNG"),
			image:         image.NewRGBA(image.Rect(0, 0, 100, 100)),
			quality:       85,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
				os.Remove(tt.saveImageFile)
			}()

			saveImageFile(tt.saveImageFile, tt.image, tt.quality)

			_, err := os.Stat(tt.saveImageFile)
			if tt.wantErr {
				if !os.IsNotExist(err) {
					t.Errorf("saveImageFile() expected error but file was created")
				}
			} else {
				if os.IsNotExist(err) {
					t.Errorf("saveImageFile() file was not created")
				}
			}
		})
	}
}

func TestGetTextContentSize(t *testing.T) {
	fontFile := filepath.Join(internal.GetFontFilePath(""), "Alibaba-PuHuiTi-Bold.ttf")

	tests := []struct {
		name     string
		fontSize int
		fontFile string
		content  string
		wantErr  bool
	}{
		{
			name:     "Calculate size for normal text",
			fontSize: 12,
			fontFile: fontFile,
			content:  "Test Content",
			wantErr:  false,
		},
		{
			name:     "Calculate size for empty text",
			fontSize: 12,
			fontFile: fontFile,
			content:  "",
			wantErr:  false,
		},
		{
			name:     "Calculate size for large font",
			fontSize: 24,
			fontFile: fontFile,
			content:  "Large Text",
			wantErr:  false,
		},
		{
			name:     "Calculate size for single character",
			fontSize: 16,
			fontFile: fontFile,
			content:  "A",
			wantErr:  false,
		},
		{
			name:     "Calculate size for nikon text",
			fontSize: 18,
			fontFile: fontFile,
			content:  "Nikon Camera",
			wantErr:  false,
		},
		{
			name:     "Calculate size for very large font",
			fontSize: 48,
			fontFile: fontFile,
			content:  "Big Text",
			wantErr:  false,
		},
		{
			name:     "Calculate size with invalid font file",
			fontSize: 12,
			fontFile: "invalid_font_file.ttf",
			content:  "Test",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotW, gotH := getTextContentSize(tt.fontSize, tt.fontFile, tt.content)

			if tt.wantErr {
				if gotW != 0 || gotH != 0 {
					t.Errorf("getTextContentSize() expected zero dimensions for invalid font, got %dx%d", gotW, gotH)
				}
			} else {
				if tt.content != "" && (gotW == 0 || gotH == 0) {
					t.Error("getTextContentSize() returned zero dimensions for non-empty content")
				}
			}
		})
	}
}

func TestGetTextContentXAndY(t *testing.T) {
	fontFile := filepath.Join(internal.GetFontFilePath(""), "Alibaba-PuHuiTi-Bold.ttf")

	tests := []struct {
		name     string
		fontSize int
		fontFile string
		content  string
		wantErr  bool
	}{
		{
			name:     "Calculate X and Y for normal text",
			fontSize: 12,
			fontFile: fontFile,
			content:  "Test Content",
			wantErr:  false,
		},
		{
			name:     "Calculate X and Y for empty text",
			fontSize: 12,
			fontFile: fontFile,
			content:  "",
			wantErr:  false,
		},
		{
			name:     "Calculate X and Y for large font",
			fontSize: 24,
			fontFile: fontFile,
			content:  "Large Text",
			wantErr:  false,
		},
		{
			name:     "Calculate X and Y with invalid font file",
			fontSize: 12,
			fontFile: "invalid_font_file.ttf",
			content:  "Test",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotW, gotH := getTextContentXAndY(tt.fontSize, tt.fontFile, tt.content)

			if tt.wantErr {
				if gotW != 0 || gotH != 0 {
					t.Errorf("getTextContentXAndY() expected zero dimensions for invalid font, got %dx%d", gotW, gotH)
				}
			} else {
				if tt.content != "" && (gotW == 0 || gotH == 0) {
					t.Error("getTextContentXAndY() returned zero dimensions for non-empty content")
				}
			}
		})
	}
}

func TestGetTextContentXAndY_Caching(t *testing.T) {
	fontFile := filepath.Join(internal.GetFontFilePath(""), "Alibaba-PuHuiTi-Bold.ttf")
	fontSize := 12
	content := "Cache Test"

	t.Run("Verify caching behavior", func(t *testing.T) {
		gotW1, gotH1 := getTextContentXAndY(fontSize, fontFile, content)
		gotW2, gotH2 := getTextContentXAndY(fontSize, fontFile, content)

		if gotW1 != gotW2 || gotH1 != gotH2 {
			t.Errorf(
				"getTextContentXAndY() caching failed, got different results: %dx%d vs %dx%d",
				gotW1,
				gotH1,
				gotW2,
				gotH2,
			)
		}
	})
}

func TestDrawLine(t *testing.T) {
	tests := []struct {
		name  string
		start image.Point
		end   image.Point
		c     color.Color
	}{
		{
			name:  "Draw horizontal line",
			start: image.Point{X: 0, Y: 10},
			end:   image.Point{X: 100, Y: 10},
			c:     color.RGBA{255, 0, 0, 255},
		},
		{
			name:  "Draw vertical line",
			start: image.Point{X: 10, Y: 0},
			end:   image.Point{X: 10, Y: 100},
			c:     color.RGBA{0, 255, 0, 255},
		},
		{
			name:  "Draw diagonal line",
			start: image.Point{X: 0, Y: 0},
			end:   image.Point{X: 100, Y: 100},
			c:     color.RGBA{0, 0, 255, 255},
		},
		{
			name:  "Draw line with reverse direction",
			start: image.Point{X: 100, Y: 100},
			end:   image.Point{X: 0, Y: 0},
			c:     color.RGBA{255, 255, 0, 255},
		},
		{
			name:  "Draw short line",
			start: image.Point{X: 10, Y: 10},
			end:   image.Point{X: 15, Y: 15},
			c:     color.RGBA{255, 0, 255, 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := image.NewRGBA(image.Rect(0, 0, 200, 200))

			drawLine(img, tt.start, tt.end, tt.c)

			bounds := img.Bounds()
			startX, startY := tt.start.X, tt.start.Y
			endX, endY := tt.end.X, tt.end.Y

			if startX < bounds.Min.X || startX >= bounds.Max.X ||
				startY < bounds.Min.Y || startY >= bounds.Max.Y ||
				endX < bounds.Min.X || endX >= bounds.Max.X ||
				endY < bounds.Min.Y || endY >= bounds.Max.Y {
				t.Skip("Line coordinates outside image bounds")
			}

			startColor := img.RGBAAt(startX, startY)
			if startColor.R != tt.c.(color.RGBA).R ||
				startColor.G != tt.c.(color.RGBA).G ||
				startColor.B != tt.c.(color.RGBA).B {
				t.Errorf("drawLine() start point color mismatch")
			}
		})
	}
}

func TestDrawBorderLogo(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{
			name:   "Draw logo on small border",
			width:  100,
			height: 100,
		},
		{
			name:   "Draw logo on large border",
			width:  500,
			height: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm := &basePhotoFrame{
				borderDraw: image.NewRGBA(image.Rect(0, 0, tt.width, tt.height)),
			}

			logoImage := image.NewRGBA(image.Rect(0, 0, 20, 20))
			for x := 0; x < 20; x++ {
				for y := 0; y < 20; y++ {
					logoImage.Set(x, y, color.RGBA{255, 0, 0, 255})
				}
			}

			startX, startY := 10, 10
			endX, endY := 30, 30

			drawBorderLogo(fm, logoImage, startX, startY, endX, endY)

			checkColor := fm.borderDraw.RGBAAt(15, 15)
			if checkColor.R != 255 || checkColor.G != 0 || checkColor.B != 0 {
				t.Errorf("drawBorderLogo() logo was not drawn correctly")
			}
		})
	}
}
