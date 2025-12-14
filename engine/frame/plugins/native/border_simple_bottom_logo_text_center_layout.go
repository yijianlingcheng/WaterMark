package native

import (
	"WaterMark/internal"
	"WaterMark/layout"
	"WaterMark/pkg"
)

// 计算布局.
func (b *simpleBottomLogoTextCenterBorder) initLayoutValue(fm *photoFrame) pkg.EError {
	// 计算边框布局
	b.setTextLayoutBorder(fm)
	// 计算文字布局与logo布局
	b.setTextLayoutTextAndLogo(fm)

	return pkg.NoError
}

// 获取最大的字体尺寸.
func (b *simpleBottomLogoTextCenterBorder) getMaxFontSize(fm *photoFrame, textOneContent, textThreeContent string) int {
	imageX := fm.opts.getSourceImageX()
	// 需要先根据图片尺寸计算出一个最大的fontSize,用于防止文字重叠
	textContent := textOneContent
	if len(textOneContent) < len(textThreeContent) {
		textContent = textThreeContent
	}

	textContentMaxFontSize := getTextContentMaxSize(
		imageX,
		textContent,
	)

	return min(fm.opts.Params.MainMarginBottom/5, textContentMaxFontSize)
}

// 计算布局信息.
func (b *simpleBottomLogoTextCenterBorder) setTextLayoutTextAndLogo(fm *photoFrame) {
	imageX := fm.opts.getSourceImageX()

	textThreeContent := changeText2ExifContent(fm.opts.getExif(), fm.opts.Params.TextThreeContent)
	textOneContent := changeText2ExifContent(fm.opts.getExif(), fm.opts.Params.TextOneContent)

	// 需要先根据图片尺寸计算出一个最大的fontSize,用于防止文字重叠
	showHeight := b.getMaxFontSize(fm, textOneContent, textThreeContent)

	// 计算logo 宽高
	logoShowInfo := layout.GetLogoXAndYByNameAndHeight(
		layout.GetLogoNameByMake(fm.opts.getMakeFromExif()),
		showHeight,
	)

	// 第二行文学信息
	fm.opts.Params.TextThreeFontSize = logoShowInfo["height"]
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

	if !b.HasLogo {
		return
	}

	b.setTextLayoutWithHasLogo(fm, logoShowInfo)
}

// 计算存在logo情况下的文字与logo布局.
func (b *simpleBottomLogoTextCenterBorder) setTextLayoutWithHasLogo(fm *photoFrame, logoShowInfo map[string]int) {
	imageX := fm.opts.getSourceImageX()

	fm.opts.Params.LogoHeight = logoShowInfo["height"]
	fm.opts.Params.LogoWidth = logoShowInfo["width"]

	oneTextWidth, oneTextHeight := getTextContentXAndY(
		fm.opts.Params.TextOneFontSize,
		internal.GetFontFilePath(fm.opts.Params.TextOneFontFile),
		changeText2ExifContent(fm.opts.getExif(), fm.opts.Params.TextOneContent),
	)

	fm.opts.Params.LogoMarginTop = fm.opts.Params.TextOneMarginTop + (fm.opts.Params.TextOneFontSize-oneTextHeight)/2

	fm.opts.Params.TextOneMarginLeft = (imageX - oneTextWidth + 2*fm.opts.Params.LogoWidth) / 2
	fm.opts.Params.TextOneMarginRight = imageX - fm.opts.Params.TextOneMarginLeft

	fm.opts.Params.LogoMarginLeft = imageX - oneTextWidth - fm.opts.Params.TextOneMarginLeft - fm.opts.Params.LogoWidth

	// 颜色
	fm.opts.Params.SeparatorColor = SEPARATOR_COLOR
	// 宽度
	fm.opts.Params.SeparatorWidth = fm.opts.Params.LogoWidth / 40
	// 高度
	fm.opts.Params.SeparatorHeight = fm.opts.Params.LogoHeight
	// 上边距
	fm.opts.Params.SeparatorMarginTop = fm.opts.Params.LogoMarginTop
	// 左边距
	fm.opts.Params.SeparatorMarginLeft = fm.opts.Params.TextOneMarginLeft - fm.opts.Params.LogoWidth
}

// 画边框.
func (b *simpleBottomLogoTextCenterBorder) drawBorder(fm *photoFrame) pkg.EError {
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
