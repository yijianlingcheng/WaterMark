package pkg

import (
	"strconv"
	"testing"
)

func TestAnyToString(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: "",
		},
		{
			name:     "string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "int",
			input:    42,
			expected: "42",
		},
		{
			name:     "int negative",
			input:    -42,
			expected: "-42",
		},
		{
			name:     "int zero",
			input:    0,
			expected: "0",
		},
		{
			name:     "int8",
			input:    int8(127),
			expected: "127",
		},
		{
			name:     "int8 negative",
			input:    int8(-128),
			expected: "-128",
		},
		{
			name:     "int16",
			input:    int16(32767),
			expected: "32767",
		},
		{
			name:     "int32",
			input:    int32(2147483647),
			expected: "2147483647",
		},
		{
			name:     "int64",
			input:    int64(9223372036854775807),
			expected: "9223372036854775807",
		},
		{
			name:     "uint",
			input:    uint(42),
			expected: "42",
		},
		{
			name:     "uint zero",
			input:    uint(0),
			expected: "0",
		},
		{
			name:     "uint8",
			input:    uint8(255),
			expected: "255",
		},
		{
			name:     "uint16",
			input:    uint16(65535),
			expected: "65535",
		},
		{
			name:     "uint32",
			input:    uint32(4294967295),
			expected: "4294967295",
		},
		{
			name:     "uint64",
			input:    uint64(18446744073709551615),
			expected: "18446744073709551615",
		},
		{
			name:     "float32",
			input:    float32(3.14),
			expected: "3.14",
		},
		{
			name:     "float32 negative",
			input:    float32(-2.5),
			expected: "-2.5",
		},
		{
			name:     "float32 zero",
			input:    float32(0.0),
			expected: "0",
		},
		{
			name:     "float64",
			input:    3.1415926535,
			expected: "3.1415926535",
		},
		{
			name:     "float64 negative",
			input:    -2.718281828,
			expected: "-2.718281828",
		},
		{
			name:     "float64 zero",
			input:    0.0,
			expected: "0",
		},
		{
			name:     "[]byte",
			input:    []byte("hello"),
			expected: "hello",
		},
		{
			name:     "[]byte empty",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "[]byte with special chars",
			input:    []byte("你好世界"),
			expected: "你好世界",
		},
		{
			name:     "struct",
			input:    struct{ Name string }{Name: "test"},
			expected: `{"Name":"test"}`,
		},
		{
			name:     "map",
			input:    map[string]int{"key": 123},
			expected: `{"key":123}`,
		},
		{
			name:     "slice",
			input:    []int{1, 2, 3},
			expected: `[1,2,3]`,
		},
		{
			name:     "bool",
			input:    true,
			expected: "true",
		},
		{
			name:     "bool false",
			input:    false,
			expected: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AnyToString(tt.input)
			if result != tt.expected {
				t.Errorf("AnyToString(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAnyToStringLargeNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "large int64",
			input:    int64(1234567890123456789),
			expected: "1234567890123456789",
		},
		{
			name:     "large uint64",
			input:    uint64(9876543210987654321),
			expected: "9876543210987654321",
		},
		{
			name:     "small float",
			input:    float64(0.000001),
			expected: "0.000001",
		},
		{
			name:     "large float",
			input:    float64(123456789.123456789),
			expected: "123456789.12345679",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AnyToString(tt.input)
			if result != tt.expected {
				t.Errorf("AnyToString(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIn(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		array    []string
		expected bool
	}{
		{
			name:     "target exists in array",
			target:   "apple",
			array:    []string{"banana", "apple", "cherry"},
			expected: true,
		},
		{
			name:     "target does not exist in array",
			target:   "orange",
			array:    []string{"banana", "apple", "cherry"},
			expected: false,
		},
		{
			name:     "empty array",
			target:   "apple",
			array:    []string{},
			expected: false,
		},
		{
			name:     "empty string in array",
			target:   "",
			array:    []string{"", "apple", "banana"},
			expected: true,
		},
		{
			name:     "empty string not in array",
			target:   "",
			array:    []string{"apple", "banana"},
			expected: false,
		},
		{
			name:     "empty string target in empty array",
			target:   "",
			array:    []string{},
			expected: false,
		},
		{
			name:     "target at beginning",
			target:   "apple",
			array:    []string{"apple", "banana", "cherry"},
			expected: true,
		},
		{
			name:     "target at end",
			target:   "cherry",
			array:    []string{"apple", "banana", "cherry"},
			expected: true,
		},
		{
			name:     "target in middle",
			target:   "banana",
			array:    []string{"apple", "banana", "cherry"},
			expected: true,
		},
		{
			name:     "case sensitive - lowercase",
			target:   "apple",
			array:    []string{"Apple", "APPLE", "apple"},
			expected: true,
		},
		{
			name:     "case sensitive - uppercase",
			target:   "APPLE",
			array:    []string{"Apple", "apple", "APPLE"},
			expected: true,
		},
		{
			name:     "case sensitive - mixed",
			target:   "Apple",
			array:    []string{"apple", "APPLE", "Apple"},
			expected: true,
		},
		{
			name:     "duplicate values",
			target:   "apple",
			array:    []string{"apple", "apple", "banana"},
			expected: true,
		},
		{
			name:     "special characters",
			target:   "hello@world.com",
			array:    []string{"test@example.com", "hello@world.com", "user@test.com"},
			expected: true,
		},
		{
			name:     "unicode characters",
			target:   "你好",
			array:    []string{"hello", "你好", "world"},
			expected: true,
		},
		{
			name:     "single element array - exists",
			target:   "apple",
			array:    []string{"apple"},
			expected: true,
		},
		{
			name:     "single element array - not exists",
			target:   "banana",
			array:    []string{"apple"},
			expected: false,
		},
		{
			name:     "numbers as strings",
			target:   "123",
			array:    []string{"123", "456", "789"},
			expected: true,
		},
		{
			name:     "mixed content",
			target:   "test",
			array:    []string{"123", "test", "@#$", "你好"},
			expected: true,
		},
		{
			name:     "whitespace",
			target:   "hello world",
			array:    []string{"hello", "hello world", "helloworld"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := In(tt.target, tt.array)
			if result != tt.expected {
				t.Errorf("In(%q, %v) = %v, want %v", tt.target, tt.array, result, tt.expected)
			}
		})
	}
}

func TestInUnsortedArray(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		array    []string
		expected bool
	}{
		{
			name:     "unsorted array - exists",
			target:   "cherry",
			array:    []string{"banana", "cherry", "apple"},
			expected: true,
		},
		{
			name:     "unsorted array - not exists",
			target:   "orange",
			array:    []string{"banana", "cherry", "apple"},
			expected: false,
		},
		{
			name:     "reverse sorted array - exists",
			target:   "apple",
			array:    []string{"cherry", "banana", "apple"},
			expected: true,
		},
		{
			name:     "random order array - exists",
			target:   "banana",
			array:    []string{"zebra", "banana", "apple", "yellow"},
			expected: true,
		},
		{
			name:     "random order array - not exists",
			target:   "orange",
			array:    []string{"zebra", "banana", "apple", "yellow"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := In(tt.target, tt.array)
			if result != tt.expected {
				t.Errorf("In(%q, %v) = %v, want %v", tt.target, tt.array, result, tt.expected)
			}
		})
	}
}

func TestInLargeArray(t *testing.T) {
	largeArray := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		largeArray[i] = "item" + strconv.Itoa(i)
	}

	tests := []struct {
		name     string
		target   string
		array    []string
		expected bool
	}{
		{
			name:     "large array - exists at beginning",
			target:   "item0",
			array:    largeArray,
			expected: true,
		},
		{
			name:     "large array - exists at end",
			target:   "item999",
			array:    largeArray,
			expected: true,
		},
		{
			name:     "large array - exists in middle",
			target:   "item500",
			array:    largeArray,
			expected: true,
		},
		{
			name:     "large array - not exists",
			target:   "item1000",
			array:    largeArray,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := In(tt.target, tt.array)
			if result != tt.expected {
				t.Errorf("In(%q, array of size %d) = %v, want %v", tt.target, len(tt.array), result, tt.expected)
			}
		})
	}
}
