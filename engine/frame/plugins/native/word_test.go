package native

import (
	"testing"

	"github.com/yijianlingcheng/go-exiftool"
)

func TestChangeText2ExifContent(t *testing.T) {
	tests := []struct {
		name string
		exif exiftool.FileMetadata
		str  string
		want string
	}{
		{
			name: "Single field with valid EXIF data",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					"Make": "Canon",
				},
			},
			str:  "Make",
			want: "Canon",
		},
		{
			name: "Single field with missing EXIF data",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{},
			},
			str:  "Make",
			want: "Make",
		},
		{
			name: "Multiple fields separated by comma",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					"Make":  "Canon",
					"Model": "EOS 5D",
				},
			},
			str:  "Make,Model",
			want: "CanonEOS 5D",
		},
		{
			name: "Fields with # separator",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					"Make":  "Canon",
					"Model": "EOS 5D",
				},
			},
			str:  "Make#Model",
			want: "CanonEOS 5D",
		},
		{
			name: "Mixed separators",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					"Make":  "Canon",
					"Model": "EOS 5D",
					"ISO":   "100",
				},
			},
			str:  "Make#Model,ISO",
			want: "CanonEOS 5D100",
		},
		{
			name: "Empty string",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{},
			},
			str:  "",
			want: "",
		},
		{
			name: "Field with numeric value",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					"ISO": 100,
				},
			},
			str:  "ISO",
			want: "100",
		},
		{
			name: "Field with float value",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					FOCAL_LENGTH: 35.0,
				},
			},
			str:  FOCAL_LENGTH,
			want: "35",
		},
		{
			name: "Field with date time value",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					DATE_TIME_ORIGINAL: "2023:12:30 15:30:00",
				},
			},
			str:  DATE_TIME_ORIGINAL,
			want: "2023:12:30 15:30",
		},
		{
			name: "Field with GPS position",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					GPS_POSITION: "35 deg 41' 22.0\" N, 139 deg 41' 30.0\" E",
				},
			},
			str:  GPS_POSITION,
			want: "35°41′N 139°41′E",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := changeText2ExifContent(tt.exif, tt.str)
			if got != tt.want {
				t.Errorf("changeText2ExifContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChangeExifShowStr(t *testing.T) {
	tests := []struct {
		name string
		sub  string
		str  string
		exif exiftool.FileMetadata
		want string
	}{
		{
			name: "GPS_OR_DATETIME with GPS data",
			sub:  GPS_OR_DATETIME,
			str:  GPS_OR_DATETIME,
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					GPS_POSITION:       "35 deg 41' 22.0\" N, 139 deg 41' 30.0\" E",
					DATE_TIME_ORIGINAL: "2023:12:30 15:30:00",
				},
			},
			want: "35°41′N 139°41′E",
		},
		{
			name: "GPS_OR_DATETIME without GPS data",
			sub:  GPS_OR_DATETIME,
			str:  GPS_OR_DATETIME,
			exif: exiftool.FileMetadata{
				Fields: map[string]any{
					DATE_TIME_ORIGINAL: "2023:12:30 15:30:00",
				},
			},
			want: "2023:12:30 15:30:00",
		},
		{
			name: "GPS_OR_DATETIME with no data",
			sub:  GPS_OR_DATETIME,
			str:  GPS_OR_DATETIME,
			exif: exiftool.FileMetadata{
				Fields: map[string]any{},
			},
			want: GPS_OR_DATETIME,
		},
		{
			name: "GPS_POSITION with valid data",
			sub:  GPS_POSITION,
			str:  "35 deg 41' 22.0\" N, 139 deg 41' 30.0\" E",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{},
			},
			want: "35°41′N 139°41′E",
		},
		{
			name: "GPS_POSITION with no data (same as sub)",
			sub:  GPS_POSITION,
			str:  GPS_POSITION,
			exif: exiftool.FileMetadata{
				Fields: map[string]any{},
			},
			want: GPS_POSITION,
		},
		{
			name: "FOCAL_LENGTH with valid data",
			sub:  FOCAL_LENGTH,
			str:  "35.0 mm",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{},
			},
			want: "35mm",
		},
		{
			name: "FOCAL_LENGTH with no decimal",
			sub:  FOCAL_LENGTH,
			str:  "35 mm",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{},
			},
			want: "35mm",
		},
		{
			name: "DATE_TIME_ORIGINAL with :00 seconds",
			sub:  DATE_TIME_ORIGINAL,
			str:  "2023:12:30 15:30:00",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{},
			},
			want: "2023:12:30 15:30",
		},
		{
			name: "DATE_TIME_ORIGINAL without :00 seconds",
			sub:  DATE_TIME_ORIGINAL,
			str:  "2023:12:30 15:30:45",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{},
			},
			want: "2023:12:30 15:30:45",
		},
		{
			name: "Regular field without special handling",
			sub:  "Make",
			str:  "Canon",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{},
			},
			want: "Canon",
		},
		{
			name: "Empty string",
			sub:  "Make",
			str:  "",
			exif: exiftool.FileMetadata{
				Fields: map[string]any{},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := changeExifShowStr(tt.sub, tt.str, tt.exif)
			if got != tt.want {
				t.Errorf("changeExifShowStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChangeFocalLength(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "Focal length with decimal and space",
			str:  "35.0 mm",
			want: "35mm",
		},
		{
			name: "Focal length with decimal no space",
			str:  "35.0mm",
			want: "35mm",
		},
		{
			name: "Focal length without decimal",
			str:  "35 mm",
			want: "35mm",
		},
		{
			name: "Focal length with multiple decimals",
			str:  "35.5 mm",
			want: "35mm",
		},
		{
			name: "Focal length with three decimals",
			str:  "35.123 mm",
			want: "35mm",
		},
		{
			name: "Focal length without unit",
			str:  "35.0",
			want: "35mm",
		},
		{
			name: "Empty string",
			str:  "",
			want: "",
		},
		{
			name: "Focal length with leading zero",
			str:  "035.0 mm",
			want: "035mm",
		},
		{
			name: "Focal length with trailing space",
			str:  "35.0  mm",
			want: "35mm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := changeFocalLength(tt.str)
			if got != tt.want {
				t.Errorf("changeFocalLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChangeDataTime(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "Date time with :00 seconds",
			str:  "2023:12:30 15:30:00",
			want: "2023:12:30 15:30",
		},
		{
			name: "Date time with non-zero seconds",
			str:  "2023:12:30 15:30:45",
			want: "2023:12:30 15:30:45",
		},
		{
			name: "Date time with :01 seconds",
			str:  "2023:12:30 15:30:01",
			want: "2023:12:30 15:30:01",
		},
		{
			name: "Date time with :59 seconds",
			str:  "2023:12:30 15:30:59",
			want: "2023:12:30 15:30:59",
		},
		{
			name: "Date time with midnight",
			str:  "2023:12:30 00:00:00",
			want: "2023:12:30 00:00",
		},
		{
			name: "Date time with :00 in minutes",
			str:  "2023:12:30 15:00:00",
			want: "2023:12:30 15:00",
		},
		{
			name: "Date time without seconds",
			str:  "2023:12:30 15:30",
			want: "2023:12:30 15:30",
		},
		{
			name: "Date time with :00 but different format",
			str:  "2023-12-30 15:30:00",
			want: "2023-12-30 15:30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := changeDataTime(tt.str)
			if got != tt.want {
				t.Errorf("changeDataTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
