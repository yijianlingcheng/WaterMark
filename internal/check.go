package internal

import (
	"strings"
	"time"

	"WaterMark/internal/cmd"
	"WaterMark/pkg"
)

// 检查与释放资源.
func checkInstallExif() pkg.EError {
	err := winRestoreExitoolZipFile()
	if pkg.HasError(err) {
		return err
	}
	err = checkExiftool()
	if pkg.HasError(err) {
		return err
	}

	return pkg.NoError
}

// 检查exif工具.
func checkExiftool() pkg.EError {
	args := []string{GetExiftoolPath(), "-ver"}
	version, err := cmd.CommandRun(5*time.Second, strings.Join(args, " "))

	if version == "" || pkg.HasError(err) {
		return pkg.ExiftoolNotExistError
	}

	return pkg.NoError
}

// 检查字体文件.
func checkFontFile() pkg.EError {
	fontDir := GetFontFilePath("")
	list, err := pkg.GetDirFiles(fontDir)

	if pkg.HasError(err) {
		return err
	}
	if len(list) == 0 {
		restoreFontFile()
	}

	return pkg.NoError
}

// 检查logo文件.
func checkLogoFile() pkg.EError {
	logoDir := GetLogosPath("")
	list, err := pkg.GetDirFiles(logoDir)

	if pkg.HasError(err) {
		return err
	}
	if len(list) == 0 {
		restoreLogoFile()
	}

	return pkg.NoError
}

// 检查ImageMagick工具.
func checkInstallImageMagick() pkg.EError {
	err := winRestoreImagemagick7zFile()
	if pkg.HasError(err) {
		return err
	}
	err = checkImageMagick()
	if pkg.HasError(err) {
		return err
	}

	return pkg.NoError
}

// 检查ImageMagick工具.
func checkImageMagick() pkg.EError {
	args := []string{GetMagickBinPath(), "-version"}
	version, err := cmd.CommandRun(5*time.Second, strings.Join(args, " "))

	if version == "" || pkg.HasError(err) {
		return pkg.ExiftoolNotExistError
	}

	return pkg.NoError
}
