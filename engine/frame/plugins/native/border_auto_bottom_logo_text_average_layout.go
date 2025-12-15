package native

import (
	"WaterMark/internal"
	"WaterMark/pkg"
)

// 计算布局.
func (b *autoBottomLogoTextAverageLayoutBorder) initLayoutValue(fm *photoFrame) pkg.EError {
	// 计算边框布局
	b.setTextLayoutBorder(fm)
	// 计算文字布局与logo布局
	b.setTextLayoutTextAndLogo(fm)

	return pkg.NoError
}

// 计算文字布局与logo布局、分割线布局.
func (b *autoBottomLogoTextAverageLayoutBorder) setTextLayoutTextAndLogo(fm *photoFrame) {
	imageX := fm.opts.getSourceImageX()

	textOneContent := changeText2ExifContent(fm.opts.getExif(), fm.opts.Params.TextOneContent)
	textTwoContent := changeText2ExifContent(fm.opts.getExif(), fm.opts.Params.TextTwoContent)
	textThreeContent := changeText2ExifContent(fm.opts.getExif(), fm.opts.Params.TextThreeContent)

	// 计算logo
	b.getTextLayoutLogoCommonData(fm)
	// 设置字体大小
	b.setFontSize(fm, textOneContent, textTwoContent, textThreeContent)

	twoTextWidth, twoTextHeight := getTextContentXAndY(
		fm.opts.Params.TextTwoFontSize,
		internal.GetFontFilePath(fm.opts.Params.TextTwoFontFile),
		textTwoContent,
	)
	// 字体布局
	fm.opts.Params.TextTwoMarginLeft = imageX - twoTextWidth - fm.opts.Params.TextTwoFontSize
	fm.opts.Params.TextTwoMarginTop = (fm.opts.Params.MainMarginBottom - fm.opts.Params.TextTwoFontSize) / 2

	_, oneTextHeight := getTextContentXAndY(
		fm.opts.Params.TextOneFontSize,
		internal.GetFontFilePath(fm.opts.Params.TextOneFontFile),
		textOneContent,
	)

	fm.opts.Params.TextOneMarginTop = (fm.opts.Params.MainMarginBottom-fm.opts.Params.TextOneFontSize)/3 +
		(fm.opts.Params.TextOneFontSize-oneTextHeight)/2

	fm.opts.Params.TextOneMarginLeft = fm.opts.Params.TextTwoFontSize
	fm.opts.Params.TextThreeMarginTop = (fm.opts.Params.MainMarginBottom - fm.opts.Params.TextThreeFontSize) / 3 * 2
	fm.opts.Params.TextThreeMarginLeft = fm.opts.Params.TextTwoFontSize

	// logo布局
	fm.opts.Params.LogoMarginTop = fm.opts.Params.TextOneMarginTop + (fm.opts.Params.TextTwoFontSize - twoTextHeight)
	fm.opts.Params.LogoMarginLeft = fm.opts.Params.TextTwoMarginLeft -
		fm.opts.Params.LogoWidth*2

	// 颜色
	fm.opts.Params.SeparatorColor = SEPARATOR_COLOR
	// 宽度
	fm.opts.Params.SeparatorWidth = fm.opts.Params.LogoWidth / 40
	// 高度
	fm.opts.Params.SeparatorHeight = fm.opts.Params.LogoHeight
	// 上边距
	fm.opts.Params.SeparatorMarginTop = fm.opts.Params.LogoMarginTop
	// 左边距
	fm.opts.Params.SeparatorMarginLeft = fm.opts.Params.LogoMarginLeft +
		fm.opts.Params.LogoWidth + fm.opts.Params.LogoWidth/2
}

// 设置字体大小.需要先根据图片尺寸计算出一个最大的fontSize,用于防止文字重叠.
func (b *autoBottomLogoTextAverageLayoutBorder) setFontSize(
	fm *photoFrame,
	textOneContent, textTwoContent, textThreeContent string,
) {
	imageX := fm.opts.getSourceImageX()

	textContent := textOneContent + textTwoContent
	textFontFile := internal.GetFontFilePath(fm.opts.Params.TextOneFontFile)
	if len(textOneContent+textTwoContent) < len(textThreeContent+textTwoContent) {
		textContent = textThreeContent + textTwoContent
	}
	// 需要先根据图片尺寸计算出一个最大的fontSize,用于防止文字重叠
	textContentMaxFontSize := getTextContentMaxSize(
		imageX-fm.opts.Params.LogoWidth*5/2,
		textFontFile,
		textContent,
	)

	// 字体大小
	fm.opts.Params.TextOneFontSize = fm.opts.Params.MainMarginBottom * fm.opts.Params.TextOneFontSize / 100
	fm.opts.Params.TextTwoFontSize = fm.opts.Params.MainMarginBottom * fm.opts.Params.TextTwoFontSize / 100
	fm.opts.Params.TextThreeFontSize = fm.opts.Params.MainMarginBottom * fm.opts.Params.TextThreeFontSize / 100

	// 重新赋值
	if fm.opts.Params.TextOneFontSize > textContentMaxFontSize {
		fm.opts.Params.TextOneFontSize = textContentMaxFontSize
	}
	if fm.opts.Params.TextTwoFontSize > textContentMaxFontSize {
		fm.opts.Params.TextTwoFontSize = textContentMaxFontSize
	}
	if fm.opts.Params.TextThreeFontSize > textContentMaxFontSize {
		fm.opts.Params.TextThreeFontSize = textContentMaxFontSize
	}
}

// 画边框.
func (b *autoBottomLogoTextAverageLayoutBorder) drawBorder(fm *photoFrame) pkg.EError {
	// 画logo
	b.drawLogo(fm)
	// 画水印文字
	b.drawWords(fm)
	// 画分隔符
	b.drawSeparator(fm)

	return pkg.NoError
}
