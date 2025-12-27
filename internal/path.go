package internal

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var rootPath string

// initRootPath 初始化项目根路径,解决go run 与go build执行时路径不一致问题.
func initRootPath() {
	if !ISRelease() {
		_, filename, _, ok := runtime.Caller(0)
		if ok {
			rootPath = path.Dir(path.Dir(filename))
		}

		return
	}
	runPath, _ := os.Executable()
	rootPath = filepath.Dir(runPath)
}

// 获取项目根路径.
func GetRootPath() string {
	return rootPath
}

// 获取项目指定路径.
func GetPwdPath(p string) string {
	return rootPath + p
}

// 获取项目配置路径.
func GetConfigPath(p string) string {
	return rootPath + appConfigsPath + "/" + p
}

// 获取日志配置路径.
func GetLogPath(p string) string {
	return rootPath + appLogsPath + "/" + p
}

// 获取logo的存放路径.
func GetLogosPath(p string) string {
	return rootPath + appLogosPath + "/" + p
}

// 获取临时文件存放路径.
func GetRuntimePath(p string) string {
	return rootPath + appRuntimePath + "/" + p
}

// 获取用户文件存放路径.
func GetUserDirectory(p string) string {
	return rootPath + appUserPath + "/" + p
}

// 获取字体文件路径.
func GetFontFilePath(p string) string {
	return rootPath + appFontFilePath + "/" + p
}

// 获取app缓存exif的文件路径.
func GetAppExifCacheFilePath() string {
	return GetRuntimePath("exifCache.cache")
}

// 获取exiftool.zip文件路径.
func GetExiftoolZipPath() string {
	return GetRootPath() + appWinExiftoolZipPath
}

// 获取exiftool.zip文件解压之后的存放路径.
func GetExiftoolUnzipPath() string {
	return GetRootPath() + appExiftoolPath + "/"
}

// 获取主要的布局文件.
func GetMainLayoutPath() string {
	return GetRootPath() + appConfigsPath + "/layout.json"
}

// 获取ImageMagick可执行文件路径.
func GetMagickPath(p string) string {
	return GetRootPath() + magickPath + "/" + p
}

// 获取ImageMagick可执行文件路径.
func GetWinMagick7zPath() string {
	return GetMagickPath("ImageMagick.7z")
}

// 获取ImageMagick可执行文件路径.
func GetMagickBinPath() string {
	if IsWindows() {
		return GetMagickPath("magick.exe")
	}

	return "magick"
}

// 获取exiftool路径.
func GetExiftoolPath() string {
	if IsWindows() {
		return GetRootPath() + appWinExiftoolPath
	}

	return appDarwinExiftoolPath
}

// 初始化程序需要的各种文件夹.
func createAppDS(list []string) {
	for _, i := range list {
		if PathExists(i) {
			continue
		}
		if err := os.MkdirAll(i, os.ModePerm); err != nil {
			panic("创建运行时文件夹失败,创建失败的文件路径为:" + i)
		}
	}
}

// 判断所给路径文件/文件夹是否存在.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	return !errors.Is(err, os.ErrNotExist)
}
