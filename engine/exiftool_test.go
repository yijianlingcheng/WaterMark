package engine

import (
	"testing"

	"WaterMark/pkg"
)

func TestInitExiftool(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Init exiftool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := initExiftool()

			if pkg.HasError(err) {
				t.Logf("initExiftool() returned error: %v", err)
				t.Skip("Exiftool not available in test environment")
			} else {
				t.Log("initExiftool() succeeded")
			}

			if !exiftoolInitFlag && !pkg.HasError(err) {
				t.Error("initExiftool() did not set exiftoolInitFlag to true")
			}
		})
	}
}

func TestCloseExiftool(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Close exiftool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			closeExiftool()
			t.Log("closeExiftool() completed")
		})
	}
}

func TestGetPhotosExifInfos(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Get photos exif infos",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !exiftoolInitFlag {
				t.Skip("Exiftool not initialized")
			}

			infos := GetPhotosExifInfos()
			if len(infos) != 0 {
				t.Errorf("GetPhotosExifInfos() with no paths returned %d items, want 0", len(infos))
			}
		})
	}
}

func TestGetPhotosExifInfo(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "Get photo exif info with non-existent file",
			path: "non_existent_file.jpg",
		},
		{
			name: "Get photo exif info with empty path",
			path: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !exiftoolInitFlag {
				t.Skip("Exiftool not initialized")
			}

			info, err := GetPhotosExifInfo(tt.path)

			if pkg.HasError(err) {
				t.Logf("GetPhotosExifInfo() returned expected error: %v", err)
			} else {
				t.Logf("GetPhotosExifInfo() returned info for %s", tt.path)
				if info.File != tt.path {
					t.Errorf("GetPhotosExifInfo() File = %v, want %v", info.File, tt.path)
				}
			}
		})
	}
}
