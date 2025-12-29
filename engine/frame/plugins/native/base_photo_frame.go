package native

import (
	"image"
	"image/draw"
	"runtime"

	"WaterMark/layout"
	"WaterMark/message"
	"WaterMark/pkg"
)

// 基础边框接口.
type (
	baseFrame interface {
		initSetSize(opts map[string]any) pkg.EError
		initFrame(opts map[string]any) pkg.EError
		drawFrame()
		getFrameSize() map[string]int
		getBorderText() []string
		getSaveImageFile() string
		getLayoutName() string
		getLayoutParams() *layout.FrameLayout
		getPhotoFrame() *basePhotoFrame
		getOptions() *frameOption
		getBorImage() *borderImage
		getSrcImage() *sourceImage
		getBorderDraw() *image.RGBA
		getFrameDraw() draw.Image
		clean()
	}

	// 基础照片边框结构体.
	basePhotoFrame struct {
		frameDraw  draw.Image
		borderDraw *image.RGBA
		srcImage   *sourceImage
		borImage   *borderImage
		finImage   *finalImage
		opts       *frameOption
		isBlur     bool
	}
)

// 获取基础照片边框对象.
func (fm *basePhotoFrame) getPhotoFrame() *basePhotoFrame {
	return fm
}

// 获取布局参数.
func (fm *basePhotoFrame) getLayoutParams() *layout.FrameLayout {
	return &fm.opts.Params
}

// 获取选项.
func (fm *basePhotoFrame) getOptions() *frameOption {
	return fm.opts
}

func (fm *basePhotoFrame) getBorImage() *borderImage {
	return fm.borImage
}

func (fm *basePhotoFrame) getSrcImage() *sourceImage {
	return fm.srcImage
}

func (fm *basePhotoFrame) getBorderDraw() *image.RGBA {
	return fm.borderDraw
}

func (fm *basePhotoFrame) getFrameDraw() draw.Image {
	return fm.frameDraw
}

// 初始化.
func (fm *basePhotoFrame) initSetSize(opts map[string]any) pkg.EError {
	// 判断是否是模糊边框
	isBlur, ok := opts["isBlur"].(bool)
	if ok && isBlur {
		fm.isBlur = true
	}
	fm.opts = newFrameOption(opts)

	sourceImagePath := fm.opts.getSourceImageFile()
	// 初始化各种对象
	fm.srcImage = newSourceImage(sourceImagePath)
	// 判断是否需要旋转图片
	fm.resetSourceImageXAndY()
	// 判断是否需要加载原图
	fm.srcImage.SetImageXAndY(fm.opts.getSourceImageX(), fm.opts.getSourceImageY())

	borImage, borderErr := getBorderImage(fm)
	if pkg.HasError(borderErr) {
		message.SendErrorMsg(borderErr.String())

		return borderErr
	}
	fm.borImage = borImage

	return pkg.NoError
}

// 获取边框尺寸.
func (fm *basePhotoFrame) getFrameSize() map[string]int {
	isBlur := 0
	if fm.isBlur {
		isBlur = 1
	}
	w := fm.opts.getSourceImageX()
	h := fm.opts.getSourceImageY()
	borderRadius := max(w, h) * fm.opts.Params.BorderRadius / 1000

	return map[string]int{
		"borderLeftWidth":    fm.borImage.leftWidth,
		"borderRightWidth":   fm.borImage.rightWidth,
		"borderTopHeight":    fm.borImage.topHeight,
		"borderBottomHeight": fm.borImage.bottomHeight,
		"sourceWidth":        fm.srcImage.width,
		"sourceHeight":       fm.srcImage.height,
		"isBlur":             isBlur,
		"borderRadius":       borderRadius,
	}
}

// 获取边框上展示的文字信息.
func (fm *basePhotoFrame) getBorderText() []string {
	data := make([]string, 0)
	for i := range fm.borImage.textLay.list {
		key := fm.borImage.textLay.list[i].words
		if fm.borImage.textLay.list[i].words == "" {
			continue
		}
		data = append(data, textWordsList[i], key, changeText2ExifContent(fm.opts.getExif(), key))
	}

	return data
}

// 初始化.
func (fm *basePhotoFrame) initFrame(opts map[string]any) pkg.EError {
	fm.opts = newFrameOption(opts)

	sourceImagePath := fm.opts.getSourceImageFile()
	// 初始化各种对象
	fm.srcImage = newSourceImage(sourceImagePath)
	// 判断是否需要旋转图片
	fm.resetSourceImageXAndY()

	borImage, borderErr := getBorderImage(fm)
	if pkg.HasError(borderErr) {
		message.SendErrorMsg(borderErr.String())

		return borderErr
	}
	fm.borImage = borImage

	// 判断是否需要加载原图
	if fm.opts.needSourceImage() {
		sourceImage, loadSourceImageErr := fm.loadSourceImage(sourceImagePath)
		if pkg.HasError(loadSourceImageErr) {
			return loadSourceImageErr
		}
		fm.srcImage.SetImage(sourceImage)
	} else {
		fm.srcImage.SetImageXAndY(fm.opts.getSourceImageX(), fm.opts.getSourceImageY())
	}
	// 初始化最终绘制对象
	fm.finImage = newFinalImage(fm.opts)

	var frame draw.Image
	var frameErr pkg.EError

	frame, frameErr = fm.createDraw()
	if pkg.HasError(frameErr) {
		return frameErr
	}
	fm.frameDraw = frame

	return pkg.NoError
}

// 获取布局名称.
func (fm *basePhotoFrame) getLayoutName() string {
	return fm.opts.Params.Name
}

// 创建模糊模板画布.
func (fm *basePhotoFrame) createDraw() (draw.Image, pkg.EError) {
	return loadImageRGBAWithColor(
		0,
		0,
		fm.finImage.width,
		fm.finImage.height,
		fm.borImage.bgColor,
	)
}

// 旋转图片.
func (fm *basePhotoFrame) resetSourceImageXAndY() {
	orientation := pkg.GetOrientation(pkg.AnyToString(fm.opts.getExif().Fields["Orientation"]))
	// 旋转图片尺寸
	if orientation == 0 {
		return
	}
	if orientation != 90 && orientation != 270 {
		return
	}
	x := fm.opts.getSourceImageX()
	y := fm.opts.getSourceImageY()
	fm.opts.resetSourceImageX(y)
	fm.opts.resetSourceImageY(x)
}

// 加载图片.
func (fm *basePhotoFrame) loadSourceImage(path string) (image.Image, pkg.EError) {
	image, loadErr := pkg.LoadImageWithDecode(path)
	if pkg.HasError(loadErr) {
		return nil, loadErr
	}
	orientation := pkg.GetOrientation(pkg.AnyToString(fm.opts.getExif().Fields["Orientation"]))
	if orientation > 0 {
		image = pkg.ImageRotate(orientation, image)
	}

	return image, pkg.NoError
}

func (fm *basePhotoFrame) drawFrame() {
}

// 获取保存图片地址.
func (fm *basePhotoFrame) getSaveImageFile() string {
	return fm.opts.SaveImageFile
}

// 清理.
//
//nolint:revive
func (fm *basePhotoFrame) clean() {
	fm.opts = nil
	fm.frameDraw = nil
	fm.borderDraw = nil
	fm.srcImage = &sourceImage{}
	fm.borImage = &borderImage{}
	fm.finImage = &finalImage{}

	runtime.GC()
}
