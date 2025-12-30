package native

import (
	"image/color"
	"testing"

	"WaterMark/pkg"
)

func TestLoadImageRGBAWithColor(t *testing.T) {
	tests := []struct {
		name    string
		x1      int
		y1      int
		x2      int
		y2      int
		color   color.RGBA
		wantErr bool
	}{
		{
			name:    "Create image with valid coordinates and color",
			x1:      0,
			y1:      0,
			x2:      100,
			y2:      100,
			color:   color.RGBA{255, 255, 255, 255},
			wantErr: false,
		},
		{
			name:    "Create image with black color",
			x1:      0,
			y1:      0,
			x2:      50,
			y2:      50,
			color:   color.RGBA{0, 0, 0, 255},
			wantErr: false,
		},
		{
			name:    "Create image with transparent color",
			x1:      0,
			y1:      0,
			x2:      200,
			y2:      150,
			color:   color.RGBA{255, 255, 255, 0},
			wantErr: false,
		},
		{
			name:    "Create image with zero size",
			x1:      0,
			y1:      0,
			x2:      0,
			y2:      0,
			color:   color.RGBA{255, 255, 255, 255},
			wantErr: false,
		},
		{
			name:    "Create image with negative coordinates",
			x1:      -10,
			y1:      -10,
			x2:      10,
			y2:      10,
			color:   color.RGBA{255, 255, 255, 255},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadImageRGBAWithColor(tt.x1, tt.y1, tt.x2, tt.y2, tt.color)

			if tt.wantErr && !pkg.HasError(err) {
				t.Errorf("loadImageRGBAWithColor() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(err) {
				t.Errorf("loadImageRGBAWithColor() unexpected error: %v", err)
			}

			if !tt.wantErr && got == nil {
				t.Error("loadImageRGBAWithColor() returned nil image")
			}

			if !tt.wantErr && got != nil {
				bounds := got.Bounds()
				expectedWidth := tt.x2 - tt.x1
				expectedHeight := tt.y2 - tt.y1

				if bounds.Dx() != expectedWidth {
					t.Errorf("loadImageRGBAWithColor() width = %v, want %v", bounds.Dx(), expectedWidth)
				}

				if bounds.Dy() != expectedHeight {
					t.Errorf("loadImageRGBAWithColor() height = %v, want %v", bounds.Dy(), expectedHeight)
				}
			}
		})
	}
}

func TestLoadImageRGBA(t *testing.T) {
	tests := []struct {
		name string
		x1   int
		y1   int
		x2   int
		y2   int
	}{
		{
			name: "Create RGBA image with valid coordinates",
			x1:   0,
			y1:   0,
			x2:   100,
			y2:   100,
		},
		{
			name: "Create RGBA image with different dimensions",
			x1:   10,
			y1:   20,
			x2:   110,
			y2:   120,
		},
		{
			name: "Create RGBA image with zero size",
			x1:   0,
			y1:   0,
			x2:   0,
			y2:   0,
		},
		{
			name: "Create RGBA image with large dimensions",
			x1:   0,
			y1:   0,
			x2:   1920,
			y2:   1080,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := loadImageRGBA(tt.x1, tt.y1, tt.x2, tt.y2)

			if got == nil {
				t.Error("loadImageRGBA() returned nil image")
			}

			if got != nil {
				bounds := got.Bounds()
				expectedWidth := tt.x2 - tt.x1
				expectedHeight := tt.y2 - tt.y1

				if bounds.Dx() != expectedWidth {
					t.Errorf("loadImageRGBA() width = %v, want %v", bounds.Dx(), expectedWidth)
				}

				if bounds.Dy() != expectedHeight {
					t.Errorf("loadImageRGBA() height = %v, want %v", bounds.Dy(), expectedHeight)
				}
			}
		})
	}
}

func TestLoadImageRGBAWithColorAndDraw(t *testing.T) {
	tests := []struct {
		name  string
		x1    int
		y1    int
		x2    int
		y2    int
		color color.RGBA
	}{
		{
			name:  "Create image with white background",
			x1:    0,
			y1:    0,
			x2:    100,
			y2:    100,
			color: color.RGBA{255, 255, 255, 255},
		},
		{
			name:  "Create image with black background",
			x1:    0,
			y1:    0,
			x2:    50,
			y2:    50,
			color: color.RGBA{0, 0, 0, 255},
		},
		{
			name:  "Create image with gray background",
			x1:    0,
			y1:    0,
			x2:    200,
			y2:    150,
			color: color.RGBA{128, 128, 128, 255},
		},
		{
			name:  "Create image with red background",
			x1:    0,
			y1:    0,
			x2:    80,
			y2:    60,
			color: color.RGBA{255, 0, 0, 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := loadImageRGBAWithColorAndDraw(tt.x1, tt.y1, tt.x2, tt.y2, tt.color)

			if got == nil {
				t.Error("loadImageRGBAWithColorAndDraw() returned nil image")
			}

			if got != nil {
				bounds := got.Bounds()
				expectedWidth := tt.x2 - tt.x1
				expectedHeight := tt.y2 - tt.y1

				if bounds.Dx() != expectedWidth {
					t.Errorf("loadImageRGBAWithColorAndDraw() width = %v, want %v", bounds.Dx(), expectedWidth)
				}

				if bounds.Dy() != expectedHeight {
					t.Errorf("loadImageRGBAWithColorAndDraw() height = %v, want %v", bounds.Dy(), expectedHeight)
				}

				if expectedWidth > 0 && expectedHeight > 0 {
					pixelColor := got.At(0, 0)
					r, g, b, a := pixelColor.RGBA()
					expectedR := uint32(tt.color.R) * 257
					expectedG := uint32(tt.color.G) * 257
					expectedB := uint32(tt.color.B) * 257
					expectedA := uint32(tt.color.A) * 257

					if r != expectedR {
						t.Errorf("loadImageRGBAWithColorAndDraw() pixel R = %v, want %v", r, expectedR)
					}

					if g != expectedG {
						t.Errorf("loadImageRGBAWithColorAndDraw() pixel G = %v, want %v", g, expectedG)
					}

					if b != expectedB {
						t.Errorf("loadImageRGBAWithColorAndDraw() pixel B = %v, want %v", b, expectedB)
					}

					if a != expectedA {
						t.Errorf("loadImageRGBAWithColorAndDraw() pixel A = %v, want %v", a, expectedA)
					}
				}
			}
		})
	}
}
