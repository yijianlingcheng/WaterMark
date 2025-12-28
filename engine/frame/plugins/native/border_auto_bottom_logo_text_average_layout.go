package native

import (
	"WaterMark/internal"
	"WaterMark/pkg"
)

// 计算布局.
func (b *autoBottomLogoTextAverageLayoutBorder) initLayoutValue(fm baseFrame) pkg.EError {
	// 计算边框布局
	b.setTextLayoutBorder(fm)
	// 计算文字布局与logo布局
	b.setTextLayoutTextAndLogo(fm)

	return pkg.NoError
}

// 计算文字布局与logo布局、分割线布局.
func (b *autoBottomLogoTextAverageLayoutBorder) setTextLayoutTextAndLogo(fm baseFrame) {
	options := fm.getOptions()

	imageX := options.getSourceImageX()

	textOneContent := changeText2ExifContent(options.getExif(), options.Params.TextOneContent)
	textTwoContent := changeText2ExifContent(options.getExif(), options.Params.TextTwoContent)
	textThreeContent := changeText2ExifContent(options.getExif(), options.Params.TextThreeContent)
	// 计算logo
	b.getTextLayoutLogoCommonData(fm)
	// 设置字体大小
	b.setFontSize(fm, textOneContent, textTwoContent, textThreeContent)

	twoTextWidth, twoTextHeight := getTextContentXAndY(
		options.Params.TextTwoFontSize,
		internal.GetFontFilePath(options.Params.TextTwoFontFile),
		textTwoContent,
	)
	// 字体布局
	options.Params.TextTwoMarginLeft = imageX - twoTextWidth - options.Params.TextTwoFontSize
	options.Params.TextTwoMarginTop = (options.Params.MainMarginBottom - options.Params.TextTwoFontSize) / 2

	_, oneTextHeight := getTextContentXAndY(
		options.Params.TextOneFontSize,
		internal.GetFontFilePath(options.Params.TextOneFontFile),
		textOneContent,
	)

	options.Params.TextOneMarginTop = (options.Params.MainMarginBottom-options.Params.TextOneFontSize)/3 +
		(options.Params.TextOneFontSize-oneTextHeight)/2
	options.Params.TextOneMarginLeft = options.Params.TextTwoFontSize
	options.Params.TextThreeMarginTop = (options.Params.MainMarginBottom - options.Params.TextThreeFontSize) / 3 * 2
	options.Params.TextThreeMarginLeft = options.Params.TextTwoFontSize

	// logo布局
	options.Params.LogoMarginTop = options.Params.TextOneMarginTop + (options.Params.TextTwoFontSize - twoTextHeight)
	options.Params.LogoMarginLeft = options.Params.TextTwoMarginLeft -
		options.Params.LogoWidth*2

	// 颜色
	options.Params.SeparatorColor = SEPARATOR_COLOR
	// 宽度
	options.Params.SeparatorWidth = options.Params.LogoWidth / 40
	// 高度
	options.Params.SeparatorHeight = options.Params.LogoHeight
	// 上边距
	options.Params.SeparatorMarginTop = options.Params.LogoMarginTop
	// 左边距
	options.Params.SeparatorMarginLeft = options.Params.LogoMarginLeft +
		options.Params.LogoWidth + options.Params.LogoWidth/2
}

// 设置字体大小.需要先根据图片尺寸计算出一个最大的fontSize,用于防止文字重叠.
func (b *autoBottomLogoTextAverageLayoutBorder) setFontSize(
	fm baseFrame,
	textOneContent, textTwoContent, textThreeContent string,
) {
	options := fm.getOptions()
	imageX := options.getSourceImageX()

	textContent := textOneContent + textTwoContent
	textFontFile := internal.GetFontFilePath(options.Params.TextOneFontFile)
	if len(textOneContent+textTwoContent) < len(textThreeContent+textTwoContent) {
		textContent = textThreeContent + textTwoContent
	}
	// 需要先根据图片尺寸计算出一个最大的fontSize,用于防止文字重叠
	textContentMaxFontSize := getTextContentMaxSize(
		imageX-options.Params.LogoWidth*5/2,
		textFontFile,
		textContent,
	)

	// 字体大小
	options.Params.TextOneFontSize = options.Params.MainMarginBottom * options.Params.TextOneFontSize / 100
	options.Params.TextTwoFontSize = options.Params.MainMarginBottom * options.Params.TextTwoFontSize / 100
	options.Params.TextThreeFontSize = options.Params.MainMarginBottom * options.Params.TextThreeFontSize / 100

	// 重新赋值
	if options.Params.TextOneFontSize > textContentMaxFontSize {
		options.Params.TextOneFontSize = textContentMaxFontSize
	}
	if options.Params.TextTwoFontSize > textContentMaxFontSize {
		options.Params.TextTwoFontSize = textContentMaxFontSize
	}
	if options.Params.TextThreeFontSize > textContentMaxFontSize {
		options.Params.TextThreeFontSize = textContentMaxFontSize
	}
}

// 画边框.
func (b *autoBottomLogoTextAverageLayoutBorder) drawBorder(fm baseFrame) pkg.EError {
	// 画logo
	b.drawLogo(fm)
	// 画水印文字
	b.drawWords(fm)
	// 画分隔符
	b.drawSeparator(fm)

	return pkg.NoError
}
