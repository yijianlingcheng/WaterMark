package native

import (
	"testing"

	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/layout"
)

func TestFrameOption_NeedSourceImage(t *testing.T) {
	tests := []struct {
		name      string
		photoType string
		want      bool
	}{
		{
			name:      "Need source image for normal photo",
			photoType: "normal",
			want:      true,
		},
		{
			name:      "Need source image for empty type",
			photoType: "",
			want:      true,
		},
		{
			name:      "Don't need source image for border type",
			photoType: PHOTO_TYPE_BORDER,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := &frameOption{
				PhotoType: tt.photoType,
			}

			got := fp.needSourceImage()

			if got != tt.want {
				t.Errorf("frameOption.needSourceImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFrameOption_GetExif(t *testing.T) {
	exif := exiftool.FileMetadata{
		Fields: map[string]any{
			"Make": "TestCamera",
		},
	}

	fp := &frameOption{
		Exif: exif,
	}

	got := fp.getExif()

	if got.Fields["Make"] != exif.Fields["Make"] {
		t.Errorf("frameOption.getExif() Make = %v, want %v", got.Fields["Make"], exif.Fields["Make"])
	}
}

func TestFrameOption_GetSourceImageFile(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "Get source image file path",
			path: "test.jpg",
			want: "test.jpg",
		},
		{
			name: "Get empty source image file path",
			path: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := &frameOption{
				SourceImageFile: tt.path,
			}

			got := fp.getSourceImageFile()

			if got != tt.want {
				t.Errorf("frameOption.getSourceImageFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFrameOption_GetSourceImageX(t *testing.T) {
	tests := []struct {
		name string
		exif exiftool.FileMetadata
		want int
	}{
		{
			name: "Get source image width",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					"ImageWidth": float64(1920),
				},
			},
			want: 1920,
		},
		{
			name: "Get source image width with zero",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					"ImageWidth": float64(0),
				},
			},
			want: 0,
		},
		{
			name: "Get source image width without field",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{},
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := &frameOption{
				Exif: tt.exif,
			}

			got := fp.getSourceImageX()

			if got != tt.want {
				t.Errorf("frameOption.getSourceImageX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFrameOption_GetSourceImageY(t *testing.T) {
	tests := []struct {
		name string
		exif exiftool.FileMetadata
		want int
	}{
		{
			name: "Get source image height",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					"ImageHeight": float64(1080),
				},
			},
			want: 1080,
		},
		{
			name: "Get source image height with zero",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					"ImageHeight": float64(0),
				},
			},
			want: 0,
		},
		{
			name: "Get source image height without field",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{},
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := &frameOption{
				Exif: tt.exif,
			}

			got := fp.getSourceImageY()

			if got != tt.want {
				t.Errorf("frameOption.getSourceImageY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFrameOption_GetMakeFromExif(t *testing.T) {
	tests := []struct {
		name string
		exif exiftool.FileMetadata
		want string
	}{
		{
			name: "Get make from exif",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					"Make": "NIKON CORPORATION",
				},
			},
			want: "NIKON CORPORATION",
		},
		{
			name: "Get empty make from exif",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := &frameOption{
				Exif: tt.exif,
			}

			got := fp.getMakeFromExif()

			if got != tt.want {
				t.Errorf("frameOption.getMakeFromExif() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFrameOption(t *testing.T) {
	tests := []struct {
		name    string
		opts    map[string]any
		wantErr bool
	}{
		{
			name: "Create frame option with valid data",
			opts: map[string]any{
				"photoType":       "normal",
				"sourceImageFile": "test.jpg",
				"saveImageFile":   "output.jpg",
				"params": layout.FrameLayout{
					Name: "test_layout",
				},
			},
			wantErr: false,
		},
		{
			name:    "Create frame option with empty map",
			opts:    map[string]any{},
			wantErr: false,
		},
		{
			name:    "Create frame option with nil map",
			opts:    nil,
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

			fp := newFrameOption(tt.opts)

			if tt.wantErr {
				return
			}

			if fp == nil {
				t.Error("newFrameOption() returned nil")
			}
		})
	}
}

func TestFrameOption_ResetSourceImageX(t *testing.T) {
	tests := []struct {
		name         string
		exif         exiftool.FileMetadata
		newWidth     int
		wantOrigin   int
		wantNewWidth int
	}{
		{
			name: "Reset source image width",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					"ImageWidth": float64(1920),
				},
			},
			newWidth:     2048,
			wantOrigin:   1920,
			wantNewWidth: 2048,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := &frameOption{
				Exif: tt.exif,
			}

			fp.resetSourceImageX(tt.newWidth)

			if fp.OriginWidth != tt.wantOrigin {
				t.Errorf("frameOption.OriginWidth = %v, want %v", fp.OriginWidth, tt.wantOrigin)
			}

			newWidth, ok := fp.Exif.Fields["ImageWidth"].(float64)
			if !ok {
				t.Error("ImageWidth is not float64")
			}

			if int(newWidth) != tt.wantNewWidth {
				t.Errorf("frameOption.Exif.Fields[ImageWidth] = %v, want %v", int(newWidth), tt.wantNewWidth)
			}
		})
	}
}

func TestFrameOption_ResetSourceImageY(t *testing.T) {
	tests := []struct {
		name          string
		exif          exiftool.FileMetadata
		newHeight     int
		wantOrigin    int
		wantNewHeight int
	}{
		{
			name: "Reset source image height",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					"ImageHeight": float64(1080),
				},
			},
			newHeight:     2048,
			wantOrigin:    1080,
			wantNewHeight: 2048,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := &frameOption{
				Exif: tt.exif,
			}

			fp.resetSourceImageY(tt.newHeight)

			if fp.OriginHeight != tt.wantOrigin {
				t.Errorf("frameOption.OriginHeight = %v, want %v", fp.OriginHeight, tt.wantOrigin)
			}

			newHeight, ok := fp.Exif.Fields["ImageHeight"].(float64)
			if !ok {
				t.Error("ImageHeight is not float64")
			}

			if int(newHeight) != tt.wantNewHeight {
				t.Errorf("frameOption.Exif.Fields[ImageHeight] = %v, want %v", int(newHeight), tt.wantNewHeight)
			}
		})
	}
}

func TestFrameOption_IsVerticalImage(t *testing.T) {
	tests := []struct {
		name string
		exif exiftool.FileMetadata
		want bool
	}{
		{
			name: "Vertical image (portrait)",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					"ImageWidth":  float64(1080),
					"ImageHeight": float64(1920),
				},
			},
			want: true,
		},
		{
			name: "Horizontal image (landscape)",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					"ImageWidth":  float64(1920),
					"ImageHeight": float64(1080),
				},
			},
			want: false,
		},
		{
			name: "Square image",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					"ImageWidth":  float64(1080),
					"ImageHeight": float64(1080),
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := &frameOption{
				Exif: tt.exif,
			}

			got := fp.isVerticalImage()

			if got != tt.want {
				t.Errorf("frameOption.isVerticalImage() = %v, want %v", got, tt.want)
			}
		})
	}
}
