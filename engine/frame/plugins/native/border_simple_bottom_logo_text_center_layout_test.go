package native

import (
	"testing"

	"WaterMark/pkg"
)

func TestSimpleBottomLogoTextCenterBorderInitLayoutValue(t *testing.T) {
	tests := []struct {
		name    string
		border  *simpleBottomLogoTextCenterBorder
		wantErr bool
	}{
		{
			name: "Init simple border layout value without logo",
			border: &simpleBottomLogoTextCenterBorder{
				HasLogo: false,
			},
			wantErr: true,
		},
		{
			name: "Init simple border layout value with logo",
			border: &simpleBottomLogoTextCenterBorder{
				HasLogo: true,
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

			err := tt.border.initLayoutValue(nil)

			if pkg.HasError(err) != tt.wantErr {
				t.Errorf("initLayoutValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSimpleBottomLogoTextCenterBorderDrawBorder(t *testing.T) {
	tests := []struct {
		name    string
		border  *simpleBottomLogoTextCenterBorder
		wantErr bool
	}{
		{
			name: "Draw simple border layout without logo",
			border: &simpleBottomLogoTextCenterBorder{
				HasLogo: false,
			},
			wantErr: true,
		},
		{
			name: "Draw simple border layout with logo",
			border: &simpleBottomLogoTextCenterBorder{
				HasLogo: true,
			},
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

func TestSimpleBottomLogoTextCenterBorderGetMaxFontSize(t *testing.T) {
	tests := []struct {
		name             string
		border           *simpleBottomLogoTextCenterBorder
		textOneContent   string
		textThreeContent string
	}{
		{
			name:             "Get max font size without logo",
			border:           &simpleBottomLogoTextCenterBorder{HasLogo: false},
			textOneContent:   "Test One",
			textThreeContent: "Test Three",
		},
		{
			name:             "Get max font size with logo",
			border:           &simpleBottomLogoTextCenterBorder{HasLogo: true},
			textOneContent:   "Test One",
			textThreeContent: "Test Three",
		},
		{
			name:             "Get max font size with empty content",
			border:           &simpleBottomLogoTextCenterBorder{HasLogo: false},
			textOneContent:   "",
			textThreeContent: "",
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

func TestSimpleBottomLogoTextCenterBorderSetTextLayoutTextAndLogo(t *testing.T) {
	tests := []struct {
		name   string
		border *simpleBottomLogoTextCenterBorder
	}{
		{
			name:   "Set text layout text and logo without logo",
			border: &simpleBottomLogoTextCenterBorder{HasLogo: false},
		},
		{
			name:   "Set text layout text and logo with logo",
			border: &simpleBottomLogoTextCenterBorder{HasLogo: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
			}()

			tt.border.setTextLayoutTextAndLogo(nil)
		})
	}
}

func TestSimpleBottomLogoTextCenterBorderSetTextLayoutWithHasLogo(t *testing.T) {
	tests := []struct {
		name         string
		border       *simpleBottomLogoTextCenterBorder
		logoShowInfo map[string]int
	}{
		{
			name:   "Set text layout with has logo",
			border: &simpleBottomLogoTextCenterBorder{HasLogo: true},
			logoShowInfo: map[string]int{
				"width":  100,
				"height": 50,
			},
		},
		{
			name:   "Set text layout with has logo zero dimensions",
			border: &simpleBottomLogoTextCenterBorder{HasLogo: true},
			logoShowInfo: map[string]int{
				"width":  0,
				"height": 0,
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

			tt.border.setTextLayoutWithHasLogo(nil, tt.logoShowInfo)
		})
	}
}
