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
func (b *baseBottomLogoTextLayoutBorder) setTextLayoutBorder(fm *photoFrame) {
	// 获取照片宽高
	imageWidth := fm.opts.getSourceImageX()
	imageHeight := fm.opts.getSourceImageY()

	// 计算四周边距
	fm.opts.Params.MainMarginLeft = fm.opts.Params.MainMarginLeft * imageWidth / 1000
	fm.opts.Params.MainMarginRight = fm.opts.Params.MainMarginRight * imageWidth / 1000
	fm.opts.Params.MainMarginTop = fm.opts.Params.MainMarginTop * imageHeight / 1000
	fm.opts.Params.MainMarginBottom = fm.opts.Params.MainMarginBottom * imageHeight / 1000

	// 竖构图判断
	if !fm.opts.isVerticalImage() {
		return
	}
	// 竖构图情况下判断,左右边距是否大于上边距
	if fm.opts.Params.MainMarginTop <= fm.opts.Params.MainMarginLeft {
		return
	}
	// 对左右边距重新赋值,防止竖构图边框展示比例不协调
	fm.opts.Params.MainMarginLeft = fm.opts.Params.MainMarginTop
	fm.opts.Params.MainMarginRight = fm.opts.Params.MainMarginTop
}

// 画分割线.
func (b *baseBottomLogoTextLayoutBorder) drawSeparator(fm *photoFrame) {
	if !fm.borImage.sepLay.isExist {
		return
	}
	// 默认是左下布局
	startX := fm.borImage.sepLay.marginLeft
	// 判断是否是右下布局
	if b.IsRight {
		startX = fm.srcImage.width - fm.borImage.sepLay.marginRight
	}
	startY := fm.borImage.sepLay.marginTop
	endX := startX + fm.borImage.sepLay.width
	endY := startY + fm.borImage.sepLay.height
	for i := startX; i <= endX; i++ {
		drawLine(fm.borderDraw, image.Point{i, startY}, image.Point{i, endY}, fm.borImage.sepLay.color)
	}
}

// 画logo.
func (b *baseBottomLogoTextLayoutBorder) drawLogo(fm *photoFrame) {
	logo := fm.borImage.logoLay.item
	// 默认是左下布局
	startX := fm.borImage.logoLay.layout.marginLeft
	// 判断是否是右下布局
	if b.IsRight {
		startX = fm.srcImage.width - fm.borImage.logoLay.layout.marginRight
	}

	startY := fm.borImage.logoLay.layout.marginTop

	if logo.Ext != ".png" {
		drawBorderLogo(fm, logo.LogoImage, startX, startY, startX+logo.Width, startY+logo.Height)

		return
	}
	// 针对png图片填色,防止logo绘制到边框时出现黑色背景
	logoImg := image.NewRGBA(logo.LogoImage.Bounds())
	draw.Draw(logoImg, logoImg.Bounds(), &image.Uniform{fm.borImage.bgColor}, image.Point{}, draw.Src)
	draw.Draw(logoImg, logoImg.Bounds(), logo.LogoImage, logo.LogoImage.Bounds().Min, draw.Over)
	drawBorderLogo(fm, logoImg, startX, startY, startX+logo.Width, startY+logo.Height)
}

// 画文字.
func (b *baseBottomLogoTextLayoutBorder) drawWords(fm *photoFrame) {
	for _, textMark := range fm.borImage.textLay.list {
		// 默认是左下布局
		margin := textMark.layout.marginLeft
		// 判断是否是右下布局
		if b.IsRight {
			margin = fm.srcImage.width - textMark.layout.marginRight
		}

		err := textMark.text.drawFontOnRGBA(
			fm.borderDraw,
			image.Pt(margin, textMark.layout.marginTop),
			changeText2ExifContent(fm.opts.getExif(), textMark.words),
		)
		// 发生错误,发送错误信息
		if pkg.HasError(err) {
			message.SendErrorMsg(err.Error.Error())
		}
	}
}
