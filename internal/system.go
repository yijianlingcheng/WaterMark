package internal

import (
	"WaterMark/pkg"
)

// exiftool win系统下面识别图片exif信息的工具,在win系统中直接内置,解决下载问题
// logs 程序运行时日志存放的目录
// logos 相机或者镜头厂家logo存放的路径
// configs 自定义配置文件的存放路径,用于保存用户自定义设置信息
// runtime 代表允许过程中存放中间文件,缓存文件的地方
// blur 代表允许过程中,生成的模糊图片存放的路径
// userData 代表用户自定义数据存放的路径
// default 程序允许过程中,下载文件默认保存的地方
// fonts 字体文件路径.
// magick ImageMagick可执行文件存放路径.
var (
	appExiftoolPath = "/exiftool"

	appConfigsPath = "/configs"

	appLogsPath = "/logs"

	appLogosPath = "/logos"

	appRuntimePath = "/runtime"

	appBlurPath = appRuntimePath + "/blur"

	appUserPath = "/userData"

	magickPath = "/magick"

	appFontFilePath = "/fonts"

	appRunNeedDS = []string{
		appExiftoolPath,
		appLogsPath,
		appLogosPath,
		appConfigsPath,
		appRuntimePath,
		appBlurPath,
		appUserPath,
		appFontFilePath,
	}

	// win系统下面的exiftool压缩文件路径.
	appWinExiftoolZipPath = appExiftoolPath + "/exiftool.zip"

	// win系统下面exiftool可执行文件的路径.
	appWinExiftoolPath = appExiftoolPath + "/exiftool.exe"

	// MacOS 系统中的exiftool可执行文件路径.
	appDarwinExiftoolPath = "exiftool"
)

// 获取APP运行时需要的全部文件夹列表.
func getAppAllRuntimeDS() []string {
	// 需要自动创建的文件夹列表
	list := appRunNeedDS
	for i := range list {
		list[i] = GetRootPath() + list[i]
	}

	return list
}

// 检查当前运行环境是否为window.
func IsWindows() bool {
	return pkg.IsWindows()
}
