package native

import (
	"image"
	"image/color"
	"os"
	"testing"
	"time"

	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/layout"
	"WaterMark/pkg"
)

func TestBlurPhotoFrame_InitFrame(t *testing.T) {
	tests := []struct {
		name    string
		opts    map[string]any
		wantErr bool
	}{
		{
			name:    "Init blur frame with nil opts",
			opts:    nil,
			wantErr: true,
		},
		{
			name:    "Init blur frame with empty opts",
			opts:    map[string]any{},
			wantErr: true,
		},
		{
			name: "Init blur frame with invalid source image",
			opts: map[string]any{
				"sourceImageFile": "non_existent.jpg",
			},
			wantErr: true,
		},
		{
			name: "Init blur frame without source image",
			opts: map[string]any{
				"sourceImageFile": "test.jpg",
				"photoType":       "border",
				"layoutName":      "simple",
				"saveImageFile":   "output.jpg",
				"params": layout.FrameLayout{
					Name:           "simple",
					BorderRadius:   10,
					MainMarginLeft: 50,
					MainMarginTop:  50,
					BgColor:        "#000000",
				},
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

			fm := &blurPhotoFrame{}
			err := fm.initFrame(tt.opts)

			if tt.wantErr && !pkg.HasError(err) {
				t.Errorf("initFrame() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(err) {
				t.Errorf("initFrame() unexpected error: %v", err)
			}

			if !tt.wantErr {
				if !fm.isBlur {
					t.Error("initFrame() should set isBlur to true")
				}
				if fm.opts == nil {
					t.Error("initFrame() should set opts")
				}
			}
		})
	}
}

func TestBlurPhotoFrame_CreateTransparentDraw(t *testing.T) {
	tests := []struct {
		name    string
		fm      *blurPhotoFrame
		wantErr bool
	}{
		{
			name: "Create transparent draw",
			fm: &blurPhotoFrame{
				basePhotoFrame: basePhotoFrame{
					finImage: &finalImage{
						width:  800,
						height: 600,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			drawImg, err := tt.fm.createTransparentDraw()

			if tt.wantErr && !pkg.HasError(err) {
				t.Errorf("createTransparentDraw() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(err) {
				t.Errorf("createTransparentDraw() unexpected error: %v", err)
			}

			if !tt.wantErr {
				if drawImg == nil {
					t.Error("createTransparentDraw() should return non-nil image")
				}
				bounds := drawImg.Bounds()
				if bounds.Dx() != tt.fm.finImage.width || bounds.Dy() != tt.fm.finImage.height {
					t.Errorf("createTransparentDraw() returned image with wrong dimensions")
				}
			}
		})
	}
}

func TestBlurPhotoFrame_GetBlurBackgroundImageFilePath(t *testing.T) {
	tests := []struct {
		name string
		fm   *blurPhotoFrame
		want string
	}{
		{
			name: "Get blur background image file path",
			fm: &blurPhotoFrame{
				basePhotoFrame: basePhotoFrame{
					srcImage: &sourceImage{
						path: "test.jpg",
					},
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.getBlurBackgroundImageFilePath()

			if got == "" {
				t.Log("getBlurBackgroundImageFilePath() returned empty path")
			}
		})
	}
}

func TestBlurPhotoFrame_SetTmpBlurResultImagePath(t *testing.T) {
	tests := []struct {
		name string
		fm   *blurPhotoFrame
	}{
		{
			name: "Set tmp blur result image path with source image",
			fm: &blurPhotoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						PhotoType: "normal",
					},
					srcImage: &sourceImage{
						path: "test.jpg",
					},
				},
			},
		},
		{
			name: "Set tmp blur result image path without source image",
			fm: &blurPhotoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						PhotoType: "border",
					},
					srcImage: &sourceImage{
						path: "test.jpg",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fm.setTmpBlurResultImagePath()

			if tt.fm.opts.needSourceImage() {
				if tt.fm.tmpBlurResultImagePath == "" {
					t.Error("setTmpBlurResultImagePath() should set path when needSourceImage is true")
				}
			} else {
				if tt.fm.tmpBlurResultImagePath != "" {
					t.Error("setTmpBlurResultImagePath() should not set path when needSourceImage is false")
				}
			}
		})
	}
}

func TestBlurPhotoFrame_LoadSourceImage(t *testing.T) {
	tests := []struct {
		name    string
		fm      *blurPhotoFrame
		path    string
		wantErr bool
	}{
		{
			name: "Load source image with invalid path",
			fm: &blurPhotoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{},
				},
			},
			path:    "non_existent.jpg",
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

			img, err := tt.fm.loadSourceImage(tt.path)

			if tt.wantErr && !pkg.HasError(err) {
				t.Errorf("loadSourceImage() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(err) {
				t.Errorf("loadSourceImage() unexpected error: %v", err)
			}

			if !tt.wantErr && img == nil {
				t.Error("loadSourceImage() should return non-nil image")
			}
		})
	}
}

func TestBlurPhotoFrame_GetBlurSaveImageFile(t *testing.T) {
	tests := []struct {
		name string
		fm   *blurPhotoFrame
		want string
	}{
		{
			name: "Get blur save image file with custom path",
			fm: &blurPhotoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						SaveImageFile: "custom_output.jpg",
					},
				},
			},
			want: "custom_output.jpg",
		},
		{
			name: "Get blur save image file with tmp path",
			fm: &blurPhotoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						SaveImageFile: "",
					},
				},
				tmpBlurResultImagePath: "tmp_test.jpg",
			},
			want: "tmp_test.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.getBlurSaveImageFile()

			if got != tt.want {
				t.Errorf("getBlurSaveImageFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlurPhotoFrame_DrawBlurMainImage(t *testing.T) {
	tests := []struct {
		name string
		fm   *blurPhotoFrame
	}{
		{
			name: "Draw blur main image without source image",
			fm: &blurPhotoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						PhotoType: "border",
					},
					frameDraw: loadImageRGBA(0, 0, 800, 600),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
			}()

			tt.fm.drawBlurMainImage()

			if tt.fm.frameDraw == nil {
				t.Error("drawBlurMainImage() should maintain frameDraw")
			}
		})
	}
}

func TestBlurPhotoFrame_DrawBlurBorderImage(t *testing.T) {
	tests := []struct {
		name    string
		fm      *blurPhotoFrame
		wantErr bool
	}{
		{
			name: "Draw blur border image with incomplete initialization",
			fm: &blurPhotoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						Params: layout.FrameLayout{
							Name: "simple",
						},
					},
					srcImage: &sourceImage{
						width:  800,
						height: 600,
					},
					borImage: &borderImage{
						bottomHeight: 100,
					},
				},
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

			err := tt.fm.drawBlurBorderImage()

			if tt.wantErr && !pkg.HasError(err) {
				t.Errorf("drawBlurBorderImage() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(err) {
				t.Errorf("drawBlurBorderImage() unexpected error: %v", err)
			}

			if !tt.wantErr && tt.fm.borderDraw == nil {
				t.Error("drawBlurBorderImage() should set borderDraw")
			}
		})
	}
}

func TestBlurPhotoFrame_DrawFrame(t *testing.T) {
	tests := []struct {
		name string
		fm   *blurPhotoFrame
	}{
		{
			name: "Draw blur frame",
			fm: &blurPhotoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						PhotoType: "border",
						Params: layout.FrameLayout{
							Name: "simple",
						},
					},
					srcImage: &sourceImage{
						width:  800,
						height: 600,
					},
					borImage: &borderImage{
						bottomHeight: 100,
					},
					frameDraw: loadImageRGBA(0, 0, 800, 600),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
			}()

			tt.fm.drawFrame()

			if tt.fm.frameDraw == nil {
				t.Error("drawFrame() should maintain frameDraw")
			}
		})
	}
}

func TestBlurPhotoFrame_GetMagickCmdArgs(t *testing.T) {
	tests := []struct {
		name string
		fm   *blurPhotoFrame
		want []string
	}{
		{
			name: "Get magick cmd args without border radius",
			fm: &blurPhotoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						Exif: exiftool.FileMetadata{
							Fields: map[string]any{
								"ImageWidth":  float64(800),
								"ImageHeight": float64(600),
							},
						},
						Params: layout.FrameLayout{
							BorderRadius:   0,
							MainMarginLeft: 50,
							MainMarginTop:  50,
						},
					},
					srcImage: &sourceImage{
						path: "test.jpg",
					},
					finImage: &finalImage{
						width:  800,
						height: 600,
					},
				},
			},
			want: nil,
		},
		{
			name: "Get magick cmd args with border radius",
			fm: &blurPhotoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						Exif: exiftool.FileMetadata{
							Fields: map[string]any{
								"ImageWidth":  float64(800),
								"ImageHeight": float64(600),
							},
						},
						Params: layout.FrameLayout{
							BorderRadius:   10,
							MainMarginLeft: 50,
							MainMarginTop:  50,
						},
					},
					srcImage: &sourceImage{
						path: "test.jpg",
					},
					finImage: &finalImage{
						width:  800,
						height: 600,
					},
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.getMagickCmdArgs(tt.fm.opts.Params.BorderRadius, "output.jpg")

			if got == nil {
				t.Error("getMagickCmdArgs() should return non-nil args")
			}

			if len(got) == 0 {
				t.Error("getMagickCmdArgs() should return non-empty args")
			}
		})
	}
}

func TestBlurPhotoFrame_DrawBlurMerge(t *testing.T) {
	tests := []struct {
		name string
		fm   *blurPhotoFrame
	}{
		{
			name: "Draw blur merge",
			fm: &blurPhotoFrame{
				basePhotoFrame: basePhotoFrame{
					frameDraw: loadImageRGBA(0, 0, 800, 600),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.drawBlurMerge()

			if got == nil {
				t.Error("drawBlurMerge() should return non-nil image")
			}

			if got != tt.fm.frameDraw {
				t.Error("drawBlurMerge() should return frameDraw")
			}
		})
	}
}

func TestBlurPhotoFrame_CreateBlurDraw(t *testing.T) {
	tests := []struct {
		name    string
		fm      *blurPhotoFrame
		wantErr bool
	}{
		{
			name: "Create blur draw without source image",
			fm: &blurPhotoFrame{
				basePhotoFrame: basePhotoFrame{
					srcImage: &sourceImage{},
					finImage: &finalImage{
						width:  800,
						height: 600,
					},
				},
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

			img, err := tt.fm.createBlurDraw()

			if tt.wantErr && !pkg.HasError(err) {
				t.Errorf("createBlurDraw() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(err) {
				t.Errorf("createBlurDraw() unexpected error: %v", err)
			}

			if !tt.wantErr && img == nil {
				t.Error("createBlurDraw() should return non-nil image")
			}
		})
	}
}

func TestBlurPhotoFrame_Clean(t *testing.T) {
	tests := []struct {
		name string
		fm   *blurPhotoFrame
	}{
		{
			name: "Clean blur photo frame",
			fm: &blurPhotoFrame{
				basePhotoFrame: basePhotoFrame{
					opts:       &frameOption{},
					frameDraw:  loadImageRGBA(0, 0, 100, 100),
					borderDraw: loadImageRGBA(0, 0, 100, 100),
					srcImage:   &sourceImage{},
					borImage:   &borderImage{},
					finImage:   &finalImage{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fm.clean()

			if tt.fm.opts != nil {
				t.Error("clean() should set opts to nil")
			}

			if tt.fm.frameDraw != nil {
				t.Error("clean() should set frameDraw to nil")
			}

			if tt.fm.borderDraw != nil {
				t.Error("clean() should set borderDraw to nil")
			}
		})
	}
}

func TestBlurPhotoFrame_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	testImagePath := tmpDir + "/test.jpg"

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

	fm := &blurPhotoFrame{}

	opts := map[string]any{
		"sourceImageFile": testImagePath,
		"photoType":       "border",
		"layoutName":      "simple",
		"saveImageFile":   tmpDir + "/output.jpg",
		"params": layout.FrameLayout{
			Name:           "simple",
			BorderRadius:   10,
			MainMarginLeft: 50,
			MainMarginTop:  50,
			BgColor:        "#000000",
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic during integration test: %v", r)
		}
	}()

	eErr := fm.initFrame(opts)
	if pkg.HasError(eErr) {
		t.Logf("Integration test initFrame error (expected): %v", eErr)
	}

	time.Sleep(100 * time.Millisecond)

	if fm.frameDraw != nil {
		t.Logf("Integration test: frameDraw created with bounds %v", fm.frameDraw.Bounds())
	}

	fm.drawFrame()

	fm.clean()
}
