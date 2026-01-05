package ui

import (
	"runtime"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"

	rt "github.com/wailsapp/wails/v2/pkg/runtime"

	"WaterMark/internal"
	"WaterMark/pkg"
)

// 这个文件的代码主要是注册整个APP的菜单项
// 整体思路是通过自定义的结构体数组进行循环添加.
type (
	// APP菜单.
	AppMenus struct {
		name  string        // 菜单的名称
		menus []appMenuItem // 菜单下面的子项
	}

	// app菜单子项.
	appMenuItem struct {
		callback    func()
		label       string
		accelerator string
	}
)

// 获取UI菜单.
func getAppMenus(app *App) []AppMenus {
	return []AppMenus{
		{
			name: "设置(S)",
			menus: []appMenuItem{
				// {label: "模板设置", accelerator: "s", callback: func() { openPhotoDSView(app) }},
				{label: "退出程序", accelerator: "q", callback: func() { rt.Quit(app.ctx) }},
			},
		},
		{
			name:  "EXIF查看(V)",
			menus: []appMenuItem{{label: "照片", accelerator: "v", callback: func() { openPhotoView() }}},
		},
		{
			name: "边框生成(F)",
			menus: []appMenuItem{
				{label: "照片", accelerator: "f", callback: func() { openFrameView() }},
				{label: "文件夹", accelerator: "w", callback: func() { openFrameDSView() }},
			},
		},
		{
			name: "关于(A)",
			menus: []appMenuItem{
				{label: "使用说明", accelerator: "a", callback: func() { openHelpView() }},
				{label: "版本", accelerator: "c", callback: func() { openAboutVersionView() }},
				{label: "贡献代码", accelerator: "g", callback: func() { openAboutCodeView() }},
			},
		},
	}
}

// 获取app的菜单.
func regiestAppMenus(app *App) *menu.Menu {
	internal.Log.Debug("开始注册菜单")

	appMenu := menu.NewMenu()
	menus := checkMenusEnabled(getAppMenus(app))
	regiestJsExec(app)
	// 循环创建菜单
	for i, m := range menus {
		if i == 0 {
			// macOS 特殊处理,这个是官方文档中的内容,原因可能是macOS平台的限制.
			// On macOS platform, this must be done right after `NewMenu()`
			// 具体示例请查看这个页面中的代码 https://wails.io/docs/reference/menus/
			if runtime.GOOS == pkg.Darwin {
				appMenu.Append(menu.AppMenu())
			}
		}
		// 添加子菜单
		t := appMenu.AddSubmenu(m.name)
		for _, item := range m.menus {
			// 添加子菜单中的选项
			// 快捷键使用Ctrl+Alt+前缀,Mac上为Ctrl+Option+
			myShortcut := keys.Combo(item.accelerator, keys.CmdOrCtrlKey, keys.OptionOrAltKey)
			t.AddText(item.label, myShortcut, func(_ *menu.CallbackData) {
				item.callback()
			})
			t.AddSeparator()
		}
		if i != 0 {
			continue
		}
		// macOS 特殊处理,这个是官方文档中的内容,原因可能是macOS平台的限制.
		// On macOS platform, EditMenu should be appended to enable Cmd+C, Cmd+V, Cmd+Z... shortcuts
		// 具体示例请查看这个页面中的代码 https://wails.io/docs/reference/menus/
		if runtime.GOOS == pkg.Darwin {
			appMenu.Append(menu.EditMenu())
		}
	}
	internal.Log.Debug("菜单注册完成")

	return appMenu
}

// 此函数是为了检查菜单是否可用,主要的检查项如下
// 1.检查菜单名称是否相同
// 2.检查选项快捷方式是否相同
// 如果检查不通过则直接报错并且停止运行.
func checkMenusEnabled(appMenus []AppMenus) []AppMenus {
	// 名称map
	namem := make(map[string]string, 5)
	// 快捷键map
	am := make(map[string]string, 20)
	ac := 0
	for _, m := range appMenus {
		namem[m.name] = m.name
		for _, i := range m.menus {
			ac++
			am[i.accelerator] = i.accelerator
		}
	}
	if len(appMenus) != len(namem) {
		internal.Log.Panic("菜单名称存在重复,请检查菜单名称")
	}
	if len(am) != ac {
		internal.Log.Panic("菜单快捷方式存在重复,请检查菜单快捷方式")
	}

	return appMenus
}
