package native

import (
	"image"
	"image/color"
	"image/draw"
	"sync"
	"testing"

	"WaterMark/internal"
	"WaterMark/pkg"
)

func TestMain(m *testing.M) {
	internal.InitAppConfigsAndRes()
	InitAllCachaAndTools()
	m.Run()
}

func TestNewTextBrush(t *testing.T) {
	tests := []struct {
		name         string
		fontFilePath string
		fontSize     float64
		fontColor    *image.Uniform
		wantErr      bool
	}{
		{
			name:         "Create text brush with empty font path",
			fontFilePath: "",
			fontSize:     12,
			fontColor:    &image.Uniform{C: color.RGBA{255, 255, 255, 255}},
			wantErr:      false,
		},
		{
			name:         "Create text brush with invalid font path",
			fontFilePath: "test_font.ttf",
			fontSize:     14,
			fontColor:    &image.Uniform{C: color.RGBA{0, 0, 0, 255}},
			wantErr:      true,
		},
		{
			name:         "Create text brush with valid font path",
			fontFilePath: "Alibaba-PuHuiTi-Bold.ttf",
			fontSize:     14,
			fontColor:    &image.Uniform{C: color.RGBA{0, 0, 0, 255}},
			wantErr:      false,
		},
		{
			name:         "Create text brush with valid font path and large size",
			fontFilePath: "Alibaba-PuHuiTi-Bold.ttf",
			fontSize:     48,
			fontColor:    &image.Uniform{C: color.RGBA{255, 0, 0, 255}},
			wantErr:      false,
		},
		{
			name:         "Create text brush with valid font path and nil color",
			fontFilePath: "Alibaba-PuHuiTi-Bold.ttf",
			fontSize:     24,
			fontColor:    nil,
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

			got, err := newTextBrush(tt.fontFilePath, tt.fontSize, tt.fontColor)

			if tt.wantErr && !pkg.HasError(err) {
				t.Errorf("newTextBrush() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(err) {
				t.Errorf("newTextBrush() unexpected error: %v", err)
			}

			if !tt.wantErr && got == nil {
				t.Error("newTextBrush() returned nil brush")
			}

			if !tt.wantErr && got != nil {
				if tt.fontFilePath == "" {
					if got.FontSize != 0 {
						t.Errorf("newTextBrush() FontSize = %v, want 0 (empty path returns empty brush)", got.FontSize)
					}
					if got.FontColor != nil {
						t.Error("newTextBrush() FontColor should be nil (empty path returns empty brush)")
					}
				} else {
					if got.FontSize != tt.fontSize {
						t.Errorf("newTextBrush() FontSize = %v, want %v", got.FontSize, tt.fontSize)
					}
					if got.FontColor != tt.fontColor {
						t.Errorf("newTextBrush() FontColor mismatch")
					}
				}
			}
		})
	}
}

func TestTextBrush_DrawFontOnRGBA_EdgeCases(t *testing.T) {
	fontFilePath := "Alibaba-PuHuiTi-Bold.ttf"

	tests := []struct {
		name      string
		fontSize  float64
		fontColor *image.Uniform
		rgba      draw.Image
		pt        image.Point
		content   string
		wantErr   bool
	}{
		{
			name:      "Draw text with very large font size",
			fontSize:  200,
			fontColor: &image.Uniform{C: color.RGBA{0, 0, 0, 255}},
			rgba:      image.NewRGBA(image.Rect(0, 0, 1000, 1000)),
			pt:        image.Point{X: 50, Y: 500},
			content:   "Large",
			wantErr:   false,
		},
		{
			name:      "Draw text with very small font size",
			fontSize:  1,
			fontColor: &image.Uniform{C: color.RGBA{0, 0, 0, 255}},
			rgba:      image.NewRGBA(image.Rect(0, 0, 100, 100)),
			pt:        image.Point{X: 10, Y: 50},
			content:   "Tiny",
			wantErr:   false,
		},
		{
			name:      "Draw text with negative coordinates",
			fontSize:  24,
			fontColor: &image.Uniform{C: color.RGBA{0, 0, 0, 255}},
			rgba:      image.NewRGBA(image.Rect(0, 0, 200, 100)),
			pt:        image.Point{X: -100, Y: -50},
			content:   "Negative",
			wantErr:   false,
		},
		{
			name:      "Draw text at edge of image",
			fontSize:  24,
			fontColor: &image.Uniform{C: color.RGBA{0, 0, 0, 255}},
			rgba:      image.NewRGBA(image.Rect(0, 0, 200, 100)),
			pt:        image.Point{X: 190, Y: 90},
			content:   "Edge",
			wantErr:   false,
		},
		{
			name:      "Draw text with zero font size",
			fontSize:  0,
			fontColor: &image.Uniform{C: color.RGBA{0, 0, 0, 255}},
			rgba:      image.NewRGBA(image.Rect(0, 0, 200, 100)),
			pt:        image.Point{X: 10, Y: 50},
			content:   "Zero",
			wantErr:   false,
		},
		{
			name:      "Draw text with transparent color",
			fontSize:  24,
			fontColor: &image.Uniform{C: color.RGBA{255, 255, 255, 0}},
			rgba:      image.NewRGBA(image.Rect(0, 0, 200, 100)),
			pt:        image.Point{X: 10, Y: 50},
			content:   "Transparent",
			wantErr:   false,
		},
		{
			name:      "Draw text with semi-transparent color",
			fontSize:  24,
			fontColor: &image.Uniform{C: color.RGBA{255, 0, 0, 128}},
			rgba:      image.NewRGBA(image.Rect(0, 0, 200, 100)),
			pt:        image.Point{X: 10, Y: 50},
			content:   "Semi",
			wantErr:   false,
		},
		{
			name:      "Draw text with newline characters",
			fontSize:  24,
			fontColor: &image.Uniform{C: color.RGBA{0, 0, 0, 255}},
			rgba:      image.NewRGBA(image.Rect(0, 0, 200, 100)),
			pt:        image.Point{X: 10, Y: 50},
			content:   "Line1\nLine2",
			wantErr:   false,
		},
		{
			name:      "Draw text with tab characters",
			fontSize:  24,
			fontColor: &image.Uniform{C: color.RGBA{0, 0, 0, 255}},
			rgba:      image.NewRGBA(image.Rect(0, 0, 200, 100)),
			pt:        image.Point{X: 10, Y: 50},
			content:   "Tab\tText",
			wantErr:   false,
		},
		{
			name:      "Draw text with very long string",
			fontSize:  12,
			fontColor: &image.Uniform{C: color.RGBA{0, 0, 0, 255}},
			rgba:      image.NewRGBA(image.Rect(0, 0, 1000, 100)),
			pt:        image.Point{X: 10, Y: 50},
			content:   "This is a very long text that should test the drawing function with a large amount of text content to ensure it handles long strings properly without any issues",
			wantErr:   false,
		},
		{
			name:      "Draw text with emoji characters",
			fontSize:  24,
			fontColor: &image.Uniform{C: color.RGBA{0, 0, 0, 255}},
			rgba:      image.NewRGBA(image.Rect(0, 0, 200, 100)),
			pt:        image.Point{X: 10, Y: 50},
			content:   "Emoji 😊",
			wantErr:   false,
		},
		{
			name:      "Draw text with numbers and symbols",
			fontSize:  24,
			fontColor: &image.Uniform{C: color.RGBA{0, 0, 0, 255}},
			rgba:      image.NewRGBA(image.Rect(0, 0, 200, 100)),
			pt:        image.Point{X: 10, Y: 50},
			content:   "123!@#$%",
			wantErr:   false,
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

			brush, err := newTextBrush(fontFilePath, tt.fontSize, tt.fontColor)
			if pkg.HasError(err) {
				t.Skipf("Failed to load font: %v", err)
				return
			}

			drawErr := brush.drawFontOnRGBA(tt.rgba, tt.pt, tt.content)

			if tt.wantErr && !pkg.HasError(drawErr) {
				t.Errorf("textBrush.drawFontOnRGBA() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(drawErr) {
				t.Errorf("textBrush.drawFontOnRGBA() unexpected error: %v", drawErr)
			}
		})
	}
}

func TestTextBrush_DrawFontOnRGBA_Concurrent(t *testing.T) {
	fontFilePath := "Alibaba-PuHuiTi-Bold.ttf"
	fontColor := &image.Uniform{C: color.RGBA{0, 0, 0, 255}}

	brush, err := newTextBrush(fontFilePath, 24, fontColor)
	if pkg.HasError(err) {
		t.Skipf("Failed to load font: %v", err)
		return
	}

	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Goroutine %d recovered from panic: %v", index, r)
				}
			}()

			rgba := image.NewRGBA(image.Rect(0, 0, 200, 100))
			pt := image.Point{X: 10, Y: 50}
			content := "Test"

			drawErr := brush.drawFontOnRGBA(rgba, pt, content)
			if pkg.HasError(drawErr) {
				t.Logf("Goroutine %d drawFontOnRGBA() returned error: %v", index, drawErr)
			}
		}(i)
	}

	wg.Wait()
}

func TestTextBrush_DrawFontOnRGBA(t *testing.T) {
	tests := []struct {
		name    string
		brush   *textBrush
		rgba    draw.Image
		pt      image.Point
		content string
		wantErr bool
	}{
		{
			name: "Draw text with nil font type",
			brush: &textBrush{
				FontType:  nil,
				FontColor: &image.Uniform{C: color.RGBA{255, 255, 255, 255}},
				FontSize:  12,
			},
			rgba:    image.NewRGBA(image.Rect(0, 0, 100, 100)),
			pt:      image.Point{X: 10, Y: 10},
			content: "Test",
			wantErr: false,
		},
		{
			name: "Draw text with empty content",
			brush: &textBrush{
				FontType:  nil,
				FontColor: &image.Uniform{C: color.RGBA{255, 255, 255, 255}},
				FontSize:  12,
			},
			rgba:    image.NewRGBA(image.Rect(0, 0, 100, 100)),
			pt:      image.Point{X: 10, Y: 10},
			content: "",
			wantErr: false,
		},
		{
			name: "Draw text with nil rgba",
			brush: &textBrush{
				FontType:  nil,
				FontColor: &image.Uniform{C: color.RGBA{255, 255, 255, 255}},
				FontSize:  12,
			},
			rgba:    nil,
			pt:      image.Point{X: 10, Y: 10},
			content: "Test",
			wantErr: false,
		},
		{
			name: "Draw text with zero point",
			brush: &textBrush{
				FontType:  nil,
				FontColor: &image.Uniform{C: color.RGBA{255, 255, 255, 255}},
				FontSize:  12,
			},
			rgba:    image.NewRGBA(image.Rect(0, 0, 100, 100)),
			pt:      image.Point{X: 0, Y: 0},
			content: "Test",
			wantErr: false,
		},
		{
			name: "Draw text with negative point",
			brush: &textBrush{
				FontType:  nil,
				FontColor: &image.Uniform{C: color.RGBA{255, 255, 255, 255}},
				FontSize:  12,
			},
			rgba:    image.NewRGBA(image.Rect(0, 0, 100, 100)),
			pt:      image.Point{X: -10, Y: -10},
			content: "Test",
			wantErr: false,
		},
		{
			name: "Draw text with large font size",
			brush: &textBrush{
				FontType:  nil,
				FontColor: &image.Uniform{C: color.RGBA{255, 255, 255, 255}},
				FontSize:  100,
			},
			rgba:    image.NewRGBA(image.Rect(0, 0, 100, 100)),
			pt:      image.Point{X: 10, Y: 10},
			content: "Test",
			wantErr: false,
		},
		{
			name: "Draw text with small font size",
			brush: &textBrush{
				FontType:  nil,
				FontColor: &image.Uniform{C: color.RGBA{255, 255, 255, 255}},
				FontSize:  1,
			},
			rgba:    image.NewRGBA(image.Rect(0, 0, 100, 100)),
			pt:      image.Point{X: 10, Y: 10},
			content: "Test",
			wantErr: false,
		},
		{
			name: "Draw text with zero font size",
			brush: &textBrush{
				FontType:  nil,
				FontColor: &image.Uniform{C: color.RGBA{255, 255, 255, 255}},
				FontSize:  0,
			},
			rgba:    image.NewRGBA(image.Rect(0, 0, 100, 100)),
			pt:      image.Point{X: 10, Y: 10},
			content: "Test",
			wantErr: false,
		},
		{
			name: "Draw text with nil font color",
			brush: &textBrush{
				FontType:  nil,
				FontColor: nil,
				FontSize:  12,
			},
			rgba:    image.NewRGBA(image.Rect(0, 0, 100, 100)),
			pt:      image.Point{X: 10, Y: 10},
			content: "Test",
			wantErr: false,
		},
		{
			name: "Draw text with special characters",
			brush: &textBrush{
				FontType:  nil,
				FontColor: &image.Uniform{C: color.RGBA{255, 255, 255, 255}},
				FontSize:  12,
			},
			rgba:    image.NewRGBA(image.Rect(0, 0, 100, 100)),
			pt:      image.Point{X: 10, Y: 10},
			content: "Test @#$%^&*()",
			wantErr: false,
		},
		{
			name: "Draw text with unicode characters",
			brush: &textBrush{
				FontType:  nil,
				FontColor: &image.Uniform{C: color.RGBA{255, 255, 255, 255}},
				FontSize:  12,
			},
			rgba:    image.NewRGBA(image.Rect(0, 0, 100, 100)),
			pt:      image.Point{X: 10, Y: 10},
			content: "测试中文",
			wantErr: false,
		},
		{
			name: "Draw text with very long content",
			brush: &textBrush{
				FontType:  nil,
				FontColor: &image.Uniform{C: color.RGBA{255, 255, 255, 255}},
				FontSize:  12,
			},
			rgba:    image.NewRGBA(image.Rect(0, 0, 100, 100)),
			pt:      image.Point{X: 10, Y: 10},
			content: "This is a very long text content that should test the drawing function with a large amount of text",
			wantErr: false,
		},
		{
			name: "Draw text with different color",
			brush: &textBrush{
				FontType:  nil,
				FontColor: &image.Uniform{C: color.RGBA{255, 0, 0, 255}},
				FontSize:  12,
			},
			rgba:    image.NewRGBA(image.Rect(0, 0, 100, 100)),
			pt:      image.Point{X: 10, Y: 10},
			content: "Test",
			wantErr: false,
		},
		{
			name: "Draw text with transparent color",
			brush: &textBrush{
				FontType:  nil,
				FontColor: &image.Uniform{C: color.RGBA{255, 255, 255, 0}},
				FontSize:  12,
			},
			rgba:    image.NewRGBA(image.Rect(0, 0, 100, 100)),
			pt:      image.Point{X: 10, Y: 10},
			content: "Test",
			wantErr: false,
		},
		{
			name: "Draw text with large rgba bounds",
			brush: &textBrush{
				FontType:  nil,
				FontColor: &image.Uniform{C: color.RGBA{255, 255, 255, 255}},
				FontSize:  12,
			},
			rgba:    image.NewRGBA(image.Rect(0, 0, 10000, 10000)),
			pt:      image.Point{X: 100, Y: 100},
			content: "Test",
			wantErr: false,
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

			err := tt.brush.drawFontOnRGBA(tt.rgba, tt.pt, tt.content)

			if tt.wantErr && !pkg.HasError(err) {
				t.Errorf("textBrush.drawFontOnRGBA() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(err) {
				t.Errorf("textBrush.drawFontOnRGBA() unexpected error: %v", err)
			}
		})
	}
}

func TestTextBrush_DrawFontOnRGBA_WithRealFont(t *testing.T) {
	fontFilePath := "Alibaba-PuHuiTi-Bold.ttf"
	fontColor := &image.Uniform{C: color.RGBA{0, 0, 0, 255}}

	tests := []struct {
		name     string
		fontSize float64
		rgba     draw.Image
		pt       image.Point
		content  string
		wantErr  bool
	}{
		{
			name:     "Draw text with real font - normal size",
			fontSize: 24,
			rgba:     image.NewRGBA(image.Rect(0, 0, 200, 100)),
			pt:       image.Point{X: 10, Y: 50},
			content:  "Hello",
			wantErr:  false,
		},
		{
			name:     "Draw text with real font - large size",
			fontSize: 48,
			rgba:     image.NewRGBA(image.Rect(0, 0, 400, 200)),
			pt:       image.Point{X: 20, Y: 100},
			content:  "Big Text",
			wantErr:  false,
		},
		{
			name:     "Draw text with real font - Chinese characters",
			fontSize: 32,
			rgba:     image.NewRGBA(image.Rect(0, 0, 300, 150)),
			pt:       image.Point{X: 10, Y: 80},
			content:  "中文测试",
			wantErr:  false,
		},
		{
			name:     "Draw text with real font - mixed content",
			fontSize: 20,
			rgba:     image.NewRGBA(image.Rect(0, 0, 300, 150)),
			pt:       image.Point{X: 10, Y: 80},
			content:  "Hello 世界",
			wantErr:  false,
		},
		{
			name:     "Draw text with real font - special characters",
			fontSize: 18,
			rgba:     image.NewRGBA(image.Rect(0, 0, 200, 100)),
			pt:       image.Point{X: 10, Y: 50},
			content:  "Test@#$%",
			wantErr:  false,
		},
		{
			name:     "Draw text with real font - small size",
			fontSize: 10,
			rgba:     image.NewRGBA(image.Rect(0, 0, 100, 50)),
			pt:       image.Point{X: 5, Y: 25},
			content:  "Small",
			wantErr:  false,
		},
		{
			name:     "Draw text with real font - empty content",
			fontSize: 24,
			rgba:     image.NewRGBA(image.Rect(0, 0, 200, 100)),
			pt:       image.Point{X: 10, Y: 50},
			content:  "",
			wantErr:  false,
		},
		{
			name:     "Draw text with real font - white color",
			fontSize: 24,
			rgba:     image.NewRGBA(image.Rect(0, 0, 200, 100)),
			pt:       image.Point{X: 10, Y: 50},
			content:  "White",
			wantErr:  false,
		},
		{
			name:     "Draw text with real font - red color",
			fontSize: 24,
			rgba:     image.NewRGBA(image.Rect(0, 0, 200, 100)),
			pt:       image.Point{X: 10, Y: 50},
			content:  "Red",
			wantErr:  false,
		},
		{
			name:     "Draw text with real font - blue color",
			fontSize: 24,
			rgba:     image.NewRGBA(image.Rect(0, 0, 200, 100)),
			pt:       image.Point{X: 10, Y: 50},
			content:  "Blue",
			wantErr:  false,
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

			brush, err := newTextBrush(fontFilePath, tt.fontSize, fontColor)
			if pkg.HasError(err) {
				t.Skipf("Failed to load font: %v", err)
				return
			}

			drawErr := brush.drawFontOnRGBA(tt.rgba, tt.pt, tt.content)

			if tt.wantErr && !pkg.HasError(drawErr) {
				t.Errorf("textBrush.drawFontOnRGBA() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(drawErr) {
				t.Errorf("textBrush.drawFontOnRGBA() unexpected error: %v", drawErr)
			}
		})
	}
}
