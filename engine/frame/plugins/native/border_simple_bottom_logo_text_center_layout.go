package native

import (
	"WaterMark/internal"
	"WaterMark/layout"
	"WaterMark/pkg"
)

// 计算布局.
func (b *simpleBottomLogoTextCenterBorder) initLayoutValue(fm baseFrame) pkg.EError {
	// 计算边框布局
	b.setTextLayoutBorder(fm)
	// 计算文字布局与logo布局
	b.setTextLayoutTextAndLogo(fm)

	return pkg.NoError
}

// 获取最大的字体尺寸.
func (b *simpleBottomLogoTextCenterBorder) getMaxFontSize(
	fm baseFrame,
	textOneContent, textThreeContent string,
) int {
	options := fm.getOptions()

	imageX := options.getSourceImageX()
	// 需要先根据图片尺寸计算出一个最大的fontSize,用于防止文字重叠
	textContent := textOneContent
	textFontFile := internal.GetFontFilePath(options.Params.TextOneFontFile)
	if len(textOneContent) < len(textThreeContent) {
		textContent = textThreeContent
		textFontFile = internal.GetFontFilePath(options.Params.TextThreeFontFile)
	}

	if b.HasLogo {
		logoName := layout.GetLogoNameByMake(options.getMakeFromExif())
		textContentMaxFontSize := getTextContentMaxSizeWithLogo(
			imageX,
			logoName,
			textFontFile,
			textContent,
		)

		return min(options.Params.MainMarginBottom/5, textContentMaxFontSize)
	}

	textContentMaxFontSize := getTextContentMaxSize(
		imageX,
		textFontFile,
		textContent,
	)

	return min(options.Params.MainMarginBottom/5, textContentMaxFontSize)
}

// 计算布局信息.
func (b *simpleBottomLogoTextCenterBorder) setTextLayoutTextAndLogo(fm baseFrame) {
	options := fm.getOptions()
	imageX := options.getSourceImageX()

	textThreeContent := changeText2ExifContent(options.getExif(), options.Params.TextThreeContent)
	textOneContent := changeText2ExifContent(options.getExif(), options.Params.TextOneContent)

	// 需要先根据图片尺寸计算出一个最大的fontSize,用于防止文字重叠
	showHeight := b.getMaxFontSize(fm, textOneContent, textThreeContent)

	// 计算logo 宽高
	logoShowInfo := layout.GetLogoXAndYByNameAndHeight(
		layout.GetLogoNameByMake(options.getMakeFromExif()),
		showHeight,
	)
	// 第二行文学信息
	options.Params.TextThreeFontSize = logoShowInfo["height"]
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

	if !b.HasLogo {
		return
	}

	b.setTextLayoutWithHasLogo(fm, logoShowInfo)
}

// 计算存在logo情况下的文字与logo布局.
func (b *simpleBottomLogoTextCenterBorder) setTextLayoutWithHasLogo(fm baseFrame, logoShowInfo map[string]int) {
	options := fm.getOptions()

	imageX := options.getSourceImageX()

	options.Params.LogoHeight = logoShowInfo["height"]
	options.Params.LogoWidth = logoShowInfo["width"]

	oneTextWidth, oneTextHeight := getTextContentXAndY(
		options.Params.TextOneFontSize,
		internal.GetFontFilePath(options.Params.TextOneFontFile),
		changeText2ExifContent(options.getExif(), options.Params.TextOneContent),
	)

	options.Params.LogoMarginTop = options.Params.TextOneMarginTop + (options.Params.TextOneFontSize-oneTextHeight)/2

	options.Params.TextOneMarginLeft = (imageX - oneTextWidth + 2*options.Params.LogoWidth) / 2
	options.Params.TextOneMarginRight = imageX - options.Params.TextOneMarginLeft

	options.Params.LogoMarginLeft = imageX - oneTextWidth - options.Params.TextOneMarginLeft - options.Params.LogoWidth

	// 颜色
	options.Params.SeparatorColor = SEPARATOR_COLOR
	// 宽度
	options.Params.SeparatorWidth = options.Params.LogoWidth / 40
	// 高度
	options.Params.SeparatorHeight = options.Params.LogoHeight
	// 上边距
	options.Params.SeparatorMarginTop = options.Params.LogoMarginTop
	// 左边距
	options.Params.SeparatorMarginLeft = options.Params.TextOneMarginLeft - options.Params.LogoWidth
}

// 画边框.
func (b *simpleBottomLogoTextCenterBorder) drawBorder(fm baseFrame) pkg.EError {
	// 画logo
	if b.HasLogo {
		b.drawLogo(fm)
	}
	// 画水印文字
	b.drawWords(fm)

	// 画分隔符
	b.drawSeparator(fm)

	return pkg.NoError
}
