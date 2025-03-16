package images

import (
	"image/color"
)

// separateTamplate 水印图片与文字间的分隔符
type separateTamplate struct {

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
