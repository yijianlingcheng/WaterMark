package images

import (
	"image/color"
)

// borderTemplate 边框模板
type borderTemplate struct {

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

// newBorderTemplate 构造一个边框模板
//
//	@return *borderTemplate
func newBorderTemplate() *borderTemplate {
	return &borderTemplate{}
}

// WithWidth 指定宽度
//
//	@param width
//	@return *borderTemplate
func (b *borderTemplate) WithWidth(width int) *borderTemplate {
	leftWidth := int(width / 2)
	b.LeftWidth = leftWidth
	b.RightWidth = width - leftWidth
	return b
}

// WithHeight 指定高度
//
//	@param height
//	@return *borderTemplate
func (b *borderTemplate) WithHeight(height int) *borderTemplate {
	topHeight := int(height / 4)
	b.TopHeight = topHeight
	b.BottomHeight = height - topHeight
	return b
}

// WithBottomHeight 只设置底部边框高度
//
//	@param height
//	@return *borderTemplate
func (b *borderTemplate) WithBottomHeight(height int) *borderTemplate {
	b.BottomHeight = height
	return b
}

// WithBoderColor 指定边框颜色颜色
//
//	@param color
//	@return *borderTemplate
func (b *borderTemplate) WithBoderColor(color color.RGBA) *borderTemplate {
	b.Color = color
	return b
}

// getWidth 获取边框宽度
//
//	@return int
func (b *borderTemplate) getWidth() int {
	return b.LeftWidth + b.RightWidth
}

// getHeight 获取边框高度
//
//	@return int
func (b *borderTemplate) getHeight() int {
	return b.TopHeight + b.BottomHeight
}
