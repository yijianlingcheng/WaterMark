package native

import (
	"image"
	"image/draw"

	"WaterMark/message"
	"WaterMark/pkg"
)

// 基础下logo布局边框.
type baseBottomLogoTextLayoutBorder struct {
	IsRight      bool
	HasSeparator bool
}

// 计算照片边距.
func (b *baseBottomLogoTextLayoutBorder) setTextLayoutBorder(fm baseFrame) {
	options := fm.getOptions()

	// 获取照片宽高
	imageWidth := options.getSourceImageX()
	imageHeight := options.getSourceImageY()

	// 计算四周边距
	options.Params.MainMarginLeft = options.Params.MainMarginLeft * imageWidth / 1000
	options.Params.MainMarginRight = options.Params.MainMarginRight * imageWidth / 1000
	options.Params.MainMarginTop = options.Params.MainMarginTop * imageHeight / 1000
	options.Params.MainMarginBottom = options.Params.MainMarginBottom * imageHeight / 1000

	// 竖构图判断
	if !options.isVerticalImage() {
		return
	}
	// 竖构图情况下判断,左右边距是否大于上边距
	if options.Params.MainMarginTop <= options.Params.MainMarginLeft {
		return
	}
	// 对左右边距重新赋值,防止竖构图边框展示比例不协调
	options.Params.MainMarginLeft = options.Params.MainMarginTop
	options.Params.MainMarginRight = options.Params.MainMarginTop
}

// 画分割线.
func (b *baseBottomLogoTextLayoutBorder) drawSeparator(fm baseFrame) {
	borImage := fm.getBorImage()
	srcImage := fm.getSrcImage()
	borderDraw := fm.getBorderDraw()
	if !borImage.sepLay.isExist {
		return
	}
	// 默认是左下布局
	startX := borImage.sepLay.marginLeft
	// 判断是否是右下布局
	if b.IsRight {
		startX = srcImage.width - borImage.sepLay.marginRight
	}
	startY := borImage.sepLay.marginTop
	endX := startX + borImage.sepLay.width
	endY := startY + borImage.sepLay.height
	for i := startX; i <= endX; i++ {
		drawLine(borderDraw, image.Point{i, startY}, image.Point{i, endY}, borImage.sepLay.color)
	}
}

// 画logo.
func (b *baseBottomLogoTextLayoutBorder) drawLogo(fm baseFrame) {
	borImage := fm.getBorImage()
	srcImage := fm.getSrcImage()

	logo := borImage.logoLay.item
	// 默认是左下布局
	startX := borImage.logoLay.layout.marginLeft
	// 判断是否是右下布局
	if b.IsRight {
		startX = srcImage.width - borImage.logoLay.layout.marginRight
	}

	startY := borImage.logoLay.layout.marginTop

	if logo.Ext != ".png" {
		drawBorderLogo(fm.getPhotoFrame(), logo.LogoImage, startX, startY, startX+logo.Width, startY+logo.Height)

		return
	}
	// 针对png图片填色,防止logo绘制到边框时出现黑色背景
	logoImg := image.NewRGBA(logo.LogoImage.Bounds())
	draw.Draw(logoImg, logoImg.Bounds(), &image.Uniform{borImage.bgColor}, image.Point{}, draw.Src)
	draw.Draw(logoImg, logoImg.Bounds(), logo.LogoImage, logo.LogoImage.Bounds().Min, draw.Over)
	drawBorderLogo(fm.getPhotoFrame(), logoImg, startX, startY, startX+logo.Width, startY+logo.Height)
}

// 画文字.
func (b *baseBottomLogoTextLayoutBorder) drawWords(fm baseFrame) {
	borImage := fm.getBorImage()
	srcImage := fm.getSrcImage()
	options := fm.getOptions()
	for _, textMark := range borImage.textLay.list {
		// 默认是左下布局
		margin := textMark.layout.marginLeft
		// 判断是否是右下布局
		if b.IsRight {
			margin = srcImage.width - textMark.layout.marginRight
		}

		err := textMark.text.drawFontOnRGBA(
			fm.getBorderDraw(),
			image.Pt(margin, textMark.layout.marginTop),
			changeText2ExifContent(options.getExif(), textMark.words),
		)
		// 发生错误,发送错误信息
		if pkg.HasError(err) {
			message.SendErrorMsg(err.Error.Error())
		}
	}
}
