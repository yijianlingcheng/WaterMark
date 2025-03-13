package images

import (
	"image/color"
)

type BorderTemplate struct {

	// OnlyBottom 是否只有底部边框
	OnlyBottom bool

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

// newBorderTemplate 构造一个Logo模板
//
//	@return *BorderTemplate
func newBorderTemplate() *BorderTemplate {
	return &BorderTemplate{}
}

// WithWidth 宽度
//
//	@param width
//	@return *BorderTemplate
func (b *BorderTemplate) WithWidth(width int) *BorderTemplate {
	leftWidth := int(width / 2)
	b.LeftWidth = leftWidth
	b.RightWidth = width - leftWidth
	return b
}

// WithHeight 高度
//
//	@param height
//	@return *BorderTemplate
func (b *BorderTemplate) WithHeight(height int) *BorderTemplate {
	topHeight := int(height / 4)
	b.TopHeight = topHeight
	b.BottomHeight = height - topHeight
	return b
}

// WithBottomHeight 只设置底部边框高度
//
//	@param height
//	@return *
func (b *BorderTemplate) WithBottomHeight(height int) *BorderTemplate {
	b.BottomHeight = height
	return b
}

// WithBoderColor 颜色
//
//	@param color
//	@return *BorderTemplate
func (b *BorderTemplate) WithBoderColor(color color.RGBA) *BorderTemplate {
	b.Color = color
	return b
}

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
