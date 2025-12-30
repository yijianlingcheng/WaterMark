package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"WaterMark/pkg"
)

func TestWinRestoreExitoolZipFile(t *testing.T) {
	// 保存原始变量
	originalRootPath := rootPath

	tests := []struct {
		name      string
		isWindows bool
		setup     func() func()
		wantErr   bool
	}{
		{
			name:      "非Windows系统",
			isWindows: false,
			setup: func() func() {
				return func() {}
			},
			wantErr: false,
		},
		{
			name:      "Windows系统，exiftool已存在",
			isWindows: true,
			setup: func() func() {
				// 创建临时目录结构
				tmpDir := t.TempDir()
				exiftoolDir := filepath.Join(tmpDir, "exiftool")
				exiftoolPath := filepath.Join(exiftoolDir, "exiftool.exe")

				// 创建exiftool.exe文件
				err := os.MkdirAll(exiftoolDir, 0o755)
				if err != nil {
					t.Fatalf("Failed to create test directory: %v", err)
				}
				err = os.WriteFile(exiftoolPath, []byte("fake exiftool"), 0o755)
				if err != nil {
					t.Fatalf("Failed to create fake exiftool: %v", err)
				}

				// 修改rootPath
				rootPath = tmpDir

				return func() {
					rootPath = originalRootPath
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置测试环境
			cleanup := tt.setup()
			defer cleanup()

			// 执行测试函数
			err := winRestoreExitoolZipFile()

			// 检查错误
			if pkg.HasError(err) != tt.wantErr {
				t.Errorf("winRestoreExitoolZipFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWinRestoreImagemagick7zFile(t *testing.T) {
	// 保存原始变量
	originalRootPath := rootPath

	tests := []struct {
		name      string
		isWindows bool
		setup     func() func()
		wantErr   bool
	}{
		{
			name:      "非Windows系统",
			isWindows: false,
			setup: func() func() {
				return func() {}
			},
			wantErr: false,
		},
		{
			name:      "Windows系统，ImageMagick已存在",
			isWindows: true,
			setup: func() func() {
				// 创建临时目录结构
				tmpDir := t.TempDir()
				magickDir := filepath.Join(tmpDir, "magick")
				magickPath := filepath.Join(magickDir, "magick.exe")
				policyPath := filepath.Join(magickDir, "policy.xml")

				// 创建magick.exe文件
				err := os.MkdirAll(magickDir, 0o755)
				if err != nil {
					t.Fatalf("Failed to create test directory: %v", err)
				}
				err = os.WriteFile(magickPath, []byte("fake magick"), 0o755)
				if err != nil {
					t.Fatalf("Failed to create fake magick: %v", err)
				}
				// 创建policy.xml文件
				err = os.WriteFile(policyPath, []byte(`<policies>
  <policy domain="resource" name="thread" value="2"/>
</policies>`), 0o644)
				if err != nil {
					t.Fatalf("Failed to create policy.xml: %v", err)
				}

				// 修改rootPath
				rootPath = tmpDir

				return func() {
					rootPath = originalRootPath
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置测试环境
			cleanup := tt.setup()
			defer cleanup()

			// 执行测试函数
			err := winRestoreImagemagick7zFile()

			// 检查错误
			if pkg.HasError(err) != tt.wantErr {
				t.Errorf("winRestoreImagemagick7zFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEditImageMagickConfig(t *testing.T) {
	// 保存原始变量
	originalRootPath := rootPath

	tests := []struct {
		name           string
		expectedThread string
		setup          func() func()
	}{
		{
			name:           "修改policy.xml配置",
			expectedThread: "4",
			setup: func() func() {
				tmpDir := t.TempDir()
				magickDir := filepath.Join(tmpDir, "magick")
				policyPath := filepath.Join(magickDir, "policy.xml")

				// 创建magick目录
				err := os.MkdirAll(magickDir, 0o755)
				if err != nil {
					t.Fatalf("Failed to create magick directory: %v", err)
				}

				// 创建测试用的policy.xml文件
				initialContent := `<policies>
  <!-- <policy domain="resource" name="thread" value="2"/> -->
</policies>`
				err = os.WriteFile(policyPath, []byte(initialContent), 0o644)
				if err != nil {
					t.Fatalf("Failed to create policy.xml: %v", err)
				}

				// 修改rootPath
				rootPath = tmpDir

				return func() {
					rootPath = originalRootPath
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置测试环境
			cleanup := tt.setup()
			defer cleanup()

			// 执行测试函数
			editImageMagickConfig()

			// 检查配置文件是否被正确修改
			policyPath := GetMagickPath("policy.xml")
			content, err := os.ReadFile(policyPath)
			if err != nil {
				t.Fatalf("Failed to read policy.xml: %v", err)
			}

			// 检查是否包含policy配置
			if !strings.Contains(string(content), `<policy domain="resource" name="thread" value="`) {
				t.Errorf("Expected policy configuration not found in content: %s", string(content))
			}
		})
	}
}

func TestRestoreFontFile(t *testing.T) {
	// 保存原始函数
	originalRootPath := rootPath

	tests := []struct {
		name    string
		setup   func() func()
		wantErr bool
	}{
		{
			name: "成功恢复字体文件",
			setup: func() func() {
				tmpDir := t.TempDir()

				// 修改rootPath
				rootPath = tmpDir

				return func() {
					rootPath = originalRootPath
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置测试环境
			cleanup := tt.setup()
			defer cleanup()

			// 执行测试函数
			err := restoreFontFile()

			// 检查错误
			if pkg.HasError(err) != tt.wantErr {
				t.Errorf("restoreFontFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRestoreLogoFile(t *testing.T) {
	// 保存原始函数
	originalRootPath := rootPath

	tests := []struct {
		name    string
		setup   func() func()
		wantErr bool
	}{
		{
			name: "成功恢复logo文件",
			setup: func() func() {
				tmpDir := t.TempDir()

				// 修改rootPath
				rootPath = tmpDir

				return func() {
					rootPath = originalRootPath
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置测试环境
			cleanup := tt.setup()
			defer cleanup()

			// 执行测试函数
			err := restoreLogoFile()

			// 检查错误
			if pkg.HasError(err) != tt.wantErr {
				t.Errorf("restoreLogoFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
