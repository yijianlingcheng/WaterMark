package config

import (
	"WaterMark/src/cmd"
	"WaterMark/src/images"
	"WaterMark/src/logs"
	"WaterMark/src/paths"
	"WaterMark/src/tool"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	// 临时目录
	TEMP_DIR = "/tmp"

	// 预览图存放路径
	PREVIEW_PATH = TEMP_DIR + "/preview/"

	// 预览缩略图存放路径
	SMALL_PREVIEW_PATH = TEMP_DIR + "/small/"
)

// Load 加载配置
func Load() {
	paths.InitRootPath()
	logs.InitLog()
	checkAndCreateDir()
	loadConfig()
	initLogoList()
	initTemplate()
	cmd.InitToolPath()
}

// loadConfig 加载主配置文件
func loadConfig() {
	// viper 加载主配置文件
	viper.SetConfigFile(paths.GetPwdPath("/config/app.yaml"))
	if err := viper.ReadInConfig(); err != nil { // 查找并读取配置文件
		panic(fmt.Errorf("fatal error load config file: %w", err))
	}

	viper.WatchConfig()

	// viper 配置重新加载记录日志
	viper.OnConfigChange(func(e fsnotify.Event) {
		logs.Info.Println("Config app.yaml file updated")
	})
}

// initLogoList 初始化logolist
func initLogoList() {
	images.LoadlogoList()
}

// initTemplate 加载模板配置文件
func initTemplate() {
	// viper 加载模板配置文件
	viper.SetConfigFile(paths.GetPwdPath("/config/tpl.yaml"))
	if err := viper.MergeInConfig(); err != nil { // 查找并读取配置文件
		panic(fmt.Errorf("fatal error load config file: %w", err))
	}

	viper.WatchConfig()

	// viper 配置重新加载记录日志
	viper.OnConfigChange(func(e fsnotify.Event) {
		logs.Info.Println("Config tpl.yaml file updated")
	})

	images.LoadTemplate()
}

// checkAndCreateDir 检查并创建文件夹
func checkAndCreateDir() {
	if !tool.Exists(paths.GetPwdPath(PREVIEW_PATH)) {
		tool.CreateDir(paths.GetPwdPath(PREVIEW_PATH))
	}
	if !tool.Exists(paths.GetPwdPath(SMALL_PREVIEW_PATH)) {
		tool.CreateDir(paths.GetPwdPath(SMALL_PREVIEW_PATH))
	}
	// 初始化图片存放路径
	images.PreviewPath = paths.GetPwdPath(PREVIEW_PATH)
	images.SmallPreviewPath = paths.GetPwdPath(SMALL_PREVIEW_PATH)
}
