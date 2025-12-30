package frame

import (
	"testing"
)

func TestNewPhotoSize(t *testing.T) {
	tests := []struct {
		name     string
		size     map[string]int
		expected PhotoSize
	}{
		{
			name: "Complete size map with all fields",
			size: map[string]int{
				"borderLeftWidth":    10,
				"borderRightWidth":   20,
				"borderTopHeight":    30,
				"borderBottomHeight": 40,
				"sourceWidth":        100,
				"sourceHeight":       200,
				"borderRadius":       5,
			},
			expected: PhotoSize{
				BorderLeftWidth:    10,
				BorderRightWidth:   20,
				BorderTopHeight:    30,
				BorderBottomHeight: 40,
				SourceWidth:        100,
				SourceHeight:       200,
				BorderRadius:       5,
			},
		},
		{
			name: "Empty size map",
			size: map[string]int{},
			expected: PhotoSize{
				BorderLeftWidth:    0,
				BorderRightWidth:   0,
				BorderTopHeight:    0,
				BorderBottomHeight: 0,
				SourceWidth:        0,
				SourceHeight:       0,
				BorderRadius:       0,
			},
		},
		{
			name: "Size map with only border fields",
			size: map[string]int{
				"borderLeftWidth":    15,
				"borderRightWidth":   25,
				"borderTopHeight":    35,
				"borderBottomHeight": 45,
			},
			expected: PhotoSize{
				BorderLeftWidth:    15,
				BorderRightWidth:   25,
				BorderTopHeight:    35,
				BorderBottomHeight: 45,
				SourceWidth:        0,
				SourceHeight:       0,
				BorderRadius:       0,
			},
		},
		{
			name: "Size map with only source fields",
			size: map[string]int{
				"sourceWidth":  300,
				"sourceHeight": 400,
			},
			expected: PhotoSize{
				BorderLeftWidth:    0,
				BorderRightWidth:   0,
				BorderTopHeight:    0,
				BorderBottomHeight: 0,
				SourceWidth:        300,
				SourceHeight:       400,
				BorderRadius:       0,
			},
		},
		{
			name: "Size map with only borderRadius",
			size: map[string]int{
				"borderRadius": 10,
			},
			expected: PhotoSize{
				BorderLeftWidth:    0,
				BorderRightWidth:   0,
				BorderTopHeight:    0,
				BorderBottomHeight: 0,
				SourceWidth:        0,
				SourceHeight:       0,
				BorderRadius:       10,
			},
		},
		{
			name: "Size map with partial fields",
			size: map[string]int{
				"borderLeftWidth":    5,
				"sourceWidth":        150,
				"borderBottomHeight": 50,
			},
			expected: PhotoSize{
				BorderLeftWidth:    5,
				BorderRightWidth:   0,
				BorderTopHeight:    0,
				BorderBottomHeight: 50,
				SourceWidth:        150,
				SourceHeight:       0,
				BorderRadius:       0,
			},
		},
		{
			name: "Size map with zero values",
			size: map[string]int{
				"borderLeftWidth":    0,
				"borderRightWidth":   0,
				"borderTopHeight":    0,
				"borderBottomHeight": 0,
				"sourceWidth":        0,
				"sourceHeight":       0,
				"borderRadius":       0,
			},
			expected: PhotoSize{
				BorderLeftWidth:    0,
				BorderRightWidth:   0,
				BorderTopHeight:    0,
				BorderBottomHeight: 0,
				SourceWidth:        0,
				SourceHeight:       0,
				BorderRadius:       0,
			},
		},
		{
			name: "Size map with large values",
			size: map[string]int{
				"borderLeftWidth":    1000,
				"borderRightWidth":   2000,
				"borderTopHeight":    3000,
				"borderBottomHeight": 4000,
				"sourceWidth":        5000,
				"sourceHeight":       6000,
				"borderRadius":       500,
			},
			expected: PhotoSize{
				BorderLeftWidth:    1000,
				BorderRightWidth:   2000,
				BorderTopHeight:    3000,
				BorderBottomHeight: 4000,
				SourceWidth:        5000,
				SourceHeight:       6000,
				BorderRadius:       500,
			},
		},
		{
			name: "Size map with negative values",
			size: map[string]int{
				"borderLeftWidth":    -10,
				"borderRightWidth":   -20,
				"borderTopHeight":    -30,
				"borderBottomHeight": -40,
				"sourceWidth":        -100,
				"sourceHeight":       -200,
				"borderRadius":       -5,
			},
			expected: PhotoSize{
				BorderLeftWidth:    -10,
				BorderRightWidth:   -20,
				BorderTopHeight:    -30,
				BorderBottomHeight: -40,
				SourceWidth:        -100,
				SourceHeight:       -200,
				BorderRadius:       -5,
			},
		},
		{
			name: "Size map with mixed positive and negative values",
			size: map[string]int{
				"borderLeftWidth":  10,
				"borderRightWidth": -20,
				"sourceWidth":      100,
				"borderRadius":     0,
			},
			expected: PhotoSize{
				BorderLeftWidth:    10,
				BorderRightWidth:   -20,
				BorderTopHeight:    0,
				BorderBottomHeight: 0,
				SourceWidth:        100,
				SourceHeight:       0,
				BorderRadius:       0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPhotoSize(tt.size)

			if got.BorderLeftWidth != tt.expected.BorderLeftWidth {
				t.Errorf(
					"NewPhotoSize().BorderLeftWidth = %v, want %v",
					got.BorderLeftWidth,
					tt.expected.BorderLeftWidth,
				)
			}
			if got.BorderRightWidth != tt.expected.BorderRightWidth {
				t.Errorf(
					"NewPhotoSize().BorderRightWidth = %v, want %v",
					got.BorderRightWidth,
					tt.expected.BorderRightWidth,
				)
			}
			if got.BorderTopHeight != tt.expected.BorderTopHeight {
				t.Errorf(
					"NewPhotoSize().BorderTopHeight = %v, want %v",
					got.BorderTopHeight,
					tt.expected.BorderTopHeight,
				)
			}
			if got.BorderBottomHeight != tt.expected.BorderBottomHeight {
				t.Errorf(
					"NewPhotoSize().BorderBottomHeight = %v, want %v",
					got.BorderBottomHeight,
					tt.expected.BorderBottomHeight,
				)
			}
			if got.SourceWidth != tt.expected.SourceWidth {
				t.Errorf("NewPhotoSize().SourceWidth = %v, want %v", got.SourceWidth, tt.expected.SourceWidth)
			}
			if got.SourceHeight != tt.expected.SourceHeight {
				t.Errorf("NewPhotoSize().SourceHeight = %v, want %v", got.SourceHeight, tt.expected.SourceHeight)
			}
			if got.BorderRadius != tt.expected.BorderRadius {
				t.Errorf("NewPhotoSize().BorderRadius = %v, want %v", got.BorderRadius, tt.expected.BorderRadius)
			}
		})
	}
}

func TestPhotoSize(t *testing.T) {
	tests := []struct {
		name string
		size PhotoSize
	}{
		{
			name: "PhotoSize with all fields set",
			size: PhotoSize{
				BorderLeftWidth:    10,
				BorderRightWidth:   20,
				BorderTopHeight:    30,
				BorderBottomHeight: 40,
				SourceWidth:        100,
				SourceHeight:       200,
				BorderRadius:       5,
			},
		},
		{
			name: "PhotoSize with zero values",
			size: PhotoSize{},
		},
		{
			name: "PhotoSize with only border widths",
			size: PhotoSize{
				BorderLeftWidth:  15,
				BorderRightWidth: 25,
			},
		},
		{
			name: "PhotoSize with only border heights",
			size: PhotoSize{
				BorderTopHeight:    35,
				BorderBottomHeight: 45,
			},
		},
		{
			name: "PhotoSize with only source dimensions",
			size: PhotoSize{
				SourceWidth:  300,
				SourceHeight: 400,
			},
		},
		{
			name: "PhotoSize with only border radius",
			size: PhotoSize{
				BorderRadius: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := tt.size

			if size.BorderLeftWidth != tt.size.BorderLeftWidth {
				t.Errorf("PhotoSize.BorderLeftWidth = %v, want %v", size.BorderLeftWidth, tt.size.BorderLeftWidth)
			}
			if size.BorderRightWidth != tt.size.BorderRightWidth {
				t.Errorf("PhotoSize.BorderRightWidth = %v, want %v", size.BorderRightWidth, tt.size.BorderRightWidth)
			}
			if size.BorderTopHeight != tt.size.BorderTopHeight {
				t.Errorf("PhotoSize.BorderTopHeight = %v, want %v", size.BorderTopHeight, tt.size.BorderTopHeight)
			}
			if size.BorderBottomHeight != tt.size.BorderBottomHeight {
				t.Errorf(
					"PhotoSize.BorderBottomHeight = %v, want %v",
					size.BorderBottomHeight,
					tt.size.BorderBottomHeight,
				)
			}
			if size.SourceWidth != tt.size.SourceWidth {
				t.Errorf("PhotoSize.SourceWidth = %v, want %v", size.SourceWidth, tt.size.SourceWidth)
			}
			if size.SourceHeight != tt.size.SourceHeight {
				t.Errorf("PhotoSize.SourceHeight = %v, want %v", size.SourceHeight, tt.size.SourceHeight)
			}
			if size.BorderRadius != tt.size.BorderRadius {
				t.Errorf("PhotoSize.BorderRadius = %v, want %v", size.BorderRadius, tt.size.BorderRadius)
			}
		})
	}
}
