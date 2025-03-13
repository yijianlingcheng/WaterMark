package main

import (
	"WaterMark/src/api"
	"WaterMark/src/config"
	"WaterMark/src/gui"
	"WaterMark/src/images"
	"WaterMark/src/log"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/src
var assets embed.FS

func main() {

	// 读取配置项
	config.Load()

	images.TestTextBrush_DrawFontOnRGBA("./test/竖图.jpg")
	return

	// 启动后端的go服务
	go api.ServerStart()

	// 创建一个APP实例
	app := gui.NewApp()

	// 启动窗口
	err := wails.Run(&options.App{
		Title:             "水印",
		Width:             1270,
		Height:            768,
		MinWidth:          1270,
		MinHeight:         768,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		Assets:            assets,
		Menu:              nil,
		Logger:            nil,
		LogLevel:          logger.DEBUG,
		OnStartup:         app.Startup,
		OnDomReady:        app.DomReady,
		OnBeforeClose:     app.BeforeClose,
		OnShutdown:        app.Shutdown,
		WindowStartState:  options.Normal,
		Bind: []interface{}{
			app,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
		},
	})
	if err != nil {
		log.ErrorLogger.Println(err)
	}
}
