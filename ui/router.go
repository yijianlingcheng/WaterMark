// Package ui 提供应用窗口内路由导航和在前端执行 JavaScript 的辅助函数。
//
// JsExec 封装了应用上下文引用，用于执行与前端交互的 JS 操作。
// jsExec 是 JsExec 的全局单例，便于在程序各处统一使用。
//
// regiestJsExec 创建并设置全局 jsExec 实例，绑定传入的 App 上下文。
//
// openPhotoView 将应用窗口导航到查看照片页面 ("/V/index.html")。
//
// openFrameView 将应用窗口导航到边框-选择照片页面 ("/F/frameView.html")。
//
// openFrameDSView 将应用窗口导航到边框-选择文件夹页面 ("/F/frameDSView.html")。
//
// openHelpView 将应用窗口导航到帮助查看页面 ("/A/helpView.html")。
//
// openAboutVersionView 将应用窗口导航到关于-版本页面 ("/A/aboutVersionView.html")。
//
// openAboutCodeView 将应用窗口导航到关于-贡献代码页面 ("/A/aboutCodeView.html")。
//
// runJsCode 在应用窗口上下文中执行传入的 JavaScript 代码字符串。
package ui

import (
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

// 封装了应用上下文引用，用于执行与前端交互的 JS 操作。
type JsExec struct {
	ctx *App
}

var jsExec *JsExec

// 创建并设置全局 jsExec 实例，绑定传入的 App 上下文。
func regiestJsExec(app *App) {
	jsExec = &JsExec{
		ctx: app,
	}
}

// 这是APP菜单与实际程序页面的关联页面
// 菜单中的按钮被点击,就会调用对应的函数打开对应的页面

// 打开照片-查看-页面.
func openPhotoView() {
	runJsCode("window.location.href = \"/V/index.html\"")
}

// 打开边框-选择照片-页面.
func openFrameView() {
	runJsCode("window.location.href = \"/F/frameView.html\"")
}

// 打开边框-选择文件夹-页面.
func openFrameDSView() {
	runJsCode("window.location.href = \"/F/frameDSView.html\"")
}

// 打开帮助-查看-页面.
func openHelpView() {
	runJsCode("window.location.href = \"/A/helpView.html\"")
}

// 打开关于-版本-页面.
func openAboutVersionView() {
	runJsCode("window.location.href = \"/A/aboutVersionView.html\"")
}

// 打开关于-贡献代码-页面.
func openAboutCodeView() {
	runJsCode("window.location.href = \"/A/aboutCodeView.html\"")
}

// 执行JS代码.
func runJsCode(code string) {
	rt.WindowExecJS(jsExec.ctx.ctx, code)
}
