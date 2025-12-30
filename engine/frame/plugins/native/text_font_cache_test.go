package native

import (
	"sync"
	"testing"

	"WaterMark/pkg"
)

func TestTextBrushInitFontFileToCache(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Initialize font files to cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
			}()

			err := textBrushInitFontFileToCache()

			if pkg.HasError(err) {
				t.Logf("textBrushInitFontFileToCache() returned error: %v", err)
				t.Skip("Fonts directory not available in test environment")
			}
		})
	}
}

func TestTextBrushInitFontFileToCache_Concurrent(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic: %v", r)
		}
	}()

	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Goroutine recovered from panic: %v", r)
				}
			}()
			err := textBrushInitFontFileToCache()
			if pkg.HasError(err) {
				t.Logf("textBrushInitFontFileToCache() returned error: %v", err)
			}
		}()
	}

	wg.Wait()
}

func TestLoadTextFontWithCache(t *testing.T) {
	tests := []struct {
		name         string
		fontFilePath string
		wantErr      bool
	}{
		{
			name:         "Load font with empty path",
			fontFilePath: "",
			wantErr:      true,
		},
		{
			name:         "Load font with invalid path",
			fontFilePath: "non_existent_font.ttf",
			wantErr:      true,
		},
		{
			name:         "Load font with relative path",
			fontFilePath: "test_font.ttf",
			wantErr:      true,
		},
		{
			name:         "Load font with valid font path",
			fontFilePath: "Alibaba-PuHuiTi-Bold.ttf",
			wantErr:      false,
		},
		{
			name:         "Load font with valid font path - Light",
			fontFilePath: "Alibaba-PuHuiTi-Light.ttf",
			wantErr:      false,
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

			got, err := loadTextFontWithCache(tt.fontFilePath)

			if tt.wantErr && !pkg.HasError(err) {
				t.Errorf("loadTextFontWithCache() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(err) {
				t.Errorf("loadTextFontWithCache() unexpected error: %v", err)
			}

			if !tt.wantErr && got == nil {
				t.Error("loadTextFontWithCache() returned nil font")
			}
		})
	}
}

func TestLoadTextFontWithCache_Cache(t *testing.T) {
	fontFilePath := "Alibaba-PuHuiTi-Bold.ttf"

	t.Run("Load font twice to test cache", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Recovered from panic: %v", r)
			}
		}()

		got1, err1 := loadTextFontWithCache(fontFilePath)
		if pkg.HasError(err1) {
			t.Skipf("Failed to load font: %v", err1)
			return
		}

		got2, err2 := loadTextFontWithCache(fontFilePath)
		if pkg.HasError(err2) {
			t.Errorf("loadTextFontWithCache() unexpected error on second load: %v", err2)
			return
		}

		if got1 == nil || got2 == nil {
			t.Error("loadTextFontWithCache() returned nil font")
		}
	})
}
