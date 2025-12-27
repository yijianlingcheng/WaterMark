package native

import (
	"image/draw"

	"WaterMark/pkg"
)

// 创建边框,并返回对应的RGBA对象.
func CreateFrameImageRGBA(opts map[string]any) (draw.Image, pkg.EError) {
	isBlur, ok := opts["isBlur"].(bool)
	if ok && isBlur {
		return createBlurFrameImageRGBA(opts)
	}

	return createNormalFrameImageRGBA(opts)
}

// 创建模糊边框,并返回对应的RGBA对象.
func createBlurFrameImageRGBA(opts map[string]any) (draw.Image, pkg.EError) {
	var fm blurPhotoFrame
	// 初始化
	err := fm.initFrame(opts)
	if pkg.HasError(err) {
		return nil, err
	}
	// 绘制
	fm.drawFrame()
	// 合并
	finalImage := fm.drawBlurMerge()
	// 保存
	imageFilePath := fm.getSaveImageFile()
	if imageFilePath != "" {
		saveImageFile(false, imageFilePath, finalImage, 100)
	}
	// 清理
	fm.clean()

	return finalImage, pkg.NoError
}

// 创建普通边框,并返回对应的RGBA对象.
func createNormalFrameImageRGBA(opts map[string]any) (draw.Image, pkg.EError) {
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
		saveImageFile(false, imageFilePath, finalImage, 100)
	}
	// 清理
	fm.clean()

	return finalImage, pkg.NoError
}

// 获取边框信息.
func GetFrameImageBorderInfo(opts map[string]any) (map[string]any, pkg.EError) {
	isBlur, ok := opts["isBlur"].(bool)
	if ok && isBlur {
		return getBlurFrameImageBorderInfo(opts)
	}

	return getNormalFrameImageBorderInfo(opts)
}

// 获取模糊边框信息.
func getBlurFrameImageBorderInfo(opts map[string]any) (map[string]any, pkg.EError) {
	var fm blurPhotoFrame
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

// 获取普通边框信息.
func getNormalFrameImageBorderInfo(opts map[string]any) (map[string]any, pkg.EError) {
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
