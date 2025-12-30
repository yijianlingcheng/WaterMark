package pkg

import "testing"

func TestECodeConstants(t *testing.T) {
	tests := []struct {
		name  string
		value int
	}{
		{"NO_ERROR", NO_ERROR},
		{"FILE_NOT_EXIST_ERROR", FILE_NOT_EXIST_ERROR},
		{"FILE_NOT_OPEN_ERROR", FILE_NOT_OPEN_ERROR},
		{"FILE_NOT_READ_ERROR", FILE_NOT_READ_ERROR},
		{"EXIFTOOL_NOTEXIST_ERROR", EXIFTOOL_NOTEXIST_ERROR},
		{"EXIFTOOL_INIT_ERROR", EXIFTOOL_INIT_ERROR},
		{"EXIFTOOL_IMAGE_EXIF_ERROR", EXIFTOOL_IMAGE_EXIF_ERROR},
		{"EXIFTOOL_IMAGE_EXIF_CACHE_ERROR", EXIFTOOL_IMAGE_EXIF_CACHE_ERROR},
		{"CSV_CREATE_ERROR", CSV_CREATE_ERROR},
		{"CSV_WRITE_HEADER_ERROR", CSV_WRITE_HEADER_ERROR},
		{"CSV_WRITE_DATA_ERROR", CSV_WRITE_DATA_ERROR},
		{"IMAGE_DECODE_ERROR", IMAGE_DECODE_ERROR},
		{"IMAGE_NO_SUPPORT_ERROR", IMAGE_NO_SUPPORT_ERROR},
		{"IMAGE_DECODE_CACHE_ERROR", IMAGE_DECODE_CACHE_ERROR},
		{"IMAGE_RGBA_CACHE_ERROR", IMAGE_RGBA_CACHE_ERROR},
		{"IMAGE_TEXT_FONT_CACHE_ERROR", IMAGE_TEXT_FONT_CACHE_ERROR},
		{"IMAGE_TEXT_DRAW_TXT_ERROR", IMAGE_TEXT_DRAW_TXT_ERROR},
		{"IMAGE_LOGO_NOT_FIND_ERROR", IMAGE_LOGO_NOT_FIND_ERROR},
		{"IMAGE_LOGO_RESIZE_ERROR", IMAGE_LOGO_RESIZE_ERROR},
		{"IMAGE_JPEG_SAVE_ERROR", IMAGE_JPEG_SAVE_ERROR},
		{"CMD_COMMAND_RUN_ERROR", CMD_COMMAND_RUN_ERROR},
		{"LAYOUT_TYPE_NOT_FIND_ERROR", LAYOUT_TYPE_NOT_FIND_ERROR},
		{"INTERNAL_ERROR", INTERNAL_ERROR},
		{"REQUEST_PARAM_ERROR", REQUEST_PARAM_ERROR},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value <= 0 && tt.name != "NO_ERROR" {
				t.Errorf("Error code %s should be positive, got %d", tt.name, tt.value)
			}
		})
	}
}

func TestErrorCodeRanges(t *testing.T) {
	tests := []struct {
		name     string
		minCode  int
		maxCode  int
		expected int
	}{
		{"File errors", 1000000, 1999999, FILE_NOT_EXIST_ERROR},
		{"Exiftool errors", 2000000, 2999999, EXIFTOOL_NOTEXIST_ERROR},
		{"CSV errors", 3000000, 3999999, CSV_CREATE_ERROR},
		{"Image errors", 4000000, 4999999, IMAGE_DECODE_ERROR},
		{"Command errors", 5000000, 5999999, CMD_COMMAND_RUN_ERROR},
		{"Layout errors", 6000000, 6999999, LAYOUT_TYPE_NOT_FIND_ERROR},
		{"Internal errors", 9000000, 9999999, INTERNAL_ERROR},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expected < tt.minCode || tt.expected > tt.maxCode {
				t.Errorf("Error code %d is not in expected range [%d, %d]", tt.expected, tt.minCode, tt.maxCode)
			}
		})
	}
}
