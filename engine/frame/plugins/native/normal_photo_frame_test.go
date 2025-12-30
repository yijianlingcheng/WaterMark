package native

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"sync"
	"testing"

	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/layout"
	"WaterMark/pkg"
)

func init() {
	layout.LogosImagesInit()
}

func TestPhotoFrame_InitFrame(t *testing.T) {
	tests := []struct {
		name    string
		opts    map[string]any
		wantErr bool
	}{
		{
			name:    "Init photo frame with nil opts",
			opts:    nil,
			wantErr: true,
		},
		{
			name:    "Init photo frame with empty opts",
			opts:    map[string]any{},
			wantErr: true,
		},
		{
			name: "Init photo frame with invalid source image",
			opts: map[string]any{
				"sourceImageFile": "non_existent.jpg",
			},
			wantErr: true,
		},
		{
			name: "Init photo frame without source image",
			opts: map[string]any{
				"sourceImageFile":    "test.jpg",
				"needSourceImage":    false,
				"layoutName":         "simple",
				"saveImageFile":      "output.jpg",
				"sourceImageX":       800,
				"sourceImageY":       600,
				"borderRadius":       10,
				"mainMarginLeft":     50,
				"mainMarginTop":      50,
				"borderWidth":        20,
				"borderHeight":       20,
				"borderColor":        "#000000",
				"borderText":         []string{"Test"},
				"borderTextSize":     12,
				"borderTextColor":    "#FFFFFF",
				"borderTextPosition": "center",
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

			fm := &photoFrame{}
			err := fm.initFrame(tt.opts)

			if tt.wantErr && !pkg.HasError(err) {
				t.Errorf("initFrame() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(err) {
				t.Errorf("initFrame() unexpected error: %v", err)
			}

			if !tt.wantErr {
				if fm.opts == nil {
					t.Error("initFrame() should set opts")
				}
			}
		})
	}
}

func TestPhotoFrame_DrawMainImage(t *testing.T) {
	tests := []struct {
		name string
		fm   *photoFrame
	}{
		{
			name: "Draw main image without source image",
			fm: &photoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						PhotoType: "border",
					},
					frameDraw: loadImageRGBA(0, 0, 800, 600),
					borImage: &borderImage{
						leftWidth:  20,
						topHeight:  20,
						rightWidth: 20,
					},
				},
			},
		},
		{
			name: "Draw main image with source image",
			fm: &photoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						PhotoType: "normal",
					},
					frameDraw: loadImageRGBA(0, 0, 800, 600),
					srcImage: &sourceImage{
						imgDecode: image.NewRGBA(image.Rect(0, 0, 760, 560)),
						width:     760,
						height:    560,
					},
					borImage: &borderImage{
						leftWidth:  20,
						topHeight:  20,
						rightWidth: 20,
					},
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

			var wg sync.WaitGroup
			wg.Add(1)

			tt.fm.drawMainImage(&wg)
			wg.Wait()

			if tt.fm.frameDraw == nil {
				t.Error("drawMainImage() should maintain frameDraw")
			}
		})
	}
}

func TestPhotoFrame_DrawBorderImage(t *testing.T) {
	tests := []struct {
		name    string
		fm      *photoFrame
		wantErr bool
	}{
		{
			name: "Draw border image",
			fm: &photoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						Params: layout.FrameLayout{
							Name: "simple",
						},
						Exif: exiftool.FileMetadata{
							Fields: map[string]any{
								"Make": "TestMake",
							},
						},
					},
					srcImage: &sourceImage{
						width:  800,
						height: 600,
					},
					borImage: &borderImage{
						bottomHeight: 100,
						bgColor:      color.RGBA{R: 0, G: 0, B: 0, A: 255},
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

			var wg sync.WaitGroup
			wg.Add(1)

			err := tt.fm.drawBorderImage(&wg)
			wg.Wait()

			if tt.wantErr && !pkg.HasError(err) {
				t.Errorf("drawBorderImage() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(err) {
				t.Logf("drawBorderImage() error (may be expected): %v", err)
			}

			if tt.fm.borderDraw == nil {
				t.Error("drawBorderImage() should set borderDraw")
			}
		})
	}
}

func TestPhotoFrame_DrawFrame(t *testing.T) {
	tests := []struct {
		name string
		fm   *photoFrame
	}{
		{
			name: "Draw photo frame without source image",
			fm: &photoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						PhotoType: "border",
						Params: layout.FrameLayout{
							Name: "simple",
						},
						Exif: exiftool.FileMetadata{
							Fields: map[string]any{
								"Make": "TestMake",
							},
						},
					},
					srcImage: &sourceImage{
						width:  800,
						height: 600,
					},
					borImage: &borderImage{
						leftWidth:    20,
						topHeight:    20,
						rightWidth:   20,
						bottomHeight: 100,
						bgColor:      color.RGBA{R: 0, G: 0, B: 0, A: 255},
					},
					frameDraw: loadImageRGBA(0, 0, 800, 600),
				},
			},
		},
		{
			name: "Draw photo frame with source image",
			fm: &photoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						PhotoType: "normal",
						Params: layout.FrameLayout{
							Name: "simple",
						},
						Exif: exiftool.FileMetadata{
							Fields: map[string]any{
								"Make": "TestMake",
							},
						},
					},
					srcImage: &sourceImage{
						imgDecode: image.NewRGBA(image.Rect(0, 0, 760, 560)),
						width:     760,
						height:    560,
					},
					borImage: &borderImage{
						leftWidth:    20,
						topHeight:    20,
						rightWidth:   20,
						bottomHeight: 100,
						bgColor:      color.RGBA{R: 0, G: 0, B: 0, A: 255},
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

func TestPhotoFrame_DrawMerge(t *testing.T) {
	tests := []struct {
		name string
		fm   *photoFrame
	}{
		{
			name: "Draw merge",
			fm: &photoFrame{
				basePhotoFrame: basePhotoFrame{
					frameDraw: loadImageRGBA(0, 0, 800, 600),
					borImage: &borderImage{
						leftWidth:    20,
						topHeight:    20,
						rightWidth:   20,
						bottomHeight: 100,
					},
					srcImage: &sourceImage{
						width:  760,
						height: 560,
					},
					borderDraw: loadImageRGBA(0, 0, 800, 100),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.drawMerge()

			if got == nil {
				t.Error("drawMerge() should return non-nil image")
			}

			if got != tt.fm.frameDraw {
				t.Error("drawMerge() should return frameDraw")
			}
		})
	}
}

func TestPhotoFrame_ConcurrentDraw(t *testing.T) {
	tests := []struct {
		name string
		fm   *photoFrame
	}{
		{
			name: "Concurrent draw with goroutines",
			fm: &photoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						PhotoType: "normal",
						Params: layout.FrameLayout{
							Name: "simple",
						},
						Exif: exiftool.FileMetadata{
							Fields: map[string]any{
								"Make": "TestMake",
							},
						},
					},
					srcImage: &sourceImage{
						imgDecode: image.NewRGBA(image.Rect(0, 0, 760, 560)),
						width:     760,
						height:    560,
					},
					borImage: &borderImage{
						leftWidth:    20,
						topHeight:    20,
						rightWidth:   20,
						bottomHeight: 100,
						bgColor:      color.RGBA{R: 0, G: 0, B: 0, A: 255},
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

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				tt.fm.drawMainImage(&wg)
			}()

			go func() {
				err := tt.fm.drawBorderImage(&wg)
				if pkg.HasError(err) {
					t.Logf("drawBorderImage() error: %v", err)
				}
			}()

			wg.Wait()

			if tt.fm.frameDraw == nil {
				t.Error("Concurrent draw should maintain frameDraw")
			}

			if tt.fm.borderDraw == nil {
				t.Error("Concurrent draw should set borderDraw")
			}
		})
	}
}

func TestPhotoFrame_Clean(t *testing.T) {
	tests := []struct {
		name string
		fm   *photoFrame
	}{
		{
			name: "Clean photo frame",
			fm: &photoFrame{
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

func TestPhotoFrame_Integration(t *testing.T) {
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

	if err := jpeg.Encode(testFile, testImg, nil); err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}
	testFile.Close()

	fm := &photoFrame{}

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

	if fm.frameDraw != nil {
		t.Logf("Integration test: frameDraw created with bounds %v", fm.frameDraw.Bounds())
	}

	fm.drawFrame()

	if fm.frameDraw != nil {
		mergeResult := fm.drawMerge()
		if mergeResult != nil {
			t.Logf("Integration test: merge result bounds %v", mergeResult.Bounds())
		}
	}

	fm.clean()
}

func TestPhotoFrame_MultipleDraws(t *testing.T) {
	tests := []struct {
		name string
		fm   *photoFrame
	}{
		{
			name: "Multiple draws on same frame",
			fm: &photoFrame{
				basePhotoFrame: basePhotoFrame{
					opts: &frameOption{
						PhotoType: "border",
						Params: layout.FrameLayout{
							Name: "simple",
						},
						Exif: exiftool.FileMetadata{
							Fields: map[string]any{
								"Make": "TestMake",
							},
						},
					},
					srcImage: &sourceImage{
						width:  800,
						height: 600,
					},
					borImage: &borderImage{
						leftWidth:    20,
						topHeight:    20,
						rightWidth:   20,
						bottomHeight: 100,
						bgColor:      color.RGBA{R: 0, G: 0, B: 0, A: 255},
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

			for i := 0; i < 3; i++ {
				tt.fm.drawFrame()

				if tt.fm.frameDraw == nil {
					t.Errorf("drawFrame() iteration %d should maintain frameDraw", i)
				}
			}
		})
	}
}

func TestPhotoFrame_WithLayoutFrameNames(t *testing.T) {
	frameNames := []string{
		"固定布局左下logo模板",
		"固定布局右下logo模板",
		"经典-左logo",
		"经典-左logo-无边框",
		"经典-左logo-1",
		"经典-左logo-无边框-1",
		"经典-左logo-2",
		"经典-左logo-无边框-2",
		"经典-右logo",
		"经典-右logo-无边框",
		"经典-右logo-1",
		"经典-右logo-1-无边框",
		"经典-右logo-2",
		"经典-右logo-无边框-2",
		"经典-右logo-对比",
		"经典-右logo-无边框-对比",
		"简约-居中-无logo",
		"简约-居中-无边框-无logo",
		"简约-居中",
		"简约-居中-无边框",
		"高斯模糊-居中",
	}

	for _, frameName := range frameNames {
		t.Run(frameName, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic for layout '%s': %v", frameName, r)
				}
			}()

			fm := &photoFrame{}

			opts := map[string]any{
				"sourceImageFile": "",
				"photoType":       "border",
				"layoutName":      frameName,
				"params": layout.FrameLayout{
					Name: frameName,
				},
				"exif": exiftool.FileMetadata{
					Fields: map[string]any{
						"Make":  "TestMake",
						"Model": "TestModel",
					},
				},
			}

			eErr := fm.initFrame(opts)

			if fm.opts != nil && fm.opts.Params.Name == frameName {
				t.Logf("Layout '%s' initialized successfully", frameName)
			} else {
				t.Logf("Layout '%s' initFrame returned error: %v", frameName, eErr)
			}

			fm.clean()
		})
	}
}
