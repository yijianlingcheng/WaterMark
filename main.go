package main

import (
	"embed"
	"runtime/debug"

	"WaterMark/api"
	"WaterMark/internal"
	"WaterMark/ui"
)

var (
	//go:embed frontend/src
	assets embed.FS

	//go:embed build/appicon.png
	icon []byte
)

// @title 照片边框工具后端接口
// @version 1.0
// @description 照片边框工具后端接口.
// @termsOfService yijianlingcheng
// @contact.name yijianlingcheng
// @contact.url https://github.com/yijianlingcheng
// @contact.email yijianlingchen@outlook.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:11079
// @BasePath /
// @schemes http https.
func main() {
	debug.SetMemoryLimit(2 * 1024 * 1024 * 1024)
	// 设置运行模式
	internal.SetAppMode(internal.APP_DEV)
	// 初始化配置与资源
	internal.InitAppConfigsAndRes()
	// 启动http server
	go api.ServerStart()
	// 启动
	ui.AppStart(assets, icon)
}
