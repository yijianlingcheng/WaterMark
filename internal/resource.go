package internal

import (
	"WaterMark/assetfs"
	"WaterMark/pkg"
)

// 释放已经保存好的exiftool.zip压缩文件,并实现自动解压,保证在win系统上exiftool工具可用.
func winRestoreExitoolZipFile() pkg.EError {
	// Windows系统逻辑
	if !IsWindows() {
		return pkg.NoError
	}
	// 判断可执行文件是否存在,不存在则从文件中输出对应的exiftool.zip文件
	// 释放完成exiftool.zip文件之后对其进行解压
	if !PathExists(GetExiftoolPath()) {
		// 从文件中释放exiftool.zip文件
		err := assetfs.RestoreAssets(GetRootPath(), "exiftool")
		if err != nil {
			Log.Panic("exiftool.zip文件释放失败,程序异常退出:" + err.Error())
		}
		// 判断文件是否释放成功,释放失败程序直接终止
		if !PathExists(GetExiftoolZipPath()) {
			Log.Panic("exiftool.zip文件释放失败,程序异常退出")
		}
		// 解压文件
		Unzip(GetExiftoolZipPath(), GetExiftoolUnzipPath())

		// 判断zip是否解压成功
		if !PathExists(GetExiftoolPath()) {
			Log.Panic("exiftool.zip文件解压失败,程序异常退出")
		}
	}

	return pkg.NoError
}

// 释放字体文件.
func restoreFontFile() pkg.EError {
	err := assetfs.RestoreAssets(GetRootPath(), "fonts")
	if err != nil {
		Log.Panic("字体文件释放失败,程序异常退出:" + err.Error())
	}

	return pkg.NoError
}

// 释放logo文件.
func restoreLogoFile() pkg.EError {
	err := assetfs.RestoreAssets(GetRootPath(), "logos")
	if err != nil {
		Log.Panic("logo文件释放失败,程序异常退出:" + err.Error())
	}

	return pkg.NoError
}
