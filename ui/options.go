package ui

import (
	"embed"

	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

// 获取APP启动时需要的参数.
func getAppOptions(app *App, assets embed.FS, icon []byte) *options.App {
	return &options.App{
		Title:             getAppTitle(),
		Width:             1024,
		Height:            768,
		MinWidth:          1024,
		MinHeight:         768,
		DisableResize:     true, // 不允许调整窗体大小
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		Assets:            assets,
		Menu:              regiestAppMenus(app),
		Logger:            nil,
		LogLevel:          1,
		OnStartup:         app.Startup,
		OnDomReady:        app.DomReady,
		OnBeforeClose:     app.BeforeClose,
		OnShutdown:        app.Shutdown,
		WindowStartState:  options.Normal, // 启动时默认最大化
		Bind:              []any{app},
		Windows: &windows.Options{ // Windows platform specific options
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			WebviewUserDataPath:  "",
		},
		Mac: &mac.Options{ // Mac platform specific options
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About:                &mac.AboutInfo{Title: getAppTitle(), Message: "", Icon: icon},
		},
	}
}

// 获取APP 标题.
func getAppTitle() string {
	return viper.GetString("app.name")
}
