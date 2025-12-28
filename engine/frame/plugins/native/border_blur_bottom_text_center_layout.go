package native

import (
	"image"

	"WaterMark/internal"
	"WaterMark/message"
	"WaterMark/pkg"
)

func (b *blurBottomTextCenterLayout) initLayoutValue(fm baseFrame) pkg.EError {
	// 计算边框布局
	b.setTextLayoutBorder(fm)
	// 计算文字布局与logo布局
	b.setTextLayoutText(fm)

	return pkg.NoError
}

// 获取最大的字体尺寸.
func (b *blurBottomTextCenterLayout) getMaxFontSize(fm baseFrame, textOneContent, textThreeContent string) int {
	options := fm.getOptions()

	imageX := options.getSourceImageX()
	// 需要先根据图片尺寸计算出一个最大的fontSize,用于防止文字重叠
	textContent := textOneContent
	textFontFile := internal.GetFontFilePath(options.Params.TextOneFontFile)
	if len(textOneContent) < len(textThreeContent) {
		textContent = textThreeContent
		textFontFile = internal.GetFontFilePath(options.Params.TextThreeFontFile)
	}

	textContentMaxFontSize := getTextContentMaxSize(
		imageX,
		textFontFile,
		textContent,
	)

	return min(options.Params.MainMarginBottom/5, textContentMaxFontSize)
}

// 计算布局信息.
func (b *blurBottomTextCenterLayout) setTextLayoutText(fm baseFrame) {
	options := fm.getOptions()

	imageX := options.getSourceImageX()

	textThreeContent := changeText2ExifContent(options.getExif(), options.Params.TextThreeContent)
	textOneContent := changeText2ExifContent(options.getExif(), options.Params.TextOneContent)

	// 需要先根据图片尺寸计算出一个最大的fontSize,用于防止文字重叠
	showHeight := b.getMaxFontSize(fm, textOneContent, textThreeContent)

	// 第二行文学信息
	options.Params.TextThreeFontSize = showHeight
	threeTextWidth, _ := getTextContentXAndY(
		options.Params.TextThreeFontSize,
		internal.GetFontFilePath(options.Params.TextThreeFontFile),
		textThreeContent,
	)
	options.Params.TextThreeMarginLeft = (imageX - threeTextWidth) / 2
	options.Params.TextThreeMarginRight = imageX - options.Params.TextThreeMarginLeft

	// 将剩余空白部分三等分,用于计算文字的上边距
	diffHeight := (options.Params.MainMarginBottom - options.Params.TextThreeFontSize*2) / 3

	options.Params.TextThreeMarginTop = diffHeight*2 + options.Params.TextThreeFontSize

	// 第一行文学信息
	options.Params.TextOneFontSize = options.Params.TextThreeFontSize
	oneTextWidth, _ := getTextContentXAndY(
		options.Params.TextOneFontSize,
		internal.GetFontFilePath(options.Params.TextOneFontFile),
		textOneContent,
	)
	options.Params.TextOneMarginLeft = (imageX - oneTextWidth) / 2
	options.Params.TextOneMarginRight = imageX - options.Params.TextOneMarginLeft
	options.Params.TextOneMarginTop = diffHeight
}

// 绘制文字.
func (b *blurBottomTextCenterLayout) drawBorder(fm baseFrame) pkg.EError {
	options := fm.getOptions()
	borImage := fm.getBorImage()
	frameDraw := fm.getFrameDraw()
	// 获取照片高度
	h := options.getSourceImageY()
	// 画水印文字
	for _, textMark := range borImage.textLay.list {
		// 默认是左下布局
		margin := textMark.layout.marginLeft + borImage.leftWidth

		err := textMark.text.drawFontOnRGBA(
			frameDraw,
			image.Pt(margin, textMark.layout.marginTop+h+borImage.topHeight),
			changeText2ExifContent(options.getExif(), textMark.words),
		)
		// 发生错误,发送错误信息
		if pkg.HasError(err) {
			message.SendErrorMsg(err.Error.Error())
		}
	}

	return pkg.NoError
}
