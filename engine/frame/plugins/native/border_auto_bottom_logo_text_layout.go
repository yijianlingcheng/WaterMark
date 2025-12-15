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
func (b *autoBottomLogoTextLayoutBorder) initLayoutValue(fm *photoFrame) pkg.EError {
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
func (b *autoBottomLogoTextLayoutBorder) getTextLayoutLogoCommonData(fm *photoFrame) autoBottomLogoShowInfo {
	// 计算logo高度
	fm.opts.Params.LogoHeight = fm.opts.Params.LogoRatio * fm.opts.Params.MainMarginBottom / 100
	// 计算logo 100%展示对应的原始宽高
	r1 := layout.GetLogoXAndYByNameAndHeight(
		layout.GetLogoNameByMake(fm.opts.getMakeFromExif()),
		fm.opts.Params.MainMarginBottom,
	)
	// 容错
	if r1["width"] == 0 {
		r1["width"] = fm.opts.Params.MainMarginBottom
	}
	// 计算logo实际展示对应的宽高
	r2 := layout.GetLogoXAndYByNameAndHeight(
		layout.GetLogoNameByMake(fm.opts.getMakeFromExif()),
		fm.opts.Params.LogoHeight,
	)
	fm.opts.Params.LogoWidth = r2["width"]
	// 容错
	if r2["width"] == 0 {
		fm.opts.Params.LogoWidth = fm.opts.Params.LogoHeight
	}

	return autoBottomLogoShowInfo{
		r1:         r1,
		r2:         r2,
		diffWidth:  r1["width"] - r2["width"],
		diffHeight: r1["height"] - r2["height"],
	}
}

// 计算logo.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutLogo(fm *photoFrame) {
	logoShowInfo := b.getTextLayoutLogoCommonData(fm)
	// logo区域没有空隙,不需要计算margin边距,或者数值过小,无需处理
	if logoShowInfo.diffHeight < 2 {
		return
	}

	fm.opts.Params.LogoMarginTop = logoShowInfo.diffHeight / 2
	fm.opts.Params.LogoMarginBottom = logoShowInfo.diffHeight - logoShowInfo.diffHeight/2

	// 数值过小,无需处理
	if logoShowInfo.diffWidth < 2 {
		return
	}
	fm.opts.Params.LogoMarginLeft = logoShowInfo.diffWidth / 2
	fm.opts.Params.LogoMarginRight = logoShowInfo.diffWidth - logoShowInfo.diffWidth/2
}

// 计算右布局下logo的展示位置.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutLogoWithRight(fm *photoFrame) {
	// 有分割线的情况下,使用分割线坐标计算logo展示位置
	if b.HasSeparator {
		fm.opts.Params.LogoMarginLeft = fm.opts.Params.SeparatorMarginLeft -
			fm.opts.Params.LogoMarginRight - fm.opts.Params.LogoWidth

		return
	}
	// 没有分割线的情况下,使用右边文字坐标计算logo展示位置
	textMarginLeft := min(fm.opts.Params.TextTwoMarginLeft, fm.opts.Params.TextFourMarginLeft)
	fm.opts.Params.LogoMarginLeft = textMarginLeft -
		fm.opts.Params.LogoMarginRight - fm.opts.Params.LogoWidth
}

// 计算分割线.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutSeparator(fm *photoFrame) {
	// 不存在分隔符
	if !b.HasSeparator {
		return
	}
	// 初始颜色
	if fm.opts.Params.SeparatorColor == "" {
		fm.opts.Params.SeparatorColor = SEPARATOR_COLOR
	}
	// 宽度
	if fm.opts.Params.SeparatorWidth == 0 {
		fm.opts.Params.SeparatorWidth = fm.opts.Params.LogoWidth / 40
	}
	// 高度
	if fm.opts.Params.SeparatorHeight == 0 {
		fm.opts.Params.SeparatorHeight = fm.opts.Params.LogoHeight + fm.opts.Params.LogoMarginTop/2
	}
	// 上边距
	if fm.opts.Params.SeparatorMarginTop == 0 {
		fm.opts.Params.SeparatorMarginTop = (fm.opts.Params.MainMarginBottom - fm.opts.Params.SeparatorHeight) / 2
	}
	// 下边距
	if fm.opts.Params.SeparatorMarginBottom == 0 {
		fm.opts.Params.SeparatorMarginBottom = fm.opts.Params.MainMarginBottom - fm.opts.Params.SeparatorMarginTop
	}

	logoShowWidth := fm.opts.Params.LogoMarginLeft + fm.opts.Params.LogoWidth + fm.opts.Params.LogoMarginRight

	// 左边距
	if fm.opts.Params.SeparatorMarginLeft == 0 {
		fm.opts.Params.SeparatorMarginLeft = fm.opts.Params.LogoWidth/5 + logoShowWidth
	}
	// 右边距
	if fm.opts.Params.SeparatorMarginRight == 0 {
		fm.opts.Params.SeparatorMarginRight = fm.opts.Params.LogoWidth / 5
	}
}

// 继续计算右布局下分割线的位置.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutSeparatorWithRight(fm *photoFrame) {
	// 不存在分隔符
	if !b.HasSeparator {
		return
	}

	textMarginLeft := min(fm.opts.Params.TextTwoMarginLeft, fm.opts.Params.TextFourMarginLeft)

	// 右边距
	fm.opts.Params.SeparatorMarginRight = fm.opts.Params.LogoWidth / 5
	// 左边距
	fm.opts.Params.SeparatorMarginLeft = textMarginLeft - fm.opts.Params.SeparatorMarginRight -
		fm.opts.Params.SeparatorWidth
}

// 计算并设置文字.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutText(fm *photoFrame) {
	// 自适应字体大小
	fontSize := fm.opts.Params.MainMarginBottom * fm.opts.Params.TextRatio / 100

	b.setTextLayoutTextFontSize(fm, fontSize)
	b.setTextLayoutTextMarginTop(fm)
	b.setTextLayoutTextMarginLeft(fm)
}

// 右布局下,计算并设置文字.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutTextMarginLeftWithRight(fm *photoFrame) {
	// 重新设置左边文字左边距
	fm.opts.Params.TextOneMarginLeft = fm.opts.Params.LogoMarginLeft
	fm.opts.Params.TextThreeMarginLeft = fm.opts.Params.LogoMarginLeft
}

// 设置字体size.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutTextFontSize(fm *photoFrame, fontSize int) {
	// 计算字体size

	imageX := fm.opts.getSourceImageX()
	textOneContent := changeText2ExifContent(fm.opts.getExif(), fm.opts.Params.TextOneContent)
	textTwoContent := changeText2ExifContent(fm.opts.getExif(), fm.opts.Params.TextTwoContent)
	textThreeContent := changeText2ExifContent(fm.opts.getExif(), fm.opts.Params.TextThreeContent)
	textFourContent := changeText2ExifContent(fm.opts.getExif(), fm.opts.Params.TextFourContent)

	textContent := textOneContent
	if len(textOneContent) < len(textTwoContent) {
		textContent = textTwoContent
	}
	if len(textThreeContent) < len(textFourContent) {
		textContent += textFourContent
	} else {
		textContent += textThreeContent
	}
	textFileFont := internal.GetFontFilePath(fm.opts.Params.TextOneFontFile)

	// 需要先根据图片尺寸计算出一个最大的fontSize,用于防止文字重叠
	leftShowWidth := fm.opts.Params.LogoMarginLeft + fm.opts.Params.LogoWidth + fm.opts.Params.LogoMarginRight
	textContentMaxFontSize := getTextContentMaxSize(
		imageX-leftShowWidth*3,
		textFileFont,
		textContent,
	)

	// 从外部字体,计算得出的最大展示字体选择一个最小的作为实际字体使用
	maxFontSize := min(fontSize, textContentMaxFontSize)

	if fm.opts.Params.TextOneFontSize == 0 {
		fm.opts.Params.TextOneFontSize = maxFontSize
	}
	if fm.opts.Params.TextTwoFontSize == 0 {
		fm.opts.Params.TextTwoFontSize = maxFontSize
	}
	if fm.opts.Params.TextThreeFontSize == 0 {
		fm.opts.Params.TextThreeFontSize = maxFontSize
	}
	if fm.opts.Params.TextFourFontSize == 0 {
		fm.opts.Params.TextFourFontSize = maxFontSize
	}
}

// 设置文字上边距.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutTextMarginTop(fm *photoFrame) {
	marginTop := fm.opts.Params.LogoMarginTop
	// 兼容logo全展示情况
	if marginTop == 0 {
		marginTop = fm.opts.Params.LogoHeight / 5
		if fm.opts.Params.TextOneMarginTop == 0 {
			fm.opts.Params.TextOneMarginTop = marginTop*2 - fm.opts.Params.TextOneFontSize*72/96
		}
		if fm.opts.Params.TextTwoMarginTop == 0 {
			fm.opts.Params.TextTwoMarginTop = marginTop*2 - fm.opts.Params.TextTwoFontSize*72/96
		}

		// 第二行文字上边距
		if fm.opts.Params.TextThreeMarginTop == 0 {
			fm.opts.Params.TextThreeMarginTop = marginTop * 3
		}
		if fm.opts.Params.TextFourMarginTop == 0 {
			fm.opts.Params.TextFourMarginTop = marginTop * 3
		}

		return
	}
	// 第一行文字上边距
	if fm.opts.Params.TextOneMarginTop == 0 {
		fm.opts.Params.TextOneMarginTop = marginTop
	}
	if fm.opts.Params.TextTwoMarginTop == 0 {
		fm.opts.Params.TextTwoMarginTop = marginTop
	}

	// 第二行文字上边距
	if fm.opts.Params.TextThreeMarginTop == 0 {
		fm.opts.Params.TextThreeMarginTop = marginTop + fm.opts.Params.LogoHeight - fm.opts.Params.TextThreeFontSize*72/96*2
	}
	if fm.opts.Params.TextFourMarginTop == 0 {
		fm.opts.Params.TextFourMarginTop = marginTop + fm.opts.Params.LogoHeight - fm.opts.Params.TextFourFontSize*72/96*2
	}
}

// 设置文字左边距.
func (b *autoBottomLogoTextLayoutBorder) setTextLayoutTextMarginLeft(fm *photoFrame) {
	// 计算logo展示区域的完整宽度
	leftShowWidth := fm.opts.Params.LogoMarginLeft + fm.opts.Params.LogoWidth + fm.opts.Params.LogoMarginRight
	if b.IsRight {
		leftShowWidth = fm.opts.Params.LogoMarginLeft
	}
	rightShowWidth := leftShowWidth
	if b.HasSeparator {
		leftShowWidth = fm.opts.Params.SeparatorMarginLeft +
			fm.opts.Params.SeparatorWidth +
			fm.opts.Params.SeparatorMarginRight +
			fm.opts.Params.LogoMarginRight
	}

	if fm.opts.Params.TextOneMarginLeft == 0 {
		fm.opts.Params.TextOneMarginLeft = leftShowWidth
	}
	if fm.opts.Params.TextThreeMarginLeft == 0 {
		fm.opts.Params.TextThreeMarginLeft = leftShowWidth
	}
	if fm.opts.Params.TextTwoMarginLeft > 0 && fm.opts.Params.TextFourMarginLeft > 0 {
		return
	}
	// 计算右边文字内容的宽度
	marginLeft := b.getRightTextMaxWidth(fm) + rightShowWidth

	// 竖构图照片
	if fm.opts.isVerticalImage() {
		rightShowWidth -= fm.opts.Params.LogoMarginLeft
		marginLeft -= rightShowWidth
	}
	imageWidth := fm.opts.getSourceImageX()
	if fm.opts.Params.TextTwoMarginLeft == 0 {
		fm.opts.Params.TextTwoMarginLeft = imageWidth - marginLeft
	}
	if fm.opts.Params.TextFourMarginLeft == 0 {
		fm.opts.Params.TextFourMarginLeft = imageWidth - marginLeft
	}
}

// 计算右边文字内容的宽度.
func (b *autoBottomLogoTextLayoutBorder) getRightTextMaxWidth(fm *photoFrame) int {
	twoTextContent := changeText2ExifContent(fm.opts.getExif(), fm.opts.Params.TextTwoContent)
	twoTextContentWidth, _ := getTextContentXAndY(
		fm.opts.Params.TextTwoFontSize,
		internal.GetFontFilePath(fm.opts.Params.TextTwoFontFile),
		twoTextContent,
	)

	fourTextContent := changeText2ExifContent(fm.opts.getExif(), fm.opts.Params.TextFourContent)
	fourTextContentWidth, _ := getTextContentXAndY(
		fm.opts.Params.TextFourFontSize,
		internal.GetFontFilePath(fm.opts.Params.TextFourFontFile),
		fourTextContent,
	)

	return max(twoTextContentWidth, fourTextContentWidth)
}

// 画边框.
func (b *autoBottomLogoTextLayoutBorder) drawBorder(fm *photoFrame) pkg.EError {
	// 画logo
	b.drawLogo(fm)

	// 画水印文字
	b.drawWords(fm)

	// 画分隔符
	b.drawSeparator(fm)

	return pkg.NoError
}
