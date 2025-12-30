package engine

import (
	"testing"

	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/pkg"
)

func TestCopyExifData(t *testing.T) {
	tests := []struct {
		name string
		exif exiftool.FileMetadata
	}{
		{
			name: "Copy exif data with fields",
			exif: exiftool.FileMetadata{
				File: "test.jpg",
				Fields: map[string]any{
					"Make":  "Canon",
					"Model": "EOS 5D",
				},
				Err: nil,
			},
		},
		{
			name: "Copy exif data with error",
			exif: exiftool.FileMetadata{
				File:   "error.jpg",
				Fields: map[string]any{},
				Err:    pkg.NewErrors(1, "test error").Error,
			},
		},
		{
			name: "Copy empty exif data",
			exif: exiftool.FileMetadata{
				File:   "empty.jpg",
				Fields: map[string]any{},
				Err:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := copyExifData(tt.exif)

			if result.File != tt.exif.File {
				t.Errorf("copyExifData() File = %v, want %v", result.File, tt.exif.File)
			}

			if len(result.Fields) != len(tt.exif.Fields) {
				t.Errorf("copyExifData() Fields length = %v, want %v", len(result.Fields), len(tt.exif.Fields))
			}

			for k, v := range tt.exif.Fields {
				if result.Fields[k] != v {
					t.Errorf("copyExifData() Fields[%v] = %v, want %v", k, result.Fields[k], v)
				}
			}

			if result.Err == nil != (tt.exif.Err == nil) {
				t.Errorf("copyExifData() Err = %v, want %v", result.Err, tt.exif.Err)
			}
		})
	}
}

func TestCopyExifDataIndependence(t *testing.T) {
	tests := []struct {
		name string
		exif exiftool.FileMetadata
	}{
		{
			name: "Test copy independence",
			exif: exiftool.FileMetadata{
				File: "test.jpg",
				Fields: map[string]any{
					"Make":  "Canon",
					"Model": "EOS 5D",
				},
				Err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := copyExifData(tt.exif)

			result.Fields["Make"] = "Nikon"
			result.Fields["NewField"] = "NewValue"

			if tt.exif.Fields["Make"] == "Nikon" {
				t.Error("copyExifData() did not create independent copy - original was modified")
			}

			if _, exists := tt.exif.Fields["NewField"]; exists {
				t.Error("copyExifData() did not create independent copy - original has new field")
			}
		})
	}
}

func TestQuitAllTools(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Quit all tools",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			QuitAllTools()
		})
	}
}

func TestExiftoolInitFlag(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Check exiftool init flag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialFlag := exiftoolInitFlag

			if initialFlag {
				t.Log("exiftoolInitFlag is true")
			} else {
				t.Log("exiftoolInitFlag is false")
			}
		})
	}
}

func TestExiftoolCache(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Check exiftool cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := 0
			exiftoolCache.Range(func(key, value any) bool {
				count++
				return true
			})

			t.Logf("Exiftool cache has %d items", count)
		})
	}
}
