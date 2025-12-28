package internal

import (
	"github.com/spf13/viper"

	"WaterMark/assetmixfs"
	"WaterMark/pkg"
)

const (
	// 开发模式.
	APP_DEV = "debug"
	// API调试模式.
	APP_API_DEV = "api_debug"
	// 发布模式.
	APP_RELEASE = "release"
)

// 程序运行模式.
var appMode string

// 初始化项目主配置文件.
func initAppConfig() pkg.EError {
	// 配置文件名称
	appConfigName := "app.yaml"

	// 配置文件路径
	appConfigFilePath := GetConfigPath(appConfigName)

	// 释放配置文件
	restoreAppConfigFile(appConfigFilePath)

	// 检查文件是否存在
	if !PathExists(appConfigFilePath) {
		return pkg.NewErrors(pkg.FILE_NOT_EXIST_ERROR, appConfigFilePath+":项目主配置文件不存在")
	}

	// 加载主配置文件
	viper.SetConfigFile(appConfigFilePath)

	if err := viper.ReadInConfig(); err != nil { // 查找并读取配置文件
		return pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, appConfigFilePath+":加载项目主配置文件失败"+err.Error())
	}

	return pkg.NoError
}

// 释放配置文件.
func restoreAppConfigFile(path string) {
	if PathExists(path) {
		return
	}
	err := assetmixfs.RestoreAssets(GetRootPath(), "configs")
	if err != nil {
		panic(path + ":配置文件创建失效,请仔细检查")
	}
}

// 设置app运行模式.
func SetAppMode(mode string) {
	appMode = mode
}

// 是否是release模式.
func ISRelease() bool {
	return appMode == APP_RELEASE
}

// 是否启用接口调试.
func ISApiDebug() bool {
	return appMode == APP_API_DEV
}

// 获取边框插件.
func GetPlugin() string {
	return viper.GetString("frame.plugin")
}
