package images

import (
	"image/color"
)

// WordsTemplate 文字模板,一共四个位置,字段由exif的字段组成,多个使用,进行分割
type WordsTemplate struct {

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

// // setDefaultWordsTemplate 默认的文字模板
// //
// //	@param w
// func setDefaultWordsTemplate(w *WordsTemplate) {
// 	w.One = "Model"
// 	w.Two = "LensModel"
// 	w.Three = "FocalLength,FNumberStr,ExposureTime,ISOStr"
// 	w.Four = "CreateDate"
// 	w.FirstFontMarginTop = 120
// 	w.SecondFontMarginTop = 240
// 	w.FirstFontSize = 60
// 	w.SecondFontSize = 60
// 	w.FirstFontMarginLeft = 300
// 	w.SecondFontMarginLeft = 300
// 	w.FirstFontMarginRight = 0
// 	w.SecondFontMarginRight = 0
// 	w.FontWidth = 35
// 	w.FirstFontColors = "0,0,0,255"
// 	w.FirstFontColor = strColor2RGBA(w.FirstFontColors)
// 	w.FirstFontFile = "./frontend/src/assets/fonts/Alibaba-PuHuiTi-Bold.ttf"
// 	w.SecondFontColors = "128,128,128,255"
// 	w.SecondFontColor = strColor2RGBA(w.SecondFontColors)
// 	w.SecondFontFile = "./frontend/src/assets/fonts/Alibaba-PuHuiTi-Light.ttf"
// }

// // newWordsTemplate 构造一个文字模板
// //
// //	@param opts 可选参数,map结构
// //	@return *WordsTemplate
// func newWordsTemplate(opts ...map[string]string) *WordsTemplate {
// 	w := WordsTemplate{}
// 	setDefaultWordsTemplate(&w)
// 	if len(opts) > 0 {
// 		val := reflect.ValueOf(&w)
// 		val = val.Elem()
// 		for i, v := range opts[0] {
// 			fieldVa := val.FieldByName(i)
// 			fieldVa.Set(reflect.ValueOf(v))
// 		}
// 	}
// 	return &w
// }

// // setOptions
// //
// //	@param opts
// func (w *WordsTemplate) setOptions(opts map[string]string) {
// 	val := reflect.ValueOf(&w)
// 	val = val.Elem()
// 	for i, v := range opts {
// 		fieldVa := val.FieldByName(i)
// 		fieldVa.Set(reflect.ValueOf(v))
// 	}
// }
