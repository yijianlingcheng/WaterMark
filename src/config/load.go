package config

import (
	"WaterMark/src/images"
	. "WaterMark/src/logs"
	"WaterMark/src/tool"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	// 临时目录
	TEMP_DIR = "./tmp"

	// 预览图存放路径
	PREVIEW_PATH = TEMP_DIR + "/preview/"

	// 预览缩略图存放路径
	SMALL_PREVIEW_PATH = TEMP_DIR + "/small/"
)

// Load 加载配置
func Load() {
	checkAndCreateDir()
	loadConfig()
	initLogoList()
	initTemplate()
}

// loadConfig 加载主配置文件
func loadConfig() {

	// viper 加载主配置文件
	viper.SetConfigFile("./config/app.yaml")
	if err := viper.ReadInConfig(); err != nil { // 查找并读取配置文件
		panic(fmt.Errorf("fatal error load config file: %w", err))
	}

	viper.WatchConfig()

	// viper 配置重新加载记录日志
	viper.OnConfigChange(func(e fsnotify.Event) {
		Info.Println("Config app.yaml file updated")
	})
}

// initLogoList 初始化logolist
func initLogoList() {
	images.LoadlogoList()
}

// initTemplate 加载模板配置文件
func initTemplate() {
	// viper 加载模板配置文件
	viper.SetConfigFile("./config/tpl.yaml")
	if err := viper.MergeInConfig(); err != nil { // 查找并读取配置文件
		panic(fmt.Errorf("fatal error load config file: %w", err))
	}

	viper.WatchConfig()

	// viper 配置重新加载记录日志
	viper.OnConfigChange(func(e fsnotify.Event) {
		Info.Println("Config tpl.yaml file updated")
	})

	images.LoadTemplate()
}

// checkAndCreateDir 检查并创建文件夹
func checkAndCreateDir() {
	if !tool.Exists(PREVIEW_PATH) {
		tool.CreateDir(PREVIEW_PATH)
	}
	if !tool.Exists(SMALL_PREVIEW_PATH) {
		tool.CreateDir(SMALL_PREVIEW_PATH)
	}
	// 初始化图片存放路径
	images.PreviewPath = PREVIEW_PATH
	images.SmallPreviewPath = SMALL_PREVIEW_PATH
}
