package pkg

import "testing"

func TestGetOrientation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Horizontal orientation",
			input:    "Horizontal (normal)",
			expected: 0,
		},
		{
			name:     "normal orientation",
			input:    "normal",
			expected: 0,
		},
		{
			name:     "Rotate 90 CCW",
			input:    "Rotate 90 CCW",
			expected: 90,
		},
		{
			name:     "Rotate 180 CCW",
			input:    "Rotate 180 CCW",
			expected: 180,
		},
		{
			name:     "Rotate 270 CCW",
			input:    "Rotate 270 CCW",
			expected: 270,
		},
		{
			name:     "Rotate 90 CW",
			input:    "Rotate 90 CW",
			expected: 270,
		},
		{
			name:     "Rotate 180 CW",
			input:    "Rotate 180 CW",
			expected: 180,
		},
		{
			name:     "Rotate 270 CW",
			input:    "Rotate 270 CW",
			expected: 90,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "no rotation keyword",
			input:    "Some random text",
			expected: 0,
		},
		{
			name:     "Rotate without numbers",
			input:    "Rotate CCW",
			expected: 0,
		},
		{
			name:     "Rotate 360 CCW",
			input:    "Rotate 360 CCW",
			expected: 360,
		},
		{
			name:     "Rotate 360 CW",
			input:    "Rotate 360 CW",
			expected: 0,
		},
		{
			name:     "Rotate 45 CCW",
			input:    "Rotate 45 CCW",
			expected: 45,
		},
		{
			name:     "Rotate 45 CW",
			input:    "Rotate 45 CW",
			expected: 315,
		},
		{
			name:     "Horizontal with extra text",
			input:    "Horizontal (normal) - extra info",
			expected: 0,
		},
		{
			name:     "Rotate 90 with extra spaces",
			input:    "Rotate   90   CCW",
			expected: 90,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetOrientation(tt.input)
			if result != tt.expected {
				t.Errorf("GetOrientation() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "single number",
			input:    "123",
			expected: []int{123},
		},
		{
			name:     "multiple numbers",
			input:    "123 456 789",
			expected: []int{123, 456, 789},
		},
		{
			name:     "numbers in text",
			input:    "Rotate 90 CCW",
			expected: []int{90},
		},
		{
			name:     "no numbers",
			input:    "hello world",
			expected: []int{},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []int{},
		},
		{
			name:     "large number",
			input:    "999999",
			expected: []int{999999},
		},
		{
			name:     "zero",
			input:    "0",
			expected: []int{0},
		},
		{
			name:     "negative numbers as text",
			input:    "-123",
			expected: []int{123},
		},
		{
			name:     "decimal numbers",
			input:    "123.456",
			expected: []int{123, 456},
		},
		{
			name:     "mixed content",
			input:    "abc123def456ghi789",
			expected: []int{123, 456, 789},
		},
		{
			name:     "leading zeros",
			input:    "00123",
			expected: []int{123},
		},
		{
			name:     "multiple same numbers",
			input:    "90 90 90",
			expected: []int{90, 90, 90},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractNumbers(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("extractNumbers() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("extractNumbers()[%d] = %v, want %v", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestGetOrientationEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "CCW without Rotate",
			input:    "CCW 90",
			expected: 0,
		},
		{
			name:     "CW without Rotate",
			input:    "CW 90",
			expected: 0,
		},
		{
			name:     "Rotate with CCW and CW",
			input:    "Rotate 90 CCW CW",
			expected: 90,
		},
		{
			name:     "Rotate 720 CCW",
			input:    "Rotate 720 CCW",
			expected: 720,
		},
		{
			name:     "Rotate 720 CW",
			input:    "Rotate 720 CW",
			expected: -360,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetOrientation(tt.input)
			if result != tt.expected {
				t.Errorf("GetOrientation() = %v, want %v", result, tt.expected)
			}
		})
	}
}
