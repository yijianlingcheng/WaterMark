package engine

import (
	"sync"

	"WaterMark/engine/frame"
	"WaterMark/message"
	"WaterMark/pkg"
)

var initAllToolsOnce sync.Once

// 初始化全部工具.
func InitAllTools() {
	initAllToolsOnce.Do(func() {
		// 初始化exif工具
		message.SendInfoMsg("初始化exiftool工具")
		err := initExiftool()
		message.SendErrorOrInfo(err, "初始化exiftool工具完成")

		// 初始化插件
		message.SendInfoMsg("插件初始化")
		err = frame.PluginInitAll()
		message.SendErrorOrInfo(err, "插件初始化完成")

		message.SendInfoMsg("加载模板文件对应图片")
		err = frame.LoadOrCreateLayoutImage()
		message.SendErrorOrInfo(err, "加载模板文件对应图片完成")

		if !pkg.HasError(err) {
			message.SendStartSuccess()
		}
	})
}
