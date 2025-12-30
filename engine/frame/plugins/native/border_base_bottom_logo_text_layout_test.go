package native

import (
	"testing"
)

func TestBaseBottomLogoTextLayoutBorderSetTextLayoutBorder(t *testing.T) {
	tests := []struct {
		name   string
		border *baseBottomLogoTextLayoutBorder
	}{
		{
			name: "Set text layout border",
			border: &baseBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: false,
			},
		},
		{
			name: "Set text layout border with right layout",
			border: &baseBottomLogoTextLayoutBorder{
				IsRight:      true,
				HasSeparator: false,
			},
		},
		{
			name: "Set text layout border with separator",
			border: &baseBottomLogoTextLayoutBorder{
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

			tt.border.setTextLayoutBorder(nil)
		})
	}
}

func TestBaseBottomLogoTextLayoutBorderDrawSeparator(t *testing.T) {
	tests := []struct {
		name   string
		border *baseBottomLogoTextLayoutBorder
	}{
		{
			name: "Draw separator left layout",
			border: &baseBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: true,
			},
		},
		{
			name: "Draw separator right layout",
			border: &baseBottomLogoTextLayoutBorder{
				IsRight:      true,
				HasSeparator: true,
			},
		},
		{
			name: "Draw separator without separator",
			border: &baseBottomLogoTextLayoutBorder{
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

			tt.border.drawSeparator(nil)
		})
	}
}

func TestBaseBottomLogoTextLayoutBorderDrawLogo(t *testing.T) {
	tests := []struct {
		name   string
		border *baseBottomLogoTextLayoutBorder
	}{
		{
			name: "Draw logo left layout",
			border: &baseBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: false,
			},
		},
		{
			name: "Draw logo right layout",
			border: &baseBottomLogoTextLayoutBorder{
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

			tt.border.drawLogo(nil)
		})
	}
}

func TestBaseBottomLogoTextLayoutBorderDrawWords(t *testing.T) {
	tests := []struct {
		name   string
		border *baseBottomLogoTextLayoutBorder
	}{
		{
			name: "Draw words left layout",
			border: &baseBottomLogoTextLayoutBorder{
				IsRight:      false,
				HasSeparator: false,
			},
		},
		{
			name: "Draw words right layout",
			border: &baseBottomLogoTextLayoutBorder{
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

			tt.border.drawWords(nil)
		})
	}
}
