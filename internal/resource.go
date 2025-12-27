package internal

import (
	"os"
	"strings"

	"WaterMark/assetexiffs"
	"WaterMark/assetmagickfs"
	"WaterMark/assetmixfs"
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
		err := assetexiffs.RestoreAssets(GetRootPath(), "exiftool")
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

// 释放已经保存好的ImageMagick.7z压缩文件,并实现自动解压,保证在win系统上ImageMagick工具可用.
func winRestoreImagemagick7zFile() pkg.EError {
	// Windows系统逻辑
	if !IsWindows() {
		return pkg.NoError
	}
	if !PathExists(GetMagickBinPath()) {
		// 从文件中释放ImageMagick.7z文件
		err := assetmagickfs.RestoreAssets(GetRootPath(), "magick")
		if err != nil {
			Log.Panic("ImageMagick.7z文件释放失败,程序异常退出:" + err.Error())
		}
		// 判断文件是否释放成功,释放失败程序直接终止
		if !PathExists(GetWinMagick7zPath()) {
			Log.Panic("ImageMagick.7z文件释放失败,程序异常退出")
		}
		// 解压文件
		Unzip7z(GetWinMagick7zPath(), GetMagickPath(""))
		editImageMagickConfig()
	}

	return pkg.NoError
}

// 修改ImageMagick配置文件,调整多线程处理能力.防止cpu占用过高.
func editImageMagickConfig() {
	policyPath := GetMagickPath("policy.xml")
	content, err := os.ReadFile(policyPath)
	if err != nil {
		Log.Panic("读取ImageMagick配置文件失败,程序异常退出:" + err.Error())
	}
	newContent := strings.ReplaceAll(
		string(content),
		"<!-- <policy domain=\"resource\" name=\"thread\" value=\"2\"/> -->",
		"<policy domain=\"resource\" name=\"thread\" value=\"8\"/>",
	)
	err = os.WriteFile(policyPath, []byte(newContent), 0o600)
	if err != nil {
		Log.Panic("写入ImageMagick配置文件失败,程序异常退出:" + err.Error())
	}
}

// 释放字体文件.
func restoreFontFile() pkg.EError {
	err := assetmixfs.RestoreAssets(GetRootPath(), "fonts")
	if err != nil {
		Log.Panic("字体文件释放失败,程序异常退出:" + err.Error())
	}

	return pkg.NoError
}

// 释放logo文件.
func restoreLogoFile() pkg.EError {
	err := assetmixfs.RestoreAssets(GetRootPath(), "logos")
	if err != nil {
		Log.Panic("logo文件释放失败,程序异常退出:" + err.Error())
	}

	return pkg.NoError
}
