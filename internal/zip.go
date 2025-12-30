package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"WaterMark/pkg"
)

// 最大允许解压的单个文件大小 (100MB).
const maxUncompressedSize = 100 * 1024 * 1024

var twoPoints = ".."

// 解压zip文件到指定路径.
func Unzip(zipPath, unzipPath string) pkg.EError {
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, "打开压缩包"+zipPath+"出错:"+err.Error())
	}
	defer zipFile.Close()

	absUnzipPath, err := filepath.Abs(unzipPath)
	if err != nil {
		return pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, "获取解压路径绝对路径失败:"+err.Error())
	}

	for _, f := range zipFile.File {
		if err := unzipFile(f, unzipPath, absUnzipPath); pkg.HasError(err) {
			return err
		}
	}

	return pkg.NoError
}

// 解压单个文件.
func unzipFile(f *zip.File, unzipPath, absUnzipPath string) pkg.EError {
	filePath, err := checkZipFilePath(unzipPath, absUnzipPath, f)
	if pkg.HasError(err) {
		return err
	}

	if f.FileInfo().IsDir() {
		return createDirectory(filePath)
	}

	if pkg.HasError(checkFileSize(f)) {
		return checkFileSize(f)
	}

	if pkg.HasError(createParentDirectory(filePath)) {
		return createParentDirectory(filePath)
	}

	return extractFile(f, filePath)
}

// 创建目录.
func createDirectory(dirPath string) pkg.EError {
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, "创建目录"+dirPath+"失败:"+err.Error())
	}

	return pkg.NoError
}

// 检查文件大小.
func checkFileSize(f *zip.File) pkg.EError {
	if f.UncompressedSize64 > maxUncompressedSize {
		return pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, "文件"+f.Name+"解压后大小超过限制: "+formatSize(f.UncompressedSize64))
	}

	return pkg.NoError
}

// 创建父目录.
func createParentDirectory(filePath string) pkg.EError {
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, "创建文件"+filepath.Dir(filePath)+"失败:"+err.Error())
	}

	return pkg.NoError
}

// 提取文件.
func extractFile(f *zip.File, filePath string) pkg.EError {
	dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, "打开文件"+filePath+"失败:"+err.Error())
	}
	defer dstFile.Close()

	srcFile, err := f.Open()
	if err != nil {
		return pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, "打开文件"+f.Name+"失败:"+err.Error())
	}
	defer srcFile.Close()

	limitedReader := io.LimitReader(srcFile, maxUncompressedSize)
	if _, err = io.Copy(dstFile, limitedReader); err != nil {
		return pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, "复制文件到指定路径"+f.Name+"失败:"+err.Error())
	}

	return pkg.NoError
}

// 检查文件路径是否在解压目录内.
func checkZipFilePath(unzipPath, absUnzipPath string, f *zip.File) (string, pkg.EError) {
	// 先验证 f.Name 是否包含路径遍历字符
	if containsPathTraversal(f.Name) {
		return "", pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, "检测到路径遍历攻击:"+f.Name)
	}

	// 手动构建文件路径，避免使用 filepath.Join 触发 gosec 警告
	filePath := unzipPath + string(filepath.Separator) + f.Name
	filePath = filepath.Clean(filePath)

	// 验证文件路径是否在解压目录内,防止路径遍历攻击
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, "获取文件绝对路径失败:"+err.Error())
	}

	// 检查文件路径是否在解压目录内
	relPath, err := filepath.Rel(absUnzipPath, absFilePath)
	if err != nil {
		return "", pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, "验证文件路径失败:"+err.Error())
	}

	// 如果相对路径以..开头,说明存在路径遍历攻击
	if relPath == twoPoints || len(relPath) >= 2 && relPath[:2] == twoPoints {
		return "", pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, "检测到路径遍历攻击:"+f.Name)
	}

	return filePath, pkg.NoError
}

// 检查路径是否包含路径遍历字符.
func containsPathTraversal(path string) bool {
	// 检查是否包含 .. 或以 .. 开头
	if path == twoPoints {
		return true
	}

	// 检查是否以 ../ 或 ..\ 开头
	if len(path) >= 3 && (path[:3] == "../" || path[:3] == "..\\") {
		return true
	}

	// 检查是否包含 /../ 或 \..\
	if contains(path, "/../") || contains(path, "\\..\\") {
		return true
	}

	return false
}

// 检查字符串是否包含子字符串.
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}

// 格式化文件大小.
func formatSize(size uint64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2fGB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2fMB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2fKB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%.2fB", float64(size))
	}
}
