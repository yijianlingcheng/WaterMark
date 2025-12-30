package native

import (
	"testing"

	"WaterMark/pkg"
)

func TestFixedBottomLogoTextLayoutBorderDrawBorder(t *testing.T) {
	tests := []struct {
		name    string
		border  *fixedBottomLogoTextLayoutBorder
		wantErr bool
	}{
		{
			name: "Draw fixed border left layout",
			border: &fixedBottomLogoTextLayoutBorder{
				IsRight: false,
			},
			wantErr: true,
		},
		{
			name: "Draw fixed border right layout",
			border: &fixedBottomLogoTextLayoutBorder{
				IsRight: true,
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

			err := tt.border.drawBorder(nil)

			if pkg.HasError(err) != tt.wantErr {
				t.Errorf("drawBorder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFixedBottomLogoTextLayoutBorderInitLayoutValue(t *testing.T) {
	tests := []struct {
		name    string
		border  *fixedBottomLogoTextLayoutBorder
		wantErr bool
	}{
		{
			name: "Init fixed border layout value",
			border: &fixedBottomLogoTextLayoutBorder{
				IsRight: false,
			},
			wantErr: false,
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
