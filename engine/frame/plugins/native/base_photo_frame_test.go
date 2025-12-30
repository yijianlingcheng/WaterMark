package native

import (
	"testing"

	"WaterMark/pkg"
)

func TestBasePhotoFrame_GetPhotoFrame(t *testing.T) {
	tests := []struct {
		name string
		fm   *basePhotoFrame
	}{
		{
			name: "Get photo frame from base photo frame",
			fm:   &basePhotoFrame{},
		},
		{
			name: "Get photo frame from initialized base photo frame",
			fm: &basePhotoFrame{
				opts: &frameOption{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.getPhotoFrame()

			if got == nil {
				t.Error("getPhotoFrame() returned nil")
			}

			if got != tt.fm {
				t.Error("getPhotoFrame() should return the same instance")
			}
		})
	}
}

func TestBasePhotoFrame_GetOptions(t *testing.T) {
	tests := []struct {
		name string
		fm   *basePhotoFrame
	}{
		{
			name: "Get options from base photo frame with nil opts",
			fm:   &basePhotoFrame{},
		},
		{
			name: "Get options from base photo frame with initialized opts",
			fm: &basePhotoFrame{
				opts: &frameOption{},
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

			got := tt.fm.getOptions()

			if tt.fm.opts != nil && got != tt.fm.opts {
				t.Error("getOptions() should return the same opts instance")
			}
		})
	}
}

func TestBasePhotoFrame_GetBorImage(t *testing.T) {
	tests := []struct {
		name string
		fm   *basePhotoFrame
	}{
		{
			name: "Get border image from base photo frame",
			fm:   &basePhotoFrame{},
		},
		{
			name: "Get border image from initialized base photo frame",
			fm: &basePhotoFrame{
				borImage: &borderImage{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.getBorImage()

			if tt.fm.borImage != nil && got != tt.fm.borImage {
				t.Error("getBorImage() should return the same border image instance")
			}
		})
	}
}

func TestBasePhotoFrame_GetSrcImage(t *testing.T) {
	tests := []struct {
		name string
		fm   *basePhotoFrame
	}{
		{
			name: "Get source image from base photo frame",
			fm:   &basePhotoFrame{},
		},
		{
			name: "Get source image from initialized base photo frame",
			fm: &basePhotoFrame{
				srcImage: &sourceImage{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.getSrcImage()

			if tt.fm.srcImage != nil && got != tt.fm.srcImage {
				t.Error("getSrcImage() should return the same source image instance")
			}
		})
	}
}

func TestBasePhotoFrame_GetBorderDraw(t *testing.T) {
	tests := []struct {
		name string
		fm   *basePhotoFrame
	}{
		{
			name: "Get border draw from base photo frame",
			fm:   &basePhotoFrame{},
		},
		{
			name: "Get border draw from initialized base photo frame",
			fm: &basePhotoFrame{
				borderDraw: loadImageRGBA(0, 0, 100, 100),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.getBorderDraw()

			if tt.fm.borderDraw != nil && got != tt.fm.borderDraw {
				t.Error("getBorderDraw() should return the same border draw instance")
			}
		})
	}
}

func TestBasePhotoFrame_GetFrameDraw(t *testing.T) {
	tests := []struct {
		name string
		fm   *basePhotoFrame
	}{
		{
			name: "Get frame draw from base photo frame",
			fm:   &basePhotoFrame{},
		},
		{
			name: "Get frame draw from initialized base photo frame",
			fm: &basePhotoFrame{
				frameDraw: loadImageRGBA(0, 0, 100, 100),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.getFrameDraw()

			if tt.fm.frameDraw != nil && got != tt.fm.frameDraw {
				t.Error("getFrameDraw() should return the same frame draw instance")
			}
		})
	}
}

func TestBasePhotoFrame_GetLayoutName(t *testing.T) {
	tests := []struct {
		name string
		fm   *basePhotoFrame
		want string
	}{
		{
			name: "Get layout name from base photo frame with nil opts",
			fm:   &basePhotoFrame{},
			want: "",
		},
		{
			name: "Get layout name from base photo frame with empty layout name",
			fm: &basePhotoFrame{
				opts: &frameOption{},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
			}()

			got := tt.fm.getLayoutName()

			if got != tt.want {
				t.Errorf("getLayoutName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasePhotoFrame_GetSaveImageFile(t *testing.T) {
	tests := []struct {
		name string
		fm   *basePhotoFrame
		want string
	}{
		{
			name: "Get save image file from base photo frame with nil opts",
			fm:   &basePhotoFrame{},
			want: "",
		},
		{
			name: "Get save image file from base photo frame with empty save path",
			fm: &basePhotoFrame{
				opts: &frameOption{},
			},
			want: "",
		},
		{
			name: "Get save image file from base photo frame with save path",
			fm: &basePhotoFrame{
				opts: &frameOption{
					SaveImageFile: "/test/path/output.jpg",
				},
			},
			want: "/test/path/output.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
			}()

			got := tt.fm.getSaveImageFile()

			if got != tt.want {
				t.Errorf("getSaveImageFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasePhotoFrame_DrawFrame(t *testing.T) {
	tests := []struct {
		name string
		fm   *basePhotoFrame
	}{
		{
			name: "Draw frame on base photo frame",
			fm:   &basePhotoFrame{},
		},
		{
			name: "Draw frame on initialized base photo frame",
			fm: &basePhotoFrame{
				opts: &frameOption{},
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
		})
	}
}

func TestBasePhotoFrame_Clean(t *testing.T) {
	tests := []struct {
		name string
		fm   *basePhotoFrame
	}{
		{
			name: "Clean base photo frame",
			fm: &basePhotoFrame{
				opts:       &frameOption{},
				frameDraw:  loadImageRGBA(0, 0, 100, 100),
				borderDraw: loadImageRGBA(0, 0, 100, 100),
				srcImage:   &sourceImage{},
				borImage:   &borderImage{},
				finImage:   &finalImage{},
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

func TestBasePhotoFrame_InitSetSize(t *testing.T) {
	tests := []struct {
		name    string
		opts    map[string]any
		wantErr bool
	}{
		{
			name:    "Init set size with nil opts",
			opts:    nil,
			wantErr: true,
		},
		{
			name:    "Init set size with empty opts",
			opts:    map[string]any{},
			wantErr: true,
		},
		{
			name: "Init set size with invalid source image",
			opts: map[string]any{
				"sourceImageFile": "non_existent.jpg",
			},
			wantErr: true,
		},
		{
			name: "Init set size with isBlur option",
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

			fm := &basePhotoFrame{}
			err := fm.initSetSize(tt.opts)

			if tt.wantErr && !pkg.HasError(err) {
				t.Errorf("initSetSize() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(err) {
				t.Errorf("initSetSize() unexpected error: %v", err)
			}
		})
	}
}

func TestBasePhotoFrame_InitFrame(t *testing.T) {
	tests := []struct {
		name    string
		opts    map[string]any
		wantErr bool
	}{
		{
			name:    "Init frame with nil opts",
			opts:    nil,
			wantErr: true,
		},
		{
			name:    "Init frame with empty opts",
			opts:    map[string]any{},
			wantErr: true,
		},
		{
			name: "Init frame with invalid source image",
			opts: map[string]any{
				"sourceImageFile": "non_existent.jpg",
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

			fm := &basePhotoFrame{}
			err := fm.initFrame(tt.opts)

			if tt.wantErr && !pkg.HasError(err) {
				t.Errorf("initFrame() expected error but got none")
			}

			if !tt.wantErr && pkg.HasError(err) {
				t.Errorf("initFrame() unexpected error: %v", err)
			}
		})
	}
}

func TestBasePhotoFrame_GetFrameSize(t *testing.T) {
	tests := []struct {
		name string
		fm   *basePhotoFrame
	}{
		{
			name: "Get frame size from base photo frame",
			fm: &basePhotoFrame{
				opts:     &frameOption{},
				borImage: &borderImage{},
				srcImage: &sourceImage{},
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

			got := tt.fm.getFrameSize()

			if got == nil {
				t.Error("getFrameSize() returned nil")
			}

			if got != nil {
				if _, ok := got["borderLeftWidth"]; !ok {
					t.Error("getFrameSize() should contain borderLeftWidth")
				}
				if _, ok := got["borderRightWidth"]; !ok {
					t.Error("getFrameSize() should contain borderRightWidth")
				}
				if _, ok := got["borderTopHeight"]; !ok {
					t.Error("getFrameSize() should contain borderTopHeight")
				}
				if _, ok := got["borderBottomHeight"]; !ok {
					t.Error("getFrameSize() should contain borderBottomHeight")
				}
				if _, ok := got["sourceWidth"]; !ok {
					t.Error("getFrameSize() should contain sourceWidth")
				}
				if _, ok := got["sourceHeight"]; !ok {
					t.Error("getFrameSize() should contain sourceHeight")
				}
				if _, ok := got["isBlur"]; !ok {
					t.Error("getFrameSize() should contain isBlur")
				}
				if _, ok := got["borderRadius"]; !ok {
					t.Error("getFrameSize() should contain borderRadius")
				}
			}
		})
	}
}

func TestBasePhotoFrame_GetBorderText(t *testing.T) {
	tests := []struct {
		name string
		fm   *basePhotoFrame
	}{
		{
			name: "Get border text from base photo frame",
			fm: &basePhotoFrame{
				opts: &frameOption{},
				borImage: &borderImage{
					textLay: textMarks{},
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

			got := tt.fm.getBorderText()

			if got == nil {
				t.Error("getBorderText() returned nil")
			}
		})
	}
}
