package native

import (
	"image"

	"WaterMark/internal"
	"WaterMark/message"
	"WaterMark/pkg"
)

func (b *blurBottomTextCenterLayout) initLayoutValue(fm *photoFrame) pkg.EError {
	// 计算边框布局
	b.setTextLayoutBorder(fm)
	// 计算文字布局与logo布局
	b.setTextLayoutText(fm)

	return pkg.NoError
}

// 获取最大的字体尺寸.
func (b *blurBottomTextCenterLayout) getMaxFontSize(fm *photoFrame, textOneContent, textThreeContent string) int {
	imageX := fm.opts.getSourceImageX()
	// 需要先根据图片尺寸计算出一个最大的fontSize,用于防止文字重叠
	textContent := textOneContent
	textFontFile := internal.GetFontFilePath(fm.opts.Params.TextOneFontFile)
	if len(textOneContent) < len(textThreeContent) {
		textContent = textThreeContent
		textFontFile = internal.GetFontFilePath(fm.opts.Params.TextThreeFontFile)
	}

	textContentMaxFontSize := getTextContentMaxSize(
		imageX,
		textFontFile,
		textContent,
	)

	return min(fm.opts.Params.MainMarginBottom/5, textContentMaxFontSize)
}

// 计算布局信息.
func (b *blurBottomTextCenterLayout) setTextLayoutText(fm *photoFrame) {
	imageX := fm.opts.getSourceImageX()

	textThreeContent := changeText2ExifContent(fm.opts.getExif(), fm.opts.Params.TextThreeContent)
	textOneContent := changeText2ExifContent(fm.opts.getExif(), fm.opts.Params.TextOneContent)

	// 需要先根据图片尺寸计算出一个最大的fontSize,用于防止文字重叠
	showHeight := b.getMaxFontSize(fm, textOneContent, textThreeContent)

	// 第二行文学信息
	fm.opts.Params.TextThreeFontSize = showHeight
	threeTextWidth, _ := getTextContentXAndY(
		fm.opts.Params.TextThreeFontSize,
		internal.GetFontFilePath(fm.opts.Params.TextThreeFontFile),
		textThreeContent,
	)
	fm.opts.Params.TextThreeMarginLeft = (imageX - threeTextWidth) / 2
	fm.opts.Params.TextThreeMarginRight = imageX - fm.opts.Params.TextThreeMarginLeft

	// 将剩余空白部分三等分,用于计算文字的上边距
	diffHeight := (fm.opts.Params.MainMarginBottom - fm.opts.Params.TextThreeFontSize*2) / 3

	fm.opts.Params.TextThreeMarginTop = diffHeight*2 + fm.opts.Params.TextThreeFontSize

	// 第一行文学信息
	fm.opts.Params.TextOneFontSize = fm.opts.Params.TextThreeFontSize
	oneTextWidth, _ := getTextContentXAndY(
		fm.opts.Params.TextOneFontSize,
		internal.GetFontFilePath(fm.opts.Params.TextOneFontFile),
		textOneContent,
	)
	fm.opts.Params.TextOneMarginLeft = (imageX - oneTextWidth) / 2
	fm.opts.Params.TextOneMarginRight = imageX - fm.opts.Params.TextOneMarginLeft
	fm.opts.Params.TextOneMarginTop = diffHeight
}

// 绘制文字.
func (b *blurBottomTextCenterLayout) drawBorder(fm *photoFrame) pkg.EError {
	// 画水印文字
	for _, textMark := range fm.borImage.textLay.list {
		// 默认是左下布局
		margin := textMark.layout.marginLeft + fm.borImage.leftWidth
		// 判断是否是右下布局
		if b.IsRight {
			margin = fm.srcImage.width - textMark.layout.marginRight
		}

		err := textMark.text.drawFontOnRGBA(
			fm.frameDraw,
			image.Pt(margin, textMark.layout.marginTop+fm.srcImage.height+fm.borImage.topHeight),
			changeText2ExifContent(fm.opts.getExif(), textMark.words),
		)
		// 发生错误,发送错误信息
		if pkg.HasError(err) {
			message.SendErrorMsg(err.Error.Error())
		}
	}

	return pkg.NoError
}
