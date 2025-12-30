package pkg

import (
	"encoding/json"
	"testing"
)

func TestIsOk(t *testing.T) {
	tests := []struct {
		name     string
		err      EError
		expected bool
	}{
		{
			name:     "NoError should return true",
			err:      NoError,
			expected: true,
		},
		{
			name:     "InternalError should return false",
			err:      InternalError,
			expected: false,
		},
		{
			name:     "ExiftoolNotExistError should return false",
			err:      ExiftoolNotExistError,
			expected: false,
		},
		{
			name:     "Custom error with code 0 should return true",
			err:      EError{Code: 0, Error: nil},
			expected: true,
		},
		{
			name:     "Custom error with code 1 should return false",
			err:      EError{Code: 1, Error: nil},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsOk(tt.err)
			if result != tt.expected {
				t.Errorf("IsOk() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestHasError(t *testing.T) {
	tests := []struct {
		name     string
		err      EError
		expected bool
	}{
		{
			name:     "NoError should return false",
			err:      NoError,
			expected: false,
		},
		{
			name:     "InternalError should return true",
			err:      InternalError,
			expected: true,
		},
		{
			name:     "ExiftoolNotExistError should return true",
			err:      ExiftoolNotExistError,
			expected: true,
		},
		{
			name:     "Custom error with code 0 should return false",
			err:      EError{Code: 0, Error: nil},
			expected: false,
		},
		{
			name:     "Custom error with code 1 should return true",
			err:      EError{Code: 1, Error: nil},
			expected: true,
		},
		{
			name:     "Custom error with negative code should return false",
			err:      EError{Code: -1, Error: nil},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasError(tt.err)
			if result != tt.expected {
				t.Errorf("HasError() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestEErrorString(t *testing.T) {
	tests := []struct {
		name     string
		err      EError
		expected string
	}{
		{
			name:     "NoError string representation",
			err:      NoError,
			expected: `{"Code":0,"Error":"no error"}`,
		},
		{
			name:     "InternalError string representation",
			err:      InternalError,
			expected: `{"Code":9000001,"Error":"发生内部错误"}`,
		},
		{
			name:     "ExiftoolNotExistError string representation",
			err:      ExiftoolNotExistError,
			expected: `{"Code":2000001,"Error":"exiftool工具不存在,请检查是否安装"}`,
		},
		{
			name:     "Custom error string representation",
			err:      NewErrors(123, "test error"),
			expected: `{"Code":123,"Error":"test error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.String()
			if result != tt.expected {
				t.Errorf("String() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestEErrorMarshalJSON(t *testing.T) {
	tests := []struct {
		name           string
		err            EError
		expectedCode   int
		expectedMsg    string
	}{
		{
			name:         "NoError JSON marshaling",
			err:          NoError,
			expectedCode: 0,
			expectedMsg:  "no error",
		},
		{
			name:         "InternalError JSON marshaling",
			err:          InternalError,
			expectedCode: 9000001,
			expectedMsg:  "发生内部错误",
		},
		{
			name:         "ExiftoolNotExistError JSON marshaling",
			err:          ExiftoolNotExistError,
			expectedCode: 2000001,
			expectedMsg:  "exiftool工具不存在,请检查是否安装",
		},
		{
			name:         "ExiftoolInitError JSON marshaling",
			err:          ExiftoolInitError,
			expectedCode: 2000002,
			expectedMsg:  "exiftool工具init失败,请检查环境",
		},
		{
			name:         "ExiftoolImageError JSON marshaling",
			err:          ExiftoolImageError,
			expectedCode: 2000003,
			expectedMsg:  "exiftool工具获取图片exif信息失败",
		},
		{
			name:         "ExiftoolCacheTypeError JSON marshaling",
			err:          ExiftoolCacheTypeError,
			expectedCode: 2000004,
			expectedMsg:  "exif cache缓存数据类型断言失败",
		},
		{
			name:         "ImageDecodeCacheTypeError JSON marshaling",
			err:          ImageDecodeCacheTypeError,
			expectedCode: 4000003,
			expectedMsg:  "缓存的图片解码数据类型断言失败",
		},
		{
			name:         "ImageRGBACacheTypeError JSON marshaling",
			err:          ImageRGBACacheTypeError,
			expectedCode: 4000004,
			expectedMsg:  "缓存的RGBA数据类型断言失败",
		},
		{
			name:         "ImageTextCacheTypeError JSON marshaling",
			err:          ImageTextCacheTypeError,
			expectedCode: 4000005,
			expectedMsg:  "缓存中获取的字体对象类型断言失败",
		},
		{
			name:         "ImageLogoNotFindError JSON marshaling",
			err:          ImageLogoNotFindError,
			expectedCode: 4000007,
			expectedMsg:  "logo图片文件没有找到",
		},
		{
			name:         "ImageLogoResizeError JSON marshaling",
			err:          ImageLogoResizeError,
			expectedCode: 4000008,
			expectedMsg:  "logo图片重置尺寸错误",
		},
		{
			name:         "ImageJpegSaveError JSON marshaling",
			err:          ImageJpegSaveError,
			expectedCode: 4000009,
			expectedMsg:  "jpeg图片保存失败",
		},
		{
			name:         "LayoutNotFindError JSON marshaling",
			err:          LayoutNotFindError,
			expectedCode: 6000001,
			expectedMsg:  "布局类型查找失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.err)
			if err != nil {
				t.Fatalf("MarshalJSON() failed: %v", err)
			}

			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(result, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			if code, ok := unmarshaled["code"].(float64); !ok || int(code) != tt.expectedCode {
				t.Errorf("JSON code = %v, expected %v", code, tt.expectedCode)
			}

			if msg, ok := unmarshaled["errmsg"].(string); !ok || msg != tt.expectedMsg {
				t.Errorf("JSON errmsg = %v, expected %v", msg, tt.expectedMsg)
			}
		})
	}
}

func TestNewErrors(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		msg      string
		expected EError
	}{
		{
			name:     "Create error with code 0",
			code:     0,
			msg:      "test error",
			expected: EError{Code: 0},
		},
		{
			name:     "Create error with code 100",
			code:     100,
			msg:      "test error",
			expected: EError{Code: 100},
		},
		{
			name:     "Create error with code 2000001",
			code:     EXIFTOOL_NOTEXIST_ERROR,
			msg:      "custom error message",
			expected: EError{Code: EXIFTOOL_NOTEXIST_ERROR},
		},
		{
			name:     "Create error with code 9000001",
			code:     INTERNAL_ERROR,
			msg:      "internal error occurred",
			expected: EError{Code: INTERNAL_ERROR},
		},
		{
			name:     "Create error with empty message",
			code:     123,
			msg:      "",
			expected: EError{Code: 123},
		},
		{
			name:     "Create error with special characters",
			code:     456,
			msg:      "错误信息：测试！@#$%",
			expected: EError{Code: 456},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewErrors(tt.code, tt.msg)
			if result.Code != tt.expected.Code {
				t.Errorf("NewErrors() Code = %v, expected %v", result.Code, tt.expected.Code)
			}
			if result.Error.Error() != tt.msg {
				t.Errorf("NewErrors() Error = %v, expected %v", result.Error.Error(), tt.msg)
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name           string
		err            EError
		expectedCode   int
		expectedMsg    string
		expectedIsOk   bool
		expectedHasErr bool
	}{
		{
			name:           "NoError",
			err:            NoError,
			expectedCode:   0,
			expectedMsg:    "no error",
			expectedIsOk:   true,
			expectedHasErr: false,
		},
		{
			name:           "InternalError",
			err:            InternalError,
			expectedCode:   INTERNAL_ERROR,
			expectedMsg:    "发生内部错误",
			expectedIsOk:   false,
			expectedHasErr: true,
		},
		{
			name:           "ExiftoolNotExistError",
			err:            ExiftoolNotExistError,
			expectedCode:   EXIFTOOL_NOTEXIST_ERROR,
			expectedMsg:    "exiftool工具不存在,请检查是否安装",
			expectedIsOk:   false,
			expectedHasErr: true,
		},
		{
			name:           "ExiftoolInitError",
			err:            ExiftoolInitError,
			expectedCode:   EXIFTOOL_INIT_ERROR,
			expectedMsg:    "exiftool工具init失败,请检查环境",
			expectedIsOk:   false,
			expectedHasErr: true,
		},
		{
			name:           "ExiftoolImageError",
			err:            ExiftoolImageError,
			expectedCode:   EXIFTOOL_IMAGE_EXIF_ERROR,
			expectedMsg:    "exiftool工具获取图片exif信息失败",
			expectedIsOk:   false,
			expectedHasErr: true,
		},
		{
			name:           "ExiftoolCacheTypeError",
			err:            ExiftoolCacheTypeError,
			expectedCode:   EXIFTOOL_IMAGE_EXIF_CACHE_ERROR,
			expectedMsg:    "exif cache缓存数据类型断言失败",
			expectedIsOk:   false,
			expectedHasErr: true,
		},
		{
			name:           "ImageDecodeCacheTypeError",
			err:            ImageDecodeCacheTypeError,
			expectedCode:   IMAGE_DECODE_CACHE_ERROR,
			expectedMsg:    "缓存的图片解码数据类型断言失败",
			expectedIsOk:   false,
			expectedHasErr: true,
		},
		{
			name:           "ImageRGBACacheTypeError",
			err:            ImageRGBACacheTypeError,
			expectedCode:   IMAGE_RGBA_CACHE_ERROR,
			expectedMsg:    "缓存的RGBA数据类型断言失败",
			expectedIsOk:   false,
			expectedHasErr: true,
		},
		{
			name:           "ImageTextCacheTypeError",
			err:            ImageTextCacheTypeError,
			expectedCode:   IMAGE_TEXT_FONT_CACHE_ERROR,
			expectedMsg:    "缓存中获取的字体对象类型断言失败",
			expectedIsOk:   false,
			expectedHasErr: true,
		},
		{
			name:           "ImageLogoNotFindError",
			err:            ImageLogoNotFindError,
			expectedCode:   IMAGE_LOGO_NOT_FIND_ERROR,
			expectedMsg:    "logo图片文件没有找到",
			expectedIsOk:   false,
			expectedHasErr: true,
		},
		{
			name:           "ImageLogoResizeError",
			err:            ImageLogoResizeError,
			expectedCode:   IMAGE_LOGO_RESIZE_ERROR,
			expectedMsg:    "logo图片重置尺寸错误",
			expectedIsOk:   false,
			expectedHasErr: true,
		},
		{
			name:           "ImageJpegSaveError",
			err:            ImageJpegSaveError,
			expectedCode:   IMAGE_JPEG_SAVE_ERROR,
			expectedMsg:    "jpeg图片保存失败",
			expectedIsOk:   false,
			expectedHasErr: true,
		},
		{
			name:           "LayoutNotFindError",
			err:            LayoutNotFindError,
			expectedCode:   LAYOUT_TYPE_NOT_FIND_ERROR,
			expectedMsg:    "布局类型查找失败",
			expectedIsOk:   false,
			expectedHasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.expectedCode {
				t.Errorf("Error Code = %v, expected %v", tt.err.Code, tt.expectedCode)
			}
			if tt.err.Error.Error() != tt.expectedMsg {
				t.Errorf("Error Message = %v, expected %v", tt.err.Error.Error(), tt.expectedMsg)
			}
			if IsOk(tt.err) != tt.expectedIsOk {
				t.Errorf("IsOk() = %v, expected %v", IsOk(tt.err), tt.expectedIsOk)
			}
			if HasError(tt.err) != tt.expectedHasErr {
				t.Errorf("HasError() = %v, expected %v", HasError(tt.err), tt.expectedHasErr)
			}
		})
	}
}

func TestEErrorUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  string
		expectErr bool
	}{
		{
			name:      "Unmarshal valid JSON",
			jsonData:  `{"errmsg":"test error","code":123}`,
			expectErr: true,
		},
		{
			name:      "Unmarshal JSON with code 0",
			jsonData:  `{"errmsg":"no error","code":0}`,
			expectErr: true,
		},
		{
			name:      "Unmarshal JSON with Chinese message",
			jsonData:  `{"errmsg":"发生内部错误","code":9000001}`,
			expectErr: true,
		},
		{
			name:      "Unmarshal invalid JSON",
			jsonData:  `{"errmsg":"test error","code":"invalid"}`,
			expectErr: true,
		},
		{
			name:      "Unmarshal malformed JSON",
			jsonData:  `{invalid json}`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err EError
			unmarshalErr := json.Unmarshal([]byte(tt.jsonData), &err)

			if tt.expectErr {
				if unmarshalErr == nil {
					t.Error("Expected unmarshal error, but got nil")
				}
				return
			}

			if unmarshalErr != nil {
				t.Fatalf("UnmarshalJSON() failed: %v", unmarshalErr)
			}
		})
	}
}

func TestEErrorWithLargeCode(t *testing.T) {
	tests := []struct {
		name string
		code int
	}{
		{
			name: "Code with maximum int value",
			code: 2147483647,
		},
		{
			name: "Code with large value",
			code: 999999999,
		},
		{
			name: "Code with minimum int value",
			code: -2147483648,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewErrors(tt.code, "test message")
			if err.Code != tt.code {
				t.Errorf("Code = %v, expected %v", err.Code, tt.code)
			}

			if tt.code > 0 && !HasError(err) {
				t.Error("HasError() should return true for positive code")
			}

			if tt.code == 0 && !IsOk(err) {
				t.Error("IsOk() should return true for code 0")
			}

			if tt.code < 0 && HasError(err) {
				t.Error("HasError() should return false for negative code")
			}
		})
	}
}

func TestEErrorJSONFieldNames(t *testing.T) {
	err := NewErrors(123, "test error")
	jsonData, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("MarshalJSON() failed: %v", marshalErr)
	}

	jsonStr := string(jsonData)

	if !contains(jsonStr, `"errmsg"`) {
		t.Error("JSON should contain 'errmsg' field")
	}

	if !contains(jsonStr, `"code"`) {
		t.Error("JSON should contain 'code' field")
	}

	if contains(jsonStr, `"Error"`) {
		t.Error("JSON should not contain 'Error' field (should be 'errmsg')")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
