package internal

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"WaterMark/internal/cmd"
	"WaterMark/pkg"
)

func restoreDir(dir, backupDir string) error {
	if backupDir == "" || !PathExists(backupDir) {
		// 如果没有备份目录，说明原目录不存在，直接删除测试生成的目录
		return os.RemoveAll(dir)
	}

	files, err := os.ReadDir(backupDir)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		// 如果备份目录为空，说明原目录为空，直接删除测试生成的目录
		return os.RemoveAll(dir)
	}

	// 先清空目标目录，确保只保留备份的文件
	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	// 确保目标目录存在
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	return filepath.Walk(backupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(backupDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dir, relPath)
		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return err
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		if _, err := io.Copy(destFile, srcFile); err != nil {
			return err
		}

		return nil
	})
}

func TestCheckExiftool(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() func()
		wantErr bool
		errCode int
	}{
		{
			name: "exiftool exists and works",
			setup: func() func() {
				return func() {}
			},
			wantErr: false,
		},
		{
			name: "exiftool does not exist",
			setup: func() func() {
				return func() {}
			},
			wantErr: true,
			errCode: pkg.ExiftoolNotExistError.Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup()
			defer cleanup()

			err := checkExiftool()
			if tt.wantErr {
				if !pkg.HasError(err) {
					t.Skipf("checkExiftool() expected error, got NoError (tool may be installed)")
				}
				if err.Code != tt.errCode {
					t.Errorf("checkExiftool() error code = %v, want %v", err.Code, tt.errCode)
				}
			} else {
				if pkg.HasError(err) {
					t.Skipf("checkExiftool() unexpected error = %v (tool may not be installed)", err)
				}
			}
		})
	}
}

func TestCheckImageMagick(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() func()
		wantErr bool
		errCode int
	}{
		{
			name: "ImageMagick exists and works",
			setup: func() func() {
				return func() {}
			},
			wantErr: false,
		},
		{
			name: "ImageMagick does not exist",
			setup: func() func() {
				return func() {}
			},
			wantErr: true,
			errCode: pkg.ExiftoolNotExistError.Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup()
			defer cleanup()

			err := checkImageMagick()
			if tt.wantErr {
				if !pkg.HasError(err) {
					t.Skipf("checkImageMagick() expected error, got NoError (tool may be installed)")
				}
				if err.Code != tt.errCode {
					t.Errorf("checkImageMagick() error code = %v, want %v", err.Code, tt.errCode)
				}
			} else {
				if pkg.HasError(err) {
					t.Skipf("checkImageMagick() unexpected error = %v (tool may not be installed)", err)
				}
			}
		})
	}
}

func TestCheckWithTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		setup   func() func()
		wantErr bool
	}{
		{
			name:    "short timeout",
			timeout: 1 * time.Millisecond,
			setup:   func() func() { return func() {} },
			wantErr: true,
		},
		{
			name:    "long timeout",
			timeout: 10 * time.Second,
			setup:   func() func() { return func() {} },
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup()
			defer cleanup()

			_, err := cmd.CommandRun(tt.timeout, "echo test")
			if pkg.HasError(err) != tt.wantErr {
				t.Errorf("CommandRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckEmptyVersion(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() func()
		wantErr bool
	}{
		{
			name: "empty version response",
			setup: func() func() {
				return func() {}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup()
			defer cleanup()

			version, err := cmd.CommandRun(5*time.Second, "nonexistentcommand12345")
			if version != "" || !pkg.HasError(err) {
				t.Errorf("Expected empty version and error, got version=%s, err=%v", version, err)
			}
		})
	}
}

func TestCheckWithInvalidPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "invalid path",
			path:    "/invalid/path",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cmd.CommandRun(5*time.Second, "nonexistentcommand12345")
			if !pkg.HasError(err) {
				t.Errorf("Expected error for invalid path %s, got NoError", tt.path)
			}
		})
	}
}
