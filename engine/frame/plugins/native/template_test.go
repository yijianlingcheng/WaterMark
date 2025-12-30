package native

import (
	"image"
	"image/color"
	"testing"

	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/layout"
	"WaterMark/pkg"
)

func TestNewSourceImage(t *testing.T) {
	tests := []struct {
		name string
		path string
		want *sourceImage
	}{
		{
			name: "Create source image with path",
			path: "/test/path/image.jpg",
			want: &sourceImage{
				path: "/test/path/image.jpg",
			},
		},
		{
			name: "Create source image with empty path",
			path: "",
			want: &sourceImage{
				path: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newSourceImage(tt.path)
			if got == nil {
				t.Fatal("newSourceImage() returned nil")
			}
			if got.path != tt.want.path {
				t.Errorf("newSourceImage().path = %v, want %v", got.path, tt.want.path)
			}
		})
	}
}

func TestNewSeparator(t *testing.T) {
	tests := []struct {
		name   string
		params *layout.FrameLayout
		want   separator
	}{
		{
			name: "Separator with valid dimensions",
			params: &layout.FrameLayout{
				SeparatorWidth:        10,
				SeparatorHeight:       2,
				SeparatorMarginTop:    5,
				SeparatorMarginRight:  5,
				SeparatorMarginBottom: 5,
				SeparatorMarginLeft:   5,
				SeparatorColor:        "255,255,255,255",
			},
			want: separator{
				isExist:      true,
				width:        10,
				height:       2,
				marginTop:    5,
				marginRight:  5,
				marginBottom: 5,
				marginLeft:   5,
				color:        color.RGBA{R: 255, G: 255, B: 255, A: 255},
			},
		},
		{
			name: "Separator with zero width",
			params: &layout.FrameLayout{
				SeparatorWidth:  0,
				SeparatorHeight: 2,
			},
			want: separator{
				isExist: false,
			},
		},
		{
			name: "Separator with zero height",
			params: &layout.FrameLayout{
				SeparatorWidth:  10,
				SeparatorHeight: 0,
			},
			want: separator{
				isExist: false,
			},
		},
		{
			name: "Separator with default color",
			params: &layout.FrameLayout{
				SeparatorWidth:  10,
				SeparatorHeight: 2,
				SeparatorColor:  SEPARATOR_COLOR,
			},
			want: separator{
				isExist: true,
				width:   10,
				height:  2,
				color:   color.RGBA{R: 203, G: 203, B: 201, A: 255},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newSeparator(tt.params)
			if got.isExist != tt.want.isExist {
				t.Errorf("newSeparator().isExist = %v, want %v", got.isExist, tt.want.isExist)
			}
			if got.isExist {
				if got.width != tt.want.width {
					t.Errorf("newSeparator().width = %v, want %v", got.width, tt.want.width)
				}
				if got.height != tt.want.height {
					t.Errorf("newSeparator().height = %v, want %v", got.height, tt.want.height)
				}
				if got.marginTop != tt.want.marginTop {
					t.Errorf("newSeparator().marginTop = %v, want %v", got.marginTop, tt.want.marginTop)
				}
				if got.marginRight != tt.want.marginRight {
					t.Errorf("newSeparator().marginRight = %v, want %v", got.marginRight, tt.want.marginRight)
				}
				if got.marginBottom != tt.want.marginBottom {
					t.Errorf("newSeparator().marginBottom = %v, want %v", got.marginBottom, tt.want.marginBottom)
				}
				if got.marginLeft != tt.want.marginLeft {
					t.Errorf("newSeparator().marginLeft = %v, want %v", got.marginLeft, tt.want.marginLeft)
				}
				if got.color != tt.want.color {
					t.Errorf("newSeparator().color = %v, want %v", got.color, tt.want.color)
				}
			}
		})
	}
}

func TestNewTextLayout(t *testing.T) {
	tests := []struct {
		name   string
		params *layout.FrameLayout
		pos    string
		want   layoutBox
	}{
		{
			name: "Text layout for position one",
			params: &layout.FrameLayout{
				TextOneMarginTop:    10,
				TextOneMarginRight:  10,
				TextOneMarginBottom: 10,
				TextOneMarginLeft:   10,
			},
			pos: textPosOne,
			want: layoutBox{
				marginTop:    10,
				marginRight:  10,
				marginBottom: 10,
				marginLeft:   10,
			},
		},
		{
			name: "Text layout for position two",
			params: &layout.FrameLayout{
				TextTwoMarginTop:    20,
				TextTwoMarginRight:  20,
				TextTwoMarginBottom: 20,
				TextTwoMarginLeft:   20,
			},
			pos: textPosTwo,
			want: layoutBox{
				marginTop:    20,
				marginRight:  20,
				marginBottom: 20,
				marginLeft:   20,
			},
		},
		{
			name: "Text layout for position three",
			params: &layout.FrameLayout{
				TextThreeMarginTop:    30,
				TextThreeMarginRight:  30,
				TextThreeMarginBottom: 30,
				TextThreeMarginLeft:   30,
			},
			pos: textPosThree,
			want: layoutBox{
				marginTop:    30,
				marginRight:  30,
				marginBottom: 30,
				marginLeft:   30,
			},
		},
		{
			name: "Text layout for position four",
			params: &layout.FrameLayout{
				TextFourMarginTop:    40,
				TextFourMarginRight:  40,
				TextFourMarginBottom: 40,
				TextFourMarginLeft:   40,
			},
			pos: textPosFour,
			want: layoutBox{
				marginTop:    40,
				marginRight:  40,
				marginBottom: 40,
				marginLeft:   40,
			},
		},
		{
			name: "Text layout for unknown position (defaults to position one)",
			params: &layout.FrameLayout{
				TextOneMarginTop:    10,
				TextOneMarginRight:  10,
				TextOneMarginBottom: 10,
				TextOneMarginLeft:   10,
			},
			pos: "unknown",
			want: layoutBox{
				marginTop:    10,
				marginRight:  10,
				marginBottom: 10,
				marginLeft:   10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newTextLayout(tt.params, tt.pos)
			if got.marginTop != tt.want.marginTop {
				t.Errorf("newTextLayout().marginTop = %v, want %v", got.marginTop, tt.want.marginTop)
			}
			if got.marginRight != tt.want.marginRight {
				t.Errorf("newTextLayout().marginRight = %v, want %v", got.marginRight, tt.want.marginRight)
			}
			if got.marginBottom != tt.want.marginBottom {
				t.Errorf("newTextLayout().marginBottom = %v, want %v", got.marginBottom, tt.want.marginBottom)
			}
			if got.marginLeft != tt.want.marginLeft {
				t.Errorf("newTextLayout().marginLeft = %v, want %v", got.marginLeft, tt.want.marginLeft)
			}
		})
	}
}

func TestNewLogoLayoutBox(t *testing.T) {
	tests := []struct {
		name   string
		params *layout.FrameLayout
		want   layoutBox
	}{
		{
			name: "Logo layout with all parameters",
			params: &layout.FrameLayout{
				LogoWidth:        100,
				LogoHeight:       50,
				LogoMarginTop:    10,
				LogoMarginRight:  10,
				LogoMarginBottom: 10,
				LogoMarginLeft:   10,
			},
			want: layoutBox{
				width:        100,
				height:       50,
				marginTop:    10,
				marginRight:  10,
				marginBottom: 10,
				marginLeft:   10,
			},
		},
		{
			name: "Logo layout with zero dimensions",
			params: &layout.FrameLayout{
				LogoWidth:  0,
				LogoHeight: 0,
			},
			want: layoutBox{
				width:  0,
				height: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newLogoLayoutBox(tt.params)
			if got.width != tt.want.width {
				t.Errorf("newLogoLayoutBox().width = %v, want %v", got.width, tt.want.width)
			}
			if got.height != tt.want.height {
				t.Errorf("newLogoLayoutBox().height = %v, want %v", got.height, tt.want.height)
			}
			if got.marginTop != tt.want.marginTop {
				t.Errorf("newLogoLayoutBox().marginTop = %v, want %v", got.marginTop, tt.want.marginTop)
			}
			if got.marginRight != tt.want.marginRight {
				t.Errorf("newLogoLayoutBox().marginRight = %v, want %v", got.marginRight, tt.want.marginRight)
			}
			if got.marginBottom != tt.want.marginBottom {
				t.Errorf("newLogoLayoutBox().marginBottom = %v, want %v", got.marginBottom, tt.want.marginBottom)
			}
			if got.marginLeft != tt.want.marginLeft {
				t.Errorf("newLogoLayoutBox().marginLeft = %v, want %v", got.marginLeft, tt.want.marginLeft)
			}
		})
	}
}

func TestNewFinalImage(t *testing.T) {
	tests := []struct {
		name string
		opts *frameOption
		want *finalImage
	}{
		{
			name: "Final image with all parameters",
			opts: &frameOption{
				Params: layout.FrameLayout{
					MainMarginLeft:   50,
					MainMarginRight:  50,
					MainMarginTop:    50,
					MainMarginBottom: 50,
				},
				Exif: exiftool.FileMetadata{
					Fields: map[string]any{
						"ImageWidth":  float64(1920),
						"ImageHeight": float64(1080),
					},
				},
				SaveImageFile: "/test/path/output.jpg",
			},
			want: &finalImage{
				width:  2020,
				height: 1180,
				path:   "/test/path/output.jpg",
			},
		},
		{
			name: "Final image with zero margins",
			opts: &frameOption{
				Params: layout.FrameLayout{
					MainMarginLeft:   0,
					MainMarginRight:  0,
					MainMarginTop:    0,
					MainMarginBottom: 0,
				},
				Exif: exiftool.FileMetadata{
					Fields: map[string]any{
						"ImageWidth":  float64(1920),
						"ImageHeight": float64(1080),
					},
				},
				SaveImageFile: "/test/path/output.jpg",
			},
			want: &finalImage{
				width:  1920,
				height: 1080,
				path:   "/test/path/output.jpg",
			},
		},
		{
			name: "Final image with empty path",
			opts: &frameOption{
				Params: layout.FrameLayout{
					MainMarginLeft:   50,
					MainMarginRight:  50,
					MainMarginTop:    50,
					MainMarginBottom: 50,
				},
				Exif: exiftool.FileMetadata{
					Fields: map[string]any{
						"ImageWidth":  float64(1920),
						"ImageHeight": float64(1080),
					},
				},
				SaveImageFile: "",
			},
			want: &finalImage{
				width:  2020,
				height: 1180,
				path:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newFinalImage(tt.opts)
			if got == nil {
				t.Fatal("newFinalImage() returned nil")
			}
			if got.width != tt.want.width {
				t.Errorf("newFinalImage().width = %v, want %v", got.width, tt.want.width)
			}
			if got.height != tt.want.height {
				t.Errorf("newFinalImage().height = %v, want %v", got.height, tt.want.height)
			}
			if got.path != tt.want.path {
				t.Errorf("newFinalImage().path = %v, want %v", got.path, tt.want.path)
			}
		})
	}
}

func TestSourceImageSetImage(t *testing.T) {
	tests := []struct {
		name string
		path string
		img  image.Image
	}{
		{
			name: "Set image with valid image",
			path: "/test/path/image.jpg",
			img:  image.NewRGBA(image.Rect(0, 0, 100, 100)),
		},
		{
			name: "Set image with different dimensions",
			path: "/test/path/image2.jpg",
			img:  image.NewRGBA(image.Rect(0, 0, 1920, 1080)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := newSourceImage(tt.path)
			src.SetImage(tt.img)

			if src.imgDecode == nil {
				t.Error("SetImage() did not set imgDecode")
			}
			if src.width != tt.img.Bounds().Dx() {
				t.Errorf("SetImage() width = %v, want %v", src.width, tt.img.Bounds().Dx())
			}
			if src.height != tt.img.Bounds().Dy() {
				t.Errorf("SetImage() height = %v, want %v", src.height, tt.img.Bounds().Dy())
			}
		})
	}
}

func TestSourceImageSetImageXAndY(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		width  int
		height int
	}{
		{
			name:   "Set dimensions with valid values",
			path:   "/test/path/image.jpg",
			width:  1920,
			height: 1080,
		},
		{
			name:   "Set dimensions with zero values",
			path:   "/test/path/image2.jpg",
			width:  0,
			height: 0,
		},
		{
			name:   "Set dimensions with negative values",
			path:   "/test/path/image3.jpg",
			width:  -100,
			height: -100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := newSourceImage(tt.path)
			src.SetImageXAndY(tt.width, tt.height)

			if src.width != tt.width {
				t.Errorf("SetImageXAndY() width = %v, want %v", src.width, tt.width)
			}
			if src.height != tt.height {
				t.Errorf("SetImageXAndY() height = %v, want %v", src.height, tt.height)
			}
		})
	}
}

func TestNewBorderImage(t *testing.T) {
	tests := []struct {
		name    string
		params  *layout.FrameLayout
		wantErr bool
	}{
		{
			name: "Border image with valid parameters",
			params: &layout.FrameLayout{
				BgColor:            "255,255,255,255",
				MainMarginLeft:     50,
				MainMarginRight:    50,
				MainMarginTop:      50,
				MainMarginBottom:   50,
				TextOneFontFile:    "",
				TextOneFontSize:    12,
				TextOneFontColor:   "255,255,255,255",
				TextTwoFontFile:    "",
				TextTwoFontSize:    12,
				TextTwoFontColor:   "255,255,255,255",
				TextThreeFontFile:  "",
				TextThreeFontSize:  12,
				TextThreeFontColor: "255,255,255,255",
				TextFourFontFile:   "",
				TextFourFontSize:   12,
				TextFourFontColor:  "255,255,255,255",
				SeparatorWidth:     10,
				SeparatorHeight:    2,
				SeparatorColor:     "203,203,201,255",
				LogoWidth:          100,
				LogoHeight:         50,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newBorderImage(tt.params)
			if pkg.HasError(err) != tt.wantErr {
				t.Errorf("newBorderImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.bgColor.R != 255 || got.bgColor.G != 255 || got.bgColor.B != 255 || got.bgColor.A != 255 {
					t.Errorf("newBorderImage().bgColor = %v, want {255,255,255,255}", got.bgColor)
				}
				if got.leftWidth != 50 {
					t.Errorf("newBorderImage().leftWidth = %v, want 50", got.leftWidth)
				}
				if got.rightWidth != 50 {
					t.Errorf("newBorderImage().rightWidth = %v, want 50", got.rightWidth)
				}
				if got.topHeight != 50 {
					t.Errorf("newBorderImage().topHeight = %v, want 50", got.topHeight)
				}
				if got.bottomHeight != 50 {
					t.Errorf("newBorderImage().bottomHeight = %v, want 50", got.bottomHeight)
				}
			}
		})
	}
}
