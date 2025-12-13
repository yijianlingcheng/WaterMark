package native

import (
	"image"
	"image/color"

	"WaterMark/layout"
	"WaterMark/pkg"
)

type (
	// 原始照片.
	sourceImage struct {
		imgDecode image.Image
		path      string
		width     int
		height    int
	}

	// 边框图片.
	borderImage struct {
		textLay      textMarks
		logoLay      logoLayout
		sepLay       separator
		bgColor      color.RGBA
		leftWidth    int
		rightWidth   int
		topHeight    int
		bottomHeight int
	}

	// logo布局.
	logoLayout struct {
		item   *layout.Logo
		layout layoutBox
	}

	// 文字水印.
	textMarks struct {
		list []textMark
	}

	// 文字水印.
	textMark struct {
		text   *textBrush
		words  string
		layout layoutBox
	}

	separator struct {
		width        int
		height       int
		marginTop    int
		marginRight  int
		marginBottom int
		marginLeft   int
		color        color.RGBA
		isExist      bool
	}

	// 布局.
	layoutBox struct {
		width        int
		height       int
		marginTop    int
		marginRight  int
		marginBottom int
		marginLeft   int
	}

	// 生成边框+文字水印后的照片.
	finalImage struct {
		path   string
		width  int
		height int
	}
)

var (
	// 文字位置一.
	textPosOne = "one"
	// 文字位置二.
	textPosTwo = "two"
	// 文字位置三.
	textPosThree = "three"
	// 文字位置四.
	textPosFour = "four"
)

// 原始照片.
func newSourceImage(path string) *sourceImage {
	return &sourceImage{
		path: path,
	}
}

// 边框布局.
func getBorderImage(fm *photoFrame) (*borderImage, pkg.EError) {
	// 根据布局类型获取对应策略
	simpleBorderFactory := &SimpleBorderFactory{}

	// 处理布局对应的初始化
	simpleBorderFactory.createBorder(fm.opts.Params.Name).initLayoutValue(fm)

	return newBorderImage(&fm.opts.Params)
}

// 返回固定布局对象.
func newBorderImage(params *layout.FrameLayout) (*borderImage, pkg.EError) {
	brushErrs := make([]pkg.EError, 4)
	oneTextBrush, oneErr := newTextBrush(
		params.TextOneFontFile, float64(params.TextOneFontSize),
		&image.Uniform{strColor2RGBA(params.TextOneFontColor)},
	)
	brushErrs[0] = oneErr
	twoTextBrush, twoErr := newTextBrush(
		params.TextTwoFontFile, float64(params.TextTwoFontSize),
		&image.Uniform{strColor2RGBA(params.TextTwoFontColor)},
	)
	brushErrs[1] = twoErr
	threeTextBrush, threeErr := newTextBrush(
		params.TextThreeFontFile, float64(params.TextThreeFontSize),
		&image.Uniform{strColor2RGBA(params.TextThreeFontColor)},
	)
	brushErrs[2] = threeErr
	fourTextBrush, fourErr := newTextBrush(
		params.TextFourFontFile, float64(params.TextFourFontSize),
		&image.Uniform{strColor2RGBA(params.TextFourFontColor)},
	)
	brushErrs[3] = fourErr
	for _, brushErr := range brushErrs {
		if pkg.HasError(brushErr) {
			return &borderImage{}, brushErr
		}
	}

	return &borderImage{
			bgColor:      strColor2RGBA(params.BgColor),
			leftWidth:    params.MainMarginLeft,
			rightWidth:   params.MainMarginRight,
			topHeight:    params.MainMarginTop,
			bottomHeight: params.MainMarginBottom,
			textLay: textMarks{
				list: newTextMarks(params, oneTextBrush, twoTextBrush, threeTextBrush, fourTextBrush),
			},
			logoLay: logoLayout{item: &layout.Logo{}, layout: newLogoLayoutBox(params)},
			sepLay:  newSeparator(params),
		},
		pkg.NoError
}

// 分割线.
func newSeparator(params *layout.FrameLayout) separator {
	isExist := params.SeparatorWidth > 0 && params.SeparatorHeight > 0
	if isExist {
		return separator{
			isExist:      isExist,
			width:        params.SeparatorWidth,
			height:       params.SeparatorHeight,
			marginTop:    params.SeparatorMarginTop,
			marginBottom: params.SeparatorMarginBottom,
			marginLeft:   params.SeparatorMarginLeft,
			marginRight:  params.SeparatorMarginRight,
			color:        strColor2RGBA(params.SeparatorColor),
		}
	}

	return separator{
		isExist: isExist,
	}
}

// 水印文字布局.
func newTextMarks(
	params *layout.FrameLayout,
	oneTextBrush, twoTextBrush, threeTextBrush, fourTextBrush *textBrush,
) []textMark {
	return []textMark{
		{
			words:  params.TextOneContent,
			text:   oneTextBrush,
			layout: newTextLayout(params, textPosOne),
		},
		{
			words:  params.TextTwoContent,
			text:   twoTextBrush,
			layout: newTextLayout(params, textPosTwo),
		},
		{
			words:  params.TextThreeContent,
			text:   threeTextBrush,
			layout: newTextLayout(params, textPosThree),
		},
		{
			words:  params.TextFourContent,
			text:   fourTextBrush,
			layout: newTextLayout(params, textPosFour),
		},
	}
}

// 获取文字布局.
func newTextLayout(params *layout.FrameLayout, pos string) layoutBox {
	lay := layoutBox{
		marginTop:    params.TextOneMarginTop,
		marginBottom: params.TextOneMarginBottom,
		marginLeft:   params.TextOneMarginLeft,
		marginRight:  params.TextOneMarginRight,
	}
	switch pos {
	case textPosTwo:
		lay = layoutBox{
			marginTop:    params.TextTwoMarginTop,
			marginBottom: params.TextTwoMarginBottom,
			marginLeft:   params.TextTwoMarginLeft,
			marginRight:  params.TextTwoMarginRight,
		}
	case textPosThree:
		lay = layoutBox{
			marginTop:    params.TextThreeMarginTop,
			marginBottom: params.TextThreeMarginBottom,
			marginLeft:   params.TextThreeMarginLeft,
			marginRight:  params.TextThreeMarginRight,
		}
	case textPosFour:
		lay = layoutBox{
			marginTop:    params.TextFourMarginTop,
			marginBottom: params.TextFourMarginBottom,
			marginLeft:   params.TextFourMarginLeft,
			marginRight:  params.TextFourMarginRight,
		}
	}

	return lay
}

// logo布局.
func newLogoLayoutBox(params *layout.FrameLayout) layoutBox {
	return layoutBox{
		width:        params.LogoWidth,
		height:       params.LogoHeight,
		marginTop:    params.LogoMarginTop,
		marginRight:  params.LogoMarginRight,
		marginBottom: params.LogoMarginBottom,
		marginLeft:   params.LogoMarginLeft,
	}
}

// 边框结果.
func newFinalImage(opts *frameOption) *finalImage {
	return &finalImage{
		width:  opts.Params.MainMarginLeft + opts.Params.MainMarginRight + opts.getSourceImageX(),
		height: opts.Params.MainMarginTop + opts.Params.MainMarginBottom + opts.getSourceImageY(),
		path:   opts.SaveImageFile,
	}
}

// 设置image Decode.
func (src *sourceImage) SetImage(img image.Image) {
	src.imgDecode = img
	src.width = img.Bounds().Dx()
	src.height = img.Bounds().Dy()
}

// 设置图片宽高.
func (src *sourceImage) SetImageXAndY(width, height int) {
	src.width = width
	src.height = height
}
