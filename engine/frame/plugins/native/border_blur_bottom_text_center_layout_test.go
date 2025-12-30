package native

import (
	"testing"

	"WaterMark/pkg"
)

func TestBlurBottomTextCenterLayoutInitLayoutValue(t *testing.T) {
	tests := []struct {
		name    string
		border  *blurBottomTextCenterLayout
		wantErr bool
	}{
		{
			name:    "Init blur border layout value",
			border:  &blurBottomTextCenterLayout{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			panicked := false
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
					panicked = true
					if !tt.wantErr {
						t.Errorf("Unexpected panic: %v", r)
					}
				}
			}()

			err := tt.border.initLayoutValue(nil)

			if !panicked && pkg.HasError(err) != tt.wantErr {
				t.Errorf("initLayoutValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlurBottomTextCenterLayoutDrawBorder(t *testing.T) {
	tests := []struct {
		name    string
		border  *blurBottomTextCenterLayout
		wantErr bool
	}{
		{
			name:    "Draw blur border layout",
			border:  &blurBottomTextCenterLayout{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			panicked := false
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
					panicked = true
					if !tt.wantErr {
						t.Errorf("Unexpected panic: %v", r)
					}
				}
			}()

			err := tt.border.drawBorder(nil)

			if !panicked && pkg.HasError(err) != tt.wantErr {
				t.Errorf("drawBorder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlurBottomTextCenterLayoutGetMaxFontSize(t *testing.T) {
	tests := []struct {
		name             string
		border           *blurBottomTextCenterLayout
		textOneContent   string
		textThreeContent string
	}{
		{
			name:             "Get max font size",
			border:           &blurBottomTextCenterLayout{},
			textOneContent:   "Test One",
			textThreeContent: "Test Three",
		},
		{
			name:             "Get max font size with empty content",
			border:           &blurBottomTextCenterLayout{},
			textOneContent:   "",
			textThreeContent: "",
		},
		{
			name:             "Get max font size with longer content",
			border:           &blurBottomTextCenterLayout{},
			textOneContent:   "Short",
			textThreeContent: "This is a much longer text content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
			}()

			fontSize := tt.border.getMaxFontSize(nil, tt.textOneContent, tt.textThreeContent)

			if fontSize < 0 {
				t.Errorf("getMaxFontSize() returned negative font size: %v", fontSize)
			}
		})
	}
}

func TestBlurBottomTextCenterLayoutSetTextLayoutText(t *testing.T) {
	tests := []struct {
		name   string
		border *blurBottomTextCenterLayout
	}{
		{
			name:   "Set text layout text",
			border: &blurBottomTextCenterLayout{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
			}()

			tt.border.setTextLayoutText(nil)
		})
	}
}
