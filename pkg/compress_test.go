package pkg

import (
	"bytes"
	"testing"
)

func TestZlibCompress(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		wantEmpty bool
	}{
		{
			name:      "simple string",
			input:     []byte("hello world"),
			wantEmpty: false,
		},
		{
			name:      "empty byte slice",
			input:     []byte{},
			wantEmpty: false,
		},
		{
			name:      "single character",
			input:     []byte("a"),
			wantEmpty: false,
		},
		{
			name: "long string",
			input: []byte(
				"This is a longer string that will be compressed to demonstrate the zlib compression functionality.",
			),
			wantEmpty: false,
		},
		{
			name:      "special characters",
			input:     []byte("!@#$%^&*()_+-=[]{}|;':\",./<>?"),
			wantEmpty: false,
		},
		{
			name:      "unicode characters",
			input:     []byte("你好世界こんにちは안녕하세요"),
			wantEmpty: false,
		},
		{
			name:      "numbers",
			input:     []byte("1234567890"),
			wantEmpty: false,
		},
		{
			name:      "whitespace",
			input:     []byte("   \t\n\r  "),
			wantEmpty: false,
		},
		{
			name:      "repeated pattern",
			input:     []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			wantEmpty: false,
		},
		{
			name:      "mixed content",
			input:     []byte("Hello 你好 123 !@#"),
			wantEmpty: false,
		},
		{
			name:      "binary data",
			input:     []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD},
			wantEmpty: false,
		},
		{
			name:      "large data",
			input:     bytes.Repeat([]byte("test data "), 1000),
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ZlibCompress(tt.input)

			if tt.wantEmpty {
				if len(result) != 0 {
					t.Errorf("ZlibCompress() expected empty result, got %d bytes", len(result))
				}
			} else {
				if len(result) == 0 {
					t.Errorf("ZlibCompress() returned empty result for non-empty input")
				}
			}
		})
	}
}

func TestZlibUnCompress(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "valid compressed data",
			input:    ZlibCompress([]byte("hello world")),
			expected: []byte("hello world"),
		},
		{
			name:     "empty compressed data",
			input:    ZlibCompress([]byte{}),
			expected: []byte{},
		},
		{
			name:     "single character",
			input:    ZlibCompress([]byte("a")),
			expected: []byte("a"),
		},
		{
			name: "long string",
			input: ZlibCompress(
				[]byte(
					"This is a longer string that will be compressed to demonstrate the zlib compression functionality.",
				),
			),
			expected: []byte(
				"This is a longer string that will be compressed to demonstrate the zlib compression functionality.",
			),
		},
		{
			name:     "special characters",
			input:    ZlibCompress([]byte("!@#$%^&*()_+-=[]{}|;':\",./<>?")),
			expected: []byte("!@#$%^&*()_+-=[]{}|;':\",./<>?"),
		},
		{
			name:     "unicode characters",
			input:    ZlibCompress([]byte("你好世界こんにちは안녕하세요")),
			expected: []byte("你好世界こんにちは안녕하세요"),
		},
		{
			name:     "numbers",
			input:    ZlibCompress([]byte("1234567890")),
			expected: []byte("1234567890"),
		},
		{
			name:     "whitespace",
			input:    ZlibCompress([]byte("   \t\n\r  ")),
			expected: []byte("   \t\n\r  "),
		},
		{
			name:     "repeated pattern",
			input:    ZlibCompress([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")),
			expected: []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		},
		{
			name:     "mixed content",
			input:    ZlibCompress([]byte("Hello 你好 123 !@#")),
			expected: []byte("Hello 你好 123 !@#"),
		},
		{
			name:     "binary data",
			input:    ZlibCompress([]byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}),
			expected: []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD},
		},
		{
			name:     "large data",
			input:    ZlibCompress(bytes.Repeat([]byte("test data "), 1000)),
			expected: bytes.Repeat([]byte("test data "), 1000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ZlibUnCompress(tt.input)

			if !bytes.Equal(result, tt.expected) {
				t.Errorf("ZlibUnCompress() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestZlibCompressUncompressRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "simple string",
			input: []byte("hello world"),
		},
		{
			name:  "empty byte slice",
			input: []byte{},
		},
		{
			name:  "single character",
			input: []byte("a"),
		},
		{
			name: "long string",
			input: []byte(
				"This is a longer string that will be compressed to demonstrate the zlib compression functionality.",
			),
		},
		{
			name:  "special characters",
			input: []byte("!@#$%^&*()_+-=[]{}|;':\",./<>?"),
		},
		{
			name:  "unicode characters",
			input: []byte("你好世界こんにちは안녕하세요"),
		},
		{
			name:  "numbers",
			input: []byte("1234567890"),
		},
		{
			name:  "whitespace",
			input: []byte("   \t\n\r  "),
		},
		{
			name:  "repeated pattern",
			input: []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		},
		{
			name:  "mixed content",
			input: []byte("Hello 你好 123 !@#"),
		},
		{
			name:  "binary data",
			input: []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD},
		},
		{
			name:  "large data",
			input: bytes.Repeat([]byte("test data "), 1000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed := ZlibCompress(tt.input)
			decompressed := ZlibUnCompress(compressed)

			if !bytes.Equal(decompressed, tt.input) {
				t.Errorf("Round trip failed: original = %v, decompressed = %v", tt.input, decompressed)
			}
		})
	}
}

func TestZlibUnCompressInvalidData(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "plain text",
			input: []byte("hello world"),
		},
		{
			name:  "invalid zlib header",
			input: []byte{0x78, 0x9C, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name:  "random bytes",
			input: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		},
		{
			name:  "empty byte slice",
			input: []byte{},
		},
		{
			name:  "partial zlib data",
			input: []byte{0x78, 0x9C},
		},
		{
			name:  "corrupted zlib data",
			input: []byte{0x78, 0xDA, 0x01, 0x02, 0x03},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ZlibUnCompress(tt.input)

			if len(result) != 0 {
				t.Errorf("ZlibUnCompress() with invalid data should return empty, got %d bytes", len(result))
			}
		})
	}
}

func TestZlibCompressionRatio(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		expectSmaller bool
	}{
		{
			name:          "repeated pattern (highly compressible)",
			input:         bytes.Repeat([]byte("a"), 1000),
			expectSmaller: true,
		},
		{
			name:          "random-like data (less compressible)",
			input:         []byte("This is a test string with some variety."),
			expectSmaller: false,
		},
		{
			name:          "large repeated text",
			input:         bytes.Repeat([]byte("test data "), 100),
			expectSmaller: true,
		},
		{
			name:          "short string",
			input:         []byte("hello"),
			expectSmaller: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed := ZlibCompress(tt.input)

			if tt.expectSmaller {
				if len(compressed) >= len(tt.input) {
					t.Logf("Warning: Compressed size %d is not smaller than original %d for %s",
						len(compressed), len(tt.input), tt.name)
				}
			}

			if len(compressed) == 0 && len(tt.input) > 0 {
				t.Errorf("Compression failed for non-empty input")
			}
		})
	}
}

func TestZlibCompressMultipleTimes(t *testing.T) {
	input := []byte("hello world")

	compressed1 := ZlibCompress(input)
	compressed2 := ZlibCompress(input)

	if !bytes.Equal(compressed1, compressed2) {
		t.Error("ZlibCompress() should produce consistent results for the same input")
	}
}

func TestZlibUnCompressMultipleTimes(t *testing.T) {
	input := ZlibCompress([]byte("hello world"))

	decompressed1 := ZlibUnCompress(input)
	decompressed2 := ZlibUnCompress(input)

	if !bytes.Equal(decompressed1, decompressed2) {
		t.Error("ZlibUnCompress() should produce consistent results for the same input")
	}
}

func TestZlibCompressWithNullBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "single null byte",
			input: []byte{0x00},
		},
		{
			name:  "multiple null bytes",
			input: []byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			name:  "mixed with null bytes",
			input: []byte{'a', 0x00, 'b', 0x00, 'c'},
		},
		{
			name:  "null bytes at start",
			input: append([]byte{0x00, 0x00}, []byte("hello")...),
		},
		{
			name:  "null bytes at end",
			input: append([]byte("hello"), []byte{0x00, 0x00}...),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed := ZlibCompress(tt.input)
			decompressed := ZlibUnCompress(compressed)

			if !bytes.Equal(decompressed, tt.input) {
				t.Errorf("Round trip with null bytes failed: original = %v, decompressed = %v",
					tt.input, decompressed)
			}
		})
	}
}

func TestZlibCompressLargeFile(t *testing.T) {
	largeData := bytes.Repeat([]byte("This is a test line that will be repeated many times.\n"), 10000)

	compressed := ZlibCompress(largeData)
	decompressed := ZlibUnCompress(compressed)

	if !bytes.Equal(decompressed, largeData) {
		t.Error("Large file compression/decompression failed")
	}

	if len(compressed) == 0 {
		t.Error("Compressed data is empty")
	}

	if len(compressed) >= len(largeData) {
		t.Logf("Warning: Large file compression did not reduce size (original: %d, compressed: %d)",
			len(largeData), len(compressed))
	}
}

func TestZlibCompressAllByteValues(t *testing.T) {
	allBytes := make([]byte, 256)
	for i := 0; i < 256; i++ {
		allBytes[i] = byte(i)
	}

	compressed := ZlibCompress(allBytes)
	decompressed := ZlibUnCompress(compressed)

	if !bytes.Equal(decompressed, allBytes) {
		t.Error("All byte values round trip failed")
	}
}
