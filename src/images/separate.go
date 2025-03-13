package images

import (
	"image/color"
)

// SeparateTamplate 水印图片与文字间的分隔符
type SeparateTamplate struct {

	// Exist  logo与相邻的文字是否存在分隔符
	Exist bool

	// Width 分隔符宽度
	Width int

	// Height 分隔符高度
	Height int

	// MarginTop 分隔符top距离
	MarginTop int

	// MarginLeft 分隔符left距离
	MarginLeft int

	// MarginRight 分隔符right距离
	MarginRight int

	// FontColor 分隔符颜色
	FontColor color.RGBA

	// FontColors 分隔符颜色,字符串表示,顺序为R,G,B,A
	FontColors string
}

// // setDefaultSeparateTamplate defaultWordsTemplate 默认的文字模板
// //
// //	@param w
// func setDefaultSeparateTamplate(w *SeparateTamplate) {
// 	w.Exist = false
// 	w.FontColors = "0,0,0,0"
// }

// // newSeparateTamplate 构造一个分隔符模板
// //
// //	@param seps
// //	@return *SeparateTamplate
// func newSeparateTamplate(seps ...map[string]string) *SeparateTamplate {
// 	w := SeparateTamplate{}
// 	setDefaultSeparateTamplate(&w)
// 	if len(seps) > 0 {
// 		val := reflect.ValueOf(&w)
// 		val = val.Elem()
// 		for i, v := range seps[0] {
// 			fieldVa := val.FieldByName(i)
// 			fieldVa.Set(reflect.ValueOf(v))
// 		}
// 		w.Exist = true
// 	}
// 	return &w
// }

// // setOptions
// //
// //	@param opts
// func (s *SeparateTamplate) setOptions(opts map[string]string) {
// 	val := reflect.ValueOf(&s)
// 	val = val.Elem()
// 	for i, v := range opts {
// 		fieldVa := val.FieldByName(i)
// 		fieldVa.Set(reflect.ValueOf(v))
// 	}
// }
