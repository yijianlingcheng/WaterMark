package images

import (
	"image/color"
)

type BorderTemplate struct {

	// IsRound  边框是否拥有圆角
	IsRound bool

	// Radius 边框半径
	Radius int

	// RightWidth 右部边框宽
	RightWidth int

	// LeftWidth 左部边框宽
	LeftWidth int

	// TopHeight 上部边框高
	TopHeight int

	// BottomHeight 下部边框高
	BottomHeight int

	// Color 边框颜色
	Color color.RGBA

	// Colors 边框颜色,字符串表示,顺序为R,G,B,A
	Colors string
}

// // setDefaultBorderTemplate 默认的边框模板
// //
// //	@param b
// func setDefaultBorderTemplate(b *BorderTemplate) {
// 	b.IsRound = false
// 	b.Radius = 0
// 	b.RightWidth = 100
// 	b.LeftWidth = 100
// 	b.TopHeight = 100
// 	b.BottomHeight = 300

// 	b.Colors = "255,255,255,255"
// 	b.Color = strColor2RGBA(b.Colors)
// }

// // newBorderTemplate 构造一个Logo模板
// //
// //	@param opts 可选参数,map结构
// //	@return *BorderTemplate
// func newBorderTemplate(opts ...map[string]string) *BorderTemplate {
// 	b := BorderTemplate{}
// 	setDefaultBorderTemplate(&b)
// 	if len(opts) > 0 {
// 		val := reflect.ValueOf(&b)
// 		val = val.Elem()
// 		for i, v := range opts[0] {
// 			fieldVa := val.FieldByName(i)
// 			fieldVa.Set(reflect.ValueOf(v))
// 		}
// 	}
// 	return &b
// }

// // setOptions
// //
// //	@param opts
// func (b *BorderTemplate) setOptions(opts map[string]string) {
// 	val := reflect.ValueOf(&b)
// 	val = val.Elem()
// 	for i, v := range opts {
// 		fieldVa := val.FieldByName(i)
// 		fieldVa.Set(reflect.ValueOf(v))
// 	}
// }

// getWidth 获取边框宽度
//
//	@return int
func (b *BorderTemplate) getWidth() int {
	return b.LeftWidth + b.RightWidth
}

// getHeight 获取边框高度
//
//	@return int
func (b *BorderTemplate) getHeight() int {
	return b.TopHeight + b.BottomHeight
}
