package native

import (
	"image"
	"image/draw"
	"sync"

	"WaterMark/internal"
	"WaterMark/layout"
	"WaterMark/message"
	"WaterMark/pkg"
)

type photoFrame struct {
	frameDraw  *image.RGBA
	borderDraw *image.RGBA
	srcImage   *sourceImage
	borImage   *borderImage
	finImage   *finalImage
	opts       *frameOption
}

func (fm *photoFrame) initSetSize(opts map[string]any) pkg.EError {
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
func (fm *photoFrame) getFrameSize() map[string]int {
	return map[string]int{
		"borderLeftWidth":    fm.borImage.leftWidth,
		"borderRightWidth":   fm.borImage.rightWidth,
		"borderTopHeight":    fm.borImage.topHeight,
		"borderBottomHeight": fm.borImage.bottomHeight,
		"sourceWidth":        fm.srcImage.width,
		"sourceHeight":       fm.srcImage.height,
	}
}

// 获取边框上展示的文字信息.
func (fm *photoFrame) getBorderText() []string {
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
func (fm *photoFrame) initFrame(opts map[string]any) pkg.EError {
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

	var frame *image.RGBA
	var frameErr pkg.EError

	frame, frameErr = fm.createDraw()
	if pkg.HasError(frameErr) {
		return frameErr
	}
	fm.frameDraw = frame

	return pkg.NoError
}

// 创建画布.
func (fm *photoFrame) createDraw() (*image.RGBA, pkg.EError) {
	return loadImageRGBAWithColor(
		0,
		0,
		fm.finImage.width,
		fm.finImage.height,
		fm.borImage.bgColor,
	)
}

// 旋转图片.
func (fm *photoFrame) resetSourceImageXAndY() {
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
func (fm *photoFrame) loadSourceImage(path string) (image.Image, pkg.EError) {
	image, loadErr := internal.CacheLoadImageWithDecode(path)
	if pkg.HasError(loadErr) {
		return nil, loadErr
	}
	orientation := pkg.GetOrientation(pkg.AnyToString(fm.opts.getExif().Fields["Orientation"]))
	if orientation > 0 {
		image = pkg.ImageRotate(orientation, image)
	}

	return image, pkg.NoError
}

// 画主体.
func (fm *photoFrame) drawMainImage(wg *sync.WaitGroup) {
	defer wg.Done()
	// 生成照片主体

	if fm.opts.needSourceImage() {
		draw.Draw(
			fm.frameDraw,
			fm.srcImage.imgDecode.Bounds().Add(image.Point{fm.borImage.leftWidth, fm.borImage.topHeight}),
			fm.srcImage.imgDecode,
			image.Pt(0, 0),
			draw.Over,
		)
	}
}

// 画边框与文字.
func (fm *photoFrame) drawBorderImage(wg *sync.WaitGroup) pkg.EError {
	defer wg.Done()

	// 生成边框对象
	fm.borderDraw = loadImageRGBA(0, 0, fm.srcImage.width, fm.borImage.bottomHeight)
	draw.Draw(
		fm.borderDraw,
		fm.borderDraw.Bounds(),
		&image.Uniform{fm.borImage.bgColor},
		image.Point{0, 0},
		draw.Src,
	)
	// 相机的logo如果没有找到,则使用特定标识的logo进行代替
	logo, logoErr := layout.GetLogoImageByNameAndWidhtAndHeight(
		layout.GetLogoNameByMake(fm.opts.getMakeFromExif()),
		fm.borImage.logoLay.layout.width,
		fm.borImage.logoLay.layout.height,
	)
	if pkg.HasError(logoErr) {
		return logoErr
	}
	fm.borImage.logoLay.item = logo
	simpleBorderFactory := &SimpleBorderFactory{}
	simpleBorderFactory.createBorder(fm.opts.Params.Name).drawBorder(fm)

	return pkg.NoError
}

// 画出照片主体与边框
// 为了性能考虑采用协程组实现.
func (fm *photoFrame) drawFrame() {
	// 协程组
	var wg sync.WaitGroup
	// 2协程
	wg.Add(2)

	// 画主体
	go fm.drawMainImage(&wg)
	// 画边框
	go fm.drawBorderImage(&wg)

	wg.Wait()
}

// 拼接这两张图片
// 将照片与生成好的边框水印图片拼接在一起.
func (fm *photoFrame) drawMerge() *image.RGBA {
	draw.Draw(
		fm.frameDraw,
		image.Rect(
			fm.borImage.leftWidth,
			fm.borImage.topHeight+fm.srcImage.height,
			fm.borImage.leftWidth+fm.srcImage.width+fm.borImage.rightWidth,
			fm.borImage.topHeight+fm.srcImage.height+fm.borImage.bottomHeight,
		),
		fm.borderDraw,
		image.Pt(0, 0),
		draw.Src,
	)

	return fm.frameDraw
}

// 获取保存图片地址.
func (fm *photoFrame) getSaveImageFile() string {
	return fm.opts.SaveImageFile
}

// 清理.
func (fm *photoFrame) clean() {
	fm.opts = nil
	fm.frameDraw = nil
	fm.borderDraw = nil
	fm.srcImage = &sourceImage{}
	fm.borImage = &borderImage{}
	fm.finImage = &finalImage{}
}
