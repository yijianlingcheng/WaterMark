package config

import (
	"WaterMark/src/images"
	"WaterMark/src/log"
	"WaterMark/src/tool"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Load 加载配置
func Load() {
	loadConfig()
	initLogoList()
	initTemplate()
	checkAndCreateDir()
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
		log.InfoLogger.Println("Config app.yaml file updated")
	})
}

// initLogoList 初始化logolist
func initLogoList() {
	images.LoadLogoList()
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
		log.InfoLogger.Println("Config tpl.yaml file updated")
	})

	images.LoadTemplate()
}

// checkAndCreateDir 检查并创建文件夹
func checkAndCreateDir() {
	dir := images.PreviewPath
	if !tool.Exists(dir) {
		tool.CreateDir(dir)
	}
	dir = images.SmallPreviewPath
	if !tool.Exists(dir) {
		tool.CreateDir(dir)
	}
}
