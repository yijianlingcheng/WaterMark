package internal

import (
	"io"
	"os"
	"path/filepath"

	"github.com/bodgit/sevenzip"

	"WaterMark/pkg"
)

// 解压7z压缩文件.
func Unzip7z(zipPath, unzipPath string) pkg.EError {
	// 创建7z读取器
	reader, err := sevenzip.OpenReader(zipPath)
	if err != nil {
		Log.Panic("无法创建7z读取器: " + err.Error())
	}

	for _, f := range reader.File {
		filePath := unzipPath + f.Name
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(filePath, os.ModePerm)

			continue
		}
		// 创建对应文件夹
		if mkErr := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); mkErr != nil {
			Log.Panic("创建文件" + filepath.Dir(filePath) + "失败:" + mkErr.Error())
		}

		if err = extract7zFile(unzipPath, f); err != nil {
			Log.Panic("无法解压7z文件: " + err.Error())
		}
	}

	return pkg.NoError
}

// 解压7z文件.
func extract7zFile(unzipPath string, f *sevenzip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	filePath := unzipPath + f.Name
	dstFile, openFileErr := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())

	if openFileErr != nil {
		Log.Panic("无法解压7z文件: " + openFileErr.Error())
	}
	_, err = io.Copy(dstFile, rc)
	if err != nil {
		return err
	}

	return nil
}
