package images

import (
	"image/color"
)

// wordsTemplate 文字模板,一共四个位置,字段由exif的字段组成,多个使用,进行分割
type wordsTemplate struct {

	// FirstFontMarginTop 水印第一行文字距离原始图片高度
	FirstFontMarginTop int

	// SecondFontMarginTop 水印第二行文字距离原始图片高度
	SecondFontMarginTop int

	// FirstFontSize 水印第一行文字size
	FirstFontSize int

	// SecondFontSize 水印第二行文字size
	SecondFontSize int

	// FirstFontMarginLeft 水印第一行文字距离原始图片左边宽度
	FirstFontMarginLeft int

	// SecondFontMarginLeft 水印第二行文字距离原始图片左边宽度
	SecondFontMarginLeft int

	// FirstFontMarginRight 水印第一行文字距离原始图片右边宽度
	FirstFontMarginRight int

	// SecondFontMarginRight 水印第二行文字距离原始图片右边宽度
	SecondFontMarginRight int

	// FontWidth 水印文字宽度,用于计算缩进
	FontWidth int

	// FirstFontColor 水印第一行文字颜色
	FirstFontColor color.RGBA

	// SecondFontColor 水印第二行文字颜色
	SecondFontColor color.RGBA

	// FirstFontFile 水印第一行文字字体
	FirstFontFile string

	// SecondFontFile 水印第二行文字字体
	SecondFontFile string

	// FirstFontColors 水印第一行文字颜色,字符串表示,顺序为R,G,B,A
	FirstFontColors string

	// SecondFontColors 水印第二行文字颜色,字符串表示,顺序为R,G,B,A
	SecondFontColors string

	// One
	One string

	// Two
	Two string

	// Three
	Three string

	// Four
	Four string
}

// setDefaultWordsTemplate 默认的文字模板
//
//	@param w
func setDefaultWordsTemplate(w *wordsTemplate) {
	w.One = "Model"
	w.Two = "LensModel"
	w.Three = "FocalLength,FNumberStr,ExposureTime,ISOStr"
	w.Four = "CreateDate"
	w.FirstFontColors = "0,0,0,255"
	w.FirstFontColor = StrColor2RGBA(w.FirstFontColors)
	w.FirstFontFile = "./fonts/Alibaba-PuHuiTi-Bold.ttf"
	w.SecondFontColors = "128,128,128,255"
	w.SecondFontColor = StrColor2RGBA(w.SecondFontColors)
	w.SecondFontFile = "./fonts/Alibaba-PuHuiTi-Light.ttf"
}

// newWordsTemplate 构造一个文字模板
//
//	@return *wordsTemplate
func newWordsTemplate() *wordsTemplate {
	w := wordsTemplate{}
	setDefaultWordsTemplate(&w)
	return &w
}

// WithFontSize 字体大小
//
//	@param fontSize
//	@return *wordsTemplate
func (t *wordsTemplate) WithFontSize(fontSize int) *wordsTemplate {
	t.FirstFontSize = fontSize
	t.SecondFontSize = fontSize
	return t
}

// WithMarginRight 右边距
//
//	@param marginRight
//	@return *wordsTemplate
func (t *wordsTemplate) WithMarginRight(marginRight int) *wordsTemplate {
	t.FirstFontMarginRight = marginRight
	t.SecondFontMarginRight = marginRight
	return t
}

// WithMarginLeft 左边距
//
//	@param marginLeft
//	@return *wordsTemplate
func (t *wordsTemplate) WithMarginLeft(marginLeft int) *wordsTemplate {
	t.FirstFontMarginLeft = marginLeft
	t.SecondFontMarginLeft = marginLeft
	return t
}

// WithMarginTop 上边距
//
//	@param marginTop
//	@return *wordsTemplate
func (t *wordsTemplate) WithMarginTop(marginTop int) *wordsTemplate {
	t.FirstFontMarginTop = marginTop
	t.SecondFontMarginTop = marginTop + int(float64(t.FirstFontSize)*1.6)
	return t
}
