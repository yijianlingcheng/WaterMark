package layout

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"WaterMark/pkg"
)

func createTestPNG(t *testing.T, width, height int, path string) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 128, A: 255})
		}
	}

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test PNG: %v", err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		t.Fatalf("Failed to encode test PNG: %v", err)
	}
}

func setupLogoTest(t *testing.T) string {
	tempDir := t.TempDir()

	createTestPNG(t, 100, 100, filepath.Join(tempDir, "camera1.png"))
	createTestPNG(t, 150, 150, filepath.Join(tempDir, "camera2.png"))
	createTestPNG(t, 200, 200, filepath.Join(tempDir, "sony.png"))
	createTestPNG(t, 120, 80, filepath.Join(tempDir, "canon.png"))

	return tempDir
}

func TestGetLogoNameByMake(t *testing.T) {
	logos = &logosMaps{
		logoMap: map[string]*Logo{
			"sony":  {Name: "sony"},
			"canon": {Name: "canon"},
			"nikon": {Name: "nikon"},
		},
	}

	tests := []struct {
		name     string
		makeName string
		expected string
	}{
		{
			name:     "Exact match - sony",
			makeName: "sony",
			expected: "sony",
		},
		{
			name:     "Exact match - canon",
			makeName: "canon",
			expected: "canon",
		},
		{
			name:     "Exact match - nikon",
			makeName: "nikon",
			expected: "nikon",
		},
		{
			name:     "Case insensitive match - Sony",
			makeName: "Sony",
			expected: "sony",
		},
		{
			name:     "Case insensitive match - CANON",
			makeName: "CANON",
			expected: "canon",
		},
		{
			name:     "Partial match - sony alpha",
			makeName: "sony alpha",
			expected: "sony",
		},
		{
			name:     "Partial match - canon eos",
			makeName: "canon eos",
			expected: "canon",
		},
		{
			name:     "No match - unknown brand",
			makeName: "unknown brand",
			expected: UNSUPPORT_LOGO,
		},
		{
			name:     "Empty string",
			makeName: "",
			expected: UNSUPPORT_LOGO,
		},
		{
			name:     "Mixed case partial match",
			makeName: "SoNy AlPhA",
			expected: "sony",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLogoNameByMake(tt.makeName)
			if result != tt.expected {
				t.Errorf("GetLogoNameByMake() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCheckLogoIsUnSupport(t *testing.T) {
	tests := []struct {
		name     string
		logoName string
		expected bool
	}{
		{
			name:     "Unsupported logo",
			logoName: UNSUPPORT_LOGO,
			expected: true,
		},
		{
			name:     "Supported logo - sony",
			logoName: "sony",
			expected: false,
		},
		{
			name:     "Supported logo - canon",
			logoName: "canon",
			expected: false,
		},
		{
			name:     "Empty string",
			logoName: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckLogoIsUnSupport(tt.logoName)
			if result != tt.expected {
				t.Errorf("CheckLogoIsUnSupport() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetLogoImageByName(t *testing.T) {
	logos = &logosMaps{
		logoMap: map[string]*Logo{
			"sony": {
				Name:   "sony",
				Width:  100,
				Height: 100,
				IsLoad: true,
			},
			"canon": {
				Name:   "canon",
				Width:  150,
				Height: 150,
				IsLoad: true,
			},
		},
	}

	tests := []struct {
		name        string
		logoName    string
		wantFound   bool
		wantWidth   int
		wantHeight  int
		wantIsLoad  bool
	}{
		{
			name:       "Find existing logo - sony",
			logoName:   "sony",
			wantFound:  true,
			wantWidth:  100,
			wantHeight: 100,
			wantIsLoad: true,
		},
		{
			name:       "Find existing logo - canon",
			logoName:   "canon",
			wantFound:  true,
			wantWidth:  150,
			wantHeight: 150,
			wantIsLoad: true,
		},
		{
			name:       "Find non-existing logo",
			logoName:   "nikon",
			wantFound:  false,
			wantWidth:  0,
			wantHeight: 0,
			wantIsLoad: false,
		},
		{
			name:       "Find with empty name",
			logoName:   "",
			wantFound:  false,
			wantWidth:  0,
			wantHeight: 0,
			wantIsLoad: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logo, err := GetLogoImageByName(tt.logoName)

			if tt.wantFound {
				if pkg.HasError(err) {
					t.Errorf("GetLogoImageByName() unexpected error: %v", err)
				}
				if logo.Width != tt.wantWidth {
					t.Errorf("GetLogoImageByName() Width = %v, want %v", logo.Width, tt.wantWidth)
				}
				if logo.Height != tt.wantHeight {
					t.Errorf("GetLogoImageByName() Height = %v, want %v", logo.Height, tt.wantHeight)
				}
				if logo.IsLoad != tt.wantIsLoad {
					t.Errorf("GetLogoImageByName() IsLoad = %v, want %v", logo.IsLoad, tt.wantIsLoad)
				}
			} else {
				if !pkg.HasError(err) {
					t.Errorf("GetLogoImageByName() expected error but got none")
				}
				if err.Code != pkg.ImageLogoNotFindError.Code {
					t.Errorf("GetLogoImageByName() error code = %v, want %v", err.Code, pkg.ImageLogoNotFindError.Code)
				}
			}
		})
	}
}

func TestGetLogoXAndYByNameAndHeight(t *testing.T) {
	logos = &logosMaps{
		logoMap: map[string]*Logo{
			"sony": {
				Name:   "sony",
				Width:  100,
				Height: 50,
				IsLoad: true,
			},
			"canon": {
				Name:   "canon",
				Width:  150,
				Height: 100,
				IsLoad: true,
			},
		},
	}

	tests := []struct {
		name          string
		logoName      string
		height        int
		wantHeight    int
		wantWidth     int
	}{
		{
			name:       "Calculate width for sony logo",
			logoName:   "sony",
			height:     100,
			wantHeight: 100,
			wantWidth:  200,
		},
		{
			name:       "Calculate width for canon logo",
			logoName:   "canon",
			height:     200,
			wantHeight: 200,
			wantWidth:  300,
		},
		{
			name:       "Calculate width with different height",
			logoName:   "sony",
			height:     50,
			wantHeight: 50,
			wantWidth:  100,
		},
		{
			name:       "Non-existing logo returns zero width",
			logoName:   "nikon",
			height:     100,
			wantHeight: 100,
			wantWidth:  0,
		},
		{
			name:       "Zero height",
			logoName:   "sony",
			height:     0,
			wantHeight: 0,
			wantWidth:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLogoXAndYByNameAndHeight(tt.logoName, tt.height)

			if result["height"] != tt.wantHeight {
				t.Errorf("GetLogoXAndYByNameAndHeight() height = %v, want %v", result["height"], tt.wantHeight)
			}
			if result["width"] != tt.wantWidth {
				t.Errorf("GetLogoXAndYByNameAndHeight() width = %v, want %v", result["width"], tt.wantWidth)
			}
		})
	}
}

func TestGetLogoImageByNameAndWidhtAndHeight(t *testing.T) {
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 50))

	logos = &logosMaps{
		logoMap: map[string]*Logo{
			"sony": {
				Name:      "sony",
				Width:     100,
				Height:    50,
				IsLoad:    true,
				LogoImage: testImg,
			},
		},
	}

	tests := []struct {
		name        string
		logoName    string
		width       int
		height      int
		wantFound   bool
		wantWidth   int
		wantHeight  int
	}{
		{
			name:       "Get existing logo with same dimensions",
			logoName:   "sony",
			width:      100,
			height:     50,
			wantFound:  true,
			wantWidth:  100,
			wantHeight: 50,
		},
		{
			name:       "Get existing logo with different dimensions",
			logoName:   "sony",
			width:      200,
			height:     100,
			wantFound:  true,
			wantWidth:  200,
			wantHeight: 100,
		},
		{
			name:       "Get cached scaled logo",
			logoName:   "200_100_sony",
			width:      200,
			height:     100,
			wantFound:  true,
			wantWidth:  200,
			wantHeight: 100,
		},
		{
			name:       "Non-existing logo",
			logoName:   "nikon",
			width:      100,
			height:     50,
			wantFound:  false,
			wantWidth:  0,
			wantHeight: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logo, err := GetLogoImageByNameAndWidhtAndHeight(tt.logoName, tt.width, tt.height)

			if tt.wantFound {
				if pkg.HasError(err) {
					t.Errorf("GetLogoImageByNameAndWidhtAndHeight() unexpected error: %v", err)
				}
				if logo.Width != tt.wantWidth {
					t.Errorf("GetLogoImageByNameAndWidhtAndHeight() Width = %v, want %v", logo.Width, tt.wantWidth)
				}
				if logo.Height != tt.wantHeight {
					t.Errorf("GetLogoImageByNameAndWidhtAndHeight() Height = %v, want %v", logo.Height, tt.wantHeight)
				}
			} else {
				if !pkg.HasError(err) {
					t.Errorf("GetLogoImageByNameAndWidhtAndHeight() expected error but got none")
				}
			}
		})
	}
}

func TestNewLogo(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_logo.png")
	createTestPNG(t, 100, 50, testPath)

	logo, err := newLogo("test_logo", testPath)

	if pkg.HasError(err) {
		t.Errorf("newLogo() unexpected error: %v", err)
	}

	if logo.Name != "test_logo" {
		t.Errorf("newLogo() Name = %v, want test_logo", logo.Name)
	}

	if logo.Width != 100 {
		t.Errorf("newLogo() Width = %v, want 100", logo.Width)
	}

	if logo.Height != 50 {
		t.Errorf("newLogo() Height = %v, want 50", logo.Height)
	}

	if logo.Ext != ".png" {
		t.Errorf("newLogo() Ext = %v, want .png", logo.Ext)
	}

	if !logo.IsLoad {
		t.Errorf("newLogo() IsLoad = false, want true")
	}

	if logo.LogoPath != testPath {
		t.Errorf("newLogo() LogoPath = %v, want %v", logo.LogoPath, testPath)
	}
}

func TestNewLogoWithImage(t *testing.T) {
	testImg := image.NewRGBA(image.Rect(0, 0, 150, 75))

	logo := newLogoWithImage("test_logo", testImg)

	if logo.Name != "test_logo" {
		t.Errorf("newLogoWithImage() Name = %v, want test_logo", logo.Name)
	}

	if logo.Width != 150 {
		t.Errorf("newLogoWithImage() Width = %v, want 150", logo.Width)
	}

	if logo.Height != 75 {
		t.Errorf("newLogoWithImage() Height = %v, want 75", logo.Height)
	}

	if logo.Ext != ".png" {
		t.Errorf("newLogoWithImage() Ext = %v, want .png", logo.Ext)
	}

	if !logo.IsLoad {
		t.Errorf("newLogoWithImage() IsLoad = false, want true")
	}

	if logo.LogoPath != "" {
		t.Errorf("newLogoWithImage() LogoPath = %v, want empty string", logo.LogoPath)
	}

	if logo.LogoImage == nil {
		t.Errorf("newLogoWithImage() LogoImage is nil")
	}
}

func TestLogosImagesInit(t *testing.T) {
	tempDir := setupLogoTest(t)

	logoFiles := []string{"camera1.png", "camera2.png", "sony.png", "canon.png"}

	logos = &logosMaps{
		logoMap: make(map[string]*Logo),
	}

	for _, f := range logoFiles {
		fullPath := filepath.Join(tempDir, f)
		baseName := filepath.Base(fullPath)
		extName := filepath.Ext(fullPath)
		logoName := baseName[:len(baseName)-len(extName)]

		logo, err := newLogo(logoName, fullPath)
		if pkg.HasError(err) {
			t.Fatalf("Failed to load logo %s: %v", logoName, err)
		}
		logos.logoMap[logoName] = logo
	}

	if len(logos.logoMap) != 4 {
		t.Errorf("LogosImagesInit() expected 4 logos, got %d", len(logos.logoMap))
	}

	expectedLogos := []string{"camera1", "camera2", "sony", "canon"}
	for _, expectedLogo := range expectedLogos {
		if _, ok := logos.logoMap[expectedLogo]; !ok {
			t.Errorf("LogosImagesInit() expected logo %s not found", expectedLogo)
		}
	}
}

func TestLogosImagesInitEmptyDir(t *testing.T) {
	tempDir := t.TempDir()

	logos = &logosMaps{
		logoMap: make(map[string]*Logo),
	}

	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read temp dir: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected empty directory, found %d files", len(files))
	}

	if len(logos.logoMap) != 0 {
		t.Errorf("LogosImagesInit() expected 0 logos, got %d", len(logos.logoMap))
	}
}
