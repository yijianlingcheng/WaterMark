package native

import (
	"image/draw"

	"WaterMark/pkg"
)

// 创建边框,并返回对应的RGBA对象.
func CreateFrameImageRGBA(opts map[string]any) (draw.Image, pkg.EError) {
	var fm photoFrame
	// 初始化
	err := fm.initFrame(opts)
	if pkg.HasError(err) {
		return nil, err
	}
	// 绘制
	fm.drawFrame()
	// 合并
	finalImage := fm.drawMerge()
	// 保存
	imageFilePath := fm.getSaveImageFile()
	if imageFilePath != "" {
		saveJpgImage(imageFilePath, finalImage, 100)
	}
	// 清理
	fm.clean()

	return finalImage, pkg.NoError
}

// 创建边框,并返回对应的RGBA对象.
func CreateFrameImageRGBABackground(opts map[string]any) pkg.EError {
	var fm photoFrame
	// 初始化
	err := fm.initFrame(opts)
	if pkg.HasError(err) {
		return err
	}
	// 绘制
	fm.drawFrame()
	// 合并
	finalImage := fm.drawMerge()
	// 保存
	imageFilePath := fm.getSaveImageFile()
	if imageFilePath != "" {
		saveJpgImage(imageFilePath, finalImage, 100)
	}
	// 清理
	fm.clean()

	return pkg.NoError
}

// 获取边框信息.
func GetFrameImageBorderInfo(opts map[string]any) (map[string]any, pkg.EError) {
	var fm photoFrame
	err := fm.initSetSize(opts)
	if pkg.HasError(err) {
		return nil, err
	}
	info := make(map[string]any)

	info["size"] = fm.getFrameSize()
	info["text"] = fm.getBorderText()

	fm.clean()

	return info, pkg.NoError
}
