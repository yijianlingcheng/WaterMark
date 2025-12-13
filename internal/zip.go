package internal

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// 解压zip文件到指定路径.
func Unzip(zipPath, unzipPath string) {
	// 第一步,打开 zip 文件
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		Log.Panic("打开压缩包" + zipPath + "出错:" + err.Error())
	}
	defer zipFile.Close()

	// 第二步,遍历 zip 中的文件
	for _, f := range zipFile.File {
		filePath := unzipPath + f.Name
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(filePath, os.ModePerm)

			continue
		}
		// 创建对应文件夹
		if mkErr := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); mkErr != nil {
			Log.Panic("创建文件" + filepath.Dir(filePath) + "失败:" + mkErr.Error())
		}
		// 解压到的目标文件
		dstFile, openFileErr := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if openFileErr != nil {
			Log.Panic("打开文件" + filePath + "失败:" + openFileErr.Error())
		}
		file, openErr := f.Open()
		if openErr != nil {
			Log.Panic("打开文件" + f.Name + "失败:" + openErr.Error())
		}
		// 写入到解压到的目标文件
		//nolint:gosec
		if _, copyErr := io.Copy(dstFile, file); copyErr != nil {
			Log.Panic("复制文件到指定路径" + f.Name + "失败:" + copyErr.Error())
		}
		dstFile.Close()
		file.Close()
	}
}
