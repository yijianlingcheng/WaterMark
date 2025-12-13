package native

import (
	"WaterMark/layout"
	"WaterMark/message"
	"WaterMark/pkg"
)

// 模板初始化.
func InitAllCachaAndTools() pkg.EError {
	// 画笔字体初始化
	message.SendInfoMsg("加载字体文件")
	err := textBrushInitFontFileToCache()
	if pkg.HasError(err) {
		return err
	}
	message.SendInfoMsg("加载字体文件完成")

	// 相机logo初始化
	message.SendInfoMsg("加载相机logo")
	err = layout.LogosImagesInit()
	if pkg.HasError(err) {
		return err
	}
	message.SendInfoMsg("加载相机logo完成")

	return pkg.NoError
}
