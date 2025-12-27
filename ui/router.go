package ui

import (
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

// 这是APP菜单与实际程序页面的关联页面
// 菜单中的按钮被点击,就会调用对应的函数打开对应的页面

// 打开照片-查看-页面.
func openPhotoView(app *App) {
	rt.WindowExecJS(app.ctx, "window.location.href = \"/V/index.html\"")
}

// 打开边框-选择照片-页面.
func openFrameView(app *App) {
	rt.WindowExecJS(app.ctx, "window.location.href = \"/F/frameView.html\"")
}

// 打开边框-选择文件夹-页面.
func openFrameDSView(app *App) {
	rt.WindowExecJS(app.ctx, "window.location.href = \"/F/frameDSView.html\"")
}

// 打开帮助-查看-页面.
func openHelpView(app *App) {
	rt.WindowExecJS(app.ctx, "window.location.href = \"/A/helpView.html\"")
}

// 打开关于-版本-页面.
func openAboutVersionView(app *App) {
	rt.WindowExecJS(app.ctx, "window.location.href = \"/A/aboutVersionView.html\"")
}

// 打开关于-贡献代码-页面.
func openAboutCodeView(app *App) {
	rt.WindowExecJS(app.ctx, "window.location.href = \"/A/aboutCodeView.html\"")
}
