package tool

import (
	"io"
	"os"
	"runtime"
	"strings"
)

// Exists 判断路径是否存在
//
//	@param path
//	@return bool
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	return !os.IsNotExist(err)
}

// CreateDir 创建目录
//
//	@param path
func CreateDir(path string) {
	os.MkdirAll(path, os.ModePerm)
}

// ReplaceDir 替换win的目录
//
//	@param path
//	@return string
func ReplaceDir(path string) string {
	if runtime.GOOS == "windows" {
		path = strings.ReplaceAll(path, "\\", "/")
	}
	return path
}

// CopyFile 复制文件
//
//	@param src
//	@param dst
//	@return int64
//	@return error
func CopyFile(src, dst string, bufferSize int) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()
	save, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer save.Close()
	buf := make([]byte, bufferSize)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := save.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}
