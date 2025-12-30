package native

import (
	"testing"
)

func TestConstants(t *testing.T) {
	tests := []struct {
		name  string
		constant string
		expected string
	}{
		{
			name:  "SEPARATOR_COLOR constant",
			constant: SEPARATOR_COLOR,
			expected: "203,203,201,255",
		},
		{
			name:  "COLOR constant",
			constant: COLOR,
			expected: "255,255,255,255",
		},
		{
			name:  "GPS_OR_DATETIME constant",
			constant: GPS_OR_DATETIME,
			expected: "GPS_OR_DATETIME",
		},
		{
			name:  "DATE_TIME_ORIGINAL constant",
			constant: DATE_TIME_ORIGINAL,
			expected: "DateTimeOriginal",
		},
		{
			name:  "GPS_POSITION constant",
			constant: GPS_POSITION,
			expected: "GPSPosition",
		},
		{
			name:  "FOCAL_LENGTH constant",
			constant: FOCAL_LENGTH,
			expected: "FocalLength",
		},
		{
			name:  "PHOTO_TYPE_BORDER constant",
			constant: PHOTO_TYPE_BORDER,
			expected: "border",
		},
		{
			name:  "JPG_FILE_TYPE constant",
			constant: JPG_FILE_TYPE,
			expected: ".jpg",
		},
		{
			name:  "JPEG_FILE_TYPE constant",
			constant: JPEG_FILE_TYPE,
			expected: ".jpeg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Constant value = %v, want %v", tt.constant, tt.expected)
			}
		})
	}
}

func TestTextWordsList(t *testing.T) {
	tests := []struct {
		name     string
		index    int
		expected string
	}{
		{
			name:     "First text word",
			index:    0,
			expected: "text_one_content",
		},
		{
			name:     "Second text word",
			index:    1,
			expected: "text_two_content",
		},
		{
			name:     "Third text word",
			index:    2,
			expected: "text_three_content",
		},
		{
			name:     "Fourth text word",
			index:    3,
			expected: "text_four_content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.index < 0 || tt.index >= len(textWordsList) {
				t.Errorf("Index %d out of bounds", tt.index)
				return
			}

			if textWordsList[tt.index] != tt.expected {
				t.Errorf("textWordsList[%d] = %v, want %v", tt.index, textWordsList[tt.index], tt.expected)
			}
		})
	}

	t.Run("TextWordsList length", func(t *testing.T) {
		expectedLength := 4
		if len(textWordsList) != expectedLength {
			t.Errorf("textWordsList length = %v, want %v", len(textWordsList), expectedLength)
		}
	})
}
