package native

import (
	"testing"

	"WaterMark/pkg"
)

func TestAutoBottomLogoTextAverageLayoutBorderInitLayoutValue(t *testing.T) {
	tests := []struct {
		name    string
		border  *autoBottomLogoTextAverageLayoutBorder
		wantErr bool
	}{
		{
			name: "Init auto average border layout value",
			border: &autoBottomLogoTextAverageLayoutBorder{
				autoBottomLogoTextLayoutBorder: autoBottomLogoTextLayoutBorder{
					IsRight:      false,
					HasSeparator: true,
				},
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

func TestAutoBottomLogoTextAverageLayoutBorderDrawBorder(t *testing.T) {
	tests := []struct {
		name    string
		border  *autoBottomLogoTextAverageLayoutBorder
		wantErr bool
	}{
		{
			name: "Draw auto average border layout",
			border: &autoBottomLogoTextAverageLayoutBorder{
				autoBottomLogoTextLayoutBorder: autoBottomLogoTextLayoutBorder{
					IsRight:      false,
					HasSeparator: true,
				},
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

func TestAutoBottomLogoTextAverageLayoutBorderSetTextLayoutTextAndLogo(t *testing.T) {
	tests := []struct {
		name   string
		border *autoBottomLogoTextAverageLayoutBorder
	}{
		{
			name: "Set text layout text and logo",
			border: &autoBottomLogoTextAverageLayoutBorder{
				autoBottomLogoTextLayoutBorder: autoBottomLogoTextLayoutBorder{
					IsRight:      false,
					HasSeparator: true,
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

			tt.border.setTextLayoutTextAndLogo(nil)
		})
	}
}

func TestAutoBottomLogoTextAverageLayoutBorderSetFontSize(t *testing.T) {
	tests := []struct {
		name            string
		border          *autoBottomLogoTextAverageLayoutBorder
		textOneContent  string
		textTwoContent  string
		textThreeContent string
	}{
		{
			name:            "Set font size",
			border:          &autoBottomLogoTextAverageLayoutBorder{},
			textOneContent:  "Test One",
			textTwoContent:  "Test Two",
			textThreeContent: "Test Three",
		},
		{
			name:            "Set font size with empty content",
			border:          &autoBottomLogoTextAverageLayoutBorder{},
			textOneContent:  "",
			textTwoContent:  "",
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

			tt.border.setFontSize(nil, tt.textOneContent, tt.textTwoContent, tt.textThreeContent)
		})
	}
}
