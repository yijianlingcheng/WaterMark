package native

import (
	"WaterMark/internal"
	"WaterMark/layout"
	"WaterMark/pkg"
)

// logo 展示宽高信息.
type autoBottomLogoShowInfo struct {
	r1         map[string]int
	r2         map[string]int
	diffWidth  int
	diffHeight int
}

// 初始化自动布局的值.
func (b *autoBottomLogoTextLayoutBorder) initLayoutValue(fm baseFrame) pkg.EError {
	// 自动计算边框
	b.setTextLayoutBorder(fm)
	// 自动计算logo
	b.setTextLayoutLogo(fm)
	// 计算分隔符
	b.setTextLayoutSeparator(fm)
	// 计算文字
	b.setTextLayoutText(fm)

	// 判断是否右布局
	if !b.IsRight {
		return pkg.NoError
	}
	// 后续继续计算文字数值
	b.setTextLayoutTextMarginLeftWithRight(fm)
	// 后续继续计算分隔符
	b.setTextLayoutSeparatorWithRight(fm)
	// 后续继续logo Left数值
	b.setTextLayoutLogoWithRight(fm)

	return pkg.NoError
}

// 获取logo展示的基础信息.
func (b *autoBottomLogoTextLayoutBorder) getTextLayoutLogoCommonData(fm baseFrame) autoBottomLogoShowInfo {
	options := fm.getOptions()
	// 计算logo高度
	options.Params.LogoHeight = options.Params.LogoRatio * options.Params.MainMarginBottom / 100
	// 计算logo 100%展示对应的原始宽高
	r1 := layout.GetLogoXAndYByNameAndHeight(
		layout.GetLogoNameByMake(options.getMakeFromExif()),
		options.Params.MainMarginBottom,
	)
	// 容错
	if r1["width"] == 0 {
		r1["width"] = options.Params.MainMarginBottom
	}
	// 计算logo实际展示对应的宽高
	r2 := layout.GetLogoXAndYByNameAndHeight(
		layout.GetLogoNameByMake(options.getMakeFromExif()),
		options.Params.LogoHeight,
	)
	options.Params.LogoWidth = r2["width"]
	// 容错
	if r2["width"] == 0 {
		options.Params.LogoWidth = options.Params.LogoHeight
	}

	return autoBottomLogoShowInfo{
		r1:         r1,
		r2:         r2,
		diffWidth:  r1["width"] - r2["width"],
		diffHeight: r1["height"] - r2["height"],
	}
}

// 计算logo.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutLogo(fm baseFrame) {
	options := fm.getOptions()

	logoShowInfo := b.getTextLayoutLogoCommonData(fm)
	// logo区域没有空隙,不需要计算margin边距,或者数值过小,无需处理
	if logoShowInfo.diffHeight < 2 {
		return
	}

	options.Params.LogoMarginTop = logoShowInfo.diffHeight / 2
	options.Params.LogoMarginBottom = logoShowInfo.diffHeight - logoShowInfo.diffHeight/2

	// 数值过小,无需处理
	if logoShowInfo.diffWidth < 2 {
		return
	}
	options.Params.LogoMarginLeft = logoShowInfo.diffWidth / 2
	options.Params.LogoMarginRight = logoShowInfo.diffWidth - logoShowInfo.diffWidth/2
}

// 计算右布局下logo的展示位置.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutLogoWithRight(fm baseFrame) {
	options := fm.getOptions()
	// 有分割线的情况下,使用分割线坐标计算logo展示位置
	if b.HasSeparator {
		options.Params.LogoMarginLeft = options.Params.SeparatorMarginLeft -
			options.Params.LogoMarginRight - options.Params.LogoWidth

		return
	}
	// 没有分割线的情况下,使用右边文字坐标计算logo展示位置
	textMarginLeft := min(options.Params.TextTwoMarginLeft, options.Params.TextFourMarginLeft)
	options.Params.LogoMarginLeft = textMarginLeft -
		options.Params.LogoMarginRight - options.Params.LogoWidth
}

// 计算分割线.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutSeparator(fm baseFrame) {
	options := fm.getOptions()
	// 不存在分隔符
	if !b.HasSeparator {
		return
	}
	// 初始颜色
	if options.Params.SeparatorColor == "" {
		options.Params.SeparatorColor = SEPARATOR_COLOR
	}
	// 宽度
	if options.Params.SeparatorWidth == 0 {
		options.Params.SeparatorWidth = options.Params.LogoWidth / 40
	}
	// 高度
	if options.Params.SeparatorHeight == 0 {
		options.Params.SeparatorHeight = options.Params.LogoHeight + options.Params.LogoMarginTop/2
	}
	// 上边距
	if options.Params.SeparatorMarginTop == 0 {
		options.Params.SeparatorMarginTop = (options.Params.MainMarginBottom - options.Params.SeparatorHeight) / 2
	}
	// 下边距
	if options.Params.SeparatorMarginBottom == 0 {
		options.Params.SeparatorMarginBottom = options.Params.MainMarginBottom - options.Params.SeparatorMarginTop
	}

	logoShowWidth := options.Params.LogoMarginLeft + options.Params.LogoWidth + options.Params.LogoMarginRight

	// 左边距
	if options.Params.SeparatorMarginLeft == 0 {
		options.Params.SeparatorMarginLeft = options.Params.LogoWidth/5 + logoShowWidth
	}
	// 右边距
	if options.Params.SeparatorMarginRight == 0 {
		options.Params.SeparatorMarginRight = options.Params.LogoWidth / 5
	}
}

// 继续计算右布局下分割线的位置.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutSeparatorWithRight(fm baseFrame) {
	options := fm.getOptions()
	// 不存在分隔符
	if !b.HasSeparator {
		return
	}

	textMarginLeft := min(options.Params.TextTwoMarginLeft, options.Params.TextFourMarginLeft)

	// 右边距
	options.Params.SeparatorMarginRight = options.Params.LogoWidth / 5
	// 左边距
	options.Params.SeparatorMarginLeft = textMarginLeft - options.Params.SeparatorMarginRight -
		options.Params.SeparatorWidth
}

// 计算并设置文字.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutText(fm baseFrame) {
	options := fm.getOptions()

	// 自适应字体大小
	fontSize := options.Params.MainMarginBottom * options.Params.TextRatio / 100

	b.setTextLayoutTextFontSize(fm, fontSize)
	b.setTextLayoutTextMarginTop(fm)
	b.setTextLayoutTextMarginLeft(fm)
}

// 右布局下,计算并设置文字.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutTextMarginLeftWithRight(fm baseFrame) {
	options := fm.getOptions()
	// 重新设置左边文字左边距
	options.Params.TextOneMarginLeft = options.Params.LogoMarginLeft
	options.Params.TextThreeMarginLeft = options.Params.LogoMarginLeft
}

// 设置字体size.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutTextFontSize(fm baseFrame, fontSize int) {
	options := fm.getOptions()
	// 计算字体size

	imageX := options.getSourceImageX()
	textOneContent := changeText2ExifContent(options.getExif(), options.Params.TextOneContent)
	textTwoContent := changeText2ExifContent(options.getExif(), options.Params.TextTwoContent)
	textThreeContent := changeText2ExifContent(options.getExif(), options.Params.TextThreeContent)
	textFourContent := changeText2ExifContent(options.getExif(), options.Params.TextFourContent)
	textContent := textOneContent
	if len(textOneContent) < len(textTwoContent) {
		textContent = textTwoContent
	}
	if len(textThreeContent) < len(textFourContent) {
		textContent += textFourContent
	} else {
		textContent += textThreeContent
	}
	textFileFont := internal.GetFontFilePath(options.Params.TextOneFontFile)

	// 需要先根据图片尺寸计算出一个最大的fontSize,用于防止文字重叠
	leftShowWidth := options.Params.LogoMarginLeft + options.Params.LogoWidth + options.Params.LogoMarginRight
	textContentMaxFontSize := getTextContentMaxSize(
		imageX-leftShowWidth*3,
		textFileFont,
		textContent,
	)

	// 从外部字体,计算得出的最大展示字体选择一个最小的作为实际字体使用
	maxFontSize := min(fontSize, textContentMaxFontSize)

	if options.Params.TextOneFontSize == 0 {
		options.Params.TextOneFontSize = maxFontSize
	}
	if options.Params.TextTwoFontSize == 0 {
		options.Params.TextTwoFontSize = maxFontSize
	}
	if options.Params.TextThreeFontSize == 0 {
		options.Params.TextThreeFontSize = maxFontSize
	}
	if options.Params.TextFourFontSize == 0 {
		options.Params.TextFourFontSize = maxFontSize
	}
}

// 设置文字上边距.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutTextMarginTop(fm baseFrame) {
	options := fm.getOptions()

	marginTop := options.Params.LogoMarginTop
	// 兼容logo全展示情况
	if marginTop == 0 {
		marginTop = options.Params.LogoHeight / 5
		if options.Params.TextOneMarginTop == 0 {
			options.Params.TextOneMarginTop = marginTop*2 - options.Params.TextOneFontSize*72/96
		}
		if options.Params.TextTwoMarginTop == 0 {
			options.Params.TextTwoMarginTop = marginTop*2 - options.Params.TextTwoFontSize*72/96
		}

		// 第二行文字上边距
		if options.Params.TextThreeMarginTop == 0 {
			options.Params.TextThreeMarginTop = marginTop * 3
		}
		if options.Params.TextFourMarginTop == 0 {
			options.Params.TextFourMarginTop = marginTop * 3
		}

		return
	}
	// 第一行文字上边距
	if options.Params.TextOneMarginTop == 0 {
		options.Params.TextOneMarginTop = marginTop
	}
	if options.Params.TextTwoMarginTop == 0 {
		options.Params.TextTwoMarginTop = marginTop
	}

	// 第二行文字上边距
	if options.Params.TextThreeMarginTop == 0 {
		options.Params.TextThreeMarginTop = marginTop + options.Params.LogoHeight - options.Params.TextThreeFontSize*72/96*2
	}
	if options.Params.TextFourMarginTop == 0 {
		options.Params.TextFourMarginTop = marginTop + options.Params.LogoHeight - options.Params.TextFourFontSize*72/96*2
	}
}

// 设置文字左边距.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutTextMarginLeft(fm baseFrame) {
	options := fm.getOptions()

	// 计算logo展示区域的完整宽度
	leftShowWidth := options.Params.LogoMarginLeft + options.Params.LogoWidth + options.Params.LogoMarginRight
	if b.IsRight {
		leftShowWidth = options.Params.LogoMarginLeft
	}
	rightShowWidth := leftShowWidth
	if b.HasSeparator {
		leftShowWidth = options.Params.SeparatorMarginLeft +
			options.Params.SeparatorWidth +
			options.Params.SeparatorMarginRight +
			options.Params.LogoMarginRight
	}

	if options.Params.TextOneMarginLeft == 0 {
		options.Params.TextOneMarginLeft = leftShowWidth
	}
	if options.Params.TextThreeMarginLeft == 0 {
		options.Params.TextThreeMarginLeft = leftShowWidth
	}
	if options.Params.TextTwoMarginLeft > 0 && options.Params.TextFourMarginLeft > 0 {
		return
	}
	// 计算右边文字内容的宽度
	marginLeft := b.getRightTextMaxWidth(fm) + rightShowWidth

	// 竖构图照片
	if options.isVerticalImage() {
		rightShowWidth -= options.Params.LogoMarginLeft
		marginLeft -= rightShowWidth
	}
	imageWidth := options.getSourceImageX()
	if options.Params.TextTwoMarginLeft == 0 {
		options.Params.TextTwoMarginLeft = imageWidth - marginLeft
	}
	if options.Params.TextFourMarginLeft == 0 {
		options.Params.TextFourMarginLeft = imageWidth - marginLeft
	}
}

// 计算右边文字内容的宽度.
func (b *autoBottomLogoTextLayoutBorder) getRightTextMaxWidth(fm baseFrame) int {
	options := fm.getOptions()
	twoTextContent := changeText2ExifContent(options.getExif(), options.Params.TextTwoContent)
	twoTextContentWidth, _ := getTextContentXAndY(
		options.Params.TextTwoFontSize,
		internal.GetFontFilePath(options.Params.TextTwoFontFile),
		twoTextContent,
	)

	fourTextContent := changeText2ExifContent(options.getExif(), options.Params.TextFourContent)
	fourTextContentWidth, _ := getTextContentXAndY(
		options.Params.TextFourFontSize,
		internal.GetFontFilePath(options.Params.TextFourFontFile),
		fourTextContent,
	)

	return max(twoTextContentWidth, fourTextContentWidth)
}

// 画边框.
func (b *autoBottomLogoTextLayoutBorder) drawBorder(fm baseFrame) pkg.EError {
	// 画logo
	b.drawLogo(fm)

	// 画水印文字
	b.drawWords(fm)

	// 画分隔符
	b.drawSeparator(fm)

	return pkg.NoError
}
