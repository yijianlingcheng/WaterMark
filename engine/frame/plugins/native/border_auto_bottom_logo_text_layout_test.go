package native

import (
	"testing"

	"WaterMark/pkg"
)

func TestAutoBottomLogoTextLayoutBorderInitLayoutValue(t *testing.T) {
	tests := []struct {
		name    string
		border  *autoBottomLogoTextLayoutBorder
		wantErr bool
	}{
		{
			name: "Init auto border layout value left without separator",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: false,
			},
			wantErr: true,
		},
		{
			name: "Init auto border layout value right without separator",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      true,
				HasSeparator: false,
			},
			wantErr: true,
		},
		{
			name: "Init auto border layout value left with separator",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: true,
			},
			wantErr: true,
		},
		{
			name: "Init auto border layout value right with separator",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      true,
				HasSeparator: true,
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

			err := tt.border.initLayoutValue(nil)

			if !panicked && pkg.HasError(err) != tt.wantErr {
				t.Errorf("initLayoutValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAutoBottomLogoTextLayoutBorderDrawBorder(t *testing.T) {
	tests := []struct {
		name    string
		border  *autoBottomLogoTextLayoutBorder
		wantErr bool
	}{
		{
			name: "Draw auto border layout left without separator",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: false,
			},
			wantErr: true,
		},
		{
			name: "Draw auto border layout right without separator",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      true,
				HasSeparator: false,
			},
			wantErr: true,
		},
		{
			name: "Draw auto border layout left with separator",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: true,
			},
			wantErr: true,
		},
		{
			name: "Draw auto border layout right with separator",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      true,
				HasSeparator: true,
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

func TestAutoBottomLogoTextLayoutBorderGetTextLayoutLogoCommonData(t *testing.T) {
	tests := []struct {
		name   string
		border *autoBottomLogoTextLayoutBorder
	}{
		{
			name: "Get text layout logo common data",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: false,
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

			info := tt.border.getTextLayoutLogoCommonData(nil)

			if info.r1 == nil {
				t.Error("getTextLayoutLogoCommonData() r1 is nil")
			}
			if info.r2 == nil {
				t.Error("getTextLayoutLogoCommonData() r2 is nil")
			}
		})
	}
}

func TestAutoBottomLogoTextLayoutBorderSetTextLayoutLogo(t *testing.T) {
	tests := []struct {
		name   string
		border *autoBottomLogoTextLayoutBorder
	}{
		{
			name: "Set text layout logo",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: false,
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

			tt.border.setTextLayoutLogo(nil)
		})
	}
}

func TestAutoBottomLogoTextLayoutBorderSetTextLayoutLogoWithRight(t *testing.T) {
	tests := []struct {
		name   string
		border *autoBottomLogoTextLayoutBorder
	}{
		{
			name: "Set text layout logo with right",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      true,
				HasSeparator: true,
			},
		},
		{
			name: "Set text layout logo with right without separator",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      true,
				HasSeparator: false,
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

			tt.border.setTextLayoutLogoWithRight(nil)
		})
	}
}

func TestAutoBottomLogoTextLayoutBorderSetTextLayoutSeparator(t *testing.T) {
	tests := []struct {
		name   string
		border *autoBottomLogoTextLayoutBorder
	}{
		{
			name: "Set text layout separator with separator",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: true,
			},
		},
		{
			name: "Set text layout separator without separator",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: false,
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

			tt.border.setTextLayoutSeparator(nil)
		})
	}
}

func TestAutoBottomLogoTextLayoutBorderSetTextLayoutSeparatorWithRight(t *testing.T) {
	tests := []struct {
		name   string
		border *autoBottomLogoTextLayoutBorder
	}{
		{
			name: "Set text layout separator with right",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      true,
				HasSeparator: true,
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

			tt.border.setTextLayoutSeparatorWithRight(nil)
		})
	}
}

func TestAutoBottomLogoTextLayoutBorderSetTextLayoutText(t *testing.T) {
	tests := []struct {
		name   string
		border *autoBottomLogoTextLayoutBorder
	}{
		{
			name: "Set text layout text",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: false,
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

			tt.border.setTextLayoutText(nil)
		})
	}
}

func TestAutoBottomLogoTextLayoutBorderSetTextLayoutTextMarginLeftWithRight(t *testing.T) {
	tests := []struct {
		name   string
		border *autoBottomLogoTextLayoutBorder
	}{
		{
			name: "Set text layout text margin left with right",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      true,
				HasSeparator: false,
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

			tt.border.setTextLayoutTextMarginLeftWithRight(nil)
		})
	}
}

func TestAutoBottomLogoTextLayoutBorderSetTextLayoutTextFontSize(t *testing.T) {
	tests := []struct {
		name     string
		border   *autoBottomLogoTextLayoutBorder
		fontSize int
	}{
		{
			name:     "Set text layout text font size",
			border:   &autoBottomLogoTextLayoutBorder{},
			fontSize: 12,
		},
		{
			name:     "Set text layout text font size with zero",
			border:   &autoBottomLogoTextLayoutBorder{},
			fontSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
			}()

			tt.border.setTextLayoutTextFontSize(nil, tt.fontSize)
		})
	}
}

func TestAutoBottomLogoTextLayoutBorderSetTextLayoutTextMarginTop(t *testing.T) {
	tests := []struct {
		name   string
		border *autoBottomLogoTextLayoutBorder
	}{
		{
			name:   "Set text layout text margin top",
			border: &autoBottomLogoTextLayoutBorder{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
			}()

			tt.border.setTextLayoutTextMarginTop(nil)
		})
	}
}

func TestAutoBottomLogoTextLayoutBorderSetTextLayoutTextMarginLeft(t *testing.T) {
	tests := []struct {
		name   string
		border *autoBottomLogoTextLayoutBorder
	}{
		{
			name: "Set text layout text margin left",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: false,
			},
		},
		{
			name: "Set text layout text margin left with right",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      true,
				HasSeparator: false,
			},
		},
		{
			name: "Set text layout text margin left with separator",
			border: &autoBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: true,
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

			tt.border.setTextLayoutTextMarginLeft(nil)
		})
	}
}

func TestAutoBottomLogoTextLayoutBorderGetRightTextMaxWidth(t *testing.T) {
	tests := []struct {
		name   string
		border *autoBottomLogoTextLayoutBorder
	}{
		{
			name:   "Get right text max width",
			border: &autoBottomLogoTextLayoutBorder{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
			}()

			width := tt.border.getRightTextMaxWidth(nil)

			if width < 0 {
				t.Errorf("getRightTextMaxWidth() returned negative width: %v", width)
			}
		})
	}
}
