package layout

import (
	"testing"
)

func TestGpsFormat(t *testing.T) {
	tests := []struct {
		name     string
		gps      string
		expected string
	}{
		{
			name:     "Valid GPS coordinates",
			gps:      "30 deg 15.5' 30.2\" N, 120 deg 45.2' 15.8\" E",
			expected: "30°15′N 120°45′E",
		},
		{
			name:     "Valid GPS coordinates with decimal degrees",
			gps:      "45.5 deg 30.0' 0.0\" N, 90.0 deg 15.0' 0.0\" W",
			expected: "45°30′N 90°15′W",
		},
		{
			name:     "Valid GPS coordinates South and West",
			gps:      "15 deg 20.0' 10.0\" S, 75 deg 30.0' 20.0\" W",
			expected: "15°20′S 75°30′W",
		},
		{
			name:     "Empty GPS string",
			gps:      "",
			expected: "",
		},
		{
			name:     "Invalid GPS format",
			gps:      "invalid gps format",
			expected: " ",
		},
		{
			name:     "GPS with only latitude",
			gps:      "30 deg 15.5' 30.2\" N",
			expected: "30°15′N ",
		},
		{
			name:     "GPS with extra spaces",
			gps:      "  30  deg  15.5'  30.2\"  N  ,  120  deg  45.2'  15.8\"  E  ",
			expected: "30°15′N 120°45′E",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GpsFormat(tt.gps)
			if result != tt.expected {
				t.Errorf("GpsFormat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetGPSOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		gps      string
		other    string
		expected string
	}{
		{
			name:     "Valid GPS returns formatted GPS",
			gps:      "30 deg 15.5' 30.2\" N, 120 deg 45.2' 15.8\" E",
			other:    "Default Location",
			expected: "30°15′N 120°45′E",
		},
		{
			name:     "Empty GPS returns default value",
			gps:      "",
			other:    "Default Location",
			expected: "Default Location",
		},
		{
			name:     "Invalid GPS returns default value",
			gps:      "invalid gps",
			other:    "Default Location",
			expected: "Default Location",
		},
		{
			name:     "Empty GPS and empty default returns empty",
			gps:      "",
			other:    "",
			expected: "",
		},
		{
			name:     "Valid GPS with empty default returns formatted GPS",
			gps:      "45 deg 30' 0\" N, 90 deg 15' 0\" W",
			other:    "",
			expected: "45°30′N 90°15′W",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetGPSOrDefault(tt.gps, tt.other)
			if result != tt.expected {
				t.Errorf("GetGPSOrDefault() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseGPSInfo(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected string
	}{
		{
			name:     "Valid latitude with N direction",
			str:      "30 deg 15.5' 30.2\" N",
			expected: "30°15′N",
		},
		{
			name:     "Valid longitude with E direction",
			str:      "120 deg 45.2' 15.8\" E",
			expected: "120°45′E",
		},
		{
			name:     "Valid coordinate with S direction",
			str:      "15 deg 20.0' 10.0\" S",
			expected: "15°20′S",
		},
		{
			name:     "Valid coordinate with W direction",
			str:      "75 deg 30.0' 20.0\" W",
			expected: "75°30′W",
		},
		{
			name:     "Empty string",
			str:      "",
			expected: "",
		},
		{
			name:     "Invalid format - missing direction",
			str:      "30 deg 15.5' 30.2\"",
			expected: "",
		},
		{
			name:     "Invalid format - wrong units",
			str:      "30 15.5 30.2 N",
			expected: "",
		},
		{
			name:     "Invalid format - missing components",
			str:      "30 deg N",
			expected: "",
		},
		{
			name:     "Valid with zero minutes",
			str:      "30 deg 0' 0\" N",
			expected: "30°00′N",
		},
		{
			name:     "Valid with decimal degrees",
			str:      "30.5 deg 15.5' 30.2\" N",
			expected: "30°15′N",
		},
		{
			name:     "Valid with decimal minutes",
			str:      "30 deg 15.75' 30.2\" N",
			expected: "30°15′N",
		},
		{
			name:     "Valid with decimal seconds",
			str:      "30 deg 15' 30.75\" N",
			expected: "30°15′N",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseGPSInfo(tt.str)
			if result != tt.expected {
				t.Errorf("parseGPSInfo() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGpsFormatEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		gps      string
		expected string
	}{
		{
			name:     "GPS with very small coordinates",
			gps:      "0 deg 0' 0.1\" N, 0 deg 0' 0.1\" E",
			expected: "0°00′N 0°00′E",
		},
		{
			name:     "GPS with large coordinates",
			gps:      "179 deg 59' 59.9\" N, 179 deg 59' 59.9\" E",
			expected: "179°59′N 179°59′E",
		},
		{
			name:     "GPS with single digit degrees",
			gps:      "5 deg 5' 5\" N, 5 deg 5' 5\" E",
			expected: "5°05′N 5°05′E",
		},
		{
			name:     "GPS with three digit degrees",
			gps:      "100 deg 30' 0\" N, 100 deg 30' 0\" E",
			expected: "100°30′N 100°30′E",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GpsFormat(tt.gps)
			if result != tt.expected {
				t.Errorf("GpsFormat() = %v, want %v", result, tt.expected)
			}
		})
	}
}
