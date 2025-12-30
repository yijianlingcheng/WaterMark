package internal

import (
	"sync"

	"WaterMark/message"
	"WaterMark/pkg"
)

const (
	msgCheckDir                = "检查程序执行需要的目录"
	msgCheckDirSuccess         = "检查程序执行需要的目录通过"
	msgLoadConfig              = "加载应用配置"
	msgLoadConfigSuccess       = "加载应用配置完成"
	msgInitLog                 = "初始化日志记录的配置"
	msgInitLogSuccess          = "初始化日志记录的配置完成"
	msgInitCustomize           = "初始化用户自定义配置"
	msgInitCustomizeSuccess    = "初始化用户自定义配置完成"
	msgCheckExiftool           = "检查exiftool工具是否安装"
	msgCheckExiftoolSuccess    = "exiftool工具已安装"
	msgCheckImageMagick        = "检查ImageMagick工具是否安装"
	msgCheckImageMagickSuccess = "ImageMagick工具已安装"
	msgCheckFont               = "检查字体文件是否释放"
	msgCheckFontSuccess        = "字体文件已存放"
	msgCheckLogo               = "检查logo是否释放"
	msgCheckLogoSuccess        = "logo已存在"
)

var once sync.Once

func InitAppConfigsAndRes() {
	once.Do(func() {
		// 初始化应用根路径
		initRootPath()

		if err := initApp(); pkg.HasError(err) {
			message.SendErrorMsg(err.String())

			return
		}

		if ISApiDebug() {
			go message.ApiDebug()
		}
	})
}

func initApp() pkg.EError {
	if err := initDirectories(); pkg.HasError(err) {
		return err
	}

	if err := initConfigurations(); pkg.HasError(err) {
		return err
	}

	if err := initTools(); pkg.HasError(err) {
		return err
	}

	if err := initResources(); pkg.HasError(err) {
		return err
	}

	return pkg.NoError
}

func initDirectories() pkg.EError {
	message.SendInfoMsg(msgCheckDir)
	createAppDS(getAppAllRuntimeDS())
	message.SendInfoMsg(msgCheckDirSuccess)

	return pkg.NoError
}

func initConfigurations() pkg.EError {
	message.SendInfoMsg(msgLoadConfig)
	err := initAppConfig()
	message.SendErrorOrInfo(err, msgLoadConfigSuccess)
	if pkg.HasError(err) {
		return err
	}

	message.SendInfoMsg(msgInitLog)
	initLogConfig()
	message.SendInfoMsg(msgInitLogSuccess)

	message.SendInfoMsg(msgInitCustomize)
	err = initCustomizeConfig()
	message.SendErrorOrInfo(err, msgInitCustomizeSuccess)
	if pkg.HasError(err) {
		return err
	}

	return pkg.NoError
}

func initTools() pkg.EError {
	message.SendInfoMsg(msgCheckExiftool)
	err := checkInstallExif()
	message.SendErrorOrInfo(err, msgCheckExiftoolSuccess)
	if pkg.HasError(err) {
		return err
	}

	message.SendInfoMsg(msgCheckImageMagick)
	err = checkInstallImageMagick()
	message.SendErrorOrInfo(err, msgCheckImageMagickSuccess)
	if pkg.HasError(err) {
		return err
	}

	return pkg.NoError
}

func initResources() pkg.EError {
	message.SendInfoMsg(msgCheckFont)
	err := checkFontFile()
	message.SendErrorOrInfo(err, msgCheckFontSuccess)
	if pkg.HasError(err) {
		return err
	}

	message.SendInfoMsg(msgCheckLogo)
	err = checkLogoFile()
	message.SendErrorOrInfo(err, msgCheckLogoSuccess)
	if pkg.HasError(err) {
		return err
	}

	return pkg.NoError
}
