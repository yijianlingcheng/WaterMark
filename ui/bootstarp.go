package ui

import (
	"embed"
	"os"
	"os/signal"
	"syscall"

	"github.com/wailsapp/wails/v2"

	"WaterMark/engine"
	"WaterMark/internal"
	"WaterMark/message"
)

// APP UI 启动.
func AppStart(assets embed.FS, icon []byte) {
	// 是否开启调试模式
	if internal.ISApiDebug() {
		// 初始化部分引擎
		engine.InitAllTools()
		// 等待退出信号
		wait2Exit()
		// 关闭开启的引擎
		engine.QuitAllTools()
		message.Close()
		internal.CleanDir()

		return
	}
	// 启动UI程序
	app := NewApp()
	err := wails.Run(getAppOptions(app, assets, icon))
	if err != nil {
		internal.Log.Fatal(err)
	}
}

// 等待发送信号退出.
func wait2Exit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
